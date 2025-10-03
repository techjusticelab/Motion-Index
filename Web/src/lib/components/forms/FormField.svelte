<script lang="ts">
	interface Props {
		label?: string;
		error?: string;
		hint?: string;
		required?: boolean;
		class?: string;
		children?: any;
	}

	let {
		label = '',
		error = '',
		hint = '',
		required = false,
		class: className = '',
		children
	}: Props = $props();

	// Generate unique ID for accessibility
	const fieldId = `field-${Math.random().toString(36).substr(2, 9)}`;
</script>

<div class="space-y-1 {className}">
	{#if label}
		<label
			for={fieldId}
			class="block text-sm font-medium text-neutral-700"
			class:text-red-700={error}
		>
			{label}
			{#if required}
				<span class="text-red-500" aria-label="required">*</span>
			{/if}
		</label>
	{/if}

	<div class="relative">
		{@render children?.(fieldId)}
	</div>

	{#if hint && !error}
		<p class="text-sm text-neutral-500">
			{hint}
		</p>
	{/if}

	{#if error}
		<p class="text-sm text-red-600" role="alert">
			{error}
		</p>
	{/if}
</div>