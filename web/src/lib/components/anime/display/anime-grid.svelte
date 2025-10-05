<script lang="ts">
	import AnimeCard from './anime-card.svelte';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import type { components } from '$lib/api/openapi';
	import { Badge } from '$lib/components/ui/badge';
	import { Pagination } from '$lib/components/ui/pagination';
	import type { Snippet } from 'svelte';
	import LibraryBtn from '../controls/library-btn.svelte';

	type AnimeResponse = components['schemas']['models.AnimeResponse'];
	type AnimeWithLibraryResponse = components['schemas']['models.AnimeWithLibraryResponse'];

	interface Props {
		anime: AnimeResponse[] | AnimeWithLibraryResponse[];
		viewMode: 'grid' | 'list';
		currentPage: number;
		totalPages: number;
		changePage: (page: number) => void;
		isLoading: boolean;
		itemsPerPage: number;
		children?: Snippet;
		empty: Snippet;
	}

	let {
		anime,
		viewMode,
		currentPage,
		totalPages,
		changePage,
		isLoading = false,
		itemsPerPage = 24,
		children,
		empty,
	}: Props = $props();
</script>

<main class="flex-1">
	{#if children}
		{@render children()}
	{/if}

	<div class="px-4">
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
		{:else if anime && anime.length > 0}
			{#if viewMode === 'grid'}
				<div
					class="grid animate-in grid-cols-2 gap-3 duration-500 fade-in sm:grid-cols-3 sm:gap-4 md:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6"
				>
					{#each anime as a, i (a.id)}
						<div
							class="animate-in duration-500 slide-in-from-bottom-5"
							style="animation-delay: {i * 50}ms"
						>
							<AnimeCard
								anime={a}
								index={i}
								libraryEntry={'library' in a && a.library
									? {
											id: a.library.id,
											anime: a,
											status: a.library.status,
											animeId: a.id,
											createdAt: '',
											updatedAt: a.library.updatedAt,
											userId: '',
											watchedEpisodes: a.library.watchedEpisodes,
										}
									: undefined}
							/>
						</div>
					{/each}
				</div>
			{:else}
				<div class="animate-in space-y-4 duration-500 fade-in">
					{#each anime as a, i (a.id)}
						<a
							href="/anime/{a.id}"
							class="flex animate-in gap-3 rounded-lg border bg-card p-3 transition-all duration-500 slide-in-from-left-5 hover:border-primary hover:shadow-lg sm:p-4"
							style="animation-delay: {i * 50}ms"
							onclick={(e) => {
								if ((e.target as HTMLElement).closest('button')) {
									e.preventDefault();
								}
							}}
						>
							<img
								src={a.imageUrl}
								alt={a.ename || a.jname}
								class="h-24 w-20 rounded-md object-cover sm:h-32 sm:w-24"
							/>
							<div class="flex-1 space-y-2">
								<div>
									<h3 class="line-clamp-2 text-sm font-semibold sm:text-lg">
										{a.ename || a.jname}
									</h3>
									{#if a.ename && a.jname}
										<p class="text-sm text-muted-foreground">{a.jname}</p>
									{/if}
								</div>
								<div class="flex flex-wrap gap-2">
									{#each a.genre.split(', ').slice(0, 5) as genre (genre)}
										<Badge variant="secondary" class="text-xs">{genre}</Badge>
									{/each}
								</div>
								<div class="flex gap-4 text-sm text-muted-foreground">
									<span class="capitalize">{a.season} {a.seasonYear}</span>
									{#if a.lastEpisode}
										<span>Episode {a.lastEpisode}</span>
									{/if}
								</div>
								{#if 'library' in a && a.library}
									<LibraryBtn
										animeId={a.id}
										libraryEntry={{
											id: a.library.id,
											anime: a,
											status: a.library.status,
											animeId: a.id,
											createdAt: '',
											updatedAt: a.library.updatedAt,
											userId: '',
											watchedEpisodes: a.library.watchedEpisodes,
										}}
									/>
								{/if}
							</div>
						</a>
					{/each}
				</div>
			{/if}

			<Pagination {totalPages} {currentPage} onPageChange={changePage} />
		{:else}
			{@render empty()}
		{/if}
	</div>
</main>

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
