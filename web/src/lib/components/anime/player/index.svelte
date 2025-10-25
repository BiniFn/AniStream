<script lang="ts">
	import { preloadData } from '$app/navigation';
	import type { components } from '$lib/api/openapi';
	import { createArtPlayer } from '$lib/components/anime/player/create-player.svelte';
	import { getAppStateContext } from '$lib/context/state.svelte';

	type StreamInfo = components['schemas']['models.StreamingDataResponse'];

	type Props = {
		playerId: string;
		info: StreamInfo;
		nextEpisodeUrl: string | null;
		updateLibrary: () => Promise<void>;
	};

	let { info, playerId, nextEpisodeUrl, updateLibrary }: Props = $props();
	const appState = getAppStateContext();

	let element: HTMLDivElement | null = $state(null);

	$effect(() => {
		if (!nextEpisodeUrl || !appState.settings?.autoNextEpisode) return;
		preloadData(nextEpisodeUrl);
	});

	$effect(() => {
		if (!element) return;

		const player = createArtPlayer({
			id: playerId,
			appState,
			container: element,
			source: info,
			nextEpisodeUrl,
			updateLibrary,
		});

		return () => {
			player.destroy();
		};
	});
</script>

<div class="h-full w-full bg-card" bind:this={element}></div>
