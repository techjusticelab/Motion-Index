<script lang="ts">
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { cubicOut, quintOut, backOut, elasticOut } from 'svelte/easing';
	import { browser } from '$app/environment';
	import { supabase } from '../../lib/db/supabase';
	const { data } = await supabase.auth.getSession();
	const session = data.session;

	let email = '';
	let password = '';
	let loading = false;
	let error: string | null = null;

	// Import the auth store
	import { user, isLoading } from '../../lib/stores/auth';

	async function handleLogin() {
		try {
			loading = true;
			error = null;

			// Set global loading state
			isLoading.set(true);

			const { data, error: err } = await $page.data.supabase.auth.signInWithPassword({
				email,
				password
			});

			if (err) throw err;

			// Extract the user data and update the store
			if (data && data.user) {
				// Set the user in the store
				user.set(data.user);
				console.log('User set in store:', data.user);
			}

			// Log successful login
			console.log('Login successful:', data);

			// Make sure to wait for the invalidation to complete
			await invalidateAll();

			// Add a console log before redirect
			console.log('Redirecting to:', $page.url.searchParams.get('redirectTo') || '/');

			// Use a slight delay before redirecting
			setTimeout(() => {
				const redirectTo = $page.url.searchParams.get('redirectTo') || '/';
				goto(redirectTo);
			}, 100);
		} catch (err: any) {
			console.error('Login error:', err);
			error = err.message || 'Failed to sign in';
		} finally {
			loading = false;
			isLoading.set(false);
		}
	}
</script>

<div
	class="flex min-h-screen flex-col justify-center bg-gray-50 py-12 sm:px-6 lg:px-8"
	in:fade={{ duration: 600, easing: cubicOut }}
>
	<div class="sm:mx-auto sm:w-full sm:max-w-md" in:fly={{ y: 20, duration: 700, easing: cubicOut }}>
		<h2
			class="mt-6 text-center text-3xl font-extrabold text-gray-900"
			in:fly={{ y: -10, duration: 700, delay: 200, easing: cubicOut }}
		>
			Sign in to Motion Index
		</h2>
		<p class="mt-2 text-center text-sm text-gray-600" in:fade={{ duration: 600, delay: 300 }}>
			Access your legal documents repository
		</p>
	</div>

	<div
		class="mt-8 sm:mx-auto sm:w-full sm:max-w-md"
		in:fly={{ y: 30, duration: 700, delay: 300, easing: cubicOut }}
	>
		<div
			class="bg-white px-4 py-8 shadow sm:rounded-lg sm:px-10"
			in:scale={{ start: 0.97, duration: 600, delay: 400, easing: cubicOut }}
		>
			<form
				class="space-y-6"
				on:submit|preventDefault={handleLogin}
				in:fade={{ duration: 500, delay: 500 }}
			>
				{#if error}
					<div class="rounded-md bg-red-50 p-4" in:fly={{ y: -5, duration: 500, easing: cubicOut }}>
						<div class="flex">
							<div class="ml-3">
								<h3 class="text-sm font-medium text-red-800" in:slide={{ duration: 400 }}>
									{error}
								</h3>
							</div>
						</div>
					</div>
				{/if}

				<div in:fly={{ y: 15, duration: 500, delay: 600, easing: cubicOut }}>
					<label for="email" class="block text-sm font-medium text-gray-700"> Email address </label>
					<div class="mt-1">
						<input
							id="email"
							name="email"
							type="email"
							autocomplete="email"
							required
							bind:value={email}
							class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
						/>
					</div>
				</div>

				<div in:fly={{ y: 15, duration: 500, delay: 700, easing: cubicOut }}>
					<label for="password" class="block text-sm font-medium text-gray-700"> Password </label>
					<div class="mt-1">
						<input
							id="password"
							name="password"
							type="password"
							autocomplete="current-password"
							required
							bind:value={password}
							class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
						/>
					</div>
				</div>

				<div
					class="flex items-center justify-between"
					in:fly={{ y: 15, duration: 500, delay: 800, easing: cubicOut }}
				>
					<div class="text-sm">
						<a
							href="/auth/forgot-password"
							class="font-medium text-indigo-600 hover:text-indigo-500"
						>
							Forgot your password?
						</a>
					</div>
				</div>

				<div in:fly={{ y: 15, duration: 500, delay: 900, easing: cubicOut }}>
					<button
						type="submit"
						disabled={loading}
						class="flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						in:scale={{ start: 0.98, duration: 600, delay: 1000, easing: backOut }}
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

			<div class="mt-6" in:fly={{ y: 20, duration: 600, delay: 1100, easing: cubicOut }}>
				<div class="relative">
					<div class="absolute inset-0 flex items-center">
						<div
							class="w-full border-t border-gray-300"
							in:scale={{ start: 0.8, duration: 500, delay: 1200, easing: cubicOut }}
						></div>
					</div>
					<div class="relative flex justify-center text-sm">
						<span class="bg-white px-2 text-gray-500" in:fade={{ duration: 400, delay: 1300 }}>
							Don't have an account?
						</span>
					</div>
				</div>

				<div class="mt-6" in:fly={{ y: 15, duration: 500, delay: 1400, easing: cubicOut }}>
					<a
						href="/auth/register"
						class="flex w-full justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
						in:scale={{ start: 0.98, duration: 600, delay: 1500, easing: backOut }}
					>
						Register
					</a>
				</div>
			</div>
		</div>
	</div>
</div>
