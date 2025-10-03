<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { fade } from 'svelte/transition';

	const dispatch = createEventDispatcher();

	interface Props {
		type?: 'success' | 'error' | 'warning' | 'info';
		title?: string;
		message?: string;
		dismissible?: boolean;
		class?: string;
		children?: any;
	}

	let {
		type = 'info',
		title = '',
		message = '',
		dismissible = false,
		class: className = '',
		children
	}: Props = $props();

	// Type-specific classes
	const typeClasses = {
		success: {
			container: 'bg-green-50 border-green-200 text-green-800',
			icon: 'text-green-400',
			iconPath: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'
		},
		error: {
			container: 'bg-red-50 border-red-200 text-red-800',
			icon: 'text-red-400',
			iconPath: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z'
		},
		warning: {
			container: 'bg-amber-50 border-amber-200 text-amber-800',
			icon: 'text-amber-400',
			iconPath: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z'
		},
		info: {
			container: 'bg-blue-50 border-blue-200 text-blue-800',
			icon: 'text-blue-400',
			iconPath: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
		}
	};

	const currentType = typeClasses[type];

	function dismiss() {
		dispatch('dismiss');
	}
</script>

<div
	class="rounded-lg border p-4 {currentType.container} {className}"
	transition:fade={{ duration: 200 }}
	role="alert"
>
	<div class="flex items-start">
		<!-- Icon -->
		<div class="flex-shrink-0">
			<svg class="h-5 w-5 {currentType.icon}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={currentType.iconPath} />
			</svg>
		</div>

		<!-- Content -->
		<div class="ml-3 flex-1">
			{#if title}
				<h3 class="text-sm font-medium">
					{title}
				</h3>
			{/if}

			{#if message || children}
				<div class="text-sm {title ? 'mt-1' : ''}">
					{#if message}
						<p>{message}</p>
					{/if}
					{@render children?.()}
				</div>
			{/if}
		</div>

		<!-- Dismiss button -->
		{#if dismissible}
			<div class="ml-auto pl-3">
				<button
					type="button"
					class="inline-flex rounded-md p-1.5 hover:bg-black hover:bg-opacity-10 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-transparent focus:ring-current"
					onclick={dismiss}
					aria-label="Dismiss alert"
				>
					<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			</div>
		{/if}
	</div>
</div>