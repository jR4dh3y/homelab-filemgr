<script lang="ts">
	/**
	 * Browse page - main file browser interface
	 * Requirements: 1.1, 1.2
	 */
	import { createQuery } from '@tanstack/svelte-query';
	import FileBrowser from '$lib/components/FileBrowser.svelte';
	import {
		pathStore,
		currentPath,
		pathSegments,
		listOptionsStore,
		fileQueryKeys
	} from '$lib/stores/files';
	import { listRoots, listDirectory, search } from '$lib/api/files';
	import type { SortField, SortDir } from '$lib/types/files';
	import type {
		FileInfo,
		FileList as FileListType,
		RootsResponse,
		SearchResponse
	} from '$lib/api/files';

	let searchQuery = $state('');
	let selectedPaths = $state(new Set<string>());

	// Get current path and options from stores
	const path = $derived($currentPath);
	const segments = $derived($pathSegments);
	const options = $derived($listOptionsStore);

	// Query for mount points (roots) - wrap in accessor function for Svelte 5
	const rootsQuery = createQuery<RootsResponse>(() => ({
		queryKey: fileQueryKeys.roots(),
		queryFn: () => listRoots()
	}));

	// Query for directory contents - wrap in accessor function for Svelte 5
	const directoryQuery = createQuery<FileListType>(() => ({
		queryKey: fileQueryKeys.list(path, options),
		queryFn: () => listDirectory(path, options),
		enabled: path !== ''
	}));

	// Query for search results - wrap in accessor function for Svelte 5
	const searchQueryResult = createQuery<SearchResponse>(() => ({
		queryKey: fileQueryKeys.search(path, searchQuery),
		queryFn: () => search(path, searchQuery),
		enabled: searchQuery.length >= 2
	}));

	// Derived state - access query results directly (no $ prefix needed in Svelte 5)
	const isLoading = $derived(directoryQuery.isLoading);
	const isSearching = $derived(searchQueryResult.isLoading);
	const fileList = $derived(directoryQuery.data ?? null);
	const searchResults = $derived(searchQueryResult.data?.results ?? []);

	// Show roots when at root path
	const showRoots = $derived(path === '');
	const roots = $derived(rootsQuery.data?.roots ?? []);

	function handleNavigate(newPath: string) {
		pathStore.navigateTo(newPath);
		searchQuery = '';
	}

	function handleFileClick(file: FileInfo) {
		// TODO: Implement file preview/download
		console.log('File clicked:', file);
	}

	function handleSearch(query: string) {
		searchQuery = query;
	}

	function handleSearchClear() {
		searchQuery = '';
	}

	function handleSortChange(field: SortField, dir: SortDir) {
		listOptionsStore.setSortBy(field);
		listOptionsStore.setSortDir(dir);
	}

	function handleSelectionChange(paths: Set<string>) {
		selectedPaths = paths;
	}

	function handleFilesDropped(files: File[]) {
		// TODO: Implement file upload
		console.log('Files dropped:', files);
	}
</script>

<svelte:head>
	<title>Browse Files</title>
</svelte:head>

<div class="browse-page">
	{#if showRoots}
		<!-- Show mount points at root -->
		<div class="roots-container">
			<h1 class="roots-title">Mount Points</h1>
			{#if rootsQuery.isLoading}
				<div class="loading">Loading mount points...</div>
			{:else if rootsQuery.error}
				<div class="error">Failed to load mount points</div>
			{:else if roots.length === 0}
				<div class="empty">No mount points configured</div>
			{:else}
				<div class="roots-grid">
					{#each roots as root (root.name)}
						<button type="button" class="root-card" onclick={() => handleNavigate(root.name)}>
							<span class="root-icon">📁</span>
							<span class="root-name">{root.name}</span>
							{#if root.readOnly}
								<span class="root-badge">Read-only</span>
							{/if}
						</button>
					{/each}
				</div>
			{/if}
		</div>
	{:else}
		<!-- Show file browser -->
		<FileBrowser
			pathSegments={segments}
			{fileList}
			{isLoading}
			{searchQuery}
			{isSearching}
			{searchResults}
			sortBy={options.sortBy}
			sortDir={options.sortDir}
			{selectedPaths}
			onNavigate={handleNavigate}
			onFileClick={handleFileClick}
			onSearch={handleSearch}
			onSearchClear={handleSearchClear}
			onSortChange={handleSortChange}
			onSelectionChange={handleSelectionChange}
			onFilesDropped={handleFilesDropped}
		/>
	{/if}
</div>

<style>
	.browse-page {
		height: 100%;
		padding: 1rem;
	}

	.roots-container {
		max-width: 800px;
		margin: 0 auto;
	}

	.roots-title {
		font-size: 1.5rem;
		font-weight: 600;
		margin-bottom: 1.5rem;
		color: #111827;
	}

	.loading,
	.error,
	.empty {
		padding: 2rem;
		text-align: center;
		color: #6b7280;
	}

	.error {
		color: #dc2626;
	}

	.roots-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 1rem;
	}

	.root-card {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		padding: 1.5rem;
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		cursor: pointer;
		transition: all 0.15s;
	}

	.root-card:hover {
		border-color: #3b82f6;
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
	}

	.root-icon {
		font-size: 2.5rem;
	}

	.root-name {
		font-weight: 500;
		color: #374151;
	}

	.root-badge {
		font-size: 0.75rem;
		padding: 0.125rem 0.5rem;
		background: #fef3c7;
		color: #92400e;
		border-radius: 9999px;
	}

	/* Dark mode */
	@media (prefers-color-scheme: dark) {
		.roots-title {
			color: #f9fafb;
		}

		.loading,
		.empty {
			color: #9ca3af;
		}

		.root-card {
			background: #1f2937;
			border-color: #374151;
		}

		.root-card:hover {
			border-color: #60a5fa;
		}

		.root-name {
			color: #e5e7eb;
		}

		.root-badge {
			background: #78350f;
			color: #fde68a;
		}
	}
</style>
