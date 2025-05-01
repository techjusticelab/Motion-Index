<!-- S3FileViewer.svelte -->
<script lang="ts">
	import { onMount, createEventDispatcher } from 'svelte';
	import { getSignedS3Url } from '../../services/s3';

	// Event dispatcher
	const dispatch = createEventDispatcher<{
		loaded: void;
		error: { message: string };
	}>();

	// Props
	export let s3Uri: string; // S3 URI in the format s3://bucket-name/object-key
	export let width: string = '100%';
	export let height: string = '500px';
	export let urlExpirationSeconds: number = 3600;

	// State
	let fileUrl: string = '';
	let isLoading: boolean = true;
	let error: string | null = null;

	onMount(async () => {
		try {
			console.log('S3 URI:', s3Uri);
			if (!s3Uri) {
				throw new Error('S3 URI is required');
			}

			isLoading = true;
			fileUrl = await getSignedS3Url(s3Uri, urlExpirationSeconds);
			console.log('File URL:', fileUrl);
			isLoading = false;
		} catch (err) {
			isLoading = false;
			error = err instanceof Error ? err.message : 'Unknown error occurred';
			console.error('Error loading file:', err);
			dispatch('error', { message: error });
		}
	});

	// Handle iframe load event
	function handleLoad() {
		isLoading = false;
		dispatch('loaded');
	}

	// Handle iframe error event
	function handleError() {
		isLoading = false;
		error = "Failed to display file. It may not be viewable in an iframe or doesn't exist.";
		dispatch('error', { message: error });
	}
</script>

<div class="s3-viewer h-full w-full">
	{#if isLoading}
		<div class="loading">
			<div class="spinner"></div>
			<p>Loading file from S3...</p>
		</div>
	{:else if error}
		<div class="error">
			<p>Error: {error}</p>
		</div>
	{:else if fileUrl}
		<iframe
			src={fileUrl}
			title="S3 File Viewer"
			{width}
			{height}
			on:load={handleLoad}
			on:error={handleError}
			class="h-full w-full"
		></iframe>
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
