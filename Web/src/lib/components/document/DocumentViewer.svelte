<!-- DocumentViewer.svelte - Refactored to orchestrate smaller components -->
<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Document } from '$lib/types';
	import DocumentViewerModal from './DocumentViewerModal.svelte';
	import DocumentPDFViewer from './DocumentPDFViewer.svelte';
	import DocumentMetadata from './DocumentMetadata.svelte';
	import CaseActions from './CaseActions.svelte';

	let {
		docData = $bindable(),
		isOpen = $bindable(false),
		supabase = null,
		session = null,
		user = null
	}: {
		docData?: Document | null;
		isOpen?: boolean;
		supabase?: any;
		session?: any;
		user?: any;
	} = $props();

	// State for PDF viewer URL (shared between components)
	let pdfViewerUrl = $state('');

	const dispatch = createEventDispatcher<{
		close: void;
		search: {
			field: string;
			value: string;
		};
	}>();

	// Handle events from child components
	function handleClose() {
		isOpen = false;
		dispatch('close');
	}

	function handleSearch(event: CustomEvent<{ field: string; value: string }>) {
		// Dispatch search event that can be caught by parent component
		dispatch('search', event.detail);
		
		// Close the viewer since we're navigating to search
		handleClose();
	}

	function handlePDFLoaded() {
		// PDF viewer has loaded successfully
		console.log('PDF viewer loaded');
	}

	function handlePDFError(event: CustomEvent<{ message: string }>) {
		// PDF viewer encountered an error
		console.error('PDF viewer error:', event.detail.message);
	}
</script>

{#if isOpen && docData}
	<DocumentViewerModal
		{docData}
		{isOpen}
		onclose={handleClose}
		{metadata}
		{caseActions}
		{pdfViewer}
	/>
{/if}

{#snippet metadata()}
	<DocumentMetadata
		{docData}
		onsearch={handleSearch}
	/>
{/snippet}

{#snippet caseActions()}
	<CaseActions
		{docData}
		{supabase}
		{session}
		{user}
		url={pdfViewerUrl}
	/>
{/snippet}

{#snippet pdfViewer()}
	<DocumentPDFViewer
		{docData}
		{session}
		onloaded={handlePDFLoaded}
		onerror={handlePDFError}
	/>
{/snippet}