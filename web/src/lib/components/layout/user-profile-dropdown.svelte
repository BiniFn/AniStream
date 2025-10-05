<script lang="ts">
	import type { components } from '$lib/api/openapi';
	import { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { cn } from '$lib/utils';
	import { Heart, LogOut, Settings, User } from 'lucide-svelte';
	import ProfilePicture from './profile-picture.svelte';

	interface Props {
		user: components['schemas']['models.UserResponse'];
		class?: string;
	}

	let { user, class: className }: Props = $props();
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class={cn(
			buttonVariants({ variant: 'ghost', size: 'lg' }),
			'flex items-center gap-2 hover:bg-primary/80',
			className,
		)}
	>
		<span class="hidden text-sm font-medium lg:block">{user.username}</span>
		<ProfilePicture class="size-8" />
		hevronDown class="h-4 w-4 opacity-50" />
	</DropdownMenu.Trigger>

	<DropdownMenu.Content class="w-56" align="end">
		<DropdownMenu.Label class="font-normal">
			<div class="flex flex-col space-y-1">
				<p class="text-sm leading-none font-medium">{user.username}</p>
				<p class="text-xs leading-none text-muted-foreground">{user.email}</p>
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
