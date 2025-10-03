<!-- CaseActions.svelte -->
<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import type { Document } from '$lib/types';
	import { CaseManager, type Case } from '$lib/supabase';
	import { fade, slide, scale } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';

	const {
		docData,
		supabase = null,
		session = null,
		user = null,
		url = ''
	}: {
		docData?: Document | null;
		supabase?: any;
		session?: any;
		user?: any;
		url?: string;
	} = $props();

	// Case management state
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

	onMount(async () => {
		if (supabase && user) {
			caseManager = new CaseManager(supabase);

			// Test database access to debug the issue
			await caseManager.testDatabaseAccess();

			loadUserCases();
		}
	});

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

{#if session && supabase}
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
	<div class="m-auto flex flex-col md:flex-row justify-center gap-2 align-middle">
		<button
			onclick={openAddToCaseModal}
			class="mt-4 flex w-full items-center justify-center rounded-lg border-primary-600 bg-white px-2 md:px-3 py-2 text-xs font-medium text-black hover:bg-primary-600 hover:text-white transition-colors"
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
			<span class="hidden md:inline">Add to Case</span>
			<span class="md:hidden">Add Case</span>
		</button>
		<button
			onclick={openNewCaseModal}
			class="mt-4 flex w-full items-center justify-center rounded-lg bg-primary-600 px-2 md:px-3 py-2 text-xs font-medium text-white hover:bg-primary-700 transition-colors"
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
			download={docData?.file_name}
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
			Download
		</a>
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