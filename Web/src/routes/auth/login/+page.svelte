<script lang="ts">
	import { enhance } from '$app/forms';

	import { invalidateAll } from '$app/navigation';

	import { page } from '$app/stores';

	import { goto } from '$app/navigation';

	import { fade, fly, scale } from 'svelte/transition';

	import { cubicOut } from 'svelte/easing';

	let email = '';

	let password = '';

	let loading = false;

	let error: string | null = null;
	let successMessage: string | null = null;
	
	// Check for success message in URL parameters
	$: {
		const message = $page.url.searchParams.get('message');
		if (message) {
			successMessage = message;
		}
	}

	const handleSubmit = async (event: SubmitEvent) => {
		event.preventDefault();

		try {
			loading = true;
			error = null;

			// Use the Supabase client from page data
			const { data, error: err } = await $page.data.supabase.auth.signInWithPassword({
				email,
				password
			});

			if (err) throw err;

			if (data?.session) {
				console.log('Login successful');

				// Invalidate all cached data to refresh auth state
				await invalidateAll();

				// Redirect to the requested page or home
				const redirectTo = $page.url.searchParams.get('redirectTo') || '/';
				goto(redirectTo);
			} else {
				throw new Error('Login successful but session data is missing');
			}
		} catch (err: any) {
			console.error('Login error:', err);
			error = err.message || 'Failed to sign in';
		} finally {
			loading = false;
		}
	};
</script>

<div class="flex min-h-screen flex-col justify-center bg-neutral-50 py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<h2 class="mt-6 text-center text-3xl font-extrabold text-neutral-900">Sign in to Motion Index</h2>
		<p class="mt-2 text-center text-sm text-neutral-600">Access your legal documents repository</p>
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

				{#if successMessage}
					<div class="rounded-md bg-green-50 p-4">
						<div class="flex">
							<div class="ml-3">
								<h3 class="text-sm font-medium text-green-800">
									{successMessage}
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
							autocomplete="current-password"
							required
							bind:value={password}
							class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 placeholder-neutral-400 shadow-sm focus:border-primary-900 focus:outline-none focus:ring-primary-900 sm:text-sm"
						/>
					</div>
				</div>

				<div class="flex items-center justify-between">
					<div class="text-sm">
						<a href="/auth/forgot" class="font-medium text-primary-900 hover:text-primary-800">
							Forgot your password?
						</a>
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
							Signing in...
						{:else}
							Sign in
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
						<span class="bg-white px-2 text-neutral-500">Don't have an account?</span>
					</div>
				</div>

				<div class="mt-6">
					<a
						href="/auth/register"
						class="flex w-full justify-center rounded-md border border-neutral-300 bg-white px-4 py-2 text-sm font-medium text-neutral-700 shadow-sm hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-900 focus:ring-offset-2"
					>
						Register
					</a>
				</div>
			</div>
		</div>
	</div>
</div>
