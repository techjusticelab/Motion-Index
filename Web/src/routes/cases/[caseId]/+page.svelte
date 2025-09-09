<script lang="ts">
	import { page } from '$app/stores';
	import { CaseManager, type Case, type CaseDocument } from '$lib/supabase';
	import { onMount } from 'svelte';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import { cubicOut, quintOut, backOut, elasticOut } from 'svelte/easing';

	let { data } = $props();
	let { session, supabase } = $derived(data);

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

	onMount(() => {
		console.log('Case detail page mounted', { session: !!session?.user, caseId });
		caseManager = new CaseManager(supabase);

		// Set initial load to false after the first render
		setTimeout(() => {
			isInitialLoad = false;
		}, 100);

		// If user is logged in, load case data
		if (session?.user && caseId) {
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
		console.log('loadCaseData called', { session: !!session?.user, caseId });
		if (!session?.user || !caseId) {
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
				caseManager.getUserCases(session.user.id),
				caseManager.getCaseDocuments(caseId)
			]);

			console.log('Data loaded', { userCases: userCases.length, documents: documents.length });

			// Find the current case
			currentCase = userCases.find(c => c.id === caseId) || null;
			caseDocuments = documents;

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

	// Handle keydown for modal escape
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			if (showEditCaseModal) {
				closeEditCaseModal();
			} else if (showAddDocumentModal) {
				closeAddDocumentModal();
			}
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="min-h-screen bg-gray-50 py-8">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		{#if isLoadingCase}
			<div class="flex items-center justify-center py-12">
				<div class="flex items-center space-x-3">
					<div class="h-8 w-8 animate-spin rounded-full border-4 border-indigo-200 border-t-indigo-600"></div>
					<span class="text-gray-600">Loading case...</span>
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
							class="text-sm text-indigo-600 hover:text-indigo-800"
						>
							← Back to Cases
						</a>
					</nav>
					<h1 class="text-3xl font-bold text-gray-900">{currentCase.case_name}</h1>
					<p class="mt-2 text-sm text-gray-600">
						Created {formatDate(currentCase.created_at)} • Updated {formatDate(currentCase.updated_at)}
					</p>
				</div>
				<div class="flex space-x-3">
					<button
						onclick={openAddDocumentModal}
						class="inline-flex items-center justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none"
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
						class="inline-flex items-center justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none"
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
				<h2 class="mb-4 text-xl font-semibold text-gray-900">Documents ({caseDocuments.length})</h2>

				{#if isLoadingDocuments}
					<div class="flex items-center justify-center py-8">
						<div class="flex items-center space-x-3">
							<div class="h-6 w-6 animate-spin rounded-full border-4 border-indigo-200 border-t-indigo-600"></div>
							<span class="text-gray-600">Loading documents...</span>
						</div>
					</div>
				{:else if caseDocuments.length === 0}
					<div
						class="rounded-lg border border-dashed border-gray-300 bg-white p-8 text-center"
						in:fly={{ y: 20, duration: 600, delay: isInitialLoad ? 300 : 0, easing: cubicOut }}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mx-auto h-12 w-12 text-gray-400"
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
						<h3 class="mt-4 text-lg font-medium text-gray-900">No documents yet</h3>
						<p class="mt-2 text-sm text-gray-500">
							Add documents to this case to start building your case file
						</p>
						<button
							onclick={openAddDocumentModal}
							class="mt-4 inline-flex items-center justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none"
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
					<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
						{#each caseDocuments as document, i}
							<div
								class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm hover:border-indigo-300 hover:shadow-md"
								in:fly={{ y: 20, duration: 400, delay: i * 100 }}
							>
								<div class="mb-3 flex items-start justify-between">
									<div class="flex items-center">
										<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-100">
											<svg
												class="h-5 w-5 text-blue-600"
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
										<div class="ml-3">
											<h4 class="text-sm font-medium text-gray-900">
												Document {document.document_ids}
											</h4>
											<p class="text-xs text-gray-500">
												Added {formatDate(document.added_at)}
											</p>
										</div>
									</div>
									<button
										onclick={() => removeDocumentFromCase(document.id)}
										class="text-red-600 hover:text-red-800"
									>
										<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
										</svg>
									</button>
								</div>

								{#if document.notes}
									<div class="mt-3 rounded-md bg-gray-50 p-3">
										<p class="text-sm text-gray-700">{document.notes}</p>
									</div>
								{/if}
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{:else}
			<!-- Case not found -->
			<div class="text-center py-12">
				<h2 class="text-2xl font-bold text-gray-900 mb-4">Case not found</h2>
				<p class="text-gray-600 mb-6">The case you're looking for doesn't exist or you don't have access to it.</p>
				<a 
					href="/cases" 
					class="inline-flex items-center justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700"
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
				class="fixed inset-0 bg-gray-900 bg-opacity-50 transition-opacity"
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
					<h3 class="text-lg font-medium text-gray-900 mb-4">Edit Case</h3>
					<div>
						<label for="edit-case-name" class="block text-sm font-medium text-gray-700">Case Name</label>
						<div class="mt-1">
							<input
								id="edit-case-name"
								type="text"
								bind:value={editCaseName}
								placeholder="Enter case name..."
								class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
								onkeydown={(e) => e.key === 'Enter' && updateCaseName()}
							/>
						</div>
					</div>
				</div>
				<div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
					<button
						type="button"
						class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
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
						class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
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
				class="fixed inset-0 bg-gray-900 bg-opacity-50 transition-opacity"
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
						<h3 class="text-lg font-medium text-gray-900">Add Document to Case</h3>
						<button 
							onclick={closeAddDocumentModal}
							class="rounded-full p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-500"
							aria-label="Close modal"
						>
							<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
							</svg>
						</button>
					</div>
					<div class="space-y-4">
						<div>
							<label for="document-id" class="block text-sm font-medium text-gray-700">Document ID</label>
							<div class="mt-1">
								<input
									id="document-id"
									type="text"
									bind:value={selectedDocumentId}
									placeholder="Enter document ID..."
									class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
								/>
							</div>
						</div>
						<div>
							<label for="document-notes" class="block text-sm font-medium text-gray-700">Notes (optional)</label>
							<div class="mt-1">
								<textarea
									id="document-notes"
									bind:value={documentNotes}
									rows="3"
									placeholder="Add any notes about this document..."
									class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
								></textarea>
							</div>
						</div>
					</div>
				</div>
				<div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
					<button
						type="button"
						class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
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
						class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
						onclick={closeAddDocumentModal}
					>
						Cancel
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}