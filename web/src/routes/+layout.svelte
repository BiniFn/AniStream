<script lang="ts">
	import '../app.css';
	import type { LayoutProps } from './$types';
	import NavBar from './nav-bar.svelte';
	import SeasonalAnime from './seasonal-anime.svelte';

	let { children, data }: LayoutProps = $props();

	let seasonalAnime = $derived(data.seasonalAnime.data);
	let seasonalAnimeError = $derived(data.seasonalAnime.error);
</script>

<NavBar />
{#if seasonalAnimeError}
	<p class="text-red-500">
		Error loading seasonal anime: {seasonalAnimeError.error ??
			'Something went wrong try again later'}
	</p>
{:else if seasonalAnime}
	<SeasonalAnime {seasonalAnime} />
{/if}
<div class="h-full w-full p-4">
	{@render children?.()}
</div>
