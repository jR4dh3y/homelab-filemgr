<script lang="ts">
	/**
	 * FilePreview - Modal component for previewing files
	 */
	import { X, Download, Maximize2, Minimize2 } from 'lucide-svelte';
	import type { FileInfo } from '$lib/api/files';
	import { getPreviewUrl, getDownloadUrl } from '$lib/api/files';
	import { getPreviewType, type PreviewType } from '$lib/utils/fileTypes';
	import { formatFileSize } from '$lib/utils/format';
	import { Button } from '$lib/components/ui';
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

	const headerBtnClass =
		'w-8 h-8 flex items-center justify-center bg-transparent border-none rounded text-text-secondary cursor-pointer transition-all duration-100 hover:bg-surface-elevated hover:text-text-primary';
</script>

<svelte:window onkeydown={handleKeydown} />

{#if file}
	<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
	<div
		class="fixed inset-0 bg-black/85 flex items-center justify-center z-[1000] {isFullscreen ? 'p-0' : 'p-10'}"
		onclick={handleBackdropClick}
		onkeydown={handleKeydown}
		role="dialog"
		aria-modal="true"
		aria-label="File preview"
		tabindex="-1"
	>
		<div
			class="flex flex-col w-full h-full bg-surface-primary rounded-lg overflow-hidden shadow-2xl {isFullscreen
				? 'max-w-none max-h-none rounded-none'
				: 'max-w-[1200px] max-h-[90vh]'}"
		>
			<!-- Header -->
			<header class="flex items-center justify-between px-4 py-3 bg-surface-secondary border-b border-border-primary shrink-0">
				<div class="flex items-center gap-3 min-w-0">
					<span class="text-sm font-medium text-text-primary overflow-hidden text-ellipsis whitespace-nowrap" title={file.name}>
						{file.name}
					</span>
					<span class="text-xs text-text-secondary shrink-0">{formatFileSize(file.size)}</span>
				</div>
				<div class="flex items-center gap-1">
					<button type="button" class={headerBtnClass} onclick={handleDownload} title="Download">
						<Download size={18} />
					</button>
					<button
						type="button"
						class={headerBtnClass}
						onclick={toggleFullscreen}
						title={isFullscreen ? 'Exit fullscreen' : 'Fullscreen'}
					>
						{#if isFullscreen}
							<Minimize2 size={18} />
						{:else}
							<Maximize2 size={18} />
						{/if}
					</button>
					<button type="button" class="{headerBtnClass} hover:bg-danger hover:text-white" onclick={onClose} title="Close">
						<X size={18} />
					</button>
				</div>
			</header>

			<!-- Content -->
			<main class="flex-1 overflow-auto flex items-center justify-center">
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
					<div class="flex flex-col items-center gap-4 text-text-secondary text-sm">
						<p>Preview not available for this file type</p>
						<Button variant="primary" onclick={handleDownload}>
							<Download size={20} />
							Download File
						</Button>
					</div>
				{/if}
			</main>
		</div>
	</div>
{/if}
