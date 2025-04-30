<script lang="ts">
	import { categoriseDocument } from '../api';
	import { onMount } from 'svelte';

	let selectedFile: File | null = null;
	let uploadStatus: string = '';
	let documentResponse: any = null;

	async function handleFileUpload() {
		if (!selectedFile) {
			uploadStatus = 'Please select a file to upload.';
			return;
		}

		uploadStatus = 'Uploading and categorising document...';

		try {
			const response = await categoriseDocument(selectedFile);
			documentResponse = response.document;
			uploadStatus = 'Document categorised successfully!';
		} catch (error) {
			uploadStatus = 'Failed to categorise document. Please try again.';
			console.error(error);
		}
	}

	function handleFileChange(event: Event) {
		const target = event.target as HTMLInputElement;
		if (target.files && target.files.length > 0) {
			selectedFile = target.files[0];
		}
	}
</script>

<div class="upload-container">
	<h1 class="text-xl font-bold">Upload and Categorise Document</h1>

	<div class="mt-4">
		<input type="file" accept=".pdf,.docx,.txt" on:change={handleFileChange} />
	</div>

	<div class="mt-4">
		<button
			class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
			on:click={handleFileUpload}
		>
			Upload and Categorise
		</button>
	</div>

	{#if uploadStatus}
		<p class="mt-4 text-gray-700">{uploadStatus}</p>
	{/if}

	{#if documentResponse}
		<div class="mt-6">
			<h2 class="text-lg font-semibold">Document Details</h2>
			<pre class="bg-gray-100 p-4 rounded">{JSON.stringify(documentResponse, null, 2)}</pre>
		</div>
	{/if}
</div>

<style>
	.upload-container {
		max-width: 600px;
		margin: 0 auto;
		padding: 20px;
	}
</style>
