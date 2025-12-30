import { BrowserWindow, ipcMain } from "electron";
import { autoUpdater } from "electron-updater";

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

function sendStatusToRenderer() {
  mainWindow?.webContents.send("update-status", updateStatus);
}

export function setupAutoUpdater(window: BrowserWindow) {
  mainWindow = window;

  // Configure auto-updater
  autoUpdater.autoDownload = false; // We'll trigger download manually when user clicks "Update"
  autoUpdater.autoInstallOnAppQuit = true;

  // Set the feed URL dynamically
  if (R2_PUBLIC_URL) {
    autoUpdater.setFeedURL({
      provider: "generic",
      url: R2_PUBLIC_URL,
    });
  }

  // Register IPC handlers first (needed for both dev and prod)
  ipcMain.handle("get-app-version", () => APP_VERSION);
  ipcMain.handle("get-update-status", () => updateStatus);

  // Don't run auto-updater in development
  if (process.env.NODE_ENV === "development") {
    console.log("Auto-updater disabled in development");
    ipcMain.handle("check-for-updates", async () => {});
    ipcMain.handle("start-update", async () => {});
    ipcMain.handle("quit-and-install", () => {});
    return;
  }

  // Auto-updater event handlers
  autoUpdater.on("checking-for-update", () => {
    console.log("Checking for updates...");
    updateStatus = { status: "checking" };
    sendStatusToRenderer();
  });

  autoUpdater.on("update-available", (info) => {
    console.log("Update available:", info.version);
    updateStatus = { status: "available", version: info.version };
    sendStatusToRenderer();
  });

  autoUpdater.on("update-not-available", () => {
    console.log("No updates available");
    updateStatus = { status: "not-available" };
    sendStatusToRenderer();
  });

  autoUpdater.on("download-progress", (progress) => {
    console.log(`Download progress: ${progress.percent.toFixed(1)}%`);
    updateStatus = { status: "downloading", percent: progress.percent };
    sendStatusToRenderer();
  });

  autoUpdater.on("update-downloaded", (info) => {
    console.log("Update downloaded:", info.version);
    updateStatus = { status: "downloaded", version: info.version };
    sendStatusToRenderer();
  });

  autoUpdater.on("error", (error) => {
    console.error("Update error:", error);
    updateStatus = { status: "error", message: error.message };
    sendStatusToRenderer();
  });

  // IPC handlers for update actions
  ipcMain.handle("check-for-updates", async () => {
    try {
      await autoUpdater.checkForUpdates();
    } catch (error) {
      console.error("Failed to check for updates:", error);
    }
  });

  ipcMain.handle("start-update", async () => {
    try {
      await autoUpdater.downloadUpdate();
    } catch (error) {
      console.error("Failed to download update:", error);
    }
  });

  ipcMain.handle("quit-and-install", () => {
    autoUpdater.quitAndInstall();
  });

  // Check for updates on startup (after a short delay)
  setTimeout(() => {
    autoUpdater.checkForUpdates().catch((error) => {
      console.error("Failed to check for updates:", error);
    });
  }, 5000);

  // Check for updates every hour
  setInterval(
    () => {
      autoUpdater.checkForUpdates().catch((error) => {
        console.error("Failed to check for updates:", error);
      });
    },
    60 * 60 * 1000,
  );
}
