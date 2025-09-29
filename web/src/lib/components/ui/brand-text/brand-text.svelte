<script lang="ts">
	import { cn } from '$lib/utils';

	export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';
	export let variant: 'default' | 'gradient' | 'anime' = 'default';

	const sizeClasses = {
		sm: 'text-lg',
		md: 'text-2xl',
		lg: 'text-4xl',
		xl: 'text-6xl',
	};

	const variantClasses = {
		default: 'text-foreground',
		gradient: 'text-primary',
		anime: 'text-primary',
	};
</script>

<span
	class={cn(
		'font-grotesk font-bold tracking-tight select-none',
		sizeClasses[size],
		variantClasses[variant],
	)}
>
	ANIWAYS
</span>

