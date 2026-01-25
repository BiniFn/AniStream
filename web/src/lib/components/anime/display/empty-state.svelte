<script lang="ts">
	import type { ComponentType, Snippet } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { cn } from '$lib/utils';

	interface Props {
		icon?: ComponentType | string;
		title: string;
		description?: string;
		action?: {
			label: string;
			href?: string;
			onClick?: () => void;
			variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link';
		};
		secondaryAction?: {
			label: string;
			href?: string;
			onClick?: () => void;
			variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link';
		};
		class?: string;
		children?: Snippet;
	}

	let {
		icon,
		title,
		description,
		action,
		secondaryAction,
		class: className = '',
		children,
	}: Props = $props();
</script>

<Card.Root class={cn('border-dashed', className)}>
	<Card.Content class="flex flex-col items-center justify-center py-20 text-center">
		{#if icon}
			<div class="mb-4">
				{#if typeof icon === 'string'}
					<div class="text-6xl">{icon}</div>
				{:else}
					{@const IconComponent = icon as ComponentType}
					<IconComponent class="h-16 w-16 text-muted-foreground" />
				{/if}
			</div>
		{/if}

		<h3 class="mb-2 text-xl font-semibold">{title}</h3>

		{#if description}
			<p class="mb-6 max-w-md text-muted-foreground">{description}</p>
		{/if}

		{#if children}
			<div class="mb-6">
				{@render children()}
			</div>
		{/if}

		{#if action || secondaryAction}
			<div class="flex flex-col gap-3 sm:flex-row">
				{#if action}
					<Button
						href={action.href}
						onclick={action.onClick}
						variant={action.variant || 'default'}
						class="gap-2"
					>
						{action.label}
					</Button>
				{/if}

				{#if secondaryAction}
					<Button
						href={secondaryAction.href}
						onclick={secondaryAction.onClick}
						variant={secondaryAction.variant || 'outline'}
						class="gap-2"
					>
						{secondaryAction.label}
					</Button>
				{/if}
			</div>
		{/if}
	</Card.Content>
</Card.Root>
