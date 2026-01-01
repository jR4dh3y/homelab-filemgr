<script lang="ts">
	/**
	 * UploadProgress component showing active uploads
	 * Displays progress bars and cancel buttons
	 */

	import type { UploadProgress as UploadProgressType } from '$lib/utils/upload';
	import { formatFileSize, formatPercentage } from '$lib/utils/format';
	import { X } from 'lucide-svelte';
	import { ProgressBar } from '$lib/components/ui';

	interface Props {
		uploads: UploadProgressType[];
		onCancel?: (uploadId: string) => void;
		onRemove?: (uploadId: string) => void;
		showCompleted?: boolean;
	}

	let { uploads = [], onCancel, onRemove, showCompleted = true }: Props = $props();

	const filteredUploads = $derived(
		showCompleted ? uploads : uploads.filter((u) => u.status !== 'complete' && u.status !== 'cancelled')
	);

	function getProgressVariant(status: UploadProgressType['status']): 'default' | 'success' | 'warning' | 'danger' {
		switch (status) {
			case 'complete': return 'success';
			case 'error': return 'danger';
			default: return 'default';
		}
	}

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

	function canCancel(status: UploadProgressType['status']): boolean {
		return status === 'pending' || status === 'uploading';
	}

	function canRemove(status: UploadProgressType['status']): boolean {
		return status === 'complete' || status === 'error' || status === 'cancelled';
	}
</script>

{#if filteredUploads.length > 0}
	<div class="bg-surface-secondary border border-border-primary rounded-lg overflow-hidden">
		<div class="px-4 py-3 bg-surface-primary border-b border-border-secondary">
			<h3 class="m-0 text-sm font-semibold text-text-primary">
				Uploads ({filteredUploads.length})
			</h3>
		</div>

		<ul class="list-none m-0 p-0 max-h-[300px] overflow-y-auto" role="list">
			{#each filteredUploads as upload (upload.uploadId)}
				<li class="flex flex-col gap-2 px-4 py-3 border-b border-border-secondary last:border-b-0">
					<div class="flex justify-between items-start gap-2">
						<div class="flex flex-col gap-0.5 min-w-0 flex-1">
							<span class="text-sm font-medium text-text-primary overflow-hidden text-ellipsis whitespace-nowrap" title={upload.fileName}>
								{upload.fileName}
							</span>
							<span class="text-xs text-text-muted">{formatFileSize(upload.totalSize)}</span>
						</div>
						<span class="text-xs shrink-0 {upload.status === 'uploading' ? 'text-accent' : ''} {upload.status === 'complete' ? 'text-success' : ''} {upload.status === 'error' ? 'text-danger' : 'text-text-muted'}">
							{getStatusText(upload)}
						</span>
					</div>

					<ProgressBar value={upload.percentage} variant={getProgressVariant(upload.status)} size="sm" />

					<div class="flex justify-end">
						{#if canCancel(upload.status)}
							<button
								type="button"
								class="flex items-center justify-center w-6 h-6 p-0 border-none bg-transparent rounded cursor-pointer text-text-muted transition-all hover:bg-danger/20 hover:text-danger"
								onclick={() => onCancel?.(upload.uploadId)}
								aria-label="Cancel upload of {upload.fileName}"
							>
								<X size={16} />
							</button>
						{:else if canRemove(upload.status)}
							<button
								type="button"
								class="flex items-center justify-center w-6 h-6 p-0 border-none bg-transparent rounded cursor-pointer text-text-muted transition-all hover:bg-surface-elevated hover:text-text-secondary"
								onclick={() => onRemove?.(upload.uploadId)}
								aria-label="Remove {upload.fileName} from list"
							>
								<X size={16} />
							</button>
						{/if}
					</div>
				</li>
			{/each}
		</ul>
	</div>
{/if}
