<script lang="ts">
	/**
	 * Browse page - main file browser interface (FilePilot style)
	 */
	import { createQuery } from '@tanstack/svelte-query';
	import { goto } from '$app/navigation';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Toolbar from '$lib/components/Toolbar.svelte';
	import FileList from '$lib/components/FileList.svelte';
	import StatusBar from '$lib/components/StatusBar.svelte';
	import DriveCard from '$lib/components/DriveCard.svelte';
	import FilePreview from '$lib/components/FilePreview.svelte';
	import { Spinner } from '$lib/components/ui';
	import { pathStore, currentPath, pathSegments, listOptionsStore, fileQueryKeys } from '$lib/stores/files';
	import { settingsStore } from '$lib/stores/settings';
	import { listRoots, listDirectory, search, getDriveStats } from '$lib/api/files';
	import type { SortField, SortDir } from '$lib/types/files';
	import type { FileInfo, FileList as FileListType, RootsResponse, SearchResponse, DriveStatsResponse } from '$lib/api/files';

	let searchQuery = $state('');
	let selectedPaths = $state(new Set<string>());
	let viewMode = $state<'list' | 'grid'>('list');
	let previewFile = $state<FileInfo | null>(null);
	let historyStack = $state<string[]>(['']);
	let historyIndex = $state(0);

	const path = $derived($currentPath);
	const segments = $derived($pathSegments);
	const options = $derived($listOptionsStore);
	const settings = $derived($settingsStore);

	const rootsQuery = createQuery<RootsResponse>(() => ({
		queryKey: fileQueryKeys.roots(),
		queryFn: () => listRoots(),
	}));

	const driveStatsQuery = createQuery<DriveStatsResponse>(() => ({
		queryKey: ['files', 'stats'],
		queryFn: () => getDriveStats(),
		enabled: path === '',
	}));

	const directoryQuery = createQuery<FileListType>(() => ({
		queryKey: fileQueryKeys.list(path, options),
		queryFn: () => listDirectory(path, options),
		enabled: path !== '',
	}));

	const searchQueryResult = createQuery<SearchResponse>(() => ({
		queryKey: fileQueryKeys.search(path, searchQuery),
		queryFn: () => search(path, searchQuery),
		enabled: searchQuery.length >= 2,
	}));

	const isLoading = $derived(directoryQuery.isLoading);
	const fileList = $derived(directoryQuery.data ?? null);
	const searchResults = $derived(searchQueryResult.data?.results ?? []);
	const roots = $derived(rootsQuery.data?.roots ?? []);
	const driveStats = $derived(driveStatsQuery.data?.drives ?? []);
	const isAtRoot = $derived(path === '');

	const displayItems = $derived.by(() => {
		let items: FileInfo[];

		if (searchQuery && searchResults.length > 0) {
			items = searchResults;
		} else {
			items = fileList?.items ?? [];
		}

		if (!settings.showHiddenFiles) {
			items = items.filter((item) => !item.name.startsWith('.'));
		}

		return items;
	});

	const itemCount = $derived(isAtRoot ? driveStats.length : displayItems.length);
	const selectedCount = $derived(selectedPaths.size);
	const canGoBack = $derived(historyIndex > 0);
	const canGoForward = $derived(historyIndex < historyStack.length - 1);
	const canGoUp = $derived(segments.length > 0);

	function handleNavigate(newPath: string) {
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

	function handleSettings() {
		goto('/settings');
	}

	function handleFileClick(file: FileInfo) {
		if (file.isDir) {
			handleNavigate(file.path);
		} else {
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

<div class="flex h-screen w-full bg-surface-primary overflow-hidden">
	<!-- Sidebar -->
	<Sidebar {roots} currentPath={path} onNavigate={handleNavigate} />

	<!-- Main content area -->
	<div class="flex-1 flex flex-col min-w-0">
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
			onSettings={handleSettings}
		/>

		<!-- File list or Drive cards -->
		<div class="flex-1 overflow-auto">
			{#if isAtRoot}
				<!-- This Server view - show drive cards -->
				<div class="p-6">
					<h2 class="text-lg font-medium text-text-primary m-0 mb-5">Storage Devices</h2>
					{#if driveStatsQuery.isLoading}
						<div class="flex items-center gap-2 text-text-secondary text-sm py-5">
							<Spinner size="sm" />
							<span>Loading drives...</span>
						</div>
					{:else if driveStats.length === 0}
						<div class="text-text-secondary text-sm py-5">No storage devices configured</div>
					{:else}
						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
							{#each driveStats as drive (drive.name)}
								<DriveCard {drive} onClick={() => handleNavigate(drive.name)} />
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
					compactMode={settings.compactMode}
					onItemClick={handleFileClick}
					onSortChange={handleSortChange}
					onSelectionChange={handleSelectionChange}
				/>
			{/if}
		</div>

		<!-- Status bar -->
		<StatusBar {itemCount} {selectedCount} {viewMode} onViewModeChange={handleViewModeChange} />
	</div>
</div>

<!-- File Preview Modal -->
<FilePreview file={previewFile} onClose={handleClosePreview} />
