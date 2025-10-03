<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	let email = '';
	let loading = false;
	let error: string | null = null;
	let success = false;

	async function handleResetPassword() {
		try {
			loading = true;
			error = null;

			const { error: resetError } = await $page.data.supabase.auth.resetPasswordForEmail(email, {
				redirectTo: `${window.location.origin}/auth/reset-password`
			});

			if (resetError) throw resetError;

			success = true;
		} catch (err: any) {
			console.error('Reset password error:', err);
			error = err.message || 'Failed to send reset password email';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen flex-col justify-center bg-neutral-50 py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<h2 class="mt-6 text-center text-3xl font-extrabold text-neutral-900">Reset your password</h2>
		<p class="mt-2 text-center text-sm text-neutral-600">
			Enter your email and we'll send you a password reset link
		</p>
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
							<h3 class="text-sm font-medium text-green-800">Reset link sent</h3>
							<p class="mt-2 text-sm text-green-700">
								Check your email for a link to reset your password. If it doesn't appear within a
								few minutes, check your spam folder.
							</p>
							<div class="mt-4">
								<a
									href="/auth/login"
									class="text-sm font-medium text-green-600 hover:text-green-500"
								>
									Return to login
								</a>
							</div>
						</div>
					</div>
				</div>
			{:else}
				<form class="space-y-6" on:submit|preventDefault={handleResetPassword}>
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
						<label for="email" class="block text-sm font-medium text-neutral-700">
							Email address
						</label>
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
						<button
							type="submit"
							disabled={loading}
							class="flex w-full justify-center rounded-md border border-transparent bg-primary-900 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-primary-800 focus:outline-none focus:ring-2 focus:ring-primary-900 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						>
							{#if loading}
								<div
									class="mr-2 h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent"
								></div>
								Sending...
							{:else}
								Send reset link
							{/if}
						</button>
					</div>
				</form>

				<div class="mt-6 text-center text-sm">
					<a href="/auth/login" class="font-medium text-primary-900 hover:text-primary-800">
						Back to login
					</a>
				</div>
			{/if}
		</div>
	</div>
</div>
