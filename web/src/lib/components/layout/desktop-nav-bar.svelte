<script lang="ts">
	import { onNavigate } from '$app/navigation';
	import { page } from '$app/state';
	import { apiClient } from '$lib/api/client';
	import type { components } from '$lib/api/openapi';
	import UserProfileDropdown from '$lib/components/layout/user-profile-dropdown.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Command from '$lib/components/ui/command';
	import { cn } from '$lib/utils';
	import { ChevronLeft, ChevronRight, Dice6, Heart, Library, Search, Swords } from 'lucide-svelte';
	import { getAppStateContext } from '$lib/context/state.svelte';
	import { getFullscreen, isMacOS, onFullscreenChange } from '$lib/hooks/is-electron';

	type AnimeResponse = components['schemas']['models.AnimeResponse'];

	let isFullscreen = $state(false);
	const isMac = isMacOS();

	// Show traffic light spacer only on macOS when not in fullscreen
	let showTrafficLightSpacer = $derived(isMac && !isFullscreen);

	// Track navigation history state
	let canGoBack = $state(false);
	let canGoForward = $state(false);

	function updateHistoryState() {
		// navigation.canGoBack/canGoForward are available in modern browsers
		if ('navigation' in window) {
			const nav = window.navigation as { canGoBack: boolean; canGoForward: boolean };
			canGoBack = nav.canGoBack;
			canGoForward = nav.canGoForward;
		} else {
			// Fallback: assume we can go back if history length > 1
			canGoBack = history.length > 1;
			canGoForward = false; // Can't reliably detect this without Navigation API
		}
	}

	$effect(() => {
		getFullscreen().then((fs) => (isFullscreen = fs));
		onFullscreenChange((fs) => (isFullscreen = fs));
	});

	$effect(() => {
		updateHistoryState();
		window.addEventListener('popstate', updateHistoryState);
		return () => window.removeEventListener('popstate', updateHistoryState);
	});

	const appState = getAppStateContext();
	let links = $derived.by(() => {
		const base = [
			{ label: 'Catalog', link: '/catalog', Icon: Library },
			{ label: 'Genres', link: '/genres', Icon: Swords },
			{ label: 'Random', link: '/random', Icon: Dice6 },
		];

		if (appState.isLoggedIn) {
			base.push({ label: 'My List', link: '/my-list', Icon: Heart });
		}

		return base;
	});

	let isSearchOpen = $state(false);
	let searchQuery = $state('');
	let searchResults = $state<AnimeResponse[]>([]);
	let isSearching = $state(false);
	let searchTimeout: NodeJS.Timeout;

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'k' && (event.metaKey || event.ctrlKey)) {
			event.preventDefault();
			isSearchOpen = !isSearchOpen;
		}
	}

	async function performSearch(query: string) {
		if (!query || query.length < 3) {
			searchResults = [];
			return;
		}

		isSearching = true;
		try {
			const response = await apiClient.GET('/anime/listings/search', {
				params: {
					query: {
						q: query,
						page: 1,
						itemsPerPage: 3,
					},
				},
			});

			if (response.data?.items) {
				searchResults = response.data.items;
			}
		} catch (error) {
			console.error('Search failed:', error);
			searchResults = [];
		}
		isSearching = false;
	}

	function handleSearchInput(value: string) {
		searchQuery = value;
		clearTimeout(searchTimeout);

		if (value.length < 3) {
			searchResults = [];
			isSearching = false;
			return;
		}

		isSearching = true;
		searchTimeout = setTimeout(() => {
			performSearch(value);
		}, 500);
	}

	$effect(() => {
		if (!isSearchOpen) {
			searchQuery = '';
			searchResults = [];
		}
	});

	onNavigate(() => {
		isSearchOpen = false;
		setTimeout(updateHistoryState, 50);
	});
</script>

<svelte:window on:keydown={handleKeydown} />

<header
	id="navbar"
	class="electron-drag border-b border-border bg-background/95 backdrop-blur-md supports-[backdrop-filter]:bg-background/60"
>
	<div class="electron-drag flex items-center gap-6 px-6 py-4">
		<!-- Traffic light spacer - only on macOS, hidden in fullscreen -->
		{#if showTrafficLightSpacer}
			<div class="electron-drag w-[70px] shrink-0"></div>
		{/if}

		<!-- Back/Forward buttons -->
		<div class="electron-no-drag flex items-center">
			<button
				onclick={() => {
					history.back();
					setTimeout(updateHistoryState, 50);
				}}
				disabled={!canGoBack}
				class="cursor-pointer rounded-md p-1 text-foreground transition-colors hover:bg-muted disabled:pointer-events-none disabled:opacity-40"
				aria-label="Go back"
			>
				<ChevronLeft class="h-6 w-6" />
			</button>
			<button
				onclick={() => {
					history.forward();
					setTimeout(updateHistoryState, 50);
				}}
				disabled={!canGoForward}
				class="cursor-pointer rounded-md p-1 text-foreground transition-colors hover:bg-muted disabled:pointer-events-none disabled:opacity-40"
				aria-label="Go forward"
			>
				<ChevronRight class="h-6 w-6" />
			</button>
		</div>

		<!-- Nav links -->
		<nav class="electron-no-drag flex items-center gap-6">
			<a
				href="/"
				class={cn(
					'font-medium text-muted-foreground transition-colors hover:text-primary',
					page.url.pathname === '/' && 'text-foreground',
				)}
			>
				Home
			</a>
			{#each links as link (link.link)}
				<a
					href={link.link}
					class={cn(
						'font-medium text-muted-foreground transition-colors hover:text-primary',
						page.url.pathname === link.link && 'text-foreground',
					)}
				>
					{link.label}
				</a>
			{/each}
		</nav>

		<!-- Spacer to push right content - this is the main drag area -->
		<div class="electron-drag flex-1"></div>

		<!-- Right: Search + Auth -->
		<div class="electron-no-drag flex items-center gap-4">
			<Button
				variant="outline"
				class="w-64 items-center justify-start gap-2 border-muted-foreground/20 bg-transparent text-muted-foreground hover:border-primary/50"
				onclick={() => (isSearchOpen = true)}
			>
				<Search class="h-4 w-4" />
				<span>Search anime...</span>
				<kbd
					class="pointer-events-none ml-auto inline-flex h-5 items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100 select-none"
				>
					<span class="text-xs">⌘</span>K
				</kbd>
			</Button>

			{#if !appState.isLoggedIn}
				<Button href="/login" variant="outline">Sign In</Button>
				<Button href="/register">Register</Button>
			{:else}
				<UserProfileDropdown />
			{/if}
		</div>
	</div>
</header>

<Command.Dialog bind:open={isSearchOpen} shouldFilter={false}>
	<Command.Input
		placeholder="Search anime..."
		value={searchQuery}
		oninput={(e) => handleSearchInput(e.currentTarget.value)}
	/>
	<Command.List class="max-h-[400px]">
		{#if searchQuery.length < 3}
			<Command.Empty>Type at least 3 characters to search...</Command.Empty>
		{:else if isSearching}
			<Command.Empty>Searching...</Command.Empty>
		{:else}
			<Command.Group>
				{#if searchResults.length === 0}
					<Command.Empty>No results found.</Command.Empty>
				{:else}
					{#each searchResults as anime (anime.id)}
						<Command.LinkItem
							value={anime.ename || anime.jname}
							class="flex cursor-pointer items-center gap-3 p-3 hover:bg-accent"
							href="/anime/{anime.id}"
						>
							<div class="relative h-16 w-12 flex-shrink-0 overflow-hidden rounded-md">
								<img
									src={anime.imageUrl}
									alt={anime.jname || anime.ename}
									class="h-full w-full object-cover"
								/>
							</div>
							<div class="min-w-0 flex-1">
								<div class="line-clamp-1 font-medium">
									{anime.jname || anime.ename}
								</div>
								{#if anime.jname && anime.ename}
									<div class="line-clamp-1 text-sm text-muted-foreground">
										{anime.ename}
									</div>
								{/if}
								<div class="flex items-center gap-2 text-xs text-muted-foreground">
									<span class="capitalize">{anime.season} {anime.seasonYear}</span>
									{#if anime.genre}
										<span>•</span>
										<span class="line-clamp-1"
											>{anime.genre.split(', ').slice(0, 2).join(', ')}</span
										>
									{/if}
								</div>
							</div>
						</Command.LinkItem>
					{/each}
				{/if}

				{#if searchQuery.length >= 3}
					{#if searchResults.length > 0}
						<Command.Separator />
					{/if}
					<Command.LinkItem
						value="view-all-{searchQuery}"
						class="flex cursor-pointer items-center justify-center gap-2 p-3 font-medium text-primary hover:bg-accent"
						href="/catalog?search={searchQuery}"
					>
						<Search class="h-4 w-4" />
						{searchResults.length > 0
							? `View all results for "${searchQuery}"`
							: `Search for "${searchQuery}" in catalog`}
					</Command.LinkItem>
				{/if}
			</Command.Group>
		{/if}
	</Command.List>
</Command.Dialog>
