<script lang="ts">
	import { apiClient } from '$lib/api/client';
	import { appState } from '$lib/context/state.svelte';

	let oldSettings = JSON.stringify(appState.settings);

	$effect(() => {
		localStorage.setItem(
			'settings',
			JSON.stringify({
				...appState.settings,
				userId: null,
			}),
		);

		if (!appState.user) return;

		const newSettings = JSON.stringify(appState.settings);
		if (newSettings === oldSettings) return;
		oldSettings = newSettings;

		apiClient.POST('/settings', {
			body: { ...appState.settings },
		});
	});
</script>
