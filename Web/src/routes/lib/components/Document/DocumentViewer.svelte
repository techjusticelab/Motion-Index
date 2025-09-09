<!-- DocumentViewer.svelte -->
<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { formatDate } from '../../../utils/utils';
	import type { Document, SearchResponse } from '../../../utils/search_types';
	import S3FileViewer from './s3viewer.svelte'; // Update this path
	import { getSignedS3Url } from '../../services/s3';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { elasticOut, backOut, quintOut, cubicOut } from 'svelte/easing';
	import { CaseManager, type Case } from '$lib/supabase';
	import { page } from '$app/stores';

	export let docData: Document | null = null;
	export let isOpen: boolean = false;
	export let supabase: any = null;
	export let session: any = null;

	let isLoading = true;
	let errorMessage = '';
	let url = '';
	
	// Case management
	let cases: Case[] = [];
	let caseManager: CaseManager | null = null;
	let showAddToCaseModal = false;
	let selectedCaseId = '';
	let documentNotes = '';
	let isAddingToCase = false;
	let showNewCaseModal = false;
	let newCaseName = '';
	let isCreatingCase = false;

	const dispatch = createEventDispatcher<{
		close: void;
	}>();

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
	onMount(() => {
		if (docData) {
			downloadFile();
		}
		if (supabase && session?.user) {
			caseManager = new CaseManager(supabase);
			loadUserCases();
		}
	});
	async function downloadFile() {
		if (docData) {
			url = await getSignedS3Url(docData.s3_uri);
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
		if (!caseManager || !session?.user) return;
		try {
			cases = await caseManager.getUserCases(session.user.id);
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
		if (!caseManager || !docData || !selectedCaseId) return;
		
		isAddingToCase = true;
		try {
			await caseManager.addDocumentToCase(
				selectedCaseId, 
				docData.id, 
				documentNotes.trim() || undefined
			);
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
		if (!caseManager || !session?.user || !docData || !newCaseName.trim()) return;
		
		isCreatingCase = true;
		try {
			const newCase = await caseManager.createCase(session.user.id, newCaseName.trim());
			if (newCase) {
				cases = [newCase, ...cases];
				await caseManager.addDocumentToCase(newCase.id, docData.id, documentNotes.trim() || undefined);
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
				class="w-1/4 overflow-auto border-r border-gray-200 bg-gray-50 p-6"
				in:fly={{ x: -20, duration: 800, delay: 300, easing: cubicOut }}
			>
				<div class="mb-4">
					<h2
						class="truncate text-xl font-semibold text-gray-800"
						in:slide={{ duration: 700, delay: 400 }}
					>
						{docData.metadata.document_name || docData.file_name}
					</h2>
					<p class="mt-1 text-sm text-gray-500" in:slide={{ duration: 700, delay: 500 }}>
						{docData.file_name}
					</p>
				</div>

				<!-- Document type tag -->
				<div class="mb-4">
					<span
						class="inline-flex rounded-full bg-blue-50 px-3 py-1 text-sm font-medium text-blue-700"
						in:scale={{ start: 0.9, duration: 600, delay: 600, easing: cubicOut }}
					>
						{docData.doc_type}
					</span>
				</div>

				<!-- Metadata list -->
				<div class="space-y-4 overflow-hidden" in:slide={{ duration: 600, delay: 650 }}>
					<h3 class="text-sm font-medium text-gray-700">Document Details</h3>

					<div class="space-y-2 rounded-lg border border-gray-200 bg-white p-3">
						<!-- Document Name - always displayed -->
						<div in:slide={{ duration: 500, delay: 700 }}>
							<p class="text-xs text-gray-500">Document Name</p>
							<p class="text-sm font-medium text-gray-800">{docData.metadata.document_name}</p>
						</div>

						{#if docData.metadata.case_number}
							<div in:slide={{ duration: 500, delay: 750 }}>
								<p class="text-xs text-gray-500">Case Number</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.case_number}</p>
							</div>
						{/if}

						{#if docData.metadata.case_name}
							<div in:slide={{ duration: 500, delay: 800 }}>
								<p class="text-xs text-gray-500">Case Name</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.case_name}</p>
							</div>
						{/if}

						{#if docData.metadata.court}
							<div in:slide={{ duration: 500, delay: 850 }}>
								<p class="text-xs text-gray-500">Court</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.court}</p>
							</div>
						{/if}

						{#if docData.metadata.judge}
							<div in:slide={{ duration: 500, delay: 900 }}>
								<p class="text-xs text-gray-500">Judge</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.judge}</p>
							</div>
						{/if}

						{#if docData.metadata.timestamp}
							<div in:slide={{ duration: 500, delay: 950 }}>
								<p class="text-xs text-gray-500">Date</p>
								<p class="text-sm font-medium text-gray-800">
									{formatDate(docData.metadata.timestamp)}
								</p>
							</div>
						{:else}
							<div in:slide={{ duration: 500, delay: 950 }}>
								<p class="text-xs text-gray-500">Created Date</p>
								<p class="text-sm font-medium text-gray-800">
									{formatDate(docData.created_at)}
								</p>
							</div>
						{/if}

						{#if docData.metadata.subject}
							<div in:slide={{ duration: 500, delay: 1000 }}>
								<p class="text-xs text-gray-500">Subject</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.subject}</p>
							</div>
						{/if}

						{#if docData.metadata.status}
							<div in:slide={{ duration: 500, delay: 1050 }}>
								<p class="text-xs text-gray-500">Status</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.status}</p>
							</div>
						{/if}

						{#if docData.metadata.author}
							<div in:slide={{ duration: 500, delay: 1100 }}>
								<p class="text-xs text-gray-500">Author</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.author}</p>
							</div>
						{/if}

						{#if docData.metadata.legal_tags && docData.metadata.legal_tags.length > 0}
							<div in:slide={{ duration: 500, delay: 1150 }}>
								<p class="text-xs text-gray-500">Legal Tags</p>
								<div class="mt-1 flex flex-wrap gap-1">
									{#each docData.metadata.legal_tags as tag, i}
										<span
											class="inline-flex rounded-full bg-gray-100 px-2 py-0.5 text-xs font-medium text-gray-800"
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

					<!-- Case Actions -->
					{#if session && supabase}
						<div class="m-auto flex flex-row justify-center gap-2 align-middle">
							<button
								onclick={openAddToCaseModal}
								class="mt-4 flex w-full items-center justify-center rounded-lg border-indigo-600 bg-white px-3 py-2 text-xs font-medium text-black hover:bg-indigo-600 hover:text-white transition-colors"
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
								class="mt-4 flex w-full items-center justify-center rounded-lg bg-indigo-600 px-3 py-2 text-xs font-medium text-white hover:bg-indigo-700 transition-colors"
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
								class="flex items-center text-xs text-indigo-600 hover:text-indigo-800 underline"
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
						class="rounded-full bg-white/90 p-2 shadow-md hover:bg-gray-100"
						onclick={closeViewer}
						in:scale={{ start: 0.9, duration: 600, delay: 900 }}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6 text-gray-500"
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
				<div class="h-full w-full overflow-auto bg-gray-100">
					{#if canViewInIframe(docData.file_name)}
						<S3FileViewer
							s3Uri={docData.s3_uri}
							height="100%"
							width="100%"
							on:loaded={handleViewerLoaded}
							on:error={handleViewerError}
						/>
					{:else}
						<div class="flex h-full w-full items-center justify-center p-6">
							<div
								class="max-w-md rounded-lg border border-blue-100 bg-blue-50 p-6 text-center"
								in:scale={{ start: 0.95, duration: 700, delay: 700, easing: cubicOut }}
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="mx-auto mb-4 h-12 w-12 text-blue-400"
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
								<p class="text-gray-700" in:slide={{ duration: 600, delay: 900 }}>
									This file type cannot be previewed directly in the browser.
								</p>
								<a
									href={url}
									download={docData.file_name}
									class="mt-4 inline-flex items-center rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
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
								<p class="text-gray-700" in:slide={{ duration: 600, delay: 300 }}>
									{errorMessage}
								</p>
								<a
									href={url}
									download={docData.file_name}
									class="mt-4 inline-flex items-center rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
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
		<div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
			<div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>

			<div 
				class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full"
				in:scale={{ start: 0.95, duration: 200 }}
				out:scale={{ start: 1, end: 0.95, duration: 200 }}
			>
				<div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
					<h3 class="text-lg font-medium text-gray-900 mb-4">Add Document to Case</h3>
					<div class="space-y-4">
						<div>
							<label for="case-select" class="block text-sm font-medium text-gray-700">Select Case</label>
							<div class="mt-1">
								<select
									id="case-select"
									bind:value={selectedCaseId}
									class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
								>
									<option value="">Choose a case...</option>
									{#each cases as caseItem}
										<option value={caseItem.id}>{caseItem.case_name}</option>
									{/each}
								</select>
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
						class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
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
		<div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
			<div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>

			<div 
				class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full"
				in:scale={{ start: 0.95, duration: 200 }}
				out:scale={{ start: 1, end: 0.95, duration: 200 }}
			>
				<div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
					<h3 class="text-lg font-medium text-gray-900 mb-4">Create New Case</h3>
					<div class="space-y-4">
						<div>
							<label for="new-case-name" class="block text-sm font-medium text-gray-700">Case Name</label>
							<div class="mt-1">
								<input
									id="new-case-name"
									type="text"
									bind:value={newCaseName}
									placeholder="Enter case name..."
									class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
									onkeydown={(e) => e.key === 'Enter' && createCaseAndAddDocument()}
								/>
							</div>
						</div>
						<div>
							<label for="new-case-notes" class="block text-sm font-medium text-gray-700">Document Notes (optional)</label>
							<div class="mt-1">
								<textarea
									id="new-case-notes"
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
						class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
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
