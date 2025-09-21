<script lang="ts">
	import { onNavigate } from '$app/navigation';
	import { page } from '$app/state';
	import { apiClient } from '$lib/api/client';
	import type { components } from '$lib/api/openapi';
	import UserProfileDropdown from '$lib/components/layout/user-profile-dropdown.svelte';
	import * as Avatar from '$lib/components/ui/avatar';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Command from '$lib/components/ui/command';
	import * as Sheet from '$lib/components/ui/sheet';
	import { appState } from '$lib/context/state.svelte';
	import { cn } from '$lib/utils';
	import {
		Heart,
		House,
		Library,
		LogOut,
		Menu,
		Search,
		Settings,
		Swords,
		User,
	} from 'lucide-svelte';

	type AnimeResponse = components['schemas']['models.AnimeResponse'];

	let user = $derived(appState.user);
	let isLoggedIn = $derived(user != null);
	let links = $derived.by(() => {
		const base = [
			{ label: 'Home', link: '/', Icon: House },
			{ label: 'Catalog', link: '/catalog', Icon: Library },
			{ label: 'Genres', link: '/genres', Icon: Swords },
		];

		if (isLoggedIn) {
			base.push({ label: 'My List', link: '/my-list', Icon: Heart });
		}

		return base;
	});

	let sheetLinks = $derived.by(() => {
		if (!isLoggedIn) {
			return links;
		}

		return [
			...links,
			{ label: 'Profile', link: '/profile', Icon: User },
			{ label: 'Settings', link: '/settings', Icon: Settings },
		];
	});

	let isSheetOpen = $state(false);
	let isSearchOpen = $state(false);
	let searchQuery = $state('');
	let searchResults = $state<AnimeResponse[]>([]);
	let isSearching = $state(false);
	let searchTimeout: NodeJS.Timeout;

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'k' && (event.metaKey || event.ctrlKey)) {
			event.preventDefault();
			isSearchOpen = !isSearchOpen;
			isSheetOpen = false;
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
						itemsPerPage: 5,
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
		isSheetOpen = false;
		isSearchOpen = false;
	});
</script>

<svelte:window on:keydown={handleKeydown} />

<header
	id="navbar"
	class="sticky top-0 z-50 border-b border-border bg-background/95 backdrop-blur-md supports-[backdrop-filter]:bg-background/60"
>
	<div class="container mx-auto p-4">
		<div class="flex items-center justify-between">
			<div class="flex items-center space-x-8">
				<h1 class="tracking-light font-serif text-3xl font-extrabold text-primary uppercase">
					<a href="/">Aniways</a>
				</h1>

				<nav class="hidden space-x-6 lg:flex">
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
			</div>

			<div class="flex items-center gap-4">
				<Button variant="outline" class="lg:hidden" onclick={() => (isSheetOpen = true)}>
					<Menu class="h-5 w-5" />
					<span class="sr-only">Menu</span>
				</Button>
				<Button
					variant="outline"
					class="hidden w-64 items-center justify-start space-x-2 border-muted-foreground/20 bg-transparent text-muted-foreground hover:border-primary/50 lg:flex"
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

				{#if !isLoggedIn}
					<Button href="/login" variant="outline" class="hidden lg:inline-flex">Sign In</Button>
					<Button href="/register" class="hidden lg:inline-flex">Register</Button>
				{:else}
					<UserProfileDropdown user={user!} class="hidden lg:flex" />
				{/if}
			</div>
		</div>
	</div>
</header>

<Sheet.Root bind:open={isSheetOpen}>
	<Sheet.Content side="right">
		<Sheet.Header>
			<Sheet.Title>Menu</Sheet.Title>
		</Sheet.Header>

		{#if isLoggedIn}
			<div class="px-4 pb-4">
				<div class="flex items-center gap-3 rounded-lg bg-muted/50 p-3">
					<UserProfileDropdown user={user!} class="hidden" />
					<Avatar.Root class="size-10">
						<Avatar.Image src={user?.profilePicture} alt={user?.username} />
						<Avatar.Fallback class="bg-primary/50 text-sm font-medium">
							{user?.username
								?.split(' ')
								.map((n) => n.charAt(0))
								.join('')
								.toUpperCase()
								.slice(0, 2)}
						</Avatar.Fallback>
					</Avatar.Root>
					<div class="flex flex-col">
						<p class="text-sm font-medium">{user?.username}</p>
						<p class="text-xs text-muted-foreground">{user?.email}</p>
					</div>
				</div>
			</div>
		{/if}

		<div class="flex flex-col gap-2 px-4">
			{#each sheetLinks as link (link.link)}
				<a
					href={link.link}
					class={cn(
						'flex items-center font-medium text-muted-foreground transition-colors hover:text-primary',
						page.url.pathname === link.link && 'text-foreground',
					)}
				>
					<link.Icon class="mr-2 h-4 w-4" />
					{link.label}
				</a>
			{/each}
		</div>

		<div class="flex flex-col gap-2 px-4">
			<Button
				variant="outline"
				class="w-full"
				onclick={() => {
					isSearchOpen = true;
					isSheetOpen = false;
				}}
			>
				<Search class="h-4 w-4" />
				<span>Search anime...</span>
			</Button>
			{#if !isLoggedIn}
				<Button href="/login" variant="outline" class="w-full">Sign In</Button>
				<Button href="/register" class="w-full">Register</Button>
			{:else}
				<Button href="/logout" variant="destructive" class="w-full">
					<LogOut class="h-4 w-4" />
					Log out
				</Button>
			{/if}
		</div>
	</Sheet.Content>
</Sheet.Root>

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
