<script lang="ts">
	import { cn } from '$lib/utils';

	type Props = {
		size?: 'sm' | 'md' | 'lg' | 'xl';
		variant?: 'default' | 'gradient' | 'anime';
	};

	let { size = 'md', variant = 'default' }: Props = $props();

	const sizeClasses = {
		sm: 'md:text-lg text-sm',
		md: 'md:text-2xl text-base',
		lg: 'md:text-4xl text-2xl',
		xl: 'md:text-6xl text-3xl',
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
