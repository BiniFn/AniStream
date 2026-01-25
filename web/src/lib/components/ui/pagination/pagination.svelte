<script lang="ts">
	import { ChevronLeft, ChevronRight } from 'lucide-svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';

	interface Props {
		totalPages: number;
		currentPage?: number;
		onPageChange?: (page: number) => void;
	}

	let { totalPages, currentPage, onPageChange }: Props = $props();

	const derivedCurrentPage = $derived(
		currentPage ?? parseInt(page.url.searchParams.get('page') || '1', 10),
	);

	const defaultPageChange = async (newPage: number) => {
		const url = new URL(page.url);
		url.searchParams.set('page', newPage.toString());
		await goto(url.toString());
		window.scrollTo({ top: 0, behavior: 'smooth' });
	};

	const handlePageChange = onPageChange ?? defaultPageChange;
</script>

{#if totalPages > 1}
	<div class="mt-12 flex flex-wrap items-center justify-center gap-2">
		<Button
			variant="outline"
			size="icon"
			disabled={derivedCurrentPage === 1}
			onclick={() => handlePageChange(derivedCurrentPage - 1)}
		>
			<ChevronLeft class="h-4 w-4" />
		</Button>

		{#if totalPages <= 7}
			{#each Array(totalPages) as _, i (i)}
				<Button
					variant={derivedCurrentPage === i + 1 ? 'default' : 'outline'}
					size="icon"
					onclick={() => handlePageChange(i + 1)}
				>
					{i + 1}
				</Button>
			{/each}
		{:else}
			{#if derivedCurrentPage > 3}
				<Button variant="outline" size="icon" onclick={() => handlePageChange(1)}>1</Button>
				{#if derivedCurrentPage > 4}
					<span class="px-2 text-muted-foreground">...</span>
				{/if}
			{/if}

			{#each Array(5) as _, i (i)}
				{@const pageNum =
					derivedCurrentPage <= 3
						? i + 1
						: derivedCurrentPage >= totalPages - 2
							? totalPages - 4 + i
							: derivedCurrentPage - 2 + i}
				{#if pageNum > 0 && pageNum <= totalPages}
					<Button
						variant={derivedCurrentPage === pageNum ? 'default' : 'outline'}
						size="icon"
						onclick={() => handlePageChange(pageNum)}
					>
						{pageNum}
					</Button>
				{/if}
			{/each}

			{#if derivedCurrentPage < totalPages - 2}
				{#if derivedCurrentPage < totalPages - 3}
					<span class="px-2 text-muted-foreground">...</span>
				{/if}
				<Button variant="outline" size="icon" onclick={() => handlePageChange(totalPages)}>
					{totalPages}
				</Button>
			{/if}
		{/if}

		<Button
			variant="outline"
			size="icon"
			disabled={derivedCurrentPage === totalPages}
			onclick={() => handlePageChange(derivedCurrentPage + 1)}
		>
			<ChevronRight class="h-4 w-4" />
		</Button>
	</div>
{/if}
