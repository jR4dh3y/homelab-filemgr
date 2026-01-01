<script lang="ts">
	/**
	 * JobMonitor component showing background jobs
	 * Displays progress, status, and cancel option
	 */

	import type { Job, JobType } from '$lib/api/jobs';
	import { isJobActive, isJobTerminal } from '$lib/api/jobs';
	import { formatFileDate, formatPercentage } from '$lib/utils/format';
	import { X, Copy, FolderInput, Trash2, Settings, ArrowRight } from 'lucide-svelte';
	import { Badge, ProgressBar } from '$lib/components/ui';

	interface Props {
		jobs: Job[];
		onCancel?: (jobId: string) => void;
		onRemove?: (jobId: string) => void;
		showCompleted?: boolean;
		maxDisplay?: number;
	}

	let { jobs = [], onCancel, onRemove, showCompleted = true, maxDisplay = 10 }: Props = $props();

	const filteredJobs = $derived(
		(showCompleted ? jobs : jobs.filter((j) => !isJobTerminal(j))).slice(0, maxDisplay)
	);

	const activeCount = $derived(jobs.filter(isJobActive).length);

	function getJobTypeLabel(type: JobType): string {
		switch (type) {
			case 'copy': return 'Copy';
			case 'move': return 'Move';
			case 'delete': return 'Delete';
			default: return 'Job';
		}
	}

	function getStatusText(job: Job): string {
		switch (job.state) {
			case 'pending': return 'Waiting...';
			case 'running': return `${formatPercentage(job.progress, 0, false)}`;
			case 'completed': return 'Completed';
			case 'failed': return job.error || 'Failed';
			case 'cancelled': return 'Cancelled';
			default: return '';
		}
	}

	function getFileName(path: string): string {
		const parts = path.split('/');
		return parts[parts.length - 1] || path;
	}
</script>

{#if filteredJobs.length > 0}
	<div class="bg-surface-secondary border border-border-primary rounded-lg overflow-hidden">
		<div class="px-4 py-3 bg-surface-primary border-b border-border-secondary">
			<h3 class="m-0 text-sm font-semibold text-text-primary flex items-center gap-2">
				Background Jobs
				{#if activeCount > 0}
					<Badge variant="info">{activeCount} active</Badge>
				{/if}
			</h3>
		</div>

		<ul class="list-none m-0 p-0 max-h-[400px] overflow-y-auto" role="list">
			{#each filteredJobs as job (job.id)}
				<li class="px-4 py-3 border-b border-border-secondary last:border-b-0">
					<div class="flex items-start gap-2">
						<span class="text-text-secondary shrink-0">
							{#if job.type === 'copy'}
								<Copy size={20} />
							{:else if job.type === 'move'}
								<FolderInput size={20} />
							{:else if job.type === 'delete'}
								<Trash2 size={20} />
							{:else}
								<Settings size={20} />
							{/if}
						</span>
						<div class="flex-1 min-w-0 flex flex-col gap-0.5">
							<span class="text-[10px] font-semibold text-text-muted uppercase tracking-wide">
								{getJobTypeLabel(job.type)}
							</span>
							<span class="text-sm font-medium text-text-primary overflow-hidden text-ellipsis whitespace-nowrap" title={job.sourcePath}>
								{getFileName(job.sourcePath)}
							</span>
							{#if job.destPath}
								<span class="text-xs text-text-muted overflow-hidden text-ellipsis whitespace-nowrap flex items-center gap-1">
									<ArrowRight size={12} /> {getFileName(job.destPath)}
								</span>
							{/if}
						</div>
						<div class="shrink-0">
							{#if isJobActive(job)}
								<button
									type="button"
									class="flex items-center justify-center w-6 h-6 p-0 border-none bg-transparent rounded cursor-pointer text-text-muted transition-all hover:bg-danger/20 hover:text-danger"
									onclick={() => onCancel?.(job.id)}
									aria-label="Cancel job"
								>
									<X size={16} />
								</button>
							{:else if isJobTerminal(job)}
								<button
									type="button"
									class="flex items-center justify-center w-6 h-6 p-0 border-none bg-transparent rounded cursor-pointer text-text-muted transition-all hover:bg-surface-elevated hover:text-text-secondary"
									onclick={() => onRemove?.(job.id)}
									aria-label="Remove from list"
								>
									<X size={16} />
								</button>
							{/if}
						</div>
					</div>

					{#if job.state === 'running' || job.state === 'pending'}
						<div class="mt-2">
							<ProgressBar value={job.progress} size="sm" />
						</div>
					{/if}

					<div class="flex justify-between items-center mt-2 text-xs">
						<span class="text-text-secondary {job.state === 'running' ? 'text-accent' : ''} {job.state === 'completed' ? 'text-success' : ''} {job.state === 'failed' ? 'text-danger' : ''}">
							{getStatusText(job)}
						</span>
						<span class="text-text-muted">{formatFileDate(job.createdAt)}</span>
					</div>
				</li>
			{/each}
		</ul>

		{#if jobs.length > maxDisplay}
			<div class="px-4 py-2 text-center text-xs text-text-muted bg-surface-primary border-t border-border-secondary">
				+{jobs.length - maxDisplay} more jobs
			</div>
		{/if}
	</div>
{/if}
