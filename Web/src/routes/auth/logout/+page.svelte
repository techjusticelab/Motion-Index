<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { invalidateAll } from '$app/navigation';

	// Handle logout on the client side
	onMount(async () => {
		if (browser) {
			try {
				// Get the Supabase client from the page data
				const { supabase } = $page.data;

				// Sign out from Supabase
				await supabase.auth.signOut();

				// Invalidate all cached data to refresh auth state
				await invalidateAll();

				console.log('User successfully logged out');

				// Redirect to login page
				await goto('/auth/login?message=You have been signed out');
			} catch (error) {
				console.error('Error signing out:', error);

				// Still try to invalidate and redirect
				await invalidateAll();
				await goto('/auth/login?error=Error signing out');
			}
		}
	});
</script>

<div class="flex h-screen w-full items-center justify-center bg-neutral-50">
	<div class="w-full max-w-md rounded-lg bg-white p-8 shadow-lg">
		<div class="text-center">
			<h2 class="mb-4 text-2xl font-bold text-neutral-800">Signing Out...</h2>
			<p class="text-neutral-600">Please wait while we sign you out.</p>
			<div class="mt-6 flex justify-center">
				<div class="h-8 w-8 animate-spin rounded-full border-4 border-neutral-300 border-t-blue-600"></div>
			</div>
		</div>
	</div>
</div>
