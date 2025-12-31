<script lang="ts">
	/**
	 * UploadDropzone component with drag-drop support
	 * Triggers chunked upload on file drop
	 * Requirements: 2.1, 2.2
	 */

	interface Props {
		/** Whether the dropzone is disabled */
		disabled?: boolean;
		/** Accepted file types (e.g., 'image/*,.pdf') */
		accept?: string;
		/** Allow multiple file selection */
		multiple?: boolean;
		/** Callback when files are selected/dropped */
		onFilesSelected?: (files: File[]) => void;
	}

	let { disabled = false, accept = '', multiple = true, onFilesSelected }: Props = $props();

	let isDragOver = $state(false);
	let fileInput: HTMLInputElement;

	/**
	 * Handle drag over event
	 */
	function handleDragOver(event: DragEvent) {
		if (disabled) return;
		event.preventDefault();
		event.stopPropagation();
		isDragOver = true;
	}

	/**
	 * Handle drag leave event
	 */
	function handleDragLeave(event: DragEvent) {
		event.preventDefault();
		event.stopPropagation();
		isDragOver = false;
	}

	/**
	 * Handle drop event
	 */
	function handleDrop(event: DragEvent) {
		event.preventDefault();
		event.stopPropagation();
		isDragOver = false;

		if (disabled) return;

		const files = event.dataTransfer?.files;
		if (files && files.length > 0) {
			const fileArray = multiple ? Array.from(files) : [files[0]];
			onFilesSelected?.(fileArray);
		}
	}

	/**
	 * Handle file input change
	 */
	function handleFileChange(event: Event) {
		const target = event.target as HTMLInputElement;
		const files = target.files;

		if (files && files.length > 0) {
			const fileArray = Array.from(files);
			onFilesSelected?.(fileArray);
		}

		// Reset input so same file can be selected again
		target.value = '';
	}

	/**
	 * Open file picker
	 */
	function openFilePicker() {
		if (!disabled) {
			fileInput?.click();
		}
	}

	/**
	 * Handle keyboard activation
	 */
	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			openFilePicker();
		}
	}
</script>

<div
	class="dropzone"
	class:drag-over={isDragOver}
	class:disabled
	ondragover={handleDragOver}
	ondragleave={handleDragLeave}
	ondrop={handleDrop}
	onclick={openFilePicker}
	onkeydown={handleKeyDown}
	tabindex={disabled ? -1 : 0}
	role="button"
	aria-label="Upload files by clicking or dragging"
	aria-disabled={disabled}
>
	<input
		bind:this={fileInput}
		type="file"
		{accept}
		{multiple}
		{disabled}
		onchange={handleFileChange}
		class="file-input"
		aria-hidden="true"
		tabindex="-1"
	/>

	<div class="dropzone-content">
		<div class="upload-icon">
			{#if isDragOver}
				<svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3-3m0 0l3 3m-3-3v12"
					/>
				</svg>
			{:else}
				<svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
					/>
				</svg>
			{/if}
		</div>

		<div class="dropzone-text">
			{#if isDragOver}
				<span class="primary-text">Drop files here</span>
			{:else}
				<span class="primary-text">
					<span class="link-text">Click to upload</span> or drag and drop
				</span>
				<span class="secondary-text">
					{#if accept}
						Accepted: {accept}
					{:else}
						Any file type
					{/if}
				</span>
			{/if}
		</div>
	</div>
</div>

<style>
	.dropzone {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 2rem;
		border: 2px dashed #d1d5db;
		border-radius: 0.5rem;
		background: #f9fafb;
		cursor: pointer;
		transition: all 0.2s ease;
		min-height: 150px;
	}

	.dropzone:hover:not(.disabled) {
		border-color: #9ca3af;
		background: #f3f4f6;
	}

	.dropzone:focus {
		outline: 2px solid #3b82f6;
		outline-offset: 2px;
	}

	.dropzone.drag-over {
		border-color: #3b82f6;
		background: #eff6ff;
		border-style: solid;
	}

	.dropzone.disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.file-input {
		position: absolute;
		width: 1px;
		height: 1px;
		padding: 0;
		margin: -1px;
		overflow: hidden;
		clip: rect(0, 0, 0, 0);
		white-space: nowrap;
		border: 0;
	}

	.dropzone-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.75rem;
		text-align: center;
	}

	.upload-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 3rem;
		height: 3rem;
		color: #9ca3af;
	}

	.drag-over .upload-icon {
		color: #3b82f6;
	}

	.icon {
		width: 100%;
		height: 100%;
	}

	.dropzone-text {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.primary-text {
		font-size: 0.875rem;
		color: #4b5563;
	}

	.link-text {
		color: #3b82f6;
		font-weight: 500;
	}

	.secondary-text {
		font-size: 0.75rem;
		color: #9ca3af;
	}

	/* Dark mode */
	@media (prefers-color-scheme: dark) {
		.dropzone {
			border-color: #4b5563;
			background: #1f2937;
		}

		.dropzone:hover:not(.disabled) {
			border-color: #6b7280;
			background: #374151;
		}

		.dropzone.drag-over {
			border-color: #60a5fa;
			background: #1e3a5f;
		}

		.upload-icon {
			color: #6b7280;
		}

		.drag-over .upload-icon {
			color: #60a5fa;
		}

		.primary-text {
			color: #d1d5db;
		}

		.link-text {
			color: #60a5fa;
		}

		.secondary-text {
			color: #6b7280;
		}
	}
</style>
