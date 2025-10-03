<script lang="ts">
	interface Props {
		variant?: 'default' | 'outline' | 'elevated' | 'gradient';
		padding?: 'none' | 'sm' | 'md' | 'lg';
		class?: string;
		children?: any;
	}

	let {
		variant = 'default',
		padding = 'md',
		class: className = '',
		children
	}: Props = $props();

	// Base card classes
	const baseClasses = 'rounded-lg';

	// Variant classes
	const variantClasses = {
		default: 'bg-white border border-neutral-200',
		outline: 'bg-white border-2 border-neutral-300',
		elevated: 'bg-white shadow-lg border border-neutral-100',
		gradient: 'bg-gradient-to-br from-primary-50 to-primary-100 border border-primary-200'
	};

	// Padding classes
	const paddingClasses = {
		none: '',
		sm: 'p-4',
		md: 'p-6',
		lg: 'p-8'
	};

	// Combine all classes
	const cardClasses = [
		baseClasses,
		variantClasses[variant],
		paddingClasses[padding],
		className
	].join(' ');
</script>

<div class={cardClasses}>
	{@render children?.()}
</div>