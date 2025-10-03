<!-- CourtFilter.svelte -->
<script lang="ts">
	import { onMount, createEventDispatcher } from 'svelte';
    import { filterOptions } from '$lib/utils';

	export let selectedCourts: string[] = [];
	export let allCourtOptions: string[] = [];

	let courtSearchInput = '';
	let filteredCourtOptions: string[] = [];
	let showCourtDropdown = false;
	let courtDropdownRef: HTMLDivElement;

	const dispatch = createEventDispatcher<{
		add: string;
		remove: string;
	}>();

	// Filter court options
	function filterCourtOptions() {
		filteredCourtOptions = filterOptions(allCourtOptions, courtSearchInput);
	}

	// Add a court to the selected courts
	function addCourt(court: string) {
		dispatch('add', court);
		courtSearchInput = '';
		filterCourtOptions();
	}

	// Remove a court from the selected courts
	function removeCourt(court: string) {
		dispatch('remove', court);
	}

	// Handle clicks outside the dropdown
	function handleClickOutside(event: MouseEvent) {
		if (courtDropdownRef && !courtDropdownRef.contains(event.target as Node)) {
			showCourtDropdown = false;
		}
	}

	onMount(() => {
		document.addEventListener('click', handleClickOutside);
		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
	});
</script>

<div>
	<label for="court-search" class="mb-1 block text-xs font-medium text-neutral-700">Courts</label>

	<!-- Selected courts tags -->
	{#if selectedCourts.length > 0}
		<div class="mb-2 flex flex-wrap gap-2">
			{#each selectedCourts as court}
				<div class="flex items-center rounded-lg bg-primary-50 px-2 py-1 text-xs text-primary-700">
					<span class="mr-1 max-w-[200px] truncate">{court}</span>
					<button
						type="button"
						onclick={() => removeCourt(court)}
						class="ml-1 text-primary-500 hover:text-primary-700"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-3 w-3"
							viewBox="0 0 20 20"
							fill="currentColor"
						>
							<path
								fill-rule="evenodd"
								d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
								clip-rule="evenodd"
							/>
						</svg>
					</button>
				</div>
			{/each}
		</div>
	{/if}

	<!-- Court search input -->
	<div class="relative" bind:this={courtDropdownRef}>
		<input
			type="text"
			id="court-search"
			bind:value={courtSearchInput}
			oninput={filterCourtOptions}
			onfocus={() => {
				showCourtDropdown = true;
				filterCourtOptions();
			}}
			placeholder="Search courts..."
			class="w-full rounded-lg border border-neutral-200 px-3 py-2 text-sm shadow-sm focus:border-primary-900 focus:ring-primary-900"
		/>

		<!-- Dropdown for court options -->
		{#if showCourtDropdown && filteredCourtOptions.length > 0}
			<div
				class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-sm shadow-lg ring-1 ring-black ring-opacity-5"
			>
				{#each filteredCourtOptions as court}
					<button
						type="button"
						class="block w-full px-4 py-2 text-left hover:bg-neutral-100 {selectedCourts.includes(
							court
						)
							? 'bg-primary-50'
							: ''}"
						onclick={() => {
							addCourt(court);
							showCourtDropdown = false;
						}}
					>
						{court}
					</button>
				{/each}
			</div>
		{/if}
	</div>
	<p class="mt-1 text-xs text-neutral-500">Search and click to add multiple courts</p>
</div>
