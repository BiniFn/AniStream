<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Sheet from '$lib/components/ui/sheet';
	import type { Snippet } from 'svelte';
	import type { FilterState, FilterActions } from '$lib/utils/filters';

	interface Props {
		open: boolean;
		genres: string[];
		filters: FilterState;
		filterActions: FilterActions;
		totalFilters: number;
		onOpenChange: (open: boolean) => void;
		children?: Snippet;
	}

	let {
		open = $bindable(),
		genres,
		filters,
		filterActions,
		totalFilters,
		onOpenChange,
		children,
	}: Props = $props();

	const seasons = ['winter', 'spring', 'summer', 'fall', 'unknown'] as const;
	const currentYear = new Date().getFullYear();
	const years = Array.from({ length: 50 }, (_, i) => currentYear - i);

	let yearMinTimeout: NodeJS.Timeout;
	let yearMaxTimeout: NodeJS.Timeout;

	function handleYearMinChange(e: Event) {
		const target = e.currentTarget as HTMLInputElement;
		const value = target.value;
		const newMin = value ? Number(value) : undefined;

		clearTimeout(yearMinTimeout);
		yearMinTimeout = setTimeout(() => {
			filterActions.setYearRange(newMin, filters.yearMax);
		}, 800);
	}

	function handleYearMaxChange(e: Event) {
		const target = e.currentTarget as HTMLInputElement;
		const value = target.value;
		const newMax = value ? Number(value) : undefined;

		clearTimeout(yearMaxTimeout);
		yearMaxTimeout = setTimeout(() => {
			filterActions.setYearRange(filters.yearMin, newMax);
		}, 800);
	}
</script>

<Sheet.Root bind:open {onOpenChange}>
	<Sheet.Content side="bottom" class="flex h-[90vh] max-h-[90vh] flex-col">
		<div class="mx-auto flex h-full w-full max-w-lg flex-col">
			<Sheet.Header class="pb-2">
				<div class="mx-auto mb-2 h-1 w-12 rounded-full bg-muted"></div>
				<Sheet.Title class="text-center text-sm">Filters & Sort</Sheet.Title>
			</Sheet.Header>

			<div class="flex-1 space-y-3 overflow-y-auto px-4">
				{#if genres && genres.length > 0}
					<div class="space-y-2">
						<Label class="text-xs font-medium">Genres</Label>
						<div class="flex max-h-32 flex-wrap gap-1.5 overflow-y-auto pr-2">
							{#each genres as genre (genre)}
								<Badge
									variant={filters.genres.map((g) => g.toLowerCase()).includes(genre.toLowerCase())
										? 'default'
										: 'outline'}
									class="h-6 cursor-pointer px-2 py-0.5 text-xs"
									onclick={() => filterActions.toggleGenre(genre)}
								>
									{genre}
								</Badge>
							{/each}
						</div>
						{#if filters.genres.length > 1}
							<div class="flex gap-2 pt-1">
								<Button
									variant={filters.genresMode === 'any' ? 'default' : 'outline'}
									size="sm"
									onclick={() => filterActions.setGenresMode('any')}
									class="h-7 flex-1 text-xs"
								>
									Any
								</Button>
								<Button
									variant={filters.genresMode === 'all' ? 'default' : 'outline'}
									size="sm"
									onclick={() => filterActions.setGenresMode('all')}
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
								variant={filters.seasons.includes(season) ? 'default' : 'outline'}
								size="sm"
								onclick={() => filterActions.toggleSeason(season)}
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
								variant={filters.years.includes(year) ? 'default' : 'outline'}
								size="sm"
								onclick={() => filterActions.toggleYear(year)}
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
							value={filters.yearMin || ''}
							class="h-8 text-xs"
							oninput={handleYearMinChange}
						/>
						<Input
							type="number"
							placeholder="To"
							min="1970"
							max={currentYear}
							value={filters.yearMax || ''}
							class="h-8 text-xs"
							oninput={handleYearMaxChange}
						/>
					</div>
				</div>

				{@render children?.()}
			</div>

			<div class="border-t px-4 py-3">
				<div class="flex gap-2">
					{#if totalFilters > 0}
						<Button variant="outline" size="sm" class="flex-1" onclick={filterActions.clearAll}>
							Clear All
						</Button>
					{/if}
					<Button
						size="sm"
						class={totalFilters > 0 ? 'flex-1' : 'w-full'}
						onclick={() => (open = false)}
					>
						Close
					</Button>
				</div>
			</div>
		</div>
	</Sheet.Content>
</Sheet.Root>
