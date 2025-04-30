<script lang="ts">
	import { page } from '$app/stores';
	import { invalidateAll } from '$app/navigation';
	import { goto } from '$app/navigation';

	let email = '';
	let password = '';
	let confirmPassword = '';
	let cpdaId = '';
	let loading = false;
	let error: string | null = null;

	async function handleRegister() {
		try {
			loading = true;
			error = null;

			// Check if passwords match
			if (password !== confirmPassword) {
				error = 'Passwords do not match';
				return;
			}

			// Make sure supabase is available
			if (!$page.data.supabase) {
				throw new Error('Authentication service unavailable');
			}

			// Register user
			const { data, error: signUpError } = await $page.data.supabase.auth.signUp({
				email,
				password,
				options: {
					data: {
						cpda_id: cpdaId
					}
				}
			});

			if (signUpError) throw signUpError;

			if (data.user) {
				// Create profile in the profiles table
				const { error: profileError } = await $page.data.supabase.from('profiles').insert({
					id: data.user.id,
					email: data.user.email!,
					cpda_id: cpdaId || null
				});

				if (profileError) throw profileError;
			}

			await invalidateAll();

			// Redirect to login or confirmation page
			goto('/auth/confirm?email=' + encodeURIComponent(email));
		} catch (err: any) {
			console.error('Registration error:', err);
			error = err.message || 'Failed to register';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen flex-col justify-center bg-gray-50 py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">Create your account</h2>
		<p class="mt-2 text-center text-sm text-gray-600">
			Join Motion Index to manage your legal documents
		</p>
	</div>

	<div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
		<div class="bg-white px-4 py-8 shadow sm:rounded-lg sm:px-10">
			<form class="space-y-6" on:submit|preventDefault={handleRegister}>
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

				<div>
					<label for="password" class="block text-sm font-medium text-gray-700"> Password </label>
					<div class="mt-1">
						<input
							id="password"
							name="password"
							type="password"
							autocomplete="new-password"
							required
							bind:value={password}
							class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
						/>
					</div>
				</div>

				<div>
					<label for="confirm-password" class="block text-sm font-medium text-gray-700">
						Confirm Password
					</label>
					<div class="mt-1">
						<input
							id="confirm-password"
							name="confirm-password"
							type="password"
							autocomplete="new-password"
							required
							bind:value={confirmPassword}
							class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
						/>
					</div>
				</div>

				<div>
					<label for="cpda-id" class="block text-sm font-medium text-gray-700">
						CPDA ID (optional)
					</label>
					<div class="mt-1">
						<input
							id="cpda-id"
							name="cpda-id"
							type="text"
							bind:value={cpdaId}
							class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
						/>
					</div>
				</div>

				<div>
					<button
						type="submit"
						disabled={loading}
						class="flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
					>
						{#if loading}
							<div
								class="mr-2 h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent"
							></div>
							Creating account...
						{:else}
							Create account
						{/if}
					</button>
				</div>
			</form>

			<div class="mt-6">
				<div class="relative">
					<div class="absolute inset-0 flex items-center">
						<div class="w-full border-t border-gray-300"></div>
					</div>
					<div class="relative flex justify-center text-sm">
						<span class="bg-white px-2 text-gray-500"> Already have an account? </span>
					</div>
				</div>

				<div class="mt-6">
					<a
						href="/auth/login"
						class="flex w-full justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
					>
						Sign in
					</a>
				</div>
			</div>
		</div>
	</div>
</div>
