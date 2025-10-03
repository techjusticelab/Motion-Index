<!-- SearchResults.svelte -->
<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Document, SearchResponse } from '$lib/types';
	import { formatDate } from '$lib/utils';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { cubicOut, quintOut } from 'svelte/easing';

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
		dispatch('goToPage', page);
	}

	function resetFilters() {
		dispatch('resetFilters');
	}
</script>

<div
	class="overflow-hidden rounded-xl border border-neutral-100 bg-white shadow-sm"
	in:fly={{ y: 15, duration: 600, easing: cubicOut }}
>
	<div
		class="flex items-center justify-between border-b border-neutral-100 p-5"
		in:fly={{ y: -10, duration: 500, delay: 100, easing: cubicOut }}
	>
		<h2 class="text-lg font-semibold text-neutral-800" in:slide={{ duration: 500, delay: 200 }}>
			Results
		</h2>

		<!-- Results Count -->
		<div
			class="text-sm font-medium"
			in:scale={{ start: 0.95, duration: 600, delay: 300, easing: cubicOut }}
		>
			{#if isLoading}
				<div class="flex items-center text-neutral-500">
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
				<span class="text-neutral-500">No documents found</span>
			{:else}
				<span
					class="rounded-full bg-secondary-200 px-3 py-1 text-secondary-700"
					in:scale={{ start: 0.9, duration: 500, easing: cubicOut }}
				>
					{searchResults.total} document{searchResults.total !== 1 ? 's' : ''}
				</span>
			{/if}
		</div>
	</div>

	<!-- Error Message -->
	{#if error}
		<div
			class="m-5 border-l-4 border-red-500 bg-red-50 p-4 text-red-700"
			role="alert"
			in:fly={{ x: -5, y: 0, duration: 600, easing: cubicOut }}
		>
			<p class="text-sm">{error}</p>
		</div>
	{/if}

	<!-- Results List -->
	<div class="p-5">
		{#if searchResults.hits.length > 0}
			<div class="space-y-4">
				{#each searchResults.hits as document, i}
					<div
						class="cursor-pointer rounded-lg border border-neutral-100 p-4 shadow-sm transition-all hover:bg-neutral-50"
						onclick={() => openDocumentViewer(document)}
						onkeydown={(e) => e.key === 'Enter' && openDocumentViewer(document)}
						role="button"
						tabindex="0"
						aria-label="View document {document.metadata?.document_name || 'Untitled'}"
						in:fly={{ y: 20, duration: 600, delay: 200 + i * 100, easing: cubicOut }}
					>
						<div class="mb-2 flex flex-wrap items-start justify-between gap-2">
							<div>
								<h3
									class="text-base font-medium text-primary-800"
									in:slide={{ duration: 500, delay: 250 + i * 100 }}
								>
									{document.metadata.subject}
								</h3>

								{#if document.metadata.subject}
									<h2
										class="text-sm font-medium text-neutral-600"
										in:slide={{ duration: 500, delay: 300 + i * 100 }}
									>
										{document.metadata.document_name || document.file_name}
									</h2>
								{/if}
							</div>
							<div class="flex flex-wrap gap-2">
								<span
									class="rounded-md bg-primary-100 px-2 py-1 text-xs font-medium text-primary-800"
									in:scale={{ start: 0.9, duration: 500, delay: 350 + i * 100, easing: cubicOut }}
								>
									{document.doc_type}
								</span>
								{#if document.metadata.status}
									<span
										class="rounded-md bg-neutral-100 px-2 py-1 text-xs font-medium text-neutral-700"
										in:scale={{ start: 0.9, duration: 500, delay: 400 + i * 100, easing: cubicOut }}
									>
										{document.metadata.status}
									</span>
								{/if}
							</div>
						</div>

						{#if document.highlight?.text}
							<div
								class="mt-2 rounded-md bg-secondary-100 p-3 text-sm text-neutral-700"
								in:fade={{ duration: 700, delay: 450 + i * 100 }}
							>
								{#each document.highlight.text as highlight}
									<p class="mb-1">...{@html highlight}...</p>
								{/each}
							</div>
						{:else}
							<p
								class="mt-2 line-clamp-2 text-sm text-neutral-600"
								in:fade={{ duration: 700, delay: 450 + i * 100 }}
							>
								{document.text.substring(0, 150)}...
							</p>
						{/if}

						<div
							class="mt-3 grid grid-cols-2 gap-x-4 gap-y-1 text-xs sm:grid-cols-3"
							in:fade={{ duration: 700, delay: 500 + i * 100 }}
						>
							{#if document.metadata.case_number}
								<div class="flex items-center">
									<span class="text-neutral-500">Case #:</span>
									<span class="ml-1 font-medium text-neutral-900">{document.metadata.case_number}</span
									>
								</div>
							{/if}

							{#if document.metadata.case_name}
								<div class="flex items-center">
									<span class="text-neutral-500">Case:</span>
									<span class="ml-1 truncate font-medium text-neutral-900"
										>{document.metadata.case_name}</span
									>
								</div>
							{/if}

							{#if document.metadata.judge}
								<div class="flex items-center">
									<span class="text-neutral-500">Judge:</span>
									<span class="ml-1 font-medium text-neutral-900">{document.metadata.judge}</span>
								</div>
							{/if}

							{#if document.metadata.court}
								<div class="flex items-center">
									<span class="text-neutral-500">Court:</span>
									<span class="ml-1 truncate font-medium text-neutral-900"
										>{document.metadata.court}</span
									>
								</div>
							{/if}

							{#if document.metadata.timestamp}
								<div class="flex items-center">
									<span class="text-neutral-500">Date:</span>
									<span class="ml-1 font-medium text-neutral-900"
										>{formatDate(document.metadata.timestamp)}</span
									>
								</div>
							{:else}
								<div class="flex items-center">
									<span class="text-neutral-500">Date:</span>
									<span class="ml-1 font-medium text-neutral-900"
										>{formatDate(document.created_at)}</span
									>
								</div>
							{/if}

							{#if document.metadata.author}
								<div class="flex items-center">
									<span class="text-neutral-500">Author:</span>
									<span class="ml-1 truncate font-medium text-neutral-900"
										>{document.metadata.author}</span
									>
								</div>
							{/if}
						</div>

						{#if document.metadata.legal_tags && document.metadata.legal_tags.length > 0}
							<div class="mt-2">
								<div class="flex flex-wrap gap-1">
									{#each document.metadata.legal_tags as tag, j}
										<span
											class="inline-flex rounded-full bg-neutral-100 px-2 py-0.5 text-xs font-medium text-neutral-800"
											in:scale={{
												start: 0.9,
												duration: 400,
												delay: 550 + i * 100 + j * 50,
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
				{/each}
			</div>

			<!-- Pagination -->
			{#if totalPages > 1}
				<div
					class="mt-6 flex justify-center"
					in:fly={{
						y: 15,
						duration: 600,
						delay: 300 + searchResults.hits.length * 50,
						easing: cubicOut
					}}
				>
					<div class="inline-flex rounded-md shadow-sm" aria-label="Pagination">
						<button
							onclick={() => goToPage(currentPage - 1)}
							disabled={currentPage === 1 || isLoading}
							aria-label="Previous page"
							class="relative inline-flex items-center rounded-l-md border border-neutral-200 bg-white px-2 py-2 text-sm font-medium text-neutral-500 hover:bg-neutral-50 disabled:opacity-50"
							in:scale={{
								start: 0.95,
								duration: 400,
								delay: 350 + searchResults.hits.length * 50,
								easing: cubicOut
							}}
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
									onclick={() => goToPage(i + 1)}
									class={`relative inline-flex items-center border px-3 py-2 text-sm font-medium ${currentPage === i + 1 ? 'z-10 border-primary-200 bg-primary-50 text-primary-800' : 'border-neutral-200 bg-white text-neutral-700 hover:bg-neutral-50'}`}
									in:scale={{
										start: 0.95,
										duration: 400,
										delay: 400 + searchResults.hits.length * 50 + i * 50,
										easing: cubicOut
									}}
								>
									{i + 1}
								</button>
							{:else if i === 2 && currentPage > 3 && currentPage < totalPages - 2}
								<button
									onclick={() => goToPage(currentPage)}
									class="relative z-10 inline-flex items-center border border-primary-200 bg-primary-50 px-3 py-2 text-sm font-medium text-primary-800"
									in:scale={{
										start: 0.95,
										duration: 400,
										delay: 400 + searchResults.hits.length * 50 + i * 50,
										easing: cubicOut
									}}
								>
									{currentPage}
								</button>
							{:else if (i === 1 && currentPage > 3) || (i === 3 && currentPage < totalPages - 2)}
								<span
									class="relative inline-flex items-center border border-neutral-200 bg-white px-3 py-2 text-sm font-medium text-neutral-700"
									in:scale={{
										start: 0.95,
										duration: 400,
										delay: 400 + searchResults.hits.length * 50 + i * 50,
										easing: cubicOut
									}}
								>
									...
								</span>
							{/if}
						{/each}

						<button
							onclick={() => goToPage(currentPage + 1)}
							disabled={currentPage === totalPages || isLoading}
							aria-label="Next page"
							class="relative inline-flex items-center rounded-r-md border border-neutral-200 bg-white px-2 py-2 text-sm font-medium text-neutral-500 hover:bg-neutral-50 disabled:opacity-50"
							in:scale={{
								start: 0.95,
								duration: 400,
								delay: 450 + searchResults.hits.length * 50,
								easing: cubicOut
							}}
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
			<div class="py-10 text-center" in:fly={{ y: 20, duration: 600, easing: cubicOut }}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mx-auto mb-4 h-12 w-12 text-neutral-300"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					in:scale={{ start: 0.8, duration: 700, delay: 200, easing: quintOut }}
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
					/>
				</svg>
				<p class="text-sm text-neutral-500" in:fade={{ duration: 600, delay: 400 }}>
					No documents found matching your search criteria
				</p>
				<button
					type="button"
					onclick={resetFilters}
					class="mt-4 rounded-lg bg-primary-50 px-4 py-2 text-sm font-medium text-primary-800 hover:bg-primary-100"
					in:scale={{ start: 0.95, duration: 600, delay: 600, easing: cubicOut }}
				>
					Reset filters
				</button>
			</div>
		{/if}
	</div>
</div>
