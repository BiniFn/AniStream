<script lang="ts">
	import FilterSidebar from '$lib/components/anime/filters/filter-sidebar.svelte';
	import MobileFilters from '$lib/components/anime/filters/mobile-filters.svelte';
	import AnimePageHeader from '$lib/components/anime/layout/anime-page-header.svelte';
	import { createFilterActions, updateUrlWithFilters } from '$lib/utils/filter-actions';
	import { getTotalFilters, type FilterState } from '$lib/utils/filters';
	import type { PageProps } from './$types';
	import AnimeGrid from '$lib/components/anime/display/anime-grid.svelte';
	import EmptyState from '$lib/components/anime/display/empty-state.svelte';
	import { Search } from 'lucide-svelte';

	let { data }: PageProps = $props();

	let filters = $state(data.initialFilters);
	let viewMode = $state<'grid' | 'list'>('grid');
	let isLoading = $state(false);
	let showMobileFilters = $state(false);

	const sortOptions = [
		{ value: 'relevance', label: 'Relevance' },
		{ value: 'updated_at', label: 'Recently Updated' },
		{ value: 'ename', label: 'Name (A-Z)' },
		{ value: 'jname', label: 'Japanese Name' },
		{ value: 'season', label: 'Season' },
		{ value: 'year', label: 'Year' },
		{ value: 'anime_updated_at', label: 'Anime Updated' },
		{ value: 'library_updated_at', label: 'Library Updated' },
	];

	async function updateFilters(newFilters: FilterState) {
		isLoading = true;
		filters = newFilters;
		await updateUrlWithFilters(filters);
		isLoading = false;
	}

	const filterActions = createFilterActions(() => filters, updateFilters);

	function handleViewModeChange(mode: 'grid' | 'list') {
		viewMode = mode;
	}

	function changePage(newPage: number) {
		filterActions.setPage(newPage);
		window.scrollTo({ top: 0, behavior: 'smooth' });
	}

	const totalFilters = $derived(getTotalFilters(filters));
	const totalPages = $derived(data.listings?.pageInfo ? data.listings.pageInfo.totalPages : 0);
</script>

<svelte:head>
	<title>Catalog - Aniways</title>
	<meta name="description" content="Browse and discover anime from our extensive catalog." />
</svelte:head>

<div class="min-h-screen bg-background">
	<AnimePageHeader
		title="Anime Catalog"
		description="Browse and discover anime from our extensive catalog"
		{filters}
		{filterActions}
		{sortOptions}
		bind:viewMode
		{totalFilters}
		pageInfo={data.listings?.pageInfo}
		onViewModeChange={handleViewModeChange}
		onMobileFiltersToggle={() => (showMobileFilters = true)}
	/>

	<MobileFilters
		bind:open={showMobileFilters}
		genres={data.genres || []}
		{filters}
		{filterActions}
		{totalFilters}
		onOpenChange={(open) => (showMobileFilters = open)}
	/>

	<div class="container mx-auto px-4 py-8">
		<div class="flex gap-8">
			<aside class="hidden w-64 shrink-0 lg:block">
				<FilterSidebar genres={data.genres || []} {filters} {filterActions} />
			</aside>

			<AnimeGrid
				anime={data.listings?.items || []}
				{viewMode}
				{isLoading}
				itemsPerPage={filters.itemsPerPage}
				{changePage}
				currentPage={filters.page}
				{totalPages}
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
