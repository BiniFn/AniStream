<script lang="ts">
	import { layoutState } from '$lib/context/layout.svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		description?: string;
		actions?: Snippet;
	}

	let { title, description, actions }: Props = $props();

	let headerElement: HTMLElement;

	const updateHeaderHeight = () => {
		if (headerElement) {
			const height = headerElement.getBoundingClientRect().height;
			layoutState.headerHeight = height;
		}
	};

	$effect(() => {
		let resizeObserver: ResizeObserver;
		if (headerElement) {
			resizeObserver = new ResizeObserver(() => {
				updateHeaderHeight();
			});
			resizeObserver.observe(headerElement);
		}

		updateHeaderHeight();

		window.addEventListener('resize', updateHeaderHeight);
		return () => {
			window.removeEventListener('resize', updateHeaderHeight);
			if (resizeObserver && headerElement) {
				resizeObserver.unobserve(headerElement);
			}
		};
	});
</script>

<div
	bind:this={headerElement}
	class="z-30 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 md:sticky"
	style="top: {layoutState.navbarHeight}px"
>
	<div class="container mx-auto px-4 py-4">
		<div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
			<div>
				<h1 class="text-2xl font-bold tracking-tight">{title}</h1>
				{#if description}
					<p class="text-sm text-muted-foreground">{description}</p>
				{/if}
			</div>

			{#if actions}
				<div class="flex items-center gap-3">
					{@render actions()}
				</div>
			{/if}
		</div>
	</div>
</div>
