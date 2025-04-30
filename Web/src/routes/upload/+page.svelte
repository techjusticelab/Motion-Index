<script lang="ts">
	import { categoriseDocument } from '../api';

	let selectedFile: File | null = null;
	let uploadStatus: string = '';
	let documentResponse: any = null;
	let fileInputLabel: string = 'Select a file'; // Added to show current selection

	async function handleFileUpload() {
		console.log('Selected file:', selectedFile);
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
		console.log('File input changed:', target.files);
		if (target.files && target.files.length > 0) {
			selectedFile = target.files[0];
			fileInputLabel = selectedFile.name; // Update label with filename
			console.log('File selected:', selectedFile.name);
		} else {
			selectedFile = null;
			fileInputLabel = 'Select a file';
			console.log('No file selected');
		}
	}
</script>

<div class="upload-container">
	<h1 class="text-xl font-bold">Upload and Categorise Document</h1>

	<div class="file-input-container mt-4">
		<!-- Custom styled file input -->
		<label for="file-upload" class="file-input-label">
			<span class="file-input-text">{fileInputLabel}</span>
			<span class="file-input-button">Browse</span>
		</label>
		<input
			id="file-upload"
			type="file"
			accept=".pdf,.docx,.txt"
			on:change={handleFileChange}
			class="hidden-file-input"
		/>
	</div>

	<div class="mt-4">
		<button
			class="rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600"
			on:click={handleFileUpload}
			disabled={!selectedFile}
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
			<pre class="rounded bg-gray-100 p-4">{JSON.stringify(documentResponse, null, 2)}</pre>
		</div>
	{/if}
</div>

<style>
	.upload-container {
		max-width: 600px;
		margin: 0 auto;
		padding: 20px;
	}

	.file-input-container {
		position: relative;
		margin-bottom: 15px;
	}

	.hidden-file-input {
		position: absolute;
		width: 0.1px;
		height: 0.1px;
		opacity: 0;
		overflow: hidden;
		z-index: -1;
	}

	.file-input-label {
		display: flex;
		cursor: pointer;
		width: 100%;
		border: 1px solid #ccc;
		border-radius: 4px;
		overflow: hidden;
	}

	.file-input-text {
		flex-grow: 1;
		padding: 8px 12px;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		color: #555;
	}

	.file-input-button {
		background-color: #e9e9e9;
		padding: 8px 16px;
		color: #333;
		font-weight: 500;
		border-left: 1px solid #ccc;
	}

	.file-input-label:hover .file-input-button {
		background-color: #ddd;
	}
</style>
