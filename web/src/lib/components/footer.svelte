<script lang="ts">
	import { Heart } from 'lucide-svelte';
</script>

<footer class="border-t border-border bg-background">
	<div class="container mx-auto px-6 py-12">
		<div class="grid grid-cols-1 gap-8 md:grid-cols-4">
			<div class="md:col-span-1">
				<h3 class="mb-4 font-serif text-2xl font-bold tracking-tight text-primary uppercase">
					Aniways
				</h3>
				<p class="text-sm leading-relaxed text-muted-foreground">
					Your ultimate destination for anime streaming. Discover, watch, and track your favorite
					series and movies.
				</p>
			</div>

			<div>
				<h4 class="mb-4 text-sm font-semibold tracking-wide text-foreground uppercase">Browse</h4>
				<ul class="space-y-3 text-sm">
					<li>
						<a href="/catalog" class="text-muted-foreground transition-colors hover:text-primary">
							All Anime
						</a>
					</li>
					<li>
						<a href="/genres" class="text-muted-foreground transition-colors hover:text-primary">
							Genres
						</a>
					</li>
					<li>
						<a href="/recent" class="text-muted-foreground transition-colors hover:text-primary">
							Recently Updated
						</a>
					</li>
				</ul>
			</div>

			<div>
				<h4 class="mb-4 text-sm font-semibold tracking-wide text-foreground uppercase">Account</h4>
				<ul class="space-y-3 text-sm">
					<li>
						<a href="/login" class="text-muted-foreground transition-colors hover:text-primary">
							Sign In
						</a>
					</li>
					<li>
						<a href="/register" class="text-muted-foreground transition-colors hover:text-primary">
							Create Account
						</a>
					</li>
					<li>
						<a href="/my-list" class="text-muted-foreground transition-colors hover:text-primary">
							My List
						</a>
					</li>
					<li>
						<a href="/settings" class="text-muted-foreground transition-colors hover:text-primary">
							Settings
						</a>
					</li>
				</ul>
			</div>

			<div>
				<h4 class="mb-4 text-sm font-semibold tracking-wide text-foreground uppercase">Support</h4>
				<ul class="space-y-3 text-sm">
					<li>
						<a href="/help" class="text-muted-foreground transition-colors hover:text-primary">
							Help Center
						</a>
					</li>
					<li>
						<a href="/contact" class="text-muted-foreground transition-colors hover:text-primary">
							Contact Us
						</a>
					</li>
					<li>
						<a href="/privacy" class="text-muted-foreground transition-colors hover:text-primary">
							Privacy Policy
						</a>
					</li>
					<li>
						<a href="/terms" class="text-muted-foreground transition-colors hover:text-primary">
							Terms of Service
						</a>
					</li>
				</ul>
			</div>
		</div>

		<div
			class="mt-12 flex flex-col items-center justify-between border-t border-border pt-8 md:flex-row"
		>
			<p class="text-sm text-muted-foreground">
				© {new Date().getFullYear()} Aniways. All rights reserved.
			</p>
			<div class="mt-4 flex items-center gap-1 text-sm text-muted-foreground md:mt-0">
				<span>Made with</span>
				<Heart class="h-4 w-4 fill-red-500 text-red-500" />
				<span>for anime fans worldwide</span>
			</div>
		</div>
	</div>
</footer>

