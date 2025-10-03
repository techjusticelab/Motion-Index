<!-- DocumentPDFViewer.svelte -->
<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Document } from '$lib/types';
	import { getDocumentUrl, getDocumentUrlWithSearch } from '$lib/api';

	const {
		docData,
		session = null
	}: {
		docData?: Document | null;
		session?: any;
	} = $props();

	let isLoading = $state(true);
	let errorMessage = $state('');
	let url = $state('');

	const dispatch = createEventDispatcher<{
		loaded: void;
		error: { message: string };
	}>();

	// Track document ID to prevent infinite loops when same document reloads
	let currentDocumentId = $state('');

	// Reactive effect to update URL when docData changes
	$effect(() => {
		if (docData) {
			const newDocumentId = docData.id || docData.file_name || JSON.stringify(docData);
			if (newDocumentId !== currentDocumentId) {
				currentDocumentId = newDocumentId;
				loadDocument();
			}
		}
	});

	async function loadDocument() {
		if (!docData) return;

		isLoading = true;
		errorMessage = '';
		
		try {
			// First try the standard URL generation
			let documentUrl = getDocumentUrl(docData);
			
			// If no URL could be generated, try the search fallback
			if (!documentUrl) {
				documentUrl = await getDocumentUrlWithSearch(docData, session);
			}
			
			if (documentUrl) {
				url = documentUrl;
				isLoading = false;
				dispatch('loaded');
			} else {
				errorMessage = 'Document file not found. This document may need to be re-uploaded.';
				isLoading = false;
				dispatch('error', { message: errorMessage });
			}
		} catch (error) {
			console.error('Error loading document:', error);
			errorMessage = 'Error loading document. Please try again.';
			isLoading = false;
			dispatch('error', { message: errorMessage });
		}
	}
</script>

<!-- PDF Viewer Content -->
<div class="h-full w-full bg-neutral-100">
	{#if isLoading}
		<div class="flex h-full w-full items-center justify-center p-6">
			<div class="flex flex-col items-center space-y-4">
				<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-200 border-t-primary-600"></div>
				<p class="text-neutral-600">Loading document...</p>
			</div>
		</div>
	{:else if url}
		<object
			data={url}
			type="application/pdf"
			class="h-full w-full"
		>
			<!-- Fallback if object tag fails -->
			<div class="flex h-full w-full items-center justify-center p-6">
				<div class="max-w-md rounded-lg border border-blue-100 bg-blue-50 p-6 text-center">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mx-auto mb-4 h-12 w-12 text-blue-400"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
						/>
					</svg>
					<p class="text-blue-700 mb-4 font-medium">
						PDF Preview Not Available
					</p>
					<p class="text-blue-600 mb-4 text-sm">
						Your browser doesn't support PDF preview or the server is blocking frame access.
					</p>
					<div class="space-y-2">
						<a
							href={url}
							target="_blank"
							rel="noopener noreferrer"
							class="block w-full rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
						>
							Open in New Tab
						</a>
						<a
							href={url}
							download={docData?.file_name}
							class="block w-full rounded-lg border border-blue-600 px-4 py-2 text-sm font-medium text-blue-600 hover:bg-blue-50"
						>
							Download PDF
						</a>
					</div>
				</div>
			</div>
		</object>
	{:else}
		<div class="flex h-full w-full items-center justify-center p-6">
			<div class="max-w-md rounded-lg border border-primary-100 bg-primary-50 p-6 text-center">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mx-auto mb-4 h-12 w-12 text-primary-400"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
					/>
				</svg>
				<p class="text-neutral-700 mb-4">
					{errorMessage || 'Unable to load document.'}
				</p>
				{#if url}
					<a
						href={url}
						download={docData?.file_name}
						class="inline-flex items-center rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mr-2 h-4 w-4"
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
						Download File
					</a>
				{/if}
			</div>
		</div>
	{/if}
</div>