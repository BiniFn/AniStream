<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Calendar, Hash, Search, Sparkles, X } from 'lucide-svelte';
	import type { FilterState, FilterActions } from '$lib/utils/filters';

	interface Props {
		filters: FilterState;
		filterActions: FilterActions;
		totalFilters: number;
		class?: string;
	}

	let { filters, filterActions, totalFilters, class: className = '' }: Props = $props();
</script>

{#if totalFilters > 0}
	<div class="flex flex-wrap items-center gap-2 {className}">
		<span class="text-sm font-medium">Active filters:</span>

		{#if filters.search}
			<Badge variant="secondary" class="gap-1">
				<Search class="h-3 w-3" />
				{filters.search}
				<button onclick={filterActions.clearSearch} class="ml-1 hover:text-destructive">
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/if}

		{#each filters.genres as genre (genre)}
			<Badge variant="secondary" class="gap-1">
				<Hash class="h-3 w-3" />
				{genre}
				<button
					onclick={() => filterActions.toggleGenre(genre)}
					class="ml-1 hover:text-destructive"
				>
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/each}

		{#each filters.seasons as season (season)}
			<Badge variant="secondary" class="gap-1 capitalize">
				<Sparkles class="h-3 w-3" />
				{season}
				<button
					onclick={() => filterActions.toggleSeason(season)}
					class="ml-1 hover:text-destructive"
				>
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/each}

		{#each filters.years as year (year)}
			<Badge variant="secondary" class="gap-1">
				<Calendar class="h-3 w-3" />
				{year}
				<button onclick={() => filterActions.toggleYear(year)} class="ml-1 hover:text-destructive">
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/each}

		<Button variant="ghost" size="sm" onclick={filterActions.clearAll}>Clear all</Button>
	</div>
{/if}
