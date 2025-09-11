<script lang="ts">
	import Footer from '$lib/components/footer.svelte';
	import NavBar from '$lib/components/nav-bar.svelte';
	import Sonner from '$lib/components/ui/sonner/sonner.svelte';
	import '../app.css';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
</script>

<NavBar />
<div class="flex min-h-screen flex-col">
	<div class="flex-1 pb-4">
		{@render children?.()}
	</div>
	<Footer />
</div>
<Sonner richColors />
