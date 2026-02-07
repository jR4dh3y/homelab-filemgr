<script lang="ts">
	/**
	 * UploadPanel - Google Drive style floating upload panel
	 * Fixed position bottom-right, minimizable, manually closed
	 */
	import { uploadStore } from '$lib/stores/upload.svelte';
	import UploadProgress from '$lib/components/UploadProgress.svelte';
	import { ChevronDown, ChevronUp, X } from 'lucide-svelte';
	import { slide } from 'svelte/transition';

	/** Whether the panel is minimized */
	let isMinimized = $state(false);

	/** Whether the panel is visible (user can close it) */
	let isVisible = $state(true);

	/** Track if we should auto-show on new uploads */
	let lastUploadCount = $state(0);

	// Auto-show panel when new uploads are added
	$effect(() => {
		const currentCount = uploadStore.uploads.length;
		if (currentCount > lastUploadCount && !isVisible) {
			isVisible = true;
			isMinimized = false;
		}
		lastUploadCount = currentCount;
	});

	function handleToggleMinimize() {
		isMinimized = !isMinimized;
	}

	function handleClose() {
		isVisible = false;
	}

	function handleCancel(uploadId: string) {
		uploadStore.cancel(uploadId);
	}

	function handleRemove(uploadId: string) {
		uploadStore.remove(uploadId);
	}

	function handleClearCompleted() {
		uploadStore.clearFinished();
	}

	// Header text shows upload status
	const headerText = $derived.by(() => {
		const active = uploadStore.activeCount;
		const total = uploadStore.uploads.length;

		if (active === 0 && total === 0) {
			return 'Uploads';
		}

		if (active === 0) {
			return `${total} upload${total !== 1 ? 's' : ''} complete`;
		}

		if (total === 1) {
			return 'Uploading 1 file';
		}

		return `Uploading ${active} of ${total}`;
	});

	// Only show if visible AND has uploads
	const shouldShow = $derived(isVisible && uploadStore.hasUploads);
</script>

{#if shouldShow}
	<div
		class="fixed bottom-4 right-4 z-50 w-[380px] bg-surface-secondary border border-border-primary rounded-lg shadow-2xl overflow-hidden"
		role="region"
		aria-label="Upload progress"
	>
		<!-- Header -->
		<div
			class="flex items-center justify-between px-4 py-3 bg-surface-primary border-b border-border-secondary cursor-pointer select-none"
			onclick={handleToggleMinimize}
			onkeydown={(e) => e.key === 'Enter' && handleToggleMinimize()}
			role="button"
			tabindex="0"
		>
			<span class="text-sm font-semibold text-text-primary">{headerText}</span>
			<div class="flex items-center gap-1">
				<!-- Minimize/Expand button -->
				<button
					type="button"
					class="p-1 rounded text-text-muted hover:text-text-primary hover:bg-surface-elevated transition-colors"
					onclick={(e) => { e.stopPropagation(); handleToggleMinimize(); }}
					aria-label={isMinimized ? 'Expand' : 'Minimize'}
				>
					{#if isMinimized}
						<ChevronUp size={16} />
					{:else}
						<ChevronDown size={16} />
					{/if}
				</button>

				<!-- Close button -->
				<button
					type="button"
					class="p-1 rounded text-text-muted hover:text-text-primary hover:bg-surface-elevated transition-colors"
					onclick={(e) => { e.stopPropagation(); handleClose(); }}
					aria-label="Close upload panel"
				>
					<X size={16} />
				</button>
			</div>
		</div>

		<!-- Content (collapsible) -->
		{#if !isMinimized}
			<div transition:slide={{ duration: 200 }}>
				<UploadProgress
					uploads={uploadStore.uploads}
					onCancel={handleCancel}
					onRemove={handleRemove}
					onClearCompleted={handleClearCompleted}
				/>
			</div>
		{/if}
	</div>
{/if}
