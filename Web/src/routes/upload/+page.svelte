<script lang="ts">
	import { categoriseDocument, updateDocumentMetadata, type Document } from '../api';
	import { onMount } from 'svelte';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { cubicOut, quintOut, elasticOut, backOut } from 'svelte/easing';

	let selectedFile: File | null = null;
	let uploadStatus: string = '';
	let documentResponse: Document | null = null;
	let fileInputLabel: string = 'Drag and drop your file here or click to browse';
	let isDragging: boolean = false;
	let uploadedDocuments: Array<{ name: string; type: string; response: Document }> = [];
	let isUploading: boolean = false;
	let isUpdatingMetadata: boolean = false; // New state for metadata update loading
	let metadataUpdateSuccess: boolean | null = null; // Track success/failure state

	// Metadata panel state
	let currentMetadata: Document['metadata'] = {
		document_name: '',
		subject: '',
		status: '',
		timestamp: '',
		case_name: '',
		case_number: '',
		author: '',
		judge: '',
		court: '',
		legal_tags: []
	};
	let legalTags: string[] = [];
	let tagInput = '';
	let currentDocumentId = '';

	// Timer to reset success/failure message
	let successMessageTimer: ReturnType<typeof setTimeout> | null = null;

	// Function to handle tag input
	function addTag() {
		if (tagInput.trim() && !legalTags.includes(tagInput.trim())) {
			legalTags = [...legalTags, tagInput.trim()];
			tagInput = '';
		}
	}

	// Function to remove a tag
	function removeTag(index: number) {
		legalTags = legalTags.filter((_, i) => i !== index);
	}

	// Function to handle keydown events in tag input
	function handleTagKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			event.preventDefault();
			addTag();
		}
	}

	// Function to clear the success/error message after a delay
	function resetMetadataUpdateStatus() {
		if (successMessageTimer) clearTimeout(successMessageTimer);

		successMessageTimer = setTimeout(() => {
			metadataUpdateSuccess = null;
		}, 5000); // Message disappears after 5 seconds
	}

	async function saveMetadata() {
		if (!currentDocumentId) {
			uploadStatus = 'Error: No document selected for metadata update';
			return;
		}

		try {
			isUpdatingMetadata = true;

			// Create the complete metadata object including tags
			const metadataToUpdate = {
				...currentMetadata,
				legal_tags: legalTags
			};

			// Call the API to update metadata
			const response = await updateDocumentMetadata(currentDocumentId, metadataToUpdate);

			// Update the document in our list
			uploadedDocuments = uploadedDocuments.map((doc) => {
				if (doc.response.id === currentDocumentId) {
					const updatedDoc = {
						...doc.response,
						metadata: metadataToUpdate
					};
					return {
						...doc,
						response: updatedDoc
					};
				}
				return doc;
			});

			// Update the current document response if it matches
			if (documentResponse && documentResponse.id === currentDocumentId) {
				documentResponse = {
					...documentResponse,
					metadata: metadataToUpdate
				};
			}

			// Show success message
			metadataUpdateSuccess = true;
			uploadStatus = 'Document metadata updated successfully!';
			
			// Set timer to clear message
			resetMetadataUpdateStatus();
		} catch (error) {
			metadataUpdateSuccess = false;
			uploadStatus = 'Failed to update document metadata. Please try again.';
			console.error(error);

			// Set timer to clear message
			resetMetadataUpdateStatus();
		} finally {
			isUpdatingMetadata = false;
		}
	}

	// Function to initialize metadata form with document data
	function initializeMetadataForm(document: Document) {
		if (!document || !document.metadata) return;

		// Set the document ID for the API call
		currentDocumentId = document.id;

		// Initialize all fields (even empty ones)
		currentMetadata = {
			document_name: document.metadata.document_name || '',
			subject: document.metadata.subject || '',
			status: document.metadata.status || '',
			timestamp: document.metadata.timestamp || '',
			case_name: document.metadata.case_name || '',
			case_number: document.metadata.case_number || '',
			author: document.metadata.author || '',
			judge: document.metadata.judge || '',
			court: document.metadata.court || '',
			legal_tags: []
		};

		// Set the legal tags
		legalTags = [...(document.metadata.legal_tags || [])];
	}

	// Show document details
	function showDocumentDetails(document: Document) {
		documentResponse = document;
		initializeMetadataForm(document);
	}

	async function handleFileUpload() {
		console.log('Selected file:', selectedFile);
		if (!selectedFile) {
			uploadStatus = 'Please select a file to upload.';
			return;
		}

		isUploading = true;
		uploadStatus = 'Uploading and categorising document...';

		try {
			// Call the actual API endpoint
			const response = await categoriseDocument(selectedFile);

			// Use the full document response
			documentResponse = response.document;

			// Initialize the metadata form
			initializeMetadataForm(response.document);

			// Add the document to our list of uploaded documents
			uploadedDocuments = [
				...uploadedDocuments,
				{
					name: selectedFile.name,
					type: selectedFile.type,
					response: response.document
				}
			];

			// Reset the file input
			selectedFile = null;
			fileInputLabel = 'Drag and drop your file here or click to browse';

			uploadStatus = 'Document categorised successfully!';
		} catch (error) {
			uploadStatus = 'Failed to categorise document. Please try again.';
			console.error(error);
		} finally {
			isUploading = false;
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
			fileInputLabel = 'Drag and drop your file here or click to browse';
			console.log('No file selected');
		}
	}

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
		isDragging = true;
	}

	function handleDragLeave(event: DragEvent) {
		event.preventDefault();
		isDragging = false;
	}

	function handleDrop(event: DragEvent) {
		event.preventDefault();
		isDragging = false;

		const files = event.dataTransfer?.files;
		if (files && files.length > 0) {
			selectedFile = files[0];
			fileInputLabel = selectedFile.name;
			console.log('File dropped:', selectedFile.name);
		}
	}

	function getFileIcon(type: string) {
		if (type.includes('pdf')) {
			return 'pdf';
		} else if (type.includes('word') || type.includes('docx') || type.includes('doc')) {
			return 'word';
		} else if (type.includes('text') || type.includes('txt')) {
			return 'text';
		} else {
			return 'generic';
		}
	}

	function getFileExtension(filename: string) {
		return filename.split('.').pop()?.toLowerCase() || '';
	}

	function getFileIconByName(filename: string) {
		const ext = getFileExtension(filename);
		if (ext === 'pdf') {
			return 'pdf';
		} else if (['doc', 'docx'].includes(ext)) {
			return 'word';
		} else if (ext === 'txt') {
			return 'text';
		} else {
			return 'generic';
		}
	}
</script>

<div class="flex min-h-[80vh] items-center justify-center p-4">
	<!-- Main container with responsive layout -->
	<div
		class="w-full max-w-7xl overflow-hidden rounded-xl bg-white shadow-xl"
		in:fly={{ y: 30, duration: 800, easing: quintOut }}
	>
		<!-- Two column layout for metadata (left) and upload (right) -->
		<div class="flex flex-col md:flex-row">
			<!-- Metadata panel (left side) - shown when document is uploaded -->
			{#if documentResponse}
				<div
					class="w-full border-r border-gray-200 bg-gray-50 p-6 md:w-2/5"
					in:fly={{ x: -20, duration: 700, easing: quintOut }}
				>
					<h2
						class="mb-4 text-xl font-semibold text-gray-800"
						in:slide={{ duration: 500, delay: 100 }}
					>
						Document Metadata
					</h2>

					<!-- Success/Error message for metadata update -->
					{#if metadataUpdateSuccess !== null}
						<div
							class="mb-4 rounded-md p-3 {metadataUpdateSuccess
								? 'bg-green-50 text-green-800'
								: 'bg-red-50 text-red-800'}"
							in:fly={{ y: -10, duration: 300, easing: cubicOut }}
							out:fade
						>
							<div class="flex">
								<div class="flex-shrink-0">
									{#if metadataUpdateSuccess}
										<svg class="h-5 w-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
											<path
												fill-rule="evenodd"
												d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
												clip-rule="evenodd"
											/>
										</svg>
									{:else}
										<svg class="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
											<path
												fill-rule="evenodd"
												d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
												clip-rule="evenodd"
											/>
										</svg>
									{/if}
								</div>
								<div class="ml-3">
									<p class="text-sm font-medium">
										{metadataUpdateSuccess
											? 'Metadata updated successfully!'
											: 'Failed to update metadata. Please try again.'}
									</p>
								</div>
							</div>
						</div>
					{/if}

					<form class="space-y-4" on:submit|preventDefault={saveMetadata}>
						{#each Object.entries(currentMetadata) as [key, value], i}
							{#if key !== 'legal_tags'}
								<div
									class="form-group"
									in:fly={{ y: 10, delay: i * 50, duration: 300, easing: cubicOut }}
								>
									<label for={key} class="mb-1 block text-sm font-medium text-gray-700">
										{key.replace(/_/g, ' ').replace(/(?:^|\s)\S/g, (a) => a.toUpperCase())}
									</label>
									<input
										type="text"
										id={key}
										bind:value={currentMetadata[key]}
										class="w-full rounded-md border border-gray-300 p-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
									/>
								</div>
							{/if}
						{/each}

						<div
							class="form-group"
							in:fly={{
								y: 10,
								delay: Object.keys(currentMetadata).length * 50,
								duration: 300,
								easing: cubicOut
							}}
						>
							<label class="mb-1 block text-sm font-medium text-gray-700">Legal Tags</label>
							<div class="mb-2 flex flex-wrap gap-2">
								{#each legalTags as tag, index}
									<div
										class="inline-flex items-center rounded-full bg-indigo-100 px-2.5 py-1 text-xs font-medium text-indigo-800"
										in:scale={{ start: 0.9, duration: 300, delay: index * 30 }}
									>
										{tag}
										<button
											type="button"
											class="ml-1.5 inline-flex h-4 w-4 items-center justify-center rounded-full bg-indigo-200 text-indigo-600 hover:bg-indigo-300"
											on:click={() => removeTag(index)}
										>
											×
										</button>
									</div>
								{/each}
							</div>
							<div class="flex gap-2">
								<input
									type="text"
									placeholder="Add a tag"
									bind:value={tagInput}
									on:keydown={handleTagKeydown}
									class="flex-1 rounded-md border border-gray-300 p-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
								/>
								<button
									type="button"
									class="inline-flex justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none"
									on:click={addTag}
								>
									Add
								</button>
							</div>
						</div>

						<button
							type="submit"
							class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none disabled:cursor-not-allowed disabled:opacity-60"
							in:scale={{
								start: 0.95,
								duration: 400,
								delay: (Object.keys(currentMetadata).length + 1) * 50,
								easing: backOut
							}}
							disabled={isUpdatingMetadata}
						>
							{#if isUpdatingMetadata}
								<div class="flex items-center">
									<div
										class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white"
									></div>
									<span>Updating...</span>
								</div>
							{:else}
								Update Metadata
							{/if}
						</button>
					</form>
				</div>
			{/if}

			<!-- Upload area (right side or full width if no document) -->
			<div
				class="w-full {documentResponse ? 'md:w-3/5' : 'md:w-full'} p-6"
				in:fly={{
					x: documentResponse ? 20 : 0,
					duration: 700,
					delay: documentResponse ? 100 : 0,
					easing: quintOut
				}}
			>
				<h1
					class="mb-6 text-center text-2xl font-bold text-indigo-700"
					in:slide={{ duration: 600, delay: 200 }}
				>
					Upload Document
				</h1>

				<div
					class="dropzone-container {isDragging
						? 'border-indigo-400 bg-indigo-50'
						: 'border-gray-300'} rounded-lg border-2 border-dashed p-6 text-center transition duration-300"
					on:dragover={handleDragOver}
					on:dragleave={handleDragLeave}
					on:drop={handleDrop}
					in:fly={{ y: 15, duration: 700, delay: 300, easing: cubicOut }}
				>
					<div class="space-y-4">
						<!-- Custom styled file input -->
						<label
							for="file-upload"
							class="flex cursor-pointer flex-col items-center"
							in:scale={{ start: 0.95, duration: 600, delay: 400, easing: cubicOut }}
						>
							{#if selectedFile}
								<div
									class="mb-3 rounded-lg bg-gray-50 p-4 shadow-sm"
									in:fly={{ y: 10, duration: 600, easing: cubicOut }}
								>
									{#if getFileIconByName(selectedFile.name) === 'pdf'}
										<div class="flex justify-center">
											<div
												class="flex h-16 w-16 items-center justify-center rounded-full bg-red-100 text-red-500"
												in:scale={{ start: 0.8, duration: 700, delay: 150, easing: elasticOut }}
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="h-10 w-10"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
													></path>
													<path d="M14 2v6h6"></path>
													<path d="M10 9H8v6h2"></path>
													<path d="M12 9h2a2 2 0 0 1 0 4h-2v2"></path>
													<path d="M16 15h2"></path>
												</svg>
											</div>
										</div>
									{:else if getFileIconByName(selectedFile.name) === 'word'}
										<div class="flex justify-center">
											<div
												class="flex h-16 w-16 items-center justify-center rounded-full bg-blue-100 text-blue-500"
												in:scale={{ start: 0.8, duration: 700, delay: 150, easing: elasticOut }}
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="h-10 w-10"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
													></path>
													<path d="M14 2v6h6"></path>
													<path d="M8 12h8"></path>
													<path d="M8 16h8"></path>
													<path d="M8 8h2"></path>
												</svg>
											</div>
										</div>
									{:else if getFileIconByName(selectedFile.name) === 'text'}
										<div class="flex justify-center">
											<div
												class="flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 text-gray-500"
												in:scale={{ start: 0.8, duration: 700, delay: 150, easing: elasticOut }}
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="h-10 w-10"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
													></path>
													<path d="M14 2v6h6"></path>
													<path d="M16 13H8"></path>
													<path d="M16 17H8"></path>
													<path d="M10 9H8"></path>
												</svg>
											</div>
										</div>
									{:else}
										<div class="flex justify-center">
											<div
												class="flex h-16 w-16 items-center justify-center rounded-full bg-purple-100 text-purple-500"
												in:scale={{ start: 0.8, duration: 700, delay: 150, easing: elasticOut }}
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="h-10 w-10"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
													></path>
													<path d="M14 2v6h6"></path>
												</svg>
											</div>
										</div>
									{/if}
									<div
										class="mt-2 text-sm font-medium text-gray-700"
										in:slide={{ duration: 500, delay: 200 }}
									>
										{selectedFile.name}
									</div>
									<div class="text-xs text-gray-500" in:slide={{ duration: 500, delay: 250 }}>
										{(selectedFile.size / 1024).toFixed(1)} KB
									</div>
								</div>
							{:else}
								<div class="flex justify-center">
									<div
										class="flex h-16 w-16 items-center justify-center text-indigo-500"
										in:scale={{ start: 0.9, duration: 700, delay: 450, easing: elasticOut }}
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											class="h-12 w-12"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											stroke-width="2"
											stroke-linecap="round"
											stroke-linejoin="round"
										>
											<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
											<polyline points="17 8 12 3 7 8"></polyline>
											<line x1="12" y1="3" x2="12" y2="15"></line>
										</svg>
									</div>
								</div>
								<div class="mt-3 text-gray-600" in:slide={{ duration: 500, delay: 500 }}>
									{fileInputLabel}
								</div>
							{/if}
						</label>
						<input
							id="file-upload"
							type="file"
							accept=".pdf,.docx,.txt"
							on:change={handleFileChange}
							class="sr-only"
						/>

						<div class="text-xs text-gray-500" in:fade={{ duration: 600, delay: 550 }}>
							Supported files: PDF, DOCX, TXT
						</div>
					</div>
				</div>

				<div class="mt-6" in:fly={{ y: 10, duration: 600, delay: 600, easing: cubicOut }}>
					<button
						class="flex w-full items-center justify-center rounded-lg bg-indigo-600 px-4 py-3 font-medium text-white shadow transition duration-200 hover:bg-indigo-700 disabled:cursor-not-allowed disabled:opacity-50"
						on:click={handleFileUpload}
						disabled={!selectedFile || isUploading}
						in:scale={{ start: 0.95, duration: 600, delay: 650, easing: backOut }}
					>
						{#if isUploading}
							<div class="flex items-center">
								<!-- Custom spinner -->
								<div
									class="mr-2 h-5 w-5 animate-spin rounded-full border-2 border-white/30 border-t-white"
								></div>
								<span>Processing...</span>
							</div>
						{:else}
							Upload and Categorise
						{/if}
					</button>
				</div>

				{#if uploadStatus && !isUploading}
					<div
						class="mt-4 rounded-lg p-3 text-center {uploadStatus.includes('Failed')
							? 'bg-red-50 text-red-700'
							: 'bg-green-50 text-green-700'}"
						in:fly={{ y: 5, duration: 600, easing: cubicOut }}
					>
						{uploadStatus}
					</div>
				{/if}

				{#if uploadedDocuments.length > 0}
					<div
						class="mt-8 border-t border-gray-200 pt-6"
						in:fly={{ y: 20, duration: 700, easing: cubicOut }}
					>
						<h2 class="mb-4 text-lg font-semibold" in:slide={{ duration: 600, delay: 100 }}>
							Uploaded Documents
						</h2>
						<div class="grid grid-cols-3 gap-4 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6">
							{#each uploadedDocuments as doc, i}
								<div
									class="flex cursor-pointer flex-col items-center rounded-lg bg-white p-3 shadow-sm transition-shadow hover:shadow-md"
									on:click={() => showDocumentDetails(doc.response)}
									in:fly={{ y: 20, x: 5, duration: 600, delay: 200 + i * 100, easing: cubicOut }}
								>
									{#if getFileIcon(doc.type) === 'pdf' || getFileIconByName(doc.name) === 'pdf'}
										<div class="flex justify-center">
											<div
												class="flex h-10 w-10 items-center justify-center rounded-full bg-red-100 text-red-500"
												in:scale={{
													start: 0.85,
													duration: 600,
													delay: 250 + i * 100,
													easing: elasticOut
												}}
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="h-6 w-6"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
													></path>
													<path d="M14 2v6h6"></path>
													<path d="M10 9H8v6h2"></path>
													<path d="M12 9h2a2 2 0 0 1 0 4h-2v2"></path>
													<path d="M16 15h2"></path>
												</svg>
											</div>
										</div>
									{:else if getFileIcon(doc.type) === 'word' || getFileIconByName(doc.name) === 'word'}
										<div class="flex justify-center">
											<div
												class="flex h-10 w-10 items-center justify-center rounded-full bg-blue-100 text-blue-500"
												in:scale={{
													start: 0.85,
													duration: 600,
													delay: 250 + i * 100,
													easing: elasticOut
												}}
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="h-6 w-6"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
													></path>
													<path d="M14 2v6h6"></path>
													<path d="M8 12h8"></path>
													<path d="M8 16h8"></path>
													<path d="M8 8h2"></path>
												</svg>
											</div>
										</div>
									{:else if getFileIcon(doc.type) === 'text' || getFileIconByName(doc.name) === 'text'}
										<div class="flex justify-center">
											<div
												class="flex h-10 w-10 items-center justify-center rounded-full bg-gray-100 text-gray-500"
												in:scale={{
													start: 0.85,
													duration: 600,
													delay: 250 + i * 100,
													easing: elasticOut
												}}
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="h-6 w-6"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
													></path>
													<path d="M14 2v6h6"></path>
													<path d="M16 13H8"></path>
													<path d="M16 17H8"></path>
													<path d="M10 9H8"></path>
												</svg>
											</div>
										</div>
									{:else}
										<div class="flex justify-center">
											<div
												class="flex h-10 w-10 items-center justify-center rounded-full bg-purple-100 text-purple-500"
												in:scale={{
													start: 0.85,
													duration: 600,
													delay: 250 + i * 100,
													easing: elasticOut
												}}
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="h-6 w-6"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													stroke-linecap="round"
													stroke-linejoin="round"
												>
													<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
													></path>
													<path d="M14 2v6h6"></path>
												</svg>
											</div>
										</div>
									{/if}
									<div
										class="mt-2 w-full truncate text-center text-xs text-gray-700"
										title={doc.name}
									>
										{doc.name.length > 15 ? doc.name.substring(0, 12) + '...' : doc.name}
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>
