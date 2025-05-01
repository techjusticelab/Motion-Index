<!-- DocumentViewer.svelte -->
<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { formatDate } from '../../../utils/utils';
	import type { Document, SearchResponse } from '../../../utils/search_types';
	import S3FileViewer from './s3viewer.svelte'; // Update this path
	import { getSignedS3Url } from '../../services/s3';

	export let docData: Document | null = null;
	export let isOpen: boolean = false;

	let isLoading = true;
	let errorMessage = '';
	let url = '';

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
</script>

{#if isOpen && docData}
	<div class="opac opac fixed inset-0 z-50 flex items-center justify-center p-4 shadow">
		<div
			class="relative flex h-[90vh] w-[90vw] max-w-6xl overflow-hidden rounded-xl bg-white shadow-2xl"
		>
			<!-- Sidebar with metadata -->
			<div class="w-80 overflow-auto border-r border-gray-200 bg-gray-50 p-6">
				<div class="mb-6">
					<h2 class="truncate text-xl font-semibold text-gray-800">
						{docData.metadata.document_name || docData.file_name}
					</h2>
					<p class="mt-1 text-sm text-gray-500">{docData.file_name}</p>
				</div>

				<!-- Document type tag -->
				<div class="mb-4">
					<span
						class="inline-flex rounded-full bg-blue-50 px-3 py-1 text-sm font-medium text-blue-700"
					>
						{docData.doc_type}
					</span>
				</div>

				<!-- Metadata list -->
				<div class="space-y-4">
					<h3 class="text-sm font-medium text-gray-700">Document Details</h3>

					<div class="space-y-2 rounded-lg border border-gray-200 bg-white p-3">
						<!-- Document Name - always displayed -->
						<div>
							<p class="text-xs text-gray-500">Document Name</p>
							<p class="text-sm font-medium text-gray-800">{docData.metadata.document_name}</p>
						</div>

						{#if docData.metadata.case_number}
							<div>
								<p class="text-xs text-gray-500">Case Number</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.case_number}</p>
							</div>
						{/if}

						{#if docData.metadata.case_name}
							<div>
								<p class="text-xs text-gray-500">Case Name</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.case_name}</p>
							</div>
						{/if}

						{#if docData.metadata.court}
							<div>
								<p class="text-xs text-gray-500">Court</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.court}</p>
							</div>
						{/if}

						{#if docData.metadata.judge}
							<div>
								<p class="text-xs text-gray-500">Judge</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.judge}</p>
							</div>
						{/if}

						{#if docData.metadata.timestamp}
							<div>
								<p class="text-xs text-gray-500">Date</p>
								<p class="text-sm font-medium text-gray-800">
									{formatDate(docData.metadata.timestamp)}
								</p>
							</div>
						{:else}
							<div>
								<p class="text-xs text-gray-500">Created Date</p>
								<p class="text-sm font-medium text-gray-800">
									{formatDate(docData.created_at)}
								</p>
							</div>
						{/if}

						{#if docData.metadata.subject}
							<div>
								<p class="text-xs text-gray-500">Subject</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.subject}</p>
							</div>
						{/if}

						{#if docData.metadata.status}
							<div>
								<p class="text-xs text-gray-500">Status</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.status}</p>
							</div>
						{/if}

						{#if docData.metadata.author}
							<div>
								<p class="text-xs text-gray-500">Author</p>
								<p class="text-sm font-medium text-gray-800">{docData.metadata.author}</p>
							</div>
						{/if}

						{#if docData.metadata.legal_tags && docData.metadata.legal_tags.length > 0}
							<div>
								<p class="text-xs text-gray-500">Legal Tags</p>
								<div class="mt-1 flex flex-wrap gap-1">
									{#each docData.metadata.legal_tags as tag}
										<span
											class="inline-flex rounded-full bg-gray-100 px-2 py-0.5 text-xs font-medium text-gray-800"
										>
											{tag}
										</span>
									{/each}
								</div>
							</div>
						{/if}
					</div>

					<!-- Download button -->
					<a
						href={url}
						download={docData.file_name}
						class="mt-6 flex w-full items-center justify-center rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
						target="_blank"
						rel="noopener noreferrer"
						on:click={downloadFile}
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
						Download Original
					</a>
				</div>
			</div>

			<!-- Document Viewer area -->
			<div class="relative flex-1 overflow-hidden">
				<!-- Header with close button -->
				<div class="absolute right-0 top-0 z-10 flex items-center justify-end p-4">
					<button
						type="button"
						class="rounded-full bg-white/90 p-2 shadow-md hover:bg-gray-100"
						on:click={closeViewer}
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
										d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
									/>
								</svg>
								<p class="text-gray-700">
									This file type cannot be previewed directly in the browser.
								</p>
								<a
									href={docData.s3_uri}
									download={docData.file_name}
									class="mt-4 inline-flex items-center rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
									target="_blank"
									rel="noopener noreferrer"
								>
									Download File
								</a>
							</div>
						</div>
					{/if}

					{#if errorMessage}
						<div class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-90">
							<div class="max-w-md rounded-lg border border-red-100 bg-red-50 p-6 text-center">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="mx-auto mb-4 h-12 w-12 text-red-400"
									fill="none"
									viewBox="0 0 24 24"
									stroke="currentColor"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
									/>
								</svg>
								<p class="text-gray-700">{errorMessage}</p>
								<a
									href={docData.s3_uri}
									download={docData.file_name}
									class="mt-4 inline-flex items-center rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
									target="_blank"
									rel="noopener noreferrer"
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

<style>
	.opac {
		background-color: rgba(0, 0, 0, 0.5);
	}
</style>
