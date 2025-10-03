<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	let password = '';
	let confirmPassword = '';
	let loading = false;
	let error: string | null = null;
	let success = false;

	async function handleUpdatePassword() {
		try {
			loading = true;
			error = null;

			if (password !== confirmPassword) {
				error = 'Passwords do not match';
				return;
			}

			const { error: updateError } = await $page.data.supabase.auth.updateUser({
				password
			});

			if (updateError) throw updateError;

			success = true;
			setTimeout(() => goto('/auth/login'), 3000);
		} catch (err: any) {
			console.error('Update password error:', err);
			error = err.message || 'Failed to update password';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen flex-col justify-center bg-neutral-50 py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<h2 class="mt-6 text-center text-3xl font-extrabold text-neutral-900">Set new password</h2>
		<p class="mt-2 text-center text-sm text-neutral-600">Choose a new password for your account</p>
	</div>

	<div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
		<div class="bg-white px-4 py-8 shadow sm:rounded-lg sm:px-10">
			{#if success}
				<div class="rounded-md bg-green-50 p-4">
					<div class="flex">
						<div class="flex-shrink-0">
							<svg
								class="h-5 w-5 text-green-400"
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 20 20"
								fill="currentColor"
							>
								<path
									fill-rule="evenodd"
									d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
									clip-rule="evenodd"
								/>
							</svg>
						</div>
						<div class="ml-3">
							<h3 class="text-sm font-medium text-green-800">Password updated successfully</h3>
							<p class="mt-2 text-sm text-green-700">
								Your password has been changed. You will be redirected to the login page.
							</p>
						</div>
					</div>
				</div>
			{:else}
				<form class="space-y-6" on:submit|preventDefault={handleUpdatePassword}>
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
						<label for="password" class="block text-sm font-medium text-neutral-700">
							New Password
						</label>
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
						<label for="confirm-password" class="block text-sm font-medium text-neutral-700">
							Confirm New Password
						</label>
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
								Updating...
							{:else}
								Update password
							{/if}
						</button>
					</div>
				</form>
			{/if}
		</div>
	</div>
</div>
