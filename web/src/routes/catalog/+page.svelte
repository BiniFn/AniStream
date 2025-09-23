<script lang="ts">
	import { goto } from '$app/navigation';
	import AnimeCard from '$lib/components/anime/anime-card.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Pagination } from '$lib/components/ui/pagination';
	import * as Select from '$lib/components/ui/select';
	import * as Sheet from '$lib/components/ui/sheet';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { cn } from '$lib/utils';
	import {
		Calendar,
		Funnel,
		Hash,
		LayoutGrid,
		LayoutList,
		Search,
		Sparkles,
		X,
	} from 'lucide-svelte';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let searchQuery = $derived(data.initialQuery.search);
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
	let showFilters = $state(false);
	let searchTimeout: NodeJS.Timeout;
	let yearMinTimeout: NodeJS.Timeout;
	let yearMaxTimeout: NodeJS.Timeout;

	const seasons = ['winter', 'spring', 'summer', 'fall', 'unknown'] as const;
	const currentYear = new Date().getFullYear();
	const years = Array.from({ length: 50 }, (_, i) => currentYear - i);

	const sortOptions = [
		{ value: 'relevance', label: 'Relevance' },
		{ value: 'updated_at', label: 'Recently Updated' },
		{ value: 'ename', label: 'Name (A-Z)' },
		{ value: 'jname', label: 'Japanese Name' },
		{ value: 'season', label: 'Season' },
		{ value: 'year', label: 'Year' },
	];

	const itemsPerPageOptions = [12, 24, 36, 48];

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

	function toggleGenre(genre: string) {
		if (selectedGenres.map((g) => g.toLowerCase()).includes(genre.toLowerCase())) {
			selectedGenres = selectedGenres.filter((g) => g.toLowerCase() !== genre.toLowerCase());
		} else {
			selectedGenres = [...selectedGenres, genre];
		}
		currentPage = 1;
		updateFilters();
	}

	function toggleSeason(season: (typeof seasons)[number]) {
		if (selectedSeasons.includes(season)) {
			selectedSeasons = selectedSeasons.filter((s) => s !== season);
		} else {
			selectedSeasons = [...selectedSeasons, season];
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

<div class="min-h-screen bg-background">
	<div
		class="sticky top-17 z-30 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
	>
		<div class="container mx-auto px-3 py-3 sm:px-4 sm:py-4">
			<div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between lg:gap-4">
				<div class="flex items-center justify-between gap-2 lg:hidden">
					<div class="flex-1">
						<h1 class="text-lg font-bold">Anime Catalog</h1>
						<p class="text-xs text-muted-foreground">
							{#if data.listings?.pageInfo}
								Page {data.listings.pageInfo.currentPage}/{data.listings.pageInfo.totalPages}
							{/if}
						</p>
					</div>
					<div class="flex items-center gap-2">
						<div class="flex gap-1 rounded-md border p-1">
							<Button
								variant={viewMode === 'grid' ? 'default' : 'ghost'}
								size="icon"
								class="h-7 w-7"
								onclick={() => (viewMode = 'grid')}
							>
								<LayoutGrid class="h-3 w-3" />
							</Button>
							<Button
								variant={viewMode === 'list' ? 'default' : 'ghost'}
								size="icon"
								class="h-7 w-7"
								onclick={() => (viewMode = 'list')}
							>
								<LayoutList class="h-3 w-3" />
							</Button>
						</div>
						<Sheet.Root bind:open={showFilters}>
							<Sheet.Trigger
								class={cn(buttonVariants({ variant: 'outline', size: 'icon' }), 'relative h-9 w-9')}
							>
								<Funnel class="h-4 w-4" />
								{#if totalFilters > 0}
									<span
										class="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[10px] text-primary-foreground"
									>
										{totalFilters}
									</span>
								{/if}
							</Sheet.Trigger>
							<Sheet.Content side="bottom" class="flex h-[90vh] max-h-[90vh] flex-col">
								<div class="mx-auto flex h-full w-full max-w-lg flex-col">
									<Sheet.Header class="pb-2">
										<div class="mx-auto mb-2 h-1 w-12 rounded-full bg-muted"></div>
										<Sheet.Title class="text-center text-sm">Filters & Sort</Sheet.Title>
									</Sheet.Header>
									<div class="flex-1 space-y-3 overflow-y-auto px-4">
										{#if data.genres && data.genres.length > 0}
											<div class="space-y-2">
												<Label class="text-xs font-medium">Genres</Label>
												<div class="flex max-h-32 flex-wrap gap-1.5 overflow-y-auto pr-2">
													{#each data.genres as genre (genre)}
														<Badge
															variant={selectedGenres
																.map((g) => g.toLowerCase())
																.includes(genre.toLowerCase())
																? 'default'
																: 'outline'}
															class="h-6 cursor-pointer px-2 py-0.5 text-xs"
															onclick={() => toggleGenre(genre)}
														>
															{genre}
														</Badge>
													{/each}
												</div>
												{#if selectedGenres.length > 1}
													<div class="flex gap-2 pt-1">
														<Button
															variant={genresMode === 'any' ? 'default' : 'outline'}
															size="sm"
															onclick={() => {
																genresMode = 'any';
																updateFilters();
															}}
															class="h-7 flex-1 text-xs"
														>
															Any
														</Button>
														<Button
															variant={genresMode === 'all' ? 'default' : 'outline'}
															size="sm"
															onclick={() => {
																genresMode = 'all';
																updateFilters();
															}}
															class="h-7 flex-1 text-xs"
														>
															All
														</Button>
													</div>
												{/if}
											</div>
										{/if}

										<div class="space-y-2">
											<Label class="text-xs font-medium">Season</Label>
											<div class="flex flex-wrap gap-1.5">
												{#each seasons as season (season)}
													<Button
														variant={selectedSeasons.includes(season) ? 'default' : 'outline'}
														size="sm"
														onclick={() => toggleSeason(season)}
														class="h-7 px-3 text-xs capitalize"
													>
														{season}
													</Button>
												{/each}
											</div>
										</div>

										<div class="space-y-2">
											<Label class="text-xs font-medium">Years</Label>
											<div class="flex max-h-24 flex-wrap gap-1.5 overflow-y-auto pr-2">
												{#each years.slice(0, 20) as year (year)}
													<Button
														variant={selectedYears.includes(year) ? 'default' : 'outline'}
														size="sm"
														onclick={() => toggleYear(year)}
														class="h-6 px-2 text-[10px]"
													>
														{year}
													</Button>
												{/each}
											</div>
										</div>

										<div class="space-y-2">
											<Label class="text-xs font-medium">Year Range</Label>
											<div class="grid grid-cols-2 gap-2">
												<Input
													type="number"
													placeholder="From"
													min="1970"
													max={currentYear}
													value={yearMin || ''}
													class="h-8 text-xs"
													oninput={(e) => {
														const value = e.currentTarget.value;
														yearMin = value ? Number(value) : undefined;
														clearTimeout(yearMinTimeout);
														yearMinTimeout = setTimeout(() => {
															currentPage = 1;
															updateFilters();
														}, 800);
													}}
												/>
												<Input
													type="number"
													placeholder="To"
													min="1970"
													max={currentYear}
													value={yearMax || ''}
													class="h-8 text-xs"
													oninput={(e) => {
														const value = e.currentTarget.value;
														yearMax = value ? Number(value) : undefined;
														clearTimeout(yearMaxTimeout);
														yearMaxTimeout = setTimeout(() => {
															currentPage = 1;
															updateFilters();
														}, 800);
													}}
												/>
											</div>
										</div>
									</div>
									<div class="border-t px-4 py-3">
										<div class="flex gap-2">
											{#if totalFilters > 0}
												<Button
													variant="outline"
													size="sm"
													class="flex-1"
													onclick={clearAllFilters}
												>
													Clear All
												</Button>
											{/if}
											<Button
												size="sm"
												class={totalFilters > 0 ? 'flex-1' : 'w-full'}
												onclick={() => (showFilters = false)}
											>
												Close
											</Button>
										</div>
									</div>
								</div>
							</Sheet.Content>
						</Sheet.Root>
					</div>
				</div>

				<div class="hidden lg:block">
					<h1 class="text-2xl font-bold tracking-tight">Anime Catalog</h1>
					<p class="text-sm text-muted-foreground">
						{#if data.listings?.items}
							{data.listings.items.length} anime found
							{#if data.listings.pageInfo}
								(Page {data.listings.pageInfo.currentPage} of {data.listings.pageInfo.totalPages})
							{/if}
						{/if}
					</p>
				</div>

				<div class="flex flex-col gap-2 lg:flex-row lg:items-center lg:gap-3">
					<div class="flex flex-col gap-2 lg:hidden">
						<div class="relative flex-1">
							<Search
								class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground"
							/>
							<Input
								type="text"
								placeholder="Search..."
								value={searchQuery}
								oninput={(e) => handleSearch(e.currentTarget.value)}
								class="h-9 pr-9 pl-9 text-sm"
							/>
							{#if searchQuery}
								<button
									onclick={() => handleSearch('')}
									class="absolute top-1/2 right-3 -translate-y-1/2 text-muted-foreground hover:text-foreground"
								>
									<X class="h-3 w-3" />
								</button>
							{/if}
						</div>
						<div class="flex items-center gap-1">
							<Select.Root
								type="single"
								value={sortBy}
								onValueChange={(value: string | undefined) => {
									if (value) {
										sortBy = value as typeof sortBy;
										updateFilters();
									}
								}}
							>
								<Select.Trigger class="h-9 flex-1 text-sm">
									{sortOptions.find((o) => o.value === sortBy)?.label || 'Sort'}
								</Select.Trigger>
								<Select.Content>
									{#each sortOptions as option (option.value)}
										<Select.Item value={option.value}>{option.label}</Select.Item>
									{/each}
								</Select.Content>
							</Select.Root>
							<Button
								variant="outline"
								size="icon"
								class="h-9 w-9"
								onclick={() => {
									sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
									updateFilters();
								}}
							>
								{sortOrder === 'asc' ? '↑' : '↓'}
							</Button>
						</div>
					</div>

					<div class="relative hidden w-full lg:block lg:w-80">
						<Search
							class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground"
						/>
						<Input
							type="text"
							placeholder="Search anime..."
							value={searchQuery}
							oninput={(e) => handleSearch(e.currentTarget.value)}
							class="pr-10 pl-10"
						/>
						{#if searchQuery}
							<button
								onclick={() => handleSearch('')}
								class="absolute top-1/2 right-3 -translate-y-1/2 text-muted-foreground hover:text-foreground"
							>
								<X class="h-4 w-4" />
							</button>
						{/if}
					</div>

					<div class="hidden items-center gap-2 lg:flex">
						<Select.Root
							type="single"
							value={sortBy}
							onValueChange={(value: string | undefined) => {
								if (value) {
									sortBy = value as typeof sortBy;
									updateFilters();
								}
							}}
						>
							<Select.Trigger class="w-[172px]">
								<span>{sortOptions.find((o) => o.value === sortBy)?.label || 'Sort by'}</span>
							</Select.Trigger>
							<Select.Content>
								{#each sortOptions as option (option.value)}
									<Select.Item value={option.value}>{option.label}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>

						<Button
							variant="outline"
							size="icon"
							onclick={() => {
								sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
								updateFilters();
							}}
						>
							{sortOrder === 'asc' ? '↑' : '↓'}
						</Button>

						<div class="hidden gap-1 rounded-md border p-1 lg:flex">
							<Button
								variant={viewMode === 'grid' ? 'default' : 'ghost'}
								size="icon"
								class="h-8 w-8"
								onclick={() => (viewMode = 'grid')}
							>
								<LayoutGrid class="h-4 w-4" />
							</Button>
							<Button
								variant={viewMode === 'list' ? 'default' : 'ghost'}
								size="icon"
								class="h-8 w-8"
								onclick={() => (viewMode = 'list')}
							>
								<LayoutList class="h-4 w-4" />
							</Button>
						</div>
					</div>
				</div>
			</div>

			{#if totalFilters > 0}
				<div class="mt-4 flex flex-wrap items-center gap-2">
					<span class="text-sm font-medium">Active filters:</span>

					{#if searchQuery}
						<Badge variant="secondary" class="gap-1">
							<Search class="h-3 w-3" />
							{searchQuery}
							<button onclick={() => handleSearch('')} class="ml-1 hover:text-destructive">
								<X class="h-3 w-3" />
							</button>
						</Badge>
					{/if}

					{#each selectedGenres as genre (genre)}
						<Badge variant="secondary" class="gap-1">
							<Hash class="h-3 w-3" />
							{genre}
							<button onclick={() => toggleGenre(genre)} class="ml-1 hover:text-destructive">
								<X class="h-3 w-3" />
							</button>
						</Badge>
					{/each}

					{#each selectedSeasons as season (season)}
						<Badge variant="secondary" class="gap-1 capitalize">
							<Sparkles class="h-3 w-3" />
							{season}
							<button
								onclick={() => toggleSeason(season as (typeof seasons)[number])}
								class="ml-1 hover:text-destructive"
							>
								<X class="h-3 w-3" />
							</button>
						</Badge>
					{/each}

					{#each selectedYears as year (year)}
						<Badge variant="secondary" class="gap-1">
							<Calendar class="h-3 w-3" />
							{year}
							<button onclick={() => toggleYear(year)} class="ml-1 hover:text-destructive">
								<X class="h-3 w-3" />
							</button>
						</Badge>
					{/each}

					<Button variant="ghost" size="sm" onclick={clearAllFilters}>Clear all</Button>
				</div>
			{/if}
		</div>
	</div>

	<div class="container mx-auto px-4 py-8">
		<div class="flex gap-8">
			<aside class="hidden w-64 shrink-0 lg:block">
				<div
					class="filter-scroll sticky top-40 max-h-[calc(100vh-11rem)] space-y-6 overflow-y-auto pr-2"
				>
					{#if data.genres && data.genres.length > 0}
						<div class="space-y-3">
							<div class="flex items-center justify-between">
								<Label class="text-base font-semibold">Genres</Label>
								{#if selectedGenres.length > 0}
									<Badge variant="secondary" class="text-xs">
										{selectedGenres.length}
									</Badge>
								{/if}
							</div>
							<div class="filter-scroll max-h-60 space-y-1 overflow-y-auto pr-2">
								{#each data.genres as genre (genre)}
									{@const isSelected = selectedGenres
										.map((g) => g.toLowerCase())
										.includes(genre.toLowerCase())}
									<button
										onclick={() => toggleGenre(genre)}
										class={cn(
											'flex w-full items-center justify-between rounded-md px-3 py-2 text-sm transition-colors hover:bg-accent',
											isSelected ? 'bg-accent font-medium' : '',
										)}
									>
										<span>{genre}</span>
										{#if isSelected}
											<X class="h-3 w-3" />
										{/if}
									</button>
								{/each}
							</div>
							{#if selectedGenres.length > 1}
								<div class="flex gap-2 pt-2">
									<Button
										variant={genresMode === 'any' ? 'default' : 'outline'}
										size="sm"
										onclick={() => {
											genresMode = 'any';
											updateFilters();
										}}
										class="flex-1"
									>
										Any
									</Button>
									<Button
										variant={genresMode === 'all' ? 'default' : 'outline'}
										size="sm"
										onclick={() => {
											genresMode = 'all';
											updateFilters();
										}}
										class="flex-1"
									>
										All
									</Button>
								</div>
							{/if}
						</div>
					{/if}

					<div class="space-y-3">
						<Label class="text-base font-semibold">Season</Label>
						<div class="flex flex-wrap gap-2">
							{#each seasons as season (season)}
								<Button
									variant={selectedSeasons.includes(season) ? 'default' : 'outline'}
									size="sm"
									onclick={() => toggleSeason(season)}
									class="capitalize"
								>
									{season}
								</Button>
							{/each}
						</div>
					</div>

					<div class="space-y-3">
						<Label class="text-base font-semibold">Years</Label>
						<div class="filter-scroll flex max-h-32 flex-wrap gap-2 overflow-y-auto pr-2">
							{#each years.slice(0, 15) as year (year)}
								<Button
									variant={selectedYears.includes(year) ? 'default' : 'outline'}
									size="sm"
									onclick={() => toggleYear(year)}
									class="h-8 px-3 text-xs"
								>
									{year}
								</Button>
							{/each}
						</div>
					</div>

					<div class="space-y-3">
						<Label class="text-base font-semibold">Year Range</Label>
						<div class="space-y-2">
							<Input
								type="number"
								placeholder="From year"
								min="1970"
								max={currentYear}
								value={yearMin || ''}
								oninput={(e) => {
									const value = e.currentTarget.value;
									yearMin = value ? Number(value) : undefined;
									clearTimeout(yearMinTimeout);
									yearMinTimeout = setTimeout(() => {
										currentPage = 1;
										updateFilters();
									}, 800);
								}}
							/>
							<Input
								type="number"
								placeholder="To year"
								min="1970"
								max={currentYear}
								value={yearMax || ''}
								oninput={(e) => {
									const value = e.currentTarget.value;
									yearMax = value ? Number(value) : undefined;
									clearTimeout(yearMaxTimeout);
									yearMaxTimeout = setTimeout(() => {
										currentPage = 1;
										updateFilters();
									}, 800);
								}}
							/>
						</div>
					</div>

					<div class="space-y-3">
						<Label class="text-base font-semibold">Items per page</Label>
						<Select.Root
							type="single"
							value={itemsPerPage.toString()}
							onValueChange={(value: string | undefined) => {
								if (value) {
									itemsPerPage = Number(value);
									currentPage = 1;
									updateFilters();
								}
							}}
						>
							<Select.Trigger class="w-full">
								<span>{itemsPerPage} items</span>
							</Select.Trigger>
							<Select.Content>
								{#each itemsPerPageOptions as option (option)}
									<Select.Item value={option.toString()}>{option} items</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>
				</div>
			</aside>

			<main class="flex-1">
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
							<Button variant="outline" onclick={clearAllFilters} class="mt-4">
								Clear all filters
							</Button>
						{/if}
					</div>
				{/if}

				<Pagination {totalPages} {currentPage} onPageChange={changePage} />
			</main>
		</div>
	</div>
</div>

<style>
	.filter-scroll::-webkit-scrollbar {
		width: 4px;
	}

	.filter-scroll::-webkit-scrollbar-track {
		background: transparent;
	}

	.filter-scroll::-webkit-scrollbar-thumb {
		background-color: hsl(var(--border));
		border-radius: 2px;
	}

	.filter-scroll::-webkit-scrollbar-thumb:hover {
		background-color: hsl(var(--muted-foreground) / 0.3);
	}

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
