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
  
  // Toggle search help
  function toggleSearchHelp() {
    showSearchHelp = !showSearchHelp;
  }
  
  // Toggle advanced filters
  function toggleAdvancedFilters() {
    showAdvancedFilters = !showAdvancedFilters;
  }
</script>

<div class="container mx-auto px-4 py-8 max-w-7xl">
  <!-- Modern Header with gradient background -->
  <div class="bg-gradient-to-r from-indigo-600 to-blue-500 rounded-xl shadow-lg p-6 mb-8">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center">
      <div class="mb-6 md:mb-0">
        <h1 class="text-3xl font-bold text-white mb-2">Motion Index Search</h1>
        <p class="text-indigo-100">Search and explore legal documents with advanced filtering</p>
      </div>
      
      {#if documentStats}
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-4 w-full md:w-auto">
          <div class="bg-white/10 backdrop-blur-sm border border-white/20 rounded-lg p-4 text-center transition-all hover:bg-white/20">
            <p class="text-sm text-indigo-100">Documents</p>
            <p class="text-2xl font-semibold text-white">{documentStats.total_documents}</p>
          </div>
          
          <div class="bg-white/10 backdrop-blur-sm border border-white/20 rounded-lg p-4 text-center transition-all hover:bg-white/20">
            <p class="text-sm text-indigo-100">Document Types</p>
            <p class="text-2xl font-semibold text-white">{Object.keys(documentTypes).length}</p>
          </div>
          
          {#if documentStats.date_range}
            <div class="bg-white/10 backdrop-blur-sm border border-white/20 rounded-lg p-4 text-center col-span-2 sm:col-span-1 transition-all hover:bg-white/20">
              <p class="text-sm text-indigo-100">Date Range</p>
              <span class="text-sm font-medium text-white">
                {formatDate(documentStats.date_range.oldest)} - {formatDate(documentStats.date_range.newest)}
              </span>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </div>
  
  <!-- Modern Mobile Tabs -->
  <div class="lg:hidden mb-6">
    <div class="flex rounded-full shadow-md bg-gray-100 p-1">
      <button 
        class={`flex-1 py-2 px-4 text-center font-medium rounded-full transition-all ${activeTab === 'search' ? 'bg-indigo-600 text-white shadow-md' : 'text-gray-700 hover:bg-gray-200'}`}
        on:click={() => activeTab = 'search'}
      >
        <div class="flex items-center justify-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          Search
        </div>
      </button>
      <button 
        class={`flex-1 py-2 px-4 text-center font-medium rounded-full transition-all ${activeTab === 'results' ? 'bg-indigo-600 text-white shadow-md' : 'text-gray-700 hover:bg-gray-200'}`}
        on:click={() => activeTab = 'results'}
      >
        <div class="flex items-center justify-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
          </svg>
          Results {searchResults.total > 0 ? `(${searchResults.total})` : ''}
        </div>
      </button>
    </div>
  </div>
  
  <div class="flex flex-col lg:flex-row gap-6">
    <!-- Search Filters -->
    <div class="w-full lg:w-1/3 {activeTab !== 'search' ? 'hidden lg:block' : ''}">
      <div class="bg-white shadow-lg rounded-xl p-5 overflow-hidden border border-gray-100">
        <!-- Search Help Toggle -->
        <div class="flex justify-between items-center mb-5">
          <h2 class="text-xl font-semibold text-gray-800">Search Filters</h2>
          <button 
            type="button" 
            class="inline-flex items-center px-3 py-2 border border-indigo-100 text-xs font-medium rounded-full text-indigo-700 bg-indigo-50 hover:bg-indigo-100 transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            on:click={toggleSearchHelp}
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Search Help
          </button>
        </div>
        
        <!-- Search Help Panel -->
        {#if showSearchHelp}
          <div class="bg-gradient-to-br from-indigo-50 to-blue-50 p-4 rounded-xl mb-5 text-sm border border-indigo-100 shadow-sm">
            <h3 class="font-semibold text-indigo-800 mb-3 flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-indigo-600" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
              </svg>
              Search Operators
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div class="bg-white/70 p-3 rounded-lg border border-indigo-100 flex items-center">
                <span class="font-mono text-indigo-700 px-2 py-1 rounded bg-indigo-50 mr-2">"exact phrase"</span>
                <span class="text-gray-700">Exact phrase match</span>
              </div>
              <div class="bg-white/70 p-3 rounded-lg border border-indigo-100 flex items-center">
                <span class="font-mono text-indigo-700 px-2 py-1 rounded bg-indigo-50 mr-2">term1 OR term2</span>
                <span class="text-gray-700">Either term</span>
              </div>
              <div class="bg-white/70 p-3 rounded-lg border border-indigo-100 flex items-center">
                <span class="font-mono text-indigo-700 px-2 py-1 rounded bg-indigo-50 mr-2">+term</span>
                <span class="text-gray-700">Must include</span>
              </div>
              <div class="bg-white/70 p-3 rounded-lg border border-indigo-100 flex items-center">
                <span class="font-mono text-indigo-700 px-2 py-1 rounded bg-indigo-50 mr-2">-term</span>
                <span class="text-gray-700">Must exclude</span>
              </div>
              <div class="bg-white/70 p-3 rounded-lg border border-indigo-100 flex items-center">
                <span class="font-mono text-indigo-700 px-2 py-1 rounded bg-indigo-50 mr-2">term*</span>
                <span class="text-gray-700">Wildcard search</span>
              </div>
              <div class="bg-white/70 p-3 rounded-lg border border-indigo-100 flex items-center">
                <span class="font-mono text-indigo-700 px-2 py-1 rounded bg-indigo-50 mr-2">term~</span>
                <span class="text-gray-700">Fuzzy search</span>
              </div>
            </div>
            <div class="mt-3 bg-blue-50 p-3 rounded-lg text-blue-800 font-medium">
              By default, all terms must match (AND logic)
            </div>
          </div>
        {/if}
        <form on:submit={handleSearch} class="space-y-5">
          <!-- Text Search -->
          <div class="relative">
            <label for="query" class="block text-sm font-medium text-gray-700 mb-2">Text Search</label>
            <div class="mt-1 flex rounded-md shadow-lg">
              <div class="relative flex-grow focus-within:z-10">
                <input 
                  type="text" 
                  id="query" 
                  bind:value={searchParams.query}
                  placeholder="Enter search terms..."
                  class="block w-full rounded-l-lg border-gray-200 focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm py-3 px-4 shadow-sm"
                />
              </div>
              <button 
                type="submit" 
                class="relative inline-flex items-center space-x-2 rounded-r-lg border border-transparent bg-indigo-600 px-5 py-3 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-colors duration-200"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </button>
            </div>
          </div>
          
          <!-- Basic Filters -->
          <div class="space-y-4">
            <!-- Document Type -->
            <div>
              <label for="doc_type" class="block text-sm font-medium text-gray-700 mb-1">Document Type</label>
              <select 
                id="doc_type" 
                bind:value={searchParams.doc_type}
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
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
            
            <!-- Case Number -->
            <div>
              <label for="case_number" class="block text-sm font-medium text-gray-700 mb-1">Case Number</label>
              <select 
                id="case_number" 
                bind:value={searchParams.case_number}
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
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
          <div>
            <button 
              type="button" 
              on:click={toggleAdvancedFilters}
              class="text-indigo-600 hover:text-indigo-800 text-sm flex items-center"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
              {showAdvancedFilters ? 'Hide Advanced Filters' : 'Show Advanced Filters'}
            </button>
          </div>
          
          <!-- Advanced Filters -->
          {#if showAdvancedFilters}
            <div class="space-y-4 border-t border-gray-200 pt-4 mt-2">
              <!-- Case Name -->
              <div>
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
              <div>
                <label for="judge" class="block text-sm font-medium text-gray-700 mb-1">Judge</label>
                <select 
                  id="judge" 
                  bind:value={searchParams.judge}
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                >
                  <option value="">All Judges</option>
                  {#if fieldOptions.judge}
                    {#each fieldOptions.judge as judge}
                      <option value={judge}>{judge}</option>
                    {/each}
                  {/if}
                </select>
              </div>
              
              <!-- Court -->
              <div>
                <label for="court" class="block text-sm font-medium text-gray-700 mb-1">Court</label>
                <select 
                  id="court" 
                  bind:value={searchParams.court}
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                >
                  <option value="">All Courts</option>
                  {#if fieldOptions.court}
                    {#each fieldOptions.court as court}
                      <option value={court}>{court}</option>
                    {/each}
                  {/if}
                </select>
              </div>
              
              <!-- Author -->
              <div>
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
              <div>
                <label for="status" class="block text-sm font-medium text-gray-700 mb-1">Status</label>
                <select 
                  id="status" 
                  bind:value={searchParams.status}
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                >
                  <option value="">All Statuses</option>
                  {#if fieldOptions.status}
                    {#each fieldOptions.status as status}
                      <option value={status}>{status}</option>
                    {/each}
                  {/if}
                </select>
              </div>
              
              <!-- Date Range -->
              <div>
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
            </div>
          {/if}
          
          <!-- Sort and Pagination Options -->
          <div class="border-t border-gray-200 pt-4 mt-2 space-y-4">
            <!-- Sort Options -->
            <div>
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
            <div>
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
          </div>
          
          <!-- Action Buttons -->
          <div class="flex gap-2 pt-2">
            <button 
              type="submit" 
              class="flex-1 bg-indigo-600 text-white py-2 px-4 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 flex items-center justify-center"
              disabled={isLoading}
            >
              {#if isLoading}
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Searching...
              {:else}
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                Search
              {/if}
            </button>
            
            <button 
              type="button" 
              on:click={resetFilters}
              class="flex-1 bg-gray-200 text-gray-800 py-2 px-4 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 flex items-center justify-center"
              disabled={isLoading}
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Reset
            </button>
          </div>
        </form>
      </div>
    </div>
    
    <!-- Search Results -->
    <div class="w-full lg:w-2/3 {activeTab !== 'results' ? 'hidden lg:block' : ''}">
      <div class="bg-white shadow rounded-lg p-4 overflow-hidden">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold text-gray-800">Search Results</h2>
          
          <!-- Results Count -->
          <div class="text-gray-600 text-sm font-medium bg-gray-100 px-3 py-1 rounded-full">
            {#if isLoading}
              <div class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Loading...
              </div>
            {:else if searchResults.total === 0}
              No documents found
            {:else}
              {searchResults.total} document{searchResults.total !== 1 ? 's' : ''}
              {#if searchParams.page > 1}
                (Page {searchParams.page} of {totalPages})
              {/if}
            {/if}
          </div>
        </div>
        
        {#if error}
          <div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4" role="alert">
            <div class="flex items-start">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-red-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm">{error}</p>
              </div>
            </div>
          </div>
        {/if}
        
        <!-- Results List -->
        {#if searchResults.hits.length > 0}
          <div class="space-y-6">
            {#each searchResults.hits as document}
              <div class="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition duration-150 ease-in-out shadow-sm">
                <div class="flex flex-wrap justify-between items-start gap-2 mb-2">
                  <h3 class="text-lg font-medium text-indigo-700">{document.metadata.document_name || document.file_name}</h3>
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
                    {document.doc_type}
                  </span>
                </div>
                
                <p class="text-gray-600">{document.metadata.subject || 'No subject'}</p>
                
                {#if document.highlight?.text}
                  <div class="mt-3 bg-yellow-50 p-3 rounded-md border border-yellow-100">
                    <p class="text-xs uppercase tracking-wide text-gray-500 font-semibold mb-1">Matching content</p>
                    {#each document.highlight.text as highlight}
                      <p class="text-sm text-gray-700 mt-1">...{@html highlight}...</p>
                    {/each}
                  </div>
                {:else}
                  <p class="mt-3 text-sm text-gray-600 line-clamp-2">{document.text.substring(0, 200)}...</p>
                {/if}
                
                <div class="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-2 text-sm">
                  {#if document.metadata.case_number}
                    <div class="flex items-center">
                      <span class="w-24 flex-shrink-0 text-gray-500">Case Number:</span> 
                      <span class="font-medium text-gray-900">{document.metadata.case_number}</span>
                    </div>
                  {/if}
                  
                  {#if document.metadata.case_name}
                    <div class="flex items-center">
                      <span class="w-24 flex-shrink-0 text-gray-500">Case Name:</span> 
                      <span class="font-medium text-gray-900">{document.metadata.case_name}</span>
                    </div>
                  {/if}
                  
                  {#if document.metadata.judge}
                    <div class="flex items-center">
                      <span class="w-24 flex-shrink-0 text-gray-500">Judge:</span> 
                      <span class="font-medium text-gray-900">{document.metadata.judge}</span>
                    </div>
                  {/if}
                  
                  {#if document.metadata.court}
                    <div class="flex items-center">
                      <span class="w-24 flex-shrink-0 text-gray-500">Court:</span> 
                      <span class="font-medium text-gray-900">{document.metadata.court}</span>
                    </div>
                  {/if}
                  
                  {#if document.metadata.author}
                    <div class="flex items-center">
                      <span class="w-24 flex-shrink-0 text-gray-500">Author:</span> 
                      <span class="font-medium text-gray-900">{document.metadata.author}</span>
                    </div>
                  {/if}
                  
                  {#if document.metadata.status}
                    <div class="flex items-center">
                      <span class="w-24 flex-shrink-0 text-gray-500">Status:</span> 
                      <span class="font-medium text-gray-900">{document.metadata.status}</span>
                    </div>
                  {/if}
                  
                  <div class="flex items-center">
                    <span class="w-24 flex-shrink-0 text-gray-500">Date:</span> 
                    <span class="font-medium text-gray-900">{formatDate(document.created_at)}</span>
                  </div>
                </div>
              </div>
            {/each}
          </div>
          
          <!-- Pagination -->
          {#if totalPages > 1}
            <div class="flex justify-center mt-8">
              <nav class="inline-flex rounded-md shadow-sm" aria-label="Pagination">
                <button 
                  on:click={() => goToPage(searchParams.page - 1)}
                  disabled={searchParams.page === 1 || isLoading}
                  class="relative inline-flex items-center px-3 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                </button>
                
                {#each Array(Math.min(5, totalPages)) as _, i}
                  {#if totalPages <= 5 || (i < 3 && searchParams.page <= 3) || (i >= 2 && searchParams.page > totalPages - 3)}
                    <button 
                      on:click={() => goToPage(i + 1)}
                      class={`relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium ${searchParams.page === i + 1 ? 'bg-indigo-50 text-indigo-600 z-10 border-indigo-500' : 'text-gray-500 hover:bg-gray-50'}`}
                    >
                      {i + 1}
                    </button>
                  {:else if i === 2 && searchParams.page > 3 && searchParams.page < totalPages - 2}
                    <button 
                      on:click={() => goToPage(searchParams.page)}
                      class="relative inline-flex items-center px-4 py-2 border border-indigo-500 bg-indigo-50 text-indigo-600 text-sm font-medium z-10"
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
                  class="relative inline-flex items-center px-3 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                  </svg>
                </button>
              </nav>
            </div>
          {/if}
        {/if}
      </div>
    </div>
  </div>
</div>
