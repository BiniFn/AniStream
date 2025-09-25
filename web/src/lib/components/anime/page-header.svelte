<script lang="ts">
	import { buttonVariants } from '$lib/components/ui/button';
	import { layoutState } from '$lib/context/layout.svelte';
	import { cn } from '$lib/utils';
	import { Funnel } from 'lucide-svelte';
	import type { Snippet } from 'svelte';
	import SearchBar from './search-bar.svelte';
	import SortControls from './sort-controls.svelte';
	import ViewModeToggle from './view-mode-toggle.svelte';

	interface SortOption {
		value: string;
		label: string;
	}

	interface Props {
		title: string;
		description?: string;
		searchQuery: string;
		sortBy: string;
		sortOrder: 'asc' | 'desc';
		sortOptions: SortOption[];
		viewMode: 'grid' | 'list';
		totalFilters: number;
		showMobileFilters?: boolean;
		pageInfo?: {
			currentPage: number;
			totalPages: number;
		};
		onSearchChange: (value: string) => void;
		onSortChange: (sortBy: string, sortOrder: 'asc' | 'desc') => void;
		onViewModeChange: (mode: 'grid' | 'list') => void;
		onMobileFiltersToggle?: () => void;
		children?: Snippet;
	}

	let {
		title,
		description,
		searchQuery = $bindable(),
		sortBy = $bindable(),
		sortOrder = $bindable(),
		sortOptions,
		viewMode = $bindable(),
		totalFilters,
		showMobileFilters = true,
		pageInfo,
		onSearchChange,
		onSortChange,
		onViewModeChange,
		onMobileFiltersToggle,
		children,
	}: Props = $props();

	let headerElement: HTMLElement;

	const updateHeaderHeight = () => {
		if (headerElement) {
			const height = headerElement.getBoundingClientRect().height;
			layoutState.headerHeight = height;
		}
	};

	$effect(() => {
		let resizeObserver: ResizeObserver;
		if (headerElement) {
			resizeObserver = new ResizeObserver(() => {
				updateHeaderHeight();
			});
			resizeObserver.observe(headerElement);
		}

		updateHeaderHeight();

		window.addEventListener('resize', updateHeaderHeight);
		return () => {
			window.removeEventListener('resize', updateHeaderHeight);
			if (resizeObserver && headerElement) {
				resizeObserver.unobserve(headerElement);
			}
		};
	});
</script>

<div
	bind:this={headerElement}
	class="sticky z-30 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
	style="top: {layoutState.navbarHeight}px"
>
	<div class="container mx-auto px-3 py-3 sm:px-4 sm:py-4">
		<div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between lg:gap-4">
			<div class="flex items-center justify-between gap-2 lg:hidden">
				<div class="flex-1">
					<h1 class="text-lg font-bold">{title}</h1>
					{#if pageInfo}
						<p class="text-xs text-muted-foreground">
							Page {pageInfo.currentPage}/{pageInfo.totalPages}
						</p>
					{:else if description}
						<p class="text-xs text-muted-foreground">{description}</p>
					{/if}
				</div>
				<div class="flex items-center gap-2">
					<ViewModeToggle {viewMode} {onViewModeChange} size="sm" />
					{#if showMobileFilters && onMobileFiltersToggle}
						<button
							onclick={onMobileFiltersToggle}
							class={cn(buttonVariants({ variant: 'outline', size: 'icon' }), 'relative h-9 w-9')}
						>
							<Funnel class="h-4 w-4" />
							{#if totalFilters > 0}
								<span
									class="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[10px] text-primary-foreground"
								>
									{totalFilters}
								</span>
							{/if}
						</button>
					{/if}
				</div>
			</div>

			<div class="hidden lg:block">
				<h1 class="text-2xl font-bold tracking-tight">{title}</h1>
				{#if pageInfo}
					<p class="text-sm text-muted-foreground">
						Page {pageInfo.currentPage} of {pageInfo.totalPages}
					</p>
				{:else if description}
					<p class="text-sm text-muted-foreground">{description}</p>
				{/if}
			</div>

			<div class="flex flex-col gap-2 lg:flex-row lg:items-center lg:gap-3">
				<div class="flex flex-col gap-2 lg:hidden">
					<SearchBar bind:value={searchQuery} onInput={onSearchChange} size="sm" />
					<div class="flex items-center gap-1">
						<SortControls
							{sortBy}
							{sortOrder}
							{sortOptions}
							{onSortChange}
							selectClass="h-9 flex-1 text-sm"
						/>
					</div>
				</div>

				<div class="hidden lg:flex lg:items-center lg:gap-3">
					<SearchBar bind:value={searchQuery} onInput={onSearchChange} class="w-80" />
					<SortControls {sortBy} {sortOrder} {sortOptions} {onSortChange} selectClass="w-[172px]" />
					<ViewModeToggle {viewMode} {onViewModeChange} />
				</div>
			</div>
		</div>

		{#if children}
			{@render children()}
		{/if}
	</div>
</div>
