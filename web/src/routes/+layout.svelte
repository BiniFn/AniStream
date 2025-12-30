<script lang="ts">
	import DesktopNavBar from '$lib/components/layout/desktop-nav-bar.svelte';
	import Footer from '$lib/components/layout/footer.svelte';
	import NavBar from '$lib/components/layout/nav-bar.svelte';
	import TopLoader from '$lib/components/layout/top-loader.svelte';
	import Sonner from '$lib/components/ui/sonner/sonner.svelte';
	import { setLayoutStateContext } from '$lib/context/layout.svelte';
	import { setAppStateContext } from '$lib/context/state.svelte';
	import { isElectron } from '$lib/hooks/is-electron';
	import { getFontUrlsForTheme } from '$lib/themes';
	import { cn } from '$lib/utils';
	import '../app.css';
	import type { LayoutProps } from './$types';

	let { children, data }: LayoutProps = $props();
	const appState = setAppStateContext(data.user, data.settings);
	const layoutState = setLayoutStateContext();
	const inElectron = isElectron();

	let theme = $derived(appState.settings?.theme);

	$effect(() => {
		if (!theme) return;
		document.documentElement.className = cn('dark', theme.className);
	});

	$effect(() => {
		appState.setUser(data.user);
		appState.setSettings(data.settings);
	});
</script>

<svelte:head>
	<link rel="preconnect" href="https://fonts.googleapis.com" />
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
	<link rel="stylesheet" href={getFontUrlsForTheme(theme?.className ?? '')} />
</svelte:head>

<TopLoader />
<div class="sticky top-0 z-50" {@attach layoutState.setHeight('navbar')}>
	{#if inElectron}
		<DesktopNavBar />
	{:else}
		<NavBar />
	{/if}
</div>
<div class="flex min-h-screen flex-col">
	<div class="flex-1 pb-4">
		{@render children?.()}
	</div>
	<Footer />
</div>
<Sonner richColors />
