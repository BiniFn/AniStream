<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import { cn } from '$lib/utils';
	import { X } from 'lucide-svelte';
	import { getLayoutStateContext } from '$lib/context/layout.svelte';
	import type { FilterManager } from '$lib/utils/filter-manager.svelte';

	interface Props {
		genres: string[];
		filterManager: FilterManager;
		showItemsPerPage?: boolean;
	}

	let { genres, filterManager, showItemsPerPage = true }: Props = $props();
	const layoutState = getLayoutStateContext();

	const sidebarTop = $derived(layoutState.totalHeight + 16);

	const seasons = ['winter', 'spring', 'summer', 'fall', 'unknown'] as const;
	const currentYear = new Date().getFullYear();
	const years = Array.from({ length: 50 }, (_, i) => currentYear - i);
	const itemsPerPageOptions = [12, 24, 36, 48];

	let yearMinTimeout: NodeJS.Timeout;
	let yearMaxTimeout: NodeJS.Timeout;

	function handleYearMinChange(e: Event) {
		const target = e.currentTarget as HTMLInputElement;
		const value = target.value;
		const newMin = value ? Number(value) : undefined;

		clearTimeout(yearMinTimeout);
		yearMinTimeout = setTimeout(() => {
			filterManager.setYearRange(newMin, filterManager.filters.yearMax);
		}, 800);
	}

	function handleYearMaxChange(e: Event) {
		const target = e.currentTarget as HTMLInputElement;
		const value = target.value;
		const newMax = value ? Number(value) : undefined;

		clearTimeout(yearMaxTimeout);
		yearMaxTimeout = setTimeout(() => {
			filterManager.setYearRange(filterManager.filters.yearMin, newMax);
		}, 800);
	}
</script>

<div
	class="filter-scroll sticky space-y-6 overflow-y-auto pr-2"
	style="top: {sidebarTop}px; max-height: calc(100vh - {sidebarTop}px - 2rem);"
>
	{#if genres && genres.length > 0}
		<div class="space-y-3">
			<div class="flex items-center justify-between">
				<Label class="text-base font-semibold">Genres</Label>
				{#if filterManager.filters.genres.length > 0}
					<Badge variant="secondary" class="text-xs">
						{filterManager.filters.genres.length}
					</Badge>
				{/if}
			</div>
			<div class="filter-scroll max-h-60 space-y-1 overflow-y-auto pr-2">
				{#each genres as genre (genre)}
					{@const isSelected = filterManager.filters.genres
						.map((g) => g.toLowerCase())
						.includes(genre.toLowerCase())}
					<button
						onclick={() => filterManager.toggleGenre(genre)}
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
			{#if filterManager.filters.genres.length > 1}
				<div class="flex gap-2 pt-2">
					<Button
						variant={filterManager.filters.genresMode === 'any' ? 'default' : 'outline'}
						size="sm"
						onclick={() => filterManager.setGenresMode('any')}
						class="flex-1"
					>
						Any
					</Button>
					<Button
						variant={filterManager.filters.genresMode === 'all' ? 'default' : 'outline'}
						size="sm"
						onclick={() => filterManager.setGenresMode('all')}
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
					variant={filterManager.filters.seasons.includes(season) ? 'default' : 'outline'}
					size="sm"
					onclick={() => filterManager.toggleSeason(season)}
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
					variant={filterManager.filters.years.includes(year) ? 'default' : 'outline'}
					size="sm"
					onclick={() => filterManager.toggleYear(year)}
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
				value={filterManager.filters.yearMin || ''}
				oninput={handleYearMinChange}
			/>
			<Input
				type="number"
				placeholder="To year"
				min="1970"
				max={currentYear}
				value={filterManager.filters.yearMax || ''}
				oninput={handleYearMaxChange}
			/>
		</div>
	</div>

	{#if showItemsPerPage}
		<div class="space-y-3">
			<Label class="text-base font-semibold">Items per page</Label>
			<Select.Root
				type="single"
				value={filterManager.filters.itemsPerPage.toString()}
				onValueChange={(value) => {
					if (value) {
						filterManager.setItemsPerPage(Number(value));
					}
				}}
			>
				<Select.Trigger class="w-full">
					<span>{filterManager.filters.itemsPerPage} items</span>
				</Select.Trigger>
				<Select.Content>
					{#each itemsPerPageOptions as option (option)}
						<Select.Item value={option.toString()}>{option} items</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>
	{/if}
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
</style>
