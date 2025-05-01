<script lang="ts">
	import { page } from '$app/stores';
	import { user, isLoading } from './lib/stores/auth';
	import { onMount } from 'svelte';
	import '../app.css';

	let { children } = $props();

	// Update the user store when session changes
	$effect(() => {
		if ($page.data) {
			user.set($page.data.session?.user || null);
			isLoading.set(false);
		}
	});

	// Handle authentication state
	$effect(() => {
		if ($page.data && !$page.data.session && !window.location.pathname.startsWith('/auth')) {
			window.location.href = `/auth/login?redirectTo=${encodeURIComponent(window.location.pathname)}`;
		}
	});
</script>

<div class="min-h-screen bg-gray-50">
	<header class="bg-indigo-600 shadow-md">
		<div class="container mx-auto px-4 py-4">
			<div class="flex flex-col items-center justify-between sm:flex-row">
				<a href="/" class="text-2xl font-bold text-white">Motion Index</a>
				<nav class="mt-3 sm:mt-0">
					<ul class="flex space-x-6 text-white">
						<li><a href="/" class="transition hover:text-indigo-200">Search</a></li>
						<li><a href="/upload" class="transition hover:text-indigo-200">Upload</a></li>

						<li>
							<a href="/account" class="transition hover:text-indigo-200"> Account</a>
						</li>
						<li>
							<a
								href="/help"
								class="rounded bg-white p-2 text-indigo-600 transition hover:bg-indigo-800 hover:text-white"
								>Help</a
							>
						</li>
					</ul>
				</nav>
			</div>
		</div>
	</header>

	<main>
		{@render children()}
	</main>

	<footer class="mt-12 border-t bg-gray-100">
		<div class="container mx-auto px-4 py-6 text-left text-sm text-gray-600">
			<p>&copy; 2025 Berkeley Technology and Justice Lab. All rights reserved.</p>
		</div>
	</footer>
</div>
