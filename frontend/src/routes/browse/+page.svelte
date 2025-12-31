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
	let viewMode = $state<'list' | 'grid'>('list');

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

	// Display items - show roots at root path, otherwise show directory contents
	const displayItems = $derived.by(() => {
		if (searchQuery && searchResults.length > 0) {
			return searchResults;
		}
		if (path === '') {
			// At root, show mount points as folders
			return roots.map(root => ({
				name: root.name,
				path: root.name,
				size: 0,
				isDir: true,
				modTime: new Date().toISOString(),
				permissions: root.readOnly ? 'r--' : 'rw-'
			})) as FileInfo[];
		}
		return fileList?.items ?? [];
	});

	const itemCount = $derived(displayItems.length);
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
		directoryQuery.refetch();
	}

	function handleFileClick(file: FileInfo) {
		if (file.isDir) {
			handleNavigate(file.path);
		} else {
			// TODO: Implement file preview/download
			console.log('File clicked:', file);
		}
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

		<!-- File list -->
		<div class="content-area">
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
		overflow: hidden;
	}
</style>
