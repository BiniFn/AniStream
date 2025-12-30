<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import Button from '$lib/components/ui/button/button.svelte';
	import {
		type UpdateStatus,
		onUpdateStatus,
		getUpdateStatus,
		startUpdate,
		quitAndInstall,
		getAppVersion,
	} from '$lib/hooks/is-electron';
	import { Download, RefreshCw, Sparkles } from 'lucide-svelte';

	let updateStatus = $state<UpdateStatus | null>(null);
	let currentVersion = $state<string | null>(null);
	let showModal = $state(false);
	let dismissedVersion = $state<string | null>(null);

	$effect(() => {
		getAppVersion().then((version) => {
			currentVersion = version;
		});

		getUpdateStatus().then((status) => {
			updateStatus = status;
			handleStatusChange(status);
		});

		onUpdateStatus((status) => {
			updateStatus = status;
			handleStatusChange(status);
		});
	});

	function handleStatusChange(status: UpdateStatus | null) {
		if (!status) return;

		if (status.status === 'available') {
			// Show modal for new update (if not dismissed)
			if (dismissedVersion !== status.version) {
				showModal = true;
			}
		} else if (status.status === 'downloaded') {
			// Always show modal when update is ready
			showModal = true;
		}
	}

	function handleSkip() {
		showModal = false;
		if (updateStatus?.status === 'available') {
			dismissedVersion = updateStatus.version;
		}
	}

	function handleStartUpdate() {
		startUpdate();
		showModal = false;
	}

	function handleRestart() {
		quitAndInstall();
	}

	const isDownloading = $derived(updateStatus?.status === 'downloading');
	const isReadyToInstall = $derived(updateStatus?.status === 'downloaded');
	const updateVersion = $derived(
		updateStatus?.status === 'available' || updateStatus?.status === 'downloaded'
			? updateStatus.version
			: null,
	);
</script>

<Dialog.Root bind:open={showModal}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title class="flex items-center gap-2">
				{#if isReadyToInstall}
					<Sparkles class="h-5 w-5 text-primary" />
					Update Ready
				{:else}
					<Download class="h-5 w-5 text-primary" />
					Update Available
				{/if}
			</Dialog.Title>
			<Dialog.Description>
				{#if isReadyToInstall}
					Version {updateVersion} has been downloaded and is ready to install.
				{:else}
					A new version of Aniways is available!
				{/if}
			</Dialog.Description>
		</Dialog.Header>

		{#if !isReadyToInstall}
			<div class="py-4">
				<div class="rounded-lg border border-border bg-muted/50 p-4">
					<div class="flex items-center justify-between">
						<div>
							<p class="text-sm text-muted-foreground">Current version</p>
							<p class="font-mono font-medium">{currentVersion}</p>
						</div>
						<div class="text-2xl text-muted-foreground">→</div>
						<div>
							<p class="text-sm text-muted-foreground">New version</p>
							<p class="font-mono font-medium text-primary">{updateVersion}</p>
						</div>
					</div>
				</div>
			</div>
		{/if}

		<Dialog.Footer class="gap-2 sm:gap-0">
			{#if isReadyToInstall}
				<Button variant="outline" onclick={() => (showModal = false)}>Later</Button>
				<Button onclick={handleRestart} class="gap-2">
					<RefreshCw class="h-4 w-4" />
					Restart Now
				</Button>
			{:else}
				<Button variant="outline" onclick={handleSkip}>Skip</Button>
				<Button onclick={handleStartUpdate} class="gap-2">
					<Download class="h-4 w-4" />
					Update
				</Button>
			{/if}
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
