<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { LayoutGrid, LayoutList } from 'lucide-svelte';

	interface Props {
		viewMode: 'grid' | 'list';
		onViewModeChange: (mode: 'grid' | 'list') => void;
		size?: 'sm' | 'md';
		class?: string;
	}

	let {
		viewMode = $bindable(),
		onViewModeChange,
		size = 'md',
		class: className = '',
	}: Props = $props();

	const sizeClasses = {
		sm: 'h-7 w-7',
		md: 'h-8 w-8',
	};

	const iconSizes = {
		sm: 'h-3 w-3',
		md: 'h-4 w-4',
	};

	function setViewMode(mode: 'grid' | 'list') {
		viewMode = mode;
		onViewModeChange(mode);
	}
</script>

<div class="flex gap-1 rounded-md border p-1 {className}">
	<Button
		variant={viewMode === 'grid' ? 'default' : 'ghost'}
		size="icon"
		class={sizeClasses[size]}
		onclick={() => setViewMode('grid')}
		title="Grid view"
	>
		<LayoutGrid class={iconSizes[size]} />
	</Button>
	<Button
		variant={viewMode === 'list' ? 'default' : 'ghost'}
		size="icon"
		class={sizeClasses[size]}
		onclick={() => setViewMode('list')}
		title="List view"
	>
		<LayoutList class={iconSizes[size]} />
	</Button>
</div>
