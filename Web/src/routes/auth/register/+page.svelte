<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { fade, fly, scale } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';

	let email = '';
	let password = '';
	let confirmPassword = '';
	let loading = false;
	let error: string | null = null;

	const handleSubmit = async (event: SubmitEvent) => {
		event.preventDefault();

		// Basic validation
		if (password !== confirmPassword) {
			error = "Passwords don't match";
			return;
		}

		if (password.length < 6) {
			error = "Password must be at least 6 characters";
			return;
		}

		try {
			loading = true;
			error = null;

			// Sign up with Supabase
			const { data, error: err } = await $page.data.supabase.auth.signUp({
				email,
				password,
				options: {
					emailRedirectTo: `${window.location.origin}/auth/confirm`
				}
			});

			if (err) throw err;

			// Add a console log to confirm the registration was successful
			console.log('Registration successful:', data);

			// Make sure to wait for the invalidation to complete
			await invalidateAll();

			// Redirect to confirmation page
			goto('/auth/login?message=Check your email to confirm your account');
		} catch (err: any) {
			console.error('Registration error:', err);
			error = err.message || 'Failed to register';
		} finally {
			loading = false;
		}
	};
</script>

<div class="flex min-h-screen flex-col justify-center bg-neutral-50 py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<h2 class="mt-6 text-center text-3xl font-extrabold text-neutral-900">Create an account</h2>
		<p class="mt-2 text-center text-sm text-neutral-600">Join Motion Index to access your legal documents</p>
	</div>

	<div
		class="mt-8 sm:mx-auto sm:w-full sm:max-w-md"
		in:fly={{ y: 30, duration: 700, delay: 300, easing: cubicOut }}
	>
		<div
			class="bg-white px-4 py-8 shadow sm:rounded-lg sm:px-10"
			in:scale={{ start: 0.97, duration: 600, delay: 400, easing: cubicOut }}
		>
			<form class="space-y-6" on:submit={handleSubmit} in:fade={{ duration: 500, delay: 500 }}>
				{#if error}
					<div class="rounded-md bg-red-50 p-4">
						<div class="flex">
							<div class="ml-3">
								<h3 class="text-sm font-medium text-red-800">
									{error}
								</h3>
							</div>
						</div>
					</div>
				{/if}

				<div>
					<label for="email" class="block text-sm font-medium text-neutral-700">Email address</label>
					<div class="mt-1">
						<input
							id="email"
							name="email"
							type="email"
							autocomplete="email"
							required
							bind:value={email}
							class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 placeholder-neutral-400 shadow-sm focus:border-primary-900 focus:outline-none focus:ring-primary-900 sm:text-sm"
						/>
					</div>
				</div>

				<div>
					<label for="password" class="block text-sm font-medium text-neutral-700">Password</label>
					<div class="mt-1">
						<input
							id="password"
							name="password"
							type="password"
							autocomplete="new-password"
							required
							bind:value={password}
							class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 placeholder-neutral-400 shadow-sm focus:border-primary-900 focus:outline-none focus:ring-primary-900 sm:text-sm"
						/>
					</div>
				</div>

				<div>
					<label for="confirm-password" class="block text-sm font-medium text-neutral-700"
						>Confirm Password</label
					>
					<div class="mt-1">
						<input
							id="confirm-password"
							name="confirm-password"
							type="password"
							autocomplete="new-password"
							required
							bind:value={confirmPassword}
							class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 placeholder-neutral-400 shadow-sm focus:border-primary-900 focus:outline-none focus:ring-primary-900 sm:text-sm"
						/>
					</div>
				</div>

				<div>
					<button
						type="submit"
						disabled={loading}
						class="flex w-full justify-center rounded-md border border-transparent bg-primary-900 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-primary-800 focus:outline-none focus:ring-2 focus:ring-primary-900 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
					>
						{#if loading}
							<div
								class="mr-2 h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent"
							></div>
							Creating account...
						{:else}
							Register
						{/if}
					</button>
				</div>
			</form>

			<div class="mt-6">
				<div class="relative">
					<div class="absolute inset-0 flex items-center">
						<div class="w-full border-t border-neutral-300"></div>
					</div>
					<div class="relative flex justify-center text-sm">
						<span class="bg-white px-2 text-neutral-500">Already have an account?</span>
					</div>
				</div>

				<div class="mt-6">
					<a
						href="/auth/login"
						class="flex w-full justify-center rounded-md border border-neutral-300 bg-white px-4 py-2 text-sm font-medium text-neutral-700 shadow-sm hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-900 focus:ring-offset-2"
					>
						Sign in
					</a>
				</div>
			</div>
		</div>
	</div>
</div>
