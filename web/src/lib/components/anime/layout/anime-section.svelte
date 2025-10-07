<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import { ChevronRight } from 'lucide-svelte';
	import type { ComponentType, Snippet } from 'svelte';

	type Props = {
		icon: ComponentType;
		title: string;
		viewAllHref?: string;
		class?: string;
		visible?: boolean;
		children: Snippet;
	};

	let { icon, title, viewAllHref, class: className, visible, children }: Props = $props();

	let Icon = $derived(icon);
	let isVisible = $derived(visible ?? true);
</script>

{#if isVisible}
	<section>
		<div class="mb-8 flex items-center justify-between">
			<div class="flex items-center gap-3">
				<Icon class="h-6 w-6 text-primary" />
				<h2 class="text-2xl font-bold sm:text-3xl">{title}</h2>
			</div>
			{#if viewAllHref}
				<Button variant="ghost" class="gap-2" href={viewAllHref}>
					View All
					<ChevronRight class="h-4 w-4" />
				</Button>
			{/if}
		</div>
		<div
			class={cn(
				'flex gap-4 overflow-x-auto pb-4 md:grid md:grid-cols-4 md:gap-6 md:overflow-visible md:pb-0 lg:grid-cols-6',
				className,
			)}
		>
			{@render children()}
		</div>
	</section>
{/if}
