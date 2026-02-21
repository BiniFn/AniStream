<script lang="ts">
	import { type } from 'arktype';
	import { ArrowRight, LoaderCircle, Mail, Sparkles } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import { defaults, setError, superForm } from 'sveltekit-superforms';
	import { arktype, arktypeClient } from 'sveltekit-superforms/adapters';
	import { apiClient } from '$lib/api/client';
	import type { components } from '$lib/api/openapi';
	import * as Form from '$lib/components/ui/form';
	import Input from '$lib/components/ui/input/input.svelte';
	import { cn } from '$lib/utils';

	const forgotPasswordSchema = type({
		email: type('string.email').describe('a valid email address'),
	});

	const sf = superForm(defaults(arktype(forgotPasswordSchema)), {
		SPA: true,
		validators: arktypeClient(forgotPasswordSchema),
		onUpdate: async ({ form, cancel }) => {
			if (!form.valid) return;
			try {
				const res = await apiClient.POST('/auth/forget-password', {
					body: {
						email: form.data.email,
					},
				});

				if (res.response.status === 400) {
					const error = res.error as components['schemas']['models.ValidationErrorResponse'];
					if (error?.details) {
						Object.entries(error.details).forEach(([field, messages]) => {
							setError(form, field as 'email', messages);
						});
						toast.error('Please fix the errors in the form and try again.');
						cancel();
						return;
					}
				}

				toast.success(
					'If an account with this email exists, a password reset link has been sent. Check your inbox.',
				);
			} catch (error) {
				console.error('Forgot password error:', error);
				toast.error('Failed to send reset email. Please try again.');
				cancel();
			}
		},
	});

	const { form, enhance, submitting, errors } = sf;
</script>

<svelte:head>
	<title>Forgot Password - AniStream</title>
	<meta
		name="description"
		content="Reset your AniStream password. Enter your email address to receive a password reset link."
	/>
</svelte:head>

<div
	class="relative min-h-screen overflow-hidden bg-gradient-to-br from-background via-background to-primary/10"
>
	<div class="pointer-events-none absolute inset-0 overflow-hidden">
		<div
			class="absolute -top-20 -right-20 h-64 w-64 animate-pulse rounded-full bg-primary/5 blur-3xl"
		></div>
		<div
			class="absolute top-1/2 -left-20 h-48 w-48 animate-pulse rounded-full bg-secondary/10 blur-2xl"
			style="animation-delay: 1s;"
		></div>
	</div>

	<div class="relative z-10 container mx-auto px-4 py-8">
		<div class="mx-auto w-full max-w-md">
			<div class="mb-8 text-center">
				<div class="mb-4 inline-flex items-center gap-2 rounded-full bg-primary/10 px-4 py-2">
					<Sparkles class="h-4 w-4 text-primary" />
					<span class="text-sm font-semibold tracking-wider text-primary uppercase">
						Password Reset
					</span>
				</div>
				<h1 class="mb-2 text-3xl font-bold tracking-tight">Forgot Password?</h1>
				<p class="text-muted-foreground">
					No worries! Enter your email address and we'll send you a reset link.
				</p>
			</div>

			<div class="rounded-2xl border border-primary/10 bg-card/80 p-8 shadow-2xl backdrop-blur-sm">
				<form method="POST" use:enhance class="space-y-6">
					<Form.Field form={sf} name="email">
						<Form.Control>
							{#snippet children({ props })}
								<Form.Label>Email Address</Form.Label>
								<div class="relative">
									<Mail
										class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground"
									/>
									<Input
										{...props}
										bind:value={$form.email}
										placeholder="your@email.com"
										class={cn(
											'pl-10',
											$errors.email &&
												$errors.email.length > 0 &&
												'border-destructive focus-visible:ring-destructive',
										)}
									/>
								</div>
							{/snippet}
						</Form.Control>
						<Form.Description>
							We'll send a password reset link to this email address
						</Form.Description>
						<Form.FieldErrors />
					</Form.Field>

					<Form.Button size="lg" class="group w-full gap-2" disabled={$submitting}>
						{#if $submitting}
							<LoaderCircle class="size-4 animate-spin" />
							Sending Reset Link...
						{:else}
							Send Reset Link
							<ArrowRight class="size-4 transition group-hover:translate-x-1" />
						{/if}
					</Form.Button>
				</form>

				<div class="my-8 flex items-center">
					<div class="flex-1 border-t border-muted"></div>
					<div class="px-4 text-sm text-muted-foreground">or</div>
					<div class="flex-1 border-t border-muted"></div>
				</div>

				<div class="text-center">
					<p class="text-sm text-muted-foreground">
						Remember your password?
						<a href="/login" class="font-medium text-primary transition-colors hover:underline">
							Sign In
						</a>
					</p>
				</div>
			</div>
		</div>
	</div>
</div>
