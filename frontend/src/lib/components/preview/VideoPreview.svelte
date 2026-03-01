<script lang="ts">
	/**
	 * VideoPreview - HTML5 video player with streaming support
	 */
	import { Button } from '$lib/components/ui';
	import { Download } from 'lucide-svelte';

	interface Props {
		url: string;
		filename: string;
		downloadUrl?: string;
	}

	let { url, filename, downloadUrl }: Props = $props();

	let videoElement: HTMLVideoElement | null = $state(null);
	let error = $state<string | null>(null);

	function handleError() {
		const mediaError = videoElement?.error;
		if (mediaError?.code === MediaError.MEDIA_ERR_DECODE) {
			error = 'This video file uses a codec your browser cannot decode. Download it or play it in a native media player.';
			return;
		}
		if (mediaError?.code === MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED) {
			error = 'This video format is not supported by your browser. Download it or play it in a native media player.';
			return;
		}
		error = 'Failed to load video. The file may be corrupted or the codec is not supported by your browser.';
	}

	function openDownload() {
		if (downloadUrl) {
			window.open(downloadUrl, '_blank');
		}
	}
</script>

<div class="w-full h-full flex items-center justify-center bg-black">
	{#if error}
		<div class="flex flex-col items-center gap-4 text-danger text-sm text-center p-5 max-w-2xl">
			<p>{error}</p>
			{#if downloadUrl}
				<Button variant="primary" onclick={openDownload}>
					<Download size={18} />
					Download Video
				</Button>
			{/if}
		</div>
	{:else}
		<video
			bind:this={videoElement}
			src={url}
			controls
			autoplay
			preload="metadata"
			playsinline
			onerror={handleError}
			class="max-w-full max-h-full outline-none"
		>
			<track kind="captions" />
			Your browser does not support the video tag.
		</video>
	{/if}
</div>
