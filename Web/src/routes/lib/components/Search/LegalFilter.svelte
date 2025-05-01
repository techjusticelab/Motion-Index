<!-- TagsFilter.svelte -->
<script lang="ts">
	import { onMount, createEventDispatcher } from 'svelte';
	import { filterOptions } from '../../../utils/utils';

	export let selectedTags: string[] = [];
	export let allTagsOptions: string[] = [];

	let TagsSearchInput = '';
	let filteredTagsOptions: string[] = [];
	let showTagsDropdown = false;
	let TagsDropdownRef: HTMLDivElement;

	const dispatch = createEventDispatcher<{
		add: string;
		remove: string;
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
	<label for="Tags-search" class="mb-1 block text-xs font-medium text-gray-700">Tags</label>

	<!-- Selected Tagss tags -->
	{#if selectedTags.length > 0}
		<div class="mb-2 flex flex-wrap gap-2">
			{#each selectedTags as Tags}
				<div class="flex items-center rounded-lg bg-blue-50 px-2 py-1 text-xs text-blue-700">
					<span class="mr-1 max-w-[200px] truncate">{Tags}</span>
					<button
						type="button"
						on:click={() => removeTags(Tags)}
						class="ml-1 text-blue-500 hover:text-blue-700"
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
			on:input={filterTagsOptions}
			on:focus={() => {
				showTagsDropdown = true;
				filterTagsOptions();
			}}
			placeholder="Search Tagss..."
			class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-blue-500"
		/>

		<!-- Dropdown for Tags options -->
		{#if showTagsDropdown && filteredTagsOptions.length > 0}
			<div
				class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-sm shadow-lg ring-1 ring-black ring-opacity-5"
			>
				{#each filteredTagsOptions as Tags}
					<button
						type="button"
						class="block w-full px-4 py-2 text-left hover:bg-gray-100 {selectedTags.includes(Tags)
							? 'bg-blue-50'
							: ''}"
						on:click={() => {
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
	<p class="mt-1 text-xs text-gray-500">Search and click to add multiple Tags</p>
</div>
