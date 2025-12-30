<script lang="ts">
	import { browser } from '$app/environment';
	import type { components } from '$lib/api/openapi';
	import Button from '$lib/components/ui/button/button.svelte';
	import { isElectron } from '$lib/hooks/is-electron';
	import { isMobile } from '$lib/hooks/is-mobile';
	import { detectOS, formatFileSize, getOSDisplayName, getPlatformDisplayName } from '$lib/utils/platform';
	import { Apple, Download, Monitor, Smartphone } from 'lucide-svelte';
	import type { PageProps } from './$types';

	type DesktopVersionResponse = components['schemas']['models.DesktopVersionResponse'];
	type DesktopPlatformRelease = components['schemas']['models.DesktopPlatformRelease'];

	let { data }: PageProps = $props();

	const inElectron = browser ? isElectron() : false;
	const onMobile = browser ? isMobile() : false;
	const shouldShow = browser ? !inElectron && !onMobile : true;

	const os = browser ? detectOS() : 'unknown';

	const latestRelease = $derived(data.latestRelease as DesktopVersionResponse | undefined);

	function getPlatformPrefix(os: string): string {
		if (os === 'macos') return 'darwin';
		if (os === 'windows') return 'win32';
		if (os === 'linux') return 'linux';
		return '';
	}

	const platformPrefix = getPlatformPrefix(os);

	const currentOsDownloads = $derived.by(() => {
		if (!latestRelease?.platforms || !platformPrefix) return [];
		return latestRelease.platforms.filter((p) => p.platform?.startsWith(platformPrefix));
	});

	const otherOsDownloads = $derived.by(() => {
		if (!latestRelease?.platforms || !platformPrefix) return latestRelease?.platforms || [];
		return latestRelease.platforms.filter((p) => !p.platform?.startsWith(platformPrefix));
	});

	function getOSIcon(platform: string) {
		if (platform.startsWith('darwin')) return Apple;
		if (platform.startsWith('win32')) return Monitor;
		return Monitor;
	}

	function groupByVersion(
		releases: components['schemas']['models.DesktopReleaseResponse'][],
	): Map<string, components['schemas']['models.DesktopReleaseResponse'][]> {
		const grouped = new Map<string, components['schemas']['models.DesktopReleaseResponse'][]>();
		for (const release of releases) {
			const version = release.version || 'unknown';
			if (!grouped.has(version)) {
				grouped.set(version, []);
			}
			grouped.get(version)!.push(release);
		}
		return grouped;
	}

	const groupedReleases = $derived(groupByVersion(data.allReleases));
</script>

<svelte:head>
	<title>Download Desktop App - Aniways</title>
	<meta
		name="description"
		content="Download the Aniways desktop app for the best anime watching experience."
	/>
</svelte:head>

{#if !shouldShow}
	<div class="container mx-auto px-4 py-16">
		<div class="mx-auto max-w-lg text-center">
			<Smartphone class="mx-auto mb-6 h-16 w-16 text-muted-foreground" />
			<h1 class="mb-4 text-2xl font-bold">Desktop App Not Available</h1>
			<p class="text-muted-foreground">
				{#if inElectron}
					You're already using the desktop app!
				{:else}
					The desktop app is only available for desktop devices. Please visit this page from a
					computer to download.
				{/if}
			</p>
		</div>
	</div>
{:else if !latestRelease}
	<div class="container mx-auto px-4 py-16">
		<div class="mx-auto max-w-lg text-center">
			<Monitor class="mx-auto mb-6 h-16 w-16 text-muted-foreground" />
			<h1 class="mb-4 text-2xl font-bold">No Releases Available</h1>
			<p class="text-muted-foreground">
				The desktop app is not yet available for download. Check back soon!
			</p>
		</div>
	</div>
{:else}
	<div class="min-h-screen bg-background">
		<div
			class="relative overflow-hidden bg-gradient-to-br from-background via-background to-primary/10 py-16"
		>
			<div class="pointer-events-none absolute inset-0 overflow-hidden">
				<div
					class="absolute -top-20 -right-20 h-64 w-64 animate-pulse rounded-full bg-primary/5 blur-3xl"
				></div>
				<div
					class="absolute top-1/2 -left-20 h-48 w-48 animate-pulse rounded-full bg-secondary/10 blur-2xl"
					style="animation-delay: 1s;"
				></div>
			</div>

			<div class="container relative z-10 mx-auto px-4">
				<div class="mx-auto max-w-2xl text-center">
					<div class="mb-6 inline-flex items-center gap-2 rounded-full bg-primary/10 px-4 py-2">
						<Download class="h-4 w-4 text-primary" />
						<span class="text-sm font-semibold tracking-wider text-primary uppercase">
							Desktop App
						</span>
					</div>

					<h1 class="mb-4 text-4xl font-bold tracking-tight md:text-5xl">
						Download Aniways for {getOSDisplayName(os)}
					</h1>

					<p class="mb-8 text-lg text-muted-foreground">
						Get the best anime watching experience with our native desktop app. Direct HLS
						streaming, no ads, and seamless library sync.
					</p>

					{#if currentOsDownloads.length > 0}
						<div class="flex flex-wrap justify-center gap-4">
							{#each currentOsDownloads as download}
								<Button
									href={download.downloadUrl}
									size="lg"
									class="gap-2 px-6 py-6 text-base"
								>
									<Download class="h-5 w-5" />
									{getPlatformDisplayName(download.platform || '')}
								</Button>
							{/each}
						</div>
						<p class="mt-4 text-sm text-muted-foreground">
							Version {latestRelease.version}
						</p>
					{:else if latestRelease.platforms && latestRelease.platforms.length > 0}
						<p class="text-sm text-muted-foreground">
							Version {latestRelease.version} • Select your platform below
						</p>
					{/if}
				</div>
			</div>
		</div>

		<div class="container mx-auto px-4 py-12">
			{#if otherOsDownloads.length > 0}
				<div class="mx-auto mb-16 max-w-3xl">
					<h2 class="mb-6 text-center text-xl font-semibold">Other Platforms</h2>
					<div class="grid gap-4 sm:grid-cols-2 md:grid-cols-3">
						{#each otherOsDownloads as download}
							{@const Icon = getOSIcon(download.platform || '')}
							<a
								href={download.downloadUrl}
								class="flex items-center gap-4 rounded-lg border border-border bg-card p-4 transition-colors hover:border-primary/50 hover:bg-muted/50"
							>
								<div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-muted">
									<Icon class="h-6 w-6" />
								</div>
								<div class="min-w-0 flex-1">
									<p class="font-medium">{getPlatformDisplayName(download.platform || '')}</p>
									<p class="text-sm text-muted-foreground">
										{formatFileSize(download.fileSize || 0)}
									</p>
								</div>
								<Download class="h-5 w-5 shrink-0 text-muted-foreground" />
							</a>
						{/each}
					</div>
				</div>
			{/if}

			{#if latestRelease.releaseNotes}
				<div class="mx-auto mb-16 max-w-2xl">
					<h2 class="mb-4 text-xl font-semibold">What's New in v{latestRelease.version}</h2>
					<div class="rounded-lg border border-border bg-card p-6">
						<p class="whitespace-pre-wrap text-muted-foreground">{latestRelease.releaseNotes}</p>
					</div>
				</div>
			{/if}

			{#if groupedReleases.size > 1}
				<div class="mx-auto max-w-2xl">
					<h2 class="mb-6 text-xl font-semibold">Previous Versions</h2>
					<div class="space-y-4">
						{#each [...groupedReleases.entries()] as [version, releases]}
							{#if version !== latestRelease.version}
								<details class="group rounded-lg border border-border bg-card">
									<summary
										class="flex cursor-pointer items-center justify-between p-4 font-medium"
									>
										<span>Version {version}</span>
										<span class="text-sm text-muted-foreground">
											{releases.length} platform{releases.length > 1 ? 's' : ''}
										</span>
									</summary>
									<div class="border-t border-border p-4">
										<div class="grid gap-2 sm:grid-cols-2">
											{#each releases as release}
												<a
													href={release.downloadUrl}
													class="flex items-center gap-3 rounded-md p-2 text-sm transition-colors hover:bg-muted"
												>
													<Download class="h-4 w-4 text-muted-foreground" />
													<span>{getPlatformDisplayName(release.platform || '')}</span>
													<span class="ml-auto text-muted-foreground">
														{formatFileSize(release.fileSize || 0)}
													</span>
												</a>
											{/each}
										</div>
									</div>
								</details>
							{/if}
						{/each}
					</div>
				</div>
			{/if}
		</div>
	</div>
{/if}
