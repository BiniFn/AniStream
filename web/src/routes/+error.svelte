<script lang="ts">
	import { type } from 'arktype';
	import { useSearchParams } from 'runed/kit';
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import { errorMessages } from '$lib/errors';

	const ErrorSchema = type({
		type: 'string = ""',
		from: 'string = ""',
	});

	let params = useSearchParams(ErrorSchema);

	let errorMessage = $derived.by(() => {
		if (params.type in errorMessages) {
			return errorMessages[params.type as keyof typeof errorMessages];
		}

		if (page.status === 404) {
			return 'The requested resource was not found.';
		}

		return 'Something went wrong. Please try again.';
	});

	let from = $derived(decodeURIComponent(params.from));
</script>

<svelte:head>
	<title>Error - Aniways</title>
	<meta name="description" content={errorMessage} />
</svelte:head>

<main class="flex min-h-screen items-center justify-center p-4 md:items-start md:pt-16">
	<div class="flex w-full flex-col items-center gap-4">
		<h1 class="text-3xl font-bold">An Error Occurred</h1>
		<img src="/error.gif" alt="Error Illustration" class="mx-auto max-w-64" />
		<p class="mt-4 text-center text-lg font-medium text-muted-foreground">{errorMessage}</p>
		<div class="flex gap-4">
			<Button onclick={() => history.back()}>Go Back</Button>
			{#if from}
				<Button variant="secondary" href={from} data-sveltekit-replacestate>Retry</Button>
			{/if}
		</div>
	</div>
</main>
