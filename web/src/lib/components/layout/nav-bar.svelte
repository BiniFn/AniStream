<script lang="ts">
	import { page } from '$app/state';
	import UserProfileDropdown from '$lib/components/layout/user-profile-dropdown.svelte';
	import * as Avatar from '$lib/components/ui/avatar';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Command from '$lib/components/ui/command';
	import * as Sheet from '$lib/components/ui/sheet';
	import { appState } from '$lib/context/state.svelte';
	import { cn } from '$lib/utils';
	import {
		Clock,
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

	let user = $derived(appState.user);
	let isLoggedIn = $derived(user != null);
	let links = $derived.by(() => {
		const base = [
			{ label: 'Home', link: '/', Icon: House },
			{ label: 'Catalog', link: '/catalog', Icon: Library },
			{ label: 'Genres', link: '/genres', Icon: Swords },
			{ label: 'Recent', link: '/recent', Icon: Clock },
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

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'k' && (event.metaKey || event.ctrlKey)) {
			event.preventDefault();
			isSearchOpen = !isSearchOpen;
			isSheetOpen = false;
		}
	}
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

<Command.Dialog bind:open={isSearchOpen}>
	<Command.Input />
	<Command.List>
		<Command.Empty>No results found.</Command.Empty>
	</Command.List>
</Command.Dialog>
