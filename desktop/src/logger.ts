import { app } from "electron";
import fs from "node:fs";
import path from "node:path";

const isDev = process.env.NODE_ENV === "development";
let logFile: string | null = null;

function getLogFilePath(): string {
  if (isDev) {
    return path.join(process.cwd(), "desktop.log");
  }
  
  const userDataPath = app.getPath("userData");
  return path.join(userDataPath, "aniways.log");
}

function writeToFile(message: string) {
  try {
    if (!logFile) {
      logFile = getLogFilePath();
    }
    
    const timestamp = new Date().toISOString();
    const logMessage = `[${timestamp}] ${message}\n`;
    
    fs.appendFileSync(logFile, logMessage, "utf-8");
  } catch (error) {
    // Silently fail if we can't write logs
  }
}

export function log(...args: any[]) {
  const message = args.map((arg) => {
    if (typeof arg === "object") {
      try {
        return JSON.stringify(arg, null, 2);
      } catch {
        return String(arg);
      }
    }
    return String(arg);
  }).join(" ");

  // Always log to console in dev
  if (isDev) {
    console.log("[Logger]", ...args);
  }
  
  // Always write to file
  writeToFile(message);
}

export function logError(...args: any[]) {
  const message = `[ERROR] ${args.map(String).join(" ")}`;
  
  if (isDev) {
    console.error("[Logger]", ...args);
  }
  
  writeToFile(message);
}

export function getLogFilePathForUser(): string {
  if (!logFile) {
    logFile = getLogFilePath();
  }
  return logFile;
}
