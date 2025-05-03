<script lang="ts">
	import { page } from '$app/stores';
	import { user } from './lib/stores/auth';
	import { onMount, onDestroy } from 'svelte';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { invalidate } from '$app/navigation';
	import { cubicOut, quintOut, backOut, elasticOut } from 'svelte/easing';
	import '../app.css';

	let { data, children } = $props();
	let { session, supabase } = $derived(data);

	// Flag to control animations after initial page load
	let isInitialLoad = true;

	// Auth token is now directly accessed from localStorage in the API client

	// Handle auth state changes
	let unsubscribe: () => void;

	$effect(() => {
		if (session) {
			const { data } = supabase.auth.onAuthStateChange((_, newSession) => {
				if (session?.access_token !== newSession?.access_token) {
					// Just invalidate the session when it changes
					invalidate('supabase:auth');
				}
			});
			unsubscribe = () => data.subscription.unsubscribe();
		}
	})

	onMount(() => {
		console.log('User session:', $page.data.user);
		// Set initial load to false after the first render
		setTimeout(() => {
			isInitialLoad = false;
		}, 100);
	});

	onDestroy(() => {
		if (unsubscribe) {
			unsubscribe();
		}
	});
</script>

<div class="min-h-screen bg-gray-50" in:fade={{ duration: 300, easing: cubicOut }}>
	<header class="bg-indigo-600 shadow-md" in:fly={{ y: -20, duration: 700, easing: cubicOut }}>
		<div class="container mx-auto px-4 py-4" in:fade={{ duration: 500, delay: 100 }}>
			<div class="flex flex-col items-center justify-between sm:flex-row">
				<a
					href="/"
					class="text-2xl font-bold text-white transition hover:text-indigo-200"
					in:scale={{ start: 0.9, duration: 600, delay: 200, easing: elasticOut }}
					>Motion Index <p class="text-gray-6000 text-sm font-normal">
						Berkeley Technology and Justice Lab
					</p>
				</a>
				<nav class="mt-3 sm:mt-0" in:fly={{ y: -10, duration: 500, delay: 300, easing: cubicOut }}>
					<ul class="flex space-x-6 text-white">
						{#each [{ href: '/', text: 'Search', isSpecial: false }, { href: '/upload', text: 'Upload', isSpecial: false }, { href: '/account', text: 'Account', isSpecial: false }, { href: '/help', text: 'Help', isSpecial: true }] as item, i}
							<li
								in:fly={{
									x: 10,
									y: -5,
									duration: 500,
									delay: 400 + i * 100,
									easing: cubicOut
								}}
							>
								{#if item.isSpecial}
									<a
										href={item.href}
										class="rounded bg-white p-2 font-bold text-indigo-600 transition hover:bg-indigo-800 hover:text-white"
										in:scale={{ start: 0.95, duration: 600, delay: 400 + i * 100, easing: backOut }}
										>{item.text}</a
									>
								{:else}
									<a href={item.href} class="transition hover:text-indigo-200">{item.text}</a>
								{/if}
							</li>
						{/each}
					</ul>
				</nav>
			</div>
		</div>
	</header>

	<main in:fade={{ duration: 600, delay: 400 }}>
		{@render children()}
	</main>

	<footer
		class="mt-12 border-t bg-gray-100"
		in:fly={{ y: 20, duration: 600, delay: 500, easing: cubicOut }}
	>
		<div
			class="container mx-auto px-4 py-6 text-left text-sm text-gray-600"
			in:fade={{ duration: 500, delay: 600 }}
		>
			<p in:slide={{ duration: 500, delay: 700 }}>
				&copy; 2025 Berkeley Technology and Justice Lab. All rights reserved.
			</p>
		</div>
	</footer>
</div>
