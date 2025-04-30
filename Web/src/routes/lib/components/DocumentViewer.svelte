<script lang="ts">
  import { onMount } from 'svelte';
  import { S3Client, GetObjectCommand } from "@aws-sdk/client-s3";
  import { getSignedUrl } from "@aws-sdk/s3-request-presigner";

  // Props
  export let document: { 
    metadata: { document_name: string, [key: string]: any },
    file_name: string, 
    s3_uri: string,
    text?: string,
    doc_type?: string,
    created_at?: string
  } | null = null;
  
  export let isOpen = false;

  // State
  let documentUrl = '';
  let isLoadingDocument = false;
  let documentError = '';
  let showMetadata = false;
  let s3Client: S3Client | null = null;
  
  // Initialize the S3 client - using the correct region for the bucket
  onMount(() => {
    try {
      s3Client = new S3Client({ 
        region: 'us-east-2', // Hardcoded to the correct region for cpda-documents bucket
        credentials: {
          accessKeyId: import.meta.env.VITE_AWS_ACCESS_KEY_ID,
          secretAccessKey: import.meta.env.VITE_AWS_SECRET_ACCESS_KEY,
        },
        forcePathStyle: false // Use virtual-hosted style URLs
      });
      console.log('S3 client initialized successfully');
    } catch (error) {
      console.error('Error initializing S3 client:', error);
    }
  });
  
  // Function to get signed URL from S3
  async function getSignedS3Url(s3Uri: string, expiresIn = 3600): Promise<string> {
    try {
      if (!s3Client) {
        throw new Error('S3 client not initialized');
      }
      
      // Parse s3:// URI to get bucket and key
      const s3UriRegex = /s3:\/\/([^\/]+)\/(.+)/;
      const match = s3Uri.match(s3UriRegex);
      
      if (!match) {
        throw new Error(`Invalid S3 URI format: ${s3Uri}`);
      }
      
      const [, bucket, key] = match;
      
      // For the specific bucket we know requires us-east-2
      if (bucket === 'cpda-documents') {
        // Construct direct URL to avoid redirect issues
        return `https://${bucket}.s3.us-east-2.amazonaws.com/${key}`;
      }
      
      // For other buckets, use signed URL
      const command = new GetObjectCommand({
        Bucket: bucket,
        Key: key,
      });
      
      // Generate a signed URL that expires in 1 hour
      const signedUrl = await getSignedUrl(s3Client, command, { expiresIn });
      return signedUrl;
    } catch (error) {
      console.error("Error generating signed URL:", error);
      throw error;
    }
  }

  // Close document popup
  function closeDocumentPopup() {
    document = null;
    isOpen = false;
    documentUrl = '';
    // Restore body scrolling
    window.document.body.style.overflow = 'auto';
  }

  // Toggle metadata display
  function toggleMetadata() {
    showMetadata = !showMetadata;
  }

  // Watch for document changes and load the document when it changes
  $: if (document && isOpen) {
    loadDocument(document);
  }

  // Load document from S3
  async function loadDocument(doc) {
    if (!doc) return;
    
    // Prevent body scrolling when popup is open
    window.document.body.style.overflow = 'hidden';
    
    documentUrl = '';
    isLoadingDocument = true;
    documentError = '';
    
    // Check if there's an S3 URI to process and S3 client is initialized
    if (doc.s3_uri) {
      try {
        // If it's from our known bucket, construct URL directly
        if (doc.s3_uri.includes('cpda-documents')) {
          const s3UriRegex = /s3:\/\/([^\/]+)\/(.+)/;
          const match = doc.s3_uri.match(s3UriRegex);
          
          if (match) {
            const [, bucket, key] = match;
            documentUrl = `https://${bucket}.s3.us-east-2.amazonaws.com/${key}`;
          } else {
            throw new Error(`Invalid S3 URI format: ${doc.s3_uri}`);
          }
        } else if (s3Client) {
          // For other buckets, use the signed URL approach
          documentUrl = await getSignedS3Url(doc.s3_uri);
        } else {
          throw new Error('S3 client not initialized');
        }
      } catch (err) {
        console.error('Error getting document URL:', err);
        documentError = 'Could not access document from S3. Showing text content if available.';
      } finally {
        isLoadingDocument = false;
      }
    } else {
      isLoadingDocument = false;
    }
  }

  // Format date (you can import the format function from date-fns if needed)
  function formatDate(dateString: string): string {
    if (!dateString) return 'N/A';
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });
    } catch (err) {
      return dateString;
    }
  }
</script>

<!-- Document Popup -->
{#if document && isOpen}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-75 p-4"
    on:click={closeDocumentPopup}
  >
    <div
      class="relative h-[85vh] w-full max-w-5xl rounded-lg bg-white shadow-2xl"
      on:click|stopPropagation
    >
      <!-- Popup Header -->
      <div class="flex items-center justify-between border-b border-gray-200 p-4">
        <div>
          <h2 class="truncate text-lg font-semibold text-gray-800">
            {document.metadata.document_name || document.file_name}
          </h2>
          {#if document.doc_type}
            <span class="text-sm text-gray-600">{document.doc_type}</span>
          {/if}
        </div>
        <div class="flex items-center space-x-2">
          <button
            class="rounded-lg bg-blue-50 px-3 py-1 text-sm font-medium text-blue-600 hover:bg-blue-100"
            on:click|stopPropagation={toggleMetadata}
          >
            {showMetadata ? 'Hide Metadata' : 'Show Metadata'}
          </button>
          <button
            class="rounded-full p-2 text-gray-500 transition-colors hover:bg-gray-100"
            on:click={closeDocumentPopup}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-6 w-6"
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
      </div>

      <!-- Metadata Panel (conditionally shown) -->
      {#if showMetadata && document}
        <div class="border-b border-gray-200 bg-gray-50 p-4">
          <h3 class="mb-2 text-sm font-medium text-gray-700">Document Metadata</h3>
          <div class="grid grid-cols-1 gap-3 sm:grid-cols-2 md:grid-cols-3">
            <!-- Document type -->
            {#if document.doc_type}
              <div class="rounded border border-gray-200 bg-white p-2 shadow-sm">
                <span class="block text-xs text-gray-500">Document Type</span>
                <span class="font-medium text-gray-800">{document.doc_type}</span>
              </div>
            {/if}
            
            <!-- Creation date -->
            {#if document.created_at}
              <div class="rounded border border-gray-200 bg-white p-2 shadow-sm">
                <span class="block text-xs text-gray-500">Created</span>
                <span class="font-medium text-gray-800">{formatDate(document.created_at)}</span>
              </div>
            {/if}
            
            <!-- S3 URI -->
            {#if document.s3_uri}
              <div class="rounded border border-gray-200 bg-white p-2 shadow-sm">
                <span class="block text-xs text-gray-500">Storage Location</span>
                <span class="block truncate font-medium text-gray-800">{document.s3_uri}</span>
              </div>
            {/if}
            
            <!-- Dynamic metadata fields -->
            {#if document.metadata}
              {#each Object.entries(document.metadata) as [key, value]}
                {#if key !== 'document_name' && value !== null && value !== undefined && value !== ''}
                  <div class="rounded border border-gray-200 bg-white p-2 shadow-sm">
                    <span class="block text-xs text-gray-500">{key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())}</span>
                    <span class="block truncate font-medium text-gray-800">{value}</span>
                  </div>
                {/if}
              {/each}
            {/if}
          </div>
        </div>
      {/if}

      <!-- Document Content -->
      <div class="h-[calc(100%-4rem-1px)] w-full {showMetadata ? 'h-[calc(100%-13rem-2px)]' : ''}">
        {#if isLoadingDocument}
          <div class="flex h-full flex-col items-center justify-center">
            <svg
              class="h-8 w-8 animate-spin text-blue-600"
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
            <p class="mt-3 text-sm text-gray-600">Loading document...</p>
          </div>
        {:else if documentUrl}
          <iframe
            src={documentUrl}
            title={document.metadata.document_name || document.file_name}
            class="h-full w-full rounded-b-lg"
            sandbox="allow-same-origin allow-scripts allow-forms"
            loading="lazy"
          ></iframe>
        {:else if document.text}
          <!-- Fallback to text content when no PDF is available -->
          <div class="h-full overflow-auto p-6 text-gray-800">
            {#if documentError}
              <div class="mb-4 rounded-md bg-yellow-50 p-3 text-sm text-yellow-700">
                <p>{documentError}</p>
              </div>
            {/if}
            <div class="mx-auto max-w-4xl whitespace-pre-wrap font-serif text-base leading-relaxed">
              {document.text}
            </div>
          </div>
        {:else if documentError}
          <div class="flex h-full flex-col items-center justify-center p-8 text-center">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="mb-4 h-16 w-16 text-yellow-300"
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
            <p class="mb-2 text-gray-600">{documentError}</p>
            <p class="text-sm text-gray-500">
              No text content available to display.
            </p>
          </div>
        {:else}
          <div class="flex h-full flex-col items-center justify-center p-8 text-center">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="mb-4 h-16 w-16 text-gray-300"
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
            <p class="mb-2 text-gray-600">Document preview not available</p>
            <p class="text-sm text-gray-500">
              No content is available to display for this document.
            </p>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}