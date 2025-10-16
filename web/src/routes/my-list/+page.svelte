<script lang="ts">
	import FilterSidebar from '$lib/components/anime/filters/filter-sidebar.svelte';
	import LibraryStatusTabs from '$lib/components/anime/controls/library-status-tabs.svelte';
	import MobileFilters from '$lib/components/anime/filters/mobile-filters.svelte';
	import AnimePageHeader from '$lib/components/anime/layout/anime-page-header.svelte';
	import { Label } from '$lib/components/ui/label';
	import { CircleCheck, CirclePlay, CircleX, Clock, Pause } from 'lucide-svelte';
	import { createFilterActions, updateUrlWithFilters } from '$lib/utils/filter-actions';
	import { getTotalFilters, type FilterState } from '$lib/utils/filters';
	import type { PageProps } from './$types';
	import AnimeGrid from '$lib/components/anime/display/anime-grid.svelte';
	import EmptyState from '$lib/components/anime/display/empty-state.svelte';
	import { getLayoutStateContext } from '$lib/context/layout.svelte';

	let { data }: PageProps = $props();
	const layoutState = getLayoutStateContext();

	let filters = $state(data.initialFilters);
	let viewMode = $state<'grid' | 'list'>('grid');
	let isLoading = $state(false);
	let showMobileFilters = $state(false);

	const sortOptions = [
		{ value: 'library_updated_at', label: 'Library Updated' },
		{ value: 'ename', label: 'Name (A-Z)' },
		{ value: 'jname', label: 'Japanese Name' },
		{ value: 'anime_updated_at', label: 'Anime Updated' },
		{ value: 'relevance', label: 'Relevance' },
		{ value: 'season', label: 'Season' },
		{ value: 'year', label: 'Year' },
	];

	const statusTabs = [
		{ value: 'watching', label: 'Watching', icon: CirclePlay, iconColor: 'text-blue-500' },
		{ value: 'planning', label: 'Plan to Watch', icon: Clock, iconColor: 'text-orange-500' },
		{ value: 'completed', label: 'Completed', icon: CircleCheck, iconColor: 'text-emerald-500' },
		{ value: 'paused', label: 'On Hold', icon: Pause, iconColor: 'text-amber-500' },
		{ value: 'dropped', label: 'Dropped', icon: CircleX, iconColor: 'text-red-500' },
	] as const;

	const currentTab = $derived(statusTabs.find((tab) => tab.value === data.status) || statusTabs[0]);
	const totalFilters = $derived(getTotalFilters(filters));
	const totalPages = $derived(data.listings?.pageInfo ? data.listings.pageInfo.totalPages : 0);

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

	function changeStatus(newStatus: FilterState['status']) {
		filterActions.setStatus(newStatus);
	}

	function changePage(newPage: number) {
		filterActions.setPage(newPage);
		window.scrollTo({ top: 0, behavior: 'smooth' });
	}
</script>

<svelte:head>
	<title>My List - Aniways</title>
	<meta name="description" content="Manage your anime watchlist and track your viewing progress" />
</svelte:head>

<div class="min-h-screen bg-background">
	<AnimePageHeader
		title="My Anime List"
		description="Track and manage your anime collection"
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
	>
		<div class="flex flex-col gap-2">
			<Label class="text-xs font-medium">Filter by Status</Label>
			<LibraryStatusTabs
				class="overflow-x-auto"
				currentStatus={data.status}
				onStatusChange={changeStatus}
			/>
		</div>
	</MobileFilters>

	<div class="container mx-auto px-4 pt-4 pb-8">
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
				<div class="sticky z-10 mb-4 hidden lg:block" style="top: {layoutState.totalHeight}px">
					<div
						class="border-b bg-background/95 px-4 py-3 backdrop-blur supports-[backdrop-filter]:bg-background/60"
					>
						<LibraryStatusTabs currentStatus={data.status} onStatusChange={changeStatus} />
					</div>
				</div>

				<div
					class="mb-4 max-w-[calc(100vw-32px)] overflow-x-auto border-b bg-background py-3 lg:hidden"
				>
					<LibraryStatusTabs currentStatus={data.status} onStatusChange={changeStatus} />
				</div>

				{#snippet empty()}
					<EmptyState
						icon={currentTab.icon}
						title="No anime in {currentTab.label}"
						description={data.status === 'watching'
							? 'Start watching anime to see them here. Your progress will be automatically tracked.'
							: data.status === 'planning'
								? "Add anime you plan to watch later. Keep track of shows you're interested in."
								: data.status === 'completed'
									? "Anime you've finished watching will appear here. Complete a series to add it."
									: data.status === 'paused'
										? "Anime you've put on hold will appear here. Take a break and come back later."
										: data.status === 'dropped'
											? "Anime you've decided not to continue will appear here."
											: ''}
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
