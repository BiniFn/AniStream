<script lang="ts">
	import { ChevronDown, Heart, LogOut, Settings, User } from 'lucide-svelte';
	import { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { getAppStateContext } from '$lib/context/state.svelte';
	import { cn } from '$lib/utils';
	import ProfilePicture from './profile-picture.svelte';

	interface Props {
		class?: string;
	}

	let { class: className }: Props = $props();
	const appState = getAppStateContext();
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class={cn(
			buttonVariants({ variant: 'ghost', size: 'lg' }),
			'flex items-center gap-2 hover:bg-primary/80',
			className,
		)}
	>
		<span class="hidden text-sm font-medium lg:block">{appState.user?.username}</span>
		<ProfilePicture class="size-8" />
		<ChevronDown class="h-4 w-4 opacity-50" />
	</DropdownMenu.Trigger>

	<DropdownMenu.Content class="w-56" align="end">
		<DropdownMenu.Label class="font-normal">
			<div class="flex flex-col space-y-1">
				<p class="text-sm leading-none font-medium">{appState.user?.username}</p>
				<p class="text-xs leading-none text-muted-foreground">{appState.user?.email}</p>
			</div>
		</DropdownMenu.Label>
		<DropdownMenu.Separator />

		<DropdownMenu.Group>
			<DropdownMenu.Item href="/profile">
				<User class="mr-2 h-4 w-4" />
				<span>Profile</span>
			</DropdownMenu.Item>
			<DropdownMenu.Item href="/my-list">
				<Heart class="mr-2 h-4 w-4" />
				<span>My List</span>
			</DropdownMenu.Item>
		</DropdownMenu.Group>

		<DropdownMenu.Separator />

		<DropdownMenu.Item href="/settings">
			<Settings class="mr-2 h-4 w-4" />
			<span>Settings</span>
		</DropdownMenu.Item>

		<DropdownMenu.Separator />

		<DropdownMenu.Item href="/logout" variant="destructive">
			<LogOut class="mr-2 h-4 w-4" />
			<span>Log out</span>
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
