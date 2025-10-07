<script lang="ts">
	import { goto } from '$app/navigation';
	import { LoaderCircle } from 'lucide-svelte';
	import type { PageProps } from './$types';

	let props: PageProps = $props();

	$effect(() => {
		props.data.stream.then((data) => {
			if (!data) {
				goto('/anime/random', { replaceState: true });
				return;
			}

			goto(`/anime/${data?.id}`, { replaceState: true });
		});
	});
</script>

<div class="flex min-h-screen items-center justify-center">
	<div class="text-center">
		<LoaderCircle class="mx-auto mb-4 h-12 w-12 animate-spin text-primary" />
		<h1 class="text-2xl font-bold">Finding random anime...</h1>
		<p class="text-muted-foreground">Redirecting you to a random anime</p>
	</div>
</div>
