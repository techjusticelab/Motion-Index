<!-- DocumentViewerModal.svelte -->
<script lang="ts">
	import type { Document } from '$lib/types';
	import { formatDocumentType } from '$lib/utils/utils';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { quintOut, cubicOut } from 'svelte/easing';

	const {
		docData,
		isOpen,
		onclose,
		metadata,
		caseActions,
		pdfViewer
	}: {
		docData?: Document | null;
		isOpen?: boolean;
		onclose?: () => void;
		metadata?: any;
		caseActions?: any;
		pdfViewer?: any;
	} = $props();

	// Function to close the viewer
	function closeViewer() {
		onclose?.();
	}

	// Handle keydown for escape
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			closeViewer();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if isOpen && docData}
	<div
		class="opac fixed inset-0 z-50 flex justify-center p-4 shadow"
		in:fade={{ duration: 300, easing: cubicOut }}
		out:fade={{ duration: 500 }}
	>
		<div
			class="relative flex h-[95vh] w-[95vw] overflow-hidden text-wrap rounded-xl bg-white shadow-2xl flex-col md:flex-row"
			transition:fly={{ y: 20, duration: 800, easing: quintOut }}
		>
			<!-- Sidebar with metadata -->
			<div
				class="w-full md:w-1/4 overflow-auto border-r border-neutral-200 bg-neutral-50 p-4 md:p-6 max-h-[40vh] md:max-h-none"
				in:fly={{ x: -20, duration: 800, delay: 300, easing: cubicOut }}
			>
				<div class="mb-4">
					<h2
						class="truncate text-lg md:text-xl font-semibold text-neutral-800"
						in:slide={{ duration: 700, delay: 400 }}
					>
						{docData.metadata?.document_name || docData.file_name || 'Untitled Document'}
					</h2>
					<p class="mt-1 text-xs md:text-sm text-neutral-500" in:slide={{ duration: 700, delay: 500 }}>
						{docData.file_name}
					</p>
				</div>

				<!-- Document type tag -->
				<div class="mb-4">
					<span
						class="inline-flex rounded-full bg-primary-50 px-2 md:px-3 py-1 text-xs md:text-sm font-medium text-primary-700"
						in:scale={{ start: 0.9, duration: 600, delay: 600, easing: cubicOut }}
					>
						{formatDocumentType(docData.doc_type)}
					</span>
				</div>

				<!-- Slot for metadata component -->
				<div in:slide={{ duration: 600, delay: 650 }}>
					{@render metadata?.()}
				</div>

				<!-- Slot for case actions component -->
				<div>
					{@render caseActions?.()}
				</div>
			</div>

			<!-- Document Viewer area -->
			<div
				class="relative w-full flex-1 overflow-hidden min-h-[60vh] md:min-h-auto"
				in:fly={{ x: 20, duration: 800, delay: 350, easing: cubicOut }}
			>
				<!-- Header with close button -->
				<div
					class="absolute right-0 top-0 z-10 flex items-center justify-end p-4"
					in:fly={{ y: -10, duration: 700, delay: 800, easing: cubicOut }}
				>
					<button
						type="button"
						class="rounded-full bg-white/90 p-2 shadow-md hover:bg-neutral-100"
						onclick={closeViewer}
						in:scale={{ start: 0.9, duration: 600, delay: 900 }}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6 text-neutral-500"
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

				<!-- Slot for PDF viewer component -->
				{@render pdfViewer?.()}
			</div>
		</div>
	</div>
{/if}

<style>
	.opac {
		background-color: rgba(0, 0, 0, 0.5);
	}
</style>