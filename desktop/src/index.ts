import { app, BrowserWindow, ipcMain, screen, shell } from "electron";
import path from "node:path";

const isDev = process.env.NODE_ENV === "development";
const BASE_URL = isDev ? "http://localhost:3000" : "https://aniways.xyz";

// Set app name for macOS menu bar
app.setName("Aniways");

let mainWindow: BrowserWindow | null = null;

const LOADING_HTML = `
<!DOCTYPE html>
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

if (isDev) {
  try {
    require("electron-reloader")(module, {
      watchRenderer: false,
    });
  } catch {}
}

app.whenReady().then(() => {
  const { width: screenWidth, height: screenHeight } =
    screen.getPrimaryDisplay().workAreaSize;

  const isMac = process.platform === "darwin";

  mainWindow = new BrowserWindow({
    width: Math.round(screenWidth * 0.75),
    height: Math.round(screenHeight * 0.75),
    show: false,
    backgroundColor: "#0a0a0a",
    // macOS: custom title bar with traffic lights
    // Windows/Linux: default title bar
    ...(isMac
      ? {
          titleBarStyle: "hiddenInset",
          trafficLightPosition: { x: 20, y: 30 },
        }
      : {
          // Use default frame on Windows/Linux
        }),
    webPreferences: {
      webSecurity: false,
      preload: path.join(__dirname, "preload.js"),
      contextIsolation: true,
    },
  });

  // Handle fullscreen events
  mainWindow.on("enter-full-screen", () => {
    mainWindow.webContents.send("fullscreen-change", true);
  });

  mainWindow.on("leave-full-screen", () => {
    mainWindow.webContents.send("fullscreen-change", false);
  });

  ipcMain.handle("get-fullscreen", () => {
    return mainWindow.isFullScreen();
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
    `data:text/html;charset=utf-8,${encodeURIComponent(LOADING_HTML)}`
  );
  mainWindow.once("ready-to-show", () => {
    mainWindow.show();
  });

  // Load the actual URL
  console.log(`Loading URL: ${BASE_URL}`);
  mainWindow.loadURL(BASE_URL);
});
