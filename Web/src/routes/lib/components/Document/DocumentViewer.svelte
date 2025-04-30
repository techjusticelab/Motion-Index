<script lang="ts">
	import { onMount, afterUpdate } from 'svelte';
	import { browser } from '$app/environment'; // Import browser check from SvelteKit

	// Props - renamed from 'document' to 'docData' to avoid conflict with global document object
	export let docData: {
		metadata: { document_name: string; [key: string]: any };
		file_name: string;
		s3_uri: string;
		text?: string;
		doc_type?: string;
		created_at?: string;
	} | null = null;

	export let isOpen = false;

	// State
	let showMetadata = true;
	let textContentFormatted = '';
	let previousDocumentId = null;
	let mounted = false;

	// Set up dispatch
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	// Watch for document changes and process text when document changes
	$: if (docData && isOpen && mounted) {
		const currentDocumentId = docData.file_name || docData.s3_uri;
		if (currentDocumentId !== previousDocumentId) {
			previousDocumentId = currentDocumentId;
			setupDocument(docData);
		}
	}

	// Reset state when popup is closed
	$: if (!isOpen && mounted) {
		resetState();
	}

	// Set up document display
	function setupDocument(doc) {
		if (!doc) return;

		// Prevent body scrolling when popup is open - only in browser
		if (browser) {
			document.body.style.overflow = 'hidden';
		}

		// Format the text content to make it more readable
		if (doc.text) {
			textContentFormatted = doc.text;
		}
	}

	// Reset component state
	function resetState() {
		showMetadata = false;
		textContentFormatted = '';
		previousDocumentId = null;

		// Ensure body scrolling is restored - only in browser
		if (browser) {
			document.body.style.overflow = 'auto';
		}
	}

	// Close document popup - dispatch event to parent
	function closeDocumentPopup() {
		console.log('closeDocumentPopup');
		isOpen = false;
		// Let the parent know to reset activeDocument and showDocumentPopup
		dispatch('close');
	}

	// Toggle metadata display
	function toggleMetadata() {
		showMetadata = !showMetadata;
	}

	// Format date
	function formatDate(dateString: string): string {
		if (!dateString) return 'N/A';
		try {
			const date = new Date(dateString);
			return date.toLocaleDateString('en-US', {
				year: 'numeric',
				month: 'short',
				day: 'numeric'
			});
		} catch (err) {
			return dateString;
		}
	}

	// Only execute client-side code after component is mounted
	onMount(() => {
		mounted = true;

		// Cleanup when component is destroyed
		return () => {
			// Ensure body scrolling is restored when component is destroyed
			if (browser) {
				document.body.style.overflow = 'auto';
			}
		};
	});
</script>

<!-- Document Popup -->
{#if docData && isOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-75 p-4"
		on:click={closeDocumentPopup}
	>
		<div
			class="relative h-[85vh] w-full max-w-5xl rounded-lg bg-white shadow-2xl"
			on:click|stopPropagation
		>
			<!-- Popup Header -->
			<div class="flex items-center justify-between border-b border-gray-200 p-4">
				<div>
					<h2 class="truncate text-lg font-semibold text-gray-800">
						{docData.metadata?.document_name || docData.file_name}
					</h2>
					{#if docData.doc_type}
						<span class="text-sm text-gray-600">{docData.doc_type}</span>
					{/if}
				</div>
				<div class="flex items-center space-x-2">
					<button
						class="rounded-lg bg-blue-50 px-3 py-1 text-sm font-medium text-blue-600 hover:bg-blue-100"
						on:click|stopPropagation={toggleMetadata}
						aria-label={showMetadata ? 'Hide document metadata' : 'Show document metadata'}
					>
						{showMetadata ? 'Hide Metadata' : 'Show Metadata'}
					</button>
					<button
						class="rounded-full p-2 text-gray-500 transition-colors hover:bg-gray-100"
						on:click={closeDocumentPopup}
						aria-label="Close document viewer"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M6 18L18 6M6 6l12 12"
							/>
						</svg>
					</button>
				</div>
			</div>

			<!-- Metadata Panel (conditionally shown) -->
			{#if showMetadata && docData}
				<div class="border-b border-gray-200 bg-gray-50 p-4">
					<h3 class="mb-2 text-sm font-medium text-gray-700">Document Metadata</h3>
					<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 md:grid-cols-3">
						<!-- Document type -->
						{#if docData.doc_type}
							<div class="rounded border border-gray-200 bg-white p-2 shadow-sm">
								<span class="block text-xs text-gray-500">Document Type</span>
								<span class="font-medium text-gray-800">{docData.doc_type}</span>
							</div>
						{/if}

						<!-- Creation date -->
						{#if docData.created_at}
							<div class="rounded border border-gray-200 bg-white p-2 shadow-sm">
								<span class="block text-xs text-gray-500">Created</span>
								<span class="font-medium text-gray-800">{formatDate(docData.created_at)}</span>
							</div>
						{/if}

						<!-- Dynamic metadata fields -->
						{#if docData.metadata}
							{#each Object.entries(docData.metadata) as [key, value]}
								{#if key !== 'document_name' && value !== null && value !== undefined && value !== ''}
									<div class="rounded border border-gray-200 bg-white p-2 shadow-sm">
										<span class="block text-xs text-gray-500"
											>{key.replace(/_/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase())}</span
										>
										<span class="block truncate font-medium text-gray-800">{value}</span>
									</div>
								{/if}
							{/each}
						{/if}
					</div>
				</div>
			{/if}

			<!-- Document Content -->
			<div class="h-[calc(100%-4rem-1px)] w-full {showMetadata ? 'h-[calc(100%-13rem-2px)]' : ''}">
				{#if docData.text}
					<!-- Text content viewer with styled layout -->
					<div class="h-full overflow-auto bg-gray-50 p-6 text-gray-800">
						<div class="mx-auto max-w-4xl rounded-lg bg-white p-8 shadow-sm">
							<pre
								class="whitespace-pre-wrap font-serif text-base leading-relaxed">{docData.text}</pre>
						</div>
					</div>
				{:else}
					<div class="flex h-full flex-col items-center justify-center p-8 text-center">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mb-4 h-16 w-16 text-gray-300"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							aria-hidden="true"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
							/>
						</svg>
						<p class="mb-2 text-gray-600">Document content not available</p>
						<p class="text-sm text-gray-500">
							This document doesn't have any text content to display.
						</p>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}
