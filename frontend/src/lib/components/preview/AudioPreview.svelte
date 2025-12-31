<script lang="ts">
	/**
	 * AudioPreview - HTML5 audio player with visualization
	 */
	import { Music } from 'lucide-svelte';

	interface Props {
		url: string;
		filename: string;
	}

	let { url, filename }: Props = $props();

	let error = $state<string | null>(null);

	function handleError() {
		error = 'Failed to load audio. The format may not be supported by your browser.';
	}
</script>

<div class="audio-preview">
	<div class="audio-icon">
		<Music size={64} />
	</div>
	<div class="audio-filename">{filename}</div>
	{#if error}
		<div class="error-message">{error}</div>
	{:else}
		<audio
			src={url}
			controls
			autoplay
			preload="metadata"
			onerror={handleError}
		>
			Your browser does not support the audio tag.
		</audio>
	{/if}
</div>

<style>
	.audio-preview {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 20px;
		padding: 40px;
		width: 100%;
		max-width: 500px;
	}

	.audio-icon {
		color: #4a9eff;
		opacity: 0.8;
	}

	.audio-filename {
		font-size: 16px;
		font-weight: 500;
		color: #e0e0e0;
		text-align: center;
		word-break: break-word;
	}

	audio {
		width: 100%;
		outline: none;
	}

	/* Style the audio controls for dark theme */
	audio::-webkit-media-controls-panel {
		background: #333;
	}

	.error-message {
		color: #f87171;
		font-size: 14px;
		text-align: center;
	}
</style>
