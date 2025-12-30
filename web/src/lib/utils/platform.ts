export type OS = 'macos' | 'windows' | 'linux' | 'unknown';

export function detectOS(): OS {
	if (typeof navigator === 'undefined') return 'unknown';

	const userAgent = navigator.userAgent.toLowerCase();
	const platform = navigator.platform?.toLowerCase() || '';

	if (userAgent.includes('mac') || platform.includes('mac')) {
		return 'macos';
	}
	if (userAgent.includes('win') || platform.includes('win')) {
		return 'windows';
	}
	if (userAgent.includes('linux') || platform.includes('linux')) {
		return 'linux';
	}

	return 'unknown';
}

export function getPlatformDisplayName(platform: string): string {
	const names: Record<string, string> = {
		'darwin-arm64': 'macOS (Apple Silicon)',
		'darwin-x64': 'macOS (Intel)',
		'win32-x64': 'Windows (64-bit)',
		'win32-arm64': 'Windows (ARM)',
		'linux-x64': 'Linux (64-bit)',
		'linux-arm64': 'Linux (ARM)',
	};
	return names[platform] || platform;
}

export function getOSDisplayName(os: OS): string {
	const names: Record<OS, string> = {
		macos: 'macOS',
		windows: 'Windows',
		linux: 'Linux',
		unknown: 'Unknown',
	};
	return names[os];
}

export function formatFileSize(bytes: number): string {
	const units = ['B', 'KB', 'MB', 'GB'];
	let size = bytes;
	let unitIndex = 0;

	while (size >= 1024 && unitIndex < units.length - 1) {
		size /= 1024;
		unitIndex++;
	}

	return `${size.toFixed(1)} ${units[unitIndex]}`;
}
