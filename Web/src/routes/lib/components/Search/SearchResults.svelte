<!-- SearchResults.svelte -->
<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Document, SearchResponse } from '../../../utils/search_types';
	import { formatDate } from '../../../utils/utils';

	export let searchResults: SearchResponse;
	export let isLoading: boolean = false;
	export let error: string = '';
	export let currentPage: number = 1;
	export let totalPages: number = 0;

	const dispatch = createEventDispatcher<{
		openDocument: Document;
		goToPage: number;
		resetFilters: void;
	}>();

	function openDocumentViewer(document: Document) {
		dispatch('openDocument', document);
	}

	function goToPage(page: number) {
		if (page < 1 || page > totalPages) return;
		dispatch('goToPage', page);
	}

	function resetFilters() {
		dispatch('resetFilters');
	}
</script>

<div class="overflow-hidden rounded-xl border border-gray-100 bg-white shadow-sm">
	<div class="flex items-center justify-between border-b border-gray-100 p-5">
		<h2 class="text-lg font-semibold text-gray-800">Results</h2>

		<!-- Results Count -->
		<div class="text-sm font-medium">
			{#if isLoading}
				<div class="flex items-center text-gray-500">
					<svg
						class="mr-2 h-4 w-4 animate-spin"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
					>
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
						></circle>
						<path
							class="opacity-75"
							fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
						></path>
					</svg>
					Loading...
				</div>
			{:else if searchResults.total === 0}
				<span class="text-gray-500">No documents found</span>
			{:else}
				<span class="rounded-full bg-blue-50 px-3 py-1 text-blue-700">
					{searchResults.total} document{searchResults.total !== 1 ? 's' : ''}
				</span>
			{/if}
		</div>
	</div>

	<!-- Error Message -->
	{#if error}
		<div class="m-5 border-l-4 border-red-500 bg-red-50 p-4 text-red-700" role="alert">
			<p class="text-sm">{error}</p>
		</div>
	{/if}

	<!-- Results List -->
	<div class="p-5">
		{#if searchResults.hits.length > 0}
			<div class="space-y-4">
				{#each searchResults.hits as document}
					<div
						class="cursor-pointer rounded-lg border border-gray-100 p-4 shadow-sm transition-all hover:bg-gray-50"
						on:click={() => openDocumentViewer(document)}
					>
						<div class="mb-2 flex flex-wrap items-start justify-between gap-2">
							<div>
								<h3 class="text-base font-medium text-blue-700">
									{document.metadata.document_name || document.file_name}
								</h3>

								{#if document.metadata.subject}
									<h2 class="text-sm font-medium text-gray-600">
										<b>Summary: </b>
										{document.metadata.subject}
									</h2>
								{/if}
							</div>
							<div class="flex flex-wrap gap-2">
								<span class="rounded-md bg-blue-50 px-2 py-1 text-xs font-medium text-blue-700">
									{document.doc_type}
								</span>
								{#if document.metadata.status}
									<span class="rounded-md bg-gray-50 px-2 py-1 text-xs font-medium text-gray-700">
										{document.metadata.status}
									</span>
								{/if}
							</div>
						</div>

						{#if document.highlight?.text}
							<div class="mt-2 rounded-md bg-yellow-50 p-3 text-sm text-gray-700">
								{#each document.highlight.text as highlight}
									<p class="mb-1">...{@html highlight}...</p>
								{/each}
							</div>
						{:else}
							<p class="mt-2 line-clamp-2 text-sm text-gray-600">
								{document.text.substring(0, 150)}...
							</p>
						{/if}

						<div class="mt-3 grid grid-cols-2 gap-x-4 gap-y-1 text-xs sm:grid-cols-3">
							{#if document.metadata.case_number}
								<div class="flex items-center">
									<span class="text-gray-500">Case #:</span>
									<span class="ml-1 font-medium text-gray-900">{document.metadata.case_number}</span
									>
								</div>
							{/if}

							{#if document.metadata.case_name}
								<div class="flex items-center">
									<span class="text-gray-500">Case:</span>
									<span class="ml-1 truncate font-medium text-gray-900"
										>{document.metadata.case_name}</span
									>
								</div>
							{/if}

							{#if document.metadata.judge}
								<div class="flex items-center">
									<span class="text-gray-500">Judge:</span>
									<span class="ml-1 font-medium text-gray-900">{document.metadata.judge}</span>
								</div>
							{/if}

							{#if document.metadata.court}
								<div class="flex items-center">
									<span class="text-gray-500">Court:</span>
									<span class="ml-1 truncate font-medium text-gray-900"
										>{document.metadata.court}</span
									>
								</div>
							{/if}

							{#if document.metadata.timestamp}
								<div class="flex items-center">
									<span class="text-gray-500">Date:</span>
									<span class="ml-1 font-medium text-gray-900"
										>{formatDate(document.metadata.timestamp)}</span
									>
								</div>
							{:else}
								<div class="flex items-center">
									<span class="text-gray-500">Date:</span>
									<span class="ml-1 font-medium text-gray-900"
										>{formatDate(document.created_at)}</span
									>
								</div>
							{/if}

							{#if document.metadata.author}
								<div class="flex items-center">
									<span class="text-gray-500">Author:</span>
									<span class="ml-1 truncate font-medium text-gray-900"
										>{document.metadata.author}</span
									>
								</div>
							{/if}
						</div>

						{#if document.metadata.legal_tags && document.metadata.legal_tags.length > 0}
							<div class="mt-2">
								<div class="flex flex-wrap gap-1">
									{#each document.metadata.legal_tags as tag}
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
				{/each}
			</div>

			<!-- Pagination -->
			{#if totalPages > 1}
				<div class="mt-6 flex justify-center">
					<div class="inline-flex rounded-md shadow-sm" aria-label="Pagination">
						<button
							on:click={() => goToPage(currentPage - 1)}
							disabled={currentPage === 1 || isLoading}
							class="relative inline-flex items-center rounded-l-md border border-gray-200 bg-white px-2 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
						>
							<svg
								class="h-5 w-5"
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 20 20"
								fill="currentColor"
							>
								<path
									fill-rule="evenodd"
									d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z"
									clip-rule="evenodd"
								/>
							</svg>
						</button>

						{#each Array(Math.min(5, totalPages)) as _, i}
							{#if totalPages <= 5 || (i < 3 && currentPage <= 3) || (i >= 2 && currentPage > totalPages - 3)}
								<button
									on:click={() => goToPage(i + 1)}
									class={`relative inline-flex items-center border px-3 py-2 text-sm font-medium ${currentPage === i + 1 ? 'z-10 border-blue-200 bg-blue-50 text-blue-700' : 'border-gray-200 bg-white text-gray-700 hover:bg-gray-50'}`}
								>
									{i + 1}
								</button>
							{:else if i === 2 && currentPage > 3 && currentPage < totalPages - 2}
								<button
									on:click={() => goToPage(currentPage)}
									class="relative z-10 inline-flex items-center border border-blue-200 bg-blue-50 px-3 py-2 text-sm font-medium text-blue-700"
								>
									{currentPage}
								</button>
							{:else if (i === 1 && currentPage > 3) || (i === 3 && currentPage < totalPages - 2)}
								<span
									class="relative inline-flex items-center border border-gray-200 bg-white px-3 py-2 text-sm font-medium text-gray-700"
								>
									...
								</span>
							{/if}
						{/each}

						<button
							on:click={() => goToPage(currentPage + 1)}
							disabled={currentPage === totalPages || isLoading}
							class="relative inline-flex items-center rounded-r-md border border-gray-200 bg-white px-2 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
						>
							<svg
								class="h-5 w-5"
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 20 20"
								fill="currentColor"
							>
								<path
									fill-rule="evenodd"
									d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
									clip-rule="evenodd"
								/>
							</svg>
						</button>
					</div>
				</div>
			{/if}
		{:else if !isLoading}
			<div class="py-10 text-center">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mx-auto mb-4 h-12 w-12 text-gray-300"
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
				<p class="text-sm text-gray-500">No documents found matching your search criteria</p>
				<button
					type="button"
					on:click={resetFilters}
					class="mt-4 rounded-lg bg-blue-50 px-4 py-2 text-sm font-medium text-blue-700 hover:bg-blue-100"
				>
					Reset filters
				</button>
			</div>
		{/if}
	</div>
</div>
