<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { Button, LoadingSpinner } from '$lib/components/ui';
	import { slide, scale } from 'svelte/transition';
	import { cubicOut, elasticOut } from 'svelte/easing';

	const dispatch = createEventDispatcher();

	interface RedactionData {
		type: string;
		text: string;
		reason?: string;
		legal_code?: string;
	}

	interface RedactionAnalysis {
		redactions_found: number;
		suggested_redactions?: RedactionData[];
		analysis_summary?: string;
	}

	interface Props {
		analysis: RedactionAnalysis | null;
		documentId?: string;
		isCreatingRedacted?: boolean;
		visible?: boolean;
		class?: string;
	}

	let {
		analysis = null,
		documentId = '',
		isCreatingRedacted = false,
		visible = false,
		class: className = ''
	}: Props = $props();

	// Handle creating redacted document
	function handleCreateRedacted() {
		dispatch('createRedacted', { documentId });
	}

	// Handle dismissing the analysis
	function handleDismiss() {
		dispatch('dismiss');
	}

	// Get redaction type styling
	function getRedactionTypeStyling(type: string): string {
		const typeMap = {
			'personal_info': 'bg-red-100 text-red-800 border-red-200',
			'financial': 'bg-orange-100 text-orange-800 border-orange-200',
			'legal': 'bg-purple-100 text-purple-800 border-purple-200',
			'medical': 'bg-blue-100 text-blue-800 border-blue-200',
			'sensitive': 'bg-yellow-100 text-yellow-800 border-yellow-200'
		};
		
		const normalizedType = type.toLowerCase().replace(/[^a-z]/g, '');
		return typeMap[normalizedType] || 'bg-neutral-100 text-neutral-800 border-neutral-200';
	}
</script>

{#if visible && analysis}
	<div 
		class="bg-orange-50 border border-orange-200 rounded-lg p-6 {className}"
		in:slide={{ duration: 600, delay: 200, easing: cubicOut }}
		out:slide={{ duration: 400, easing: cubicOut }}
	>
		<!-- Header -->
		<div class="flex items-start justify-between mb-4">
			<div class="flex items-center space-x-2">
				<div class="flex-shrink-0">
					<svg class="h-6 w-6 text-orange-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
					</svg>
				</div>
				<h3 class="text-lg font-semibold text-orange-800">
					⚠️ Redaction Analysis
				</h3>
			</div>
			
			<!-- Close button -->
			<button
				type="button"
				onclick={handleDismiss}
				class="text-orange-400 hover:text-orange-600 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2 focus:ring-offset-orange-50 rounded-full p-1"
				aria-label="Dismiss redaction analysis"
			>
				<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>
		</div>

		<!-- Summary -->
		<div class="mb-4">
			<p class="text-orange-800 font-medium">
				Found {analysis.redactions_found} potential redaction{analysis.redactions_found !== 1 ? 's' : ''}
			</p>
			{#if analysis.analysis_summary}
				<p class="text-sm text-orange-700 mt-1">
					{analysis.analysis_summary}
				</p>
			{/if}
		</div>

		<!-- Redaction Details -->
		{#if analysis.suggested_redactions && analysis.suggested_redactions.length > 0}
			<div class="mb-6">
				<p class="text-sm text-orange-700 mb-3 font-medium">
					The following sensitive information was detected and should be redacted:
				</p>
				
				<div class="max-h-64 overflow-y-auto space-y-2">
					{#each analysis.suggested_redactions as redaction, index}
						<div 
							class="bg-white border rounded-lg p-3 {getRedactionTypeStyling(redaction.type)}"
							in:scale={{ start: 0.95, duration: 300, delay: index * 100, easing: elasticOut }}
						>
							<div class="flex flex-wrap items-start gap-2">
								<!-- Type badge -->
								<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-white/60">
									{redaction.type}
								</span>
								
								<!-- Redacted text -->
								<code class="flex-1 text-sm font-mono bg-neutral-100 px-2 py-1 rounded border">
									{redaction.text}
								</code>
							</div>
							
							<!-- Additional details -->
							{#if redaction.reason}
								<p class="text-xs mt-2 opacity-80">
									<strong>Reason:</strong> {redaction.reason}
								</p>
							{/if}
							
							{#if redaction.legal_code}
								<p class="text-xs mt-1 opacity-70">
									<strong>Legal Code:</strong> {redaction.legal_code}
								</p>
							{/if}
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<!-- Action Buttons -->
		<div class="flex flex-col sm:flex-row gap-3">
			<Button
				variant="warning"
				onclick={handleCreateRedacted}
				disabled={isCreatingRedacted || !documentId}
				loading={isCreatingRedacted}
				class="flex-1"
			>
				{#if isCreatingRedacted}
					Creating Redacted Version...
				{:else}
					Create Redacted Document
				{/if}
			</Button>

			<Button
				variant="secondary"
				onclick={handleDismiss}
				class="sm:w-auto"
			>
				Dismiss
			</Button>
		</div>

		<!-- Warning Notice -->
		<div class="mt-4 p-3 bg-orange-100 border border-orange-200 rounded-md">
			<p class="text-xs text-orange-700">
				<strong>Important:</strong> This analysis is automated and may not catch all sensitive information. 
				Please review the document manually before sharing or publishing.
			</p>
		</div>
	</div>
{/if}