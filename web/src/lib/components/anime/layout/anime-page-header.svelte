<script lang="ts">
	import PageHeader from '$lib/components/layout/page-header.svelte';
	import { cn } from '$lib/utils';
	import { Funnel } from 'lucide-svelte';
	import ActiveFilters from '$lib/components/anime/filters/active-filters.svelte';
	import SearchBar from '$lib/components/anime/controls/search-bar.svelte';
	import SortControls from '$lib/components/anime/controls/sort-controls.svelte';
	import ViewModeToggle from '$lib/components/anime/controls/view-mode-toggle.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { FilterManager } from '$lib/utils/filter-manager.svelte';
	import { useDebounce } from 'runed';

	interface SortOption {
		value: string;
		label: string;
	}

	interface Props {
		title: string;
		description?: string;
		filterManager: FilterManager;
		sortOptions: SortOption[];
		pageInfo?: {
			currentPage: number;
			totalPages: number;
		};
	}

	let { title, description, filterManager, sortOptions, pageInfo }: Props = $props();

	const displayDescription = $derived(
		pageInfo ? `Page ${pageInfo.currentPage} of ${pageInfo.totalPages}` : description,
	);

	const handleSearch = useDebounce((value: string) => {
		filterManager.updateSearch(value);
	}, 500);
</script>

<PageHeader {title} description={displayDescription}>
	{#snippet actions()}
		<div class="flex w-full flex-col gap-2 lg:flex-row lg:items-center lg:gap-3">
			<div class="flex flex-col gap-2 lg:hidden">
				<SearchBar value={filterManager.filters.search} onInput={handleSearch} />

				<Button
					onclick={() => filterManager.setMobileFiltersVisibility(true)}
					variant="outline"
					class={cn('relative lg:hidden')}
				>
					<Funnel class="h-4 w-4" />
					Open Filters
					{#if filterManager.totalFilters > 0}
						<span
							class="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[10px] text-primary-foreground"
						>
							{filterManager.totalFilters}
						</span>
					{/if}
				</Button>

				<SortControls {sortOptions} selectClass="h-9 w-full text-sm" {filterManager} />
			</div>

			<div class="hidden lg:flex lg:items-center lg:gap-3">
				<SearchBar value={filterManager.filters.search} onInput={handleSearch} class="w-80" />
				<SortControls {sortOptions} selectClass="w-[172px]" {filterManager} />
				<ViewModeToggle
					viewMode={filterManager.viewMode}
					onViewModeChange={filterManager.handleViewModeChange}
				/>
			</div>
		</div>
	{/snippet}
</PageHeader>

<div class="container mx-auto px-4">
	<ActiveFilters {filterManager} class="mt-4" />
</div>
