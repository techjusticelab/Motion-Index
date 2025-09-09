<script lang="ts">
	import { page } from '$app/stores';
	import { CaseManager, type Case } from '$lib/supabase';
	import { onMount } from 'svelte';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { invalidate, goto } from '$app/navigation';
	import { cubicOut, quintOut, backOut, elasticOut } from 'svelte/easing';
	import { enhance } from '$app/forms';

	let { data } = $props();
	let { session, supabase } = $derived(data);

	// Flag to control animations after initial page load
	let isInitialLoad = $state(true);

	// User data
	let isLoadingUserDetails = $state(false);
	let userDetails: { email: string; user_metadata?: any } | null = null;
	let cases = $state<Case[]>([]);
	let showNewCaseModal = $state(false);
	let newCaseName = $state('');
	let caseManager: CaseManager;
	let caseDocuments = $page.data.caseDocuments || [];
	let isCreatingCase = $state(false);
	let selectedCase: { id: any } | null = null;
	let updateSuccess = $state<boolean | null>(null);
	let updateMessage = $state('');

	// Timer for success/error message
	let messageTimer: ReturnType<typeof setTimeout> | null = null;

	onMount(() => {
		console.log('User session:', session);
		caseManager = new CaseManager(supabase);
		
		// Set initial load to false after the first render
		setTimeout(() => {
			isInitialLoad = false;
		}, 100);

		// Auth state change listener
		const { data } = supabase.auth.onAuthStateChange((_, newSession) => {
			if (newSession?.expires_at !== session?.expires_at) {
				invalidate('supabase:auth');
			}
		});

		// If user is logged in, load their details
		if (session?.user) {
			loadUserDetails();
			loadUserCases();
		}

		return () => data.subscription.unsubscribe();
	});

	async function loadUserDetails() {
		isLoadingUserDetails = true;

		try {
			// Use the session user data directly
			userDetails = {
				email: session.user.email,
				user_metadata: session.user.user_metadata || {}
			};
		} catch (error) {
			console.error('Error loading user details:', error);
		} finally {
			isLoadingUserDetails = false;
		}
	}

	// Get user display name from metadata
	function getUserDisplayName() {
		if (!userDetails) return 'User';

		const metadata = userDetails.user_metadata;
		if (!metadata) return 'User';

		// Check for various name fields that might exist in metadata
		return (
			metadata.full_name || metadata.name || metadata.display_name || userDetails.email || 'User'
		);
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

	// Get document count for a case
	function getDocumentCount(caseId: any) {
		return caseDocuments.filter((doc: { case_id: any }) => doc.case_id === caseId).length;
	}

	// Format date for display
	function formatDate(dateString: string | number | Date) {
		if (!dateString) return 'N/A';
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	// Handle form submission results
	function handleFormResult(result: any) {
		if (result?.type === 'success') {
			const formResult = result.data;
			updateSuccess = formResult.success;
			updateMessage = formResult.message;
		} else {
			updateSuccess = false;
			updateMessage = 'An error occurred. Please try again.';
		}
		resetUpdateStatus();
	}

	// Reset the update status message after a delay
	function resetUpdateStatus() {
		if (messageTimer) clearTimeout(messageTimer);

		messageTimer = setTimeout(() => {
			updateSuccess = null;
			updateMessage = '';
		}, 5000); // Message disappears after 5 seconds
	}

	// Sign out function
	async function signOut() {
		try {
			// Use the dedicated logout route which handles proper session cleanup
			await goto('/auth/logout');
		} catch (error) {
			console.error('Error signing out:', error);
		}
	}

	// Navigate to case view
	function viewCase(caseId: any) {
		goto(`/cases/${caseId}`);
	}

	// Get case name or placeholder
	function getCaseName(caseItem: any) {
		if (caseItem.case_name) {
			return caseItem.case_name;
		}

		if (caseItem.id) {
			return `Case #${caseItem.id.substring(0, 6)}`;
		}
		return 'Untitled Case';
	}

	// Load user's cases
	async function loadUserCases() {
		if (!session?.user) return;
		cases = await caseManager.getUserCases(session.user.id);
	}

	// Create new case
	async function createCase() {
		if (!session?.user || !newCaseName.trim()) return;
		
		isCreatingCase = true;
		try {
			const newCase = await caseManager.createCase(session.user.id, newCaseName.trim());
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

	// Handle keydown for modal escape
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && showNewCaseModal) {
			closeNewCaseModal();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="flex min-h-[80vh] items-center justify-center p-4">
	<!-- Main container with responsive layout -->
	<div
		class="w-full max-w-7xl overflow-hidden rounded-xl bg-white shadow-xl"
		in:fly={{ y: 30, duration: 800, easing: quintOut, delay: isInitialLoad ? 0 : 0 }}
	>
		<!-- Two column layout for user profile (left) and cases (right) -->
		<div class="flex flex-col md:flex-row">
			<!-- User profile panel (left side) -->
			<div
				class="w-full border-r border-gray-200 bg-gray-50 p-6 md:w-2/5"
				in:fly={{ x: -20, duration: 700, easing: quintOut, delay: isInitialLoad ? 100 : 0 }}
			>
				<h2
					class="mb-4 text-xl font-semibold text-gray-800"
					in:slide={{ duration: 500, delay: isInitialLoad ? 200 : 0 }}
				>
					Account Information
				</h2>

				<!-- User profile card -->
				<div class="mb-6 overflow-hidden rounded-lg bg-white p-4 shadow-sm">
					<div class="flex items-center">
						<div
							class="flex h-16 w-16 items-center justify-center rounded-full bg-indigo-100 text-indigo-600"
							in:scale={{
								start: 0.9,
								duration: 600,
								delay: isInitialLoad ? 300 : 0,
								easing: elasticOut
							}}
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-8 w-8"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="2"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
								/>
							</svg>
						</div>
						<div class="ml-4" in:slide={{ duration: 500, delay: isInitialLoad ? 400 : 0 }}>
							<h3 class="text-lg font-medium text-gray-800">
								{session?.user?.email || 'User'}
							</h3>
							<p class="text-sm text-gray-500">
								{isLoadingUserDetails ? 'Loading details...' : getUserDisplayName()}
							</p>
						</div>
					</div>
				</div>

				<!-- Account stats -->
				<div class="mb-6 grid grid-cols-2 gap-4">
					<div
						class="rounded-lg bg-white p-4 shadow-sm"
						in:fly={{ y: 15, duration: 500, delay: isInitialLoad ? 500 : 0, easing: cubicOut }}
					>
						<h4 class="text-sm font-medium text-gray-500">Total Cases</h4>
						<p class="mt-1 text-2xl font-semibold text-indigo-600">{cases.length}</p>
					</div>
					<div
						class="rounded-lg bg-white p-4 shadow-sm"
						in:fly={{ y: 15, duration: 500, delay: isInitialLoad ? 600 : 0, easing: cubicOut }}
					>
						<h4 class="text-sm font-medium text-gray-500">Total Documents</h4>
						<p class="mt-1 text-2xl font-semibold text-indigo-600">{caseDocuments.length}</p>
					</div>
				</div>

				<!-- Account actions -->
				<div class="space-y-3" in:slide={{ duration: 500, delay: isInitialLoad ? 700 : 0 }}>
					<h3 class="text-md font-semibold text-gray-700">Account Actions</h3>
					<button
						class="flex w-full items-center justify-between rounded-lg border border-gray-300 bg-white p-3 text-left text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50"
						in:scale={{
							start: 0.95,
							duration: 400,
							delay: isInitialLoad ? 800 : 0,
							easing: backOut
						}}
					>
						<span class="flex items-center">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="mr-2 h-5 w-5 text-gray-400"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="2"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
								/>
							</svg>
							Edit Profile
						</span>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5 text-gray-400"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
						</svg>
					</button>
					<button
						class="flex w-full items-center justify-between rounded-lg border border-gray-300 bg-white p-3 text-left text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50"
						in:scale={{
							start: 0.95,
							duration: 400,
							delay: isInitialLoad ? 900 : 0,
							easing: backOut
						}}
					>
						<span class="flex items-center">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="mr-2 h-5 w-5 text-gray-400"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="2"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
								/>
							</svg>
							Change Password
						</span>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5 text-gray-400"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
						</svg>
					</button>
					<button
						onclick={signOut}
						class="flex w-full items-center justify-between rounded-lg border border-red-200 bg-white p-3 text-left text-sm font-medium text-red-600 shadow-sm hover:bg-red-50"
						in:scale={{
							start: 0.95,
							duration: 400,
							delay: isInitialLoad ? 1000 : 0,
							easing: backOut
						}}
					>
						<span class="flex items-center">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="mr-2 h-5 w-5"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="2"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
								/>
							</svg>
							Sign Out
						</span>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
						</svg>
					</button>
				</div>
			</div>

			<!-- Cases management (right side) -->
			<div
				class="w-full p-6 md:w-3/5"
				in:fly={{ x: 20, duration: 700, delay: isInitialLoad ? 200 : 0, easing: quintOut }}
			>
				<h1
					class="mb-6 text-center text-2xl font-bold text-indigo-700"
					in:slide={{ duration: 600, delay: isInitialLoad ? 300 : 0 }}
				>
					Your Cases
				</h1>

				<!-- Success/Error message -->
				{#if updateSuccess !== null}
					<div
						class="mb-4 rounded-md p-3 {updateSuccess
							? 'bg-green-50 text-green-800'
							: 'bg-red-50 text-red-800'}"
						in:fly={{ y: -10, duration: 300, easing: cubicOut }}
						out:fade
					>
						<div class="flex">
							<div class="flex-shrink-0">
								{#if updateSuccess}
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
									{updateMessage}
								</p>
							</div>
						</div>
					</div>
				{/if}

				<!-- Create new case button -->
				<div
					class="mb-6 rounded-lg border border-gray-200 bg-white p-4 shadow-sm"
					in:fly={{ y: 15, duration: 600, delay: isInitialLoad ? 400 : 0, easing: cubicOut }}
				>
					<div class="flex items-center justify-between">
						<h3 class="text-lg font-medium text-gray-800">Case Management</h3>
						<button
							onclick={openNewCaseModal}
							class="inline-flex items-center justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none"
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

				<!-- Cases list -->
				{#if cases.length === 0}
					<div
						class="mb-4 rounded-lg border border-dashed border-gray-300 bg-gray-50 p-6 text-center"
						in:fly={{ y: 20, duration: 600, delay: isInitialLoad ? 500 : 0, easing: cubicOut }}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mx-auto h-12 w-12 text-gray-400"
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
						<h3 class="mt-2 text-sm font-medium text-gray-900">No cases yet</h3>
						<p class="mt-1 text-sm text-gray-500">Create your first case to get started</p>
					</div>
				{:else}
					<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
						{#each cases as caseItem, i}
							<div
								class="bg-white border border-gray-200 rounded-lg p-4 hover:border-indigo-300 transition-colors cursor-pointer"
								in:fly={{ y: 20, duration: 400, delay: i * 100 }}
							>
								<a href="/cases/{caseItem.id}" class="block">
									<div class="flex items-center mb-3">
										<div class="h-8 w-8 rounded-full bg-indigo-100 flex items-center justify-center">
											<svg class="h-5 w-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
											</svg>
										</div>
										<div class="ml-3 flex-1">
											<h3 class="text-sm font-medium text-gray-900 truncate">
												{caseItem.case_name}
											</h3>
											<p class="text-xs text-gray-500">
												{new Date(caseItem.created_at).toLocaleDateString()}
											</p>
										</div>
									</div>
								</a>
								<div class="flex justify-between items-center pt-2 border-t border-gray-100">
									<span class="text-xs text-gray-400">
										Updated {new Date(caseItem.updated_at).toLocaleDateString()}
									</span>
									<button
										class="text-red-600 hover:text-red-800 text-xs"
										onclick={(e) => { e.preventDefault(); deleteCase(caseItem.id); }}
									>
										Delete
									</button>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>
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
				<h3 class="text-xl font-medium text-gray-900">Create New Case</h3>
				<button
					onclick={closeNewCaseModal}
					class="rounded-full p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-500"
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
				<label for="case-name" class="block text-sm font-medium text-gray-700">Case Name</label>
				<input
					type="text"
					id="case-name"
					bind:value={newCaseName}
					class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
					placeholder="Enter a name for your case"
					onkeydown={(e) => e.key === 'Enter' && createCase()}
					required
				/>
				<p class="mt-1 text-xs text-gray-500">
					Give your case a descriptive name to help you identify it later.
				</p>
			</div>

			<div class="mt-5 flex justify-end space-x-3">
				<button
					type="button"
					onclick={closeNewCaseModal}
					class="inline-flex justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none"
				>
					Cancel
				</button>
				<button
					type="button"
					class="inline-flex justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
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


