<script lang="ts">
	import { appState } from '$lib/context/state.svelte';
	import { Heart } from 'lucide-svelte';

	let isLoggedIn = $derived(appState.user != null);
</script>

<footer class="border-t border-border bg-background">
	<div class="container mx-auto px-6 py-12">
		<div class="grid grid-cols-1 gap-8 md:grid-cols-4">
			<div class="md:col-span-1">
				<h3 class="mb-4 font-serif text-2xl font-bold tracking-tight text-primary uppercase">
					Aniways
				</h3>
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
					<li>
						<a href="/recent" class="text-muted-foreground transition-colors hover:text-primary">
							Recent
						</a>
					</li>
					{#if isLoggedIn}
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
					{#if isLoggedIn}
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
					{/if}
				</ul>
			</div>

			<div>
				<h4 class="mb-4 text-sm font-semibold tracking-wide text-foreground uppercase">More</h4>
				<ul class="space-y-3 text-sm">
					<li>
						<a href="/about" class="text-muted-foreground transition-colors hover:text-primary">
							About
						</a>
					</li>
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
						<a href="/status" class="text-muted-foreground transition-colors hover:text-primary">
							Status
						</a>
					</li>
				</ul>
			</div>
		</div>

		<div
			class="mt-12 flex flex-col items-center justify-between border-t border-border pt-8 md:flex-row"
		>
			<p class="text-sm text-muted-foreground">
				© {new Date().getFullYear()} Aniways. Open source project.
			</p>
			<div class="mt-4 flex items-center gap-1 text-sm text-muted-foreground md:mt-0">
				<span>Made with</span>
				<Heart class="h-4 w-4 fill-red-500 text-red-500" />
				<span>for anime fans worldwide</span>
			</div>
		</div>
	</div>
</footer>
