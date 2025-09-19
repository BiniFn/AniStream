<script lang="ts">
	import Footer from '$lib/components/layout/footer.svelte';
	import NavBar from '$lib/components/layout/nav-bar.svelte';
	import Sonner from '$lib/components/ui/sonner/sonner.svelte';
	import { setSettings, setUser } from '$lib/context/state.svelte';
	import '../app.css';
	import type { LayoutProps } from './$types';

	let { children, data }: LayoutProps = $props();

	$effect(() => {
		setUser(data.user);
		setSettings(data.settings);
	});
</script>

<NavBar />
<div class="flex min-h-screen flex-col">
	<div class="flex-1 pb-4">
		{@render children?.()}
	</div>
	<Footer />
</div>
<Sonner richColors />
