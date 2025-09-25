<script lang="ts">
	import AnimeCard from '$lib/components/anime/anime-card.svelte';
	import { Pagination } from '$lib/components/ui/pagination';
	import { layoutState } from '$lib/context/layout.svelte';
	import { Play } from 'lucide-svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Continue Watching - Aniways</title>
	<meta name="description" content="Resume watching your anime where you left off" />
</svelte:head>

<div class="min-h-screen bg-background">
	<div
		class="sticky z-30 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
		style="top: {layoutState.navbarHeight}px"
	>
		<div class="container mx-auto px-4 py-4">
			<h1 class="text-2xl font-bold tracking-tight">Continue Watching</h1>
			<p class="text-sm text-muted-foreground">Resume watching your anime where you left off</p>
		</div>
	</div>

	<div class="container mx-auto px-4 py-8">
		{#if data.listings.items.length === 0}
			<div class="flex flex-col items-center justify-center py-16 text-center">
				<div class="mb-4 rounded-full bg-muted p-4">
					<Play class="h-8 w-8 text-muted-foreground" />
				</div>
				<h3 class="mb-2 text-lg font-semibold">No anime to continue watching</h3>
				<p class="text-muted-foreground">
					Start watching some anime to see them appear here. You can resume from where you left off.
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
