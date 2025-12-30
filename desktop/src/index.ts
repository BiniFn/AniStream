import { app, BrowserWindow, ipcMain, screen, shell } from "electron";
import path from "node:path";
import { setupAutoUpdater } from "./updater";

// These values are replaced at build time by tsup
const BASE_URL = process.env.ANIWAYS_URL!;

let mainWindow: BrowserWindow | null = null;

const originalLog = console.log;
console.log = (...args: any[]) => {
  if (!args.length) return;
  if (process.env.NODE_ENV === "development") {
    originalLog("[Main][Dev]", ...args);
    return;
  }
  originalLog("[Main]", ...args);
};

const LOADING_HTML = `<!DOCTYPE html>
<html>
<head>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #0a0a0a;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    }
    .loader {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 24px;
    }
    .spinner {
      width: 48px;
      height: 48px;
      border: 3px solid #333;
      border-top-color: #fff;
      border-radius: 50%;
      animation: spin 1s linear infinite;
    }
    .text {
      color: #888;
      font-size: 14px;
    }
    @keyframes spin {
      to { transform: rotate(360deg); }
    }
  </style>
</head>
<body>
  <div class="loader">
    <div class="spinner"></div>
    <div class="text">Loading...</div>
  </div>
</body>
</html>
`;

if (process.env.NODE_ENV === "development") {
  try {
    require("electron-reloader")(module, {
      watchRenderer: false,
    });
  } catch {}
}

function createWindow() {
  const { width: screenWidth, height: screenHeight } =
    screen.getPrimaryDisplay().workAreaSize;

  const isMac = process.platform === "darwin";

  const options: Electron.BrowserWindowConstructorOptions = {
    width: Math.round(screenWidth * 0.75),
    height: Math.round(screenHeight * 0.75),
    minWidth: Math.round(screenWidth * 0.75),
    minHeight: Math.round(screenHeight * 0.75),
    show: false,
    backgroundColor: "#0a0a0a",
    webPreferences: {
      webSecurity: false,
      preload: path.join(__dirname, "preload.js"),
      contextIsolation: true,
    },
  };

  if (isMac) {
    options.titleBarStyle = "hiddenInset";
    options.trafficLightPosition = { x: 20, y: 29 };
  }

  mainWindow = new BrowserWindow(options);

  // Handle fullscreen events
  mainWindow.on("enter-full-screen", () => {
    mainWindow?.webContents.send("fullscreen-change", true);
  });

  mainWindow.on("leave-full-screen", () => {
    mainWindow?.webContents.send("fullscreen-change", false);
  });

  ipcMain.handle("get-fullscreen", () => {
    return mainWindow?.isFullScreen();
  });

  // Security: Only allow navigation to the base URL we loaded
  mainWindow.webContents.on("will-navigate", (event, url) => {
    const isAllowed = url.startsWith(BASE_URL) || url.startsWith("data:");
    if (!isAllowed) {
      event.preventDefault();
      shell.openExternal(url);
    }
  });

  // Open external links in system browser
  mainWindow.webContents.setWindowOpenHandler(({ url }) => {
    if (!url.startsWith(BASE_URL)) {
      shell.openExternal(url);
    }
    return { action: "deny" };
  });

  // Show loading screen first
  mainWindow.loadURL(
    `data:text/html;charset=utf-8,${encodeURIComponent(LOADING_HTML)}`,
  );
  mainWindow.once("ready-to-show", () => {
    mainWindow?.show();
  });

  // Load the actual URL
  console.log(`Loading URL: ${BASE_URL}`);
  mainWindow.loadURL(BASE_URL);
}

app.setName("Aniways");

app.whenReady().then(() => {
  createWindow();

  if (mainWindow) {
    setupAutoUpdater(mainWindow);
  }
});

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});

app.on("activate", () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    // Create window if none are open (macOS)
    createWindow();
  }
});

app.on("before-quit", () => {
  mainWindow = null;
});
