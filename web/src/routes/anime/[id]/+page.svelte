<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import { cn } from '$lib/utils';
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let img: HTMLImageElement | null = $state(null);

	const onScroll = () => {
		if (!img) return;
		const navbar = document.getElementById('navbar');
		if (!navbar) return;
		const navbarHeight = navbar.offsetHeight;
		const imageHeight = img.offsetHeight;
		const scrollY = window.scrollY;
		if (scrollY < imageHeight - navbarHeight) {
			img.style.position = 'relative';
			img.style.top = '0px';
		} else {
			img.style.position = 'sticky';
			img.style.top = `-${imageHeight - navbarHeight}px`;
		}
	};

	onMount(() => {
		onScroll();
	});

	const ratings: Record<string, string> = {
		g: 'G - All Ages',
		pg: 'PG - Children',
		pg_13: 'PG-13 - Teens 13 or older',
		r: 'R - 17+ (violence & profanity)',
		r_plus: 'R+ - Mild Nudity',
		rx: 'Rx - Hentai',
		unknown: 'Unknown Rating',
	};

	let infoPills = $derived.by(() => {
		return [
			data.anime?.data?.metadata?.mediaType.toUpperCase(),
			ratings[data.anime?.data?.metadata?.rating ?? 'unknown'],
			data.anime?.data?.metadata?.airingStatus
				.split('_')
				.map((text) => text.charAt(0).toUpperCase() + text.slice(1).toLowerCase())
				.join(' '),
		];
	});
</script>

<svelte:window onscroll={onScroll} />

<svelte:head>
	<title>{data.anime?.data?.jname} - Aniways</title>
	<meta
		name="description"
		content={data.anime?.data?.metadata?.description ??
			`Details and information about ${data.anime?.data?.jname}`}
	/>
</svelte:head>

<img
	bind:this={img}
	src={data.banner?.data?.url ?? data.anime?.data?.imageUrl}
	alt="Banner"
	class={cn('z-20 -mt-17 aspect-[19/4] w-full overflow-hidden object-center')}
/>

<div class="px-8 py-4">
	<div class="relative z-10 flex gap-4">
		<img
			src={data.anime?.data?.metadata?.mainPictureUrl ?? data.anime?.data?.imageUrl}
			alt="Anime Cover"
			class="sticky top-20 h-fit w-1/2 max-w-sm rounded-lg object-cover object-center shadow-lg"
		/>
		<div class="p-4">
			<h1 class="text-3xl font-bold">
				{data.anime?.data?.jname}
			</h1>
			<h2 class="text-xl text-muted-foreground">
				{data.anime?.data?.ename}
			</h2>
			<div class="mt-2 flex flex-wrap gap-2">
				{#each infoPills as info (info)}
					<span
						class="rounded-md border border-primary bg-primary/30 px-2 py-1 text-sm font-medium text-primary"
					>
						{info}
					</span>
				{/each}
			</div>
			{#if data.anime?.data?.genre?.split(', ').length}
				<div class="mt-2 flex flex-wrap gap-2">
					{#each data.anime?.data?.genre?.split(', ') as genre (genre)}
						<Button size="sm" variant="outline" href="/catalog?genres={genre}">{genre}</Button>
					{/each}
				</div>
			{/if}
			<p class="mt-2 text-muted-foreground">{data.anime?.data?.metadata?.description}</p>
		</div>
	</div>
</div>

<div class="h-17"></div>
