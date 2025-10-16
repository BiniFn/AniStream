<script lang="ts">
	import Footer from '$lib/components/layout/footer.svelte';
	import NavBar from '$lib/components/layout/nav-bar.svelte';
	import Sonner from '$lib/components/ui/sonner/sonner.svelte';
	import { setLayoutStateContext } from '$lib/context/layout.svelte';
	import { setAppStateContext } from '$lib/context/state.svelte';
	import '../app.css';
	import type { LayoutProps } from './$types';

	let { children, data }: LayoutProps = $props();
	const appState = setAppStateContext(data.user, data.settings);
	const layoutState = setLayoutStateContext();

	$effect(() => {
		appState.setUser(data.user);
		appState.setSettings(data.settings);
	});
</script>

<div class="sticky top-0 z-50" {@attach layoutState.setHeight('navbar')}>
	<NavBar />
</div>
<div class="flex min-h-screen flex-col">
	<div class="flex-1 pb-4">
		{@render children?.()}
	</div>
	<Footer />
</div>
<Sonner richColors />
