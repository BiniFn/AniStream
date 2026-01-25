<script lang="ts">
	import { type } from 'arktype';
	import {
		ArrowLeft,
		ArrowRight,
		Check,
		Eye,
		EyeOff,
		LoaderCircle,
		Lock,
		Sparkles,
	} from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import { defaults, setError, superForm } from 'sveltekit-superforms';
	import { arktype, arktypeClient } from 'sveltekit-superforms/adapters';
	import { goto } from '$app/navigation';
	import { apiClient } from '$lib/api/client';
	import type { components } from '$lib/api/openapi';
	import * as Form from '$lib/components/ui/form';
	import Input from '$lib/components/ui/input/input.svelte';
	import { cn } from '$lib/utils';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const resetPasswordSchema = type({
		password: type('string>=8&/^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).*$/').describe(
			'a strong password with uppercase, lowercase, and number',
		),
		confirmPassword: 'string',
	}).narrow((data, ctx) => {
		if (data.password !== data.confirmPassword) {
			return ctx.mustBe('Passwords do not match');
		}
		return true;
	});

	const sf = superForm(defaults(arktype(resetPasswordSchema)), {
		SPA: true,
		validators: arktypeClient(resetPasswordSchema),
		onUpdate: async ({ form, cancel }) => {
			if (!form.valid) return;
			try {
				const res = await apiClient.PUT('/auth/reset-password/{token}', {
					params: {
						path: { token: data.token },
					},
					body: {
						password: form.data.password,
					},
				});

				if (res.response.status === 200) {
					await goto('/login');
					toast.success('Password reset successfully! You can now sign in with your new password.');
					return;
				}

				if (res.response.status === 400) {
					const error = res.error as components['schemas']['models.ValidationErrorResponse'];
					if (error?.details) {
						Object.entries(error.details).forEach(([field, messages]) => {
							setError(form, field as 'password', messages);
						});
						toast.error('Please fix the errors in the form and try again.');
						cancel();
						return;
					}
				}

				if (res.response.status === 401) {
					toast.error('Invalid or expired reset token. Please request a new password reset.');
					cancel();
					return;
				}

				throw new Error('Unexpected response from server');
			} catch (error) {
				console.error('Reset password error:', error);
				toast.error('Failed to reset password. Please try again.');
				cancel();
			}
		},
	});

	const { form, enhance, submitting, errors } = sf;
	let showPassword = $state(false);
	let showConfirmPassword = $state(false);

	let passwordValidation = $derived.by(() => ({
		length: $form.password.length >= 8,
		hasUppercase: /[A-Z]/.test($form.password),
		hasLowercase: /[a-z]/.test($form.password),
		hasNumber: /\d/.test($form.password),
	}));

	let passwordMatch = $derived(
		$form.password === $form.confirmPassword && $form.confirmPassword.length > 0,
	);

	let expirationTimer: NodeJS.Timeout | null = null;

	$effect(() => {
		if (data.expiresAt && data.isValidToken) {
			const now = Date.now();
			const remaining = data.expiresAt - now;

			if (remaining > 0) {
				expirationTimer = setTimeout(() => {
					window.location.reload();
				}, remaining);
			}
		}

		return () => {
			if (expirationTimer) {
				clearTimeout(expirationTimer);
			}
		};
	});
</script>

<svelte:head>
	<title>Reset Password - Aniways</title>
	<meta
		name="description"
		content="Reset your Aniways password using the secure reset token sent to your email."
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
				<h1 class="mb-2 text-3xl font-bold tracking-tight">Reset Your Password</h1>
				{#if data.isValidToken && data.user}
					<p class="text-muted-foreground">
						Hello <span class="font-medium">{data.user.username}</span>! Enter your new password
						below.
					</p>
				{:else}
					<p class="text-muted-foreground">
						This reset link is invalid or has expired. Please request a new password reset.
					</p>
				{/if}
			</div>

			{#if data.isValidToken && data.user}
				<div
					class="rounded-2xl border border-primary/10 bg-card/80 p-8 shadow-2xl backdrop-blur-sm"
				>
					<form method="POST" use:enhance class="space-y-6">
						<Form.Field form={sf} name="password">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>New Password</Form.Label>
									<div class="relative">
										<Lock
											class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground"
										/>
										<Input
											{...props}
											type={showPassword ? 'text' : 'password'}
											bind:value={$form.password}
											placeholder="Create a strong password"
											class={cn(
												'pr-10 pl-10',
												$errors.password &&
													$errors.password.length > 0 &&
													'border-destructive focus-visible:ring-destructive',
											)}
											disabled={$submitting}
										/>
										<button
											type="button"
											onclick={() => (showPassword = !showPassword)}
											class="absolute top-1/2 right-3 -translate-y-1/2 text-muted-foreground transition-colors hover:text-foreground"
											disabled={$submitting}
										>
											{#if showPassword}
												<EyeOff class="h-4 w-4" />
											{:else}
												<Eye class="h-4 w-4" />
											{/if}
										</button>
									</div>

									{#if $form.password.length > 0}
										<div class="mt-3 space-y-2 rounded-lg bg-muted/30 p-4 backdrop-blur-sm">
											<p class="text-sm font-medium text-muted-foreground">
												Password Requirements:
											</p>
											<div class="grid grid-cols-2 gap-2 text-sm">
												<div class="flex items-center gap-2">
													<Check
														class="h-4 w-4 {passwordValidation.length
															? 'text-green-500'
															: 'text-muted-foreground'}"
													/>
													<span
														class={passwordValidation.length
															? 'text-green-600'
															: 'text-muted-foreground'}
													>
														8+ characters
													</span>
												</div>
												<div class="flex items-center gap-2">
													<Check
														class="h-4 w-4 {passwordValidation.hasUppercase
															? 'text-green-500'
															: 'text-muted-foreground'}"
													/>
													<span
														class={passwordValidation.hasUppercase
															? 'text-green-600'
															: 'text-muted-foreground'}
													>
														Uppercase
													</span>
												</div>
												<div class="flex items-center gap-2">
													<Check
														class="h-4 w-4 {passwordValidation.hasLowercase
															? 'text-green-500'
															: 'text-muted-foreground'}"
													/>
													<span
														class={passwordValidation.hasLowercase
															? 'text-green-600'
															: 'text-muted-foreground'}
													>
														Lowercase
													</span>
												</div>
												<div class="flex items-center gap-2">
													<Check
														class="h-4 w-4 {passwordValidation.hasNumber
															? 'text-green-500'
															: 'text-muted-foreground'}"
													/>
													<span
														class={passwordValidation.hasNumber
															? 'text-green-600'
															: 'text-muted-foreground'}
													>
														Number
													</span>
												</div>
											</div>
										</div>
									{/if}
								{/snippet}
							</Form.Control>
							<Form.FieldErrors />
						</Form.Field>

						<Form.Field form={sf} name="confirmPassword">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Confirm New Password</Form.Label>
									<div class="relative">
										<Lock
											class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground"
										/>
										<Input
											{...props}
											type={showConfirmPassword ? 'text' : 'password'}
											bind:value={$form.confirmPassword}
											placeholder="Confirm your new password"
											class={cn(
												'pr-10 pl-10',
												$errors.confirmPassword &&
													$errors.confirmPassword.length > 0 &&
													'border-destructive focus-visible:ring-destructive',
											)}
											disabled={$submitting}
										/>
										<button
											type="button"
											onclick={() => (showConfirmPassword = !showConfirmPassword)}
											class="absolute top-1/2 right-3 -translate-y-1/2 text-muted-foreground transition-colors hover:text-foreground"
											disabled={$submitting}
										>
											{#if showConfirmPassword}
												<EyeOff class="h-4 w-4" />
											{:else}
												<Eye class="h-4 w-4" />
											{/if}
										</button>
									</div>

									{#if $form.confirmPassword.length > 0}
										<div class="mt-2 flex items-center gap-2 text-sm">
											<Check
												class="h-4 w-4 {passwordMatch ? 'text-green-500' : 'text-muted-foreground'}"
											/>
											<span class={passwordMatch ? 'text-green-600' : 'text-muted-foreground'}>
												{#if passwordMatch}
													Passwords match
												{:else}
													Passwords do not match
												{/if}
											</span>
										</div>
									{/if}
								{/snippet}
							</Form.Control>
							<Form.FieldErrors />
						</Form.Field>

						<Form.Button size="lg" class="group w-full gap-2" disabled={$submitting}>
							{#if $submitting}
								<LoaderCircle class="size-4 animate-spin" />
								Resetting Password...
							{:else}
								Reset Password
								<ArrowRight class="size-4 transition group-hover:translate-x-1" />
							{/if}
						</Form.Button>
					</form>
				</div>
			{:else}
				<div
					class="rounded-2xl border border-destructive/20 bg-destructive/5 p-8 shadow-2xl backdrop-blur-sm"
				>
					<div class="text-center">
						<div
							class="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10"
						>
							<Lock class="h-6 w-6 text-destructive" />
						</div>
						<h2 class="mb-2 text-xl font-semibold text-destructive">Invalid Reset Link</h2>
						<p class="mb-6 text-sm text-muted-foreground">
							This password reset link is invalid or has expired. Please request a new password
							reset.
						</p>
						<a
							href="/forgot-password"
							class="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
						>
							Request New Reset Link
						</a>
					</div>
				</div>
			{/if}

			<div class="mt-6 text-center">
				<a
					href="/login"
					class="inline-flex items-center gap-2 text-sm text-muted-foreground transition-colors hover:text-primary"
				>
					<ArrowLeft class="h-4 w-4" />
					Back to Sign In
				</a>
			</div>
		</div>
	</div>
</div>
