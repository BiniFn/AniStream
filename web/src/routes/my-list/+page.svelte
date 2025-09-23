<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import AnimeCard from '$lib/components/anime/anime-card.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import { Pagination } from '$lib/components/ui/pagination';
	import Skeleton from '$lib/components/ui/skeleton/skeleton.svelte';
	import { cn } from '$lib/utils';
	import { ChevronRight, CircleCheck, CirclePlay, CircleX, Clock, Pause } from 'lucide-svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const statusTabs = [
		{ value: 'watching', label: 'Watching', icon: CirclePlay, iconColor: 'text-blue-500' },
		{ value: 'planning', label: 'Plan to Watch', icon: Clock, iconColor: 'text-orange-500' },
		{ value: 'completed', label: 'Completed', icon: CircleCheck, iconColor: 'text-emerald-500' },
		{ value: 'paused', label: 'On Hold', icon: Pause, iconColor: 'text-amber-500' },
		{ value: 'dropped', label: 'Dropped', icon: CircleX, iconColor: 'text-red-500' },
	] as const;

	const currentTab = $derived(statusTabs.find((tab) => tab.value === data.status) || statusTabs[0]);

	function changeStatus(newStatus: string) {
		const url = new URL(page.url);
		url.searchParams.set('status', newStatus);
		url.searchParams.set('page', '1');
		goto(url.toString());
	}

	const hasContent = $derived(data.library.items.length > 0);
</script>

<svelte:head>
	<title>My List - Aniways</title>
	<meta name="description" content="Manage your anime watchlist and track your viewing progress" />
</svelte:head>

<div class="min-h-screen bg-background">
	<div
		class="sticky top-17 z-30 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
	>
		<div class="container mx-auto px-4 py-4">
			<div class="mb-4">
				<h1 class="text-2xl font-bold tracking-tight">My Anime List</h1>
				<p class="text-sm text-muted-foreground">Track and manage your anime collection</p>
			</div>

			<div class="overflow-x-auto">
				<div class="flex gap-2 pb-2">
					{#each statusTabs as tab (tab.value)}
						<Button
							variant={data.status === tab.value ? 'default' : 'outline'}
							size="sm"
							onclick={() => changeStatus(tab.value)}
							class="gap-2 whitespace-nowrap"
						>
							{@const Icon = tab.icon}
							<Icon class={cn('h-4 w-4', tab.iconColor)} />
							{tab.label}
						</Button>
					{/each}
				</div>
			</div>
		</div>
	</div>

	<div class="container mx-auto px-4 py-8">
		<main>
			{#if hasContent}
				<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
					{#each data.library.items as item, index (item.id)}
						<AnimeCard anime={item.anime} {index} libraryEntry={item} showLibraryInfo={true} />
					{/each}
				</div>

				<Pagination totalPages={data.library.pageInfo.totalPages} />
			{:else}
				<Card.Root class="border-dashed">
					<Card.Content class="flex flex-col items-center justify-center py-20 text-center">
						{@const Icon = currentTab.icon}
						<Icon class={cn('mb-4 h-16 w-16', currentTab.iconColor)} />
						<h3 class="mb-2 text-xl font-semibold">No anime in {currentTab.label}</h3>
						<p class="mb-6 max-w-md text-muted-foreground">
							{#if data.status === 'watching'}
								Start watching anime to see them here. Your progress will be automatically tracked.
							{:else if data.status === 'planning'}
								Add anime you plan to watch later. Keep track of shows you're interested in.
							{:else if data.status === 'completed'}
								Anime you've finished watching will appear here. Complete a series to add it.
							{:else if data.status === 'paused'}
								Anime you've put on hold will appear here. Take a break and come back later.
							{:else if data.status === 'dropped'}
								Anime you've decided not to continue will appear here.
							{/if}
						</p>
						<div class="flex gap-3">
							<Button href="/catalog" variant="default" class="gap-2">
								Browse Catalog
								<ChevronRight class="h-4 w-4" />
							</Button>
							<Button href="/genres" variant="outline" class="gap-2">Explore Genres</Button>
						</div>
					</Card.Content>
				</Card.Root>
			{/if}
		</main>
	</div>
</div>

{#if !data}
	<div class="container mx-auto px-4 py-8">
		<Skeleton class="mb-8 h-12 w-64" />
		<div class="mb-8 flex gap-2">
			{#each Array(5) as _, i (i)}
				<Skeleton class="h-10 w-32" />
			{/each}
		</div>
		<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
			{#each Array(12) as _, i (i)}
				<div>
					<Skeleton class="aspect-[3/4] w-full" />
					<Skeleton class="mt-2 h-4 w-full" />
				</div>
			{/each}
		</div>
	</div>
{/if}
