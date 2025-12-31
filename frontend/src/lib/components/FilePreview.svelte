<script lang="ts">
	/**
	 * FilePreview - Modal component for previewing files
	 */
	import { X, Download, Maximize2, Minimize2 } from 'lucide-svelte';
	import type { FileInfo } from '$lib/api/files';
	import { getPreviewUrl, getDownloadUrl } from '$lib/api/files';
	import { getPreviewType, type PreviewType } from '$lib/utils/fileTypes';
	import { formatFileSize } from '$lib/utils/format';
	import VideoPreview from './preview/VideoPreview.svelte';
	import AudioPreview from './preview/AudioPreview.svelte';
	import ImagePreview from './preview/ImagePreview.svelte';
	import CodePreview from './preview/CodePreview.svelte';
	import PdfPreview from './preview/PdfPreview.svelte';

	interface Props {
		file: FileInfo | null;
		onClose: () => void;
	}

	let { file, onClose }: Props = $props();

	let isFullscreen = $state(false);

	const previewType = $derived<PreviewType>(file ? getPreviewType(file.name) : 'unsupported');
	const previewUrl = $derived(file ? getPreviewUrl(file.path) : '');
	const downloadUrl = $derived(file ? getDownloadUrl(file.path) : '');

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			if (isFullscreen) {
				isFullscreen = false;
			} else {
				onClose();
			}
		}
	}

	function handleBackdropClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			onClose();
		}
	}

	function toggleFullscreen() {
		isFullscreen = !isFullscreen;
	}

	function handleDownload() {
		if (downloadUrl) {
			window.open(downloadUrl, '_blank');
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if file}
	<div 
		class="preview-overlay" 
		class:fullscreen={isFullscreen}
		onclick={handleBackdropClick}
		role="dialog"
		aria-modal="true"
		aria-label="File preview"
	>
		<div class="preview-container">
			<!-- Header -->
			<header class="preview-header">
				<div class="file-info">
					<span class="file-name" title={file.name}>{file.name}</span>
					<span class="file-size">{formatFileSize(file.size)}</span>
				</div>
				<div class="header-actions">
					<button 
						type="button" 
						class="header-btn" 
						onclick={handleDownload}
						title="Download"
					>
						<Download size={18} />
					</button>
					<button 
						type="button" 
						class="header-btn" 
						onclick={toggleFullscreen}
						title={isFullscreen ? 'Exit fullscreen' : 'Fullscreen'}
					>
						{#if isFullscreen}
							<Minimize2 size={18} />
						{:else}
							<Maximize2 size={18} />
						{/if}
					</button>
					<button 
						type="button" 
						class="header-btn close-btn" 
						onclick={onClose}
						title="Close"
					>
						<X size={18} />
					</button>
				</div>
			</header>

			<!-- Content -->
			<main class="preview-content">
				{#if previewType === 'video'}
					<VideoPreview url={previewUrl} filename={file.name} />
				{:else if previewType === 'audio'}
					<AudioPreview url={previewUrl} filename={file.name} />
				{:else if previewType === 'image'}
					<ImagePreview url={previewUrl} filename={file.name} />
				{:else if previewType === 'pdf'}
					<PdfPreview url={previewUrl} filename={file.name} />
				{:else if previewType === 'code' || previewType === 'text'}
					<CodePreview url={previewUrl} filename={file.name} />
				{:else}
					<div class="unsupported">
						<p>Preview not available for this file type</p>
						<button type="button" class="download-btn" onclick={handleDownload}>
							<Download size={20} />
							Download File
						</button>
					</div>
				{/if}
			</main>
		</div>
	</div>
{/if}

<style>
	.preview-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.85);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 40px;
	}

	.preview-overlay.fullscreen {
		padding: 0;
	}

	.preview-container {
		display: flex;
		flex-direction: column;
		width: 100%;
		height: 100%;
		max-width: 1200px;
		max-height: 90vh;
		background: #1e1e1e;
		border-radius: 8px;
		overflow: hidden;
		box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
	}

	.preview-overlay.fullscreen .preview-container {
		max-width: none;
		max-height: none;
		border-radius: 0;
	}

	.preview-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 12px 16px;
		background: #252525;
		border-bottom: 1px solid #333;
		flex-shrink: 0;
	}

	.file-info {
		display: flex;
		align-items: center;
		gap: 12px;
		min-width: 0;
	}

	.file-name {
		font-size: 14px;
		font-weight: 500;
		color: #e0e0e0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.file-size {
		font-size: 12px;
		color: #888;
		flex-shrink: 0;
	}

	.header-actions {
		display: flex;
		align-items: center;
		gap: 4px;
	}

	.header-btn {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: none;
		border-radius: 4px;
		color: #888;
		cursor: pointer;
		transition: all 0.1s ease;
	}

	.header-btn:hover {
		background: #333;
		color: #ccc;
	}

	.close-btn:hover {
		background: #dc3545;
		color: #fff;
	}

	.preview-content {
		flex: 1;
		overflow: auto;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.unsupported {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		color: #888;
		font-size: 14px;
	}

	.download-btn {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 10px 20px;
		background: #2d4a6f;
		border: none;
		border-radius: 6px;
		color: #fff;
		font-size: 14px;
		cursor: pointer;
		transition: background-color 0.15s ease;
	}

	.download-btn:hover {
		background: #345580;
	}
</style>
