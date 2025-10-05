<script lang="ts">
	import * as Avatar from '$lib/components/ui/avatar';
	import { appState } from '$lib/context/state.svelte';
	import { cn } from '$lib/utils';

	type Props = {
		class?: string;
	};

	let { class: className }: Props = $props();
	let user = $derived(appState.user);
</script>

<Avatar.Root class={cn('size-10', className)}>
	<Avatar.Image src={user?.profilePicture} alt={user?.username} />
	<Avatar.Fallback class="bg-primary/50 text-sm font-medium">
		{user?.username
			?.split(' ')
			.map((n) => n.charAt(0))
			.join('')
			.toUpperCase()
			.slice(0, 2)}
	</Avatar.Fallback>
</Avatar.Root>
