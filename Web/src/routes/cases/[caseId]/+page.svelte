<script lang="ts">
	import { page } from '$app/stores';
	import { CaseManager, type Case, type CaseDocument } from '$lib/supabase';
	import { onMount } from 'svelte';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import { cubicOut, quintOut, backOut, elasticOut } from 'svelte/easing';
	import DocumentViewer from '$lib/components/document/DocumentViewer.svelte';
	import type { Document } from '$lib/api';
	import * as api from '$lib/api';

	let { data } = $props();
	let { session, user, supabase } = $derived(data);

	// Get case ID from URL
	let caseId = $derived($page.params.caseId);

	// Flag to control animations after initial page load
	let isInitialLoad = $state(true);

	// Case data
	let currentCase = $state<Case | null>(null);
	let caseDocuments = $state<CaseDocument[]>([]);
	let caseManager: CaseManager;
	let isLoadingCase = $state(true);
	let isLoadingDocuments = $state(true);

	// Edit case modal
	let showEditCaseModal = $state(false);
	let editCaseName = $state('');
	let isUpdatingCase = $state(false);

	// Add document modal
	let showAddDocumentModal = $state(false);
	let selectedDocumentId = $state('');
	let documentNotes = $state('');
	let isAddingDocument = $state(false);

	// Document viewer
	let activeDocument = $state<Document | null>(null);
	let showDocumentViewer = $state(false);
	let documentDetails = $state<Map<string, Document>>(new Map());

	onMount(() => {
		console.log('Case detail page mounted', { user: !!user, caseId });
		caseManager = new CaseManager(supabase);

		// Set initial load to false after the first render
		setTimeout(() => {
			isInitialLoad = false;
		}, 100);

		// If user is logged in, load case data
		if (user && caseId) {
			console.log('Loading case data for:', caseId);
			loadCaseData();
		} else {
			console.log('No session or caseId, stopping loading');
			// If no session or caseId, stop loading
			isLoadingCase = false;
			isLoadingDocuments = false;
		}
	});

	// Load case data
	async function loadCaseData() {
		console.log('loadCaseData called', { user: !!user, caseId });
		if (!user || !caseId) {
			console.log('No session or caseId in loadCaseData');
			isLoadingCase = false;
			isLoadingDocuments = false;
			return;
		}
		
		console.log('Setting loading states to true');
		isLoadingCase = true;
		isLoadingDocuments = true;
		
		try {
			console.log('Starting to load case data...');
			// Load case details and documents in parallel
			const [userCases, documents] = await Promise.all([
				caseManager.getUserCases(user.id),
				caseManager.getCaseDocuments(caseId)
			]);

			console.log('Data loaded', { userCases: userCases.length, documents: documents.length });

			// Find the current case
			currentCase = userCases.find(c => c.id === caseId) || null;
			caseDocuments = documents;

			// Preload document details for better UX (run in background)
			if (documents.length > 0) {
				preloadDocumentDetails(documents).catch(error =>
					console.error('Error preloading document details:', error)
				);
			}

			// If case not found, set current case to null but don't redirect
			// Let the template handle showing the "not found" state
			if (!currentCase) {
				console.log('Case not found:', caseId);
			} else {
				console.log('Case found:', currentCase.case_name);
			}

		} catch (error) {
			console.error('Error loading case data:', error);
			// Set currentCase to null so the "not found" state is shown
			currentCase = null;
			caseDocuments = [];
		} finally {
			console.log('Setting loading states to false');
			isLoadingCase = false;
			isLoadingDocuments = false;
		}
	}

	// Open edit case modal
	function openEditCaseModal() {
		if (!currentCase) return;
		editCaseName = currentCase.case_name;
		showEditCaseModal = true;
	}

	// Close edit case modal
	function closeEditCaseModal() {
		showEditCaseModal = false;
		editCaseName = '';
	}

	// Update case name
	async function updateCaseName() {
		if (!currentCase || !editCaseName.trim()) return;
		
		isUpdatingCase = true;
		try {
			const success = await caseManager.updateCaseName(currentCase.id, editCaseName.trim());
			if (success && currentCase) {
				currentCase.case_name = editCaseName.trim();
				closeEditCaseModal();
			}
		} catch (error) {
			console.error('Error updating case:', error);
		} finally {
			isUpdatingCase = false;
		}
	}

	// Open add document modal
	function openAddDocumentModal() {
		selectedDocumentId = '';
		documentNotes = '';
		showAddDocumentModal = true;
	}

	// Close add document modal
	function closeAddDocumentModal() {
		showAddDocumentModal = false;
		selectedDocumentId = '';
		documentNotes = '';
	}

	// Add document to case
	async function addDocumentToCase() {
		if (!caseId || !selectedDocumentId.trim()) return;
		
		isAddingDocument = true;
		try {
			const caseDocument = await caseManager.addDocumentToCase(
				caseId, 
				selectedDocumentId.trim(), 
				documentNotes.trim() || undefined
			);
			if (caseDocument) {
				caseDocuments = [caseDocument, ...caseDocuments];
				closeAddDocumentModal();
			}
		} catch (error) {
			console.error('Error adding document to case:', error);
		} finally {
			isAddingDocument = false;
		}
	}

	// Remove document from case
	async function removeDocumentFromCase(caseDocumentId: string) {
		if (!confirm('Are you sure you want to remove this document from the case?')) return;
		
		const success = await caseManager.removeDocumentFromCase(caseDocumentId);
		if (success) {
			caseDocuments = caseDocuments.filter(d => d.id !== caseDocumentId);
		}
	}

	// Delete entire case
	async function deleteCase() {
		if (!currentCase) return;
		if (!confirm('Are you sure you want to delete this case and all its documents?')) return;
		
		const success = await caseManager.deleteCase(currentCase.id);
		if (success) {
			goto('/cases');
		}
	}

	// Format date for display
	function formatDate(dateString: string) {
		if (!dateString) return 'N/A';
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	// Fetch document details by ID and open viewer
	async function openDocumentViewer(documentId: string) {
		console.log('Opening document viewer for ID:', documentId);

		// Check if we already have the document details
		const cachedDocument = documentDetails.get(documentId);
		if (cachedDocument) {
			activeDocument = cachedDocument;
			showDocumentViewer = true;
			return;
		}

		// If we don't have cached details, we need to search for the document
		// This is a limitation since we only have the document ID, not full details
		try {
			// Try to search for documents that might match this ID
			const searchResults = await api.searchDocuments({
				query: documentId,
				size: 10
			}, session);

			// Look for an exact match or the best match
			const matchedDocument = searchResults.hits.find(doc =>
				doc.s3_uri?.includes(documentId) ||
				doc.file_name?.includes(documentId) ||
				JSON.stringify(doc).includes(documentId)
			);

			if (matchedDocument) {
				documentDetails.set(documentId, matchedDocument);
				activeDocument = matchedDocument;
				showDocumentViewer = true;
			} else {
				console.error('Document not found in search results');
				// Create a placeholder document with limited info
				const placeholderDoc: Document = {
					metadata: {
						document_name: `Document ${documentId}`,
						case_name: currentCase?.case_name
					},
					file_name: `Document_${documentId}`,
					s3_uri: documentId,
					doc_type: 'Unknown',
					text: '',
					created_at: new Date().toISOString()
				};
				activeDocument = placeholderDoc;
				showDocumentViewer = true;
			}
		} catch (error) {
			console.error('Error fetching document details:', error);
		}
	}

	// Close document viewer
	function closeDocumentViewer() {
		activeDocument = null;
		showDocumentViewer = false;
	}

	// Preload document details for all case documents
	async function preloadDocumentDetails(caseDocuments: CaseDocument[]) {
		console.log('Preloading document details for', caseDocuments.length, 'documents');

		for (const caseDoc of caseDocuments) {
			try {
				// Try to search for the document to get full details
				const searchResults = await api.searchDocuments({
					query: caseDoc.document_ids,
					size: 5
				}, session);

				// Look for a matching document
				const matchedDocument = searchResults.hits.find(doc =>
					doc.id === caseDoc.document_ids ||
					doc.s3_uri?.includes(caseDoc.document_ids) ||
					doc.file_name?.includes(caseDoc.document_ids)
				);

				if (matchedDocument) {
					documentDetails.set(caseDoc.document_ids, matchedDocument);
					console.log('Preloaded details for document:', caseDoc.document_ids);
				} else {
					console.log('Could not find details for document:', caseDoc.document_ids);
				}
			} catch (error) {
				console.error('Error preloading document details for', caseDoc.document_ids, ':', error);
			}
		}
	}

	// Handle keydown for modal escape
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			if (showEditCaseModal) {
				closeEditCaseModal();
			} else if (showAddDocumentModal) {
				closeAddDocumentModal();
			} else if (showDocumentViewer) {
				closeDocumentViewer();
			}
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="min-h-screen bg-neutral-50 py-8">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		{#if isLoadingCase}
			<div class="flex items-center justify-center py-12">
				<div class="flex items-center space-x-3">
					<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-200 border-t-indigo-600"></div>
					<span class="text-neutral-600">Loading case...</span>
				</div>
			</div>
		{:else if currentCase}
			<!-- Header -->
			<div 
				class="mb-8 flex items-start justify-between"
				in:fly={{ y: -30, duration: 800, easing: quintOut, delay: isInitialLoad ? 0 : 0 }}
			>
				<div>
					<nav class="mb-4">
						<a 
							href="/cases" 
							class="text-sm text-primary-600 hover:text-primary-800"
						>
							← Back to Cases
						</a>
					</nav>
					<h1 class="text-3xl font-bold text-neutral-900">{currentCase.case_name}</h1>
					<p class="mt-2 text-sm text-neutral-600">
						Created {formatDate(currentCase.created_at)} • Updated {formatDate(currentCase.updated_at)}
					</p>
				</div>
				<div class="flex space-x-3">
					<button
						onclick={openAddDocumentModal}
						class="inline-flex items-center justify-center rounded-md bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 focus:outline-none"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mr-2 h-4 w-4"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
						</svg>
						Add Document
					</button>
					<button
						onclick={openEditCaseModal}
						class="inline-flex items-center justify-center rounded-md border border-neutral-300 bg-white px-4 py-2 text-sm font-medium text-neutral-700 hover:bg-neutral-50 focus:outline-none"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mr-2 h-4 w-4"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
						</svg>
						Edit Case
					</button>
					<button
						onclick={deleteCase}
						class="inline-flex items-center justify-center rounded-md border border-red-200 bg-white px-4 py-2 text-sm font-medium text-red-600 hover:bg-red-50 focus:outline-none"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mr-2 h-4 w-4"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
						</svg>
						Delete Case
					</button>
				</div>
			</div>

			<!-- Case Documents -->
			<div class="mb-8">
				<h2 class="mb-4 text-xl font-semibold text-neutral-900">Documents ({caseDocuments.length})</h2>

				{#if isLoadingDocuments}
					<div class="flex items-center justify-center py-8">
						<div class="flex items-center space-x-3">
							<div class="h-6 w-6 animate-spin rounded-full border-4 border-primary-200 border-t-indigo-600"></div>
							<span class="text-neutral-600">Loading documents...</span>
						</div>
					</div>
				{:else if caseDocuments.length === 0}
					<div
						class="rounded-lg border border-dashed border-neutral-300 bg-white p-8 text-center"
						in:fly={{ y: 20, duration: 600, delay: isInitialLoad ? 300 : 0, easing: cubicOut }}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mx-auto h-12 w-12 text-neutral-400"
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
						<h3 class="mt-4 text-lg font-medium text-neutral-900">No documents yet</h3>
						<p class="mt-2 text-sm text-neutral-500">
							Add documents to this case to start building your case file
						</p>
						<button
							onclick={openAddDocumentModal}
							class="mt-4 inline-flex items-center justify-center rounded-md bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 focus:outline-none"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="mr-2 h-4 w-4"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="2"
							>
								<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
							</svg>
							Add Your First Document
						</button>
					</div>
				{:else}
					<div class="grid gap-6 md:grid-cols-1 lg:grid-cols-2">
						{#each caseDocuments as document, i}
							<div
								class="group rounded-xl border border-neutral-200 bg-white p-6 shadow-sm hover:border-primary-300 hover:shadow-lg transition-all duration-300 cursor-pointer"
								in:fly={{ y: 20, duration: 400, delay: i * 100 }}
								onclick={() => openDocumentViewer(document.document_ids)}
							>
								<!-- Document Header -->
								<div class="mb-4 flex items-start justify-between">
									<div class="flex items-start space-x-3">
										<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-primary-500 to-primary-600 shadow-md">
											<svg
												class="h-6 w-6 text-white"
												fill="none"
												stroke="currentColor"
												viewBox="0 0 24 24"
											>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
												/>
											</svg>
										</div>
										<div class="flex-1 min-w-0">
											<h4 class="text-lg font-semibold text-neutral-900 group-hover:text-primary-600 transition-colors truncate">
												{documentDetails.get(document.document_ids)?.metadata?.document_name ||
												 documentDetails.get(document.document_ids)?.file_name ||
												 `Document ${document.document_ids.slice(0, 8)}...`}
											</h4>
											<div class="mt-1 flex items-center space-x-4 text-sm text-neutral-500">
												<div class="flex items-center">
													<svg class="mr-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
													</svg>
													Added {formatDate(document.added_at)}
												</div>
											</div>
										</div>
									</div>
									<div class="flex items-center space-x-2">
										<button
											onclick={(e) => { e.stopPropagation(); openDocumentViewer(document.document_ids); }}
											class="inline-flex items-center rounded-lg bg-primary-50 px-3 py-2 text-sm font-medium text-primary-700 hover:bg-primary-100 transition-colors"
											title="View Document"
											aria-label="View Document"
										>
											<svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
											</svg>
											View
										</button>
										<button
											onclick={(e) => { e.stopPropagation(); removeDocumentFromCase(document.id); }}
											class="inline-flex items-center rounded-lg bg-red-50 px-3 py-2 text-sm font-medium text-red-700 hover:bg-red-100 transition-colors"
											title="Remove from Case"
											aria-label="Remove from Case"
										>
											<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
											</svg>
										</button>
									</div>
								</div>

								<!-- Document Details -->
								<div class="space-y-3">
									<!-- Document ID -->
									<div class="flex items-center justify-between">
										<span class="text-sm font-medium text-neutral-600">Document ID</span>
										<div class="flex items-center space-x-2">
											<code class="rounded bg-neutral-100 px-2 py-1 text-xs font-mono text-neutral-800">
												{document.document_ids.slice(0, 16)}...
											</code>
											<button
												onclick={(e) => { e.stopPropagation(); navigator.clipboard.writeText(document.document_ids); }}
												class="text-neutral-400 hover:text-neutral-600"
												title="Copy Document ID"
												aria-label="Copy Document ID"
											>
												<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
												</svg>
											</button>
										</div>
									</div>

									<!-- Document Type Badge -->
									{#if documentDetails.get(document.document_ids)?.doc_type}
										<div class="flex items-center justify-between">
											<span class="text-sm font-medium text-neutral-600">Type</span>
											<span class="inline-flex rounded-full bg-blue-50 px-3 py-1 text-xs font-medium text-blue-700">
												{documentDetails.get(document.document_ids)?.doc_type}
											</span>
										</div>
									{/if}

									<!-- Notes Section -->
									{#if document.notes}
										<div class="mt-4 rounded-lg bg-amber-50 border border-amber-200 p-4">
											<div class="flex items-start">
												<svg class="mr-2 mt-0.5 h-4 w-4 text-amber-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
												</svg>
												<div>
													<h5 class="text-sm font-medium text-amber-800">Notes</h5>
													<p class="mt-1 text-sm text-amber-700">{document.notes}</p>
												</div>
											</div>
										</div>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{:else}
			<!-- Case not found -->
			<div class="text-center py-12">
				<h2 class="text-2xl font-bold text-neutral-900 mb-4">Case not found</h2>
				<p class="text-neutral-600 mb-6">The case you're looking for doesn't exist or you don't have access to it.</p>
				<a 
					href="/cases" 
					class="inline-flex items-center justify-center rounded-md bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700"
				>
					Back to Cases
				</a>
			</div>
		{/if}
	</div>
</div>

<!-- Edit Case Modal -->
{#if showEditCaseModal}
	<div 
		class="fixed inset-0 z-[60] overflow-y-auto"
		in:fade={{ duration: 200 }}
		out:fade={{ duration: 200 }}
	>
		<div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:p-0">
			<div 
				class="fixed inset-0 bg-neutral-900 bg-opacity-50 transition-opacity"
				onclick={closeEditCaseModal}
				aria-hidden="true"
			></div>

			<div 
				class="relative inline-block bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all w-full max-w-lg mx-auto"
				in:scale={{ start: 0.95, duration: 200 }}
				out:scale={{ start: 1, end: 0.95, duration: 200 }}
				onclick={(e) => e.stopPropagation()}
			>
				<div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
					<h3 class="text-lg font-medium text-neutral-900 mb-4">Edit Case</h3>
					<div>
						<label for="edit-case-name" class="block text-sm font-medium text-neutral-700">Case Name</label>
						<div class="mt-1">
							<input
								id="edit-case-name"
								type="text"
								bind:value={editCaseName}
								placeholder="Enter case name..."
								class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
								onkeydown={(e) => e.key === 'Enter' && updateCaseName()}
							/>
						</div>
					</div>
				</div>
				<div class="bg-neutral-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
					<button
						type="button"
						class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-primary-600 text-base font-medium text-white hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
						disabled={!editCaseName.trim() || isUpdatingCase}
						onclick={updateCaseName}
					>
						{#if isUpdatingCase}
							<div class="flex items-center">
								<div class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white"></div>
								<span>Updating...</span>
							</div>
						{:else}
							Update Case
						{/if}
					</button>
					<button
						type="button"
						class="mt-3 w-full inline-flex justify-center rounded-md border border-neutral-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-neutral-700 hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
						onclick={closeEditCaseModal}
					>
						Cancel
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Add Document Modal -->
{#if showAddDocumentModal}
	<div 
		class="fixed inset-0 z-[60] overflow-y-auto"
		in:fade={{ duration: 200 }}
		out:fade={{ duration: 200 }}
	>
		<div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:p-0">
			<div 
				class="fixed inset-0 bg-neutral-900 bg-opacity-50 transition-opacity"
				onclick={closeAddDocumentModal}
				aria-hidden="true"
			></div>

			<div 
				class="relative inline-block bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all w-full max-w-lg mx-auto"
				in:scale={{ start: 0.95, duration: 200 }}
				out:scale={{ start: 1, end: 0.95, duration: 200 }}
				onclick={(e) => e.stopPropagation()}
			>
				<div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
					<div class="mb-4 flex items-center justify-between">
						<h3 class="text-lg font-medium text-neutral-900">Add Document to Case</h3>
						<button 
							onclick={closeAddDocumentModal}
							class="rounded-full p-1 text-neutral-400 hover:bg-neutral-100 hover:text-neutral-500"
							aria-label="Close modal"
						>
							<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
							</svg>
						</button>
					</div>
					<div class="space-y-4">
						<div>
							<label for="document-id" class="block text-sm font-medium text-neutral-700">Document ID</label>
							<div class="mt-1">
								<input
									id="document-id"
									type="text"
									bind:value={selectedDocumentId}
									placeholder="Enter document ID..."
									class="block w-full appearance-none rounded-md border border-neutral-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
								/>
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
						disabled={!selectedDocumentId.trim() || isAddingDocument}
						onclick={addDocumentToCase}
					>
						{#if isAddingDocument}
							<div class="flex items-center">
								<div class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white"></div>
								<span>Adding...</span>
							</div>
						{:else}
							Add Document
						{/if}
					</button>
					<button
						type="button"
						class="mt-3 w-full inline-flex justify-center rounded-md border border-neutral-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-neutral-700 hover:bg-neutral-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
						onclick={closeAddDocumentModal}
					>
						Cancel
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Document Viewer -->
<DocumentViewer
	docData={activeDocument}
	isOpen={showDocumentViewer}
	{supabase}
	{session}
	{user}
	onclose={closeDocumentViewer}
/>