<script lang="ts">
	import { Heart } from 'lucide-svelte';
</script>

<footer class="border-t border-border bg-background">
	<div class="container mx-auto px-6 py-12">
		<div class="grid grid-cols-1 gap-8 md:grid-cols-4">
			<!-- Brand Section -->
			<div class="md:col-span-1">
				<h3 class="mb-4 text-2xl font-bold tracking-tight text-primary uppercase font-serif">
					Aniways
				</h3>
				<p class="text-sm text-muted-foreground leading-relaxed">
					Your ultimate destination for anime streaming. Discover, watch, and track your favorite series and movies.
				</p>
			</div>

			<!-- Browse Section -->
			<div>
				<h4 class="mb-4 text-sm font-semibold uppercase tracking-wide text-foreground">Browse</h4>
				<ul class="space-y-3 text-sm">
					<li>
						<a href="/catalog" class="text-muted-foreground hover:text-primary transition-colors">
							All Anime
						</a>
					</li>
					<li>
						<a href="/genres" class="text-muted-foreground hover:text-primary transition-colors">
							Genres
						</a>
					</li>
					<li>
						<a href="/recent" class="text-muted-foreground hover:text-primary transition-colors">
							Recently Updated
						</a>
					</li>
				</ul>
			</div>

			<!-- Account Section -->
			<div>
				<h4 class="mb-4 text-sm font-semibold uppercase tracking-wide text-foreground">Account</h4>
				<ul class="space-y-3 text-sm">
					<li>
						<a href="/login" class="text-muted-foreground hover:text-primary transition-colors">
							Sign In
						</a>
					</li>
					<li>
						<a href="/register" class="text-muted-foreground hover:text-primary transition-colors">
							Create Account
						</a>
					</li>
					<li>
						<a href="/my-list" class="text-muted-foreground hover:text-primary transition-colors">
							My List
						</a>
					</li>
					<li>
						<a href="/settings" class="text-muted-foreground hover:text-primary transition-colors">
							Settings
						</a>
					</li>
				</ul>
			</div>

			<!-- Support Section -->
			<div>
				<h4 class="mb-4 text-sm font-semibold uppercase tracking-wide text-foreground">Support</h4>
				<ul class="space-y-3 text-sm">
					<li>
						<a href="/help" class="text-muted-foreground hover:text-primary transition-colors">
							Help Center
						</a>
					</li>
					<li>
						<a href="/contact" class="text-muted-foreground hover:text-primary transition-colors">
							Contact Us
						</a>
					</li>
					<li>
						<a href="/privacy" class="text-muted-foreground hover:text-primary transition-colors">
							Privacy Policy
						</a>
					</li>
					<li>
						<a href="/terms" class="text-muted-foreground hover:text-primary transition-colors">
							Terms of Service
						</a>
					</li>
				</ul>
			</div>
		</div>

		<!-- Bottom Section -->
		<div class="mt-12 flex flex-col items-center justify-between border-t border-border pt-8 md:flex-row">
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