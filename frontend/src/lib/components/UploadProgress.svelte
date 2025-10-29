<script lang="ts">
	/**
	 * UploadProgress component showing active uploads
	 * Displays progress bars and cancel buttons
	 * Requirements: 2.2
	 */

	import type { UploadProgress as UploadProgressType } from '$lib/utils/upload';
	import { formatFileSize, formatPercentage } from '$lib/utils/format';

	interface Props {
		/** List of active uploads */
		uploads: UploadProgressType[];
		/** Callback when cancel is clicked */
		onCancel?: (uploadId: string) => void;
		/** Callback when remove is clicked (for completed/failed) */
		onRemove?: (uploadId: string) => void;
		/** Whether to show completed uploads */
		showCompleted?: boolean;
	}

	let {
		uploads = [],
		onCancel,
		onRemove,
		showCompleted = true
	}: Props = $props();

	const filteredUploads = $derived(
		showCompleted 
			? uploads 
			: uploads.filter(u => u.status !== 'complete' && u.status !== 'cancelled')
	);

	/**
	 * Get status color class
	 */
	function getStatusColor(status: UploadProgressType['status']): string {
		switch (status) {
			case 'uploading': return 'status-uploading';
			case 'complete': return 'status-complete';
			case 'error': return 'status-error';
			case 'cancelled': return 'status-cancelled';
			default: return 'status-pending';
		}
	}

	/**
	 * Get status text
	 */
	function getStatusText(upload: UploadProgressType): string {
		switch (upload.status) {
			case 'pending': return 'Waiting...';
			case 'uploading': return `${formatPercentage(upload.percentage, 0, false)} - Chunk ${upload.currentChunk + 1}/${upload.totalChunks}`;
			case 'complete': return 'Complete';
			case 'error': return upload.error || 'Failed';
			case 'cancelled': return 'Cancelled';
			default: return '';
		}
	}

	/**
	 * Check if upload can be cancelled
	 */
	function canCancel(status: UploadProgressType['status']): boolean {
		return status === 'pending' || status === 'uploading';
	}

	/**
	 * Check if upload can be removed
	 */
	function canRemove(status: UploadProgressType['status']): boolean {
		return status === 'complete' || status === 'error' || status === 'cancelled';
	}
</script>

{#if filteredUploads.length > 0}
	<div class="upload-progress-container">
		<div class="header">
			<h3 class="title">Uploads ({filteredUploads.length})</h3>
		</div>

		<ul class="upload-list" role="list">
			{#each filteredUploads as upload (upload.uploadId)}
				<li class="upload-item {getStatusColor(upload.status)}">
					<div class="upload-info">
						<div class="file-info">
							<span class="file-name" title={upload.fileName}>{upload.fileName}</span>
							<span class="file-size">{formatFileSize(upload.totalSize)}</span>
						</div>
						<div class="status-info">
							<span class="status-text">{getStatusText(upload)}</span>
						</div>
					</div>

					<div class="progress-bar-container">
						<div 
							class="progress-bar" 
							style="width: {upload.percentage}%"
							role="progressbar"
							aria-valuenow={upload.percentage}
							aria-valuemin={0}
							aria-valuemax={100}
							aria-label="Upload progress for {upload.fileName}"
						></div>
					</div>

					<div class="actions">
						{#if canCancel(upload.status)}
							<button
								type="button"
								class="action-btn cancel-btn"
								onclick={() => onCancel?.(upload.uploadId)}
								aria-label="Cancel upload of {upload.fileName}"
							>
								<svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
								</svg>
							</button>
						{:else if canRemove(upload.status)}
							<button
								type="button"
								class="action-btn remove-btn"
								onclick={() => onRemove?.(upload.uploadId)}
								aria-label="Remove {upload.fileName} from list"
							>
								<svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
								</svg>
							</button>
						{/if}
					</div>
				</li>
			{/each}
		</ul>
	</div>
{/if}

<style>
	.upload-progress-container {
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
	}

	.upload-list {
		list-style: none;
		margin: 0;
		padding: 0;
		max-height: 300px;
		overflow-y: auto;
	}

	.upload-item {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		border-bottom: 1px solid #e5e7eb;
	}

	.upload-item:last-child {
		border-bottom: none;
	}

	.upload-info {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 0.5rem;
	}

	.file-info {
		display: flex;
		flex-direction: column;
		gap: 0.125rem;
		min-width: 0;
		flex: 1;
	}

	.file-name {
		font-size: 0.875rem;
		font-weight: 500;
		color: #1f2937;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.file-size {
		font-size: 0.75rem;
		color: #6b7280;
	}

	.status-info {
		flex-shrink: 0;
	}

	.status-text {
		font-size: 0.75rem;
		color: #6b7280;
	}

	.status-uploading .status-text { color: #3b82f6; }
	.status-complete .status-text { color: #10b981; }
	.status-error .status-text { color: #ef4444; }
	.status-cancelled .status-text { color: #6b7280; }

	.progress-bar-container {
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

	.status-complete .progress-bar { background: #10b981; }
	.status-error .progress-bar { background: #ef4444; }
	.status-cancelled .progress-bar { background: #9ca3af; }

	.actions {
		display: flex;
		justify-content: flex-end;
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

	/* Dark mode */
	@media (prefers-color-scheme: dark) {
		.upload-progress-container {
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

		.upload-item {
			border-bottom-color: #374151;
		}

		.file-name {
			color: #f3f4f6;
		}

		.file-size,
		.status-text {
			color: #9ca3af;
		}

		.status-uploading .status-text { color: #60a5fa; }
		.status-complete .status-text { color: #34d399; }
		.status-error .status-text { color: #f87171; }

		.progress-bar-container {
			background: #374151;
		}

		.progress-bar {
			background: #60a5fa;
		}

		.status-complete .progress-bar { background: #34d399; }
		.status-error .progress-bar { background: #f87171; }
		.status-cancelled .progress-bar { background: #6b7280; }

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
	}
</style>
