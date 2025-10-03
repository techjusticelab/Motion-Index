<!-- TagsFilter.svelte -->
<script lang="ts">
	import { onMount, createEventDispatcher } from 'svelte';
	import { filterOptions } from '$lib/utils';

	export let selectedTags: string[] = [];
	export let allTagsOptions: string[] = [];
	export let matchAll: boolean = false; // Default to OR behavior (any tag)

	let TagsSearchInput = '';
	let filteredTagsOptions: string[] = [];
	let showTagsDropdown = false;
	let TagsDropdownRef: HTMLDivElement;

	const dispatch = createEventDispatcher<{
		add: string;
		remove: string;
		matchAllToggle: boolean;
	}>();

	// Filter Tags options
	function filterTagsOptions() {
		filteredTagsOptions = filterOptions(allTagsOptions, TagsSearchInput);
	}

	// Add a Tags to the selected Tagss
	function addTags(Tags: string) {
		dispatch('add', Tags);
		TagsSearchInput = '';
		filterTagsOptions();
	}

	// Remove a Tags from the selected Tagss
	function removeTags(Tags: string) {
		dispatch('remove', Tags);
	}

	// Toggle between match all (AND) and match any (OR) behavior
	function toggleMatchAll() {
		matchAll = !matchAll;
		dispatch('matchAllToggle', matchAll);
	}

	// Handle clicks outside the dropdown
	function handleClickOutside(event: MouseEvent) {
		if (TagsDropdownRef && !TagsDropdownRef.contains(event.target as Node)) {
			showTagsDropdown = false;
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
	<div class="flex justify-between items-center mb-1">
		<label for="Tags-search" class="block text-xs font-medium text-neutral-700">Tags</label>
		
		<!-- Toggle switch for match all/any -->
		<div class="flex items-center">
			<span class="text-xs text-neutral-500 mr-2">{matchAll ? 'Match ALL tags' : 'Match ANY tag'}</span>
			<button 
				type="button" 
				onclick={toggleMatchAll}
				class="relative inline-flex h-5 w-10 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent bg-neutral-200 transition-colors duration-200 ease-in-out focus:outline-none {matchAll ? 'bg-primary-600' : ''}"
			>
				<span 
					class="pointer-events-none relative inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out {matchAll ? 'translate-x-5' : 'translate-x-0'}"
				></span>
			</button>
		</div>
	</div>

	<!-- Selected Tagss tags -->
	{#if selectedTags.length > 0}
		<div class="mb-2 flex flex-wrap gap-2">
			{#each selectedTags as Tags}
				<div class="flex items-center rounded-lg bg-primary-50 px-2 py-1 text-xs text-primary-700">
					<span class="mr-1 max-w-[200px] truncate">{Tags}</span>
					<button
						type="button"
						onclick={() => removeTags(Tags)}
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

	<!-- Tags search input -->
	<div class="relative" bind:this={TagsDropdownRef}>
		<input
			type="text"
			id="Tags-search"
			bind:value={TagsSearchInput}
			oninput={filterTagsOptions}
			onfocus={() => {
				showTagsDropdown = true;
				filterTagsOptions();
			}}
			placeholder="Search Tags..."
			class="w-full rounded-lg border border-neutral-200 px-3 py-2 text-sm shadow-sm focus:border-primary-900 focus:ring-primary-900"
		/>

		<!-- Dropdown for Tags options -->
		{#if showTagsDropdown && filteredTagsOptions.length > 0}
			<div
				class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-sm shadow-lg ring-1 ring-black ring-opacity-5"
			>
				{#each filteredTagsOptions as Tags}
					<button
						type="button"
						class="block w-full px-4 py-2 text-left hover:bg-neutral-100 {selectedTags.includes(Tags)
							? 'bg-primary-50'
							: ''}"
						onclick={() => {
							addTags(Tags);
							showTagsDropdown = false;
						}}
					>
						{Tags}
					</button>
				{/each}
			</div>
		{/if}
	</div>
	<p class="mt-1 text-xs text-neutral-500">Search and click to add multiple Tags</p>
</div>
