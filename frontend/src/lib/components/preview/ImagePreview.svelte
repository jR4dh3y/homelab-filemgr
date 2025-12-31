<script lang="ts">
	/**
	 * ImagePreview - Image viewer with zoom support
	 */
	import { ZoomIn, ZoomOut, RotateCw } from 'lucide-svelte';

	interface Props {
		url: string;
		filename: string;
	}

	let { url, filename }: Props = $props();

	let scale = $state(1);
	let rotation = $state(0);
	let error = $state<string | null>(null);
	let loading = $state(true);

	function zoomIn() {
		scale = Math.min(scale + 0.25, 5);
	}

	function zoomOut() {
		scale = Math.max(scale - 0.25, 0.25);
	}

	function rotate() {
		rotation = (rotation + 90) % 360;
	}

	function resetView() {
		scale = 1;
		rotation = 0;
	}

	function handleLoad() {
		loading = false;
	}

	function handleError() {
		loading = false;
		error = 'Failed to load image.';
	}

	function handleWheel(event: WheelEvent) {
		event.preventDefault();
		if (event.deltaY < 0) {
			zoomIn();
		} else {
			zoomOut();
		}
	}
</script>

<div class="image-preview">
	<!-- Controls -->
	<div class="image-controls">
		<button type="button" class="control-btn" onclick={zoomOut} title="Zoom out">
			<ZoomOut size={18} />
		</button>
		<span class="zoom-level">{Math.round(scale * 100)}%</span>
		<button type="button" class="control-btn" onclick={zoomIn} title="Zoom in">
			<ZoomIn size={18} />
		</button>
		<button type="button" class="control-btn" onclick={rotate} title="Rotate">
			<RotateCw size={18} />
		</button>
		<button type="button" class="control-btn reset-btn" onclick={resetView}>
			Reset
		</button>
	</div>

	<!-- Image container -->
	<div class="image-container" onwheel={handleWheel}>
		{#if loading}
			<div class="loading">Loading...</div>
		{/if}
		{#if error}
			<div class="error-message">{error}</div>
		{:else}
			<img
				src={url}
				alt={filename}
				style="transform: scale({scale}) rotate({rotation}deg);"
				class:hidden={loading}
				onload={handleLoad}
				onerror={handleError}
				draggable="false"
			/>
		{/if}
	</div>
</div>

<style>
	.image-preview {
		display: flex;
		flex-direction: column;
		width: 100%;
		height: 100%;
	}

	.image-controls {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		padding: 8px;
		background: #252525;
		border-bottom: 1px solid #333;
		flex-shrink: 0;
	}

	.control-btn {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #333;
		border: none;
		border-radius: 4px;
		color: #ccc;
		cursor: pointer;
		transition: all 0.1s ease;
	}

	.control-btn:hover {
		background: #444;
		color: #fff;
	}

	.reset-btn {
		width: auto;
		padding: 0 12px;
		font-size: 12px;
	}

	.zoom-level {
		font-size: 12px;
		color: #888;
		min-width: 50px;
		text-align: center;
	}

	.image-container {
		flex: 1;
		overflow: auto;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #1a1a1a;
	}

	img {
		max-width: 100%;
		max-height: 100%;
		object-fit: contain;
		transition: transform 0.15s ease;
		user-select: none;
	}

	img.hidden {
		opacity: 0;
	}

	.loading {
		color: #888;
		font-size: 14px;
	}

	.error-message {
		color: #f87171;
		font-size: 14px;
	}
</style>
