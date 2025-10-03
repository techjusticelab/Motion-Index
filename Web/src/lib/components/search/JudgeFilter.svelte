<!-- JudgeFilter.svelte -->
<script lang="ts">
	import { onMount, createEventDispatcher } from 'svelte';
	import { filterOptions } from '$lib/utils';

	export let selectedJudges: string[] = [];
	export let allJudgeOptions: string[] = [];

	let judgeSearchInput = '';
	let filteredJudgeOptions: string[] = [];
	let showJudgeDropdown = false;
	let judgeDropdownRef: HTMLDivElement;

	const dispatch = createEventDispatcher<{
		add: string;
		remove: string;
	}>();

	// Filter judge options
	function filterJudgeOptions() {
		filteredJudgeOptions = filterOptions(allJudgeOptions, judgeSearchInput);
	}

	// Add a judge to the selected judges
	function addJudge(judge: string) {
		dispatch('add', judge);
		judgeSearchInput = '';
		filterJudgeOptions();
	}

	// Remove a judge from the selected judges
	function removeJudge(judge: string) {
		dispatch('remove', judge);
	}

	// Handle clicks outside the dropdown
	function handleClickOutside(event: MouseEvent) {
		if (judgeDropdownRef && !judgeDropdownRef.contains(event.target as Node)) {
			showJudgeDropdown = false;
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
	<label for="judge-search" class="mb-1 block text-xs font-medium text-neutral-700">Judges</label>

	<!-- Selected judges tags -->
	{#if selectedJudges.length > 0}
		<div class="mb-2 flex flex-wrap gap-2">
			{#each selectedJudges as judge}
				<div class="flex items-center rounded-lg bg-primary-50 px-2 py-1 text-xs text-primary-700">
					<span class="mr-1 max-w-[200px] truncate">{judge}</span>
					<button
						type="button"
						onclick={() => removeJudge(judge)}
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

	<!-- Judge search input -->
	<div class="relative" bind:this={judgeDropdownRef}>
		<input
			type="text"
			id="judge-search"
			bind:value={judgeSearchInput}
			oninput={filterJudgeOptions}
			onfocus={() => {
				showJudgeDropdown = true;
				filterJudgeOptions();
			}}
			placeholder="Search judges..."
			class="w-full rounded-lg border border-neutral-200 px-3 py-2 text-sm shadow-sm focus:border-primary-900 focus:ring-primary-900"
		/>

		<!-- Dropdown for judge options -->
		{#if showJudgeDropdown && filteredJudgeOptions.length > 0}
			<div
				class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-sm shadow-lg ring-1 ring-black ring-opacity-5"
			>
				{#each filteredJudgeOptions as judge}
					<button
						type="button"
						class="block w-full px-4 py-2 text-left hover:bg-neutral-100 {selectedJudges.includes(
							judge
						)
							? 'bg-primary-50'
							: ''}"
						onclick={() => {
							addJudge(judge);
							showJudgeDropdown = false;
						}}
					>
						{judge}
					</button>
				{/each}
			</div>
		{/if}
	</div>
	<p class="mt-1 text-xs text-neutral-500">Search and click to add multiple judges</p>
</div>
