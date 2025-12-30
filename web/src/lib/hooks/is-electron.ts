export type UpdateStatus =
	| { status: 'checking' }
	| { status: 'available'; version: string }
	| { status: 'not-available' }
	| { status: 'downloading'; percent: number }
	| { status: 'downloaded'; version: string }
	| { status: 'error'; message: string };

type ElectronAPI = {
	isElectron: boolean;
	platform: 'darwin' | 'win32' | 'linux';
	onFullscreenChange: (callback: (isFullscreen: boolean) => void) => void;
	getFullscreen: () => Promise<boolean>;
	getAppVersion: () => Promise<string>;
	getUpdateStatus: () => Promise<UpdateStatus>;
	checkForUpdates: () => Promise<void>;
	startUpdate: () => Promise<void>;
	quitAndInstall: () => Promise<void>;
	onUpdateStatus: (callback: (status: UpdateStatus) => void) => void;
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
 * Get current update status (only works in Electron)
 */
export async function getUpdateStatus(): Promise<UpdateStatus | null> {
	return (await getElectronAPI()?.getUpdateStatus()) ?? null;
}

/**
 * Check for updates (only works in Electron)
 */
export async function checkForUpdates(): Promise<void> {
	await getElectronAPI()?.checkForUpdates();
}

/**
 * Start downloading update (only works in Electron)
 */
export async function startUpdate(): Promise<void> {
	await getElectronAPI()?.startUpdate();
}

/**
 * Quit and install the downloaded update (only works in Electron)
 */
export async function quitAndInstall(): Promise<void> {
	await getElectronAPI()?.quitAndInstall();
}

/**
 * Subscribe to update status changes (only works in Electron)
 */
export function onUpdateStatus(callback: (status: UpdateStatus) => void): void {
	getElectronAPI()?.onUpdateStatus(callback);
}
