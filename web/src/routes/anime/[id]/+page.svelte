<script lang="ts">
	import { page } from '$app/state';
	import { cn } from '$lib/utils';
	import type { PageProps } from './$types';
	import HeroSection from './hero-section.svelte';
	import InfoSidebar from './info-sidebar.svelte';
	import TabContent from './tab-content.svelte';

	let { data }: PageProps = $props();

	const tabs = ['overview', 'episodes', 'relations'] as const;
	let selectedTab: (typeof tabs)[number] = $derived.by(() => {
		const urlTab = page.url.searchParams.get('tab') || 'overview';
		return urlTab as (typeof tabs)[number];
	});
</script>

<svelte:head>
	<title>{data.anime.jname || data.anime.ename} - Aniways</title>
	<meta
		name="description"
		content={data.anime.metadata?.description ||
			`Watch ${data.anime.jname || data.anime.ename} on Aniways`}
	/>
</svelte:head>

<HeroSection
	anime={data.anime}
	banner={data.banner?.data?.url ?? null}
	trailer={data.trailer?.data?.trailer ?? null}
	ratingLabel={data.ratingLabel}
	totalEpisodes={data.episodes?.data?.length}
/>

<div
	class="sticky top-17 z-20 container mx-auto w-full border-b bg-background/95 px-4 backdrop-blur-md supports-[backdrop-filter]:bg-background/60"
>
	<nav class="mx-auto flex justify-center gap-8 md:justify-start">
		{#each tabs as tab (tab)}
			<a
				href={`?tab=${tab}`}
				class={cn(
					'relative py-4 text-sm font-medium transition-colors hover:text-primary',
					selectedTab === tab ? 'text-primary' : 'text-muted-foreground',
				)}
				data-sveltekit-noscroll
				data-sveltekit-replacestate
			>
				{tab.charAt(0).toUpperCase() + tab.slice(1)}
				{#if selectedTab === tab}
					<div class="absolute right-0 bottom-0 left-0 h-0.5 bg-primary"></div>
				{/if}
			</a>
		{/each}
	</nav>
</div>
<div class="container mx-auto px-4 py-8">
	<div class={cn('grid gap-8 lg:grid-cols-3')}>
		<TabContent
			{selectedTab}
			anime={data.anime}
			episodes={data.episodes?.data ?? []}
			franchise={data.franchise?.data ?? null}
		/>

		<InfoSidebar anime={data.anime} ratingLabel={data.ratingLabel} {selectedTab} />
	</div>
</div>
