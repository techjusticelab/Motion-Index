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

	// Download document function
	function downloadDocument() {
		if (!docData) return;

		// This is a placeholder - implement actual download logic
		dispatch('download', { document: docData });
		console.log('Download document:', docData.file_name);
	}

	// Add to case function
	function addToCase() {
		if (!docData) return;

		// This is a placeholder - implement actual add to case logic
		dispatch('addToCase', { document: docData });
		console.log('Add to case:', docData.file_name);
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
		class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-80 p-4 backdrop-blur-sm"
		on:click={closeDocumentPopup}
	>
		<div
			class="relative flex h-[90vh] w-[95%] max-w-6xl overflow-hidden rounded-xl bg-white/95 shadow-2xl"
			on:click|stopPropagation
		>
			<!-- Metadata Panel (Left Side) -->
			<div
				class="flex h-full w-72 flex-shrink-0 flex-col border-r border-gray-200/60 bg-gray-50/70"
			>
				<!-- Action Buttons -->
				<div class="flex space-x-2 border-b border-gray-200/60 p-3">
					<button
						class="flex flex-1 items-center justify-center rounded-lg bg-blue-500 px-3 py-1.5 text-xs font-medium text-white shadow-sm transition-colors hover:bg-blue-600"
						on:click={downloadDocument}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mr-1 h-3.5 w-3.5"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
							/>
						</svg>
						Download
					</button>
					<button
						class="flex flex-1 items-center justify-center rounded-lg bg-green-500 px-3 py-1.5 text-xs font-medium text-white shadow-sm transition-colors hover:bg-green-600"
						on:click={addToCase}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mr-1 h-3.5 w-3.5"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 6v6m0 0v6m0-6h6m-6 0H6"
							/>
						</svg>
						Add to Case
					</button>
				</div>

				<!-- Metadata Content -->
				<div class="flex-grow overflow-auto p-3">
					<div class="space-y-2.5">
						<!-- Document type -->
						{#if docData.doc_type}
							<div class="rounded-lg border border-gray-200/60 bg-white/80 p-2.5 shadow-sm">
								<span class="mb-0.5 block text-xs text-gray-500">Document Type</span>
								<span class="text-sm font-medium text-gray-800">{docData.doc_type}</span>
							</div>
						{/if}

						<!-- Creation date -->
						{#if docData.created_at}
							<div class="rounded-lg border border-gray-200/60 bg-white/80 p-2.5 shadow-sm">
								<span class="mb-0.5 block text-xs text-gray-500">Created</span>
								<span class="text-sm font-medium text-gray-800"
									>{formatDate(docData.created_at)}</span
								>
							</div>
						{/if}

						<!-- File name -->
						<div class="rounded-lg border border-gray-200/60 bg-white/80 p-2.5 shadow-sm">
							<span class="mb-0.5 block text-xs text-gray-500">File Name</span>
							<span class="break-all text-sm font-medium text-gray-800">{docData.file_name}</span>
						</div>

						<!-- Dynamic metadata fields -->
						{#if docData.metadata}
							{#each Object.entries(docData.metadata) as [key, value]}
								{#if key !== 'document_name' && value !== null && value !== undefined && value !== ''}
									<div class="rounded-lg border border-gray-200/60 bg-white/80 p-2.5 shadow-sm">
										<span class="mb-0.5 block text-xs text-gray-500"
											>{key.replace(/_/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase())}</span
										>
										<span class="block break-words text-sm font-medium text-gray-800">{value}</span>
									</div>
								{/if}
							{/each}
						{/if}
					</div>
				</div>
			</div>

			<!-- Document Content (Right Side) -->
			<div class="flex h-full min-w-0 flex-grow flex-col">
				<!-- Popup Header -->
				<div class="flex items-center justify-between border-b border-gray-200/60 p-3">
					<div class="max-w-[calc(100%-40px)] truncate px-2">
						<h2 class="truncate text-base font-medium text-gray-800">
							{docData.metadata?.document_name || docData.file_name}
						</h2>
						{#if docData.doc_type}
							<span class="text-xs text-gray-600">{docData.doc_type}</span>
						{/if}
					</div>
					<div class="flex items-center">
						<button
							class="rounded-full p-1.5 text-gray-500 transition-colors hover:bg-gray-100/70"
							on:click={closeDocumentPopup}
							aria-label="Close document viewer"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-5 w-5"
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

				<!-- Document Content -->
				<div class="flex-grow overflow-auto">
					{#if docData.text}
						<!-- Text content viewer with styled layout -->
						<div class="h-full overflow-auto bg-gray-50/40 p-5 text-gray-800">
							<div class="mx-auto max-w-4xl rounded-lg bg-white/90 p-6 shadow-sm">
								<pre
									class="whitespace-pre-wrap font-sans text-sm leading-relaxed">{docData.text}</pre>
							</div>
						</div>
					{:else}
						<div class="flex h-full flex-col items-center justify-center p-8 text-center">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="mb-4 h-14 w-14 text-gray-300"
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
	</div>
{/if}
