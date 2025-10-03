<script lang="ts">
	import { scale, slide, fly } from 'svelte/transition';
	import { elasticOut, cubicOut, backOut } from 'svelte/easing';
	import type { Snippet } from 'svelte';
	
	interface Props {
		user?: { email?: string; user_metadata?: any };
		casesCount?: number;
		documentsCount?: number;
		isLoadingDetails?: boolean;
		displayName?: string;
		isInitialLoad?: boolean;
		class?: string;
		actions?: Snippet;
	}

	let {
		user = null,
		casesCount = 0,
		documentsCount = 0,
		isLoadingDetails = false,
		displayName = '',
		isInitialLoad = true,
		class: className = '',
		actions
	}: Props = $props();

	// Get user display name
	function getUserDisplayName(): string {
		if (displayName) return displayName;
		if (user?.user_metadata?.full_name) return user.user_metadata.full_name;
		if (user?.user_metadata?.name) return user.user_metadata.name;
		return 'User Account';
	}
</script>

<div class="space-y-6 {className}">
	<!-- Header -->
	<h2
		class="text-xl font-semibold text-neutral-800"
		in:slide={{ duration: 500, delay: isInitialLoad ? 200 : 0 }}
	>
		Account Information
	</h2>

	<!-- User profile card -->
	<div class="overflow-hidden rounded-lg bg-white p-4 shadow-sm">
		<div class="flex items-center">
			<div
				class="flex h-16 w-16 items-center justify-center rounded-full bg-primary-100 text-primary-600"
				in:scale={{
					start: 0.9,
					duration: 600,
					delay: isInitialLoad ? 300 : 0,
					easing: elasticOut
				}}
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-8 w-8"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
					aria-hidden="true"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
					/>
				</svg>
			</div>
			<div class="ml-4" in:slide={{ duration: 500, delay: isInitialLoad ? 400 : 0 }}>
				<h3 class="text-lg font-medium text-neutral-800">
					{user?.email || 'User'}
				</h3>
				<p class="text-sm text-neutral-500">
					{isLoadingDetails ? 'Loading details...' : getUserDisplayName()}
				</p>
			</div>
		</div>
	</div>

	<!-- Account stats -->
	<div class="grid grid-cols-2 gap-4">
		<div
			class="rounded-lg bg-white p-4 shadow-sm"
			in:fly={{ y: 15, duration: 500, delay: isInitialLoad ? 500 : 0, easing: cubicOut }}
		>
			<h4 class="text-sm font-medium text-neutral-500">Total Cases</h4>
			<p class="mt-1 text-2xl font-semibold text-primary-600">{casesCount}</p>
		</div>
		<div
			class="rounded-lg bg-white p-4 shadow-sm"
			in:fly={{ y: 15, duration: 500, delay: isInitialLoad ? 600 : 0, easing: cubicOut }}
		>
			<h4 class="text-sm font-medium text-neutral-500">Total Documents</h4>
			<p class="mt-1 text-2xl font-semibold text-primary-600">{documentsCount}</p>
		</div>
	</div>

	<!-- Account actions -->
	<div class="space-y-3" in:slide={{ duration: 500, delay: isInitialLoad ? 700 : 0 }}>
		<h3 class="text-md font-semibold text-neutral-700">Account Actions</h3>
		
		{#if actions}
			{@render actions()}
		{:else}
			<!-- Default actions content -->
			<button
				class="flex w-full items-center justify-between rounded-lg border border-neutral-300 bg-white p-3 text-left text-sm font-medium text-neutral-700 shadow-sm hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
				in:scale={{
					start: 0.95,
					duration: 400,
					delay: isInitialLoad ? 800 : 0,
					easing: backOut
				}}
			>
				<span class="flex items-center">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mr-2 h-5 w-5 text-neutral-400"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
						/>
						<path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
					</svg>
					Account Settings
				</span>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-5 w-5 text-neutral-400"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
					aria-hidden="true"
				>
					<path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
				</svg>
			</button>
		{/if}
	</div>
</div>