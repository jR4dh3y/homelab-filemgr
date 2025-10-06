<script lang="ts">
	import type { FileInfo } from '$lib/api/files';
	import FileList from './FileList.svelte';

	interface Props {
		currentPath?: string;
		files?: FileInfo[];
		isLoading?: boolean;
		onNavigate?: (path: string) => void;
	}

	let {
		currentPath = '',
		files = [],
		isLoading = false,
		onNavigate
	}: Props = $props();

	function handleItemClick(item: FileInfo) {
		if (item.isDir) {
			onNavigate?.(item.path);
		}
	}
</script>

<div class="file-browser">
	<header class="browser-header">
		<span class="current-path">{currentPath || '/'}</span>
	</header>

	<main class="browser-content">
		{#if isLoading}
			<p>Loading...</p>
		{:else}
			<FileList items={files} onItemClick={handleItemClick} />
		{/if}
	</main>
</div>

<style>
	.file-browser {
		display: flex;
		flex-direction: column;
		height: 100%;
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
	}

	.browser-header {
		padding: 1rem;
		border-bottom: 1px solid #e5e7eb;
		background: #f9fafb;
	}

	.current-path {
		font-family: monospace;
		color: #374151;
	}

	.browser-content {
		flex: 1;
		overflow: auto;
		padding: 1rem;
	}
</style>