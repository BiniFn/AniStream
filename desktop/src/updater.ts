import { BrowserWindow, ipcMain } from "electron";
import { autoUpdater } from "electron-updater";
import { log, logError, getLogFilePathForUser } from "./logger";

const APP_VERSION = process.env.APP_VERSION!;
const R2_PUBLIC_URL = process.env.R2_PUBLIC_URL!;

export type UpdateStatus =
  | { status: "checking" }
  | { status: "available"; version: string }
  | { status: "not-available" }
  | { status: "downloading"; percent: number }
  | { status: "downloaded"; version: string }
  | { status: "error"; message: string };

let updateStatus: UpdateStatus = { status: "not-available" };
let mainWindow: BrowserWindow | null = null;
let feedUrl: string | null = null;

function sendStatusToRenderer() {
  mainWindow?.webContents.send("update-status", updateStatus);
}

export function setupAutoUpdater(window: BrowserWindow) {
  mainWindow = window;

  // Configure auto-updater
  autoUpdater.autoDownload = false; // We'll trigger download manually when user clicks "Update"
  autoUpdater.autoInstallOnAppQuit = true;

  // Set the feed URL dynamically
  if (!R2_PUBLIC_URL) {
    logError("R2_PUBLIC_URL is not set - auto-updater will not work");
    return;
  }

  // Ensure URL ends with / for electron-updater
  feedUrl = R2_PUBLIC_URL.endsWith("/") ? R2_PUBLIC_URL : `${R2_PUBLIC_URL}/`;
  
  log(`Setting update feed URL: ${feedUrl}`);
  log(`Current app version: ${APP_VERSION}`);
  autoUpdater.setFeedURL({
    provider: "generic",
    url: feedUrl,
  });

  // Register IPC handlers first (needed for both dev and prod)
  ipcMain.handle("get-app-version", () => APP_VERSION);
  ipcMain.handle("get-update-status", () => updateStatus);
  ipcMain.handle("get-log-file-path", () => getLogFilePathForUser());

  // Don't run auto-updater in development
  if (process.env.NODE_ENV === "development") {
    log("Auto-updater disabled in development");
    ipcMain.handle("check-for-updates", async () => {});
    ipcMain.handle("start-update", async () => {});
    ipcMain.handle("quit-and-install", () => {});
    return;
  }

  // Auto-updater event handlers
  autoUpdater.on("checking-for-update", () => {
    log("Checking for updates...");
    updateStatus = { status: "checking" };
    sendStatusToRenderer();
  });

  autoUpdater.on("update-available", (info) => {
    log("Update available:", info.version);
    updateStatus = { status: "available", version: info.version };
    sendStatusToRenderer();
  });

  autoUpdater.on("update-not-available", () => {
    log("No updates available");
    updateStatus = { status: "not-available" };
    sendStatusToRenderer();
  });

  autoUpdater.on("download-progress", (progress) => {
    log(`Download progress: ${progress.percent.toFixed(1)}%`);
    updateStatus = { status: "downloading", percent: progress.percent };
    sendStatusToRenderer();
  });

  autoUpdater.on("update-downloaded", (info) => {
    log("Update downloaded:", info.version);
    updateStatus = { status: "downloaded", version: info.version };
    sendStatusToRenderer();
  });

  autoUpdater.on("error", (error) => {
    logError("Update error:", error);
    logError("Error details:", {
      message: error.message,
      stack: error.stack,
      feedUrl: feedUrl,
      currentVersion: APP_VERSION,
    });
    updateStatus = { status: "error", message: error.message };
    sendStatusToRenderer();
  });

  // IPC handlers for update actions
  ipcMain.handle("check-for-updates", async () => {
    try {
      log("Manually checking for updates...");
      log("Feed URL:", feedUrl);
      log("Current version:", APP_VERSION);
      const result = await autoUpdater.checkForUpdates();
      log("Update check result:", result);
    } catch (error) {
      logError("Failed to check for updates:", error);
      if (error instanceof Error) {
        logError("Error stack:", error.stack);
      }
    }
  });

  ipcMain.handle("start-update", async () => {
    try {
      log("Starting update download...");
      await autoUpdater.downloadUpdate();
    } catch (error) {
      logError("Failed to download update:", error);
    }
  });

  ipcMain.handle("quit-and-install", () => {
    autoUpdater.quitAndInstall();
  });

  // Verify feed URL is accessible (for debugging)
  if (feedUrl) {
    const testUrl = `${feedUrl}latest.yml`;
    log(`Testing update feed accessibility: ${testUrl}`);
    // This is just for logging - electron-updater will handle the actual check
  }

  // Check for updates on startup (after a short delay)
  setTimeout(() => {
    log("Auto-checking for updates on startup...");
    autoUpdater.checkForUpdates().catch((error) => {
      logError("Failed to check for updates on startup:", error);
      if (error instanceof Error) {
        logError("Error details:", error.message, error.stack);
      }
    });
  }, 5000);

  // Check for updates every hour
  setInterval(
    () => {
      log("Periodic update check...");
      autoUpdater.checkForUpdates().catch((error) => {
        logError("Failed to check for updates:", error);
      });
    },
    60 * 60 * 1000,
  );
}
