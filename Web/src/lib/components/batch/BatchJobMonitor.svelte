<script lang="ts">
	import { onMount, onDestroy, createEventDispatcher } from 'svelte';
	import { getBatchJobStatus, getBatchJobResults, cancelBatchJob, type BatchJob, type BatchProgress } from '$lib/api';
	import { fade, fly, scale, slide } from 'svelte/transition';
	import { cubicOut, backOut } from 'svelte/easing';

	const dispatch = createEventDispatcher();

	interface Props {
		job?: BatchJob | null;
		session?: any;
		autoStart?: boolean;
		pollInterval?: number;
		class?: string;
	}

	let {
		job = null,
		session,
		autoStart = true,
		pollInterval = 2000,
		class: className = ''
	}: Props = $props();

	// State
	let currentJob = $state<BatchJob | null>(job);
	let progress = $state<BatchProgress | null>(null);
	let isMonitoring = $state(false);
	let results = $state<any[]>([]);
	let error = $state<string>('');
	let monitorTimer: ReturnType<typeof setInterval> | null = null;

	// Watch for job changes
	$effect(() => {
		if (job && job !== currentJob) {
			currentJob = job;
			if (autoStart) {
				startMonitoring();
			}
		}
	});

	// Start monitoring job progress
	async function startMonitoring() {
		if (!currentJob || isMonitoring) return;

		isMonitoring = true;
		error = '';
		
		try {
			await updateJobStatus();
			
			// Start polling if job is still running
			if (currentJob && (currentJob.status === 'running' || currentJob.status === 'pending')) {
				monitorTimer = setInterval(updateJobStatus, pollInterval);
			}
		} catch (err) {
			error = 'Failed to start monitoring job';
			console.error('Error starting job monitoring:', err);
			isMonitoring = false;
		}
	}

	// Stop monitoring
	function stopMonitoring() {
		if (monitorTimer) {
			clearInterval(monitorTimer);
			monitorTimer = null;
		}
		isMonitoring = false;
	}

	// Update job status
	async function updateJobStatus() {
		if (!currentJob) return;

		try {
			const status = await getBatchJobStatus(currentJob.id, session);
			
			// Update job status
			currentJob = { ...currentJob, ...status };
			progress = status.progress;

			// Check if job is completed or failed
			if (status.status === 'completed' || status.status === 'failed' || status.status === 'cancelled') {
				stopMonitoring();
				
				if (status.status === 'completed') {
					await loadResults();
					dispatch('completed', { job: currentJob, results });
				} else if (status.status === 'failed') {
					dispatch('failed', { job: currentJob, error: status.error });
				} else if (status.status === 'cancelled') {
					dispatch('cancelled', { job: currentJob });
				}
			} else {
				dispatch('progress', { job: currentJob, progress });
			}
		} catch (err) {
			console.error('Error updating job status:', err);
			error = 'Failed to update job status';
		}
	}

	// Load job results
	async function loadResults() {
		if (!currentJob) return;

		try {
			results = await getBatchJobResults(currentJob.id, session);
		} catch (err) {
			console.error('Error loading results:', err);
			error = 'Failed to load job results';
		}
	}

	// Cancel job
	async function handleCancelJob() {
		if (!currentJob) return;

		try {
			await cancelBatchJob(currentJob.id, session);
			stopMonitoring();
			dispatch('cancelled', { job: currentJob });
		} catch (err) {
			console.error('Error cancelling job:', err);
			error = 'Failed to cancel job';
		}
	}

	// Get status color
	function getStatusColor(status: string): string {
		switch (status) {
			case 'running':
				return 'text-blue-600 bg-blue-50 border-blue-200';
			case 'completed':
				return 'text-green-600 bg-green-50 border-green-200';
			case 'failed':
				return 'text-red-600 bg-red-50 border-red-200';
			case 'cancelled':
				return 'text-orange-600 bg-orange-50 border-orange-200';
			default:
				return 'text-neutral-600 bg-neutral-50 border-neutral-200';
		}
	}

	// Get status icon
	function getStatusIcon(status: string): string {
		switch (status) {
			case 'running':
				return '🔄';
			case 'completed':
				return '✅';
			case 'failed':
				return '❌';
			case 'cancelled':
				return '⏹️';
			default:
				return '⏳';
		}
	}

	// Cleanup on destroy
	onDestroy(() => {
		stopMonitoring();
	});

	// Public methods
	export function start(newJob: BatchJob) {
		currentJob = newJob;
		startMonitoring();
	}

	export function stop() {
		stopMonitoring();
	}
</script>

{#if currentJob}
	<div class="space-y-4 {className}" transition:fly={{ y: -10, duration: 400 }}>
		<!-- Job Header -->
		<div class="flex items-center justify-between rounded-lg border p-4 {getStatusColor(currentJob.status || 'pending')}">
			<div class="flex items-center space-x-3">
				<div class="text-2xl">
					{getStatusIcon(currentJob.status || 'pending')}
				</div>
				<div>
					<h3 class="font-semibold">Batch Job: {currentJob.id}</h3>
					<p class="text-sm opacity-80">
						Status: <span class="font-medium">{(currentJob.status || 'pending').toUpperCase()}</span>
					</p>
				</div>
			</div>
			
			{#if currentJob.status === 'running' || currentJob.status === 'pending'}
				<button
					class="rounded-lg bg-red-100 px-3 py-1 text-sm font-medium text-red-700 hover:bg-red-200 transition-colors"
					onclick={handleCancelJob}
					transition:scale={{ start: 0.95, duration: 200 }}
				>
					Cancel Job
				</button>
			{/if}
		</div>

		<!-- Progress Bar -->
		{#if progress && currentJob.status === 'running'}
			<div class="rounded-lg border border-blue-200 bg-blue-50 p-4" transition:slide={{ duration: 300 }}>
				<div class="space-y-3">
					<div class="flex justify-between text-sm text-blue-700">
						<span>Processing: {progress.processed}/{progress.total}</span>
						<span>{Math.round((progress.processed / progress.total) * 100)}%</span>
					</div>
					
					<!-- Progress Bar -->
					<div class="h-2 w-full overflow-hidden rounded-full bg-blue-200">
						<div 
							class="h-full bg-blue-600 transition-all duration-500 ease-out"
							style="width: {(progress.processed / progress.total) * 100}%"
						></div>
					</div>
					
					<!-- Additional Stats -->
					<div class="grid grid-cols-3 gap-4 text-xs text-blue-600">
						<div class="text-center">
							<div class="font-semibold">{progress.processed}</div>
							<div>Processed</div>
						</div>
						<div class="text-center">
							<div class="font-semibold">{progress.errors || 0}</div>
							<div>Errors</div>
						</div>
						<div class="text-center">
							<div class="font-semibold">{progress.total - progress.processed}</div>
							<div>Remaining</div>
						</div>
					</div>

					{#if isMonitoring}
						<div class="flex items-center justify-center text-xs text-blue-600">
							<div class="mr-2 h-3 w-3 animate-spin rounded-full border-2 border-blue-300 border-t-blue-600"></div>
							Monitoring progress...
						</div>
					{/if}
				</div>
			</div>
		{/if}

		<!-- Job Results -->
		{#if currentJob.status === 'completed' && results.length > 0}
			<div class="rounded-lg border border-green-200 bg-green-50 p-4" transition:slide={{ duration: 400 }}>
				<h4 class="mb-3 font-semibold text-green-800">Job Results ({results.length} documents)</h4>
				
				<div class="max-h-64 space-y-2 overflow-y-auto">
					{#each results.slice(0, 20) as result, index}
						<div class="flex items-center justify-between rounded-lg bg-white p-3 text-sm shadow-sm" 
							 transition:slide={{ delay: index * 50, duration: 300 }}>
							<div class="min-w-0 flex-1">
								<div class="font-medium text-neutral-900 truncate">
									{result.filename || result.document_name || result.id}
								</div>
								{#if result.doc_type}
									<div class="text-xs text-neutral-600">Type: {result.doc_type}</div>
								{/if}
							</div>
							<div class="ml-3 flex-shrink-0">
								{#if result.error}
									<span class="text-red-600 font-medium">❌ Failed</span>
								{:else}
									<span class="text-green-600 font-medium">✅ Success</span>
								{/if}
							</div>
						</div>
					{/each}

					{#if results.length > 20}
						<div class="text-center text-xs text-neutral-600 pt-2">
							... and {results.length - 20} more documents
						</div>
					{/if}
				</div>

				<!-- Export Results Button -->
				<div class="mt-3 flex justify-end">
					<button
						class="rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 transition-colors"
						onclick={() => dispatch('export', { results })}
					>
						Export Results
					</button>
				</div>
			</div>
		{/if}

		<!-- Error Display -->
		{#if error}
			<div class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-700" transition:fade>
				<div class="flex items-center">
					<svg class="mr-2 h-5 w-5 text-red-500" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
					</svg>
					<span class="font-medium">{error}</span>
				</div>
			</div>
		{/if}

		<!-- Job Failed -->
		{#if currentJob.status === 'failed'}
			<div class="rounded-lg border border-red-200 bg-red-50 p-4" transition:slide={{ duration: 300 }}>
				<h4 class="mb-2 font-semibold text-red-800">Job Failed</h4>
				<p class="text-sm text-red-700">
					{currentJob.error || 'The batch job failed to complete successfully.'}
				</p>
				
				<div class="mt-3 flex justify-end">
					<button
						class="rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 transition-colors"
						onclick={() => dispatch('retry', { job: currentJob })}
					>
						Retry Job
					</button>
				</div>
			</div>
		{/if}

		<!-- Manual Controls (Development/Debug) -->
		{#if currentJob && !isMonitoring}
			<div class="flex justify-center space-x-3">
				<button
					class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 transition-colors"
					onclick={startMonitoring}
				>
					Start Monitoring
				</button>
				
				{#if currentJob.status === 'completed'}
					<button
						class="rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 transition-colors"
						onclick={loadResults}
					>
						Load Results
					</button>
				{/if}
			</div>
		{/if}
	</div>
{:else}
	<div class="rounded-lg border border-neutral-200 bg-neutral-50 p-8 text-center text-neutral-600">
		<div class="text-4xl mb-3">📊</div>
		<h3 class="font-medium">No Batch Job</h3>
		<p class="text-sm">Start a batch processing job to monitor its progress here.</p>
	</div>
{/if}