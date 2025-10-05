<script lang="ts">
	import AnimeCard from '$lib/components/anime/display/anime-card.svelte';
	import EmptyState from '$lib/components/anime/display/empty-state.svelte';
	import PageHeader from '$lib/components/layout/page-header.svelte';
	import { Pagination } from '$lib/components/ui/pagination';
	import { Clock } from 'lucide-svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Planning - Aniways</title>
	<meta name="description" content="Anime you plan to watch" />
</svelte:head>

<div class="min-h-screen bg-background">
	<PageHeader title="Planning" description="Anime you plan to watch" />

	<div class="container mx-auto px-4 py-8">
		{#if data.listings.items.length === 0}
			<EmptyState
				icon={Clock}
				title="No anime in your planning list"
				description="Add anime to your planning list to see them appear here. You can start watching anytime."
				action={{
					label: 'Browse Catalog',
					href: '/catalog',
					variant: 'default',
				}}
			/>
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
