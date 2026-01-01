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

<div class="w-full h-full flex items-center justify-center bg-black">
	{#if error}
		<div class="text-danger text-sm text-center p-5">{error}</div>
	{:else}
		<video
			bind:this={videoElement}
			src={url}
			controls
			autoplay
			preload="metadata"
			onerror={handleError}
			class="max-w-full max-h-full outline-none"
		>
			<track kind="captions" />
			Your browser does not support the video tag.
		</video>
	{/if}
</div>
