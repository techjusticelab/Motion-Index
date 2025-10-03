<!-- DocumentViewer.svelte -->
<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { formatDate } from '$lib/utils';
	import type { Document, SearchResponse } from '$lib/types';
	import { getDocumentUrl } from '$lib/api';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { elasticOut, backOut, quintOut, cubicOut } from 'svelte/easing';
	import { CaseManager, type Case } from '$lib/supabase';
	import { page } from '$app/stores';

	const {
		docData = $bindable(null),
		isOpen = $bindable(false),
		supabase = null,
		session = null,
		user = null
	}: {
		docData?: Document | null;
		isOpen?: boolean;
		supabase?: any;
		session?: any;
		user?: any;
	} = $props();

	let isLoading = $state(true);
	let errorMessage = $state('');
	let url = $state('');

	// Case management
	let cases = $state<Case[]>([]);
	let caseManager = $state<CaseManager | null>(null);
	let showAddToCaseModal = $state(false);
	let selectedCaseId = $state('');
	let documentNotes = $state('');
	let isAddingToCase = $state(false);
	let showNewCaseModal = $state(false);
	let newCaseName = $state('');
	let isCreatingCase = $state(false);
	let successMessage = $state('');

	const dispatch = createEventDispatcher<{
		close: void;
	}>();

	// Generate a UUID v5 from document data for database compatibility
	function generateDocumentUUID(docData: Document): string {
		// Create a consistent identifier from document data
		const identifier = docData.s3_uri || docData.file_name || JSON.stringify(docData);

		// Simple hash function to create consistent UUID-like string
		let hash = 0;
		for (let i = 0; i < identifier.length; i++) {
			const char = identifier.charCodeAt(i);
			hash = ((hash << 5) - hash) + char;
			hash = hash & hash; // Convert to 32-bit integer
		}

		// Convert hash to positive number and format as UUID
		const positiveHash = Math.abs(hash).toString(16).padStart(8, '0');
		const uuid = `${positiveHash.slice(0, 8)}-${positiveHash.slice(0, 4)}-4${positiveHash.slice(1, 4)}-8${positiveHash.slice(1, 4)}-${positiveHash}${positiveHash.slice(0, 4)}`;

		return uuid;
	}

	// Extract document ID from S3 URI or generate a consistent ID
	function getDocumentId(docData: Document): string {
		// Try to use id if available and it's a valid UUID format
		if (docData.id && docData.id.match(/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i)) {
			return docData.id;
		}

		// If database expects UUIDs, generate one consistently
		return generateDocumentUUID(docData);
	}

	// Function to determine if file can be viewed in iframe
	function canViewInIframe(fileName: string): boolean {
		if (!fileName) return false;
		console.log(fileName);
		console.log(docData);
		console.log('DSAD', docData?.s3_uri);
		const extension = fileName.split('.').pop()?.toLowerCase() || '';
		// Most browsers can display these formats directly in an iframe
		return ['pdf', 'docx', 'txt', 'jpg', 'jpeg', 'png', 'gif'].includes(extension);
	}

	// Function to close the viewer
	function closeViewer() {
		dispatch('close');
	}
	// Reactive effect to update URL when docData changes
	$effect(() => {
		if (docData) {
			downloadFile();
		}
	});

	onMount(async () => {
		if (supabase && user) {
			caseManager = new CaseManager(supabase);

			// Test database access to debug the issue
			await caseManager.testDatabaseAccess();

			loadUserCases();
		}
	});
	async function downloadFile() {
		if (docData) {
			// Use the new getDocumentUrl function that handles local document server
			const documentUrl = getDocumentUrl(docData);
			if (documentUrl) {
				url = documentUrl;
				console.log('Document URL from local server:', url);
				isLoading = false; // Set loading to false when URL is ready
			} else {
				console.error('Failed to get document URL for:', docData);
				console.log('Document data:', {
					file_path: docData.file_path,
					s3_uri: docData.s3_uri,
					file_url: docData.file_url,
					file_name: docData.file_name
				});
				errorMessage = 'Document file not found in local storage. This document may need to be re-uploaded or the file path may be outdated.';
				isLoading = false;
			}
		}
	}
	// Handle S3FileViewer events
	function handleViewerLoaded() {
		isLoading = false;
	}

	function handleViewerError(event: CustomEvent<{ message: string }>) {
		isLoading = false;
		errorMessage =
			event.detail.message ||
			'Unable to display this document in the browser. Please download the file to view it.';
	}

	// Load user's cases
	async function loadUserCases() {
		if (!caseManager || !user) return;
		try {
			cases = await caseManager.getUserCases(user.id);
		} catch (error) {
			console.error('Error loading cases:', error);
		}
	}

	// Open add to case modal
	function openAddToCaseModal() {
		selectedCaseId = '';
		documentNotes = '';
		showAddToCaseModal = true;
	}

	// Close add to case modal
	function closeAddToCaseModal() {
		showAddToCaseModal = false;
		selectedCaseId = '';
		documentNotes = '';
	}

	// Add document to selected case
	async function addDocumentToCase() {
		if (!caseManager || !docData || !selectedCaseId) {
			console.error('Missing required data:', {
				caseManager: !!caseManager,
				docData: !!docData,
				selectedCaseId: !!selectedCaseId,
				docDataId: docData?.id
			});
			return;
		}

		console.log('Adding document to case:', {
			caseId: selectedCaseId,
			documentId: docData.id,
			notes: documentNotes.trim()
		});

		console.log('Full docData object:', docData);
		console.log('Available docData properties:', Object.keys(docData));

		isAddingToCase = true;
		try {
			// Extract proper document ID
			const documentId = getDocumentId(docData);
			console.log('Using document ID:', documentId);

			const result = await caseManager.addDocumentToCase(
				selectedCaseId,
				documentId,
				documentNotes.trim() || undefined
			);
			console.log('Document added successfully:', result);

			// Find the case name for the success message
			const selectedCase = cases.find(c => c.id === selectedCaseId);
			successMessage = `Document added to case "${selectedCase?.case_name || 'Unknown'}" successfully!`;

			// Clear success message after 3 seconds
			setTimeout(() => {
				successMessage = '';
			}, 3000);

			closeAddToCaseModal();
		} catch (error) {
			console.error('Error adding document to case:', error);
		} finally {
			isAddingToCase = false;
		}
	}

	// Open new case modal
	function openNewCaseModal() {
		newCaseName = '';
		showNewCaseModal = true;
	}

	// Close new case modal
	function closeNewCaseModal() {
		showNewCaseModal = false;
		newCaseName = '';
	}

	// Create new case and add document
	async function createCaseAndAddDocument() {
		if (!caseManager || !user || !docData || !newCaseName.trim()) return;
		
		isCreatingCase = true;
		try {
			const newCase = await caseManager.createCase(user.id, newCaseName.trim());
			if (newCase) {
				cases = [newCase, ...cases];

				// Extract proper document ID
				try {
					const documentId = getDocumentId(docData);
					console.log('Using document ID for new case:', documentId);
					await caseManager.addDocumentToCase(newCase.id, documentId, documentNotes.trim() || undefined);
				} catch (error) {
					console.error('Error getting document ID for new case:', error);
				}

				closeNewCaseModal();
			}
		} catch (error) {
			console.error('Error creating case:', error);
		} finally {
			isCreatingCase = false;
		}
	}


	// Handle keydown for modal escape
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			if (showAddToCaseModal) {
				closeAddToCaseModal();
			} else if (showNewCaseModal) {
				closeNewCaseModal();
			}
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- TODO: ADD SIMILAR FEATURES HERE SO WHEN YOU CLICK ON METADATA IT WILL AUTO SEARCH FOR SIMILAR SHIT -->
{#if isOpen && docData}
	<div
		class="opac fixed inset-0 z-50 flex justify-center p-4 shadow"
		in:fade={{ duration: 300, easing: cubicOut }}
		out:fade={{ duration: 500 }}
	>
		<div
			class="relative flex h-[95vh] w-[95vw] overflow-hidden text-wrap rounded-xl bg-white shadow-2xl"
			transition:fly={{ y: 20, duration: 800, easing: quintOut }}
		>
			<!-- Sidebar with metadata -->
			<div
				class="w-1/4 overflow-auto border-r border-neutral-200 bg-neutral-50 p-6"
				in:fly={{ x: -20, duration: 800, delay: 300, easing: cubicOut }}
			>
				<div class="mb-4">
					<h2
						class="truncate text-xl font-semibold text-neutral-800"
						in:slide={{ duration: 700, delay: 400 }}
					>
						{docData.metadata.document_name || docData.file_name}
					</h2>
					<p class="mt-1 text-sm text-neutral-500" in:slide={{ duration: 700, delay: 500 }}>
						{docData.file_name}
					</p>
				</div>

				<!-- Document type tag -->
				<div class="mb-4">
					<span
						class="inline-flex rounded-full bg-primary-50 px-3 py-1 text-sm font-medium text-primary-700"
						in:scale={{ start: 0.9, duration: 600, delay: 600, easing: cubicOut }}
					>
						{docData.doc_type}
					</span>
				</div>

				<!-- Metadata list -->
				<div class="space-y-4 overflow-hidden" in:slide={{ duration: 600, delay: 650 }}>
					<h3 class="text-sm font-medium text-neutral-700">Document Details</h3>

					<div class="space-y-2 rounded-lg border border-neutral-200 bg-white p-3">
						<!-- Document Name - always displayed -->
						<div in:slide={{ duration: 500, delay: 700 }}>
							<p class="text-xs text-neutral-500">Document Name</p>
							<p class="text-sm font-medium text-neutral-800">{docData.metadata.document_name}</p>
						</div>

						{#if docData.metadata.case_number}
							<div in:slide={{ duration: 500, delay: 750 }}>
								<p class="text-xs text-neutral-500">Case Number</p>
								<p class="text-sm font-medium text-neutral-800">{docData.metadata.case_number}</p>
							</div>
						{/if}

						{#if docData.metadata.case_name}
							<div in:slide={{ duration: 500, delay: 800 }}>
								<p class="text-xs text-neutral-500">Case Name</p>
								<p class="text-sm font-medium text-neutral-800">{docData.metadata.case_name}</p>
							</div>
						{/if}

						{#if docData.metadata.court}
							<div in:slide={{ duration: 500, delay: 850 }}>
								<p class="text-xs text-neutral-500">Court</p>
								<p class="text-sm font-medium text-neutral-800">{docData.metadata.court}</p>
							</div>
						{/if}

						{#if docData.metadata.judge}
							<div in:slide={{ duration: 500, delay: 900 }}>
								<p class="text-xs text-neutral-500">Judge</p>
								<p class="text-sm font-medium text-neutral-800">{docData.metadata.judge}</p>
							</div>
						{/if}

						{#if docData.metadata.timestamp}
							<div in:slide={{ duration: 500, delay: 950 }}>
								<p class="text-xs text-neutral-500">Date</p>
								<p class="text-sm font-medium text-neutral-800">
									{formatDate(docData.metadata.timestamp)}
								</p>
							</div>
						{:else}
							<div in:slide={{ duration: 500, delay: 950 }}>
								<p class="text-xs text-neutral-500">Created Date</p>
								<p class="text-sm font-medium text-neutral-800">
									{formatDate(docData.created_at)}
								</p>
							</div>
						{/if}

						{#if docData.metadata.subject}
							<div in:slide={{ duration: 500, delay: 1000 }}>
								<p class="text-xs text-neutral-500">Subject</p>
								<p class="text-sm font-medium text-neutral-800">{docData.metadata.subject}</p>
							</div>
						{/if}

						{#if docData.metadata.status}
							<div in:slide={{ duration: 500, delay: 1050 }}>
								<p class="text-xs text-neutral-500">Status</p>
								<p class="text-sm font-medium text-neutral-800">{docData.metadata.status}</p>
							</div>
						{/if}

						{#if docData.metadata.author}
							<div in:slide={{ duration: 500, delay: 1100 }}>
								<p class="text-xs text-neutral-500">Author</p>
								<p class="text-sm font-medium text-neutral-800">{docData.metadata.author}</p>
							</div>
						{/if}

						{#if docData.metadata.legal_tags && docData.metadata.legal_tags.length > 0}
							<div in:slide={{ duration: 500, delay: 1150 }}>
								<p class="text-xs text-neutral-500">Legal Tags</p>
								<div class="mt-1 flex flex-wrap gap-1">
									{#each docData.metadata.legal_tags as tag, i}
										<span
											class="inline-flex rounded-full bg-neutral-100 px-2 py-0.5 text-xs font-medium text-neutral-800"
											in:scale={{
												start: 0.9,
												duration: 500,
												delay: 1200 + i * 100,
												easing: cubicOut
											}}
										>
											{tag}
										</span>
									{/each}
								</div>
							</div>
						{/if}
					</div>

					<!-- Success Message -->
					{#if successMessage}
						<div
							class="mb-4 rounded-lg bg-green-50 border border-green-200 p-3"
							in:slide={{ duration: 300 }}
							out:slide={{ duration: 300 }}
						>
							<div class="flex items-center">
								<svg
									class="h-4 w-4 text-green-400 mr-2"
									fill="currentColor"
									viewBox="0 0 20 20"
								>
									<path
										fill-rule="evenodd"
										d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
										clip-rule="evenodd"
									/>
								</svg>
								<p class="text-sm font-medium text-green-800">{successMessage}</p>
							</div>
						</div>
					{/if}

					<!-- Case Actions -->
					{#if session && supabase}
						<div class="m-auto flex flex-row justify-center gap-2 align-middle">
							<button
								onclick={openAddToCaseModal}
								class="mt-4 flex w-full items-center justify-center rounded-lg border-primary-600 bg-white px-3 py-2 text-xs font-medium text-black hover:bg-primary-600 hover:text-white transition-colors"
								in:scale={{ start: 0.5, duration: 300, delay: 1300, easing: cubicOut }}
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="mr-1 h-3 w-3"
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
								Add to Case
							</button>
							<button
								onclick={openNewCaseModal}
								class="mt-4 flex w-full items-center justify-center rounded-lg bg-primary-600 px-3 py-2 text-xs font-medium text-white hover:bg-primary-700 transition-colors"
								in:scale={{ start: 0.5, duration: 300, delay: 1350, easing: cubicOut }}
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="mr-1 h-3 w-3"
									fill="none"
									viewBox="0 0 24 24"
									stroke="currentColor"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M12 4v16m8-8H4"
									/>
								</svg>
								New Case
							</button>
						</div>
						<div class="mt-2 flex justify-center">
							<a
								href={url}
								download={docData.file_name}
								class="flex items-center text-xs text-primary-600 hover:text-primary-800 underline"
								in:scale={{ start: 0.5, duration: 300, delay: 1400, easing: cubicOut }}
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="mr-1 h-3 w-3"
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
								Download Original
							</a>
						</div>
					{/if}
				</div>
			</div>

			<!-- Document Viewer area -->
			<div
				class="realtive w-full flex-1 overflow-hidden"
				in:fly={{ x: 20, duration: 800, delay: 350, easing: cubicOut }}
			>
				<!-- Header with close button -->
				<div
					class="absolute right-0 top-0 z-10 flex items-center justify-end p-4"
					in:fly={{ y: -10, duration: 700, delay: 800, easing: cubicOut }}
				>
					<button
						type="button"
						class="rounded-full bg-white/90 p-2 shadow-md hover:bg-neutral-100"
						onclick={closeViewer}
						in:scale={{ start: 0.9, duration: 600, delay: 900 }}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6 text-neutral-500"
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

				<!-- Document content -->
				<div class="h-full w-full overflow-auto bg-neutral-100">
					{#if isLoading}
						<div class="flex h-full w-full items-center justify-center p-6">
							<div class="flex flex-col items-center space-y-4">
								<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-200 border-t-primary-600"></div>
								<p class="text-neutral-600">Loading document...</p>
							</div>
						</div>
					{:else if url}
						<!-- Direct PDF/document display -->
						<embed
							src={url}
							type="application/pdf"
							width="100%"
							height="100%"
							class="h-full w-full"
							onload={handleViewerLoaded}
							onerror={handleViewerError}
						/>
					{:else}
						<div class="flex h-full w-full items-center justify-center p-6">
							<div
								class="max-w-md rounded-lg border border-primary-100 bg-primary-50 p-6 text-center"
								in:scale={{ start: 0.95, duration: 700, delay: 700, easing: cubicOut }}
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="mx-auto mb-4 h-12 w-12 text-primary-400"
									fill="none"
									viewBox="0 0 24 24"
									stroke="currentColor"
									in:scale={{ start: 0.8, duration: 800, delay: 800, easing: cubicOut }}
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
									/>
								</svg>
								<p class="text-neutral-700" in:slide={{ duration: 600, delay: 900 }}>
									This file type cannot be previewed directly in the browser.
								</p>
								<a
									href={url}
									download={docData.file_name}
									class="mt-4 inline-flex items-center rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700"
									in:scale={{ start: 0.98, duration: 700, delay: 1000, easing: cubicOut }}
								>
									Download File
								</a>
							</div>
						</div>
					{/if}

					{#if errorMessage}
						<div
							class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-90"
							transition:fade={{ duration: 600 }}
						>
							<div
								class="max-w-md rounded-lg border border-red-100 bg-red-50 p-6 text-center"
								in:scale={{ start: 0.95, duration: 700, easing: cubicOut }}
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="mx-auto mb-4 h-12 w-12 text-red-400"
									fill="none"
									viewBox="0 0 24 24"
									stroke="currentColor"
									in:scale={{ start: 0.8, duration: 800, delay: 150, easing: cubicOut }}
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
									/>
								</svg>
								<p class="text-neutral-700" in:slide={{ duration: 600, delay: 300 }}>
									{errorMessage}
								</p>
								<a
									href={url}
									download={docData.file_name}
									class="mt-4 inline-flex items-center rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700"
									in:scale={{ start: 0.98, duration: 700, delay: 450, easing: cubicOut }}
								>
									Download File
								</a>
							</div>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Add to Case Modal -->
{#if showAddToCaseModal}
	<div
		class="fixed inset-0 z-[60] overflow-y-auto"
		in:fade={{ duration: 200 }}
		out:fade={{ duration: 200 }}
	>
		<div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center">
			<!-- Backdrop -->
			<div class="fixed inset-0 bg-neutral-500 bg-opacity-75 transition-opacity" onclick={closeAddToCaseModal}></div>

			<!-- Modal Content -->
			<div
				class="relative inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:align-middle sm:max-w-lg sm:w-full z-10"
				in:scale={{ start: 0.95, duration: 200 }}
				out:scale={{ start: 1, end: 0.95, duration: 200 }}
			>
				<div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
					<h3 class="text-lg font-medium text-neutral-900 mb-4">Add Document to Case</h3>
					<div class="space-y-4">
						<div>
							<label for="case-select" class="block text-sm font-medium text-neutral-700">Select Case</label>
							<div class="mt-1">
								<select
									id="case-select"
									bind:value={selectedCaseId}
									class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
								>
									<option value="">Choose a case...</option>
									{#each cases as caseItem}
										<option value={caseItem.id}>{caseItem.case_name}</option>
									{/each}
								</select>
							</div>
						</div>
						<div>
							<label for="document-notes" class="block text-sm font-medium text-neutral-700">Notes (optional)</label>
							<div class="mt-1">
								<textarea
									id="document-notes"
									bind:value={documentNotes}
									rows="3"
									placeholder="Add any notes about this document..."
									class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
								></textarea>
							</div>
						</div>
					</div>
				</div>
				<div class="bg-neutral-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
					<button
						type="button"
						class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-primary-600 text-base font-medium text-white hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
						disabled={!selectedCaseId || isAddingToCase}
						onclick={addDocumentToCase}
					>
						{#if isAddingToCase}
							<div class="flex items-center">
								<div class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white"></div>
								<span>Adding...</span>
							</div>
						{:else}
							Add to Case
						{/if}
					</button>
					<button
						type="button"
						class="mt-3 w-full inline-flex justify-center rounded-md border border-neutral-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-neutral-700 hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
						onclick={closeAddToCaseModal}
					>
						Cancel
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Create New Case Modal -->
{#if showNewCaseModal}
	<div
		class="fixed inset-0 z-[60] overflow-y-auto"
		in:fade={{ duration: 200 }}
		out:fade={{ duration: 200 }}
	>
		<div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center">
			<!-- Backdrop -->
			<div class="fixed inset-0 bg-neutral-500 bg-opacity-75 transition-opacity" onclick={closeNewCaseModal}></div>

			<!-- Modal Content -->
			<div
				class="relative inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:align-middle sm:max-w-lg sm:w-full z-10"
				in:scale={{ start: 0.95, duration: 200 }}
				out:scale={{ start: 1, end: 0.95, duration: 200 }}
			>
				<div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
					<h3 class="text-lg font-medium text-neutral-900 mb-4">Create New Case</h3>
					<div class="space-y-4">
						<div>
							<label for="new-case-name" class="block text-sm font-medium text-neutral-700">Case Name</label>
							<div class="mt-1">
								<input
									id="new-case-name"
									type="text"
									bind:value={newCaseName}
									placeholder="Enter case name..."
									class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
									onkeydown={(e) => e.key === 'Enter' && createCaseAndAddDocument()}
								/>
							</div>
						</div>
						<div>
							<label for="new-case-notes" class="block text-sm font-medium text-neutral-700">Document Notes (optional)</label>
							<div class="mt-1">
								<textarea
									id="new-case-notes"
									bind:value={documentNotes}
									rows="3"
									placeholder="Add any notes about this document..."
									class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
								></textarea>
							</div>
						</div>
					</div>
				</div>
				<div class="bg-neutral-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
					<button
						type="button"
						class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-primary-600 text-base font-medium text-white hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
						disabled={!newCaseName.trim() || isCreatingCase}
						onclick={createCaseAndAddDocument}
					>
						{#if isCreatingCase}
							<div class="flex items-center">
								<div class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white"></div>
								<span>Creating...</span>
							</div>
						{:else}
							Create Case & Add Document
						{/if}
					</button>
					<button
						type="button"
						class="mt-3 w-full inline-flex justify-center rounded-md border border-neutral-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-neutral-700 hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
						onclick={closeNewCaseModal}
					>
						Cancel
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.opac {
		background-color: rgba(0, 0, 0, 0.5);
	}
</style>
