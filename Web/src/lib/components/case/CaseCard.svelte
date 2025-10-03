<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { fly } from 'svelte/transition';
	import type { Case } from '$lib/supabase';

	const dispatch = createEventDispatcher();

	interface Props {
		case: Case;
		index?: number;
		showActions?: boolean;
		class?: string;
	}

	let {
		case: caseItem,
		index = 0,
		showActions = true,
		class: className = ''
	}: Props = $props();

	// Format date for display
	function formatDate(dateString: string | number | Date) {
		if (!dateString) return 'N/A';
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	// Handle case click
	function handleClick() {
		dispatch('view', { id: caseItem.id });
	}

	// Handle view action
	function handleView(event: MouseEvent) {
		event.stopPropagation();
		dispatch('view', { id: caseItem.id });
	}

	// Handle delete action
	function handleDelete(event: MouseEvent) {
		event.stopPropagation();
		dispatch('delete', { id: caseItem.id, name: caseItem.case_name });
	}
</script>

<div
	class="group relative p-3 hover:bg-neutral-50 transition-colors cursor-pointer {className}"
	in:fly={{ y: 20, duration: 400, delay: index * 50 }}
	onclick={handleClick}
	role="button"
	tabindex="0"
	aria-label="View case {caseItem.case_name}"
	onkeydown={(e) => e.key === 'Enter' && handleClick()}
>
	<div class="flex items-center justify-between">
		<div class="flex items-center space-x-3 flex-1 min-w-0">
			<!-- Case Icon -->
			<div class="flex h-6 w-6 items-center justify-center rounded-lg bg-primary-100 flex-shrink-0">
				<svg
					class="h-3 w-3 text-primary-600"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
					aria-hidden="true"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
					/>
				</svg>
			</div>

			<!-- Case Info -->
			<div class="flex-1 min-w-0">
				<h3 class="text-xs font-medium text-neutral-900 group-hover:text-primary-600 break-words">
					{caseItem.case_name}
				</h3>
				<div class="mt-1 flex items-center space-x-3 text-xs text-neutral-500">
					<span>Created {formatDate(caseItem.created_at)}</span>
					<span>Updated {formatDate(caseItem.updated_at)}</span>
				</div>
			</div>
		</div>

		<!-- Actions -->
		{#if showActions}
			<div class="flex items-center space-x-2 flex-shrink-0">
				<button
					onclick={handleView}
					class="inline-flex items-center text-xs font-medium text-primary-600 hover:text-primary-800 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 rounded-sm"
					aria-label="View case {caseItem.case_name}"
				>
					View
					<svg class="ml-1 h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
					</svg>
				</button>
				<button
					onclick={handleDelete}
					class="text-xs text-red-600 hover:text-red-800 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 rounded-sm px-1 py-0.5"
					aria-label="Delete case {caseItem.case_name}"
				>
					Delete
				</button>
			</div>
		{/if}
	</div>
</div>