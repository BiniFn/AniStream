<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Search, X } from 'lucide-svelte';

	interface Props {
		value: string;
		placeholder?: string;
		onInput: (value: string) => void;
		class?: string;
		size?: 'sm' | 'md' | 'lg';
	}

	let {
		value = $bindable(),
		placeholder = 'Search anime...',
		onInput,
		class: className = '',
		size = 'md',
	}: Props = $props();

	const sizeClasses = {
		sm: 'h-8 text-sm',
		md: 'h-10',
		lg: 'h-12 text-lg',
	};

	const iconSizes = {
		sm: 'h-3 w-3',
		md: 'h-4 w-4',
		lg: 'h-5 w-5',
	};

	function handleInput(e: Event) {
		const target = e.currentTarget as HTMLInputElement;
		value = target.value;
		onInput(target.value);
	}

	function clearSearch() {
		value = '';
		onInput('');
	}
</script>

<div class="relative {className}">
	<Search
		class="absolute top-1/2 left-3 {iconSizes[size]} -translate-y-1/2 text-muted-foreground"
	/>
	<Input
		type="text"
		{placeholder}
		{value}
		oninput={handleInput}
		class="pr-10 pl-10 {sizeClasses[size]}"
	/>
	{#if value}
		<button
			onclick={clearSearch}
			class="absolute top-1/2 right-3 -translate-y-1/2 text-muted-foreground transition-colors hover:text-foreground"
		>
			<X class={iconSizes[size]} />
		</button>
	{/if}
</div>
