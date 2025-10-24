<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Select from '$lib/components/ui/select';
	import { cn } from '$lib/utils';
	import type { FilterManager } from '$lib/utils/filter-manager.svelte';
	import type { FilterState } from '$lib/utils/filters';

	interface SortOption {
		value: string;
		label: string;
	}

	interface Props {
		sortOptions: SortOption[];
		class?: string;
		selectClass?: string;
		filterManager: FilterManager;
	}

	let { filterManager, sortOptions, class: className = '', selectClass = '' }: Props = $props();

	function handleSortByChange(value: FilterState['sortBy'] | undefined) {
		if (value) {
			filterManager.setSort(value, filterManager.filters.sortOrder);
		}
	}

	function toggleSortOrder() {
		const newOrder = filterManager.filters.sortOrder === 'asc' ? 'desc' : 'asc';
		filterManager.setSort(filterManager.filters.sortBy, newOrder);
	}
</script>

<div class={cn('flex items-center gap-2', className)}>
	<Select.Root
		type="single"
		value={filterManager.filters.sortBy}
		onValueChange={(value) => handleSortByChange(value as FilterState['sortBy'])}
	>
		<Select.Trigger class={selectClass}>
			<span>
				{sortOptions.find((o) => o.value === filterManager.filters.sortBy)?.label || 'Sort by'}
			</span>
		</Select.Trigger>
		<Select.Content>
			{#each sortOptions as option (option.value)}
				<Select.Item value={option.value}>{option.label}</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>

	<Button
		variant="outline"
		size="icon"
		onclick={toggleSortOrder}
		title={filterManager.filters.sortOrder === 'asc' ? 'Sort ascending' : 'Sort descending'}
	>
		{filterManager.filters.sortOrder === 'asc' ? '↑' : '↓'}
		<span class="sr-only">
			{filterManager.filters.sortOrder === 'asc' ? 'Sort ascending' : 'Sort descending'}
		</span>
	</Button>
</div>
