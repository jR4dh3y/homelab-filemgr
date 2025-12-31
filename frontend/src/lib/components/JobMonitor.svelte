<script lang="ts">
	/**
	 * JobMonitor component showing background jobs
	 * Displays progress, status, and cancel option
	 * Requirements: 4.2, 4.4, 4.5
	 */

	import type { Job, JobState, JobType } from '$lib/api/jobs';
	import { isJobActive, isJobTerminal } from '$lib/api/jobs';
	import { formatFileDate, formatPercentage } from '$lib/utils/format';

	interface Props {
		/** List of jobs to display */
		jobs: Job[];
		/** Callback when cancel is clicked */
		onCancel?: (jobId: string) => void;
		/** Callback when remove is clicked (for completed/failed) */
		onRemove?: (jobId: string) => void;
		/** Whether to show completed jobs */
		showCompleted?: boolean;
		/** Maximum jobs to display */
		maxDisplay?: number;
	}

	let { jobs = [], onCancel, onRemove, showCompleted = true, maxDisplay = 10 }: Props = $props();

	const filteredJobs = $derived(
		(showCompleted ? jobs : jobs.filter((j) => !isJobTerminal(j))).slice(0, maxDisplay)
	);

	const activeCount = $derived(jobs.filter(isJobActive).length);

	/**
	 * Get job type icon
	 */
	function getJobIcon(type: JobType): string {
		switch (type) {
			case 'copy':
				return '📋';
			case 'move':
				return '📦';
			case 'delete':
				return '🗑️';
			default:
				return '⚙️';
		}
	}

	/**
	 * Get job type label
	 */
	function getJobTypeLabel(type: JobType): string {
		switch (type) {
			case 'copy':
				return 'Copy';
			case 'move':
				return 'Move';
			case 'delete':
				return 'Delete';
			default:
				return 'Job';
		}
	}

	/**
	 * Get status color class
	 */
	function getStatusClass(state: JobState): string {
		switch (state) {
			case 'pending':
				return 'status-pending';
			case 'running':
				return 'status-running';
			case 'completed':
				return 'status-completed';
			case 'failed':
				return 'status-failed';
			case 'cancelled':
				return 'status-cancelled';
			default:
				return '';
		}
	}

	/**
	 * Get status text
	 */
	function getStatusText(job: Job): string {
		switch (job.state) {
			case 'pending':
				return 'Waiting...';
			case 'running':
				return `${formatPercentage(job.progress, 0, false)}`;
			case 'completed':
				return 'Completed';
			case 'failed':
				return job.error || 'Failed';
			case 'cancelled':
				return 'Cancelled';
			default:
				return '';
		}
	}

	/**
	 * Get source file name from path
	 */
	function getFileName(path: string): string {
		const parts = path.split('/');
		return parts[parts.length - 1] || path;
	}
</script>

{#if filteredJobs.length > 0}
	<div class="job-monitor">
		<div class="header">
			<h3 class="title">
				Background Jobs
				{#if activeCount > 0}
					<span class="active-badge">{activeCount} active</span>
				{/if}
			</h3>
		</div>

		<ul class="job-list" role="list">
			{#each filteredJobs as job (job.id)}
				<li class="job-item {getStatusClass(job.state)}">
					<div class="job-header">
						<span class="job-icon">{getJobIcon(job.type)}</span>
						<div class="job-info">
							<span class="job-type">{getJobTypeLabel(job.type)}</span>
							<span class="job-source" title={job.sourcePath}>{getFileName(job.sourcePath)}</span>
							{#if job.destPath}
								<span class="job-dest">→ {getFileName(job.destPath)}</span>
							{/if}
						</div>
						<div class="job-actions">
							{#if isJobActive(job)}
								<button
									type="button"
									class="action-btn cancel-btn"
									onclick={() => onCancel?.(job.id)}
									aria-label="Cancel job"
								>
									<svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M6 18L18 6M6 6l12 12"
										/>
									</svg>
								</button>
							{:else if isJobTerminal(job)}
								<button
									type="button"
									class="action-btn remove-btn"
									onclick={() => onRemove?.(job.id)}
									aria-label="Remove from list"
								>
									<svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M6 18L18 6M6 6l12 12"
										/>
									</svg>
								</button>
							{/if}
						</div>
					</div>

					{#if job.state === 'running' || job.state === 'pending'}
						<div class="progress-container">
							<div
								class="progress-bar"
								style="width: {job.progress}%"
								role="progressbar"
								aria-valuenow={job.progress}
								aria-valuemin={0}
								aria-valuemax={100}
								aria-label="Job progress"
							></div>
						</div>
					{/if}

					<div class="job-footer">
						<span class="status-text">{getStatusText(job)}</span>
						<span class="job-time">{formatFileDate(job.createdAt)}</span>
					</div>
				</li>
			{/each}
		</ul>

		{#if jobs.length > maxDisplay}
			<div class="more-jobs">
				+{jobs.length - maxDisplay} more jobs
			</div>
		{/if}
	</div>
{/if}

<style>
	.job-monitor {
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		overflow: hidden;
	}

	.header {
		padding: 0.75rem 1rem;
		background: #f9fafb;
		border-bottom: 1px solid #e5e7eb;
	}

	.title {
		margin: 0;
		font-size: 0.875rem;
		font-weight: 600;
		color: #374151;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.active-badge {
		font-size: 0.75rem;
		font-weight: 500;
		padding: 0.125rem 0.5rem;
		background: #dbeafe;
		color: #1d4ed8;
		border-radius: 9999px;
	}

	.job-list {
		list-style: none;
		margin: 0;
		padding: 0;
		max-height: 400px;
		overflow-y: auto;
	}

	.job-item {
		padding: 0.75rem 1rem;
		border-bottom: 1px solid #e5e7eb;
	}

	.job-item:last-child {
		border-bottom: none;
	}

	.job-header {
		display: flex;
		align-items: flex-start;
		gap: 0.5rem;
	}

	.job-icon {
		font-size: 1.25rem;
		flex-shrink: 0;
	}

	.job-info {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
		gap: 0.125rem;
	}

	.job-type {
		font-size: 0.75rem;
		font-weight: 600;
		color: #6b7280;
		text-transform: uppercase;
		letter-spacing: 0.025em;
	}

	.job-source {
		font-size: 0.875rem;
		font-weight: 500;
		color: #1f2937;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.job-dest {
		font-size: 0.75rem;
		color: #6b7280;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.job-actions {
		flex-shrink: 0;
	}

	.action-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 1.5rem;
		height: 1.5rem;
		padding: 0;
		border: none;
		background: transparent;
		border-radius: 0.25rem;
		cursor: pointer;
		color: #9ca3af;
		transition: all 0.15s;
	}

	.action-btn:hover {
		background: #f3f4f6;
		color: #6b7280;
	}

	.cancel-btn:hover {
		background: #fef2f2;
		color: #ef4444;
	}

	.btn-icon {
		width: 1rem;
		height: 1rem;
	}

	.progress-container {
		margin-top: 0.5rem;
		height: 4px;
		background: #e5e7eb;
		border-radius: 2px;
		overflow: hidden;
	}

	.progress-bar {
		height: 100%;
		background: #3b82f6;
		border-radius: 2px;
		transition: width 0.3s ease;
	}

	.status-running .progress-bar {
		background: linear-gradient(90deg, #3b82f6 0%, #60a5fa 50%, #3b82f6 100%);
		background-size: 200% 100%;
		animation: shimmer 1.5s infinite;
	}

	@keyframes shimmer {
		0% {
			background-position: 200% 0;
		}
		100% {
			background-position: -200% 0;
		}
	}

	.job-footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-top: 0.5rem;
		font-size: 0.75rem;
	}

	.status-text {
		color: #6b7280;
	}

	.status-running .status-text {
		color: #3b82f6;
	}
	.status-completed .status-text {
		color: #10b981;
	}
	.status-failed .status-text {
		color: #ef4444;
	}
	.status-cancelled .status-text {
		color: #6b7280;
	}

	.job-time {
		color: #9ca3af;
	}

	.more-jobs {
		padding: 0.5rem 1rem;
		text-align: center;
		font-size: 0.75rem;
		color: #6b7280;
		background: #f9fafb;
		border-top: 1px solid #e5e7eb;
	}

	/* Dark mode */
	@media (prefers-color-scheme: dark) {
		.job-monitor {
			background: #1f2937;
			border-color: #374151;
		}

		.header {
			background: #111827;
			border-bottom-color: #374151;
		}

		.title {
			color: #e5e7eb;
		}

		.active-badge {
			background: #1e3a5f;
			color: #93c5fd;
		}

		.job-item {
			border-bottom-color: #374151;
		}

		.job-type {
			color: #9ca3af;
		}

		.job-source {
			color: #f3f4f6;
		}

		.job-dest {
			color: #9ca3af;
		}

		.action-btn {
			color: #6b7280;
		}

		.action-btn:hover {
			background: #374151;
			color: #9ca3af;
		}

		.cancel-btn:hover {
			background: #450a0a;
			color: #f87171;
		}

		.progress-container {
			background: #374151;
		}

		.progress-bar {
			background: #60a5fa;
		}

		.status-text {
			color: #9ca3af;
		}

		.status-running .status-text {
			color: #60a5fa;
		}
		.status-completed .status-text {
			color: #34d399;
		}
		.status-failed .status-text {
			color: #f87171;
		}

		.job-time {
			color: #6b7280;
		}

		.more-jobs {
			background: #111827;
			border-top-color: #374151;
			color: #9ca3af;
		}
	}
</style>
