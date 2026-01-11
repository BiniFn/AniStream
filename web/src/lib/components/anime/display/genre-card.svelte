<script lang="ts">
	import { onMount } from 'svelte';
	import type { components } from '$lib/api/openapi';

	type Props = components['schemas']['models.GenrePreview'];

	let { name, previews = [] }: Props = $props();

	let cardElement = $state<HTMLDivElement | null>(null);
	let shouldLoadImages = $state(false);

	onMount(() => {
		if (!cardElement || typeof IntersectionObserver === 'undefined') {
			// Fallback: load immediately if IntersectionObserver not available
			shouldLoadImages = true;
			return;
		}

		const observer = new IntersectionObserver(
			(entries) => {
				entries.forEach((entry) => {
					if (entry.isIntersecting) {
						shouldLoadImages = true;
						observer.disconnect();
					}
				});
			},
			{
				// Start loading when card is 100px away from viewport
				rootMargin: '100px',
			},
		);

		observer.observe(cardElement);

		return () => observer.disconnect();
	});
</script>

{#if name != null && previews.length}
	<div
		bind:this={cardElement}
		class="group relative overflow-hidden rounded-xl border bg-card transition hover:scale-[1.05] hover:shadow-lg"
	>
		<a
			href={`/catalog?genres=${name}`}
			class="absolute inset-0 z-20"
			aria-label={`Open ${name} genre`}
		></a>

		{#if previews.length > 0}
			<div class="grid aspect-video grid-cols-3 grid-rows-2 gap-0.5">
				{#if shouldLoadImages}
					{#each previews as url (url)}
						<img
							src={url}
							alt={name}
							class="h-full w-full object-cover"
							loading="lazy"
							decoding="async"
						/>
					{/each}
				{:else}
					<!-- Placeholder while waiting to load -->
					{#each previews as _ (name + _)}
						<div class="h-full w-full bg-muted animate-pulse"></div>
					{/each}
				{/if}
			</div>
		{:else}
			<div class="flex aspect-video items-center justify-center bg-muted">
				<span class="text-sm text-muted-foreground">{name}</span>
			</div>
		{/if}

		<div
			class="pointer-events-none absolute inset-0 bg-gradient-to-t from-background/80 via-background/20 to-transparent opacity-90 transition-opacity group-hover:opacity-70"
		></div>
		<div class="pointer-events-none absolute right-0 bottom-0 left-0 z-10 p-3 sm:p-4">
			<div
				class="inline-flex rounded-md bg-background/80 px-2 py-1 text-xs font-medium text-foreground capitalize sm:text-sm"
			>
				{name}
			</div>
		</div>
	</div>
{/if}
