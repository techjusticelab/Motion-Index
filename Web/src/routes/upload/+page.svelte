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

	// Metadata panel state
	let isMetadataPanelOpen = false;
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

	// Add these to your script section
	let isEditingMetadata = false;
	let editableMetadata = null;

	function startEditingMetadata(document: Document) {
		// Create a deep copy of the document for editing
		editableMetadata = JSON.parse(JSON.stringify(document));
		isEditingMetadata = true;
	}

	function openMetadataPanel(document: Document) {
		if (!document || !document.metadata) return;

		// Set the documentId for the API call
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

		// Open the panel
		isMetadataPanelOpen = true;
	}

	function closeMetadataPanel() {
		isMetadataPanelOpen = false;
	}

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

	async function saveMetadataFromPanel() {
		if (!currentDocumentId) {
			uploadStatus = 'Error: No document selected for metadata update';
			return;
		}

		try {
			uploadStatus = 'Updating document metadata...';

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

			// Close the panel and show success message
			closeMetadataPanel();
			uploadStatus = 'Document metadata updated successfully!';
		} catch (error) {
			uploadStatus = 'Failed to update document metadata. Please try again.';
			console.error(error);
		}
	}

	async function saveMetadata() {
		if (!editableMetadata || !editableMetadata.id) {
			uploadStatus = 'Error: No document selected for metadata update';
			return;
		}

		try {
			uploadStatus = 'Updating document metadata...';

			// Extract the metadata fields we want to update
			const metadataToUpdate = editableMetadata.metadata;

			// Update document metadata through API
			const response = await updateDocumentMetadata(editableMetadata.id, metadataToUpdate);

			// Update the document in our list
			uploadedDocuments = uploadedDocuments.map((doc) => {
				if (doc.response.id === editableMetadata.id) {
					return {
						...doc,
						response: response.document
					};
				}
				return doc;
			});

			// Update the current document response
			documentResponse = response.document;

			// Exit edit mode
			isEditingMetadata = false;
			uploadStatus = 'Document metadata updated successfully!';
		} catch (error) {
			uploadStatus = 'Failed to update document metadata. Please try again.';
			console.error(error);
		}
	}

	// Modify showDocumentDetails to initialize editable metadata
	function showDocumentDetails(document: Document) {
		documentResponse = document;

		// Open the metadata panel when a document is selected
		openMetadataPanel(document);

		// Reset the editing state
		isEditingMetadata = false;
		editableMetadata = null;
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

			// Add the document to our list of uploaded documents
			uploadedDocuments = [
				...uploadedDocuments,
				{
					name: selectedFile.name,
					type: selectedFile.type,
					response: response.document
				}
			];

			// Automatically open the metadata panel when a document is uploaded
			openMetadataPanel(response.document);

			// Reset the file input
			selectedFile = null;
			fileInputLabel = 'Drag and drop your file here or click to browse';

			uploadStatus = response.data;
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

<div
	class="flex min-h-[80vh] items-center justify-center"
	in:fly={{ y: 30, duration: 800, easing: cubicOut }}
>
	<div
		class="upload-container rounded-lg bg-white shadow-lg"
		in:fly={{ y: 20, duration: 700, delay: 200, easing: cubicOut }}
	>
		<h1
			class="mb-6 text-center text-2xl font-bold text-indigo-700"
			in:slide={{ duration: 600, delay: 300 }}
		>
			Upload Document
		</h1>

		<div
			class="dropzone-container {isDragging ? 'dragging' : ''}"
			on:dragover={handleDragOver}
			on:dragleave={handleDragLeave}
			on:drop={handleDrop}
			in:fly={{ y: 15, duration: 700, delay: 400, easing: cubicOut }}
		>
			<div class="file-input-container">
				<!-- Custom styled file input -->
				<label
					for="file-upload"
					class="file-input-label"
					in:scale={{ start: 0.95, duration: 600, delay: 500, easing: cubicOut }}
				>
					{#if selectedFile}
						<div class="document-preview" in:fly={{ y: 10, duration: 600, easing: cubicOut }}>
							{#if getFileIconByName(selectedFile.name) === 'pdf'}
								<div
									class="document-icon pdf"
									in:scale={{ start: 0.8, duration: 700, delay: 150, easing: elasticOut }}
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
										<path d="M14 2v6h6"></path>
										<path d="M10 9H8v6h2"></path>
										<path d="M12 9h2a2 2 0 0 1 0 4h-2v2"></path>
										<path d="M16 15h2"></path>
									</svg>
								</div>
							{:else if getFileIconByName(selectedFile.name) === 'word'}
								<div
									class="document-icon word"
									in:scale={{ start: 0.8, duration: 700, delay: 150, easing: elasticOut }}
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
										<path d="M14 2v6h6"></path>
										<path d="M8 12h8"></path>
										<path d="M8 16h8"></path>
										<path d="M8 8h2"></path>
									</svg>
								</div>
							{:else if getFileIconByName(selectedFile.name) === 'text'}
								<div
									class="document-icon text"
									in:scale={{ start: 0.8, duration: 700, delay: 150, easing: elasticOut }}
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
										<path d="M14 2v6h6"></path>
										<path d="M16 13H8"></path>
										<path d="M16 17H8"></path>
										<path d="M10 9H8"></path>
									</svg>
								</div>
							{:else}
								<div
									class="document-icon generic"
									in:scale={{ start: 0.8, duration: 700, delay: 150, easing: elasticOut }}
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
										<path d="M14 2v6h6"></path>
									</svg>
								</div>
							{/if}
							<div class="document-name selected" in:slide={{ duration: 500, delay: 200 }}>
								{selectedFile.name}
							</div>
							<div class="document-size" in:slide={{ duration: 500, delay: 250 }}>
								{(selectedFile.size / 1024).toFixed(1)} KB
							</div>
						</div>
					{:else}
						<div
							class="upload-icon"
							in:scale={{ start: 0.9, duration: 700, delay: 550, easing: elasticOut }}
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="48"
								height="48"
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
						<div class="file-input-text" in:slide={{ duration: 500, delay: 600 }}>
							{fileInputLabel}
						</div>
					{/if}
				</label>
				<input
					id="file-upload"
					type="file"
					accept=".pdf,.docx,.txt"
					on:change={handleFileChange}
					class="hidden-file-input"
				/>
			</div>

			<div class="file-types" in:fade={{ duration: 600, delay: 650 }}>
				Supported files: PDF, DOCX, TXT
			</div>
		</div>

		<div
			class="upload-button-container"
			in:fly={{ y: 10, duration: 600, delay: 700, easing: cubicOut }}
		>
			<button
				class="upload-button {!selectedFile ? 'disabled' : ''}"
				on:click={handleFileUpload}
				disabled={!selectedFile || isUploading}
				in:scale={{ start: 0.95, duration: 600, delay: 750, easing: backOut }}
			>
				{#if isUploading}
					<div class="spinner-container">
						<!-- Simple custom spinner -->
						<div class="spinner"></div>
						<span class="ml-2">Processing...</span>
					</div>
				{:else}
					Upload and Categorise
				{/if}
			</button>
		</div>

		{#if uploadStatus && !isUploading}
			<div
				class="status-message {uploadStatus.includes('Failed') ? 'error' : 'success'}"
				in:fly={{ y: 5, duration: 600, easing: cubicOut }}
			>
				{uploadStatus}
			</div>
		{/if}

		{#if uploadedDocuments.length > 0}
			<div class="uploaded-documents" in:fly={{ y: 20, duration: 700, easing: cubicOut }}>
				<h2 class="mb-3 text-lg font-semibold" in:slide={{ duration: 600, delay: 100 }}>
					Uploaded Documents
				</h2>
				<div class="document-grid">
					{#each uploadedDocuments as doc, i}
						<div
							class="document-item"
							on:click={() => showDocumentDetails(doc.response)}
							in:fly={{ y: 20, x: 5, duration: 600, delay: 200 + i * 100, easing: cubicOut }}
						>
							{#if getFileIcon(doc.type) === 'pdf' || getFileIconByName(doc.name) === 'pdf'}
								<div
									class="document-icon pdf"
									in:scale={{
										start: 0.85,
										duration: 600,
										delay: 250 + i * 100,
										easing: elasticOut
									}}
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
										<path d="M14 2v6h6"></path>
										<path d="M10 9H8v6h2"></path>
										<path d="M12 9h2a2 2 0 0 1 0 4h-2v2"></path>
										<path d="M16 15h2"></path>
									</svg>
								</div>
							{:else if getFileIcon(doc.type) === 'word' || getFileIconByName(doc.name) === 'word'}
								<div
									class="document-icon word"
									in:scale={{
										start: 0.85,
										duration: 600,
										delay: 250 + i * 100,
										easing: elasticOut
									}}
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
										<path d="M14 2v6h6"></path>
										<path d="M8 12h8"></path>
										<path d="M8 16h8"></path>
										<path d="M8 8h2"></path>
									</svg>
								</div>
							{:else if getFileIcon(doc.type) === 'text' || getFileIconByName(doc.name) === 'text'}
								<div class="document-icon text">
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
										<path d="M14 2v6h6"></path>
										<path d="M16 13H8"></path>
										<path d="M16 17H8"></path>
										<path d="M10 9H8"></path>
									</svg>
								</div>
							{:else}
								<div class="document-icon generic">
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
										<path d="M14 2v6h6"></path>
									</svg>
								</div>
							{/if}
							<div class="document-name" title={doc.name}>
								{doc.name.length > 20 ? doc.name.substring(0, 17) + '...' : doc.name}
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		{#if documentResponse}
			<div class="document-details">
				<h2 class="text-lg font-semibold">Document Details</h2>

				<!-- Add edit button and metadata form -->
				<div class="metadata-controls mb-4">
					<button
						class="edit-metadata-button"
						on:click={() => (isEditingMetadata = !isEditingMetadata)}
					>
						{isEditingMetadata ? 'Cancel' : 'Edit Metadata'}
					</button>
				</div>

				{#if isEditingMetadata}
					<div class="metadata-form">
						<div class="form-group">
							<label for="doc-type">Document Type</label>
							<input id="doc-type" type="text" bind:value={editableMetadata.doc_type} />
						</div>
						<div class="form-group">
							<label for="category">Category</label>
							<input id="category" type="text" bind:value={editableMetadata.category} />
						</div>
						<div class="form-group">
							<label for="case-name">Case Name</label>
							<input id="case-name" type="text" bind:value={editableMetadata.metadata.case_name} />
						</div>
						<div class="form-group">
							<label for="case-number">Case Number</label>
							<input
								id="case-number"
								type="text"
								bind:value={editableMetadata.metadata.case_number}
							/>
						</div>
						<div class="form-group">
							<label for="court">Court</label>
							<input id="court" type="text" bind:value={editableMetadata.metadata.court} />
						</div>
						<div class="form-group">
							<label for="judge">Judge</label>
							<input id="judge" type="text" bind:value={editableMetadata.metadata.judge} />
						</div>
						<div class="form-group">
							<label for="status">Status</label>
							<input id="status" type="text" bind:value={editableMetadata.metadata.status} />
						</div>

						<div class="form-actions">
							<button class="save-metadata-button" on:click={saveMetadata}> Save Changes </button>
						</div>
					</div>
				{:else}
					<pre class="details-json">{JSON.stringify(documentResponse, null, 2)}</pre>
				{/if}
			</div>
		{/if}
	</div>
</div>

<style>
    .metadata-form {
    background-color: #f9fafb;
    border-radius: 0.5rem;
    padding: 1rem;
    margin-bottom: 1rem;
    border: 1px solid #e5e7eb;
}

.form-group {
    margin-bottom: 1rem;
}

.form-group label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    color: #4b5563;
    margin-bottom: 0.25rem;
}

.form-group input {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #d1d5db;
    border-radius: 0.375rem;
    font-size: 0.875rem;
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 1.5rem;
}

.edit-metadata-button, .save-metadata-button {
    padding: 0.5rem 1rem;
    border-radius: 0.375rem;
    font-size: 0.875rem;
    font-weight: 500;
    transition: all 0.2s;
}

.edit-metadata-button {
    background-color: #f3f4f6;
    color: #4b5563;
    border: 1px solid #d1d5db;
    margin-right: 0.5rem;
}

.edit-metadata-button:hover {
    background-color: #e5e7eb;
}

.save-metadata-button {
    background-color: #6366f1;
    color: white;
    border: none;
}

.save-metadata-button:hover {
    background-color: #4f46e5;
}
	.upload-container {
		width: 90%;
		max-width: 700px;
		padding: 2.5rem;
	}

	.dropzone-container {
		border: 2px dashed #cbd5e0;
		border-radius: 8px;
		padding: 2rem;
		text-align: center;
		transition: all 0.3s ease;
		margin-bottom: 1.5rem;
	}

	.dropzone-container.dragging {
		border-color: #6366f1;
		background-color: #f5f5ff;
	}

	.file-input-container {
		margin-bottom: 1rem;
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
		flex-direction: column;
		align-items: center;
		cursor: pointer;
		width: 100%;
	}

	.upload-icon {
		color: #6366f1;
		margin-bottom: 1rem;
	}

	.file-input-text {
		color: #4b5563;
		margin-bottom: 0.5rem;
		font-size: 1rem;
	}

	.document-preview {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 1.5rem;
		border-radius: 8px;
		background-color: #f9fafb;
		margin-bottom: 1rem;
		width: 80%;
		max-width: 300px;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
	}

	.file-name {
		font-size: 0.875rem;
		color: #6366f1;
		font-weight: 500;
		margin-top: 0.5rem;
	}

	.file-types {
		font-size: 0.75rem;
		color: #6b7280;
		margin-top: 0.5rem;
	}

	.upload-button-container {
		display: flex;
		justify-content: center;
		margin-top: 1.5rem;
	}

	.upload-button {
		background-color: #6366f1;
		color: white;
		font-weight: 500;
		padding: 0.75rem 2rem;
		border-radius: 0.375rem;
		transition: all 0.2s ease;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
		min-width: 200px;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.upload-button:hover:not(.disabled) {
		background-color: #4f46e5;
		transform: translateY(-1px);
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	}

	.upload-button.disabled {
		background-color: #9ca3af
		cursor: not-allowed;
	}

	.spinner-container {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	/* Custom spinner */
	.spinner {
		width: 20px;
		height: 20px;
		border: 3px solid rgba(255, 255, 255, 0.3);
		border-radius: 50%;
		border-top-color: white;
		animation: spin 1s ease-in-out infinite;
		margin-right: 8px;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.ml-2 {
		margin-left: 0.5rem;
	}

	.status-message {
		margin-top: 1.5rem;
		padding: 0.75rem;
		border-radius: 0.375rem;
		text-align: center;
	}

	.status-message.success {
		background-color: #ecfdf5;
		color: #047857;
	}

	.status-message.error {
		background-color: #fef2f2;
		color: #b91c1c;
	}

	.uploaded-documents {
		margin-top: 2rem;
		padding-top: 1.5rem;
		border-top: 1px solid #e5e7eb;
	}

	.document-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
		gap: 1rem;
		margin-top: 1rem;
	}

	.document-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		cursor: pointer;
		transition: transform 0.2s ease;
	}

	.document-item:hover {
		transform: scale(1.05);
	}

	.document-icon {
		width: 50px;
		height: 50px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 8px;
		margin-bottom: 0.5rem;
		padding: 0.5rem;
	}

	.document-icon.pdf {
		background-color: #fee2e2;
		color: #dc2626;
	}

	.document-icon.word {
		background-color: #e0f2fe;
		color: #0284c7;
	}

	.document-icon.text {
		background-color: #f3f4f6;
		color: #4b5563;
	}

	.document-icon.generic {
		background-color: #f3e8ff;
		color: #7c3aed;
	}

	.document-name {
		font-size: 0.75rem;
		text-align: center;
		color: #4b5563;
		max-width: 100%;
		word-break: break-word;
	}

	.document-name.selected {
		font-size: 0.875rem;
		font-weight: 500;
		color: #1f2937;
		margin-top: 0.75rem;
		margin-bottom: 0.25rem;
	}

	.document-size {
		font-size: 0.75rem;
		color: #6b7280;
	}

	.document-details {
		margin-top: 2rem;
	}

	.details-json {
		background-color: #f9fafb;
		border-radius: 0.5rem;
		padding: 1rem;
		overflow-x: auto;
		margin-top: 0.5rem;
		font-size: 0.875rem;
		border: 1px solid #e5e7eb;
	}
</style>
