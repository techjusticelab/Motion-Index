<script lang="ts">
	import { page } from '$app/stores';
	import { CaseManager, type Case } from '$lib/supabase';
	import { onMount } from 'svelte';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import { cubicOut, quintOut, backOut, elasticOut } from 'svelte/easing';

	let { data } = $props();
	let { session, user, supabase } = $derived(data);

	// Flag to control animations after initial page load
	let isInitialLoad = $state(true);

	// Case management data
	let cases = $state<Case[]>([]);
	let caseManager: CaseManager;
	let isLoadingCases = $state(true);
	let showNewCaseModal = $state(false);
	let newCaseName = $state('');
	let isCreatingCase = $state(false);

	onMount(() => {
		caseManager = new CaseManager(supabase);

		// Set initial load to false after the first render
		setTimeout(() => {
			isInitialLoad = false;
		}, 100);

		// If user is logged in, load their cases
		if (user) {
			loadUserCases();
		} else {
			// If no session, stop loading
			isLoadingCases = false;
		}
	});

	// Load user's cases
	async function loadUserCases() {
		if (!user) {
			isLoadingCases = false;
			return;
		}
		
		isLoadingCases = true;
		try {
			cases = await caseManager.getUserCases(user.id);
		} catch (error) {
			console.error('Error loading cases:', error);
		} finally {
			isLoadingCases = false;
		}
	}

	// Open case creation modal
	function openNewCaseModal() {
		newCaseName = '';
		showNewCaseModal = true;
	}

	// Close case creation modal
	function closeNewCaseModal() {
		showNewCaseModal = false;
		newCaseName = '';
	}

	// Create new case
	async function createCase() {
		if (!user || !newCaseName.trim()) return;
		
		isCreatingCase = true;
		try {
			const newCase = await caseManager.createCase(user.id, newCaseName.trim());
			if (newCase) {
				cases = [newCase, ...cases];
				closeNewCaseModal();
			}
		} catch (error) {
			console.error('Error creating case:', error);
		} finally {
			isCreatingCase = false;
		}
	}

	// Delete case
	async function deleteCase(caseId: string) {
		if (!confirm('Are you sure you want to delete this case and all its documents?')) return;
		
		const success = await caseManager.deleteCase(caseId);
		if (success) {
			cases = cases.filter(c => c.id !== caseId);
		}
	}

	// Navigate to case detail
	function viewCase(caseId: string) {
		goto(`/cases/${caseId}`);
	}

	// Format date for display
	function formatDate(dateString: string) {
		if (!dateString) return 'N/A';
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	// Handle keydown for modal escape
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && showNewCaseModal) {
			closeNewCaseModal();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="min-h-screen bg-neutral-50 py-8">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Header -->
		<div 
			class="mb-8"
			in:fly={{ y: -30, duration: 800, easing: quintOut, delay: isInitialLoad ? 0 : 0 }}
		>
			<div class="flex items-center justify-between">
				<div>
					<h1 class="text-3xl font-bold text-neutral-900">Case Management</h1>
					<p class="mt-2 text-sm text-neutral-600">
						Organize and manage your legal cases and documents
					</p>
				</div>
				<button
					onclick={openNewCaseModal}
					class="inline-flex items-center justify-center rounded-md bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mr-2 h-4 w-4"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
					</svg>
					New Case
				</button>
			</div>
		</div>

		<!-- Cases Grid -->
		{#if isLoadingCases}
			<div 
				class="flex items-center justify-center py-12"
				in:fade={{ duration: 300 }}
			>
				<div class="flex items-center space-x-3">
					<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-200 border-t-indigo-600"></div>
					<span class="text-neutral-600">Loading cases...</span>
				</div>
			</div>
		{:else if cases.length === 0}
			<div
				class="rounded-lg border border-dashed border-neutral-300 bg-white p-12 text-center"
				in:fly={{ y: 20, duration: 600, delay: isInitialLoad ? 300 : 0, easing: cubicOut }}
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mx-auto h-12 w-12 text-neutral-400"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
					/>
				</svg>
				<h3 class="mt-4 text-lg font-medium text-neutral-900">No cases yet</h3>
				<p class="mt-2 text-sm text-neutral-500">
					Create your first case to start organizing your legal documents
				</p>
				<button
					onclick={openNewCaseModal}
					class="mt-4 inline-flex items-center justify-center rounded-md bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 focus:outline-none"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mr-2 h-4 w-4"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
					</svg>
					Create Your First Case
				</button>
			</div>
		{:else}
			<div class="max-h-96 overflow-y-auto rounded-lg border border-neutral-200 bg-white">
				<div class="divide-y divide-neutral-200">
					{#each cases as caseItem, i}
						<div
							class="group relative p-4 hover:bg-neutral-50 transition-colors cursor-pointer"
							in:fly={{ y: 20, duration: 400, delay: i * 50 }}
							onclick={() => viewCase(caseItem.id)}
						>
							<div class="flex items-center justify-between">
								<div class="flex items-center space-x-3 flex-1 min-w-0">
									<!-- Case Icon -->
									<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary-100 flex-shrink-0">
										<svg
											class="h-4 w-4 text-primary-600"
											fill="none"
											stroke="currentColor"
											viewBox="0 0 24 24"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
											/>
										</svg>
									</div>

									<!-- Case Info -->
									<div class="flex-1 min-w-0">
										<h3 class="text-sm font-medium text-neutral-900 group-hover:text-primary-600 break-words">
											{caseItem.case_name}
										</h3>
										<div class="mt-1 flex items-center space-x-4 text-xs text-neutral-500">
											<span>Created {formatDate(caseItem.created_at)}</span>
											<span>Updated {formatDate(caseItem.updated_at)}</span>
										</div>
									</div>
								</div>

								<!-- Actions -->
								<div class="flex items-center space-x-2 flex-shrink-0">
									<button
										onclick={(e) => { e.stopPropagation(); viewCase(caseItem.id); }}
										class="inline-flex items-center text-xs font-medium text-primary-600 hover:text-primary-800"
									>
										View
										<svg class="ml-1 h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
										</svg>
									</button>
									<button
										onclick={(e) => { e.stopPropagation(); deleteCase(caseItem.id); }}
										class="text-xs text-red-600 hover:text-red-800"
									>
										Delete
									</button>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	</div>
</div>

<!-- New Case Modal -->
{#if showNewCaseModal}
	<div 
		class="fixed inset-0 z-50 flex items-center justify-center p-4"
		style="background-color: rgba(0, 0, 0, 0.5);"
		in:fade={{ duration: 300, easing: cubicOut }}
		out:fade={{ duration: 200 }}
	>
		<div 
			class="w-full max-w-md rounded-lg bg-white p-6 shadow-xl"
			in:scale={{ start: 0.95, duration: 300, easing: cubicOut }}
			out:scale={{ start: 1, end: 0.95, duration: 200 }}
		>
			<div class="mb-4 flex items-center justify-between">
				<h3 class="text-xl font-medium text-neutral-900">Create New Case</h3>
				<button
					onclick={closeNewCaseModal}
					class="rounded-full p-1 text-neutral-400 hover:bg-neutral-100 hover:text-neutral-500"
					aria-label="Close modal"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-6 w-6"
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

			<div class="mb-4">
				<label for="case-name" class="block text-sm font-medium text-neutral-700">Case Name</label>
				<input
					type="text"
					id="case-name"
					bind:value={newCaseName}
					class="mt-1 block w-full rounded-md border border-neutral-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
					placeholder="Enter a name for your case"
					onkeydown={(e) => e.key === 'Enter' && createCase()}
					required
				/>
				<p class="mt-1 text-xs text-neutral-500">
					Give your case a descriptive name to help you identify it later.
				</p>
			</div>

			<div class="mt-5 flex justify-end space-x-3">
				<button
					type="button"
					onclick={closeNewCaseModal}
					class="inline-flex justify-center rounded-md border border-neutral-300 bg-white px-4 py-2 text-sm font-medium text-neutral-700 shadow-sm hover:bg-neutral-50 focus:outline-none"
				>
					Cancel
				</button>
				<button
					type="button"
					class="inline-flex justify-center rounded-md bg-primary-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
					disabled={!newCaseName.trim() || isCreatingCase}
					onclick={createCase}
				>
					{#if isCreatingCase}
						<div class="flex items-center">
							<div class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white"></div>
							<span>Creating...</span>
						</div>
					{:else}
						Create Case
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}