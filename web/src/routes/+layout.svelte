<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar';
	import '../app.css';
	import NavBar from './nav-bar.svelte';

	let { children } = $props();
</script>

<Sidebar.Provider>
	<NavBar />
	<div class="h-full w-full p-4">
		{@render children?.()}
	</div>
</Sidebar.Provider>
