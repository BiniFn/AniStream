<script lang="ts">
	import { invalidate } from '$app/navigation';
	import { PUBLIC_API_URL } from '$env/static/public';
	import { apiClient } from '$lib/api/client';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Import, TriangleAlert } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import type { PageProps } from './$types';
	import { getAppStateContext } from '$lib/context/state.svelte';

	let { data }: PageProps = $props();
	const appState = getAppStateContext();

	let oauthProviders = $derived(data.oauthProviders);
	let isDisconnecting = $state<string | null>(null);

	let showImportDialog = $state(false);
	let selectedProvider = $state<'myanimelist' | 'anilist' | null>(null);
	let isImporting = $state(false);

	const availableProviders = [
		{ name: 'myanimelist', displayName: 'MyAnimeList', icon: '/mal.svg' },
		{ name: 'anilist', displayName: 'AniList', icon: '/anilist.svg' },
	] as const;

	async function connectOAuth(provider: 'myanimelist' | 'anilist') {
		window.location.href = `${PUBLIC_API_URL}/auth/oauth/${provider}?redirect=${encodeURIComponent(window.location.href)}`;
	}

	async function disconnectOAuth(provider: string) {
		isDisconnecting = provider;
		try {
			await apiClient.DELETE('/auth/providers/{provider}', {
				params: {
					path: { provider },
				},
			});
			await invalidate(
				(url) => url.pathname.startsWith('/auth') || url.pathname.startsWith('/settings'),
			);
			toast.success('OAuth provider disconnected successfully');
		} catch (error) {
			console.error('Failed to disconnect OAuth provider:', error);
			toast.error('Failed to disconnect OAuth provider');
		}
		isDisconnecting = null;
	}

	async function importLibrary(provider: 'myanimelist' | 'anilist') {
		isImporting = true;
		try {
			const response = await apiClient.POST('/library/import', {
				params: {
					query: { provider },
				},
			});

			if (response.data?.id) {
				appState.setImportJobId(response.data.id);
				toast.success('Library import started. You will be notified when it completes.');
				showImportDialog = false;
			}
		} catch (error) {
			console.error('Failed to start library import:', error);
			toast.error('Failed to start library import. Please try again.');
		}
		isImporting = false;
	}
</script>

<div class="space-y-6">
	<Card.Root>
		<Card.Header>
			<Card.Title>External Connections</Card.Title>
			<Card.Description>Connect your accounts to sync data and import libraries</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			{#each availableProviders as provider (provider)}
				{@const isConnected = oauthProviders.includes(provider.name)}
				<div class="flex items-center justify-between rounded-lg border p-4">
					<div class="flex items-center gap-3">
						<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
							<img src={provider.icon} alt={provider.displayName} class="h-6 w-6" />
						</div>
						<div>
							<h3 class="font-medium">{provider.displayName}</h3>
							<p class="text-sm text-muted-foreground">
								{isConnected ? 'Connected' : 'Not connected'}
							</p>
						</div>
					</div>
					<div class="flex items-center gap-2">
						{#if isConnected}
							<Button
								variant="outline"
								size="sm"
								onclick={() => disconnectOAuth(provider.name)}
								disabled={isDisconnecting === provider.name}
							>
								{isDisconnecting === provider.name ? 'Disconnecting...' : 'Disconnect'}
							</Button>
						{:else}
							<Button size="sm" onclick={() => connectOAuth(provider.name)}>Connect</Button>
						{/if}
					</div>
				</div>
			{/each}
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title class="flex items-center gap-2">
				<Import class="h-5 w-5" />
				Library Import
			</Card.Title>
			<Card.Description>Import your anime library from external services</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="rounded-lg border border-amber-800 bg-amber-950 p-4">
				<div class="flex items-start gap-3">
					<TriangleAlert class="mt-0.5 h-5 w-5 text-amber-400" />
					<div class="space-y-2">
						<h4 class="font-medium text-amber-200">Important Notes</h4>
						<ul class="list-disc space-y-1 text-sm text-amber-300">
							<li>Importing will add anime to your library but won't remove existing entries</li>
							<li>AniList tokens expire yearly and require re-authentication</li>
							<li>Make sure you're connected to the provider before importing</li>
							<li>Large libraries may take several minutes to import</li>
						</ul>
					</div>
				</div>
			</div>

			<div class="grid gap-4 md:grid-cols-2">
				<div class="rounded-lg border p-4">
					<h3 class="mb-2 font-medium">MyAnimeList</h3>
					<p class="mb-4 text-sm text-muted-foreground">Import your anime list from MyAnimeList</p>
					<Button
						variant="outline"
						class="w-full"
						onclick={() => {
							selectedProvider = 'myanimelist';
							showImportDialog = true;
						}}
						disabled={!oauthProviders.includes('myanimelist')}
					>
						Import from myanimelist
					</Button>
				</div>

				<div class="rounded-lg border p-4">
					<h3 class="mb-2 font-medium">AniList</h3>
					<p class="mb-4 text-sm text-muted-foreground">Import your anime list from AniList</p>
					<Button
						variant="outline"
						class="w-full"
						onclick={() => {
							selectedProvider = 'anilist';
							showImportDialog = true;
						}}
						disabled={!oauthProviders.includes('anilist')}
					>
						Import from AniList
					</Button>
				</div>
			</div>
		</Card.Content>
	</Card.Root>
</div>

<Dialog.Root bind:open={showImportDialog}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title
				>Import Library from {selectedProvider === 'myanimelist'
					? 'MyAnimeList'
					: 'AniList'}</Dialog.Title
			>
			<Dialog.Description>
				This will import your anime library. The process may take several minutes for large
				libraries.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-4">
			<div class="rounded-lg border border-blue-800 bg-blue-950 p-4">
				<div class="flex items-start gap-3">
					<TriangleAlert class="mt-0.5 h-5 w-5 text-blue-400" />
					<div class="space-y-2">
						<h4 class="font-medium text-blue-200">Before Importing</h4>
						<ul class="space-y-1 text-sm text-blue-300">
							<li>
								• Make sure you're connected to {selectedProvider === 'myanimelist'
									? 'MyAnimeList'
									: 'AniList'}
							</li>
							<li>• Import will add anime to your library without removing existing entries</li>
							<li>• Large libraries may take 5-10 minutes to complete</li>
							{#if selectedProvider === 'anilist'}
								<li>• AniList tokens expire yearly - you may need to reconnect</li>
							{/if}
						</ul>
					</div>
				</div>
			</div>
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showImportDialog = false)}>Cancel</Button>
			<Button onclick={() => importLibrary(selectedProvider!)} disabled={isImporting}>
				{isImporting
					? 'Importing...'
					: `Import from ${selectedProvider === 'myanimelist' ? 'myanimelist' : 'AniList'}`}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
