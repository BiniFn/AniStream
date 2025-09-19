<script lang="ts">
	import { goto, preloadData } from '$app/navigation';
	import { apiClient } from '$lib/api/client';
	import type { components } from '$lib/api/openapi';
	import LibraryBtn from '$lib/components/anime/library-btn.svelte';
	import Player from '$lib/components/anime/player/index.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import { appState } from '$lib/context/state.svelte';
	import { cn } from '$lib/utils';
	import {
		ArrowLeft,
		ChevronLeft,
		ChevronRight,
		CircleAlert,
		Info,
		LoaderCircle,
		Play,
		RefreshCcw,
		Server,
		Star,
	} from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { PageData } from './$types';

	type StreamingData = components['schemas']['models.StreamingDataResponse'];
	type EpisodeServer = components['schemas']['models.EpisodeServerResponse'];

	let { data }: { data: PageData } = $props();

	let selectedServer: EpisodeServer | null = $state(data.servers[0] || null);
	let streamInfo: StreamingData | null = $state(null);
	let isLoading = $state(true);
	let videoError = $state(false);
	let errorMessage = $state('');
	let isFullscreen = $state(false);

	let episodesSearch = $state('');
	let showEpisodesSidebar = $state(true);
	let isMobile = $state(false);

	let groupedServers = $derived.by(() => {
		const groups: Record<string, EpisodeServer[]> = {};
		for (const server of data.servers) {
			const t = server.type.toLowerCase();
			(groups[t] ??= []).push(server);
		}
		return groups;
	});

	const hasNextEpisode = $derived(data.episodeNumber < data.episodes.length);
	const hasPrevEpisode = $derived(data.episodeNumber > 1);
	const nextEpisodeUrl = $derived.by(() =>
		hasNextEpisode ? `/anime/${data.anime.id}/watch?ep=${data.episodeNumber + 1}` : undefined,
	);

	$effect(() => {
		if (nextEpisodeUrl) preloadData(nextEpisodeUrl);
	});

	$effect(() => {
		const check = () => {
			isMobile = window.innerWidth < 768;
			if (isMobile) showEpisodesSidebar = false;
		};
		check();
		window.addEventListener('resize', check);
		return () => window.removeEventListener('resize', check);
	});

	let filteredEpisodes = $derived.by(() => {
		if (!episodesSearch) return data.episodes;
		return data.episodes.filter((ep) => {
			const s = (ep.number.toString() + (ep.title || `Episode ${ep.number}`)).toLowerCase();
			return s.includes(episodesSearch.toLowerCase());
		});
	});

	let loadToken = 0;
	$effect(() => {
		if (!selectedServer) return;

		const serverNow =
			data.servers.find((s) => s.serverId === selectedServer?.serverId) ?? data.servers[0];
		if (!serverNow) return;

		const token = ++loadToken;
		(async () => {
			isLoading = true;
			videoError = false;
			errorMessage = '';

			try {
				const res = await apiClient.GET('/anime/{id}/episodes/servers/{serverID}', {
					params: {
						path: { id: data.anime.id, serverID: selectedServer.serverId },
						query: {
							server: selectedServer.serverName.toLowerCase().replace(/\s+/g, '-'),
							type: selectedServer.type.toLowerCase(),
						},
					},
				});
				if (token !== loadToken) return;
				if (!res.data) throw new Error('No stream data received from server');
				streamInfo = res.data;
			} catch (e) {
				videoError = true;
				errorMessage = e instanceof Error ? e.message : 'Failed to load video stream';
				streamInfo = null;
			} finally {
				if (token === loadToken) isLoading = false;
			}
		})();
	});

	$effect(() => {
		const _deps = `${data.anime.id}:${data.episodeNumber}:${data.servers.length}`;

		const list = data.servers ?? [];
		if (!selectedServer || !list.some((s) => s.serverId === selectedServer?.serverId)) {
			selectedServer = list[0] || null;
		}
		streamInfo = null;
		videoError = false;
		errorMessage = '';
	});

	function changeEpisode(n: number) {
		goto(`/anime/${data.anime.id}/watch?ep=${n}`, { replaceState: true });
	}

	function selectServer(server: EpisodeServer) {
		selectedServer = server;
	}

	onMount(() => {
		const handleFS = () => (isFullscreen = !!document.fullscreenElement);
		document.addEventListener('fullscreenchange', handleFS);
		document.addEventListener('webkitfullscreenchange', handleFS);
		return () => {
			document.removeEventListener('fullscreenchange', handleFS);
			document.removeEventListener('webkitfullscreenchange', handleFS);
		};
	});
</script>

<svelte:head>
	<title>
		{data.anime.jname || data.anime.ename} - Episode {data.episodeNumber} - Aniways
	</title>
	<meta
		name="description"
		content="Watch {data.anime.jname || data.anime.ename} Episode {data.episodeNumber} on Aniways"
	/>
</svelte:head>

<div class="min-h-screen bg-background">
	<header
		class="sticky top-0 z-40 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
	>
		<div class="container mx-auto flex h-16 items-center gap-4 px-4">
			<Button href="/anime/{data.anime.id}" variant="ghost" size="icon">
				<ArrowLeft class="h-5 w-5" />
			</Button>
			<div class="min-w-0 flex-1">
				<h1 class="line-clamp-1 text-lg font-semibold">
					{data.anime.jname || data.anime.ename}
				</h1>
				<p class="text-sm text-muted-foreground">
					Episode {data.episodeNumber}: {data.currentEpisode.title ||
						`Episode ${data.episodeNumber}`}
				</p>
			</div>
			<Button
				href="/anime/{data.anime.id}"
				variant="outline"
				size="sm"
				class="hidden md:inline-flex"
			>
				<Info class="mr-2 h-4 w-4" />
				Anime Details
			</Button>
		</div>
	</header>

	<div
		class={cn(
			'flex flex-col p-4 lg:grid lg:gap-6 lg:p-6',
			showEpisodesSidebar ? 'lg:grid-cols-[1fr,380px]' : 'lg:grid-cols-1',
		)}
	>
		<div class="space-y-4">
			<div
				class={cn(
					'relative w-full overflow-hidden rounded-lg bg-black',
					isFullscreen ? 'fixed inset-0 z-50' : 'aspect-video',
				)}
				role="presentation"
			>
				{#if videoError}
					<div
						class="absolute inset-0 flex flex-col items-center justify-center gap-4 p-4 text-white"
					>
						<CircleAlert class="h-12 w-12 text-destructive" />
						<p class="text-center text-lg font-medium">{errorMessage || 'Failed to load video'}</p>
						<div class="flex gap-2">
							<Button
								onclick={() => {
									videoError = false;
									selectedServer = selectedServer && { ...selectedServer };
								}}
								variant="secondary"
								size="sm"
							>
								<RefreshCcw class="mr-2 h-4 w-4" />
								Try Again
							</Button>
							<Button
								onclick={() => {
									const next = data.servers.find((s) => s.serverId !== selectedServer?.serverId);
									if (next) {
										videoError = false;
										selectServer(next);
									}
								}}
								variant="secondary"
								size="sm"
								disabled={data.servers.length <= 1}
							>
								Try Different Server
							</Button>
						</div>
					</div>
				{:else if isLoading || !streamInfo}
					<div class="absolute inset-0 grid place-items-center">
						<LoaderCircle class="h-12 w-12 animate-spin text-primary" />
					</div>
				{:else}
					{#key `${data.anime.id}:${data.episodeNumber}:${selectedServer?.serverId ?? ''}`}
						<Player
							playerId={`anime-${data.anime.id}-ep-${data.episodeNumber}`}
							info={streamInfo}
							{nextEpisodeUrl}
							updateLibrary={async () => {
								if (!appState.user) return;
								if (data.libraryEntry && data.episodeNumber <= data.libraryEntry.watchedEpisodes)
									return;

								const id = toast.loading('Updating library status...');
								try {
									await apiClient.PUT('/library/{animeID}', {
										params: { path: { animeID: data.anime.id } },
										body: { watchedEpisodes: data.episodeNumber, status: 'watching' },
									});

									toast.success('Library status updated');
								} catch {
									toast.error('Failed to update library status');
								} finally {
									toast.dismiss(id);
								}
							}}
						/>
					{/key}
				{/if}
			</div>

			<div class="space-y-4">
				<div class="flex items-center justify-between rounded-lg border bg-card p-4">
					<div class="flex items-center gap-2">
						<Button
							size="sm"
							variant="outline"
							disabled={!hasPrevEpisode}
							onclick={() => changeEpisode(data.episodeNumber - 1)}
						>
							<ChevronLeft class="mr-1 h-4 w-4" />
							Previous
						</Button>

						<Button
							size="sm"
							variant="outline"
							disabled={!hasNextEpisode}
							onclick={() => changeEpisode(data.episodeNumber + 1)}
						>
							Next
							<ChevronRight class="ml-1 h-4 w-4" />
						</Button>
					</div>

					<Button
						size="sm"
						variant="ghost"
						onclick={() => (showEpisodesSidebar = !showEpisodesSidebar)}
						class="hidden lg:inline-flex"
					>
						{showEpisodesSidebar ? 'Hide' : 'Show'} Episodes
					</Button>
				</div>

				{#if Object.keys(groupedServers).length > 0}
					<div class="space-y-3 rounded-lg border bg-card p-4">
						<h3 class="text-base font-semibold">Servers</h3>
						{#each Object.entries(groupedServers) as [type, servers], i (i)}
							<div class="space-y-2">
								<p class="text-sm font-medium text-muted-foreground uppercase">{type}</p>
								<div class="grid grid-cols-2 gap-2 sm:grid-cols-3 md:grid-cols-4">
									{#each servers as server (server.serverId)}
										<Button
											variant={selectedServer?.serverId === server.serverId &&
											selectedServer.type === server.type
												? 'default'
												: 'outline'}
											size="sm"
											onclick={() => selectServer(server)}
											class="justify-start"
										>
											<Server class="mr-2 h-3 w-3" />
											<span class="truncate capitalize">{server.serverName}</span>
										</Button>
									{/each}
								</div>
							</div>
						{/each}
					</div>
				{/if}

				<div class="rounded-lg border bg-card p-4">
					<div class="flex items-start gap-4">
						<img
							src={data.anime.imageUrl}
							alt={data.anime.jname || data.anime.ename}
							class="h-24 w-16 rounded object-cover"
						/>
						<div class="flex-1 space-y-2">
							<div>
								<h2 class="line-clamp-1 text-lg font-bold">
									{data.anime.jname || data.anime.ename}
								</h2>
								{#if data.anime.ename && data.anime.jname}
									<p class="line-clamp-1 text-sm text-muted-foreground">{data.anime.ename}</p>
								{/if}
							</div>
							<div class="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
								{#if data.anime.metadata?.mean}
									<div class="flex items-center gap-1">
										<Star class="h-3 w-3 fill-yellow-500 text-yellow-500" />
										<span>{data.anime.metadata.mean.toFixed(1)}</span>
									</div>
									<span>•</span>
								{/if}
								<span class="capitalize">{data.anime.season} {data.anime.seasonYear}</span>
								<span>•</span>
								<span>{data.anime.metadata?.totalEpisodes || data.anime.lastEpisode} Episodes</span>
							</div>
							<div class="mt-2 flex flex-wrap gap-2">
								<Button href="/anime/{data.anime.id}" variant="outline">
									<Info class="mr-2 h-4 w-4" />
									View Details
								</Button>
								<LibraryBtn animeId={data.anime.id} libraryEntry={data.libraryEntry ?? null} />
							</div>
						</div>
					</div>
				</div>

				<div class="rounded-lg border bg-card p-4">
					<h3 class="text-base font-bold">
						Episode {data.episodeNumber}: {data.currentEpisode.title ||
							`Episode ${data.episodeNumber}`}
					</h3>
					{#if data.currentEpisode.isFiller}
						<span class="mt-2 inline-block rounded bg-muted px-2 py-1 text-xs font-medium">
							Filler Episode
						</span>
					{/if}
				</div>
			</div>
		</div>

		{#if showEpisodesSidebar}
			<div class="hidden lg:block">
				<div class="sticky top-20">
					<div class="rounded-lg border bg-card">
						<div class="border-b p-4">
							<div class="mb-3 flex items-center justify-between">
								<h3 class="font-semibold">Episodes</h3>
								<span class="text-sm text-muted-foreground">
									{data.episodes.length} Total
								</span>
							</div>
							<Input
								type="text"
								placeholder="Search episodes..."
								bind:value={episodesSearch}
								class="h-9"
							/>
						</div>

						<div class="max-h-[calc(100vh-200px)] overflow-y-auto p-2">
							<div class="space-y-1">
								{#each filteredEpisodes as episode (episode.id)}
									{@const isActive = episode.number === data.episodeNumber}
									<button
										onclick={() => changeEpisode(episode.number)}
										class={cn(
											'group flex w-full items-center gap-3 rounded-md p-2.5 text-left transition-colors',
											isActive
												? 'bg-primary text-primary-foreground'
												: 'hover:bg-accent hover:text-accent-foreground',
										)}
									>
										<div
											class={cn(
												'flex h-7 w-7 flex-shrink-0 items-center justify-center rounded text-xs font-bold',
												isActive ? 'bg-primary-foreground/20' : 'bg-muted',
											)}
										>
											{episode.number}
										</div>
										<div class="min-w-0 flex-1">
											<p class="line-clamp-1 text-sm">
												{episode.title || `Episode ${episode.number}`}
											</p>
											{#if episode.isFiller}
												<span class="text-xs opacity-70">Filler</span>
											{/if}
										</div>
										{#if isActive}
											<Play class="h-3 w-3 flex-shrink-0" />
										{/if}
									</button>
								{/each}
							</div>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>

	<div class="border-t bg-card p-4 lg:hidden">
		<div class="mb-3 flex items-center justify-between">
			<h3 class="font-semibold">Episodes</h3>
			<span class="text-sm text-muted-foreground">
				{data.episodes.length} Total
			</span>
		</div>
		<Input
			type="text"
			placeholder="Search episodes..."
			bind:value={episodesSearch}
			class="mb-3 h-9"
		/>
		<div class="max-h-64 overflow-y-auto rounded-lg border">
			<div class="space-y-1 p-2">
				{#each filteredEpisodes as episode (episode.id)}
					{@const isActive = episode.number === data.episodeNumber}
					<button
						onclick={() => changeEpisode(episode.number)}
						class={cn(
							'group flex w-full items-center gap-3 rounded-md p-2.5 text-left transition-colors',
							isActive
								? 'bg-primary text-primary-foreground'
								: 'hover:bg-accent hover:text-accent-foreground',
						)}
					>
						<div
							class={cn(
								'flex h-7 w-7 flex-shrink-0 items-center justify-center rounded text-xs font-bold',
								isActive ? 'bg-primary-foreground/20' : 'bg-muted',
							)}
						>
							{episode.number}
						</div>
						<div class="min-w-0 flex-1">
							<p class="line-clamp-1 text-sm">
								{episode.title || `Episode ${episode.number}`}
							</p>
						</div>
						{#if isActive}
							<Play class="h-3 w-3 flex-shrink-0" />
						{/if}
					</button>
				{/each}
			</div>
		</div>
	</div>
</div>
