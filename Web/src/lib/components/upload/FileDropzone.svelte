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
	}

	let {
		accept = ['.pdf', '.docx', '.txt'],
		maxSizeMB = 50,
		disabled = false,
		class: className = ''
	}: Props = $props();

	// Initialize file upload hook
	const fileUpload = useFileUpload({
		allowedTypes: accept,
		maxSizeBytes: maxSizeMB * 1024 * 1024
	});

	// Subscribe to file upload state
	let { selectedFile, isDragging, isUploading, uploadStatus } = $derived($fileUpload);

	// Generate unique ID for accessibility
	const dropzoneId = `dropzone-${Math.random().toString(36).substr(2, 9)}`;

	// Handle file selection and emit event
	function onFileSelected(file: File) {
		dispatch('fileSelected', { file });
	}

	// Watch for file changes and emit events
	$effect(() => {
		if (selectedFile) {
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
		class="relative overflow-hidden rounded-xl border-2 border-dashed transition-all duration-300 {isDragging
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
			ondragover={disabled ? undefined : fileUpload.handleDragOver}
			ondragleave={disabled ? undefined : fileUpload.handleDragLeave}
			ondrop={disabled ? undefined : fileUpload.handleDrop}
		>
			{#if isUploading}
				<!-- Upload in progress -->
				<div class="flex flex-col items-center space-y-3" transition:fade={{ duration: 300 }}>
					<LoadingSpinner size="lg" color="primary" />
					<div class="text-sm font-medium text-neutral-700">Processing file...</div>
				</div>
			{:else if selectedFile}
				<!-- File selected -->
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
							{isDragging ? 'Drop file here' : 'Drag and drop your file here'}
						</div>
						<div class="text-sm text-neutral-500" transition:slide={{ duration: 300, delay: 400 }}>
							or click to browse files
						</div>
					</div>

					<!-- Supported formats -->
					<div class="text-xs text-neutral-400" transition:fade={{ duration: 300, delay: 500 }}>
						Supported formats: {accept.join(', ').toUpperCase()}
						{#if maxSizeMB}
							• Max size: {maxSizeMB}MB
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
			onchange={fileUpload.handleFileChange}
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