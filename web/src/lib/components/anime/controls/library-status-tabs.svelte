<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import type { FilterState } from '$lib/utils/filters';
	import { CircleCheck, CirclePlay, CircleX, Clock, Pause } from 'lucide-svelte';

	interface Props {
		currentStatus: string;
		onStatusChange: (status: FilterState['status']) => void;
		class?: string;
	}

	let { currentStatus, onStatusChange, class: className = '' }: Props = $props();

	const statusTabs = [
		{ value: 'watching', label: 'Watching', icon: CirclePlay, iconColor: 'text-blue-500' },
		{ value: 'planning', label: 'Plan to Watch', icon: Clock, iconColor: 'text-orange-500' },
		{ value: 'completed', label: 'Completed', icon: CircleCheck, iconColor: 'text-emerald-500' },
		{ value: 'paused', label: 'On Hold', icon: Pause, iconColor: 'text-amber-500' },
		{ value: 'dropped', label: 'Dropped', icon: CircleX, iconColor: 'text-red-500' },
	] as const;
</script>

<div class={cn('flex gap-2 pb-2', className)}>
	{#each statusTabs as tab (tab.value)}
		<Button
			variant={currentStatus === tab.value ? 'default' : 'outline'}
			size="sm"
			onclick={() => onStatusChange(tab.value)}
			class="gap-2 whitespace-nowrap"
		>
			{@const Icon = tab.icon}
			<Icon class={cn('h-4 w-4', tab.iconColor)} />
			{tab.label}
		</Button>
	{/each}
</div>
