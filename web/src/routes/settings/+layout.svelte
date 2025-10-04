<script lang="ts">
	import { page } from '$app/state';
	import { layoutState } from '$lib/context/layout.svelte';
	import { cn } from '$lib/utils';
	import { Settings, TriangleAlert, User, Users } from 'lucide-svelte';
	import type { LayoutProps } from './$types';
	import * as Tabs from '$lib/components/ui/tabs';
	import { goto } from '$app/navigation';

	let { children }: LayoutProps = $props();

	const tabs = [
		{ id: 'preferences', label: 'Preferences', icon: Settings, href: '/settings' },
		{ id: 'account', label: 'Account', icon: User, href: '/settings/account' },
		{ id: 'integrations', label: 'Integrations', icon: Users, href: '/settings/integrations' },
		{ id: 'danger', label: 'Danger Zone', icon: TriangleAlert, href: '/settings/danger' },
	];

	let currentTab = $derived(
		tabs.find((tab) => page.url.pathname === tab.href)?.id || 'preferences',
	);

	let tabsContainer: HTMLDivElement;

	function scrollToActiveTab() {
		if (tabsContainer) {
			const activeElement = tabsContainer.querySelector(`[data-state="active"]`);

			if (activeElement) {
				activeElement.scrollIntoView({
					behavior: 'smooth',
					block: 'nearest',
					inline: 'center',
				});
			}
		}
	}

	$effect(() => {
		if (currentTab) {
			setTimeout(scrollToActiveTab, 100);
		}
	});
</script>

<svelte:head>
	<title>Settings - Aniways</title>
</svelte:head>

<div class="min-h-screen bg-background">
	<div
		class="z-30 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 md:sticky"
		style="top: {layoutState.navbarHeight}px"
	>
		<div class="container mx-auto px-3 py-3 sm:px-4 sm:py-4">
			<h1 class="text-lg font-bold">Settings</h1>
			<p class="text-xs text-muted-foreground">Manage your account preferences and integrations</p>
		</div>
	</div>

	<div class="container mx-auto px-4 py-8">
		<div class="flex flex-col gap-6 lg:flex-row">
			<aside class="hidden w-64 flex-shrink-0 lg:block">
				<nav class="space-y-1">
					{#each tabs as tab (tab.id)}
						<a
							href={tab.href}
							class={cn(
								'flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left text-sm font-medium transition-colors',
								currentTab === tab.id
									? 'bg-primary text-primary-foreground'
									: 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
							)}
						>
							<tab.icon class="h-4 w-4" />
							{tab.label}
						</a>
					{/each}
				</nav>
			</aside>

			<div class="w-full overflow-x-auto lg:hidden" bind:this={tabsContainer}>
				<Tabs.Root value={currentTab}>
					<Tabs.List>
						{#each tabs as tab (tab.id)}
							<Tabs.Trigger value={tab.id} onclick={() => goto(tab.href)}>
								<tab.icon class="h-4 w-4" />
								{tab.label}
							</Tabs.Trigger>
						{/each}
					</Tabs.List>
				</Tabs.Root>
			</div>

			<main class="flex-1">
				{@render children()}
			</main>
		</div>
	</div>
</div>
