<script lang="ts">
	import { browser } from '$app/environment';
	import { getAppStateContext } from '$lib/context/state.svelte';
	import { isElectron } from '$lib/hooks/is-electron';
	import { isMobile } from '$lib/hooks/is-mobile';
	import { Heart } from 'lucide-svelte';
	import BrandText from './brand-text.svelte';

	const appState = getAppStateContext();
	const showDownload = browser ? !isElectron() && !isMobile() : false;
</script>

<footer class="border-t border-border bg-background">
	<div class="container mx-auto px-6 py-12">
		<div class="grid grid-cols-1 gap-8 md:grid-cols-4">
			<div class="md:col-span-1">
				<div class="mb-4">
					<BrandText size="md" variant="anime" />
				</div>
				<p class="text-sm leading-relaxed text-muted-foreground">
					Discover, track, and organize your favorite anime series. Built by fans, for fans.
				</p>
			</div>

			<div>
				<h4 class="mb-4 text-sm font-semibold tracking-wide text-foreground uppercase">Browse</h4>
				<ul class="space-y-3 text-sm">
					<li>
						<a href="/" class="text-muted-foreground transition-colors hover:text-primary">
							Home
						</a>
					</li>
					<li>
						<a href="/catalog" class="text-muted-foreground transition-colors hover:text-primary">
							Catalog
						</a>
					</li>
					<li>
						<a href="/genres" class="text-muted-foreground transition-colors hover:text-primary">
							Genres
						</a>
					</li>
					{#if appState.isLoggedIn}
						<li>
							<a href="/my-list" class="text-muted-foreground transition-colors hover:text-primary">
								My List
							</a>
						</li>
					{/if}
				</ul>
			</div>

			<div>
				<h4 class="mb-4 text-sm font-semibold tracking-wide text-foreground uppercase">Account</h4>
				<ul class="space-y-3 text-sm">
					{#if appState.isLoggedIn}
						<li>
							<a href="/profile" class="text-muted-foreground transition-colors hover:text-primary">
								Profile
							</a>
						</li>
						<li>
							<a
								href="/settings"
								class="text-muted-foreground transition-colors hover:text-primary"
							>
								Settings
							</a>
						</li>
						<li>
							<a href="/logout" class="text-muted-foreground transition-colors hover:text-primary">
								Log out
							</a>
						</li>
					{:else}
						<li>
							<a href="/login" class="text-muted-foreground transition-colors hover:text-primary">
								Sign In
							</a>
						</li>
						<li>
							<a
								href="/register"
								class="text-muted-foreground transition-colors hover:text-primary"
							>
								Register
							</a>
						</li>
						<li>
							<a
								href="/forgot-password"
								class="text-muted-foreground transition-colors hover:text-primary"
							>
								Forgot Password
							</a>
						</li>
					{/if}
				</ul>
			</div>

			<div>
				<h4 class="mb-4 text-sm font-semibold tracking-wide text-foreground uppercase">More</h4>
				<ul class="space-y-3 text-sm">
					{#if showDownload}
						<li>
							<a
								href="/download"
								class="text-muted-foreground transition-colors hover:text-primary"
							>
								Desktop App
							</a>
						</li>
					{/if}
					<li>
						<a
							href="https://github.com/coeeter/aniways"
							target="_blank"
							rel="noopener noreferrer"
							class="text-muted-foreground transition-colors hover:text-primary"
						>
							GitHub
						</a>
					</li>
					<li>
						<a
							href="https://github.com/coeeter/aniways/issues"
							target="_blank"
							rel="noopener noreferrer"
							class="text-muted-foreground transition-colors hover:text-primary"
						>
							Report Bug
						</a>
					</li>
				</ul>
			</div>
		</div>

		<div
			class="mt-12 flex flex-col items-center justify-between border-t border-border pt-8 md:flex-row"
		>
			<div class="text-center md:text-left">
				<p class="text-sm text-muted-foreground">
					© {new Date().getFullYear()} Aniways. Open source project.
				</p>
				<p class="mt-2 text-xs text-muted-foreground">
					AniWays is for educational purposes only. All content is sourced from third-party sites.
					We do not host or store any copyrighted material. Use at your own risk.
				</p>
			</div>
			<div class="mt-4 flex items-center gap-1 text-sm text-muted-foreground md:mt-0">
				<span>Made with</span>
				<Heart class="h-4 w-4 fill-red-500 text-red-500" />
				<span>for anime fans worldwide</span>
			</div>
		</div>
	</div>
</footer>
