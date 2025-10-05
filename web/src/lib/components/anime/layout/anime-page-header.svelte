<script lang="ts">
	import PageHeader from '$lib/components/layout/page-header.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import { Funnel } from 'lucide-svelte';
	import type { FilterState, FilterActions } from '$lib/utils/filters';
	import ActiveFilters from '$lib/components/anime/filters/active-filters.svelte';
	import SearchBar from '$lib/components/anime/controls/search-bar.svelte';
	import SortControls from '$lib/components/anime/controls/sort-controls.svelte';
	import ViewModeToggle from '$lib/components/anime/controls/view-mode-toggle.svelte';

	interface SortOption {
		value: string;
		label: string;
	}

	interface Props {
		title: string;
		description?: string;
		filters: FilterState;
		filterActions: FilterActions;
		sortOptions: SortOption[];
		viewMode: 'grid' | 'list';
		totalFilters: number;
		showMobileFilters?: boolean;
		pageInfo?: {
			currentPage: number;
			totalPages: number;
		};
		onViewModeChange: (mode: 'grid' | 'list') => void;
		onMobileFiltersToggle?: () => void;
	}

	let {
		title,
		description,
		filters,
		filterActions,
		sortOptions,
		viewMode = $bindable(),
		totalFilters,
		showMobileFilters = true,
		pageInfo,
		onViewModeChange,
		onMobileFiltersToggle,
	}: Props = $props();

	const displayDescription = $derived(
		pageInfo ? `Page ${pageInfo.currentPage} of ${pageInfo.totalPages}` : description,
	);

	let searchTimeout: NodeJS.Timeout;

	function handleSearch(value: string) {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			filterActions.updateSearch(value);
		}, 500);
	}

	function handleSortChange(newSortBy: FilterState['sortBy'], newSortOrder: 'asc' | 'desc') {
		filterActions.setSort(newSortBy, newSortOrder);
	}
</script>

<PageHeader {title} description={displayDescription}>
	{#snippet actions()}
		<div class="flex flex-col gap-2 lg:flex-row lg:items-center lg:gap-3">
			<div class="flex flex-col gap-2 lg:hidden">
				<SearchBar value={filters.search} onInput={handleSearch} size="sm" />
				<div class="flex items-center gap-1">
					<SortControls
						sortBy={filters.sortBy}
						sortOrder={filters.sortOrder}
						{sortOptions}
						onSortChange={handleSortChange}
						selectClass="h-9 flex-1 text-sm"
					/>
				</div>
			</div>

			<div class="hidden lg:flex lg:items-center lg:gap-3">
				<SearchBar value={filters.search} onInput={handleSearch} class="w-80" />
				<SortControls
					sortBy={filters.sortBy}
					sortOrder={filters.sortOrder}
					{sortOptions}
					onSortChange={handleSortChange}
					selectClass="w-[172px]"
				/>
				<ViewModeToggle {viewMode} {onViewModeChange} />
			</div>

			{#if showMobileFilters && onMobileFiltersToggle}
				<button
					onclick={onMobileFiltersToggle}
					class={cn(
						buttonVariants({ variant: 'outline', size: 'icon' }),
						'relative h-9 w-9 lg:hidden',
					)}
				>
					<Funnel class="h-4 w-4" />
					{#if totalFilters > 0}
						<span
							class="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[10px] text-primary-foreground"
						>
							{totalFilters}
						</span>
					{/if}
				</button>
			{/if}
		</div>
	{/snippet}
</PageHeader>

<div class="container mx-auto px-4">
	<ActiveFilters {filters} {filterActions} {totalFilters} class="mt-4" />
</div>
