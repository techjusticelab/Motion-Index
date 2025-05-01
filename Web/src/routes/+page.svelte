<!-- SearchPage.svelte -->
<script lang="ts">
	import { onMount } from 'svelte';
	import * as api from './api';
	import DocumentViewer from './lib/components/Document/DocumentViewer.svelte';
	import type {
		SearchParams,
		SearchResponse,
		Document,
		MetadataField,
		DocumentStats
	} from './utils/search_types';

	import SearchForm from './lib/components/Search/SearchForm.svelte';

	import SearchResults from './lib/components/Search/SearchResults.svelte';
	import { formatDate } from './utils/utils';
	import { fade } from 'svelte/transition';

	// Document popup state
	let activeDocument: Document | null = null;
	let showDocumentPopup = false;

	// State variables
	let searchParams: SearchParams = {
		query: '',
		doc_type: '',
		case_number: '',
		case_name: '',
		judge: [],
		court: [],
		author: '',
		status: '',
		date_range: {
			start: '',
			end: ''
		},
		legal_tags: [],
		legal_tags_match_all: false, // Default to OR behavior (match any tag)
		size: 10,
		sort_by: 'created_at',
		sort_order: 'desc',
		page: 1,
		use_fuzzy: false
	};

	let searchResults: SearchResponse = { total: 0, hits: [] };
	let isLoading = false;
	let error = '';
	let documentTypes: Record<string, number> = {};
	let metadataFields: MetadataField[] = [];
	let documentStats: DocumentStats | null = null;
	let fieldOptions: Record<string, string[]> = {};
	// UI state
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
			console.log('Document stats:', documentStats);
			if (documentStats?.date_range) {
				searchParams.date_range.start = documentStats.date_range.oldest;
				searchParams.date_range.end = documentStats.date_range.newest;
			}
			if (documentStats?.total_documents) {
				totalPages = Math.min(searchParams.size, documentStats.total_documents);
			}
			// Fetch all field options
			fieldOptions = await api.getAllFieldOptions();
			console.log('Field options:', fieldOptions);

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
			console.log('Total pages:', totalPages);
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
	function openDocumentViewer(document: Document) {
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
		console.log(searchParams);
		console.log('Going to page:', page);
		performSearch();
	}

	// Reset all filters
	function resetFilters() {
		searchParams = {
			query: '',
			doc_type: '',
			case_number: '',
			case_name: '',
			judge: [],
			court: [],
			author: '',
			status: '',
			date_range: {
				start: documentStats?.date_range?.oldest || '',
				end: documentStats?.date_range?.newest || ''
			},
			legal_tags: [],
			legal_tags_match_all: false, // Reset to OR behavior
			size: 10,
			sort_by: 'created_at',
			sort_order: 'desc',
			page: 1,
			use_fuzzy: false
		};
		performSearch();
	}
</script>

<div class="max-w-7/8 container mx-auto px-4 py-6" transition:fade>
	<!-- Header -->
	<div class="mb-6 rounded-xl bg-gradient-to-r from-blue-600 to-indigo-700 p-6 shadow-lg">
		<div class="flex flex-col items-start justify-between md:flex-row md:items-center">
			<div class="mb-6 md:mb-0">
				<h1 class="mb-2 text-2xl font-bold text-white md:text-3xl">Motion Index</h1>
				<p class="text-sm text-blue-100 md:text-base">Search legal documents with precision</p>
			</div>

			{#if documentStats}
				<div class="grid w-full grid-cols-2 gap-3 sm:grid-cols-2 md:w-auto">
					<div class="rounded-lg bg-white/10 p-3 text-center backdrop-blur-sm">
						<p class="text-xs text-blue-100">Documents</p>
						<p class="text-xl font-semibold text-white">{documentStats.total_documents}</p>
					</div>

					<!-- <div class="rounded-lg bg-white/10 p-3 text-center backdrop-blur-sm">
						<p class="text-xs text-blue-100">Types</p>
						<p class="text-xl font-semibold text-white">{Object.keys(documentTypes).length}</p>
					</div> -->

					{#if documentStats.date_range}
						<div class="rounded-lg bg-white/10 p-3 text-center backdrop-blur-sm">
							<p class="text-xs text-blue-100">Date Range</p>
							<p class="truncate text-sm font-medium text-white">
								{formatDate(documentStats.date_range.oldest)} - {formatDate(
									documentStats.date_range.newest
								)}
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
			<SearchForm
				{searchParams}
				{isLoading}
				{documentTypes}
				{fieldOptions}
				on:search={handleSearch}
				on:reset={resetFilters}
			/>
		</div>

		<!-- Search Results -->
		<div class="w-full lg:w-2/3 {activeTab !== 'results' ? 'hidden lg:block' : ''}">
			<SearchResults
				{searchResults}
				{isLoading}
				{error}
				currentPage={searchParams.page}
				{totalPages}
				on:openDocument={(e) => openDocumentViewer(e.detail)}
				on:goToPage={(e) => goToPage(e.detail)}
				on:resetFilters={resetFilters}
			/>
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
