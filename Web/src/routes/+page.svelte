<!-- SearchPage.svelte -->
<script lang="ts">
	import { onMount } from 'svelte';
	import * as api from '$lib/api';
	
	let { data } = $props();
	let { session, user, supabase } = $derived(data);
	import DocumentViewer from '$lib/components/document/DocumentViewer.svelte';
	import type {
		SearchParams,
		SearchResponse,
		Document,
		MetadataField,
		DocumentStats
	} from '$lib/types';

	import SearchForm from '$lib/components/search/SearchForm.svelte';

	import SearchResults from '$lib/components/search/SearchResults.svelte';
	import { formatDate } from '$lib/utils/utils';
	import { fade } from 'svelte/transition';

	// Document popup state
	let activeDocument = $state<Document | null>(null);
	let showDocumentPopup = $state(false);

	// State variables
	let searchParams = $state({
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
		date_field_type: 'created_at', // Enhanced date field selection
		legal_tags: [],
		legal_tags_match_all: false, // Default to OR behavior (match any tag)
		size: 10,
		sort_by: 'created_at',
		sort_order: 'desc',
		page: 1,
		use_fuzzy: false,
		include_highlights: true // Enhanced search with highlights
	});

	let searchResults = $state<SearchResponse>({ total: 0, hits: [] });
	let isLoading = $state(false);
	let error = $state('');
	let documentTypes = $state<Record<string, number>>({});
	let metadataFields = $state<MetadataField[]>([]);
	let documentStats = $state<DocumentStats | null>(null);
	let fieldOptions = $state<Record<string, string[]>>({});
	let legalTags = $state<string[]>([]);
	// UI state
	let activeTab = $state('search'); // 'search' or 'results'

	// Pagination
	let totalPages = $state(0);

	// Fetch initial data
	onMount(async () => {
		try {
			isLoading = true;

			// Fetch document types for filter dropdown
			documentTypes = await api.getDocumentTypes();
			
			// Skip metadata fields - endpoint doesn't exist in Go Fiber API
			// The API only has /api/v1/metadata-fields/:field with field parameter
			console.log('Skipping metadata fields - using field-options instead');

			// Fetch document stats
			try {
				documentStats = await api.getDocumentStats();
				console.log('Document stats:', documentStats);
				if (documentStats?.date_range) {
					searchParams.date_range.start = documentStats.date_range.oldest;
					searchParams.date_range.end = documentStats.date_range.newest;
				}
				if (documentStats?.total_documents) {
					totalPages = Math.min(searchParams.size, documentStats.total_documents);
				}
			} catch (err) {
				console.warn('Could not fetch document stats:', err);
			}

			// Fetch all field options
			try {
				fieldOptions = await api.getAllFieldOptions();
				console.log('Field options:', fieldOptions);
			} catch (err) {
				console.warn('Could not fetch field options:', err);
			}

			// Fetch legal tags
			try {
				legalTags = await api.getLegalTags(session);
				console.log('Legal tags:', legalTags);
			} catch (err) {
				console.warn('Could not fetch legal tags:', err);
			}

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

			searchResults = await api.searchDocuments(cleanParams, session);
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
			date_field_type: 'created_at',
			legal_tags: [],
			legal_tags_match_all: false, // Reset to OR behavior
			size: 10,
			sort_by: 'created_at',
			sort_order: 'desc',
			page: 1,
			use_fuzzy: false,
			include_highlights: true
		};
		performSearch();
	}
</script>

<div class="max-w-7/8 container mx-auto px-4 py-6" transition:fade>
	<!-- Header -->
	<div class="mb-6 rounded-xl bg-gradient-to-r from-primary-800 to-primary-900 p-6 shadow-lg">
		<div class="flex flex-col items-start justify-between md:flex-row md:items-center">
			<div class="mb-6 md:mb-0">
				<h1 class="mb-2 text-2xl font-bold text-white md:text-3xl">Motion Index</h1>
				<p class="text-sm text-primary-100 md:text-base">Search legal documents with precision</p>
			</div>

			{#if documentStats}
				<div class="grid w-full grid-cols-2 gap-3 sm:grid-cols-3 md:w-auto">
					<div class="rounded-lg bg-white/10 p-3 text-center backdrop-blur-sm">
						<p class="text-xs text-primary-100">Documents</p>
						<p class="text-xl font-semibold text-white">{documentStats.total_documents?.toLocaleString()}</p>
					</div>

					{#if documentStats.total_size_gb}
						<div class="rounded-lg bg-white/10 p-3 text-center backdrop-blur-sm">
							<p class="text-xs text-primary-100">Storage</p>
							<p class="text-xl font-semibold text-white">{documentStats.total_size_gb.toFixed(1)}GB</p>
						</div>
					{/if}

					{#if documentStats.date_range}
						<div class="rounded-lg bg-white/10 p-3 text-center backdrop-blur-sm">
							<p class="text-xs text-primary-100">Date Range</p>
							<p class="truncate text-sm font-medium text-white">
								{formatDate(documentStats.date_range.oldest)} - {formatDate(
									documentStats.date_range.newest
								)}
							</p>
						</div>
					{/if}

					{#if Object.keys(documentTypes).length > 0}
						<div class="rounded-lg bg-white/10 p-3 text-center backdrop-blur-sm">
							<p class="text-xs text-primary-100">Doc Types</p>
							<p class="text-xl font-semibold text-white">{Object.keys(documentTypes).length}</p>
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
				class={`flex-1 py-3 text-center font-medium transition-all ${activeTab === 'search' ? 'border-b-2 border-primary-900 bg-primary-50 text-primary-900' : 'text-neutral-600'}`}
				onclick={() => (activeTab = 'search')}
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
				class={`flex-1 py-3 text-center font-medium transition-all ${activeTab === 'results' ? 'border-b-2 border-primary-900 bg-primary-50 text-primary-900' : 'text-neutral-600'}`}
				onclick={() => (activeTab = 'results')}
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
				{legalTags}
				{documentStats}
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
	{supabase}
	{session}
	{user}
	on:close={() => {
		activeDocument = null;
		showDocumentPopup = false;
	}}
/>
