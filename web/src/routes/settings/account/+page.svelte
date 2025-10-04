<script lang="ts">
	import { apiClient } from '$lib/api/client';
	import type { components } from '$lib/api/openapi';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import Input from '$lib/components/ui/input/input.svelte';
	import { cn } from '$lib/utils';
	import { type } from 'arktype';
	import { Eye, EyeOff, LoaderCircle, Shield } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import { defaults, setError, superForm } from 'sveltekit-superforms';
	import { arktype, arktypeClient } from 'sveltekit-superforms/adapters';

	const passwordFormSchema = type({
		oldPassword: type('string>=1').describe('current password'),
		newPassword: type('string>=1').describe('new password'),
		confirmPassword: type('string>=1').describe('password confirmation'),
	});

	const sf = superForm(defaults(arktype(passwordFormSchema)), {
		SPA: true,
		resetForm: true,
		validators: arktypeClient(passwordFormSchema),
		onUpdate: async ({ form, cancel }) => {
			if (!form.valid) return;

			if (form.data.newPassword !== form.data.confirmPassword) {
				setError(form, 'confirmPassword', ['Passwords do not match']);
				toast.error('Passwords do not match');
				cancel();
				return;
			}

			try {
				const res = await apiClient.PUT('/users/password', {
					body: {
						oldPassword: form.data.oldPassword,
						newPassword: form.data.newPassword,
					},
				});

				if (res.response.status === 200) {
					toast.success('Password changed successfully');
					return;
				}

				if (res.response.status === 400) {
					const error = res.error as components['schemas']['models.ValidationErrorResponse'];
					if (error?.details) {
						Object.entries(error.details).forEach(([field, messages]) => {
							setError(form, field as 'oldPassword' | 'newPassword', messages);
						});
						toast.error('Please fix the errors in the form and try again.');
						cancel();
						return;
					}
				}

				toast.error('Failed to change password. Please check your current password.');
				cancel();
			} catch {
				toast.error('An error occurred. Please try again.');
				cancel();
			}
		},
	});

	const { form, enhance, submitting, errors } = sf;

	let showCurrentPassword = $state(false);
	let showNewPassword = $state(false);
	let showConfirmPassword = $state(false);
</script>

<Card.Root>
	<Card.Header>
		<Card.Title class="flex items-center gap-2">
			<Shield class="h-5 w-5" />
			Change Password
		</Card.Title>
		<Card.Description>Update your account password for better security</Card.Description>
	</Card.Header>
	<Card.Content>
		<form method="POST" use:enhance class="space-y-4">
			<Form.Field form={sf} name="oldPassword">
				<Form.Control>
					{#snippet children({ props })}
						<Form.Label>Current Password</Form.Label>
						<div class="relative">
							<Input
								{...props}
								type={showCurrentPassword ? 'text' : 'password'}
								bind:value={$form.oldPassword}
								placeholder="Enter your current password"
								class={cn(
									$errors.oldPassword &&
										$errors.oldPassword.length > 0 &&
										'border-destructive focus-visible:ring-destructive',
								)}
							/>
							<Button
								variant="ghost"
								size="sm"
								type="button"
								class="absolute top-0 right-0 h-full px-3 py-2 hover:bg-transparent"
								onclick={() => (showCurrentPassword = !showCurrentPassword)}
							>
								{#if showCurrentPassword}
									<EyeOff class="h-4 w-4" />
								{:else}
									<Eye class="h-4 w-4" />
								{/if}
							</Button>
						</div>
					{/snippet}
				</Form.Control>
				<Form.Description>Enter your current password to verify your identity</Form.Description>
				<Form.FieldErrors />
			</Form.Field>

			<Form.Field form={sf} name="newPassword">
				<Form.Control>
					{#snippet children({ props })}
						<Form.Label>New Password</Form.Label>
						<div class="relative">
							<Input
								{...props}
								type={showNewPassword ? 'text' : 'password'}
								bind:value={$form.newPassword}
								placeholder="Enter your new password"
								class={cn(
									$errors.newPassword &&
										$errors.newPassword.length > 0 &&
										'border-destructive focus-visible:ring-destructive',
								)}
							/>
							<Button
								variant="ghost"
								size="sm"
								type="button"
								class="absolute top-0 right-0 h-full px-3 py-2 hover:bg-transparent"
								onclick={() => (showNewPassword = !showNewPassword)}
							>
								{#if showNewPassword}
									<EyeOff class="h-4 w-4" />
								{:else}
									<Eye class="h-4 w-4" />
								{/if}
							</Button>
						</div>
					{/snippet}
				</Form.Control>
				<Form.Description>Choose a strong password with at least 8 characters</Form.Description>
				<Form.FieldErrors />
			</Form.Field>

			<Form.Field form={sf} name="confirmPassword">
				<Form.Control>
					{#snippet children({ props })}
						<Form.Label>Confirm New Password</Form.Label>
						<div class="relative">
							<Input
								{...props}
								type={showConfirmPassword ? 'text' : 'password'}
								bind:value={$form.confirmPassword}
								placeholder="Confirm your new password"
								class={cn(
									$errors.confirmPassword &&
										$errors.confirmPassword.length > 0 &&
										'border-destructive focus-visible:ring-destructive',
								)}
							/>
							<Button
								variant="ghost"
								size="sm"
								type="button"
								class="absolute top-0 right-0 h-full px-3 py-2 hover:bg-transparent"
								onclick={() => (showConfirmPassword = !showConfirmPassword)}
							>
								{#if showConfirmPassword}
									<EyeOff class="h-4 w-4" />
								{:else}
									<Eye class="h-4 w-4" />
								{/if}
							</Button>
						</div>
					{/snippet}
				</Form.Control>
				<Form.Description>Re-enter your new password to confirm</Form.Description>
				<Form.FieldErrors />
			</Form.Field>

			<Button type="submit" disabled={$submitting} class="w-full">
				{#if $submitting}
					<LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Update Password
			</Button>
		</form>
	</Card.Content>
</Card.Root>
