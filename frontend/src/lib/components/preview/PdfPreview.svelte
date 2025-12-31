<script lang="ts">
	/**
	 * PdfPreview - PDF viewer using iframe embed
	 */

	interface Props {
		url: string;
		filename: string;
	}

	let { url, filename }: Props = $props();

	let error = $state<string | null>(null);

	function handleError() {
		error = 'Failed to load PDF. Try downloading the file instead.';
	}
</script>

<div class="pdf-preview">
	{#if error}
		<div class="error-message">{error}</div>
	{:else}
		<iframe
			src={url}
			title={filename}
			onerror={handleError}
		></iframe>
	{/if}
</div>

<style>
	.pdf-preview {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #2a2a2a;
	}

	iframe {
		width: 100%;
		height: 100%;
		border: none;
		background: #fff;
	}

	.error-message {
		color: #f87171;
		font-size: 14px;
		text-align: center;
		padding: 20px;
	}
</style>
