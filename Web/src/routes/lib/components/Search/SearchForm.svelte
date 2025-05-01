<!-- SearchForm.svelte -->
<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { SearchParams } from '../../../utils/search_types';
	import CourtFilter from './CourtFilter.svelte';
	import JudgeFilter from './JudgeFilter.svelte';
	import LegalFilter from './LegalFilter.svelte';

	export let searchParams: SearchParams;
	export let isLoading: boolean = false;
	export let documentTypes: Record<string, number> = {};
	export let fieldOptions: Record<string, string[]> = {};

	let showSearchHelp = false;
	let showAdvancedFilters = false;

	const dispatch = createEventDispatcher<{
		search: Event;
		reset: void;
	}>();

	// Handle search form submission
	function handleSearch(event: Event) {
		event.preventDefault();
		searchParams.page = 1; // Reset to first page
		dispatch('search', event);
	}

	// Reset all filters
	function resetFilters() {
		dispatch('reset');
	}

	// Toggle search help
	function toggleSearchHelp() {
		showSearchHelp = !showSearchHelp;
	}

	// Toggle advanced filters
	function toggleAdvancedFilters() {
		showAdvancedFilters = !showAdvancedFilters;
	}

	// Handle court selection
	function handleAddCourt(event: CustomEvent<string>) {
		const court = event.detail;
		if (!searchParams.court.includes(court)) {
			searchParams.court = [...searchParams.court, court];
		}
	}

	function handleRemoveCourt(event: CustomEvent<string>) {
		const court = event.detail;
		searchParams.court = searchParams.court.filter((c) => c !== court);
	}

	// Handle judge selection
	function handleAddJudge(event: CustomEvent<string>) {
		const judge = event.detail;
		if (!searchParams.judge.includes(judge)) {
			searchParams.judge = [...searchParams.judge, judge];
		}
	}

	function handleRemoveTag(event: CustomEvent<string>) {
		const tag = event.detail;
		searchParams.legal_tags = searchParams.legal_tags.filter((t) => t !== tag);
	}

	function handleAddTag(event: CustomEvent<string>) {
		const tag = event.detail;
		if (!searchParams.legal_tags.includes(tag)) {
			searchParams.legal_tags = [...searchParams.legal_tags, tag];
		}
	}
	function handleRemoveJudge(event: CustomEvent<string>) {
		const judge = event.detail;
		searchParams.judge = searchParams.judge.filter((j) => j !== judge);
	}
</script>

<div class="overflow-hidden rounded-xl border border-gray-100 bg-white shadow-sm">
	<!-- Search Form Header -->
	<div class="flex items-center justify-between px-5 pb-3 pt-5">
		<h2 class="text-lg font-semibold text-gray-800">Search</h2>
		<button
			type="button"
			class="flex items-center text-xs font-medium text-blue-600 hover:text-blue-800"
			on:click={toggleSearchHelp}
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="mr-1 h-4 w-4"
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
			Search Tips
		</button>
	</div>

	<!-- Search Help Panel -->
	{#if showSearchHelp}
		<div class="mx-5 mb-4 rounded-lg bg-blue-50 p-4 text-xs">
			<h3 class="mb-2 font-medium text-blue-800">Search Operators</h3>
			<div class="grid grid-cols-2 gap-2">
				<div class="rounded border border-blue-100 bg-white p-2">
					<code class="text-blue-700">"exact phrase"</code>
					<span class="mt-1 block text-gray-600">Exact match</span>
				</div>
				<div class="rounded border border-blue-100 bg-white p-2">
					<code class="text-blue-700">term1 OR term2</code>
					<span class="mt-1 block text-gray-600">Either term</span>
				</div>
				<div class="rounded border border-blue-100 bg-white p-2">
					<code class="text-blue-700">+required</code>
					<span class="mt-1 block text-gray-600">Must include</span>
				</div>
				<div class="rounded border border-blue-100 bg-white p-2">
					<code class="text-blue-700">-excluded</code>
					<span class="mt-1 block text-gray-600">Must exclude</span>
				</div>
			</div>
		</div>
	{/if}

	<form on:submit={handleSearch} class="px-5 pb-5">
		<!-- Text Search -->
		<div class="mb-4">
			<div class="relative">
				<input
					type="text"
					id="query"
					bind:value={searchParams.query}
					placeholder="Search documents..."
					class="block w-full rounded-lg border-gray-200 py-3 pl-10 pr-4 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
				/>
				<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-5 w-5 text-gray-400"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
						/>
					</svg>
				</div>
			</div>
		</div>

		<!-- Primary Filters -->
		<div class="mb-4 space-y-3">
			<div>
				<label for="doc_type" class="mb-1 block text-xs font-medium text-gray-700"
					>Document Type</label
				>
				<select
					id="doc_type"
					bind:value={searchParams.doc_type}
					class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
				>
					<option value="">All Document Types</option>
					{#if fieldOptions.doc_type}
						{#each fieldOptions.doc_type as type}
							<option value={type}>{type}</option>
						{/each}
					{:else}
						{#each Object.entries(documentTypes) as [type, count]}
							<option value={type}>{type} ({count})</option>
						{/each}
					{/if}
				</select>
			</div>

			<div>
				<label for="legal_tags" class="mb-1 block text-xs font-medium text-gray-700">Tags</label>
				<select
					id="legal_tags"
					bind:value={searchParams.legal_tags}
					class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
				>
					<option value="">All Tags</option>
					{#if fieldOptions.legal_tags}
						{#each fieldOptions.legal_tags as tag}
							<option value={tag}>{tag}</option>
						{/each}
					{/if}
				</select>
			</div>

			<div>
				<label for="case_number" class="mb-1 block text-xs font-medium text-gray-700"
					>Case Number</label
				>
				<select
					id="case_number"
					bind:value={searchParams.case_number}
					class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
				>
					<option value="">All Case Numbers</option>
					{#if fieldOptions.case_number}
						{#each fieldOptions.case_number as caseNum}
							<option value={caseNum}>{caseNum}</option>
						{/each}
					{/if}
				</select>
			</div>
		</div>

		<!-- Advanced Filters Toggle -->
		<button
			type="button"
			on:click={toggleAdvancedFilters}
			class="mb-4 flex items-center text-xs font-medium text-blue-600 hover:text-blue-800"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="mr-1 h-4 w-4"
				viewBox="0 0 20 20"
				fill="currentColor"
			>
				<path
					fill-rule="evenodd"
					d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
					clip-rule="evenodd"
				/>
			</svg>
			{showAdvancedFilters ? 'Hide Advanced Filters' : 'Show Advanced Filters'}
		</button>

		<!-- Advanced Filters -->
		{#if showAdvancedFilters}
			<div class="mb-4 space-y-3 border-t border-gray-100 pt-2">
				<div>
					<label for="case_name" class="mb-1 block text-xs font-medium text-gray-700"
						>Case Name</label
					>
					<input
						type="text"
						id="case_name"
						bind:value={searchParams.case_name}
						placeholder="Enter case name"
						class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
					/>
				</div>

				<!-- Judge Filter Component -->
				<JudgeFilter
					selectedJudges={searchParams.judge}
					allJudgeOptions={fieldOptions.judge || []}
					on:add={handleAddJudge}
					on:remove={handleRemoveJudge}
				/>

				<!-- Court Filter Component -->
				<CourtFilter
					selectedCourts={searchParams.court}
					allCourtOptions={fieldOptions.court || []}
					on:add={handleAddCourt}
					on:remove={handleRemoveCourt}
				/>

				<LegalFilter
					selectedTags={searchParams.legal_tags}
					allTagsOptions={fieldOptions.legal_tags || []}
					on:add={handleAddTag}
					on:remove={handleRemoveTag}
				/>
				<!-- Date Range -->
				<div>
					<label class="mb-1 block text-xs font-medium text-gray-700">Date Range</label>
					<div class="grid grid-cols-2 gap-2">
						<div>
							<input
								type="date"
								id="date_start"
								bind:value={searchParams.date_range.start}
								class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
							/>
						</div>
						<div>
							<input
								type="date"
								id="date_end"
								bind:value={searchParams.date_range.end}
								class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
							/>
						</div>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-2">
					<div>
						<label for="author" class="mb-1 block text-xs font-medium text-gray-700">Author</label>
						<input
							type="text"
							id="author"
							bind:value={searchParams.author}
							placeholder="Document author"
							class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
						/>
					</div>

					<div>
						<label for="status" class="mb-1 block text-xs font-medium text-gray-700">Status</label>
						<select
							id="status"
							bind:value={searchParams.status}
							class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
						>
							<option value="">All Statuses</option>
							{#if fieldOptions.status}
								{#each fieldOptions.status as status}
									<option value={status}>{status}</option>
								{/each}
							{/if}
						</select>
					</div>
				</div>
			</div>
		{/if}

		<!-- Options -->
		<div class="mb-5 grid grid-cols-2 gap-2">
			<div>
				<label for="sort_by" class="mb-1 block text-xs font-medium text-gray-700">Sort By</label>
				<select
					bind:value={searchParams.sort_by}
					class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
				>
					<option value="created_at">Date</option>
					<option value="metadata.document_name">Name</option>
					<option value="doc_type">Type</option>
					<option value="metadata.case_number">Case #</option>
				</select>
			</div>

			<div>
				<label for="sort_order" class="mb-1 block text-xs font-medium text-gray-700">Order</label>
				<select
					bind:value={searchParams.sort_order}
					class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
				>
					<option value="desc">Newest first</option>
					<option value="asc">Oldest first</option>
				</select>
			</div>
		</div>

		<!-- Action Buttons -->
		<div class="flex gap-2">
			<button
				type="submit"
				class="flex flex-1 items-center justify-center rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
				disabled={isLoading}
			>
				{#if isLoading}
					<svg
						class="mr-2 h-4 w-4 animate-spin text-white"
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
					Searching...
				{:else}
					Search
				{/if}
			</button>

			<button
				type="button"
				on:click={resetFilters}
				class="flex-1 rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-200"
				disabled={isLoading}
			>
				Reset
			</button>
		</div>
	</form>
</div>
