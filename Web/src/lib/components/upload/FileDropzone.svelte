<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { useFileUpload } from '$lib/hooks/useFileUpload';
	import { LoadingSpinner } from '$lib/components/ui';
	import { fade, scale, slide } from 'svelte/transition';
	import { elasticOut, cubicOut } from 'svelte/easing';

	const dispatch = createEventDispatcher();

	interface Props {
		accept?: string[];
		maxSizeMB?: number;
		disabled?: boolean;
		class?: string;
		multiple?: boolean;
		batchMode?: boolean;
	}

	let {
		accept = ['.pdf', '.docx', '.txt'],
		maxSizeMB = 50,
		disabled = false,
		class: className = '',
		multiple = false,
		batchMode = false
	}: Props = $props();

	// Initialize file upload hook
	const fileUpload = useFileUpload({
		allowedTypes: accept,
		maxSizeBytes: maxSizeMB * 1024 * 1024
	});

	// Subscribe to file upload state
	let { selectedFile, isDragging, isUploading, uploadStatus } = $derived($fileUpload);

	// Batch mode state
	let selectedFiles = $state<File[]>([]);

	// Generate unique ID for accessibility
	const dropzoneId = `dropzone-${Math.random().toString(36).substr(2, 9)}`;

	// Handle file selection and emit event
	function onFileSelected(file: File) {
		dispatch('fileSelected', { file });
	}

	// Handle multiple file selection for batch mode
	function onFilesSelected(files: File[]) {
		selectedFiles = files;
		dispatch('filesSelected', { files });
	}

	// Handle file input change for batch mode
	function handleBatchFileChange(event: Event) {
		const target = event.target as HTMLInputElement;
		if (target.files && target.files.length > 0) {
			const files = Array.from(target.files);
			onFilesSelected(files);
		}
	}

	// Batch mode drag state
	let batchIsDragging = $state(false);

	// Handle drag events for batch mode
	function handleBatchDragOver(event: DragEvent) {
		event.preventDefault();
		batchIsDragging = true;
	}

	function handleBatchDragLeave(event: DragEvent) {
		event.preventDefault();
		batchIsDragging = false;
	}

	function handleBatchDrop(event: DragEvent) {
		event.preventDefault();
		batchIsDragging = false;
		
		if (disabled) return;
		
		const files = event.dataTransfer?.files;
		if (files && files.length > 0) {
			const fileArray = Array.from(files).filter(file => {
				const extension = '.' + file.name.split('.').pop()?.toLowerCase();
				return accept.includes(extension) && file.size <= maxSizeMB * 1024 * 1024;
			});
			onFilesSelected(fileArray);
		}
	}

	// Remove file from batch selection
	function removeFile(index: number) {
		selectedFiles = selectedFiles.filter((_, i) => i !== index);
		dispatch('filesSelected', { files: selectedFiles });
	}

	// Watch for file changes and emit events (single file mode)
	$effect(() => {
		if (selectedFile && !batchMode) {
			onFileSelected(selectedFile);
		}
	});

	// Get file size display
	function getFileSizeDisplay(file: File): string {
		const size = file.size;
		if (size < 1024) return `${size} B`;
		if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
		return `${(size / (1024 * 1024)).toFixed(1)} MB`;
	}

	// Get status styling
	function getStatusStyling(status: string): string {
		if (status.toLowerCase().includes('error') || status.toLowerCase().includes('failed')) {
			return 'bg-red-50 text-red-700 border-red-200';
		}
		if (status.toLowerCase().includes('success') || status.toLowerCase().includes('completed')) {
			return 'bg-green-50 text-green-700 border-green-200';
		}
		return 'bg-blue-50 text-blue-700 border-blue-200';
	}
</script>

<div class="space-y-4 {className}">
	<!-- Main dropzone -->
	<div
		class="relative overflow-hidden rounded-xl border-2 border-dashed transition-all duration-300 {(batchMode ? batchIsDragging : isDragging)
			? 'border-primary-400 bg-primary-50'
			: disabled
				? 'border-neutral-200 bg-neutral-50'
				: 'border-neutral-300 bg-white hover:border-primary-300 hover:bg-primary-50/30'}"
		class:cursor-not-allowed={disabled}
		class:cursor-pointer={!disabled}
	>
		<label
			for={dropzoneId}
			class="block p-8 text-center"
			ondragover={disabled ? undefined : (batchMode ? handleBatchDragOver : fileUpload.handleDragOver)}
			ondragleave={disabled ? undefined : (batchMode ? handleBatchDragLeave : fileUpload.handleDragLeave)}
			ondrop={disabled ? undefined : (batchMode ? handleBatchDrop : fileUpload.handleDrop)}
		>
			{#if isUploading}
				<!-- Upload in progress -->
				<div class="flex flex-col items-center space-y-3" transition:fade={{ duration: 300 }}>
					<LoadingSpinner size="lg" color="primary" />
					<div class="text-sm font-medium text-neutral-700">Processing file{batchMode && selectedFiles.length > 1 ? 's' : ''}...</div>
				</div>
			{:else if batchMode && selectedFiles.length > 0}
				<!-- Batch files selected -->
				<div class="space-y-4" transition:scale={{ start: 0.95, duration: 300 }}>
					<div class="flex items-center justify-between">
						<div class="font-medium text-neutral-900">
							{selectedFiles.length} file{selectedFiles.length !== 1 ? 's' : ''} selected
						</div>
						<button
							type="button"
							class="text-sm text-red-600 hover:text-red-700 font-medium"
							onclick={() => onFilesSelected([])}
						>
							Clear all
						</button>
					</div>
					
					<!-- File list -->
					<div class="max-h-48 overflow-y-auto space-y-2">
						{#each selectedFiles as file, index}
							<div class="flex items-center justify-between bg-neutral-50 rounded-lg p-3" transition:slide={{ duration: 300 }}>
								<div class="flex items-center space-x-3">
									<!-- File icon -->
									<div class="flex h-8 w-8 items-center justify-center rounded-full bg-primary-100 text-primary-600 flex-shrink-0">
										{#if file.type.includes('pdf')}
											<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
											</svg>
										{:else if file.type.includes('word')}
											<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
											</svg>
										{:else}
											<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
											</svg>
										{/if}
									</div>
									
									<!-- File info -->
									<div class="min-w-0">
										<div class="text-sm font-medium text-neutral-900 truncate">{file.name}</div>
										<div class="text-xs text-neutral-500">{getFileSizeDisplay(file)}</div>
									</div>
								</div>
								
								<!-- Remove button -->
								<button
									type="button"
									class="text-neutral-400 hover:text-red-500 flex-shrink-0"
									onclick={() => removeFile(index)}
								>
									<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
									</svg>
								</button>
							</div>
						{/each}
					</div>
					
					<!-- Add more files button -->
					<button
						type="button"
						class="text-sm text-primary-600 hover:text-primary-700 font-medium underline"
						onclick={() => document.getElementById(dropzoneId)?.click()}
					>
						Add more files
					</button>
				</div>
			{:else if selectedFile && !batchMode}
				<!-- Single file selected -->
				<div class="flex flex-col items-center space-y-3" transition:scale={{ start: 0.95, duration: 300 }}>
					<!-- File icon -->
					<div class="flex h-16 w-16 items-center justify-center rounded-full bg-primary-100 text-primary-600">
						{#if selectedFile.type.includes('pdf')}
							<svg class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
							</svg>
						{:else if selectedFile.type.includes('word')}
							<svg class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
							</svg>
						{:else}
							<svg class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
							</svg>
						{/if}
					</div>

					<!-- File info -->
					<div class="space-y-1">
						<div class="font-medium text-neutral-900" transition:slide={{ duration: 300 }}>
							{selectedFile.name}
						</div>
						<div class="text-sm text-neutral-500" transition:slide={{ duration: 300, delay: 100 }}>
							{getFileSizeDisplay(selectedFile)}
						</div>
					</div>

					<!-- Change file button -->
					<button
						type="button"
						class="text-sm text-primary-600 hover:text-primary-700 font-medium underline"
						onclick={fileUpload.clearFile}
					>
						Choose different file
					</button>
				</div>
			{:else}
				<!-- Default state -->
				<div class="flex flex-col items-center space-y-4" transition:fade={{ duration: 300 }}>
					<!-- Upload icon -->
					<div
						class="flex h-16 w-16 items-center justify-center text-neutral-400"
						transition:scale={{ start: 0.9, duration: 400, delay: 200, easing: elasticOut }}
					>
						<svg class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
						</svg>
					</div>

					<!-- Instructions -->
					<div class="space-y-2">
						<div class="text-lg font-medium text-neutral-700" transition:slide={{ duration: 300, delay: 300 }}>
							{(batchMode ? batchIsDragging : isDragging)
								? `Drop file${batchMode ? 's' : ''} here` 
								: `Drag and drop your file${batchMode ? 's' : ''} here`}
						</div>
						<div class="text-sm text-neutral-500" transition:slide={{ duration: 300, delay: 400 }}>
							or click to browse files
						</div>
					</div>

					<!-- Supported formats -->
					<div class="text-xs text-neutral-400" transition:fade={{ duration: 300, delay: 500 }}>
						Supported formats: {accept.join(', ').toUpperCase()}
						{#if maxSizeMB}
							• Max size: {maxSizeMB}MB{batchMode ? ' per file' : ''}
						{/if}
					</div>
				</div>
			{/if}
		</label>

		<!-- Hidden file input -->
		<input
			id={dropzoneId}
			type="file"
			accept={accept.join(',')}
			{disabled}
			multiple={batchMode}
			onchange={batchMode ? handleBatchFileChange : fileUpload.handleFileChange}
			class="sr-only"
		/>
	</div>

	<!-- Upload status -->
	{#if uploadStatus}
		<div
			class="rounded-lg border p-3 text-sm {getStatusStyling(uploadStatus)}"
			transition:slide={{ duration: 300, easing: cubicOut }}
		>
			{uploadStatus}
		</div>
	{/if}
</div>