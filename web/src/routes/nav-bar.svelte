<script lang="ts">
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Command from '$lib/components/ui/command';
	import * as Sheet from '$lib/components/ui/sheet';
	import { cn } from '$lib/utils';
	import { Menu, Search } from 'lucide-svelte';

	const links = [
		{
			label: 'Home',
			link: '/'
		},
		{
			label: 'Catalog',
			link: '/catalog'
		},
		{
			label: 'Genres',
			link: '/genres'
		},
		{
			label: 'My List',
			link: '/my-list'
		}
	];

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
								page.url.pathname === link.link && 'text-foreground'
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

				<Button href="/login" variant="outline" class="hidden lg:inline-flex">Sign In</Button>
				<Button href="/register" class="hidden lg:inline-flex">Register</Button>
			</div>
		</div>
	</div>
</header>

<Sheet.Root bind:open={isSheetOpen}>
	<Sheet.Content side="right">
		<Sheet.Header>
			<Sheet.Title>Menu</Sheet.Title>
		</Sheet.Header>
		<div class="flex flex-col gap-2 px-4">
			{#each links as link (link.link)}
				<a
					href={link.link}
					class={cn(
						'font-medium text-muted-foreground transition-colors hover:text-primary',
						page.url.pathname === link.link && 'text-foreground'
					)}
				>
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
			<Button href="/login" variant="outline" class="w-full">Sign In</Button>
			<Button href="/register" class="w-full">Register</Button>
		</div>
	</Sheet.Content>
</Sheet.Root>

<Command.Dialog bind:open={isSearchOpen}>
	<Command.Input />
	<Command.List>
		<Command.Empty>No results found.</Command.Empty>
		<Command.Group heading="Suggestions">
			<Command.Item>Calendar</Command.Item>
			<Command.Item>Search Emoji</Command.Item>
			<Command.Item>Calculator</Command.Item>
		</Command.Group>
	</Command.List>
</Command.Dialog>
