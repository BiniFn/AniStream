<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Calendar, Hash, Search, Sparkles, X } from 'lucide-svelte';
	import type { FilterManager } from '$lib/utils/filter-manager.svelte';

	interface Props {
		filterManager: FilterManager;
		class?: string;
	}

	let { filterManager, class: className = '' }: Props = $props();
</script>

{#if filterManager.totalFilters > 0}
	<div class="flex flex-wrap items-center gap-2 {className}">
		<span class="text-sm font-medium">Active filters:</span>

		{#if filterManager.filters.search}
			<Badge variant="secondary" class="gap-1">
				<Search class="h-3 w-3" />
				{filterManager.filters.search}
				<button onclick={filterManager.clearSearch} class="ml-1 hover:text-destructive">
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/if}

		{#each filterManager.filters.genres as genre (genre)}
			<Badge variant="secondary" class="gap-1">
				<Hash class="h-3 w-3" />
				{genre}
				<button
					onclick={() => filterManager.toggleGenre(genre)}
					class="ml-1 hover:text-destructive"
				>
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/each}

		{#each filterManager.filters.seasons as season (season)}
			<Badge variant="secondary" class="gap-1 capitalize">
				<Sparkles class="h-3 w-3" />
				{season}
				<button
					onclick={() => filterManager.toggleSeason(season)}
					class="ml-1 hover:text-destructive"
				>
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/each}

		{#each filterManager.filters.years as year (year)}
			<Badge variant="secondary" class="gap-1">
				<Calendar class="h-3 w-3" />
				{year}
				<button onclick={() => filterManager.toggleYear(year)} class="ml-1 hover:text-destructive">
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/each}

		<Button variant="ghost" size="sm" onclick={filterManager.clearAll}>Clear all</Button>
	</div>
{/if}
