<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Select from '$lib/components/ui/select';

	interface SortOption {
		value: string;
		label: string;
	}

	interface Props {
		sortBy: string;
		sortOrder: 'asc' | 'desc';
		sortOptions: SortOption[];
		onSortChange: (sortBy: string, sortOrder: 'asc' | 'desc') => void;
		class?: string;
		selectClass?: string;
	}

	let {
		sortBy = $bindable(),
		sortOrder = $bindable(),
		sortOptions,
		onSortChange,
		class: className = '',
		selectClass = '',
	}: Props = $props();

	function handleSortByChange(value: string | undefined) {
		if (value) {
			sortBy = value;
			onSortChange(sortBy, sortOrder);
		}
	}

	function toggleSortOrder() {
		sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
		onSortChange(sortBy, sortOrder);
	}
</script>

<div class="flex items-center gap-2 {className}">
	<Select.Root type="single" value={sortBy} onValueChange={handleSortByChange}>
		<Select.Trigger class={selectClass}>
			<span>{sortOptions.find((o) => o.value === sortBy)?.label || 'Sort by'}</span>
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
		title={sortOrder === 'asc' ? 'Sort ascending' : 'Sort descending'}
	>
		{sortOrder === 'asc' ? '↑' : '↓'}
	</Button>
</div>
