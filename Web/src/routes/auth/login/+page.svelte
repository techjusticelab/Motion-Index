<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/stores';
	import { fade, fly, scale } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import type { ActionData } from './$types';
	import { PUBLIC_SUPABASE_URL } from '$env/static/public';
	import { checkNetworkConnectivity, getNetworkErrorMessage, detectAdBlockers } from '$lib/utils/network-check';

	let { form }: { form: ActionData } = $props();
	
	let loading = $state(false);
	let successMessage: string | null = $state(null);
	let networkStatus: string | null = $state(null);
	let checkingNetwork = $state(false);
	
	// Check for success/error messages in URL parameters
	$effect(() => {
		const message = $page.url.searchParams.get('message');
		const error = $page.url.searchParams.get('error');
		
		if (message) {
			successMessage = decodeURIComponent(message);
		}
		
		// Clear URL parameters after showing message
		if (message || error) {
			const url = new URL($page.url);
			url.searchParams.delete('message');
			url.searchParams.delete('error');
			history.replaceState({}, '', url.toString());
		}
	});

	// Handle form submission with enhanced progressive enhancement
	function handleEnhance() {
		return ({ formElement, formData, action, cancel, submitter }) => {
			loading = true;
			successMessage = null;
			networkStatus = null;
			
			return async ({ result, update }) => {
				loading = false;
				
				if (result.type === 'redirect') {
					// Let SvelteKit handle the redirect
					await update();
				} else if (result.type === 'failure') {
					// Form validation or auth errors - update will show them
					await update();
				} else {
					// Success or other result types
					await update();
				}
			};
		};
	}

	// Network connectivity checker
	async function checkNetwork() {
		checkingNetwork = true;
		networkStatus = null;
		
		try {
			const result = await checkNetworkConnectivity(PUBLIC_SUPABASE_URL);
			const detectedBlockers = detectAdBlockers();
			
			let statusMessage = getNetworkErrorMessage(result);
			
			if (detectedBlockers.length > 0) {
				statusMessage += `\n\nDetected extensions: ${detectedBlockers.join(', ')}`;
			}
			
			if (result.supabaseReachable) {
				statusMessage = '✅ Network connectivity OK - Supabase is reachable';
			}
			
			networkStatus = statusMessage;
		} catch (error) {
			networkStatus = `Network check failed: ${error instanceof Error ? error.message : 'Unknown error'}`;
		} finally {
			checkingNetwork = false;
		}
	}
</script>

<div class="flex min-h-screen flex-col justify-center bg-neutral-50 py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<h2 class="mt-6 text-center text-3xl font-extrabold text-neutral-900">Sign in to Motion Index</h2>
		<p class="mt-2 text-center text-sm text-neutral-600">Access your legal documents repository</p>
	</div>

	<div
		class="mt-8 sm:mx-auto sm:w-full sm:max-w-md"
		in:fly={{ y: 30, duration: 700, delay: 300, easing: cubicOut }}
	>
		<div
			class="bg-white px-4 py-8 shadow sm:rounded-lg sm:px-10"
			in:scale={{ start: 0.97, duration: 600, delay: 400, easing: cubicOut }}
		>
			<form method="POST" class="space-y-6" use:enhance={handleEnhance} in:fade={{ duration: 500, delay: 500 }}>
				<!-- Server-side form errors -->
				{#if form?.error}
					<div class="rounded-md bg-red-50 p-4">
						<div class="flex">
							<div class="ml-3">
								<h3 class="text-sm font-medium text-red-800 whitespace-pre-line">
									{form.error}
								</h3>
								{#if form.error.includes('ad blocker') || form.error.includes('Connection failed')}
									<div class="mt-2 text-xs text-red-600">
										<p class="font-medium">Common solutions:</p>
										<ul class="list-disc list-inside mt-1 space-y-1">
											<li>Disable uBlock Origin, AdBlock Plus, or similar extensions</li>
											<li>Add *.supabase.co to your allowlist</li>
											<li>Try incognito/private browsing mode</li>
											<li>Check your network/firewall settings</li>
										</ul>
									</div>
								{/if}
							</div>
						</div>
					</div>
				{/if}

				<!-- URL parameter errors (e.g., from OAuth) -->
				{#if $page.url.searchParams.get('error')}
					<div class="rounded-md bg-red-50 p-4">
						<div class="flex">
							<div class="ml-3">
								<h3 class="text-sm font-medium text-red-800">
									{decodeURIComponent($page.url.searchParams.get('error') || '')}
								</h3>
								<p class="text-xs text-red-600 mt-1">
									If you're using an ad blocker, please allow Supabase domains or try disabling it temporarily.
								</p>
							</div>
						</div>
					</div>
				{/if}

				{#if successMessage}
					<div class="rounded-md bg-green-50 p-4">
						<div class="flex">
							<div class="ml-3">
								<h3 class="text-sm font-medium text-green-800">
									{successMessage}
								</h3>
							</div>
						</div>
					</div>
				{/if}

				<!-- Network status information -->
				{#if networkStatus}
					<div class="rounded-md bg-blue-50 p-4">
						<div class="flex">
							<div class="ml-3">
								<h3 class="text-sm font-medium text-blue-800 whitespace-pre-line">
									{networkStatus}
								</h3>
							</div>
						</div>
					</div>
				{/if}

				<div>
					<label for="email" class="block text-sm font-medium text-neutral-700">Email address</label>
					<div class="mt-1">
						<input
							id="email"
							name="email"
							type="email"
							autocomplete="email"
							required
							value={form?.email ?? ''}
							class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 placeholder-neutral-400 shadow-sm focus:border-primary-900 focus:outline-none focus:ring-primary-900 sm:text-sm"
						/>
					</div>
				</div>

				<div>
					<label for="password" class="block text-sm font-medium text-neutral-700">Password</label>
					<div class="mt-1">
						<input
							id="password"
							name="password"
							type="password"
							autocomplete="current-password"
							required
							class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 placeholder-neutral-400 shadow-sm focus:border-primary-900 focus:outline-none focus:ring-primary-900 sm:text-sm"
						/>
					</div>
				</div>

				<div class="flex items-center justify-between">
					<div class="text-sm">
						<a href="/auth/forgot" class="font-medium text-primary-900 hover:text-primary-800">
							Forgot your password?
						</a>
					</div>
					<div class="text-sm">
						<button
							type="button"
							onclick={checkNetwork}
							disabled={checkingNetwork}
							class="font-medium text-blue-600 hover:text-blue-500 disabled:opacity-50"
						>
							{#if checkingNetwork}
								Checking...
							{:else}
								Test Connection
							{/if}
						</button>
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
							Signing in...
						{:else}
							Sign in
						{/if}
					</button>
				</div>
			</form>

			<div class="mt-6">
				<div class="relative">
					<div class="absolute inset-0 flex items-center">
						<div class="w-full border-t border-neutral-300"></div>
					</div>
					<div class="relative flex justify-center text-sm">
						<span class="bg-white px-2 text-neutral-500">Don't have an account?</span>
					</div>
				</div>

				<div class="mt-6">
					<a
						href="/auth/register"
						class="flex w-full justify-center rounded-md border border-neutral-300 bg-white px-4 py-2 text-sm font-medium text-neutral-700 shadow-sm hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-900 focus:ring-offset-2"
					>
						Register
					</a>
				</div>
			</div>
		</div>
	</div>
</div>
