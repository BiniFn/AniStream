<script lang="ts">
	import '../app.css';
	import type { LayoutProps } from './$types';
	import NavBar from './nav-bar.svelte';

	let { children }: LayoutProps = $props();
</script>

<NavBar />
<div class="h-full w-full p-4">
	{@render children?.()}
</div>
