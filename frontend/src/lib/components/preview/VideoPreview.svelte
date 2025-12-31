<script lang="ts">
	/**
	 * VideoPreview - HTML5 video player with streaming support
	 */

	interface Props {
		url: string;
		filename: string;
	}

	let { url, filename }: Props = $props();

	let videoElement: HTMLVideoElement | null = $state(null);
	let error = $state<string | null>(null);

	function handleError() {
		error = 'Failed to load video. The format may not be supported by your browser.';
	}
</script>

<div class="video-preview">
	{#if error}
		<div class="error-message">{error}</div>
	{:else}
		<video
			bind:this={videoElement}
			src={url}
			controls
			autoplay
			preload="metadata"
			onerror={handleError}
		>
			<track kind="captions" />
			Your browser does not support the video tag.
		</video>
	{/if}
</div>

<style>
	.video-preview {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #000;
	}

	video {
		max-width: 100%;
		max-height: 100%;
		outline: none;
	}

	.error-message {
		color: #f87171;
		font-size: 14px;
		text-align: center;
		padding: 20px;
	}
</style>
