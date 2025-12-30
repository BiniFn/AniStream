import { contextBridge, ipcRenderer } from "electron";

export type UpdateStatus =
  | { status: "checking" }
  | { status: "available"; version: string }
  | { status: "not-available" }
  | { status: "downloading"; percent: number }
  | { status: "downloaded"; version: string }
  | { status: "error"; message: string };

contextBridge.exposeInMainWorld("electron", {
  isElectron: true,
  platform: process.platform,
  onFullscreenChange: (callback: (isFullscreen: boolean) => void) => {
    ipcRenderer.on("fullscreen-change", (_event, isFullscreen) => {
      callback(isFullscreen);
    });
  },
  getFullscreen: () => ipcRenderer.invoke("get-fullscreen"),

  // Auto-update APIs
  getAppVersion: () => ipcRenderer.invoke("get-app-version"),
  getUpdateStatus: () => ipcRenderer.invoke("get-update-status"),
  checkForUpdates: () => ipcRenderer.invoke("check-for-updates"),
  startUpdate: () => ipcRenderer.invoke("start-update"),
  quitAndInstall: () => ipcRenderer.invoke("quit-and-install"),
  onUpdateStatus: (callback: (status: UpdateStatus) => void) => {
    ipcRenderer.on("update-status", (_event, status) => {
      callback(status);
    });
  },
});
