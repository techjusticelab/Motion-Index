<script lang="ts">
	interface Props {
		size?: 'sm' | 'md' | 'lg' | 'xl';
		color?: 'primary' | 'secondary' | 'white' | 'neutral';
		text?: string;
		class?: string;
	}

	let {
		size = 'md',
		color = 'primary',
		text = '',
		class: className = ''
	}: Props = $props();

	// Size classes
	const sizeClasses = {
		sm: 'h-4 w-4',
		md: 'h-6 w-6',
		lg: 'h-8 w-8',
		xl: 'h-12 w-12'
	};

	// Color classes
	const colorClasses = {
		primary: 'text-primary-600',
		secondary: 'text-neutral-600',
		white: 'text-white',
		neutral: 'text-neutral-400'
	};

	// Text size based on spinner size
	const textSizeClasses = {
		sm: 'text-xs',
		md: 'text-sm',
		lg: 'text-base',
		xl: 'text-lg'
	};
</script>

<div class="flex items-center justify-center {className}">
	<div class="flex flex-col items-center gap-2">
		<!-- Spinner -->
		<svg
			class="animate-spin {sizeClasses[size]} {colorClasses[color]}"
			fill="none"
			viewBox="0 0 24 24"
			aria-hidden="true"
		>
			<circle
				class="opacity-25"
				cx="12"
				cy="12"
				r="10"
				stroke="currentColor"
				stroke-width="4"
			></circle>
			<path
				class="opacity-75"
				fill="currentColor"
				d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
			></path>
		</svg>

		<!-- Loading text -->
		{#if text}
			<p class="{textSizeClasses[size]} {colorClasses[color]} font-medium">
				{text}
			</p>
		{/if}
	</div>
</div>