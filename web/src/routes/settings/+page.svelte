<script lang="ts">
	import type { components } from '$lib/api/openapi';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import Label from '$lib/components/ui/label/label.svelte';
	import { appState } from '$lib/context/state.svelte';

	type Settings = Omit<components['schemas']['models.SettingsResponse'], 'userId'>;

	function updateSetting(key: keyof Settings, value: boolean) {
		if (!appState.settings) return;
		appState.settings[key] = value;
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>App Preferences</Card.Title>
		<Card.Description>Customize your viewing experience and app behavior</Card.Description>
	</Card.Header>
	<Card.Content class="space-y-6">
		<div class="space-y-4">
			<h3 class="text-lg font-medium">Playback Settings</h3>
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<div class="space-y-0.5">
						<Label>Auto Play Episodes</Label>
						<p class="text-sm text-muted-foreground">
							Automatically start playing episodes when opened
						</p>
					</div>
					<Button
						variant={appState.settings?.autoPlayEpisode ? 'default' : 'outline'}
						size="sm"
						onclick={() => updateSetting('autoPlayEpisode', !appState.settings?.autoPlayEpisode)}
					>
						{appState.settings?.autoPlayEpisode ? 'Enabled' : 'Disabled'}
					</Button>
				</div>

				<div class="flex items-center justify-between">
					<div class="space-y-0.5">
						<Label>Auto Next Episode</Label>
						<p class="text-sm text-muted-foreground">
							Automatically play the next episode when current one ends
						</p>
					</div>
					<Button
						variant={appState.settings?.autoNextEpisode ? 'default' : 'outline'}
						size="sm"
						onclick={() => updateSetting('autoNextEpisode', !appState.settings?.autoNextEpisode)}
					>
						{appState.settings?.autoNextEpisode ? 'Enabled' : 'Disabled'}
					</Button>
				</div>

				<div class="flex items-center justify-between">
					<div class="space-y-0.5">
						<Label>Auto Resume Episodes</Label>
						<p class="text-sm text-muted-foreground">
							Automatically resume from where you left off
						</p>
					</div>
					<Button
						variant={appState.settings?.autoResumeEpisode ? 'default' : 'outline'}
						size="sm"
						onclick={() =>
							updateSetting('autoResumeEpisode', !appState.settings?.autoResumeEpisode)}
					>
						{appState.settings?.autoResumeEpisode ? 'Enabled' : 'Disabled'}
					</Button>
				</div>

				<div class="flex items-center justify-between">
					<div class="space-y-0.5">
						<Label>Incognito Mode</Label>
						<p class="text-sm text-muted-foreground">
							Don't track viewing history or update library
						</p>
					</div>
					<Button
						variant={appState.settings?.incognitoMode ? 'default' : 'outline'}
						size="sm"
						onclick={() => updateSetting('incognitoMode', !appState.settings?.incognitoMode)}
					>
						{appState.settings?.incognitoMode ? 'Enabled' : 'Disabled'}
					</Button>
				</div>
			</div>
		</div>
	</Card.Content>
</Card.Root>

