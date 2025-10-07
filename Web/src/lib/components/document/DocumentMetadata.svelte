<!-- DocumentMetadata.svelte -->
<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { 
		formatDate, 
		formatDateTime, 
		formatFileSize, 
		formatPercentage, 
		formatCount, 
		formatDocumentType, 
		formatPersonName, 
		formatCourtName, 
		formatCaseInfo, 
		isEmptyValue 
	} from '$lib/utils/utils';
	import type { Document } from '$lib/types';
	import { slide } from 'svelte/transition';

	const {
		docData
	}: {
		docData?: Document | null;
	} = $props();

	// Section collapse state
	let sectionsCollapsed = $state({
		document: false,
		case: false,
		content: false,
		properties: false,
		legal: false
	});

	const dispatch = createEventDispatcher<{
		search: {
			field: string;
			value: string;
		};
	}>();

	// Section management functions
	function toggleSection(section: keyof typeof sectionsCollapsed) {
		sectionsCollapsed[section] = !sectionsCollapsed[section];
	}

	function getSectionIcon(collapsed: boolean): string {
		return collapsed ? 'M9 5l7 7-7 7' : 'M19 9l-7 7-7-7';
	}

	// Helper function to create search queries
	function createSearchForField(field: string, value: string) {
		if (!value) return;
		
		// Dispatch search event that can be caught by parent component
		dispatch('search', {
			field,
			value
		});
	}
</script>

{#if docData}
	<div class="space-y-4 overflow-hidden">
		<h3 class="text-sm font-medium text-neutral-700">Document Details</h3>

		<!-- Document Information Section -->
		<div class="rounded-lg border border-neutral-200 bg-white">
			<button
				type="button"
				onclick={() => toggleSection('document')}
				class="flex w-full items-center justify-between p-3 text-left hover:bg-neutral-50"
			>
				<h4 class="text-sm font-medium text-neutral-800">Document Information</h4>
				<svg
					class="h-4 w-4 transform transition-transform {sectionsCollapsed.document ? '' : 'rotate-180'}"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getSectionIcon(sectionsCollapsed.document)} />
				</svg>
			</button>
			{#if !sectionsCollapsed.document}
				<div class="space-y-3 border-t border-neutral-100 p-3" transition:slide={{ duration: 300 }}>
					<div>
						<p class="text-xs text-neutral-500">Document Name</p>
						<p class="text-sm font-medium text-neutral-800">{docData.metadata?.document_name || docData.file_name || 'Untitled Document'}</p>
					</div>

					{#if !isEmptyValue(docData.metadata?.subject)}
						<div>
							<p class="text-xs text-neutral-500">Subject</p>
							<p class="text-sm font-medium text-neutral-800">{docData.metadata?.subject}</p>
						</div>
					{/if}

					<div>
						<p class="text-xs text-neutral-500">Document Type</p>
						<span class="inline-flex rounded-full bg-primary-100 px-2 py-1 text-xs font-medium text-primary-800">
							{formatDocumentType(docData.doc_type)}
						</span>
					</div>

					{#if !isEmptyValue(docData.metadata?.summary)}
						<div>
							<p class="text-xs text-neutral-500">Summary</p>
							<p class="mt-1 text-sm text-neutral-700 leading-relaxed">{docData.metadata.summary}</p>
						</div>
					{/if}

					<!-- Legal Classification Information -->
					{#if docData.metadata?.legal_tags && docData.metadata.legal_tags.length > 0}
						<div>
							<p class="text-xs text-neutral-500">Legal Tags</p>
							<div class="mt-1 flex flex-wrap gap-1">
								{#each docData.metadata?.legal_tags || [] as tag, i}
									<button
										type="button"
										onclick={() => createSearchForField('legal_tags', tag)}
										class="inline-flex rounded-full bg-neutral-100 px-2 py-0.5 text-xs font-medium text-neutral-800 hover:bg-neutral-200 transition-colors"
									>
										{tag}
									</button>
								{/each}
							</div>
						</div>
					{/if}

					{#if docData.metadata?.authorities}
						<div>
							<p class="text-xs text-neutral-500">Authorities</p>
							<div class="mt-1 flex flex-wrap gap-1">
								{#if Array.isArray(docData.metadata.authorities)}
									{#each docData.metadata.authorities as authority, i}
										<span class="inline-flex rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-800">
											{typeof authority === 'object' ? authority.citation || authority.case_title : authority}
										</span>
									{/each}
								{:else}
									<span class="text-sm font-medium text-neutral-800">
										{typeof docData.metadata.authorities === 'object' ? JSON.stringify(docData.metadata.authorities) : docData.metadata.authorities}
									</span>
								{/if}
							</div>
						</div>
					{/if}

					{#if docData.metadata?.charges && Array.isArray(docData.metadata.charges) && docData.metadata.charges.length > 0}
						<div>
							<p class="text-xs text-neutral-500">Charges</p>
							<div class="mt-1 space-y-1">
								{#each docData.metadata.charges as charge}
									<div class="text-sm">
										<span class="font-medium text-neutral-800">{charge.statute}</span>
										{#if charge.description}
											<div class="text-xs text-neutral-600">{charge.description}</div>
										{/if}
										{#if charge.grade || charge.class}
											<div class="text-xs text-neutral-500">
												{charge.grade ? `Grade: ${charge.grade}` : ''}
												{charge.class ? `Class: ${charge.class}` : ''}
											</div>
										{/if}
									</div>
								{/each}
							</div>
						</div>
					{/if}

					{#if docData.metadata?.sensitive_terms && docData.metadata.sensitive_terms.length > 0}
						<div>
							<p class="text-xs text-neutral-500">Sensitive Terms</p>
							<div class="mt-1 flex flex-wrap gap-1">
								{#each docData.metadata?.sensitive_terms || [] as term, i}
									<span class="inline-flex rounded-full bg-red-100 px-2 py-0.5 text-xs font-medium text-red-800">
										{term}
									</span>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			{/if}
		</div>

		<!-- Case & Legal Information Section -->
		{#if docData.metadata?.case || docData.metadata?.case_number || docData.metadata?.case_name || docData.metadata?.court || docData.metadata?.judge}
			<div class="rounded-lg border border-neutral-200 bg-white">
				<button
					type="button"
					onclick={() => toggleSection('case')}
					class="flex w-full items-center justify-between p-3 text-left hover:bg-neutral-50"
				>
					<h4 class="text-sm font-medium text-neutral-800">Case & Legal Information</h4>
					<svg
						class="h-4 w-4 transform transition-transform {sectionsCollapsed.case ? '' : 'rotate-180'}"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getSectionIcon(sectionsCollapsed.case)} />
					</svg>
				</button>
				{#if !sectionsCollapsed.case}
					<div class="space-y-3 border-t border-neutral-100 p-3" transition:slide={{ duration: 300 }}>
						{#if docData.metadata?.case || docData.metadata?.case_number || docData.metadata?.case_name}
							{@const caseInfo = formatCaseInfo(docData.metadata?.case)}
							
							{#if !isEmptyValue(caseInfo.number) || !isEmptyValue(docData.metadata?.case_number)}
							<div>
								<p class="text-xs text-neutral-500">Case Number</p>
								<button
									type="button"
									onclick={() => createSearchForField('case_number', caseInfo.number || docData.metadata?.case_number)}
									class="text-sm font-medium text-primary-600 hover:text-primary-800 hover:underline"
								>
									{caseInfo.number || docData.metadata?.case_number}
								</button>
							</div>
						{/if}

						{#if !isEmptyValue(caseInfo.name) || !isEmptyValue(docData.metadata?.case_name)}
							<div>
								<p class="text-xs text-neutral-500">Case Name</p>
								<button
									type="button"
									onclick={() => createSearchForField('case_name', caseInfo.name || docData.metadata?.case_name)}
									class="text-sm font-medium text-primary-600 hover:text-primary-800 hover:underline"
								>
									{caseInfo.name || docData.metadata?.case_name}
								</button>
							</div>
						{/if}

						{#if !isEmptyValue(docData.metadata?.court)}
							<div>
								<p class="text-xs text-neutral-500">Court</p>
								<button
									type="button"
									onclick={() => createSearchForField('court', formatCourtName(docData.metadata?.court))}
									class="text-sm font-medium text-primary-600 hover:text-primary-800 hover:underline"
								>
									{formatCourtName(docData.metadata?.court)}
								</button>
							</div>
						{/if}

						{#if !isEmptyValue(docData.metadata?.judge)}
							<div>
								<p class="text-xs text-neutral-500">Judge</p>
								<button
									type="button"
									onclick={() => createSearchForField('judge', formatPersonName(docData.metadata?.judge))}
									class="text-sm font-medium text-primary-600 hover:text-primary-800 hover:underline"
								>
									{formatPersonName(docData.metadata?.judge)}
								</button>
							</div>
						{/if}

						{#if !isEmptyValue(docData.metadata?.author)}
							<div>
								<p class="text-xs text-neutral-500">Author</p>
								<button
									type="button"
									onclick={() => createSearchForField('author', docData.metadata?.author)}
									class="text-sm font-medium text-primary-600 hover:text-primary-800 hover:underline"
								>
									{docData.metadata?.author}
								</button>
							</div>
						{/if}

						{#if !isEmptyValue(docData.metadata?.status)}
							<div>
								<p class="text-xs text-neutral-500">Status</p>
								<span class="inline-flex rounded-full bg-blue-100 px-2 py-1 text-xs font-medium text-blue-800">
									{docData.metadata?.status}
								</span>
							</div>
						{/if}

						<!-- Parties Information -->
						{#if docData.metadata?.parties && Array.isArray(docData.metadata.parties) && docData.metadata.parties.length > 0}
							<div>
								<p class="text-xs text-neutral-500">Parties</p>
								<div class="mt-1 space-y-2">
									{#each docData.metadata.parties as party}
										<div class="text-sm">
											<div class="font-medium text-neutral-800 break-words">{party.name}</div>
											{#if party.role}
												<div class="text-xs text-neutral-500 mt-0.5">{party.role}</div>
											{/if}
										</div>
									{/each}
								</div>
							</div>
						{/if}

						<!-- Attorneys Information -->
						{#if docData.metadata?.attorneys && Array.isArray(docData.metadata.attorneys) && docData.metadata.attorneys.length > 0}
							<div>
								<p class="text-xs text-neutral-500">Attorneys</p>
								<div class="mt-1 space-y-2">
									{#each docData.metadata.attorneys as attorney}
										<div class="text-sm">
											<div class="font-medium text-neutral-800 break-words">{attorney.name}</div>
											{#if attorney.role}
												<div class="text-xs text-neutral-500 mt-0.5">{attorney.role}</div>
											{/if}
											{#if attorney.organization}
												<div class="text-xs text-neutral-500 mt-0.5 break-words">{attorney.organization}</div>
											{/if}
										</div>
									{/each}
								</div>
							</div>
						{/if}
						{/if}
					</div>
				{/if}
			</div>
		{/if}

		<!-- Enhanced Content Analysis & Privacy Section -->
		{#if docData.metadata?.classification_confidence !== undefined || docData.metadata?.has_redactions !== undefined || docData.metadata?.redaction_score !== undefined || !isEmptyValue(docData.metadata?.extraction_method) || docData.metadata?.sensitive_terms}
		<div class="rounded-lg border border-neutral-200 bg-white">
			<button
				type="button"
				onclick={() => toggleSection('content')}
				class="flex w-full items-center justify-between p-3 text-left hover:bg-neutral-50"
			>
				<h4 class="text-sm font-medium text-neutral-800">Content Analysis & Privacy</h4>
				<svg
					class="h-4 w-4 transform transition-transform {sectionsCollapsed.content ? '' : 'rotate-180'}"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getSectionIcon(sectionsCollapsed.content)} />
				</svg>
			</button>
			{#if !sectionsCollapsed.content}
				<div class="space-y-4 border-t border-neutral-100 p-3" transition:slide={{ duration: 300 }}>
					<!-- Privacy & Security Status -->
					{#if docData.metadata?.has_redactions !== undefined || docData.metadata?.sensitive_terms || docData.metadata?.redaction_score !== undefined}
						<div class="rounded-lg bg-neutral-50 p-3">
							<p class="text-xs font-medium text-neutral-700 mb-2">🔐 Privacy & Security Status</p>
							
							{#if docData.metadata?.has_redactions !== undefined}
								<div class="mb-2">
									<span class={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${docData.metadata?.has_redactions ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'}`}>
										{docData.metadata?.has_redactions ? '🔒 Contains Redactions' : '✅ No Redactions'}
									</span>
								</div>
							{/if}

							{#if docData.metadata?.sensitive_terms && docData.metadata.sensitive_terms.length > 0}
								<div class="mb-2">
									<div class="text-xs text-neutral-600 mb-1">Sensitive Terms Detected:</div>
									<div class="flex flex-wrap gap-1">
										{#each docData.metadata.sensitive_terms.slice(0, 5) as term}
											<span class="inline-flex rounded-full bg-yellow-100 px-2 py-0.5 text-xs font-medium text-yellow-800">
												⚠️ {term}
											</span>
										{/each}
										{#if docData.metadata.sensitive_terms.length > 5}
											<span class="inline-flex rounded-full bg-neutral-100 px-2 py-0.5 text-xs font-medium text-neutral-600">
												+{docData.metadata.sensitive_terms.length - 5} more
											</span>
										{/if}
									</div>
								</div>
							{/if}

							{#if docData.metadata?.redaction_score !== undefined}
								<div>
									<div class="flex items-center justify-between text-xs text-neutral-600 mb-1">
										<span>Privacy Score:</span>
										<span class="font-medium">{formatPercentage(docData.metadata.redaction_score)}</span>
									</div>
									<div class="flex items-center space-x-2">
										<div class="flex-1 bg-neutral-200 rounded-full h-2">
											<div 
												class={`h-2 rounded-full ${docData.metadata.redaction_score > 0.7 ? 'bg-red-500' : docData.metadata.redaction_score > 0.3 ? 'bg-yellow-500' : 'bg-green-500'}`}
												style="width: {formatPercentage(docData.metadata.redaction_score).replace('%', '')}%"
											></div>
										</div>
									</div>
								</div>
							{/if}
						</div>
					{/if}

					<!-- Classification Analysis -->
					{#if docData.metadata?.classification_confidence !== undefined}
						<div>
							<p class="text-xs text-neutral-500">AI Classification Confidence</p>
							<div class="flex items-center space-x-2">
								<div class="flex-1 bg-neutral-200 rounded-full h-2">
									<div 
										class="bg-blue-600 h-2 rounded-full" 
										style="width: {formatPercentage(docData.metadata.classification_confidence).replace('%', '')}%"
									></div>
								</div>
								<span class="text-sm font-medium text-neutral-800">{formatPercentage(docData.metadata.classification_confidence)}</span>
							</div>
						</div>
					{/if}

					{#if !isEmptyValue(docData.metadata?.extraction_method)}
						<div>
							<p class="text-xs text-neutral-500">Text Extraction Method</p>
							<span class="inline-flex rounded-full bg-purple-100 px-2 py-1 text-xs font-medium text-purple-800">
								{docData.metadata?.extraction_method}
							</span>
						</div>
					{/if}
				</div>
			{/if}
		</div>
		{/if}

		<!-- Document Properties Section -->
		<div class="rounded-lg border border-neutral-200 bg-white">
			<button
				type="button"
				onclick={() => toggleSection('properties')}
				class="flex w-full items-center justify-between p-3 text-left hover:bg-neutral-50"
			>
				<h4 class="text-sm font-medium text-neutral-800">Document Properties</h4>
				<svg
					class="h-4 w-4 transform transition-transform {sectionsCollapsed.properties ? '' : 'rotate-180'}"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getSectionIcon(sectionsCollapsed.properties)} />
				</svg>
			</button>
			{#if !sectionsCollapsed.properties}
				<div class="space-y-3 border-t border-neutral-100 p-3" transition:slide={{ duration: 300 }}>
					{#if docData.size || docData.metadata?.size}
						<div>
							<p class="text-xs text-neutral-500">File Size</p>
							<p class="text-sm font-medium text-neutral-800">{formatFileSize(docData.size || docData.metadata?.size)}</p>
						</div>
					{/if}

					{#if !isEmptyValue(docData.content_type) || !isEmptyValue(docData.metadata?.file_type)}
						<div>
							<p class="text-xs text-neutral-500">File Type</p>
							<p class="text-sm font-medium text-neutral-800">{docData.content_type || docData.metadata?.file_type}</p>
						</div>
					{/if}

					{#if !isEmptyValue(docData.metadata?.word_count)}
						<div>
							<p class="text-xs text-neutral-500">Word Count</p>
							<p class="text-sm font-medium text-neutral-800">{formatCount(docData.metadata?.word_count)}</p>
						</div>
					{/if}

					{#if !isEmptyValue(docData.metadata?.page_count) || !isEmptyValue(docData.metadata?.pages)}
						<div>
							<p class="text-xs text-neutral-500">Page Count</p>
							<p class="text-sm font-medium text-neutral-800">{docData.metadata?.page_count || docData.metadata?.pages}</p>
						</div>
					{/if}

					{#if !isEmptyValue(docData.metadata?.language)}
						<div>
							<p class="text-xs text-neutral-500">Language</p>
							<p class="text-sm font-medium text-neutral-800">{docData.metadata?.language}</p>
						</div>
					{/if}

					<!-- AI Classification and Confidence Score on same line -->
					{#if docData.metadata?.ai_classified !== undefined || docData.metadata?.confidence !== undefined}
						<div>
							<p class="text-xs text-neutral-500 mb-2">AI Classification</p>
							<div class="flex items-center justify-between space-x-3">
								{#if docData.metadata?.ai_classified !== undefined}
									<div class="flex items-center space-x-2">
										<span class="text-xs text-neutral-500">Status:</span>
										<span class={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${docData.metadata?.ai_classified ? 'bg-green-100 text-green-800' : 'bg-neutral-100 text-neutral-800'}`}>
											{docData.metadata?.ai_classified ? 'Yes' : 'No'}
										</span>
									</div>
								{/if}
								{#if docData.metadata?.confidence !== undefined}
									<div class="flex items-center space-x-2 flex-1">
										<span class="text-xs text-neutral-500 whitespace-nowrap">Confidence:</span>
										<div class="flex items-center space-x-2 flex-1">
											<div class="flex-1 bg-neutral-200 rounded-full h-2 min-w-16">
												<div 
													class="bg-primary-600 h-2 rounded-full" 
													style="width: {formatPercentage(docData.metadata.confidence).replace('%', '')}%"
												></div>
											</div>
											<span class="text-xs font-medium text-neutral-800 whitespace-nowrap">{formatPercentage(docData.metadata.confidence)}</span>
										</div>
									</div>
								{/if}
							</div>
						</div>
					{/if}

					<!-- Enhanced Legal Date Information -->
					{#if !isEmptyValue(docData.metadata?.filing_date)}
						<div>
							<p class="text-xs text-neutral-500">Filing Date</p>
							<span class="inline-flex items-center rounded bg-blue-50 px-2 py-1 text-xs font-medium text-blue-700">
								📅 {formatDate(docData.metadata?.filing_date)}
							</span>
						</div>
					{/if}

					{#if !isEmptyValue(docData.metadata?.event_date)}
						<div>
							<p class="text-xs text-neutral-500">Event Date</p>
							<span class="inline-flex items-center rounded bg-green-50 px-2 py-1 text-xs font-medium text-green-700">
								⚖️ {formatDate(docData.metadata?.event_date)}
							</span>
						</div>
					{/if}

					{#if !isEmptyValue(docData.metadata?.hearing_date)}
						<div>
							<p class="text-xs text-neutral-500">Hearing Date</p>
							<span class="inline-flex items-center rounded bg-purple-50 px-2 py-1 text-xs font-medium text-purple-700">
								🏛️ {formatDate(docData.metadata?.hearing_date)}
							</span>
						</div>
					{/if}

					{#if !isEmptyValue(docData.metadata?.decision_date)}
						<div>
							<p class="text-xs text-neutral-500">Decision Date</p>
							<span class="inline-flex items-center rounded bg-orange-50 px-2 py-1 text-xs font-medium text-orange-700">
								⚡ {formatDate(docData.metadata?.decision_date)}
							</span>
						</div>
					{/if}

					{#if !isEmptyValue(docData.metadata?.served_date)}
						<div>
							<p class="text-xs text-neutral-500">Served Date</p>
							<span class="inline-flex items-center rounded bg-indigo-50 px-2 py-1 text-xs font-medium text-indigo-700">
								📬 {formatDate(docData.metadata?.served_date)}
							</span>
						</div>
					{/if}

					{#if !isEmptyValue(docData.metadata?.timestamp)}
						<div>
							<p class="text-xs text-neutral-500">Document Date</p>
							<p class="text-sm font-medium text-neutral-800">{formatDate(docData.metadata?.timestamp)}</p>
						</div>
					{/if}

					<div>
						<p class="text-xs text-neutral-500">Created</p>
						<p class="text-sm font-medium text-neutral-800">{formatDateTime(docData.created_at)}</p>
					</div>

					{#if !isEmptyValue(docData.updated_at)}
						<div>
							<p class="text-xs text-neutral-500">Updated</p>
							<p class="text-sm font-medium text-neutral-800">{formatDateTime(docData.updated_at)}</p>
						</div>
					{/if}

					{#if !isEmptyValue(docData.metadata?.processed_at)}
						<div>
							<p class="text-xs text-neutral-500">Processed</p>
							<p class="text-sm font-medium text-neutral-800">{formatDateTime(docData.metadata?.processed_at)}</p>
						</div>
					{/if}
				</div>
			{/if}
		</div>
	</div>
{/if}