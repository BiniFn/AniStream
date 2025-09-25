<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import ActiveFilters from '$lib/components/anime/active-filters.svelte';
	import AnimeCard from '$lib/components/anime/anime-card.svelte';
	import FilterSidebar from '$lib/components/anime/filter-sidebar.svelte';
	import LibraryStatusTabs from '$lib/components/anime/library-status-tabs.svelte';
	import MobileFilters from '$lib/components/anime/mobile-filters.svelte';
	import PageHeader from '$lib/components/anime/page-header.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import { Label } from '$lib/components/ui/label';
	import { Pagination } from '$lib/components/ui/pagination';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { layoutState } from '$lib/context/layout.svelte';
	import { cn } from '$lib/utils';
	import { ChevronRight, CircleCheck, CirclePlay, CircleX, Clock, Pause } from 'lucide-svelte';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let searchQuery = $state(data.initialQuery.search);
	let selectedGenres = $state(data.initialQuery.genres);
	let genresMode = $state(data.initialQuery.genresMode);
	let selectedSeasons = $state(data.initialQuery.seasons);
	let selectedYears = $state(data.initialQuery.years);
	let sortBy = $state(data.initialQuery.sortBy);
	let sortOrder = $state(data.initialQuery.sortOrder);
	let itemsPerPage = $state(data.initialQuery.itemsPerPage);
	let currentPage = $state(data.initialQuery.page);
	let yearMin = $state(data.initialQuery.yearMin);
	let yearMax = $state(data.initialQuery.yearMax);

	let viewMode = $state<'grid' | 'list'>('grid');
	let isLoading = $state(false);
	let showMobileFilters = $state(false);
	let searchTimeout: NodeJS.Timeout;

	const statusTabsTop = $derived(layoutState.navbarHeight + layoutState.headerHeight);

	const sortOptions = [
		{ value: 'library_updated_at', label: 'Recently Updated' },
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
	const hasContent = $derived(data.listings.items.length > 0);

	const totalFilters = $derived(
		selectedGenres.length +
			selectedSeasons.length +
			selectedYears.length +
			(yearMin ? 1 : 0) +
			(yearMax ? 1 : 0) +
			(searchQuery ? 1 : 0),
	);

	const totalPages = $derived(data.listings?.pageInfo ? data.listings.pageInfo.totalPages : 0);

	async function updateFilters() {
		isLoading = true;
		const params = new SvelteURLSearchParams();

		params.set('status', data.status);
		if (searchQuery) params.set('search', searchQuery);
		selectedGenres.forEach((g) => params.append('genres', g));
		if (genresMode !== 'any') params.set('genresMode', genresMode);
		selectedSeasons.forEach((s) => params.append('seasons', s));
		selectedYears.forEach((y) => params.append('years', y.toString()));
		if (yearMin) params.set('yearMin', yearMin.toString());
		if (yearMax) params.set('yearMax', yearMax.toString());
		if (sortBy !== 'library_updated_at') params.set('sortBy', sortBy);
		if (sortOrder !== 'desc') params.set('sortOrder', sortOrder);
		if (itemsPerPage !== 24) params.set('itemsPerPage', itemsPerPage.toString());
		if (currentPage !== 1) params.set('page', currentPage.toString());

		await goto(`?${params.toString()}`, { invalidateAll: true });
		isLoading = false;
	}

	function handleSearch(value: string) {
		searchQuery = value;
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			currentPage = 1;
			updateFilters();
		}, 500);
	}

	function handleSortChange(newSortBy: string, newSortOrder: 'asc' | 'desc') {
		sortBy = newSortBy as typeof sortBy;
		sortOrder = newSortOrder;
		updateFilters();
	}

	function handleViewModeChange(mode: 'grid' | 'list') {
		viewMode = mode;
	}

	function changeStatus(newStatus: string) {
		const url = new URL(page.url);
		url.searchParams.set('status', newStatus);
		url.searchParams.set('page', '1');
		url.searchParams.delete('search');
		url.searchParams.delete('genres');
		url.searchParams.delete('seasons');
		url.searchParams.delete('years');
		url.searchParams.delete('yearMin');
		url.searchParams.delete('yearMax');
		goto(url.toString());
	}

	function toggleGenre(genre: string) {
		if (selectedGenres.map((g) => g.toLowerCase()).includes(genre.toLowerCase())) {
			selectedGenres = selectedGenres.filter((g) => g.toLowerCase() !== genre.toLowerCase());
		} else {
			selectedGenres = [...selectedGenres, genre];
		}
		currentPage = 1;
		updateFilters();
	}

	function toggleSeason(season: string) {
		const typedSeason = season as 'winter' | 'spring' | 'summer' | 'fall' | 'unknown';
		if (selectedSeasons.includes(typedSeason)) {
			selectedSeasons = selectedSeasons.filter((s) => s !== typedSeason);
		} else {
			selectedSeasons = [...selectedSeasons, typedSeason];
		}
		currentPage = 1;
		updateFilters();
	}

	function toggleYear(year: number) {
		if (selectedYears.includes(year)) {
			selectedYears = selectedYears.filter((y) => y !== year);
		} else {
			selectedYears = [...selectedYears, year];
		}
		currentPage = 1;
		updateFilters();
	}

	function handleYearRangeChange(min?: number, max?: number) {
		yearMin = min;
		yearMax = max;
		currentPage = 1;
		updateFilters();
	}

	function handleGenresModeChange(mode: 'any' | 'all') {
		genresMode = mode;
		updateFilters();
	}

	function handleItemsPerPageChange(count: number) {
		itemsPerPage = count;
		currentPage = 1;
		updateFilters();
	}

	function clearAllFilters() {
		searchQuery = '';
		selectedGenres = [];
		selectedSeasons = [];
		selectedYears = [];
		yearMin = undefined;
		yearMax = undefined;
		genresMode = 'any';
		sortBy = 'library_updated_at';
		sortOrder = 'desc';
		currentPage = 1;
		updateFilters();
	}

	function changePage(newPage: number) {
		currentPage = newPage;
		updateFilters();
		window.scrollTo({ top: 0, behavior: 'smooth' });
	}
</script>

<svelte:head>
	<title>My List - Aniways</title>
	<meta name="description" content="Manage your anime watchlist and track your viewing progress" />
</svelte:head>

<div class="min-h-screen bg-background">
	<PageHeader
		title="My Anime List"
		description="Track and manage your anime collection"
		bind:searchQuery
		bind:sortBy
		bind:sortOrder
		bind:viewMode
		{sortOptions}
		{totalFilters}
		pageInfo={data.listings?.pageInfo}
		onSearchChange={handleSearch}
		onSortChange={handleSortChange}
		onViewModeChange={handleViewModeChange}
		onMobileFiltersToggle={() => (showMobileFilters = true)}
	>
		<ActiveFilters
			{searchQuery}
			{selectedGenres}
			{selectedSeasons}
			{selectedYears}
			{totalFilters}
			onClearSearch={() => {
				searchQuery = '';
				handleSearch('');
			}}
			onRemoveGenre={toggleGenre}
			onRemoveSeason={toggleSeason}
			onRemoveYear={toggleYear}
			onClearAll={clearAllFilters}
			class="mt-4"
		/>
	</PageHeader>

	<MobileFilters
		bind:open={showMobileFilters}
		genres={data.genres || []}
		{selectedGenres}
		{selectedSeasons}
		{selectedYears}
		{yearMin}
		{yearMax}
		{genresMode}
		{totalFilters}
		onOpenChange={(open) => (showMobileFilters = open)}
		onGenreToggle={toggleGenre}
		onSeasonToggle={toggleSeason}
		onYearToggle={toggleYear}
		onYearRangeChange={handleYearRangeChange}
		onGenresModeChange={handleGenresModeChange}
		onClearAll={clearAllFilters}
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
				<FilterSidebar
					genres={data.genres || []}
					{selectedGenres}
					{selectedSeasons}
					{selectedYears}
					{yearMin}
					{yearMax}
					{genresMode}
					{itemsPerPage}
					onGenreToggle={toggleGenre}
					onSeasonToggle={toggleSeason}
					onYearToggle={toggleYear}
					onYearRangeChange={handleYearRangeChange}
					onGenresModeChange={handleGenresModeChange}
					onItemsPerPageChange={handleItemsPerPageChange}
				/>
			</aside>

			<main class="flex-1">
				<div class="sticky z-10 hidden lg:block" style="top: {statusTabsTop}px">
					<div
						class="border-b bg-background/95 px-4 py-3 backdrop-blur supports-[backdrop-filter]:bg-background/60"
					>
						<LibraryStatusTabs currentStatus={data.status} onStatusChange={changeStatus} />
					</div>
				</div>

				<div class="max-w-[calc(100vw-32px)] overflow-x-auto border-b bg-background py-3 lg:hidden">
					<LibraryStatusTabs currentStatus={data.status} onStatusChange={changeStatus} />
				</div>

				<div class="px-4 pt-4">
					{#if isLoading}
						<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
							{#each Array(itemsPerPage) as _, i (i)}
								<div class="space-y-2">
									<Skeleton class="aspect-[3/4] rounded-xl" />
									<Skeleton class="h-4 w-3/4" />
									<Skeleton class="h-3 w-1/2" />
								</div>
							{/each}
						</div>
					{:else if hasContent}
						{#if viewMode === 'grid'}
							<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
								{#each data.listings.items as item, index (item.id)}
									<div
										class="animate-in duration-500 slide-in-from-bottom-5"
										style="animation-delay: {index * 50}ms"
									>
										<AnimeCard
											anime={item}
											{index}
											libraryEntry={item.library
												? {
														id: item.library.id,
														userId: '',
														animeId: item.id,
														status: item.library.status,
														watchedEpisodes: item.library.watchedEpisodes,
														createdAt: '',
														updatedAt: item.library.updatedAt,
														anime: item,
													}
												: null}
											showLibraryInfo={true}
										/>
									</div>
								{/each}
							</div>
						{:else}
							<div class="animate-in space-y-4 duration-500 fade-in">
								{#each data.listings.items as item, i (item.id)}
									<a
										href="/anime/{item.id}"
										class="flex animate-in gap-3 rounded-lg border bg-card p-3 transition-all duration-500 slide-in-from-left-5 hover:shadow-lg sm:p-4"
										style="animation-delay: {i * 50}ms"
									>
										<img
											src={item.imageUrl}
											alt={item.ename || item.jname}
											class="h-24 w-20 rounded-md object-cover sm:h-32 sm:w-24"
										/>
										<div class="flex-1 space-y-2">
											<div>
												<h3 class="line-clamp-2 text-sm font-semibold sm:text-lg">
													{item.ename || item.jname}
												</h3>
												{#if item.ename && item.jname}
													<p class="text-sm text-muted-foreground">{item.jname}</p>
												{/if}
											</div>

											{#if item.library}
												<div class="flex items-center gap-2">
													<Badge variant="secondary" class="text-xs capitalize">
														{item.library.status.replace('_', ' ')}
													</Badge>
													{#if item.library.watchedEpisodes > 0}
														<span class="text-xs text-muted-foreground">
															Ep {item.library.watchedEpisodes}
															{#if item.lastEpisode}
																/ {item.lastEpisode}
															{/if}
														</span>
													{/if}
												</div>
											{/if}

											<div class="flex flex-wrap gap-2">
												{#each item.genre.split(', ').slice(0, 5) as genre (genre)}
													<Badge variant="secondary" class="text-xs">{genre}</Badge>
												{/each}
											</div>
											<div class="flex gap-4 text-sm text-muted-foreground">
												<span class="capitalize">{item.season} {item.seasonYear}</span>
												{#if item.lastEpisode}
													<span>Episode {item.lastEpisode}</span>
												{/if}
											</div>
										</div>
									</a>
								{/each}
							</div>
						{/if}

						<Pagination {totalPages} {currentPage} onPageChange={changePage} />
					{:else}
						<Card.Root class="border-dashed">
							<Card.Content class="flex flex-col items-center justify-center py-20 text-center">
								{@const Icon = currentTab.icon}
								<Icon class={cn('mb-4 h-16 w-16', currentTab.iconColor)} />
								<h3 class="mb-2 text-xl font-semibold">No anime in {currentTab.label}</h3>
								<p class="mb-6 max-w-md text-muted-foreground">
									{#if data.status === 'watching'}
										Start watching anime to see them here. Your progress will be automatically
										tracked.
									{:else if data.status === 'planning'}
										Add anime you plan to watch later. Keep track of shows you're interested in.
									{:else if data.status === 'completed'}
										Anime you've finished watching will appear here. Complete a series to add it.
									{:else if data.status === 'paused'}
										Anime you've put on hold will appear here. Take a break and come back later.
									{:else if data.status === 'dropped'}
										Anime you've decided not to continue will appear here.
									{/if}
								</p>
								<div class="flex gap-3">
									<Button href="/catalog" variant="default" class="gap-2">
										Browse Catalog
										<ChevronRight class="h-4 w-4" />
									</Button>
									<Button href="/genres" variant="outline" class="gap-2">Explore Genres</Button>
								</div>
							</Card.Content>
						</Card.Root>
					{/if}
				</div>
			</main>
		</div>
	</div>
</div>

<style>
	@keyframes slide-in-from-bottom-5 {
		from {
			transform: translateY(20px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	@keyframes slide-in-from-left-5 {
		from {
			transform: translateX(-20px);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}

	.animate-in {
		animation-fill-mode: both;
	}

	.fade-in {
		animation: fade-in 0.5s ease-out;
	}

	.slide-in-from-bottom-5 {
		animation: slide-in-from-bottom-5 0.5s ease-out;
		animation-fill-mode: both;
	}

	.slide-in-from-left-5 {
		animation: slide-in-from-left-5 0.5s ease-out;
		animation-fill-mode: both;
	}

	@keyframes fade-in {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}
</style>
