<script lang="ts">
	/**
	 * UploadProgress component showing active uploads
	 * Displays progress bars with icon that transforms to X on hover
	 */

	import type { UploadProgress as UploadProgressType } from '$lib/utils/upload';
	import { formatFileSize, formatPercentage } from '$lib/utils/format';
	import { X, Upload } from 'lucide-svelte';
	import { Badge, Button, ProgressBar } from '$lib/components/ui';

	interface Props {
		uploads: UploadProgressType[];
		onCancel?: (uploadId: string) => void;
		onRemove?: (uploadId: string) => void;
		onClearCompleted?: () => void;
		showCompleted?: boolean;
	}

	let { uploads = [], onCancel, onRemove, onClearCompleted, showCompleted = true }: Props = $props();

	const filteredUploads = $derived(
		showCompleted ? uploads : uploads.filter((u) => u.status !== 'complete' && u.status !== 'cancelled')
	);

	const activeCount = $derived(uploads.filter((u) => u.status === 'pending' || u.status === 'uploading').length);
	const completedCount = $derived(uploads.filter((u) => u.status === 'complete' || u.status === 'error' || u.status === 'cancelled').length);

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
			case 'uploading': return formatPercentage(upload.percentage, 0, false);
			case 'complete': return 'Complete';
			case 'error': return upload.error || 'Failed';
			case 'cancelled': return 'Cancelled';
			default: return '';
		}
	}

	function isTerminal(status: UploadProgressType['status']): boolean {
		return status === 'complete' || status === 'error' || status === 'cancelled';
	}
</script>

{#if filteredUploads.length > 0}
	<div class="bg-surface-secondary border border-border-primary rounded-lg overflow-hidden">
		<div class="px-4 py-3 bg-surface-primary border-b border-border-secondary flex items-center justify-between">
			<h3 class="m-0 text-sm font-semibold text-text-primary flex items-center gap-2">
				Uploads
				{#if activeCount > 0}
					<Badge variant="info">{activeCount} active</Badge>
				{/if}
			</h3>
			{#if completedCount > 0 && onClearCompleted}
				<Button variant="ghost" size="sm" onclick={onClearCompleted}>
					Clear done
				</Button>
			{/if}
		</div>

		<ul class="list-none m-0 p-0 max-h-[300px] overflow-y-auto" role="list">
			{#each filteredUploads as upload (upload.uploadId)}
				<li class="px-4 py-3 border-b border-border-secondary last:border-b-0 transition-all hover:bg-surface-tertiary">
					<div class="flex items-stretch gap-3">
						<!-- Icon that transforms to X on hover -->
						<button
							type="button"
							class="group shrink-0 w-16 flex items-center justify-center rounded bg-surface-elevated text-text-secondary border-none cursor-pointer transition-all hover:bg-danger/20 hover:text-danger"
							onclick={() => isTerminal(upload.status) ? onRemove?.(upload.uploadId) : onCancel?.(upload.uploadId)}
							aria-label={isTerminal(upload.status) ? 'Remove from list' : 'Cancel upload'}
						>
							<span class="group-hover:hidden">
								<Upload size={20} />
							</span>
							<span class="hidden group-hover:block">
								<X size={20} />
							</span>
						</button>
						
						<div class="flex-1 min-w-0 flex flex-col gap-1 py-0.5">
							<div class="flex items-center justify-between gap-2">
								<span class="text-sm font-medium text-text-primary overflow-hidden text-ellipsis whitespace-nowrap" title={upload.fileName}>
									{upload.fileName}
								</span>
								<span class="text-xs text-text-muted shrink-0">
									{formatFileSize(upload.totalSize)}
								</span>
							</div>
							<div class="text-xs text-text-muted">
								{#if upload.status === 'uploading'}
									Chunk {upload.currentChunk + 1}/{upload.totalChunks}
								{:else}
									{formatFileSize(upload.uploadedSize)} / {formatFileSize(upload.totalSize)}
								{/if}
							</div>
							
							<!-- Progress bar with inline status -->
							<div class="flex items-center gap-3">
								<div class="flex-1">
									<ProgressBar 
										value={upload.percentage} 
										size="sm" 
										variant={getProgressVariant(upload.status)}
									/>
								</div>
								<span class="text-[11px] shrink-0 {upload.status === 'complete' ? 'text-success' : ''} {upload.status === 'error' ? 'text-danger' : ''} {upload.status === 'uploading' ? 'text-accent' : 'text-text-muted'}">
									{getStatusText(upload)}
								</span>
							</div>
						</div>
					</div>
				</li>
			{/each}
		</ul>
	</div>
{/if}
