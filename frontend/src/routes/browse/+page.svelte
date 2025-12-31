<script lang="ts">
	/**
	 * Browse page - main file browser interface (FilePilot style)
	 * Requirements: 1.1, 1.2
	 */
	import { createQuery } from '@tanstack/svelte-query';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Toolbar from '$lib/components/Toolbar.svelte';
	import FileList from '$lib/components/FileList.svelte';
	import StatusBar from '$lib/components/StatusBar.svelte';
	import DriveCard from '$lib/components/DriveCard.svelte';
	import FilePreview from '$lib/components/FilePreview.svelte';
	import {
		pathStore,
		currentPath,
		pathSegments,
		listOptionsStore,
		fileQueryKeys
	} from '$lib/stores/files';
	import { listRoots, listDirectory, search, getDriveStats } from '$lib/api/files';
	import { canPreview } from '$lib/utils/fileTypes';
	import type { SortField, SortDir } from '$lib/types/files';
	import type {
		FileInfo,
		FileList as FileListType,
		RootsResponse,
		SearchResponse,
		DriveStatsResponse
	} from '$lib/api/files';

	let searchQuery = $state('');
	let selectedPaths = $state(new Set<string>());
	let viewMode = $state<'list' | 'grid'>('list');

	// File preview state
	let previewFile = $state<FileInfo | null>(null);

	// Navigation history for back/forward
	let historyStack = $state<string[]>(['']);
	let historyIndex = $state(0);

	// Get current path and options from stores
	const path = $derived($currentPath);
	const segments = $derived($pathSegments);
	const options = $derived($listOptionsStore);

	// Query for mount points (roots)
	const rootsQuery = createQuery<RootsResponse>(() => ({
		queryKey: fileQueryKeys.roots(),
		queryFn: () => listRoots()
	}));

	// Query for drive stats
	const driveStatsQuery = createQuery<DriveStatsResponse>(() => ({
		queryKey: ['files', 'stats'],
		queryFn: () => getDriveStats(),
		enabled: path === ''
	}));

	// Query for directory contents
	const directoryQuery = createQuery<FileListType>(() => ({
		queryKey: fileQueryKeys.list(path, options),
		queryFn: () => listDirectory(path, options),
		enabled: path !== ''
	}));

	// Query for search results
	const searchQueryResult = createQuery<SearchResponse>(() => ({
		queryKey: fileQueryKeys.search(path, searchQuery),
		queryFn: () => search(path, searchQuery),
		enabled: searchQuery.length >= 2
	}));

	// Derived state
	const isLoading = $derived(directoryQuery.isLoading);
	const fileList = $derived(directoryQuery.data ?? null);
	const searchResults = $derived(searchQueryResult.data?.results ?? []);
	const roots = $derived(rootsQuery.data?.roots ?? []);
	const driveStats = $derived(driveStatsQuery.data?.drives ?? []);

	// Check if we're at root (This Server view)
	const isAtRoot = $derived(path === '');

	// Display items for file list (when not at root)
	const displayItems = $derived.by(() => {
		if (searchQuery && searchResults.length > 0) {
			return searchResults;
		}
		return fileList?.items ?? [];
	});

	const itemCount = $derived(isAtRoot ? driveStats.length : displayItems.length);
	const selectedCount = $derived(selectedPaths.size);

	// Navigation helpers
	const canGoBack = $derived(historyIndex > 0);
	const canGoForward = $derived(historyIndex < historyStack.length - 1);
	const canGoUp = $derived(segments.length > 0);

	function handleNavigate(newPath: string) {
		// Add to history
		const newHistory = historyStack.slice(0, historyIndex + 1);
		newHistory.push(newPath);
		historyStack = newHistory;
		historyIndex = newHistory.length - 1;

		pathStore.navigateTo(newPath);
		searchQuery = '';
		selectedPaths = new Set();
	}

	function handleBack() {
		if (canGoBack) {
			historyIndex--;
			pathStore.navigateTo(historyStack[historyIndex]);
			selectedPaths = new Set();
		}
	}

	function handleForward() {
		if (canGoForward) {
			historyIndex++;
			pathStore.navigateTo(historyStack[historyIndex]);
			selectedPaths = new Set();
		}
	}

	function handleUp() {
		if (canGoUp) {
			const parentPath = segments.slice(0, -1).join('/');
			handleNavigate(parentPath);
		}
	}

	function handleRefresh() {
		if (isAtRoot) {
			driveStatsQuery.refetch();
		} else {
			directoryQuery.refetch();
		}
	}

	function handleFileClick(file: FileInfo) {
		if (file.isDir) {
			handleNavigate(file.path);
		} else {
			// Open file preview
			previewFile = file;
		}
	}

	function handleClosePreview() {
		previewFile = null;
	}

	function handleSortChange(field: SortField, dir: SortDir) {
		listOptionsStore.setSortBy(field);
		listOptionsStore.setSortDir(dir);
	}

	function handleSelectionChange(paths: Set<string>) {
		selectedPaths = paths;
	}

	function handleViewModeChange(mode: 'list' | 'grid') {
		viewMode = mode;
	}
</script>

<svelte:head>
	<title>File Manager</title>
</svelte:head>

<div class="file-manager">
	<!-- Sidebar -->
	<Sidebar
		{roots}
		currentPath={path}
		onNavigate={handleNavigate}
	/>

	<!-- Main content area -->
	<div class="main-area">
		<!-- Toolbar with navigation and path bar -->
		<Toolbar
			pathSegments={segments}
			{canGoBack}
			{canGoForward}
			{canGoUp}
			onBack={handleBack}
			onForward={handleForward}
			onUp={handleUp}
			onNavigate={handleNavigate}
			onRefresh={handleRefresh}
		/>

		<!-- File list or Drive cards -->
		<div class="content-area">
			{#if isAtRoot}
				<!-- This Server view - show drive cards -->
				<div class="drives-container">
					<h2 class="drives-title">Storage Devices</h2>
					{#if driveStatsQuery.isLoading}
						<div class="loading-message">Loading drives...</div>
					{:else if driveStats.length === 0}
						<div class="empty-message">No storage devices configured</div>
					{:else}
						<div class="drives-grid">
							{#each driveStats as drive (drive.name)}
								<DriveCard 
									{drive} 
									onClick={() => handleNavigate(drive.name)} 
								/>
							{/each}
						</div>
					{/if}
				</div>
			{:else}
				<FileList
					items={displayItems}
					sortBy={options.sortBy}
					sortDir={options.sortDir}
					{selectedPaths}
					{isLoading}
					onItemClick={handleFileClick}
					onSortChange={handleSortChange}
					onSelectionChange={handleSelectionChange}
				/>
			{/if}
		</div>

		<!-- Status bar -->
		<StatusBar
			{itemCount}
			{selectedCount}
			{viewMode}
			onViewModeChange={handleViewModeChange}
		/>
	</div>
</div>

<!-- File Preview Modal -->
<FilePreview file={previewFile} onClose={handleClosePreview} />

<style>
	.file-manager {
		display: flex;
		height: 100vh;
		width: 100%;
		background: #1e1e1e;
		overflow: hidden;
	}

	.main-area {
		flex: 1;
		display: flex;
		flex-direction: column;
		min-width: 0;
	}

	.content-area {
		flex: 1;
		overflow: auto;
	}

	.drives-container {
		padding: 24px;
	}

	.drives-title {
		font-size: 18px;
		font-weight: 500;
		color: #e0e0e0;
		margin: 0 0 20px 0;
	}

	.drives-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: 16px;
	}

	.loading-message,
	.empty-message {
		color: #888;
		font-size: 14px;
		padding: 20px 0;
	}
</style>
