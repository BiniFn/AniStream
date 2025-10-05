<script lang="ts">
	import GenreCard from '$lib/components/anime/display/genre-card.svelte';
	import PageHeader from '$lib/components/layout/page-header.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	let query = $state('');
	let genres = $derived.by(() => {
		if (!query) return data.genres;
		return data.genres.filter((g) => g.name?.toLowerCase().includes(query.toLowerCase()));
	});
</script>

<svelte:head>
	<title>Genres - Aniways</title>
	<meta name="description" content="Browse anime by genre and discover new favorites." />
</svelte:head>

<div class="min-h-screen bg-background">
	<PageHeader title="Browse by Genre" description="Pick a genre to explore the catalog">
		{#snippet actions()}
			<Input type="search" bind:value={query} placeholder="Search genres..." class="max-w-sm" />
		{/snippet}
	</PageHeader>

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
