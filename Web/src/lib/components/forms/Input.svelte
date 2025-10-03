<script lang="ts">
	import type { InputProps } from '$lib/types';

	interface Props extends InputProps {
		id?: string;
		name?: string;
		label?: string;
		class?: string;
		value?: string;
	}

	let {
		id = '',
		name = '',
		label = '',
		type = 'text',
		value = $bindable(''),
		placeholder = '',
		disabled = false,
		required = false,
		error = '',
		class: className = ''
	}: Props = $props();

	// Generate unique ID if not provided
	const inputId = id || `input-${Math.random().toString(36).substr(2, 9)}`;

	// Base input classes
	const baseClasses = 'block w-full rounded-md border px-3 py-2 text-sm placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-offset-0 disabled:cursor-not-allowed disabled:opacity-50';

	// State-specific classes
	const stateClasses = error
		? 'border-red-300 text-red-900 focus:border-red-500 focus:ring-red-500'
		: 'border-neutral-300 text-neutral-900 focus:border-primary-500 focus:ring-primary-500';

	// Combine all classes
	const inputClasses = [baseClasses, stateClasses, className].join(' ');
</script>

<div class="space-y-1">
	{#if label}
		<label
			for={inputId}
			class="block text-sm font-medium text-neutral-700"
			class:text-red-700={error}
		>
			{label}
			{#if required}
				<span class="text-red-500" aria-label="required">*</span>
			{/if}
		</label>
	{/if}

	<input
		{id}
		{name}
		{type}
		{placeholder}
		{disabled}
		{required}
		class={inputClasses}
		bind:value
		aria-invalid={error ? 'true' : 'false'}
		aria-describedby={error ? `${inputId}-error` : undefined}
	/>

	{#if error}
		<p id="{inputId}-error" class="text-sm text-red-600" role="alert">
			{error}
		</p>
	{/if}
</div>