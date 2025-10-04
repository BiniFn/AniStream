<script lang="ts">
	import AnimeCard from '$lib/components/anime/anime-card.svelte';
	import { Pagination } from '$lib/components/ui/pagination';
	import { layoutState } from '$lib/context/layout.svelte';
	import { Clock } from 'lucide-svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Planning - Aniways</title>
	<meta name="description" content="Anime you plan to watch" />
</svelte:head>

<div class="min-h-screen bg-background">
	<div
		class="z-30 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 md:sticky"
		style="top: {layoutState.navbarHeight}px"
	>
		<div class="container mx-auto px-4 py-4">
			<h1 class="text-2xl font-bold tracking-tight">Planning</h1>
			<p class="text-sm text-muted-foreground">Anime you plan to watch</p>
		</div>
	</div>

	<div class="container mx-auto px-4 py-8">
		{#if data.listings.items.length === 0}
			<div class="flex flex-col items-center justify-center py-16 text-center">
				<div class="mb-4 rounded-full bg-muted p-4">
					<Clock class="h-8 w-8 text-muted-foreground" />
				</div>
				<h3 class="mb-2 text-lg font-semibold">No anime in your planning list</h3>
				<p class="text-muted-foreground">
					Add anime to your planning list to see them appear here. You can start watching anytime.
				</p>
			</div>
		{:else}
			<div
				class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
			>
				{#each data.listings.items as item, i (item.id)}
					<AnimeCard anime={item.anime} index={i} episodeLink={item.watchedEpisodes + 1} />
				{/each}
			</div>

			{#if data.listings.pageInfo.totalPages > 1}
				<div class="mt-8 flex justify-center">
					<Pagination
						currentPage={data.listings.pageInfo.currentPage}
						totalPages={data.listings.pageInfo.totalPages}
					/>
				</div>
			{/if}
		{/if}
	</div>
</div>
