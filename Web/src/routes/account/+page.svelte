<script lang="ts">
	import { page } from '$app/stores';
	import { invalidateAll } from '$app/navigation';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let user = $page.data.session?.user;
	let profile = null;
	let loading = false;
	let updateStatus: { success?: boolean; message?: string } = {};
	let cpdaId = '';

	onMount(async () => {
		if (user) {
			loading = true;
			const { data, error } = await $page.data.supabase
				.from('Users')
				.select('*')
				.eq('id', user.id)
				.single();

			if (data && !error) {
				profile = data;
				cpdaId = profile.cpda_id || '';
			}
			loading = false;
		}
	});

	async function updateProfile() {
		try {
			loading = true;
			updateStatus = {};

			const { error } = await $page.data.supabase
				.from('profiles')
				.update({
					cpda_id: cpdaId || null,
					updated_at: new Date().toISOString()
				})
				.eq('id', user.id);

			if (error) throw error;

			updateStatus = { success: true, message: 'Profile updated successfully' };
		} catch (err: any) {
			console.error('Update profile error:', err);
			updateStatus = { success: false, message: err.message || 'Failed to update profile' };
		} finally {
			loading = false;
		}
	}

	async function handleSignOut() {
		const { error } = await $page.data.supabase.auth.signOut();
		if (!error) {
			await invalidateAll();
			goto('/');
		}
	}
</script>

<div class="mx-auto max-w-3xl px-4 py-10 sm:px-6 lg:px-8">
	<header>
		<h1 class="text-3xl font-bold text-gray-900">Account Settings</h1>
	</header>

	<div class="mt-10 overflow-hidden bg-white shadow sm:rounded-lg">
		<div class="px-4 py-5 sm:px-6">
			<h2 class="text-lg font-medium leading-6 text-gray-900">Your Profile</h2>
			<p class="mt-1 max-w-2xl text-sm text-gray-500">Manage your personal information</p>
		</div>

		{#if loading}
			<div class="px-4 py-10 text-center sm:px-6">
				<div
					class="mx-auto h-8 w-8 animate-spin rounded-full border-4 border-indigo-500 border-t-transparent"
				></div>
				<p class="mt-2 text-sm text-gray-500">Loading your information...</p>
			</div>
		{:else if profile}
			<div class="border-t border-gray-200 px-4 py-5 sm:p-6">
				<form on:submit|preventDefault={updateProfile} class="space-y-6">
					{#if updateStatus.message}
						<div class={`rounded-md p-4 ${updateStatus.success ? 'bg-green-50' : 'bg-red-50'}`}>
							<div class="flex">
								<div class="ml-3">
									<h3
										class={`text-sm font-medium ${updateStatus.success ? 'text-green-800' : 'text-red-800'}`}
									>
										{updateStatus.message}
									</h3>
								</div>
							</div>
						</div>
					{/if}

					<div>
						<label for="email" class="block text-sm font-medium text-gray-700"> Email </label>
						<div class="mt-1">
							<input
								type="email"
								id="email"
								value={user.email}
								disabled
								class="block w-full rounded-md border-gray-300 bg-gray-100 px-3 py-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
							/>
							<p class="mt-1 text-xs text-gray-500">Email cannot be changed</p>
						</div>
					</div>

					<div>
						<label for="cpda-id" class="block text-sm font-medium text-gray-700"> CPDA ID </label>
						<div class="mt-1">
							<input
								type="text"
								id="cpda-id"
								bind:value={cpdaId}
								class="block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
							/>
						</div>
					</div>

					<div class="flex justify-between">
						<button
							type="submit"
							class="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
							disabled={loading}
						>
							{loading ? 'Saving...' : 'Save Changes'}
						</button>

						<button
							type="button"
							on:click={handleSignOut}
							class="inline-flex justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
						>
							Sign Out
						</button>
					</div>
				</form>
			</div>
		{:else}
			<div class="px-4 py-5 text-center text-red-600 sm:px-6">
				Failed to load profile information. Please try refreshing the page.
			</div>
		{/if}
	</div>

	<div class="mt-6 overflow-hidden bg-white shadow sm:rounded-lg">
		<div class="flex items-center justify-between px-4 py-5 sm:px-6">
			<div>
				<h2 class="text-lg font-medium leading-6 text-gray-900">Security</h2>
				<p class="mt-1 max-w-2xl text-sm text-gray-500">Manage your password</p>
			</div>
			<a
				href="/auth/reset-password"
				class="text-sm font-medium text-indigo-600 hover:text-indigo-900"
			>
				Change Password
			</a>
		</div>
	</div>
</div>
