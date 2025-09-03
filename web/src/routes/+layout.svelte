<script lang="ts">
	import { setUser } from '$lib/auth.svelte';
	import Footer from '$lib/components/footer.svelte';
	import NavBar from '$lib/components/nav-bar.svelte';
	import '../app.css';
	import type { LayoutProps } from './$types';

	let { children, data }: LayoutProps = $props();

	if (data?.user) setUser(data?.user);
</script>

<NavBar />
<div class="flex min-h-screen flex-col">
	<div class="flex-1 p-4">
		{@render children?.()}
	</div>
	<Footer />
</div>
