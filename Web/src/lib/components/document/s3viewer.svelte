<!-- LocalFileViewer.svelte -->
<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	// Event dispatcher
	const dispatch = createEventDispatcher<{
		loaded: void;
		error: { message: string };
	}>();

	// Props
	const {
		fileUrl,
		width = '100%',
		height = '500px'
	}: {
		fileUrl: string;
		width?: string;
		height?: string;
	} = $props();

	// State
	let isLoading = $state(true);
	let error = $state<string | null>(null);

	// React to changes in fileUrl
	$effect(() => {
		try {
			console.log('File URL from local server:', fileUrl);
			if (!fileUrl) {
				// Reset state when no URL
				isLoading = true;
				error = null;
				return;
			}
			// URL is ready, let iframe handle loading state
			error = null;
		} catch (err) {
			isLoading = false;
			error = err instanceof Error ? err.message : 'Unknown error occurred';
			console.error('Error loading file:', err);
			dispatch('error', { message: error });
		}
	});

	// Handle iframe load event
	function handleLoad(event: Event) {
		console.log('Iframe loaded successfully:', event);
		isLoading = false;
		dispatch('loaded');
	}

	// Handle iframe error event
	function handleError(event: Event) {
		console.error('Iframe failed to load:', event);
		isLoading = false;
		error = "Failed to display file. Browser security restrictions may prevent viewing localhost URLs in iframes. Try opening in a new tab.";
		dispatch('error', { message: error });
	}
</script>

<div class="s3-viewer h-full w-full">
	{#if isLoading}
		<div class="loading">
			<div class="spinner"></div>
			<p>Loading document...</p>
		</div>
	{:else if error}
		<div class="error">
			<p>Error: {error}</p>
			<div class="mt-4 flex gap-2 justify-center">
				<a
					href={fileUrl}
					target="_blank"
					rel="noopener noreferrer"
					class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700"
				>
					Open in New Tab
				</a>
				<a
					href={fileUrl}
					download
					class="inline-flex items-center px-4 py-2 text-sm font-medium text-blue-600 bg-white border border-blue-600 rounded-md hover:bg-blue-50"
				>
					Download File
				</a>
			</div>
		</div>
	{:else if fileUrl}
		<!-- Try iframe first for PDF display -->
		<iframe
			src={fileUrl}
			title="Document Viewer"
			{width}
			{height}
			on:load={handleLoad}
			on:error={handleError}
			class="h-full w-full"
			allow="same-origin"
			sandbox="allow-same-origin allow-scripts"
		></iframe>

		<!-- Fallback object tag if iframe fails (commented out for now) -->
		<!--
		<object
			data={fileUrl}
			type="application/pdf"
			{width}
			{height}
			class="h-full w-full"
		>
			<p>Unable to display PDF file. <a href={fileUrl} target="_blank">Click here to download</a></p>
		</object>
		-->
	{/if}
</div>

<style>
	.s3-viewer {
		position: relative;
		overflow: hidden;
		background-color: #f9f9f9;
	}

	iframe {
		border: none;
		display: block;
		background-color: white;
	}

	.loading,
	.error {
		padding: 2rem;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 200px;
		height: 100%;
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 4px solid rgba(0, 0, 0, 0.1);
		border-radius: 50%;
		border-top-color: #3498db;
		animation: spin 1s linear infinite;
		margin-bottom: 1rem;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.error {
		color: #e74c3c;
	}

	.error p {
		margin: 0;
		text-align: center;
	}
</style>
