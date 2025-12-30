type ElectronAPI = {
	isElectron: boolean;
	platform: 'darwin' | 'win32' | 'linux';
	onFullscreenChange: (callback: (isFullscreen: boolean) => void) => void;
	getFullscreen: () => Promise<boolean>;
	getAppVersion: () => Promise<string>;
	getLogFilePath: () => Promise<string>;
};

function getElectronAPI(): ElectronAPI | undefined {
	if (typeof window === 'undefined') return undefined;
	return (window as unknown as { electron?: ElectronAPI }).electron;
}

/**
 * Check if the app is running inside Electron
 */
export function isElectron(): boolean {
	return !!getElectronAPI()?.isElectron;
}

/**
 * Get the platform (darwin, win32, linux). Returns null if not in Electron.
 */
export function getPlatform(): 'darwin' | 'win32' | 'linux' | null {
	return getElectronAPI()?.platform ?? null;
}

/**
 * Check if running on macOS in Electron
 */
export function isMacOS(): boolean {
	return getPlatform() === 'darwin';
}

/**
 * Subscribe to fullscreen changes (only works in Electron)
 */
export function onFullscreenChange(callback: (isFullscreen: boolean) => void): void {
	getElectronAPI()?.onFullscreenChange(callback);
}

/**
 * Get current fullscreen state (only works in Electron)
 */
export async function getFullscreen(): Promise<boolean> {
	return (await getElectronAPI()?.getFullscreen()) ?? false;
}

/**
 * Get the app version (only works in Electron)
 */
export async function getAppVersion(): Promise<string | null> {
	return (await getElectronAPI()?.getAppVersion()) ?? null;
}

/**
 * Get log file path (only works in Electron)
 */
export async function getLogFilePath(): Promise<string | null> {
	return (await getElectronAPI()?.getLogFilePath()) ?? null;
}
