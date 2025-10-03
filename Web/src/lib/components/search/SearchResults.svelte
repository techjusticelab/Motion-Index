<!-- SearchResults.svelte -->
<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Document, SearchResponse } from '$lib/types';
	import { formatDate } from '$lib/utils';
	import { 
		extractJudgeName, 
		extractCourtInfo, 
		formatParties, 
		formatAttorneys, 
		formatCharges,
		getDocumentPriority, 
		getRelevantDate, 
		getMotionStatus, 
		formatCaseInfo,
		sortDocumentsByImportance,
		calculateDocumentImportanceScore
	} from '$lib/utils/legal-metadata';
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

	// Sort documents by importance for public defenders
	$: sortedDocuments = sortDocumentsByImportance(searchResults.hits || []);
</script>

<div
	class="overflow-hidden rounded-xl border border-neutral-100 bg-white shadow-sm"
	in:fly={{ y: 15, duration: 600, easing: cubicOut }}
>
	<div
		class="flex items-center justify-between border-b border-neutral-100 p-5"
		in:fly={{ y: -10, duration: 500, delay: 100, easing: cubicOut }}
	>
		<div class="flex items-center gap-3">
			<h2 class="text-lg font-semibold text-neutral-800" in:slide={{ duration: 500, delay: 200 }}>
				Results
			</h2>
			{#if sortedDocuments.length > 0}
				<span class="rounded-full bg-blue-100 px-2 py-1 text-xs font-medium text-blue-800" title="Documents sorted by relevance to public defenders">
					🏆 Ranked by Importance
				</span>
			{/if}
		</div>

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
		{#if sortedDocuments.length > 0}
			<div class="space-y-4">
				{#each sortedDocuments as document, i}
					{@const priority = getDocumentPriority(document)}
					{@const relevantDate = getRelevantDate(document)}
					{@const motionStatus = getMotionStatus(document)}
					{@const caseInfo = formatCaseInfo(document)}
					{@const judgeName = extractJudgeName(document.metadata?.judge)}
					{@const courtInfo = extractCourtInfo(document.metadata?.court)}
					{@const parties = formatParties(document.metadata?.parties)}
					{@const attorneys = formatAttorneys(document.metadata?.attorneys)}
					{@const charges = formatCharges(document.metadata?.charges)}
					{@const importanceScore = calculateDocumentImportanceScore(document)}
					
					<div
						class="relative cursor-pointer rounded-lg border transition-all hover:shadow-md {priority === 'high' ? 'border-red-200 bg-red-50' : priority === 'medium' ? 'border-yellow-200 bg-yellow-50' : 'border-neutral-100 bg-white'}"
						onclick={() => openDocumentViewer(document)}
						onkeydown={(e) => e.key === 'Enter' && openDocumentViewer(document)}
						role="button"
						tabindex="0"
						aria-label="View document {document.metadata?.subject || document.file_name || 'Untitled'}"
						in:fly={{ y: 20, duration: 600, delay: 200 + i * 100, easing: cubicOut }}
					>
						<!-- Priority Indicator -->
						{#if priority === 'high'}
							<div class="absolute -left-1 top-4 h-8 w-1 rounded-r bg-red-500"></div>
						{/if}
						
						<div class="p-4">
							<!-- Header: Document Type, Status, and Date -->
							<div class="mb-3 flex flex-wrap items-start justify-between gap-2">
								<div class="flex items-center gap-2">
									<!-- Importance Indicator -->
									{#if i < 3}
										<span
											class="flex h-6 w-6 items-center justify-center rounded-full text-xs font-bold {i === 0 ? 'bg-red-500 text-white' : i === 1 ? 'bg-orange-500 text-white' : 'bg-yellow-500 text-white'}"
											title="Ranked #{i + 1} by importance (Score: {importanceScore})"
											in:scale={{ start: 0.8, duration: 500, delay: 300 + i * 100, easing: cubicOut }}
										>
											{i + 1}
										</span>
									{:else if importanceScore > 800}
										<span
											class="flex h-6 w-6 items-center justify-center rounded-full bg-red-100 text-xs font-bold text-red-800"
											title="High importance (Score: {importanceScore})"
											in:scale={{ start: 0.8, duration: 500, delay: 300 + i * 100, easing: cubicOut }}
										>
											!
										</span>
									{/if}
									
									<span
										class="rounded-md px-2 py-1 text-xs font-medium {motionStatus.color === 'green' ? 'bg-green-100 text-green-800' : motionStatus.color === 'red' ? 'bg-red-100 text-red-800' : motionStatus.color === 'yellow' ? 'bg-yellow-100 text-yellow-800' : motionStatus.color === 'blue' ? 'bg-blue-100 text-blue-800' : motionStatus.color === 'purple' ? 'bg-purple-100 text-purple-800' : motionStatus.color === 'indigo' ? 'bg-indigo-100 text-indigo-800' : 'bg-neutral-100 text-neutral-800'}"
										in:scale={{ start: 0.9, duration: 500, delay: 350 + i * 100, easing: cubicOut }}
									>
										{motionStatus.type}
									</span>
									{#if motionStatus.outcome}
										<span
											class="rounded-md px-2 py-1 text-xs font-semibold {motionStatus.outcome === 'Granted' ? 'bg-green-200 text-green-900' : motionStatus.outcome === 'Denied' ? 'bg-red-200 text-red-900' : 'bg-yellow-200 text-yellow-900'}"
											in:scale={{ start: 0.9, duration: 500, delay: 400 + i * 100, easing: cubicOut }}
										>
											{motionStatus.outcome}
										</span>
									{/if}
								</div>
								<div class="flex flex-col items-end text-right">
									<span class="text-xs text-neutral-500">{relevantDate.label}</span>
									<span class="text-sm font-medium text-neutral-900">{formatDate(relevantDate.date)}</span>
								</div>
							</div>

							<!-- Title and Case Information -->
							<div class="mb-3">
								<h3
									class="text-lg font-semibold text-neutral-900 leading-tight"
									in:slide={{ duration: 500, delay: 250 + i * 100 }}
								>
									{document.metadata?.subject || document.metadata?.document_name || document.file_name || 'Untitled Document'}
								</h3>
								
								{#if caseInfo.name || caseInfo.number}
									<div class="mt-1 flex flex-wrap gap-x-4 text-sm text-neutral-700">
										{#if caseInfo.number}
											<span><strong>Case #:</strong> {caseInfo.number}</span>
										{/if}
										{#if caseInfo.name}
											<span><strong>Case:</strong> {caseInfo.name}</span>
										{/if}
									</div>
								{/if}
							</div>

							<!-- Court and Judge Information -->
							{#if courtInfo.name || judgeName}
								<div class="mb-3 grid grid-cols-1 gap-2 text-sm sm:grid-cols-2">
									{#if courtInfo.name}
										<div>
											<span class="text-neutral-500">Court:</span>
											<span class="ml-1 font-medium text-neutral-900">{courtInfo.name}</span>
											{#if courtInfo.jurisdiction}
												<span class="text-neutral-500">({courtInfo.jurisdiction})</span>
											{/if}
										</div>
									{/if}
									{#if judgeName}
										<div>
											<span class="text-neutral-500">Judge:</span>
											<span class="ml-1 font-medium text-neutral-900">{judgeName}</span>
										</div>
									{/if}
								</div>
							{/if}

							<!-- Charges Information -->
							{#if charges.length > 0}
								<div class="mb-3">
									<span class="text-sm text-neutral-500">Charges:</span>
									<div class="mt-1">
										{#each charges as charge, j}
											<div class="text-sm text-neutral-700">{charge}</div>
										{/each}
									</div>
								</div>
							{/if}

							<!-- Parties and Attorneys -->
							{#if parties || attorneys}
								<div class="mb-3 text-sm">
									{#if parties}
										<div class="text-neutral-700">{parties}</div>
									{/if}
									{#if attorneys}
										<div class="text-neutral-700">{attorneys}</div>
									{/if}
								</div>
							{/if}

							<!-- Document Content Preview -->
							{#if document.highlight?.text}
								<div
									class="mb-3 rounded-md bg-yellow-100 p-3 text-sm text-neutral-700"
									in:fade={{ duration: 700, delay: 450 + i * 100 }}
								>
									<div class="font-medium text-yellow-800 mb-1">Relevant Content:</div>
									{#each document.highlight.text as highlight}
										<p class="mb-1">...{@html highlight}...</p>
									{/each}
								</div>
							{:else if document.metadata?.summary}
								<p
									class="mb-3 text-sm text-neutral-600 line-clamp-2"
									in:fade={{ duration: 700, delay: 450 + i * 100 }}
								>
									{document.metadata.summary.substring(0, 200)}{document.metadata.summary.length > 200 ? '...' : ''}
								</p>
							{/if}

							<!-- Document Metadata -->
							<div class="mt-3 flex flex-wrap gap-x-4 gap-y-1 text-xs text-neutral-500">
								{#if document.metadata?.pages}
									<span>{document.metadata.pages} pages</span>
								{/if}
								{#if document.metadata?.word_count}
									<span>{document.metadata.word_count.toLocaleString()} words</span>
								{/if}
								{#if document.metadata?.author}
									<span>Author: {document.metadata.author}</span>
								{/if}
								{#if document.metadata?.language && document.metadata.language !== 'en'}
									<span>Language: {document.metadata.language.toUpperCase()}</span>
								{/if}
							</div>

							<!-- Legal Tags -->
							{#if document.metadata?.legal_tags && document.metadata.legal_tags.length > 0}
								<div class="mt-3">
									<div class="flex flex-wrap gap-1">
										{#each document.metadata.legal_tags.slice(0, 5) as tag, j}
											<span
												class="inline-flex rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-800"
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
										{#if document.metadata.legal_tags.length > 5}
											<span class="inline-flex rounded-full bg-neutral-100 px-2 py-0.5 text-xs font-medium text-neutral-600">
												+{document.metadata.legal_tags.length - 5} more
											</span>
										{/if}
									</div>
								</div>
							{/if}
						</div>
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
						delay: 300 + sortedDocuments.length * 50,
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
								delay: 350 + sortedDocuments.length * 50,
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
										delay: 400 + sortedDocuments.length * 50 + i * 50,
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
										delay: 400 + sortedDocuments.length * 50 + i * 50,
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
										delay: 400 + sortedDocuments.length * 50 + i * 50,
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
								delay: 450 + sortedDocuments.length * 50,
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
