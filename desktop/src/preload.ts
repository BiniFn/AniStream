import { contextBridge, ipcRenderer } from "electron";

contextBridge.exposeInMainWorld("electron", {
  isElectron: true,
  platform: process.platform,
  onFullscreenChange: (callback: (isFullscreen: boolean) => void) => {
    ipcRenderer.on("fullscreen-change", (_event, isFullscreen) => {
      callback(isFullscreen);
    });
  },
  getFullscreen: () => ipcRenderer.invoke("get-fullscreen"),

  // Minimal APIs (no update UI)
  getAppVersion: () => ipcRenderer.invoke("get-app-version"),
  getLogFilePath: () => ipcRenderer.invoke("get-log-file-path"),
});
