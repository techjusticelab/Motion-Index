<script lang="ts">
	import { fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import type { Snippet } from 'svelte';

	interface Props {
		icon?: 'cases' | 'documents' | 'search' | 'generic';
		title: string;
		description?: string;
		actionText?: string;
		delay?: number;
		class?: string;
		action?: Snippet;
	}

	let {
		icon = 'generic',
		title,
		description = '',
		actionText = '',
		delay = 0,
		class: className = '',
		action
	}: Props = $props();

	// Icon map for different empty states
	const icons = {
		cases: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />`,
		documents: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />`,
		search: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />`,
		generic: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />`
	};
</script>

<div
	class="rounded-lg border border-dashed border-neutral-300 bg-neutral-50 p-6 text-center {className}"
	in:fly={{ y: 20, duration: 600, delay, easing: cubicOut }}
>
	<!-- Icon -->
	<svg
		xmlns="http://www.w3.org/2000/svg"
		class="mx-auto h-12 w-12 text-neutral-400"
		fill="none"
		viewBox="0 0 24 24"
		stroke="currentColor"
		aria-hidden="true"
	>
		{@html icons[icon]}
	</svg>

	<!-- Content -->
	<h3 class="mt-2 text-sm font-medium text-neutral-900">{title}</h3>
	
	{#if description}
		<p class="mt-1 text-sm text-neutral-500">{description}</p>
	{/if}

	{#if actionText}
		<div class="mt-4">
			{#if action}
				{@render action()}
			{:else}
				<button
					class="inline-flex items-center rounded-md border border-transparent bg-primary-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
				>
					{actionText}
				</button>
			{/if}
		</div>
	{/if}
</div>