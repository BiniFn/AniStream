<script lang="ts">
	import Footer from '$lib/components/layout/footer.svelte';
	import NavBar from '$lib/components/layout/nav-bar.svelte';
	import ImportPolling from '$lib/components/settings/import-polling.svelte';
	import Sync from '$lib/components/settings/sync.svelte';
	import Sonner from '$lib/components/ui/sonner/sonner.svelte';
	import { layoutState } from '$lib/context/layout.svelte';
	import { setSettings, setUser } from '$lib/context/state.svelte';
	import '../app.css';
	import type { LayoutProps } from './$types';

	let { children, data }: LayoutProps = $props();

	let navbarElement: HTMLElement;

	setUser(data.user);
	setSettings(data.settings);

	$effect(() => {
		setUser(data.user);
		setSettings(data.settings);
	});

	const updateNavbarHeight = () => {
		if (navbarElement) {
			const height = navbarElement.getBoundingClientRect().height;
			layoutState.navbarHeight = height;
		}
	};

	$effect(() => {
		let resizeObserver: ResizeObserver;
		if (navbarElement) {
			resizeObserver = new ResizeObserver(() => {
				updateNavbarHeight();
			});
			resizeObserver.observe(navbarElement);
		}

		updateNavbarHeight();

		window.addEventListener('resize', updateNavbarHeight);
		return () => {
			window.removeEventListener('resize', updateNavbarHeight);
			if (resizeObserver && navbarElement) {
				resizeObserver.unobserve(navbarElement);
			}
		};
	});
</script>

<div class="sticky top-0 z-50" bind:this={navbarElement}>
	<NavBar />
</div>
<div class="flex min-h-screen flex-col">
	<div class="flex-1 pb-4">
		{@render children?.()}
	</div>
	<Footer />
</div>
<Sonner richColors />

<!-- script only components -->
<ImportPolling />
<Sync />
