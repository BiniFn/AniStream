<script lang="ts">
	import { browser } from '$app/environment';
	import { getUser } from '$lib/auth.svelte';
	import AnimeCard from '$lib/components/anime-card.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Calendar, ChevronRight, CircleAlert, Play, Star, TrendingUp } from 'lucide-svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const { trending, popular, recentlyUpdated, seasonal, featuredAnime } = data;
	const user = getUser();

	const hasData =
		trending.length > 0 || popular.length > 0 || recentlyUpdated.length > 0 || seasonal.length > 0;

	let popularAnime = $state(popular);

	function updatePopularAnime() {
		if (!browser) return;

		const width = window.innerWidth;
		if (width < 640) {
			popularAnime = popular.slice(0, 3);
			return;
		}
		if (width < 1024) {
			popularAnime = popular.slice(0, 4);
			return;
		}
		if (width < 1280) {
			popularAnime = popular.slice(0, 6);
			return;
		}
		popularAnime = popular.slice(0, 8);
	}

	$effect(() => {
		updatePopularAnime();
		if (!browser) return;

		window.addEventListener('resize', updatePopularAnime);
		return () => window.removeEventListener('resize', updatePopularAnime);
	});
</script>

{#if featuredAnime}
	<section class="absolute top-0 right-0 left-0 z-10 mb-16 h-screen w-screen overflow-hidden">
		<div class="absolute inset-0">
			<img
				src={featuredAnime.imageUrl}
				alt={featuredAnime.ename || featuredAnime.jname}
				class="absolute inset-0 h-full w-full scale-110 object-cover blur-sm"
			/>
			<div class="absolute inset-0 bg-gradient-to-r from-black/70 via-black/40 to-black/20"></div>
			<div
				class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent"
			></div>
		</div>

		<div class="absolute inset-0 z-20 flex items-center">
			<div class="container mx-auto px-6">
				<div
					class="flex flex-col items-center gap-8 lg:flex-row lg:justify-between lg:gap-16 xl:gap-20"
				>
					<div class="max-w-3xl flex-1">
						<div class="mb-6 flex items-center gap-3 pt-20 md:pt-0">
							<div
								class="flex items-center gap-2 rounded-full border border-primary/30 bg-primary/20 px-4 py-2 backdrop-blur-sm"
							>
								<TrendingUp class="h-4 w-4 text-primary" />
								<span class="text-sm font-semibold tracking-wider text-primary uppercase">
									Trending Now
								</span>
							</div>
						</div>
						<h1
							class="mb-4 text-4xl leading-none font-black tracking-tight text-white sm:mb-6 sm:text-6xl md:text-8xl xl:text-9xl"
						>
							{featuredAnime.ename || featuredAnime.jname}
						</h1>
						{#if featuredAnime.jname && featuredAnime.ename}
							<p class="mb-4 hidden text-2xl font-light text-gray-300 sm:mb-6 sm:block">
								{featuredAnime.jname}
							</p>
						{/if}
						<div class="mb-4 flex flex-wrap items-center gap-6 sm:mb-8">
							<div
								class="flex items-center gap-2 rounded-full bg-white/10 px-4 py-2 text-gray-300 backdrop-blur-sm"
							>
								<Calendar class="h-4 w-4" />
								<span class="font-medium capitalize"
									>{featuredAnime.season} {featuredAnime.seasonYear}</span
								>
							</div>
							<div class="flex flex-wrap gap-2">
								{#each featuredAnime.genre.split(', ').slice(0, 3) as genre (genre)}
									<span
										class="rounded-full border border-white/20 bg-white/10 px-3 py-1 text-sm text-white backdrop-blur-sm"
									>
										{genre}
									</span>
								{/each}
							</div>
						</div>
						<div class="flex flex-col gap-4 sm:flex-row">
							<Button
								size="lg"
								class="gap-3 bg-primary px-8 py-4 text-lg shadow-xl shadow-primary/25 hover:bg-primary/90"
							>
								<Play class="h-6 w-6" />
								Watch Now
							</Button>
							<Button
								variant="outline"
								size="lg"
								class="gap-3 border-white/30 bg-white/10 px-8 py-4 text-lg text-white shadow-xl backdrop-blur-sm hover:bg-white/20"
							>
								More Info
								<ChevronRight class="h-5 w-5" />
							</Button>
						</div>
					</div>

					<div class="lg:flex-shrink-0">
						<div class="relative">
							<div class="relative overflow-hidden rounded-2xl shadow-2xl">
								<img
									src={featuredAnime.imageUrl}
									alt={featuredAnime.ename || featuredAnime.jname}
									class="h-96 w-72 rounded-2xl object-cover lg:h-[550px] lg:w-96 xl:h-[600px] xl:w-[400px]"
								/>
							</div>
							<div
								class="absolute -right-4 -bottom-4 rounded-xl border border-white/20 bg-black/80 px-4 py-3 backdrop-blur-sm"
							>
								<div class="flex items-center gap-2 text-white">
									<Star class="h-4 w-4 fill-yellow-400 text-yellow-400" />
									<span class="font-semibold">#1 Trending</span>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>

		<div class="absolute bottom-4 left-1/2 z-30 -translate-x-1/2 transform sm:bottom-8">
			<div class="animate-bounce">
				<ChevronRight class="h-6 w-6 rotate-90 text-white/60" />
			</div>
		</div>
	</section>
{/if}

{#if !hasData}
	<div class="mt-screen container mx-auto px-4 py-20">
		<div class="mx-auto max-w-md text-center">
			<div class="mb-6">
				<CircleAlert class="mx-auto mb-4 h-16 w-16 text-muted-foreground" />
				<h2 class="mb-2 text-2xl font-bold">No Content Available</h2>
				<p class="text-muted-foreground">
					We're having trouble loading anime data right now. Please try again later.
				</p>
			</div>
			<Button onclick={() => window.location.reload()} class="gap-2">
				<ChevronRight class="h-4 w-4 rotate-180" />
				Try Again
			</Button>
		</div>
	</div>
{:else}
	<div class="container mx-auto mt-[100vh] space-y-12 px-4">
		{#if trending.length > 0}
			<section>
				<div class="mb-8 flex items-center gap-3">
					<TrendingUp class="h-6 w-6 text-primary" />
					<h2 class="text-2xl font-bold sm:text-3xl">Trending Anime</h2>
				</div>
				<div
					class="flex gap-4 overflow-x-auto pb-4 md:grid md:grid-cols-4 md:gap-6 md:overflow-visible md:pb-0 lg:grid-cols-6"
				>
					{#each trending.slice(1, 7) as anime, index (anime.id)}
						<div class="w-40 flex-shrink-0 md:w-auto">
							<AnimeCard {anime} {index}>
								{#snippet topLeftBadge()}
									<div
										class="flex items-center gap-1 rounded-md bg-primary/90 px-2 py-1 text-xs font-semibold text-primary-foreground backdrop-blur-sm"
									>
										<Star class="h-3 w-3" />
										<span>#{index + 2}</span>
									</div>
								{/snippet}
							</AnimeCard>
						</div>
					{/each}
				</div>
			</section>
		{/if}

		{#if popularAnime.length > 0}
			<section>
				<div class="mb-8 flex items-center gap-3">
					<Star class="h-6 w-6 text-primary" />
					<h2 class="text-2xl font-bold sm:text-3xl">Popular Anime</h2>
				</div>
				<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
					{#each popularAnime as anime, index (anime.id)}
						<a
							href="/anime/{anime.id}"
							class="group block transform transition-all duration-300 hover:scale-[1.02]"
							style="animation-delay: {index * 50}ms"
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
				</div>
			</section>
		{/if}

		{#if recentlyUpdated.length > 0}
			<section>
				<div class="mb-8 flex items-center justify-between">
					<div class="flex items-center gap-3">
						<Calendar class="h-6 w-6 text-primary" />
						<h2 class="text-2xl font-bold sm:text-3xl">Recently Updated</h2>
					</div>
					<Button variant="ghost" class="gap-2" href="/recent">
						View All
						<ChevronRight class="h-4 w-4" />
					</Button>
				</div>
				<div class="grid grid-cols-2 gap-6 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
					{#each recentlyUpdated as anime, index (anime.id)}
						<AnimeCard {anime} {index}>
							{#snippet topLeftBadge()}
								<div
									class="animate-pulse rounded-full bg-red-500/90 px-2 py-1 text-xs font-semibold text-white backdrop-blur-sm"
								>
									NEW
								</div>
							{/snippet}
						</AnimeCard>
					{/each}
				</div>
			</section>
		{/if}

		{#if seasonal.length > 0}
			<section>
				<div class="mb-8 flex items-center gap-3">
					<Star class="h-6 w-6 text-primary" />
					<h2 class="text-2xl font-bold sm:text-3xl">This Season</h2>
				</div>
				<div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
					{#each seasonal.slice(0, 6) as seasonalAnime, index (seasonalAnime.id)}
						<a
							href="/anime/{seasonalAnime.anime.id}"
							class="group block transform transition-all duration-300 hover:scale-[1.02]"
							style="animation-delay: {index * 100}ms"
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

									<div class="absolute top-4 left-4">
										<span
											class="rounded-full bg-black/60 px-3 py-1 text-sm font-medium text-white backdrop-blur-sm"
										>
											{seasonalAnime.type}
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
													<span
														class="rounded-full bg-muted px-2 py-1 text-xs text-muted-foreground"
													>
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
				</div>
			</section>
		{/if}

		{#if !user}
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
