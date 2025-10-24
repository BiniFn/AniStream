<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils';
	import { Search, X } from 'lucide-svelte';

	interface Props {
		value: string;
		placeholder?: string;
		onInput: (value: string) => void;
		class?: string;
	}

	let { value, placeholder = 'Search anime...', onInput, class: className = '' }: Props = $props();

	function handleInput(e: Event) {
		const target = e.currentTarget as HTMLInputElement;
		onInput(target.value);
	}

	function clearSearch() {
		onInput('');
	}
</script>

<div class={cn('relative', className)}>
	<Search class={cn('absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground')} />
	<Input
		type="text"
		{placeholder}
		{value}
		oninput={handleInput}
		class={cn('h-9 pr-10 pl-10 text-sm')}
	/>
	{#if value}
		<button
			onclick={clearSearch}
			class="absolute top-1/2 right-3 -translate-y-1/2 cursor-pointer text-muted-foreground transition-colors hover:text-foreground"
		>
			<X class="size-4" />
		</button>
	{/if}
</div>
