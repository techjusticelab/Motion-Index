<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { useForm } from '$lib/hooks/useForm';
	import { Input, FormField } from '$lib/components/forms';
	import { Button, Alert } from '$lib/components/ui';
	import type { Document } from '$lib/types';
	import { fly, scale } from 'svelte/transition';
	import { cubicOut, backOut } from 'svelte/easing';

	const dispatch = createEventDispatcher();

	interface Props {
		document?: Document | null;
		isUpdating?: boolean;
		class?: string;
	}

	let {
		document = null,
		isUpdating = false,
		class: className = ''
	}: Props = $props();

	// Form initial values
	const initialMetadata = {
		document_name: '',
		subject: '',
		status: '',
		timestamp: '',
		case_name: '',
		case_number: '',
		author: '',
		judge: '',
		court: ''
	};

	// Legal tags state
	let legalTags = $state<string[]>([]);
	let tagInput = $state('');

	// Form validation rules
	const validationRules = {
		document_name: { required: true, minLength: 2 },
		subject: { required: false },
		status: { required: false },
		timestamp: { required: false },
		case_name: { required: false },
		case_number: { required: false },
		author: { required: false },
		judge: { required: false },
		court: { required: false }
	};

	// Initialize form
	const form = useForm(initialMetadata, validationRules);

	// Update form when document changes
	$effect(() => {
		if (document?.metadata) {
			const metadata = document.metadata;
			form.setFields({
				document_name: metadata.document_name || '',
				subject: metadata.subject || '',
				status: metadata.status || '',
				timestamp: metadata.timestamp || '',
				case_name: metadata.case_name || '',
				case_number: metadata.case_number || '',
				author: metadata.author || '',
				judge: metadata.judge || '',
				court: metadata.court || ''
			});
			legalTags = [...(metadata.legal_tags || [])];
		}
	});

	// Format field labels
	function formatLabel(key: string): string {
		return key.replace(/_/g, ' ').replace(/(?:^|\s)\S/g, (a) => a.toUpperCase());
	}

	// Tag management
	function addTag() {
		const trimmedTag = tagInput.trim();
		if (trimmedTag && !legalTags.includes(trimmedTag)) {
			legalTags = [...legalTags, trimmedTag];
			tagInput = '';
		}
	}

	function removeTag(index: number) {
		legalTags = legalTags.filter((_, i) => i !== index);
	}

	function handleTagKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			event.preventDefault();
			addTag();
		}
	}

	// Submit handler
	async function handleSubmit(values: typeof initialMetadata) {
		const metadataToSave = {
			...values,
			legal_tags: legalTags
		};

		dispatch('save', {
			documentId: document?.id,
			metadata: metadataToSave
		});
	}

	// Success/error state
	let showSuccess = $state(false);
	let showError = $state(false);
	let statusMessage = $state('');

	// Watch for external update status
	let previousIsUpdating = isUpdating;
	$effect(() => {
		// If updating changes from true to false, show result
		if (previousIsUpdating && !isUpdating) {
			// This would need to be passed in as props or events
			// For now, just show generic success
			showSuccess = true;
			statusMessage = 'Metadata updated successfully!';
			setTimeout(() => {
				showSuccess = false;
				statusMessage = '';
			}, 5000);
		}
		previousIsUpdating = isUpdating;
	});
</script>

<div class="space-y-6 {className}">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<h3 class="text-lg font-semibold text-neutral-900">Document Metadata</h3>
		{#if document}
			<div class="text-sm text-neutral-500">
				ID: {document.id}
			</div>
		{/if}
	</div>

	<!-- Status messages -->
	{#if showSuccess}
		<Alert type="success" dismissible bind:showSuccess>
			{statusMessage}
		</Alert>
	{/if}

	{#if showError}
		<Alert type="error" dismissible bind:showError>
			{statusMessage}
		</Alert>
	{/if}

	{#if document}
		<form class="space-y-4" use:form.handleSubmit={handleSubmit}>
			<!-- Basic metadata fields -->
			{#each Object.entries($form.values) as [key, value], i}
				<div in:fly={{ y: 10, delay: i * 50, duration: 300, easing: cubicOut }}>
					<FormField
						label={formatLabel(key)}
						error={$form.errors[key]}
						required={validationRules[key]?.required}
					>
						{#snippet children(fieldId)}
							<Input
								id={fieldId}
								name={key}
								bind:value={$form.values[key]}
								error={$form.errors[key]}
								disabled={isUpdating}
								placeholder={`Enter ${formatLabel(key).toLowerCase()}`}
								on:blur={() => form.setTouched(key)}
							/>
						{/snippet}
					</FormField>
				</div>
			{/each}

			<!-- Legal tags section -->
			<div
				class="space-y-3"
				in:fly={{
					y: 10,
					delay: Object.keys($form.values).length * 50,
					duration: 300,
					easing: cubicOut
				}}
			>
				<label class="block text-sm font-medium text-neutral-700">
					Legal Tags
				</label>
				
				<!-- Existing tags -->
				{#if legalTags.length > 0}
					<div class="flex flex-wrap gap-2">
						{#each legalTags as tag, index}
							<div
								class="inline-flex items-center rounded-full bg-primary-100 px-3 py-1.5 text-sm font-medium text-primary-800"
								in:scale={{ start: 0.9, duration: 300, delay: index * 30 }}
							>
								{tag}
								<button
									type="button"
									class="ml-2 inline-flex h-5 w-5 items-center justify-center rounded-full bg-primary-200 text-primary-700 hover:bg-primary-300 focus:outline-none focus:ring-2 focus:ring-primary-500"
									onclick={() => removeTag(index)}
									disabled={isUpdating}
									aria-label="Remove tag {tag}"
								>
									<svg class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
									</svg>
								</button>
							</div>
						{/each}
					</div>
				{/if}

				<!-- Add new tag -->
				<div class="flex gap-2">
					<Input
						bind:value={tagInput}
						placeholder="Add a legal tag"
						disabled={isUpdating}
						onkeydown={handleTagKeydown}
						class="flex-1"
					/>
					<Button
						type="button"
						onclick={addTag}
						disabled={!tagInput.trim() || isUpdating}
						variant="secondary"
					>
						Add
					</Button>
				</div>
			</div>

			<!-- Submit button -->
			<div class="pt-4">
				<Button
					type="submit"
					loading={isUpdating}
					disabled={!$form.isValid || isUpdating}
					class="w-full"
				>
					{isUpdating ? 'Updating...' : 'Update Metadata'}
				</Button>
			</div>
		</form>
	{:else}
		<!-- No document selected state -->
		<div class="rounded-lg border-2 border-dashed border-neutral-300 p-8 text-center">
			<div class="text-neutral-500">
				<svg class="mx-auto h-12 w-12 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
				</svg>
				<p class="text-lg font-medium">No Document Selected</p>
				<p class="text-sm text-neutral-400 mt-1">
					Upload a document to edit its metadata
				</p>
			</div>
		</div>
	{/if}
</div>