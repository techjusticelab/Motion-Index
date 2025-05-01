<script lang="ts">
	import { onMount } from 'svelte';
	import { fade, fly, slide, scale } from 'svelte/transition';
	import { cubicOut, quintOut, elasticOut, backOut } from 'svelte/easing';
	import { page } from '$app/stores';
	import { supabaseClient } from '$lib/supabase';

	// User data
	let user = $page.data.session?.user;
	let isLoadingUserDetails = false;
	let userDetails = null;

	// Case management
	let cases = $page.data.cases || [];
	let caseDocuments = $page.data.caseDocuments || [];
	let isCreatingCase = false;
	let isUpdatingCase = false;
	let newCaseName = '';
	let selectedCase = null;
	let updateSuccess = null;
	let updateMessage = '';

	// Form states for case editing
	let editingCaseId = null;
	let editCaseName = '';

	// Timer for success/error message
	let messageTimer: ReturnType<typeof setTimeout> | null = null;

	// Get user details on mount
	onMount(async () => {
		if (user) {
			await loadUserDetails();
		}
	});

	async function loadUserDetails() {
		isLoadingUserDetails = true;

		try {
			const { data, error } = await supabaseClient
				.from('profiles')
				.select('*')
				.eq('id', user.id)
				.single();

			if (error) throw error;
			userDetails = data;
		} catch (error) {
			console.error('Error loading user details:', error);
		} finally {
			isLoadingUserDetails = false;
		}
	}

	// Create a new case
	async function createCase() {
		if (!newCaseName.trim()) {
			updateSuccess = false;
			updateMessage = 'Please enter a case name';
			resetUpdateStatus();
			return;
		}

		isCreatingCase = true;

		try {
			const { data, error } = await supabaseClient
				.from('cases')
				.insert({
					user_id: user.id,
					case_docs: [],
					created_at: new Date().toISOString(),
					updated_at: new Date().toISOString()
				})
				.select();

			if (error) throw error;

			// Add to cases list
			cases = [...cases, data[0]];
			newCaseName = '';

			updateSuccess = true;
			updateMessage = 'Case created successfully!';
		} catch (error) {
			console.error('Error creating case:', error);
			updateSuccess = false;
			updateMessage = 'Failed to create case. Please try again.';
		} finally {
			isCreatingCase = false;
			resetUpdateStatus();
		}
	}

	// Start editing a case
	function startEditing(caseItem) {
		editingCaseId = caseItem.id;
		editCaseName = caseItem.name || '';
		selectedCase = caseItem;
	}

	// Cancel editing
	function cancelEditing() {
		editingCaseId = null;
		editCaseName = '';
	}

	// Update case details
	async function updateCase() {
		if (!editingCaseId) return;

		isUpdatingCase = true;

		try {
			const { data, error } = await supabaseClient
				.from('cases')
				.update({
					name: editCaseName,
					updated_at: new Date().toISOString()
				})
				.eq('id', editingCaseId)
				.select();

			if (error) throw error;

			// Update case in the list
			cases = cases.map((c) => (c.id === editingCaseId ? data[0] : c));

			updateSuccess = true;
			updateMessage = 'Case updated successfully!';
			editingCaseId = null;
		} catch (error) {
			console.error('Error updating case:', error);
			updateSuccess = false;
			updateMessage = 'Failed to update case. Please try again.';
		} finally {
			isUpdatingCase = false;
			resetUpdateStatus();
		}
	}

	// Delete a case
	async function deleteCase(caseId) {
		if (!confirm('Are you sure you want to delete this case? This action cannot be undone.')) {
			return;
		}

		try {
			// First delete all case documents associated with this case
			const { error: docDeleteError } = await supabaseClient
				.from('case_documents')
				.delete()
				.eq('case_id', caseId);

			if (docDeleteError) throw docDeleteError;

			// Then delete the case itself
			const { error } = await supabaseClient.from('cases').delete().eq('id', caseId);

			if (error) throw error;

			// Update the lists
			cases = cases.filter((c) => c.id !== caseId);
			caseDocuments = caseDocuments.filter((d) => d.case_id !== caseId);

			if (selectedCase && selectedCase.id === caseId) {
				selectedCase = null;
			}

			updateSuccess = true;
			updateMessage = 'Case deleted successfully!';
		} catch (error) {
			console.error('Error deleting case:', error);
			updateSuccess = false;
			updateMessage = 'Failed to delete case. Please try again.';
		} finally {
			resetUpdateStatus();
		}
	}

	// Get document count for a case
	function getDocumentCount(caseId) {
		return caseDocuments.filter((doc) => doc.case_id === caseId).length;
	}

	// Format date for display
	function formatDate(dateString) {
		if (!dateString) return 'N/A';
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
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
			await supabaseClient.auth.signOut();
			window.location.href = '/login';
		} catch (error) {
			console.error('Error signing out:', error);
		}
	}
</script>

<div class="flex min-h-[80vh] items-center justify-center p-4">
	<!-- Main container with responsive layout -->
	<div
		class="w-full max-w-7xl overflow-hidden rounded-xl bg-white shadow-xl"
		in:fly={{ y: 30, duration: 800, easing: quintOut }}
	>
		<!-- Two column layout for user profile (left) and cases (right) -->
		<div class="flex flex-col md:flex-row">
			<!-- User profile panel (left side) -->
			<div
				class="w-full border-r border-gray-200 bg-gray-50 p-6 md:w-2/5"
				in:fly={{ x: -20, duration: 700, easing: quintOut }}
			>
				<h2
					class="mb-4 text-xl font-semibold text-gray-800"
					in:slide={{ duration: 500, delay: 100 }}
				>
					Account Information
				</h2>

				<!-- User profile card -->
				<div class="mb-6 overflow-hidden rounded-lg bg-white p-4 shadow-sm">
					<div class="flex items-center">
						<div
							class="flex h-16 w-16 items-center justify-center rounded-full bg-indigo-100 text-indigo-600"
							in:scale={{ start: 0.9, duration: 600, delay: 200, easing: elasticOut }}
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
						<div class="ml-4" in:slide={{ duration: 500, delay: 300 }}>
							<h3 class="text-lg font-medium text-gray-800">
								{user?.email || 'User'}
							</h3>
							<p class="text-sm text-gray-500">
								{isLoadingUserDetails
									? 'Loading details...'
									: userDetails?.full_name || 'No name set'}
							</p>
						</div>
					</div>
				</div>

				<!-- Account stats -->
				<div class="mb-6 grid grid-cols-2 gap-4">
					<div
						class="rounded-lg bg-white p-4 shadow-sm"
						in:fly={{ y: 15, duration: 500, delay: 400, easing: cubicOut }}
					>
						<h4 class="text-sm font-medium text-gray-500">Total Cases</h4>
						<p class="mt-1 text-2xl font-semibold text-indigo-600">{cases.length}</p>
					</div>
					<div
						class="rounded-lg bg-white p-4 shadow-sm"
						in:fly={{ y: 15, duration: 500, delay: 500, easing: cubicOut }}
					>
						<h4 class="text-sm font-medium text-gray-500">Total Documents</h4>
						<p class="mt-1 text-2xl font-semibold text-indigo-600">{caseDocuments.length}</p>
					</div>
				</div>

				<!-- Account actions -->
				<div class="space-y-3" in:slide={{ duration: 500, delay: 600 }}>
					<h3 class="text-md font-semibold text-gray-700">Account Actions</h3>
					<button
						class="flex w-full items-center justify-between rounded-lg border border-gray-300 bg-white p-3 text-left text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50"
						in:scale={{ start: 0.95, duration: 400, delay: 700, easing: backOut }}
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
						in:scale={{ start: 0.95, duration: 400, delay: 800, easing: backOut }}
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
						on:click={signOut}
						class="flex w-full items-center justify-between rounded-lg border border-red-200 bg-white p-3 text-left text-sm font-medium text-red-600 shadow-sm hover:bg-red-50"
						in:scale={{ start: 0.95, duration: 400, delay: 900, easing: backOut }}
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
				in:fly={{ x: 20, duration: 700, delay: 100, easing: quintOut }}
			>
				<h1
					class="mb-6 text-center text-2xl font-bold text-indigo-700"
					in:slide={{ duration: 600, delay: 200 }}
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

				<!-- Create new case form -->
				<div
					class="mb-6 rounded-lg border border-gray-200 bg-white p-4 shadow-sm"
					in:fly={{ y: 15, duration: 600, delay: 300, easing: cubicOut }}
				>
					<h3 class="mb-3 text-lg font-medium text-gray-800">Create New Case</h3>
					<form class="flex gap-2" on:submit|preventDefault={createCase}>
						<input
							type="text"
							bind:value={newCaseName}
							placeholder="Enter case name"
							class="flex-1 rounded-md border border-gray-300 p-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
						/>
						<button
							type="submit"
							class="inline-flex justify-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none disabled:cursor-not-allowed disabled:opacity-60"
							disabled={isCreatingCase}
						>
							{#if isCreatingCase}
								<div class="flex items-center">
									<div
										class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white"
									></div>
									<span>Creating...</span>
								</div>
							{:else}
								Create Case
							{/if}
						</button>
					</form>
				</div>

				<!-- Cases list -->
				{#if cases.length === 0}
					<div
						class="mb-4 rounded-lg border border-dashed border-gray-300 bg-gray-50 p-6 text-center"
						in:fly={{ y: 20, duration: 600, delay: 400, easing: cubicOut }}
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
						<h3 class="mt-2 text-sm font-medium text-gray-900">No cases</h3>
						<p class="mt-1 text-sm text-gray-500">Get started by creating a new case.</p>
					</div>
				{:else}
					<div class="space-y-4">
						{#each cases as caseItem, i}
							<div
								class="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm transition-shadow hover:shadow-md"
								in:fly={{ y: 20, duration: 600, delay: 400 + i * 100, easing: cubicOut }}
							>
								<div class="p-4">
									{#if editingCaseId === caseItem.id}
										<!-- Edit mode -->
										<form class="flex gap-2" on:submit|preventDefault={updateCase}>
											<input
												type="text"
												bind:value={editCaseName}
												class="flex-1 rounded-md border border-gray-300 p-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
											/>
											<div class="flex gap-1">
												<button
													type="submit"
													class="inline-flex justify-center rounded-md bg-indigo-600 px-3 py-1 text-sm font-medium text-white hover:bg-indigo-700"
													disabled={isUpdatingCase}
												>
													{#if isUpdatingCase}
														<div
															class="h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white"
														></div>
													{:else}
														Save
													{/if}
												</button>
												<button
													type="button"
													class="inline-flex justify-center rounded-md bg-gray-200 px-3 py-1 text-sm font-medium text-gray-700 hover:bg-gray-300"
													on:click={cancelEditing}
												>
													Cancel
												</button>
											</div>
										</form>
									{:else}
										<!-- View mode -->
										<div class="flex justify-between">
											<div>
												<h3 class="text-lg font-medium text-gray-800">
													{caseItem.name || `Case #${caseItem.id.substring(0, 6)}`}
												</h3>
												<div class="mt-1 flex items-center space-x-4 text-sm text-gray-500">
													<span>
														<svg
															xmlns="http://www.w3.org/2000/svg"
															class="mr-1 inline h-4 w-4"
															fill="none"
															viewBox="0 0 24 24"
															stroke="currentColor"
														>
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																stroke-width="2"
																d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
															/>
														</svg>
														{formatDate(caseItem.created_at)}
													</span>
													<span>
														<svg
															xmlns="http://www.w3.org/2000/svg"
															class="mr-1 inline h-4 w-4"
															fill="none"
															viewBox="0 0 24 24"
															stroke="currentColor"
														>
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																stroke-width="2"
																d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
															/>
														</svg>
														{getDocumentCount(caseItem.id)} document{getDocumentCount(
															caseItem.id
														) !== 1
															? 's'
															: ''}
													</span>
												</div>
											</div>
											<div class="flex space-x-2">
												<button
													class="inline-flex h-8 w-8 items-center justify-center rounded-full bg-gray-100 text-gray-600 hover:bg-gray-200"
													on:click={() => startEditing(caseItem)}
												>
													<svg
														xmlns="http://www.w3.org/2000/svg"
														class="h-4 w-4"
														fill="none"
														viewBox="0 0 24 24"
														stroke="currentColor"
													>
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
														/>
													</svg>
												</button>
												<button
													class="inline-flex h-8 w-8 items-center justify-center rounded-full bg-red-100 text-red-600 hover:bg-red-200"
													on:click={() => deleteCase(caseItem.id)}
												>
													<svg
														xmlns="http://www.w3.org/2000/svg"
														class="h-4 w-4"
														fill="none"
														viewBox="0 0 24 24"
														stroke="currentColor"
													>
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
														/>
													</svg>
												</button>
											</div>
										</div>
									{/if}
								</div>
								<div class="border-t border-gray-100 bg-gray-50 px-4 py-3">
									<button
										class="flex w-full items-center justify-center text-sm font-medium text-indigo-600 hover:text-indigo-800"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											class="mr-1 h-4 w-4"
											fill="none"
											viewBox="0 0 24 24"
											stroke="currentColor"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
											/>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
											/>
										</svg>
										View Documents
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
