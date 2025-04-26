<script lang="ts">
  import { onMount } from 'svelte';
  import { format } from 'date-fns';
  import * as api from './api';

  // State variables
  let searchParams = {
    query: '',
    doc_type: '',
    case_number: '',
    case_name: '',
    judge: '',
    court: '',
    author: '',
    status: '',
    date_range: {
      start: '',
      end: ''
    },
    size: 10,
    sort_by: 'created_at',
    sort_order: 'desc' as 'asc' | 'desc',
    page: 1
  };
  
  let searchResults: api.SearchResponse = { total: 0, hits: [] };
  let isLoading = false;
  let error = '';
  let documentTypes: Record<string, number> = {};
  let metadataFields: api.MetadataField[] = [];
  let documentStats: any = null;
  
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
      Object.keys(cleanParams).forEach(key => {
        if (key !== 'date_range' && !cleanParams[key]) {
          delete cleanParams[key];
        }
      });
      
      // Clean up empty date range
      if (!cleanParams.date_range?.start && !cleanParams.date_range?.end) {
        delete cleanParams.date_range;
      }
      
      searchResults = await api.searchDocuments(cleanParams);
      totalPages = Math.ceil(searchResults.total / searchParams.size);
    } catch (err) {
      console.error('Search error:', err);
      error = 'An error occurred while searching. Please try again.';
      searchResults = { total: 0, hits: [] };
    } finally {
      isLoading = false;
    }
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
      judge: '',
      court: '',
      author: '',
      status: '',
      date_range: {
        start: '',
        end: ''
      },
      size: 10,
      sort_by: 'created_at',
      sort_order: 'desc',
      page: 1
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
</script>

<div class="container mx-auto px-4 py-8">
  <h1 class="text-3xl font-bold mb-6">Motion Index Search</h1>
  
  <!-- Stats Summary -->
  {#if documentStats}
    <div class="bg-gray-100 p-4 rounded-lg mb-6">
      <h2 class="text-lg font-semibold mb-2">Document Statistics</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <span class="font-medium">Total Documents:</span> {documentStats.total_documents}
        </div>
        {#if documentStats.date_range}
          <div>
            <span class="font-medium">Date Range:</span> 
            {formatDate(documentStats.date_range.oldest)} - {formatDate(documentStats.date_range.newest)}
          </div>
        {/if}
        <div>
          <span class="font-medium">Document Types:</span> {Object.keys(documentStats.document_types || {}).length}
        </div>
      </div>
    </div>
  {/if}
  
  <div class="flex flex-col lg:flex-row gap-6">
    <!-- Search Filters -->
    <div class="w-full lg:w-1/3">
      <div class="bg-white shadow rounded-lg p-4">
        <h2 class="text-xl font-semibold mb-4">Search Filters</h2>
        
        <form on:submit={handleSearch}>
          <!-- Text Search -->
          <div class="mb-4">
            <label for="query" class="block text-sm font-medium text-gray-700 mb-1">Text Search</label>
            <input 
              type="text" 
              id="query" 
              bind:value={searchParams.query} 
              placeholder="Search document content..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          
          <!-- Document Type -->
          <div class="mb-4">
            <label for="doc_type" class="block text-sm font-medium text-gray-700 mb-1">Document Type</label>
            <select 
              id="doc_type" 
              bind:value={searchParams.doc_type}
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">All Document Types</option>
              {#each Object.entries(documentTypes) as [type, count]}
                <option value={type}>{type} ({count})</option>
              {/each}
            </select>
          </div>
          
          <!-- Case Number -->
          <div class="mb-4">
            <label for="case_number" class="block text-sm font-medium text-gray-700 mb-1">Case Number</label>
            <input 
              type="text" 
              id="case_number" 
              bind:value={searchParams.case_number} 
              placeholder="Enter case number..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          
          <!-- Case Name -->
          <div class="mb-4">
            <label for="case_name" class="block text-sm font-medium text-gray-700 mb-1">Case Name</label>
            <input 
              type="text" 
              id="case_name" 
              bind:value={searchParams.case_name} 
              placeholder="Enter case name..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          
          <!-- Judge -->
          <div class="mb-4">
            <label for="judge" class="block text-sm font-medium text-gray-700 mb-1">Judge</label>
            <input 
              type="text" 
              id="judge" 
              bind:value={searchParams.judge} 
              placeholder="Enter judge name..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          
          <!-- Court -->
          <div class="mb-4">
            <label for="court" class="block text-sm font-medium text-gray-700 mb-1">Court</label>
            <input 
              type="text" 
              id="court" 
              bind:value={searchParams.court} 
              placeholder="Enter court..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          
          <!-- Author -->
          <div class="mb-4">
            <label for="author" class="block text-sm font-medium text-gray-700 mb-1">Author</label>
            <input 
              type="text" 
              id="author" 
              bind:value={searchParams.author} 
              placeholder="Enter author..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          
          <!-- Status -->
          <div class="mb-4">
            <label for="status" class="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <input 
              type="text" 
              id="status" 
              bind:value={searchParams.status} 
              placeholder="Enter status..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          
          <!-- Date Range -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">Date Range</label>
            <div class="grid grid-cols-2 gap-2">
              <div>
                <label for="date_start" class="block text-xs text-gray-500">From</label>
                <input 
                  type="date" 
                  id="date_start" 
                  bind:value={searchParams.date_range.start}
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                />
              </div>
              <div>
                <label for="date_end" class="block text-xs text-gray-500">To</label>
                <input 
                  type="date" 
                  id="date_end" 
                  bind:value={searchParams.date_range.end}
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                />
              </div>
            </div>
          </div>
          
          <!-- Sort Options -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">Sort By</label>
            <div class="grid grid-cols-2 gap-2">
              <select 
                bind:value={searchParams.sort_by}
                class="px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="created_at">Date</option>
                <option value="metadata.document_name">Document Name</option>
                <option value="doc_type">Document Type</option>
                <option value="metadata.case_number">Case Number</option>
                <option value="metadata.case_name">Case Name</option>
              </select>
              
              <select 
                bind:value={searchParams.sort_order}
                class="px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="desc">Descending</option>
                <option value="asc">Ascending</option>
              </select>
            </div>
          </div>
          
          <!-- Results Per Page -->
          <div class="mb-6">
            <label for="size" class="block text-sm font-medium text-gray-700 mb-1">Results Per Page</label>
            <select 
              id="size" 
              bind:value={searchParams.size}
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="10">10</option>
              <option value="25">25</option>
              <option value="50">50</option>
              <option value="100">100</option>
            </select>
          </div>
          
          <!-- Action Buttons -->
          <div class="flex gap-2">
            <button 
              type="submit" 
              class="flex-1 bg-indigo-600 text-white py-2 px-4 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
              disabled={isLoading}
            >
              {isLoading ? 'Searching...' : 'Search'}
            </button>
            
            <button 
              type="button" 
              on:click={resetFilters}
              class="flex-1 bg-gray-200 text-gray-800 py-2 px-4 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
              disabled={isLoading}
            >
              Reset
            </button>
          </div>
        </form>
      </div>
    </div>
    
    <!-- Search Results -->
    <div class="w-full lg:w-2/3">
      <div class="bg-white shadow rounded-lg p-4">
        <h2 class="text-xl font-semibold mb-4">Search Results</h2>
        
        {#if error}
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4" role="alert">
            <p>{error}</p>
          </div>
        {/if}
        
        <!-- Results Count -->
        <div class="mb-4 text-gray-600">
          {#if isLoading}
            Loading...
          {:else if searchResults.total === 0}
            No documents found.
          {:else}
            Found {searchResults.total} document{searchResults.total !== 1 ? 's' : ''}
            {#if searchParams.page > 1}
              (Page {searchParams.page} of {totalPages})
            {/if}
          {/if}
        </div>
        
        <!-- Results List -->
        {#if searchResults.hits.length > 0}
          <div class="space-y-4">
            {#each searchResults.hits as document}
              <div class="border rounded-lg p-4 hover:bg-gray-50">
                <div class="flex justify-between items-start">
                  <h3 class="text-lg font-medium text-indigo-600">{document.metadata.document_name || document.file_name}</h3>
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                    {document.doc_type}
                  </span>
                </div>
                
                <p class="mt-1 text-gray-500">{document.metadata.subject || 'No subject'}</p>
                
                {#if document.highlight?.text}
                  <div class="mt-2">
                    <p class="text-sm text-gray-700 font-medium">Matching content:</p>
                    {#each document.highlight.text as highlight}
                      <p class="text-sm text-gray-600 mt-1">...{@html highlight}...</p>
                    {/each}
                  </div>
                {:else}
                  <p class="mt-2 text-sm text-gray-600 line-clamp-2">{document.text.substring(0, 200)}...</p>
                {/if}
                
                <div class="mt-3 grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
                  {#if document.metadata.case_number}
                    <div>
                      <span class="font-medium">Case Number:</span> {document.metadata.case_number}
                    </div>
                  {/if}
                  
                  {#if document.metadata.case_name}
                    <div>
                      <span class="font-medium">Case Name:</span> {document.metadata.case_name}
                    </div>
                  {/if}
                  
                  {#if document.metadata.judge}
                    <div>
                      <span class="font-medium">Judge:</span> {document.metadata.judge}
                    </div>
                  {/if}
                  
                  {#if document.metadata.court}
                    <div>
                      <span class="font-medium">Court:</span> {document.metadata.court}
                    </div>
                  {/if}
                  
                  {#if document.metadata.author}
                    <div>
                      <span class="font-medium">Author:</span> {document.metadata.author}
                    </div>
                  {/if}
                  
                  {#if document.metadata.status}
                    <div>
                      <span class="font-medium">Status:</span> {document.metadata.status}
                    </div>
                  {/if}
                  
                  <div>
                    <span class="font-medium">Date:</span> {formatDate(document.created_at)}
                  </div>
                </div>
              </div>
            {/each}
          </div>
          
          <!-- Pagination -->
          {#if totalPages > 1}
            <div class="flex justify-center mt-6">
              <nav class="inline-flex rounded-md shadow">
                <button 
                  on:click={() => goToPage(searchParams.page - 1)}
                  disabled={searchParams.page === 1 || isLoading}
                  class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Previous
                </button>
                
                {#each Array(Math.min(5, totalPages)) as _, i}
                  {#if totalPages <= 5 || (i < 3 && searchParams.page <= 3) || (i >= 2 && searchParams.page > totalPages - 3)}
                    <button 
                      on:click={() => goToPage(i + 1)}
                      class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium 
                        ${searchParams.page === i + 1 ? 'bg-indigo-50 text-indigo-600 z-10' : 'text-gray-500 hover:bg-gray-50'}"
                    >
                      {i + 1}
                    </button>
                  {:else if i === 2 && searchParams.page > 3 && searchParams.page < totalPages - 2}
                    <button 
                      on:click={() => goToPage(searchParams.page)}
                      class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-indigo-50 text-indigo-600 text-sm font-medium z-10"
                    >
                      {searchParams.page}
                    </button>
                  {:else if (i === 1 && searchParams.page > 3) || (i === 3 && searchParams.page < totalPages - 2)}
                    <span class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700">
                      ...
                    </span>
                  {/if}
                {/each}
                
                <button 
                  on:click={() => goToPage(searchParams.page + 1)}
                  disabled={searchParams.page === totalPages || isLoading}
                  class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Next
                </button>
              </nav>
            </div>
          {/if}
        {/if}
      </div>
    </div>
  </div>
</div>
