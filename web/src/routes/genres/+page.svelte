<script lang="ts">
	import GenreCard from '$lib/components/anime/genre-card.svelte';
	import { layoutState } from '$lib/context/layout.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	let genres = $derived(data.genres);
</script>

<svelte:head>
	<title>Genres - Aniways</title>
	<meta name="description" content="Browse anime by genre and discover new favorites." />
</svelte:head>

<div class="min-h-screen bg-background">
	<div
		class="sticky z-30 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
		style="top: {layoutState.navbarHeight}px"
	>
		<div class="container mx-auto px-4 py-4">
			<h1 class="text-2xl font-bold tracking-tight">Browse by Genre</h1>
			<p class="text-sm text-muted-foreground">Pick a genre to explore the catalog</p>
		</div>
	</div>

	<div class="container mx-auto px-4 py-8">
		{#if genres.length === 0}
			<div class="rounded-lg border p-6 text-center text-muted-foreground">
				No genres available.
			</div>
		{:else}
			<div class="grid grid-cols-2 gap-3 sm:grid-cols-3 sm:gap-6">
				{#each genres as g (g.name)}
					<GenreCard {...g} />
				{/each}
			</div>
		{/if}
	</div>
</div>
