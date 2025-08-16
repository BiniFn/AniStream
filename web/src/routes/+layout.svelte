<script lang="ts">
	import '../app.css';
	import NavBar from './nav-bar.svelte';

	let { children } = $props();
</script>

<NavBar />
<div class="h-full w-full p-4">
	{@render children?.()}
</div>
