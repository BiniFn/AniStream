<script lang="ts">
	import { goto } from '$app/navigation';
	import ActiveFilters from '$lib/components/anime/active-filters.svelte';
	import AnimeCard from '$lib/components/anime/anime-card.svelte';
	import FilterSidebar from '$lib/components/anime/filter-sidebar.svelte';
	import MobileFilters from '$lib/components/anime/mobile-filters.svelte';
	import PageHeader from '$lib/components/anime/page-header.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Pagination } from '$lib/components/ui/pagination';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let searchQuery = $derived(data.initialQuery.search);
	let selectedGenres = $derived(data.initialQuery.genres);
	let genresMode = $derived(data.initialQuery.genresMode);
	let selectedSeasons = $derived(data.initialQuery.seasons);
	let selectedYears = $derived(data.initialQuery.years);
	let sortBy = $derived(data.initialQuery.sortBy);
	let sortOrder = $derived(data.initialQuery.sortOrder);
	let itemsPerPage = $derived(data.initialQuery.itemsPerPage);
	let currentPage = $derived(data.initialQuery.page);
	let yearMin = $derived(data.initialQuery.yearMin);
	let yearMax = $derived(data.initialQuery.yearMax);

	let viewMode = $state<'grid' | 'list'>('grid');
	let isLoading = $state(false);
	let showMobileFilters = $state(false);
	let searchTimeout: NodeJS.Timeout;

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

	async function updateFilters() {
		isLoading = true;
		const params = new SvelteURLSearchParams();

		if (searchQuery) params.set('search', searchQuery);
		selectedGenres.forEach((g) => params.append('genres', g));
		if (genresMode !== 'any') params.set('genresMode', genresMode);
		selectedSeasons.forEach((s) => params.append('seasons', s));
		selectedYears.forEach((y) => params.append('years', y.toString()));
		if (yearMin) params.set('yearMin', yearMin.toString());
		if (yearMax) params.set('yearMax', yearMax.toString());
		if (sortBy !== 'relevance') params.set('sortBy', sortBy);
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
		sortBy = 'relevance';
		sortOrder = 'desc';
		currentPage = 1;
		updateFilters();
	}

	function changePage(newPage: number) {
		currentPage = newPage;
		updateFilters();
		window.scrollTo({ top: 0, behavior: 'smooth' });
	}

	const totalFilters = $derived(
		selectedGenres.length +
			selectedSeasons.length +
			selectedYears.length +
			(yearMin ? 1 : 0) +
			(yearMax ? 1 : 0) +
			(searchQuery ? 1 : 0),
	);

	const totalPages = $derived(data.listings?.pageInfo ? data.listings.pageInfo.totalPages : 0);
</script>

<svelte:head>
	<title>Catalog - Aniways</title>
	<meta name="description" content="Browse and discover anime from our extensive catalog." />
</svelte:head>

<div class="min-h-screen bg-background">
	<PageHeader
		title="Anime Catalog"
		description="Browse and discover anime from our extensive catalog"
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
	/>

	<div class="container mx-auto px-4 py-8">
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

			<main class="flex-1 px-4">
				{#if isLoading}
					<div
						class={viewMode === 'grid'
							? 'grid grid-cols-2 gap-3 sm:grid-cols-3 sm:gap-4 md:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6'
							: 'space-y-3 sm:space-y-4'}
					>
						{#each Array(itemsPerPage) as _, i (i)}
							<div class="space-y-2">
								<Skeleton class="aspect-[3/4] rounded-xl" />
								<Skeleton class="h-4 w-3/4" />
								<Skeleton class="h-3 w-1/2" />
							</div>
						{/each}
					</div>
				{:else if data.listings?.items && data.listings.items.length > 0}
					{#if viewMode === 'grid'}
						<div
							class="grid animate-in grid-cols-2 gap-3 duration-500 fade-in sm:grid-cols-3 sm:gap-4 md:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6"
						>
							{#each data.listings.items as anime, i (anime.id)}
								<div
									class="animate-in duration-500 slide-in-from-bottom-5"
									style="animation-delay: {i * 50}ms"
								>
									<AnimeCard {anime} index={i} />
								</div>
							{/each}
						</div>
					{:else}
						<div class="animate-in space-y-4 duration-500 fade-in">
							{#each data.listings.items as anime, i (anime.id)}
								<a
									href="/anime/{anime.id}"
									class="flex animate-in gap-3 rounded-lg border bg-card p-3 transition-all duration-500 slide-in-from-left-5 hover:shadow-lg sm:p-4"
									style="animation-delay: {i * 50}ms"
								>
									<img
										src={anime.imageUrl}
										alt={anime.ename || anime.jname}
										class="h-24 w-20 rounded-md object-cover sm:h-32 sm:w-24"
									/>
									<div class="flex-1 space-y-2">
										<div>
											<h3 class="line-clamp-2 text-sm font-semibold sm:text-lg">
												{anime.ename || anime.jname}
											</h3>
											{#if anime.ename && anime.jname}
												<p class="text-sm text-muted-foreground">{anime.jname}</p>
											{/if}
										</div>
										<div class="flex flex-wrap gap-2">
											{#each anime.genre.split(', ').slice(0, 5) as genre (genre)}
												<Badge variant="secondary" class="text-xs">{genre}</Badge>
											{/each}
										</div>
										<div class="flex gap-4 text-sm text-muted-foreground">
											<span class="capitalize">{anime.season} {anime.seasonYear}</span>
											{#if anime.lastEpisode}
												<span>Episode {anime.lastEpisode}</span>
											{/if}
										</div>
									</div>
								</a>
							{/each}
						</div>
					{/if}
				{:else}
					<div class="flex flex-col items-center justify-center py-20 text-center">
						<div class="mb-4 text-6xl">🔍</div>
						<h3 class="mb-2 text-xl font-semibold">No anime found</h3>
						<p class="text-muted-foreground">Try adjusting your filters or search query</p>
						{#if totalFilters > 0}
							<button onclick={clearAllFilters} class="mt-4 text-primary hover:underline">
								Clear all filters
							</button>
						{/if}
					</div>
				{/if}

				<Pagination {totalPages} {currentPage} onPageChange={changePage} />
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
