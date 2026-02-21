<script lang="ts">
	import { CircleCheck, CirclePlay, CircleX, Clock, Pause } from 'lucide-svelte';
	import { watch } from 'runed';
	import LibraryStatusTabs from '$lib/components/anime/controls/library-status-tabs.svelte';
	import AnimeGrid from '$lib/components/anime/display/anime-grid.svelte';
	import EmptyState from '$lib/components/anime/display/empty-state.svelte';
	import FilterSidebar from '$lib/components/anime/filters/filter-sidebar.svelte';
	import MobileFilters from '$lib/components/anime/filters/mobile-filters.svelte';
	import AnimePageHeader from '$lib/components/anime/layout/anime-page-header.svelte';
	import { Label } from '$lib/components/ui/label';
	import { getLayoutStateContext } from '$lib/context/layout.svelte';
	import { FilterManager } from '$lib/utils/filter-manager.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const layoutState = getLayoutStateContext();

	const filterManager = new FilterManager(data.initialFilters);

	const sortOptions = [
		{ value: 'library_updated_at', label: 'Library Updated' },
		{ value: 'ename', label: 'Name (A-Z)' },
		{ value: 'jname', label: 'Japanese Name' },
		{ value: 'anime_updated_at', label: 'Anime Updated' },
		{ value: 'relevance', label: 'Relevance' },
		{ value: 'season', label: 'Season' },
		{ value: 'year', label: 'Year' },
	];

	const statusTabs = [
		{ value: 'watching', label: 'Watching', icon: CirclePlay, iconColor: 'text-blue-500' },
		{ value: 'planning', label: 'Plan to Watch', icon: Clock, iconColor: 'text-orange-500' },
		{ value: 'completed', label: 'Completed', icon: CircleCheck, iconColor: 'text-emerald-500' },
		{ value: 'paused', label: 'On Hold', icon: Pause, iconColor: 'text-amber-500' },
		{ value: 'dropped', label: 'Dropped', icon: CircleX, iconColor: 'text-red-500' },
	] as const;

	let currentTab = $derived.by(() => {
		return statusTabs.find((tab) => tab.value === data.status) || statusTabs[0];
	});

	let description = $derived.by(() => {
		switch (data.status) {
			case 'watching':
				return 'Track the anime you are currently watching and stay updated with your progress.';
			case 'planning':
				return 'Keep a list of anime you plan to watch in the future and never miss out on new shows.';
			case 'completed':
				return 'View and manage the anime you have completed watching.';
			case 'paused':
				return 'Manage the anime you have put on hold and plan to resume later.';
			case 'dropped':
				return 'Keep track of the anime you have decided to stop watching.';
			default:
				return '';
		}
	});

	watch(
		() => filterManager.filters.page,
		() => {
			window.scrollTo({ top: 0, behavior: 'smooth' });
		},
	);
</script>

<svelte:head>
	<title>My List - AniStream</title>
	<meta name="description" content="Manage your anime watchlist and track your viewing progress" />
</svelte:head>

<div class="min-h-screen bg-background">
	<AnimePageHeader
		title="My Anime List"
		description="Track and manage your anime collection"
		{filterManager}
		{sortOptions}
	/>

	<MobileFilters genres={data.genres || []} {filterManager}>
		<div class="flex flex-col gap-2">
			<Label class="text-xs font-medium">Filter by Status</Label>
			<LibraryStatusTabs
				class="overflow-x-auto"
				currentStatus={data.status}
				onStatusChange={filterManager.setStatus}
			/>
		</div>
	</MobileFilters>

	<div class="container mx-auto px-4 pt-4 pb-8">
		<div class="flex gap-8">
			<aside class="hidden w-64 shrink-0 lg:block">
				<FilterSidebar genres={data.genres || []} {filterManager} />
			</aside>

			<AnimeGrid
				anime={data.listings?.items || []}
				{filterManager}
				totalPages={data.listings?.pageInfo.totalPages || 1}
			>
				<div class="sticky z-20 mb-4 hidden lg:block" style="top: {layoutState.totalHeight}px">
					<div
						class="border-b bg-background/95 px-4 py-3 backdrop-blur supports-[backdrop-filter]:bg-background/60"
					>
						<LibraryStatusTabs
							currentStatus={data.status}
							onStatusChange={filterManager.setStatus}
						/>
					</div>
				</div>

				<div
					class="mb-4 max-w-[calc(100vw-32px)] overflow-x-auto border-b bg-background py-3 lg:hidden"
				>
					<LibraryStatusTabs currentStatus={data.status} onStatusChange={filterManager.setStatus} />
				</div>

				{#snippet empty()}
					<EmptyState
						icon={currentTab.icon}
						title="No anime in {currentTab.label}"
						{description}
						action={{
							label: 'Browse Catalog',
							href: '/catalog',
							variant: 'default',
						}}
						secondaryAction={{
							label: 'Explore Genres',
							href: '/genres',
							variant: 'outline',
						}}
					/>
				{/snippet}
			</AnimeGrid>
		</div>
	</div>
</div>
