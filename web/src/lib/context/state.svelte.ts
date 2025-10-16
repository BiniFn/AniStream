import { apiClient } from '$lib/api/client';
import type { components } from '$lib/api/openapi';
import { getContext, setContext } from 'svelte';
import { toast } from 'svelte-sonner';

type User = components['schemas']['models.UserResponse'];
type Settings = Omit<components['schemas']['models.SettingsResponse'], 'userId'>;

export class AppState {
	user: User | null = $state(null);
	settings: Settings | null = $state(null);
	searchOpen = $state(false);
	isLoggedIn = $derived(this.user !== null);
	importJobId: string | null = $state(null);

	private oldSettings: string | null = null;
	private importJobPollingTimeout: NodeJS.Timeout | null = null;

	constructor(
		user?: User | null,
		settings?: Settings | null,
		searchOpen?: boolean,
		importJobId?: string | null,
	) {
		if (user) this.user = user;
		if (searchOpen !== undefined) this.searchOpen = searchOpen;
		this.settings = settings ?? this.getDefaultSettings();
		this.importJobId = importJobId ?? this.getDefaultImportJobId();

		$effect(() => {
			localStorage.setItem('settings', JSON.stringify(this.settings));

			if (!this.isLoggedIn) return;

			const newSettings = JSON.stringify(this.settings);
			if (this.oldSettings === newSettings) return;
			this.oldSettings = newSettings;

			apiClient.POST('/settings', {
				body: { ...this.settings },
			});
		});

		$effect(() => {
			if (this.importJobId) {
				localStorage.setItem('importJobId', this.importJobId);
			} else {
				localStorage.removeItem('importJobId');
			}

			if (this.importJobId && !this.importJobPollingTimeout) {
				this.importJobPollingTimeout = setInterval(() => {
					this.checkImportStatus();
				}, 5000);
			} else if (!this.importJobId && this.importJobPollingTimeout) {
				clearInterval(this.importJobPollingTimeout);
				this.importJobPollingTimeout = null;
			}

			return () => {
				if (this.importJobPollingTimeout) {
					clearInterval(this.importJobPollingTimeout);
					this.importJobPollingTimeout = null;
				}
			};
		});
	}

	private async checkImportStatus() {
		if (!this.importJobId) return;

		try {
			const response = await apiClient.GET('/library/import/{id}', {
				params: {
					path: { id: this.importJobId },
				},
			});

			if (response.data) {
				const status = response.data.status;

				if (status === 'completed') {
					clearInterval(this.importJobPollingTimeout!);
					this.importJobPollingTimeout = null;
					toast.success('Library import completed successfully!');
				} else if (status === 'failed') {
					clearInterval(this.importJobPollingTimeout!);
					this.importJobPollingTimeout = null;
					toast.error('Library import failed. Please try again.');
				}
			}
		} catch (error) {
			console.error('Failed to check import status:', error);
		}
	}

	getDefaultSettings(): Settings {
		if (typeof localStorage !== 'undefined' && localStorage.getItem('settings')) {
			const s = JSON.parse(localStorage.getItem('settings')!);
			return s;
		}

		return {
			autoNextEpisode: true,
			autoPlayEpisode: true,
			incognitoMode: false,
			autoResumeEpisode: true,
		};
	}

	getDefaultImportJobId(): string | null {
		if (typeof localStorage === 'undefined') return null;
		return localStorage ? localStorage.getItem('importJobId') : null;
	}

	logout() {
		this.user = null;
		this.settings = this.getDefaultSettings();
	}

	setUser(u: User | null) {
		this.user = u;
	}

	setSettings(s: Settings | null) {
		if (s === null) {
			this.settings = this.getDefaultSettings();
			return;
		}
		this.settings = s;
	}

	toggleSetting<K extends keyof Settings>(key: K) {
		if (!this.settings) return;
		this.settings[key] = !this.settings[key];
		return this.settings[key];
	}

	setImportJobId(id: string | null) {
		this.importJobId = id;
	}
}

const APPSTATE_CONTEXT_KEY = Symbol('AppState');

export const setAppStateContext = (
	user?: User | null,
	settings?: Settings | null,
	searchOpen?: boolean,
	key = APPSTATE_CONTEXT_KEY,
) => {
	const state = new AppState(user, settings, searchOpen);
	return setContext(key, state);
};

export const getAppStateContext = (key = APPSTATE_CONTEXT_KEY) => {
	const appState = getContext<AppState>(key);
	if (!appState) {
		throw new Error('AppState context not found');
	}
	return appState;
};
