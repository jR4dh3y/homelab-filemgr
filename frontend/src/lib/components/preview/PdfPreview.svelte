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

<div class="w-full h-full flex items-center justify-center bg-surface-tertiary">
	{#if error}
		<div class="text-danger text-sm text-center p-5">{error}</div>
	{:else}
		<iframe
			src={url}
			title={filename}
			onerror={handleError}
			class="w-full h-full border-none bg-white"
		></iframe>
	{/if}
</div>
