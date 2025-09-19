import type { components } from '$lib/api/openapi';

type State = {
	user: components['schemas']['models.UserResponse'] | null;
	settings: Omit<components['schemas']['models.SettingsResponse'], 'userId'> | null;
	searchOpen: boolean;
};

function getDefaultSettings(): State['settings'] {
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

export const appState = $state<State>({
	settings: getDefaultSettings(),
	user: null,
	searchOpen: false,
});

export function setUser(u: State['user']) {
	appState.user = u;
}

export function setSettings(s: State['settings'] | null) {
	if (s === null) {
		appState.settings = getDefaultSettings();
		return;
	}
	appState.settings = s;
}
