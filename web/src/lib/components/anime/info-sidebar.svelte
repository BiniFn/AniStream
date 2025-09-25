<script lang="ts">
	import type { components } from '$lib/api/openapi';
	import Button from '$lib/components/ui/button/button.svelte';
	import { cn } from '$lib/utils';
	import { formatDate } from 'date-fns';
	import { ExternalLink, TrendingUp } from 'lucide-svelte';

	type AnimeResponse = components['schemas']['models.AnimeWithMetadataResponse'];
	interface Props {
		anime: AnimeResponse;
		ratingLabel: string;
		selectedTab?: string;
	}

	let { anime, ratingLabel, selectedTab }: Props = $props();
	let mediaType = $derived.by(() => {
		const type = anime.metadata?.mediaType || 'tv';
		return type === 'tv' ? 'TV' : type.charAt(0).toUpperCase() + type.slice(1).toLowerCase();
	});
</script>

<div class={cn('sticky top-36 h-fit space-y-6', selectedTab !== 'overview' && 'hidden md:block')}>
	<div class="rounded-xl border bg-card p-6 shadow-sm">
		<h3 class="mb-4 text-lg font-bold">Information</h3>
		<dl class="space-y-3 text-sm">
			<div class="border-b pb-3">
				<dt class="mb-2 font-semibold text-muted-foreground">Media</dt>
				<div class="space-y-2">
					{#if anime.metadata?.mediaType}
						<div class="flex justify-between">
							<span class="text-muted-foreground">Type</span>
							<span class="font-medium">{mediaType}</span>
						</div>
					{/if}
					{#if anime.metadata?.totalEpisodes}
						<div class="flex justify-between">
							<span class="text-muted-foreground">Episodes</span>
							<span class="font-medium">{anime.metadata.totalEpisodes}</span>
						</div>
					{/if}
					{#if anime.metadata?.avgEpisodeDuration}
						<div class="flex justify-between">
							<span class="text-muted-foreground">Duration</span>
							<span class="font-medium">
								{Math.ceil(anime.metadata.avgEpisodeDuration / 60)} min
							</span>
						</div>
					{/if}
					{#if anime.metadata?.airingStatus}
						<div class="flex justify-between">
							<span class="text-muted-foreground">Status</span>
							<span class={cn('font-medium')}>
								{anime.metadata.airingStatus
									.replace(/_/g, ' ')
									.replace(/\b\w/g, (l) => l.toUpperCase())}
							</span>
						</div>
					{/if}
				</div>
			</div>

			<div class="border-b pb-3">
				<dt class="mb-2 font-semibold text-muted-foreground">Airing</dt>
				<div class="space-y-2">
					{#if anime.season && anime.seasonYear}
						<div class="flex justify-between">
							<span class="text-muted-foreground">Season</span>
							<span class="font-medium capitalize">{anime.season} {anime.seasonYear}</span>
						</div>
					{/if}
					{#if anime.metadata?.airingStartDate}
						<div class="flex justify-between">
							<span class="text-muted-foreground">Start Date</span>
							<span class="font-medium">
								{formatDate(anime.metadata.airingStartDate, 'dd/MM/yyyy')}
							</span>
						</div>
					{/if}
					{#if anime.metadata?.airingEndDate}
						<div class="flex justify-between">
							<span class="text-muted-foreground">End Date</span>
							<span class="font-medium">
								{#if !anime.metadata.airingEndDate}
									???
								{:else}
									{formatDate(anime.metadata.airingEndDate, 'dd/MM/yyyy')}
								{/if}
							</span>
						</div>
					{/if}
				</div>
			</div>

			<div class="border-b pb-3">
				<dt class="mb-2 font-semibold text-muted-foreground">Production</dt>
				<div class="space-y-2">
					{#if anime.metadata?.studio}
						<div class="flex justify-between">
							<span class="text-muted-foreground">Studio</span>
							<span class="font-medium capitalize">{anime.metadata.studio}</span>
						</div>
					{/if}
					{#if anime.metadata?.source}
						<div class="flex justify-between">
							<span class="text-muted-foreground">Source</span>
							<span class="font-medium capitalize">{anime.metadata.source}</span>
						</div>
					{/if}
					{#if anime.metadata?.rating}
						<div class="flex justify-between">
							<span class="text-muted-foreground">Rating</span>
							<span class="font-medium">{ratingLabel}</span>
						</div>
					{/if}
				</div>
			</div>

			<div class="border-b pb-3 md:hidden">
				<dt class="mb-2 font-semibold text-muted-foreground">Synopsis</dt>
				<div class="space-y-2">
					{#if anime.metadata?.description}
						<p class="text-sm text-muted-foreground">{anime.metadata.description}</p>
					{:else}
						<p class="text-sm text-muted-foreground italic">No synopsis available.</p>
					{/if}
				</div>
			</div>

			{#if anime.malId || anime.anilistId}
				<div>
					<dt class="mb-2 font-semibold text-muted-foreground">External Links</dt>
					<div class="flex gap-2">
						{#if anime.malId}
							<Button
								size="sm"
								variant="outline"
								href={`https://myanimelist.net/anime/${anime.malId}`}
								target="_blank"
								class="gap-2"
							>
								<ExternalLink class="h-3 w-3" />
								MAL
							</Button>
						{/if}
						{#if anime.anilistId}
							<Button
								size="sm"
								variant="outline"
								href={`https://anilist.co/anime/${anime.anilistId}`}
								target="_blank"
								class="gap-2"
							>
								<ExternalLink class="h-3 w-3" />
								AniList
							</Button>
						{/if}
					</div>
				</div>
			{/if}
		</dl>
	</div>

	<div class="rounded-xl border bg-card p-6 shadow-sm">
		<h3 class="mb-4 text-lg font-bold">You Might Also Like</h3>
		<div class="text-center text-sm text-muted-foreground">
			<TrendingUp class="mx-auto mb-2 h-8 w-8" />
			<p>Recommendations coming soon</p>
		</div>
	</div>
</div>
