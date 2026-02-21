<script lang="ts">
	import { Search } from 'lucide-svelte';
	import { watch } from 'runed';
	import AnimeGrid from '$lib/components/anime/display/anime-grid.svelte';
	import EmptyState from '$lib/components/anime/display/empty-state.svelte';
	import FilterSidebar from '$lib/components/anime/filters/filter-sidebar.svelte';
	import MobileFilters from '$lib/components/anime/filters/mobile-filters.svelte';
	import AnimePageHeader from '$lib/components/anime/layout/anime-page-header.svelte';
	import { FilterManager } from '$lib/utils/filter-manager.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const filterManager = new FilterManager(data.initialFilters);

	const sortOptions = [
		{ value: 'relevance', label: 'Relevance' },
		{ value: 'updated_at', label: 'Recently Updated' },
		{ value: 'ename', label: 'Name (A-Z)' },
		{ value: 'jname', label: 'Japanese Name' },
		{ value: 'season', label: 'Season' },
		{ value: 'year', label: 'Year' },
	];

	watch(
		() => filterManager.filters.page,
		() => {
			window.scrollTo({ top: 0, behavior: 'smooth' });
		},
	);
</script>

<svelte:head>
	<title>Catalog - AniStream</title>
	<meta name="description" content="Browse and discover anime from our extensive catalog." />
</svelte:head>

<div class="min-h-screen bg-background">
	<AnimePageHeader
		title="Anime Catalog"
		description="Browse and discover anime from our extensive catalog"
		{filterManager}
		{sortOptions}
		pageInfo={data.listings?.pageInfo}
	/>

	<MobileFilters genres={data.genres || []} {filterManager} />

	<div class="container mx-auto px-4 py-8">
		<div class="flex gap-8">
			<aside class="hidden w-64 shrink-0 lg:block">
				<FilterSidebar genres={data.genres || []} {filterManager} />
			</aside>

			<AnimeGrid
				anime={data.listings?.items || []}
				{filterManager}
				totalPages={data.listings?.pageInfo.totalPages || 1}
			>
				{#snippet empty()}
					<EmptyState
						icon={Search}
						title="No anime found with your current filters."
						description="Try adjusting your filters or explore our catalog and genres to find something interesting."
						action={{
							label: 'Browse Catalog',
							href: '/catalog',
							variant: 'default',
						}}
						secondaryAction={{
							label: 'Explore Genres',
							href: '/genres',
							variant: 'outline',
						}}
					/>
				{/snippet}
			</AnimeGrid>
		</div>
	</div>
</div>
