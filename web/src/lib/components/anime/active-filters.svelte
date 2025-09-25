<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Calendar, Hash, Search, Sparkles, X } from 'lucide-svelte';

	interface Props {
		searchQuery?: string;
		selectedGenres: string[];
		selectedSeasons: string[];
		selectedYears: number[];
		totalFilters: number;
		onClearSearch: () => void;
		onRemoveGenre: (genre: string) => void;
		onRemoveSeason: (season: string) => void;
		onRemoveYear: (year: number) => void;
		onClearAll: () => void;
		class?: string;
	}

	let {
		searchQuery,
		selectedGenres,
		selectedSeasons,
		selectedYears,
		totalFilters,
		onClearSearch,
		onRemoveGenre,
		onRemoveSeason,
		onRemoveYear,
		onClearAll,
		class: className = '',
	}: Props = $props();
</script>

{#if totalFilters > 0}
	<div class="flex flex-wrap items-center gap-2 {className}">
		<span class="text-sm font-medium">Active filters:</span>

		{#if searchQuery}
			<Badge variant="secondary" class="gap-1">
				<Search class="h-3 w-3" />
				{searchQuery}
				<button onclick={onClearSearch} class="ml-1 hover:text-destructive">
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/if}

		{#each selectedGenres as genre (genre)}
			<Badge variant="secondary" class="gap-1">
				<Hash class="h-3 w-3" />
				{genre}
				<button onclick={() => onRemoveGenre(genre)} class="ml-1 hover:text-destructive">
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/each}

		{#each selectedSeasons as season (season)}
			<Badge variant="secondary" class="gap-1 capitalize">
				<Sparkles class="h-3 w-3" />
				{season}
				<button onclick={() => onRemoveSeason(season)} class="ml-1 hover:text-destructive">
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/each}

		{#each selectedYears as year (year)}
			<Badge variant="secondary" class="gap-1">
				<Calendar class="h-3 w-3" />
				{year}
				<button onclick={() => onRemoveYear(year)} class="ml-1 hover:text-destructive">
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/each}

		<Button variant="ghost" size="sm" onclick={onClearAll}>Clear all</Button>
	</div>
{/if}
