<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Sheet from '$lib/components/ui/sheet';
	import type { FilterManager } from '$lib/utils/filter-manager.svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		genres: string[];
		filterManager: FilterManager;
		children?: Snippet;
	}

	let { genres, filterManager, children }: Props = $props();

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

<Sheet.Root
	open={filterManager.showMobileFilters}
	onOpenChange={(open) => filterManager.setMobileFiltersVisibility(open)}
>
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
									variant={filterManager.filters.genres
										.map((g) => g.toLowerCase())
										.includes(genre.toLowerCase())
										? 'default'
										: 'outline'}
									class="h-6 cursor-pointer px-2 py-0.5 text-xs"
									onclick={() => filterManager.toggleGenre(genre)}
								>
									{genre}
								</Badge>
							{/each}
						</div>
						{#if filterManager.filters.genres.length > 1}
							<div class="flex gap-2 pt-1">
								<Button
									variant={filterManager.filters.genresMode === 'any' ? 'default' : 'outline'}
									size="sm"
									onclick={() => filterManager.setGenresMode('any')}
									class="h-7 flex-1 text-xs"
								>
									Any
								</Button>
								<Button
									variant={filterManager.filters.genresMode === 'all' ? 'default' : 'outline'}
									size="sm"
									onclick={() => filterManager.setGenresMode('all')}
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
								variant={filterManager.filters.seasons.includes(season) ? 'default' : 'outline'}
								size="sm"
								onclick={() => filterManager.toggleSeason(season)}
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
								variant={filterManager.filters.years.includes(year) ? 'default' : 'outline'}
								size="sm"
								onclick={() => filterManager.toggleYear(year)}
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
							value={filterManager.filters.yearMin || ''}
							class="h-8 text-xs"
							oninput={handleYearMinChange}
						/>
						<Input
							type="number"
							placeholder="To"
							min="1970"
							max={currentYear}
							value={filterManager.filters.yearMax || ''}
							class="h-8 text-xs"
							oninput={handleYearMaxChange}
						/>
					</div>
				</div>

				{@render children?.()}
			</div>

			<div class="border-t px-4 py-3">
				<div class="flex gap-2">
					{#if filterManager.totalFilters > 0}
						<Button variant="outline" size="sm" class="flex-1" onclick={filterManager.clearAll}>
							Clear All
						</Button>
					{/if}
					<Button
						size="sm"
						class={filterManager.totalFilters > 0 ? 'flex-1' : 'w-full'}
						onclick={() => filterManager.setMobileFiltersVisibility(false)}
					>
						Close
					</Button>
				</div>
			</div>
		</div>
	</Sheet.Content>
</Sheet.Root>
