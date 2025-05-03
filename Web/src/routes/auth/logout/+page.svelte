<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { clearAuth } from '../../../lib/auth';
	
	// Handle logout on the client side
	onMount(async () => {
		if (browser) {
			try {
				// Get the Supabase client from the page data
				const { supabase } = $page.data;
				
				// Sign out from Supabase
				await supabase.auth.signOut();
				
				// Clear our custom auth state
				clearAuth();
				
				console.log('User successfully logged out');
				
				// Redirect to login page
				await goto('/auth/login?message=You have been signed out');
			} catch (error) {
				console.error('Error signing out:', error);
				
				// Still clear auth even if Supabase logout fails
				clearAuth();
				
				// Redirect to login
				await goto('/auth/login?error=Error signing out');
			}
		}
	});
</script>

<div class="flex h-screen w-full items-center justify-center bg-gray-50">
	<div class="w-full max-w-md rounded-lg bg-white p-8 shadow-lg">
		<div class="text-center">
			<h2 class="mb-4 text-2xl font-bold text-gray-800">Signing Out...</h2>
			<p class="text-gray-600">Please wait while we sign you out.</p>
			<div class="mt-6 flex justify-center">
				<div class="h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-blue-600"></div>
			</div>
		</div>
	</div>
</div>
