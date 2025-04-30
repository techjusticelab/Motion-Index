<script lang="ts">
	import { onMount } from 'svelte';
	import { format } from 'date-fns';
	import * as api from './api';
	import DocumentViewer from './lib/components/DocumentViewer.svelte';

	// Filter state for dropdowns
	let courtSearchInput = '';
	let filteredCourtOptions: string[] = [];
	let showCourtDropdown = false;
	let courtDropdownRef: HTMLDivElement;

	let judgeSearchInput = '';
	let filteredJudgeOptions: string[] = [];
	let showJudgeDropdown = false;
	let judgeDropdownRef: HTMLDivElement;

	// Document popup state and functions
	let activeDocument: { metadata: { document_name: any }; file_name: any; s3_uri: any } | null =
		null;
	let showDocumentPopup = false;

	// $: if (showDocumentPopup) {
	// 	activeDocument = null;
	// }

	// State variables
	let searchParams = {
		query: '',
		doc_type: '',
		case_number: '',
		case_name: '',
		judge: [] as string[], // Changed to array for multi-select
		court: [] as string[], // Changed to array for multi-select
		author: '',
		status: '',
		date_range: {
			start: '',
			end: ''
		},
		size: 10,
		sort_by: 'created_at',
		sort_order: 'desc' as 'asc' | 'desc',
		page: 1,
		use_fuzzy: false
	};

	let searchResults: api.SearchResponse = { total: 0, hits: [] };
	let isLoading = false;
	let error = '';
	let documentTypes: Record<string, number> = {};
	let metadataFields: api.MetadataField[] = [];
	let documentStats: any = null;
	let fieldOptions: Record<string, string[]> = {};

	// UI state
	let showSearchHelp = false;
	let showAdvancedFilters = false;
	let activeTab = 'search'; // 'search' or 'results'

	// Pagination
	let totalPages = 0;

	// Fetch initial data
	onMount(async () => {
		try {
			isLoading = true;

			// Fetch document types for filter dropdown
			documentTypes = await api.getDocumentTypes();

			// Fetch metadata fields
			const fieldsResponse = await api.getMetadataFields();
			metadataFields = fieldsResponse.fields;

			// Fetch document stats
			documentStats = await api.getDocumentStats();

			// Fetch all field options
			fieldOptions = await api.getAllFieldOptions();

			// Initial search
			await performSearch();
		} catch (err) {
			console.error('Error initializing search page:', err);
			error = 'Failed to load initial data. Please try refreshing the page.';
		} finally {
			isLoading = false;
		}
	});

	// Search function
	async function performSearch() {
		try {
			isLoading = true;
			error = '';

			// Clean up empty filters
			const cleanParams = { ...searchParams };
			Object.keys(cleanParams).forEach((key) => {
				if (key !== 'date_range' && !cleanParams[key]) {
					delete cleanParams[key];
				}
			});

			// Clean up empty date range
			if (!cleanParams.date_range?.start && !cleanParams.date_range?.end) {
				delete cleanParams.date_range;
			}

			searchResults = await api.searchDocuments(cleanParams);
			console.log('Search results:', searchResults);
			totalPages = Math.ceil(searchResults.total / searchParams.size);

			// Switch to results tab if we have results and on mobile
			if (searchResults.total > 0 && window.innerWidth < 1024) {
				activeTab = 'results';
			}
		} catch (err) {
			console.error('Search error:', err);
			error = 'An error occurred while searching. Please try again.';
			searchResults = { total: 0, hits: [] };
		} finally {
			isLoading = false;
		}
	}

	// Open document viewer
	function openDocumentViewer(document) {
		activeDocument = document;
		showDocumentPopup = true;
	}

	// Handle search form submission
	function handleSearch(event: Event) {
		event.preventDefault();
		searchParams.page = 1; // Reset to first page
		performSearch();
	}

	// Handle pagination
	function goToPage(page: number) {
		if (page < 1 || page > totalPages) return;
		searchParams.page = page;
		performSearch();
	}

	// Reset all filters
	function resetFilters() {
		searchParams = {
			query: '',
			doc_type: '',
			case_number: '',
			case_name: '',
			judge: [] as string[], // Changed to array for multi-select
			court: [] as string[], // Changed to array for multi-select
			author: '',
			status: '',
			date_range: {
				start: '',
				end: ''
			},
			size: 10,
			sort_by: 'created_at',
			sort_order: 'desc',
			page: 1,
			use_fuzzy: false
		};
		performSearch();
	}

	// Format date for display
	function formatDate(dateString: string): string {
		if (!dateString) return 'N/A';
		try {
			return format(new Date(dateString), 'MMM d, yyyy');
		} catch (err) {
			return dateString;
		}
	}

	// Handle clicks outside the dropdowns
	function handleClickOutside(event: MouseEvent) {
		if (courtDropdownRef && !courtDropdownRef.contains(event.target as Node)) {
			showCourtDropdown = false;
		}
		if (judgeDropdownRef && !judgeDropdownRef.contains(event.target as Node)) {
			showJudgeDropdown = false;
		}
	}

	// Add and remove event listeners for clicks outside the dropdown
	onMount(() => {
		document.addEventListener('click', handleClickOutside);
		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
	});

	// Filter options based on search input
	function filterOptions(field: 'court' | 'judge', searchInput: string) {
		if (!fieldOptions[field]) return [];

		if (!searchInput) {
			return [...fieldOptions[field]];
		}

		const searchLower = searchInput.toLowerCase();
		return fieldOptions[field].filter((item) => item.toLowerCase().includes(searchLower));
	}

	// Filter court options
	function filterCourtOptions() {
		filteredCourtOptions = filterOptions('court', courtSearchInput);
	}

	// Filter judge options
	function filterJudgeOptions() {
		filteredJudgeOptions = filterOptions('judge', judgeSearchInput);
	}

	// Add an item to a multi-select field
	function addItem(field: 'court' | 'judge', item: string) {
		if (!searchParams[field].includes(item)) {
			searchParams[field] = [...searchParams[field], item];
		}
	}

	// Add a court to the selected courts
	function addCourt(court: string) {
		addItem('court', court);
		courtSearchInput = '';
		filterCourtOptions();
	}

	// Add a judge to the selected judges
	function addJudge(judge: string) {
		addItem('judge', judge);
		judgeSearchInput = '';
		filterJudgeOptions();
	}

	// Remove an item from a multi-select field
	function removeItem(field: 'court' | 'judge', item: string) {
		searchParams[field] = searchParams[field].filter((i) => i !== item);
	}

	// Remove a court from the selected courts
	function removeCourt(court: string) {
		removeItem('court', court);
	}

	// Remove a judge from the selected judges
	function removeJudge(judge: string) {
		removeItem('judge', judge);
	}

	// Toggle search help
	function toggleSearchHelp() {
		showSearchHelp = !showSearchHelp;
	}

	// Toggle advanced filters
	function toggleAdvancedFilters() {
		showAdvancedFilters = !showAdvancedFilters;
	}
</script>

<div class="max-w-7/8 container mx-auto px-4 py-6">
	<!-- Header -->
	<div class="mb-6 rounded-xl bg-gradient-to-r from-blue-600 to-indigo-700 p-6 shadow-lg">
		<div class="flex flex-col items-start justify-between md:flex-row md:items-center">
			<div class="mb-6 md:mb-0">
				<h1 class="mb-2 text-2xl font-bold text-white md:text-3xl">Motion Index</h1>
				<p class="text-sm text-blue-100 md:text-base">Search legal documents with precision</p>
			</div>

			{#if documentStats}
				<div class="grid w-full grid-cols-2 gap-3 sm:grid-cols-3 md:w-auto">
					<div class="rounded-lg bg-white/10 p-3 text-center backdrop-blur-sm">
						<p class="text-xs text-blue-100">Documents</p>
						<p class="text-xl font-semibold text-white">{documentStats.total_documents}</p>
					</div>

					<div class="rounded-lg bg-white/10 p-3 text-center backdrop-blur-sm">
						<p class="text-xs text-blue-100">Types</p>
						<p class="text-xl font-semibold text-white">{Object.keys(documentTypes).length}</p>
					</div>

					{#if documentStats.date_range}
						<div class="rounded-lg bg-white/10 p-3 text-center backdrop-blur-sm">
							<p class="text-xs text-blue-100">Date Range</p>
							<p class="truncate text-sm font-medium text-white">
								{formatDate(documentStats.date_range.oldest).slice(0, 6)} - {formatDate(
									documentStats.date_range.newest
								).slice(0, 6)}
							</p>
						</div>
					{/if}
				</div>
			{/if}
		</div>
	</div>

	<!-- Mobile Tabs -->
	<div class="mb-4 lg:hidden">
		<div class="flex overflow-hidden rounded-lg bg-white shadow-sm">
			<button
				class={`flex-1 py-3 text-center font-medium transition-all ${activeTab === 'search' ? 'border-b-2 border-blue-600 bg-blue-50 text-blue-700' : 'text-gray-600'}`}
				on:click={() => (activeTab = 'search')}
			>
				<div class="flex items-center justify-center">
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
							d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
						/>
					</svg>
					Search
				</div>
			</button>
			<button
				class={`flex-1 py-3 text-center font-medium transition-all ${activeTab === 'results' ? 'border-b-2 border-blue-600 bg-blue-50 text-blue-700' : 'text-gray-600'}`}
				on:click={() => (activeTab = 'results')}
			>
				<div class="flex items-center justify-center">
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
							d="M4 6h16M4 10h16M4 14h16M4 18h16"
						/>
					</svg>
					Results {searchResults.total > 0 ? `(${searchResults.total})` : ''}
				</div>
			</button>
		</div>
	</div>

	<div class="flex flex-col gap-5 lg:flex-row">
		<!-- Search Filters -->
		<div class="w-full lg:w-1/3 {activeTab !== 'search' ? 'hidden lg:block' : ''}">
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

							<div>
								<label for="judge-search" class="mb-1 block text-xs font-medium text-gray-700"
									>Judges</label
								>

								<!-- Selected judges tags -->
								{#if searchParams.judge.length > 0}
									<div class="mb-2 flex flex-wrap gap-2">
										{#each searchParams.judge as judge}
											<div
												class="flex items-center rounded-lg bg-blue-50 px-2 py-1 text-xs text-blue-700"
											>
												<span class="mr-1 max-w-[200px] truncate">{judge}</span>
												<button
													type="button"
													on:click={() => removeJudge(judge)}
													class="ml-1 text-blue-500 hover:text-blue-700"
												>
													<svg
														xmlns="http://www.w3.org/2000/svg"
														class="h-3 w-3"
														viewBox="0 0 20 20"
														fill="currentColor"
													>
														<path
															fill-rule="evenodd"
															d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
															clip-rule="evenodd"
														/>
													</svg>
												</button>
											</div>
										{/each}
									</div>
								{/if}

								<!-- Judge search input -->
								<div class="relative" bind:this={judgeDropdownRef}>
									<input
										type="text"
										id="judge-search"
										bind:value={judgeSearchInput}
										on:input={filterJudgeOptions}
										on:focus={() => {
											showJudgeDropdown = true;
											filterJudgeOptions();
										}}
										placeholder="Search judges..."
										class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
									/>

									<!-- Dropdown for judge options -->
									{#if showJudgeDropdown && filteredJudgeOptions.length > 0}
										<div
											class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-sm shadow-lg ring-1 ring-black ring-opacity-5"
										>
											{#each filteredJudgeOptions as judge}
												<button
													type="button"
													class="block w-full px-4 py-2 text-left hover:bg-gray-100 {searchParams.judge.includes(
														judge
													)
														? 'bg-blue-50'
														: ''}"
													on:click={() => {
														addJudge(judge);
														showJudgeDropdown = false;
													}}
												>
													{judge}
												</button>
											{/each}
										</div>
									{/if}
								</div>
								<p class="mt-1 text-xs text-gray-500">Search and click to add multiple judges</p>
							</div>

							<div>
								<label for="court-search" class="mb-1 block text-xs font-medium text-gray-700"
									>Courts</label
								>

								<!-- Selected courts tags -->
								{#if searchParams.court.length > 0}
									<div class="mb-2 flex flex-wrap gap-2">
										{#each searchParams.court as court}
											<div
												class="flex items-center rounded-lg bg-blue-50 px-2 py-1 text-xs text-blue-700"
											>
												<span class="mr-1 max-w-[200px] truncate">{court}</span>
												<button
													type="button"
													on:click={() => removeCourt(court)}
													class="ml-1 text-blue-500 hover:text-blue-700"
												>
													<svg
														xmlns="http://www.w3.org/2000/svg"
														class="h-3 w-3"
														viewBox="0 0 20 20"
														fill="currentColor"
													>
														<path
															fill-rule="evenodd"
															d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
															clip-rule="evenodd"
														/>
													</svg>
												</button>
											</div>
										{/each}
									</div>
								{/if}

								<!-- Court search input -->
								<div class="relative" bind:this={courtDropdownRef}>
									<input
										type="text"
										id="court-search"
										bind:value={courtSearchInput}
										on:input={filterCourtOptions}
										on:focus={() => {
											showCourtDropdown = true;
											filterCourtOptions();
										}}
										placeholder="Search courts..."
										class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
									/>

									<!-- Dropdown for court options -->
									{#if showCourtDropdown && filteredCourtOptions.length > 0}
										<div
											class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-sm shadow-lg ring-1 ring-black ring-opacity-5"
										>
											{#each filteredCourtOptions as court}
												<button
													type="button"
													class="block w-full px-4 py-2 text-left hover:bg-gray-100 {searchParams.court.includes(
														court
													)
														? 'bg-blue-50'
														: ''}"
													on:click={() => {
														addCourt(court);
														showCourtDropdown = false;
													}}
												>
													{court}
												</button>
											{/each}
										</div>
									{/if}
								</div>
								<p class="mt-1 text-xs text-gray-500">Search and click to add multiple courts</p>
							</div>

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
									<label for="author" class="mb-1 block text-xs font-medium text-gray-700"
										>Author</label
									>
									<input
										type="text"
										id="author"
										bind:value={searchParams.author}
										placeholder="Document author"
										class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
									/>
								</div>

								<div>
									<label for="status" class="mb-1 block text-xs font-medium text-gray-700"
										>Status</label
									>
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
							<label for="sort_by" class="mb-1 block text-xs font-medium text-gray-700"
								>Sort By</label
							>
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
							<label for="sort_order" class="mb-1 block text-xs font-medium text-gray-700"
								>Order</label
							>
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
									<circle
										class="opacity-25"
										cx="12"
										cy="12"
										r="10"
										stroke="currentColor"
										stroke-width="4"
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
		</div>

		<!-- Search Results -->
		<div class="w-full lg:w-2/3 {activeTab !== 'results' ? 'hidden lg:block' : ''}">
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
									<circle
										class="opacity-25"
										cx="12"
										cy="12"
										r="10"
										stroke="currentColor"
										stroke-width="4"
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
										<h3 class="text-base font-medium text-blue-700">
											{document.metadata.document_name || document.file_name}
										</h3>
										<span class="rounded-md bg-blue-50 px-2 py-1 text-xs font-medium text-blue-700">
											{document.doc_type}
										</span>
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

									<div class="mt-3 grid grid-cols-2 gap-x-4 gap-y-1 text-xs">
										{#if document.metadata.case_number}
											<div class="flex items-center">
												<span class="text-gray-500">Case #:</span>
												<span class="ml-1 font-medium text-gray-900"
													>{document.metadata.case_number}</span
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
												<span class="ml-1 font-medium text-gray-900">{document.metadata.judge}</span
												>
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

										<div class="flex items-center">
											<span class="text-gray-500">Date:</span>
											<span class="ml-1 font-medium text-gray-900"
												>{formatDate(document.metadata.timestamp || document.created_at)}</span
											>
										</div>
									</div>
								</div>
							{/each}
						</div>

						<!-- Pagination -->
						{#if totalPages > 1}
							<div class="mt-6 flex justify-center">
								<div class="inline-flex rounded-md shadow-sm" aria-label="Pagination">
									<button
										on:click={() => goToPage(searchParams.page - 1)}
										disabled={searchParams.page === 1 || isLoading}
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
										{#if totalPages <= 5 || (i < 3 && searchParams.page <= 3) || (i >= 2 && searchParams.page > totalPages - 3)}
											<button
												on:click={() => goToPage(i + 1)}
												class={`relative inline-flex items-center border px-3 py-2 text-sm font-medium ${searchParams.page === i + 1 ? 'z-10 border-blue-200 bg-blue-50 text-blue-700' : 'border-gray-200 bg-white text-gray-700 hover:bg-gray-50'}`}
											>
												{i + 1}
											</button>
										{:else if i === 2 && searchParams.page > 3 && searchParams.page < totalPages - 2}
											<button
												on:click={() => goToPage(searchParams.page)}
												class="relative z-10 inline-flex items-center border border-blue-200 bg-blue-50 px-3 py-2 text-sm font-medium text-blue-700"
											>
												{searchParams.page}
											</button>
										{:else if (i === 1 && searchParams.page > 3) || (i === 3 && searchParams.page < totalPages - 2)}
											<span
												class="relative inline-flex items-center border border-gray-200 bg-white px-3 py-2 text-sm font-medium text-gray-700"
											>
												...
											</span>
										{/if}
									{/each}

									<button
										on:click={() => goToPage(searchParams.page + 1)}
										disabled={searchParams.page === totalPages || isLoading}
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
		</div>
	</div>
</div>

<!-- Document Viewer Component -->
<DocumentViewer
	docData={activeDocument}
	isOpen={showDocumentPopup}
	on:close={() => {
		activeDocument = null;
		showDocumentPopup = false;
	}}
/>
