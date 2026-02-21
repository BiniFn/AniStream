<script lang="ts">
	import { Calendar, ChevronRight, Clock, Play, RefreshCcw, Star, TrendingUp } from 'lucide-svelte';
	import { resource } from 'runed';
	import { apiClient } from '$lib/api/client';
	import AnimeCard from '$lib/components/anime/display/anime-card.svelte';
	import AnimeSection from '$lib/components/anime/layout/anime-section.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { cn } from '$lib/utils';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const isLoggedIn = $derived(data.isLoggedIn);

	// Fetch homepage data client-side using runed resource
	const homeResource = resource(
		() => null,
		async (_, __, { signal }) => {
			const response = await apiClient.GET('/home', { signal });
			if (response.error || !response.data) {
				throw response.error || new Error('Failed to fetch homepage data');
			}
			return response.data;
		},
		{ once: true },
	);

	// Derived data from resource
	const homeData = $derived(homeResource.current);
	const hasError = $derived(!!homeResource.error);
	const fallbackHomeData = {
		trending: [
			{
				id: 16498,
				ename: 'Attack on Titan',
				jname: 'Shingeki no Kyojin',
				imageUrl:
					'https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx16498-buvcRTBx8iAt.jpg',
				season: 'spring',
				seasonYear: 2013,
				genre: 'Action, Drama, Mystery',
				score: 8.5,
			},
		],
		popular: [],
		recentlyUpdated: [],
		seasonal: [],
		featuredAnime: {
			id: 21,
			ename: 'One Piece',
			jname: 'One Piece',
			synopsis:
				"Monkey D. Luffy and his crew sail the Grand Line in search of the legendary treasure, the One Piece.",
			imageUrl:
				'https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx21-YCDoj1EkAxFn.jpg',
			season: 'fall',
			seasonYear: 1999,
			genre: 'Action, Adventure, Comedy, Fantasy',
			type: 'TV',
			lastEpisode: 1100,
			score: 8.7,
		},
		continueWatching: [],
		planning: [],
	};
	const pageData = $derived(hasError ? fallbackHomeData : homeData);

	// Extract data with fallbacks
	const trending = $derived(pageData?.trending || []);
	const popular = $derived(pageData?.popular || []);
	const recentlyUpdated = $derived(pageData?.recentlyUpdated || []);
	const seasonal = $derived(pageData?.seasonal || []);
	const featuredAnime = $derived(pageData?.featuredAnime || null);
	const continueWatching = $derived(pageData?.continueWatching || []);
	const planning = $derived(pageData?.planning || []);

	const popularAnime = $derived(popular.slice(0, 8));
</script>

<svelte:head>
	<title>Aniways</title>
	<meta
		name="description"
		content="Discover, watch, and track your favorite anime series and movies. Stay updated with the latest releases and trending shows."
	/>
</svelte:head>

{#if !homeData && !hasError}
	<!-- Full Page Skeleton - Shows immediately while data loads -->
	<div class="min-h-screen">
		<!-- Hero Skeleton - Takes up full viewport height naturally -->
		<section class="relative h-screen w-full overflow-hidden">
			<div
				class="absolute inset-0 animate-pulse bg-gradient-to-b from-muted/30 via-muted/50 to-muted"
			></div>

			<!-- Hero content skeleton -->
			<div class="relative z-10 flex h-full items-center">
				<div class="container mx-auto px-6">
					<div
						class="flex flex-col-reverse items-center gap-4 md:gap-8 lg:flex-row lg:justify-between lg:gap-16 xl:gap-20"
					>
						<div class="max-w-3xl flex-1 space-y-6">
							<!-- Trending badge skeleton -->
							<div class="mb-6 hidden md:block">
								<Skeleton class="h-10 w-32 rounded-full" />
							</div>

							<!-- Title skeleton -->
							<Skeleton class="h-16 w-3/4 sm:h-20 md:h-24 lg:h-32" />

							<!-- Subtitle skeleton -->
							<Skeleton class="hidden h-8 w-1/2 sm:block md:h-10" />

							<!-- Season badge skeleton -->
							<Skeleton class="hidden h-10 w-40 rounded-full sm:block" />

							<!-- Genre buttons skeleton -->
							<div class="flex flex-wrap gap-2">
								{#each Array(4) as _, i (i)}
									<Skeleton class="h-8 w-20 rounded-full" />
								{/each}
							</div>

							<!-- Action buttons skeleton -->
							<div class="flex flex-col gap-3 sm:flex-row sm:gap-4">
								<Skeleton class="h-12 w-40" />
								<Skeleton class="h-12 w-32" />
							</div>
						</div>

						<!-- Featured image skeleton -->
						<div class="lg:flex-shrink-0">
							<Skeleton
								class="aspect-[2/3] w-[45vw] rounded-2xl md:w-[300px] lg:h-[550px] lg:w-96"
							/>
						</div>
					</div>
				</div>
			</div>
		</section>

		<!-- Content Skeletons - Now naturally follows the hero -->
		<div class="container mx-auto space-y-12 px-4 py-12">
			<!-- Trending Skeleton -->
			<div class="space-y-4">
				<Skeleton class="h-8 w-40" />
				<div
					class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
				>
					{#each Array(6) as _, i (i)}
						<div class="space-y-3">
							<div
								class="relative overflow-hidden rounded-xl bg-gradient-to-br from-muted to-muted/50 shadow-lg"
							>
								<Skeleton class="aspect-[3/4] w-full" />
								<div class="absolute top-3 left-3">
									<Skeleton class="h-6 w-12 rounded-md" />
								</div>
							</div>
							<Skeleton class="h-4 w-full" />
							<Skeleton class="h-3 w-2/3" />
						</div>
					{/each}
				</div>
			</div>

			<!-- Continue Watching Skeleton (only if logged in) -->
			{#if isLoggedIn}
				<div class="space-y-4">
					<Skeleton class="h-8 w-48" />
					<div
						class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
					>
						{#each Array(6) as _, i (i)}
							<div class="space-y-3">
								<div
									class="relative overflow-hidden rounded-xl bg-gradient-to-br from-muted to-muted/50 shadow-lg"
								>
									<Skeleton class="aspect-[3/4] w-full" />
								</div>
								<Skeleton class="h-4 w-full" />
								<Skeleton class="h-3 w-2/3" />
							</div>
						{/each}
					</div>
				</div>

				<!-- Planning Skeleton -->
				<div class="space-y-4">
					<Skeleton class="h-8 w-32" />
					<div
						class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
					>
						{#each Array(6) as _, i (i)}
							<div class="space-y-3">
								<div
									class="relative overflow-hidden rounded-xl bg-gradient-to-br from-muted to-muted/50 shadow-lg"
								>
									<Skeleton class="aspect-[3/4] w-full" />
								</div>
								<Skeleton class="h-4 w-full" />
								<Skeleton class="h-3 w-2/3" />
							</div>
						{/each}
					</div>
				</div>
			{/if}

			<!-- Popular Skeleton -->
			<div class="space-y-4">
				<Skeleton class="h-8 w-36" />
				<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
					{#each Array(8) as _, i (i)}
						<div
							class="flex gap-4 rounded-xl border bg-card/50 p-5 transition-all hover:border-primary/20"
						>
							<div class="relative overflow-hidden rounded-lg shadow-md">
								<Skeleton class="aspect-[3/4] w-20" />
								<div class="absolute top-1 right-1">
									<Skeleton class="h-5 w-8 rounded" />
								</div>
							</div>
							<div class="flex-1 space-y-3">
								<Skeleton class="h-4 w-3/4" />
								<Skeleton class="h-3 w-1/2" />
								<Skeleton class="h-3 w-1/3" />
							</div>
						</div>
					{/each}
				</div>
			</div>

			<!-- Seasonal Skeleton -->
			<div class="space-y-4">
				<Skeleton class="h-8 w-28" />
				<div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
					{#each Array(6) as _, i (i)}
						<div
							class="overflow-hidden rounded-2xl bg-gradient-to-br from-muted to-muted/50 shadow-lg"
						>
							<div class="relative">
								<Skeleton class="h-48 w-full" />
								<div class="absolute top-4 right-4">
									<Skeleton class="h-8 w-28 rounded-md" />
								</div>
							</div>
							<div class="space-y-4 p-6">
								<div class="flex gap-4">
									<div class="relative overflow-hidden rounded-lg shadow-md">
										<Skeleton class="h-24 w-16" />
									</div>
									<div class="flex-1 space-y-3">
										<Skeleton class="h-4 w-3/4" />
										<Skeleton class="h-3 w-1/2" />
										<div class="flex flex-wrap gap-1">
											<Skeleton class="h-6 w-16 rounded-full" />
											<Skeleton class="h-6 w-16 rounded-full" />
										</div>
									</div>
								</div>
								<Skeleton class="h-20 w-full" />
							</div>
						</div>
					{/each}
				</div>
			</div>

			<!-- Recently Updated Skeleton -->
			<div class="space-y-4">
				<Skeleton class="h-8 w-44" />
				<div
					class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
				>
					{#each Array(6) as _, i (i)}
						<div class="space-y-3">
							<div
								class="relative overflow-hidden rounded-xl bg-gradient-to-br from-muted to-muted/50 shadow-lg"
							>
								<Skeleton class="aspect-[3/4] w-full" />
								<div class="absolute top-3 left-3">
									<Skeleton class="h-6 w-10 rounded-full" />
								</div>
							</div>
							<Skeleton class="h-4 w-full" />
							<Skeleton class="h-3 w-2/3" />
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>

{:else}
	{#if hasError}
		<div class="container mx-auto px-4 pt-4">
			<div
				class="flex flex-col gap-3 rounded-lg border border-amber-500/30 bg-amber-500/10 p-4 text-sm text-amber-100 md:flex-row md:items-center md:justify-between"
			>
				<p>Live data is unavailable right now, so you're seeing demo content.</p>
				<Button onclick={() => homeResource.refetch()} size="sm" variant="outline">
					<RefreshCcw class="mr-2 h-4 w-4" />
					Retry connection
				</Button>
			</div>
		</div>
	{/if}

	<!-- Real Content - Only renders when data is ready -->
	{#if featuredAnime}
		<section class="absolute top-0 right-0 left-0 z-10 mb-16 h-screen w-screen overflow-hidden">
			<div class="absolute inset-0">
				<img
					src={featuredAnime.metadata?.mainPictureUrl || featuredAnime.imageUrl}
					alt={featuredAnime.ename || featuredAnime.jname}
					class="absolute inset-0 h-full w-full scale-110 object-cover blur-sm"
				/>
				<div class="absolute inset-0 bg-gradient-to-r from-black/80 via-black/50 to-black/20"></div>
				<div
					class="absolute inset-0 bg-gradient-to-b from-black/60 via-transparent to-transparent md:bg-gradient-to-t"
				></div>
			</div>

			<div class="absolute inset-0 z-20 flex items-center">
				<div class="container mx-auto px-6">
					<div
						class="flex flex-col-reverse items-center gap-4 md:gap-8 lg:flex-row lg:justify-between lg:gap-16 xl:gap-20"
					>
						<div class="max-w-3xl flex-1">
							<div class="mb-6 flex items-center gap-3 md:pt-0">
								<div
									class="hidden items-center gap-2 rounded-full border border-primary/50 bg-primary/20 px-4 py-2 backdrop-blur-xl md:flex"
								>
									<TrendingUp class="h-4 w-4 text-primary" />
									<span class="text-sm font-bold tracking-wider text-primary uppercase">
										Trending Now
									</span>
								</div>
							</div>
							<h1
								class="mb-3 line-clamp-2 text-3xl leading-tight font-black tracking-tight text-foreground sm:mb-6 sm:text-5xl md:text-5xl lg:text-6xl xl:text-7xl"
							>
								{featuredAnime.jname || featuredAnime.ename}
							</h1>
							{#if featuredAnime.ename && featuredAnime.jname}
								<p
									class="mb-4 hidden text-xl font-light text-muted-foreground sm:mb-6 sm:block md:text-2xl lg:text-3xl"
								>
									{featuredAnime.ename}
								</p>
							{/if}
							<a
								href="/catalog?seasons={featuredAnime.season}&years={featuredAnime.seasonYear}"
								class="mb-3 hidden w-fit items-center gap-2 rounded-full bg-white/10 px-4 py-2 text-gray-300 backdrop-blur-sm transition-colors hover:bg-white/20 hover:text-white sm:mb-6 sm:flex"
							>
								<Calendar class="h-4 w-4" />
								<span class="font-medium capitalize">
									{featuredAnime.season}
									{featuredAnime.seasonYear}
								</span>
							</a>
							<div class="mb-4 flex flex-wrap items-center gap-6 sm:mb-6">
								<div class="flex flex-wrap gap-2">
									{#each featuredAnime.genre.split(', ') as genre, i (genre)}
										<Button
											size="sm"
											class={cn(
												'rounded-full text-xs md:text-sm',
												i > 3 ? 'hidden md:inline-flex' : '',
											)}
											variant="outline"
											href="/catalog?genres={genre}"
										>
											{genre}
										</Button>
									{/each}
								</div>
							</div>
							<div class="mt-6 flex flex-col gap-3 sm:mt-0 sm:flex-row sm:gap-4">
								<Button
									class="gap-3 shadow-xl shadow-primary/25 hover:bg-primary/90"
									href={`/anime/${featuredAnime.id}/watch`}
								>
									<Play class="h-6 w-6" />
									Watch Now
								</Button>
								<Button
									variant="outline"
									class="gap-3 backdrop-blur-lg"
									href={`/anime/${featuredAnime.id}`}
								>
									More Info
									<ChevronRight class="h-5 w-5" />
								</Button>
							</div>
						</div>

						<div class="lg:flex-shrink-0">
							<div class="relative">
								<div class="relative overflow-hidden rounded-2xl pt-12 shadow-2xl md:pt-0">
									<img
										src={featuredAnime.metadata?.mainPictureUrl || featuredAnime.imageUrl}
										alt={featuredAnime.ename || featuredAnime.jname}
										class="aspect-[2/3] w-[45vw] overflow-hidden rounded-2xl object-cover md:w-[300px] lg:h-[550px] lg:w-96"
									/>
								</div>
								<div
									class="absolute -right-4 -bottom-4 rounded-xl border border-white/20 bg-black/80 px-4 py-3 backdrop-blur-sm"
								>
									<div class="flex items-center gap-2">
										<Star class="h-4 w-4 fill-yellow-400 text-yellow-400" />
										<span class="font-semibold">#1 Trending</span>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>

			<div class="absolute bottom-2 left-1/2 z-30 -translate-x-1/2 transform sm:bottom-8">
				<div class="animate-bounce">
					<ChevronRight class="size-6 rotate-90 text-muted-foreground" />
				</div>
			</div>
		</section>
	{/if}

	<div class="container mx-auto mt-[100vh] space-y-12 px-4">
		<AnimeSection icon={TrendingUp} title="Trending Anime" visible={trending.length > 0}>
			{#each trending.slice(1, 7) as anime, index (anime.id)}
				<AnimeCard {anime} class="w-40 flex-shrink-0 md:w-auto">
					{#snippet topLeftBadge()}
						<div
							class="flex items-center gap-1 rounded-md bg-background/90 px-2 py-1 text-xs font-semibold text-primary-foreground backdrop-blur-sm"
						>
							<Star class="h-3 w-3 fill-yellow-400 text-yellow-400" />
							<span>#{index + 2}</span>
						</div>
					{/snippet}
				</AnimeCard>
			{/each}
		</AnimeSection>

		<AnimeSection
			icon={Play}
			title="Continue Watching"
			viewAllHref="/continue-watching"
			visible={isLoggedIn && continueWatching.length > 0}
		>
			{#each continueWatching as item (item.id)}
				<AnimeCard
					anime={item.anime}
					episodeLink={item.watchedEpisodes + 1}
					class="w-40 flex-shrink-0 md:w-auto"
				/>
			{/each}
		</AnimeSection>

		<AnimeSection
			icon={Clock}
			title="Planning"
			viewAllHref="/planning"
			visible={isLoggedIn && planning.length > 0}
		>
			{#each planning as item (item.id)}
				<AnimeCard anime={item.anime} episodeLink={1} class="w-40 flex-shrink-0 md:w-auto" />
			{/each}
		</AnimeSection>

		<AnimeSection
			icon={Star}
			title="Popular Anime"
			class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
			visible={popularAnime.length > 0}
		>
			{#each popularAnime as anime, index (anime.id)}
				<a
					href="/anime/{anime.id}"
					class="group block transform transition-all duration-300 hover:scale-[1.02]"
				>
					<div
						class="flex gap-4 rounded-xl border bg-card/50 p-5 backdrop-blur-sm transition-all duration-300 group-hover:border-primary/20 hover:bg-card hover:shadow-xl hover:shadow-primary/5"
					>
						<div
							class="relative aspect-[3/4] w-20 flex-shrink-0 overflow-hidden rounded-lg bg-gradient-to-br from-muted to-muted/50 shadow-md"
						>
							<img
								src={anime.imageUrl}
								alt={anime.ename || anime.jname}
								class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
							/>
							<div
								class="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100"
							></div>

							<div class="absolute top-1 right-1">
								<div
									class="flex items-center gap-1 rounded-md bg-yellow-500/90 px-1.5 py-0.5 text-xs font-semibold text-white"
								>
									<Star class="h-2.5 w-2.5 fill-current" />
									<span>#{index + 1}</span>
								</div>
							</div>
						</div>
						<div class="min-w-0 flex-1 space-y-3">
							<div>
								<h3
									class="mb-1 line-clamp-1 text-base font-semibold transition-colors duration-300 group-hover:text-primary"
								>
									{anime.ename || anime.jname}
								</h3>
								<div class="text-sm text-muted-foreground capitalize">
									{anime.season}
									{anime.seasonYear}
								</div>
							</div>
							<div class="flex flex-wrap gap-1">
								{#each anime.genre.split(', ').slice(0, 2) as genre (genre)}
									<span
										class="rounded-full border bg-muted/80 px-2 py-1 text-xs text-muted-foreground"
									>
										{genre}
									</span>
								{/each}
							</div>
							{#if anime.lastEpisode}
								<div class="text-xs font-medium text-primary">
									Episode {anime.lastEpisode} Available
								</div>
							{/if}
						</div>
						<div
							class="flex flex-shrink-0 items-center opacity-0 transition-all duration-300 group-hover:opacity-100"
						>
							<ChevronRight class="h-5 w-5 text-primary" />
						</div>
					</div>
				</a>
			{/each}
		</AnimeSection>

		<AnimeSection
			icon={Star}
			title="This Season"
			class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3"
			visible={seasonal.length > 0}
		>
			{#each seasonal.slice(0, 6) as seasonalAnime (seasonalAnime.id)}
				<a
					href="/anime/{seasonalAnime.anime.id}"
					class="group block transform transition-all duration-300 hover:scale-[1.02]"
				>
					<div
						class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-muted to-muted/50 shadow-lg transition-all duration-300 group-hover:shadow-xl group-hover:shadow-primary/10"
					>
						<div class="relative h-48 overflow-hidden">
							<img
								src={seasonalAnime.bannerImageUrl}
								alt={seasonalAnime.anime.ename || seasonalAnime.anime.jname}
								class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
							/>
							<div
								class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent"
							></div>

							<div class="absolute top-4 right-4">
								<span
									class="rounded-md bg-primary/90 px-3 py-1 text-sm font-semibold text-primary-foreground backdrop-blur-sm"
								>
									{seasonalAnime.episodes} Episodes
								</span>
							</div>
						</div>
						<div class="p-6">
							<div class="mb-4 flex items-start gap-4">
								<div class="relative h-24 w-16 flex-shrink-0 overflow-hidden rounded-lg">
									<img
										src={seasonalAnime.anime.imageUrl}
										alt={seasonalAnime.anime.ename || seasonalAnime.anime.jname}
										class="h-full w-full object-cover"
									/>
								</div>

								<div class="min-w-0 flex-1">
									<h3
										class="mb-2 line-clamp-1 text-lg font-bold transition-colors duration-300 group-hover:text-primary"
									>
										{seasonalAnime.anime.ename || seasonalAnime.anime.jname}
									</h3>
									<div class="mb-2 text-sm text-muted-foreground capitalize">
										{seasonalAnime.anime.season}
										{seasonalAnime.anime.seasonYear}
									</div>
									<div class="flex flex-wrap gap-1">
										{#each seasonalAnime.anime.genre.split(', ').slice(0, 2) as genre (genre)}
											<span class="rounded-full bg-muted px-2 py-1 text-xs text-muted-foreground">
												{genre}
											</span>
										{/each}
									</div>
								</div>
							</div>

							<p class="line-clamp-3 text-sm text-muted-foreground">
								<!-- eslint-disable-next-line -->
								{@html seasonalAnime.description?.replace(/<br\s*\/?>/gi, ' ') || ''}
							</p>
						</div>
					</div>
				</a>
			{/each}
		</AnimeSection>

		<AnimeSection
			icon={Calendar}
			title="Recently Updated"
			viewAllHref="/catalog?sortBy=updated_at&sortOrder=desc"
			visible={recentlyUpdated.length > 0}
		>
			{#each recentlyUpdated as anime (anime.id)}
				<AnimeCard
					{anime}
					episodeLink={anime.lastEpisode || 1}
					class="w-40 flex-shrink-0 md:w-auto"
				>
					{#snippet topLeftBadge()}
						<div
							class="animate-pulse rounded-full bg-red-500/90 px-2 py-1 text-xs font-semibold text-white backdrop-blur-sm"
						>
							NEW
						</div>
					{/snippet}
				</AnimeCard>
			{/each}
		</AnimeSection>

		{#if !isLoggedIn}
			<section
				class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-primary/20 via-primary/10 to-secondary/20 p-8 text-center md:p-16"
			>
				<div class="bg-grid-pattern absolute inset-0 opacity-10"></div>
				<div
					class="absolute inset-0 bg-gradient-to-r from-transparent via-white/5 to-transparent"
				></div>
				<div class="relative z-10">
					<div
						class="mb-8 inline-flex items-center gap-2 rounded-full border border-primary/20 bg-primary/10 px-6 py-2"
					>
						<Star class="h-4 w-4 text-primary" />
						<span class="text-sm font-semibold tracking-wide text-primary uppercase"
							>Join the Community</span
						>
					</div>
					<h2
						class="mb-6 bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text text-4xl font-black text-transparent md:text-6xl"
					>
						Ready to Start Your<br />Anime Journey?
					</h2>
					<p class="mx-auto mb-12 max-w-3xl text-xl leading-relaxed text-muted-foreground">
						Join thousands of anime fans worldwide. Create your personalized watchlist, track your
						progress, discover new favorites, and never miss an episode again.
					</p>
					<div class="mb-8 flex flex-col items-center justify-center gap-6 sm:flex-row">
						<Button
							href="/register"
							size="lg"
							class="transform gap-3 bg-primary px-10 py-4 text-lg font-semibold shadow-2xl shadow-primary/25 transition-all duration-300 hover:scale-105 hover:bg-primary/90"
						>
							<Play class="h-5 w-5" />
							Get Started Free
							<ChevronRight class="h-4 w-4" />
						</Button>
						<Button
							href="/catalog"
							variant="outline"
							size="lg"
							class="transform gap-3 border-2 px-10 py-4 text-lg font-semibold transition-all duration-300 hover:scale-105 hover:border-primary hover:text-primary"
						>
							Browse Catalog
						</Button>
					</div>
					<div
						class="flex flex-wrap items-center justify-center gap-8 text-sm text-muted-foreground"
					>
						<div class="flex items-center gap-2">
							<div class="h-2 w-2 rounded-full bg-green-500"></div>
							<span>Free Forever</span>
						</div>
						<div class="flex items-center gap-2">
							<div class="h-2 w-2 rounded-full bg-blue-500"></div>
							<span>No Ads</span>
						</div>
						<div class="flex items-center gap-2">
							<div class="h-2 w-2 rounded-full bg-purple-500"></div>
							<span>HD Quality</span>
						</div>
					</div>
				</div>
			</section>
		{/if}
	</div>
{/if}
