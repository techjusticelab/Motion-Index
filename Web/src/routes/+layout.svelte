<script lang="ts">
	import { page } from '$app/stores';
	import { onMount, onDestroy } from 'svelte';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { invalidate, goto } from '$app/navigation';
	import { cubicOut, quintOut, backOut, elasticOut } from 'svelte/easing';
	import { browser } from '$app/environment';
	import '../app.css';

	let { data, children } = $props();
	let { session, user, supabase } = $derived(data);

	// Flag to control animations after initial page load
	let isInitialLoad = true;

	let unsubscribe: (() => void) | null = null;

	onMount(() => {
		const {
			data: { subscription },
		} = supabase.auth.onAuthStateChange((event, newSession) => {
			if (newSession?.expires_at !== session?.expires_at) {
				invalidate('supabase:auth')
			}
		})
		unsubscribe = () => subscription.unsubscribe()

		// Set initial load to false after the first render
		setTimeout(() => {
			isInitialLoad = false;
		}, 100);

		return () => {
			if (unsubscribe) unsubscribe()
		}
	})
</script>

<div class="min-h-screen bg-neutral-50" in:fade={{ duration: 300, easing: cubicOut }}>
	<header class="bg-primary-900 shadow-md" in:fly={{ y: -20, duration: 700, easing: cubicOut }}>
		<div class="container mx-auto px-4 py-4" in:fade={{ duration: 500, delay: 100 }}>
			<div class="flex flex-col items-center justify-between sm:flex-row">
				<div
					class="text-2xl font-bold text-white"
					in:scale={{ start: 0.9, duration: 600, delay: 200, easing: elasticOut }}
				>
					{#if session}
						<a href="/" class="text-white transition hover:text-secondary-300">
							Motion Index
							<p class="text-neutral-200 text-sm font-normal">
								Berkeley Technology and Justice Lab
							</p>
						</a>
					{:else}
						<span>
							Motion Index
							<p class="text-neutral-200 text-sm font-normal">
								Berkeley Technology and Justice Lab
							</p>
						</span>
					{/if}
				</div>
				<nav class="mt-3 sm:mt-0" in:fly={{ y: -10, duration: 500, delay: 300, easing: cubicOut }}>
					{#if session}
						<ul class="flex space-x-6 text-white">
							<li><a href="/" class="transition hover:text-secondary-300">Search</a></li>
							<li><a href="/upload" class="transition hover:text-secondary-300">Upload</a></li>
							<li><a href="/cases" class="transition hover:text-secondary-300">Cases</a></li>
							<li><a href="/account" class="transition hover:text-secondary-300">Account</a></li>
							<li><a href="/help" class="rounded bg-secondary-700 p-2 font-bold text-neutral-50 transition hover:bg-secondary-600 hover:text-neutral-50">Help</a></li>
						</ul>
					{:else}
						<a href="/auth/login" class="rounded bg-secondary-700 p-2 font-bold text-neutral-50 transition hover:bg-secondary-600 hover:text-neutral-50">
							Sign In
						</a>
					{/if}
				</nav>
			</div>
		</div>
	</header>

	<main in:fade={{ duration: 600, delay: 400 }}>
		{@render children()}
	</main>

	<footer class="mt-12 border-t bg-neutral-100" in:fly={{ y: 20, duration: 600, delay: 500, easing: cubicOut }}>
		<div class="container mx-auto px-4 py-6 text-left text-sm text-neutral-600" in:fade={{ duration: 500, delay: 600 }}>
			<p in:slide={{ duration: 500, delay: 700 }}>
				&copy; 2025 Berkeley Technology and Justice Lab. All rights reserved.
			</p>
		</div>
	</footer>
</div>
