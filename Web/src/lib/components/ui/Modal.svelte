<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { fade, scale } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';

	const dispatch = createEventDispatcher();

	interface Props {
		isOpen?: boolean;
		title?: string;
		size?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
		closable?: boolean;
		class?: string;
		children?: any;
	}

	let {
		isOpen = false,
		title = '',
		size = 'md',
		closable = true,
		class: className = '',
		children
	}: Props = $props();

	// Size classes
	const sizeClasses = {
		sm: 'max-w-sm',
		md: 'max-w-md',
		lg: 'max-w-lg',
		xl: 'max-w-2xl',
		full: 'max-w-full mx-4'
	};

	function closeModal() {
		if (!closable) return;
		dispatch('close');
	}

	function handleBackdropClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			closeModal();
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && closable) {
			closeModal();
		}
	}
</script>

{#if isOpen}
	<!-- Modal overlay -->
	<div
		class="fixed inset-0 z-50 overflow-y-auto"
		role="dialog"
		aria-modal="true"
		aria-labelledby={title ? 'modal-title' : undefined}
		onkeydown={handleKeydown}
		transition:fade={{ duration: 200 }}
	>
		<!-- Background overlay -->
		<div
			class="fixed inset-0 bg-neutral-900 bg-opacity-50 transition-opacity"
			onclick={handleBackdropClick}
		></div>

		<!-- Modal container -->
		<div class="flex min-h-full items-center justify-center p-4">
			<!-- Modal content -->
			<div
				class="relative w-full {sizeClasses[size]} rounded-lg bg-white shadow-xl {className}"
				transition:scale={{ duration: 200, easing: cubicOut, start: 0.95 }}
			>
				{#if title || closable}
					<!-- Modal header -->
					<div class="flex items-center justify-between p-6 pb-4">
						{#if title}
							<h2 id="modal-title" class="text-lg font-semibold text-neutral-900">
								{title}
							</h2>
						{/if}
						
						{#if closable}
							<button
								type="button"
								class="rounded-full p-1 text-neutral-400 hover:bg-neutral-100 hover:text-neutral-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
								onclick={closeModal}
								aria-label="Close modal"
							>
								<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
								</svg>
							</button>
						{/if}
					</div>
				{/if}

				<!-- Modal body -->
				<div class="px-6 {title || closable ? 'pb-6' : 'py-6'}">
					{@render children?.()}
				</div>
			</div>
		</div>
	</div>
{/if}