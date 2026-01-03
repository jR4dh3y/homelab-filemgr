<script lang="ts">
	/**
	 * Browse page - main file browser interface (FilePilot style)
	 */
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { goto } from '$app/navigation';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Toolbar from '$lib/components/Toolbar.svelte';
	import FileList from '$lib/components/FileList.svelte';
	import StatusBar from '$lib/components/StatusBar.svelte';
	import DriveCard from '$lib/components/DriveCard.svelte';
	import FilePreview from '$lib/components/FilePreview.svelte';
	import { Spinner, Modal, Input, Button } from '$lib/components/ui';
	import { pathStore, currentPath, pathSegments, listOptionsStore, fileQueryKeys } from '$lib/stores/files';
	import { settingsStore } from '$lib/stores/settings';
	import { clipboardStore } from '$lib/stores/clipboard.svelte';
	import { listRoots, listDirectory, search, getDriveStats, rename, deleteFile, getDownloadUrl } from '$lib/api/files';
	import { createCopyJob, createMoveJob, createDeleteJob } from '$lib/api/jobs';
	import { formatFileSize, formatFileDate } from '$lib/utils/format';
	import type { SortField, SortDir } from '$lib/types/files';
	import type { FileInfo, FileList as FileListType, RootsResponse, SearchResponse, DriveStatsResponse } from '$lib/api/files';
	import { SvelteSet } from 'svelte/reactivity';

	const queryClient = useQueryClient();

	let searchQuery = $state('');
	let selectedPaths = $state(new Set<string>());
	let viewMode = $state<'list' | 'grid'>('list');
	let previewFile = $state<FileInfo | null>(null);
	let historyStack = $state<string[]>(['']);
	let historyIndex = $state(0);

	// Rename dialog state
	let renameDialog = $state<{ open: boolean; file: FileInfo | null; newName: string }>({
		open: false,
		file: null,
		newName: '',
	});

	// Delete confirmation dialog state
	let deleteDialog = $state<{ open: boolean; items: FileInfo[] }>({
		open: false,
		items: [],
	});

	// Properties dialog state
	let propertiesDialog = $state<{ open: boolean; file: FileInfo | null }>({
		open: false,
		file: null,
	});

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

	// Clipboard state for context menu
	const canPaste = $derived(clipboardStore.hasItems);
	const cutPaths = $derived.by(() => {
		if (clipboardStore.operation === 'cut') {
			return new SvelteSet(clipboardStore.items.map((i) => i.path));
		}
		return new SvelteSet<string>();
	});

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

	/**
	 * Handle context menu actions
	 */
	async function handleContextMenuAction(action: string, items: FileInfo[]) {
		switch (action) {
			case 'copy':
				clipboardStore.copy(items);
				break;

			case 'cut':
				clipboardStore.cut(items);
				break;

			case 'paste':
				await handlePaste();
				break;

			case 'rename':
				if (items.length === 1) {
					renameDialog = {
						open: true,
						file: items[0],
						newName: items[0].name,
					};
				}
				break;

			case 'delete':
				deleteDialog = {
					open: true,
					items: items,
				};
				break;

			case 'download':
				handleDownload(items);
				break;

			case 'properties':
				if (items.length === 1) {
					propertiesDialog = {
						open: true,
						file: items[0],
					};
				}
				break;
		}
	}

	/**
	 * Handle paste operation
	 */
	async function handlePaste() {
		if (!clipboardStore.hasItems || !path) return;

		const operation = clipboardStore.operation;
		const items = clipboardStore.items;

		try {
			for (const item of items) {
				const destPath = `${path}/${item.name}`;
				if (operation === 'copy') {
					await createCopyJob(item.path, destPath);
				} else if (operation === 'cut') {
					await createMoveJob(item.path, destPath);
				}
			}

			// Clear clipboard after cut operation
			if (operation === 'cut') {
				clipboardStore.clear();
			}

			// Refresh directory listing
			directoryQuery.refetch();
		} catch (error) {
			console.error('Paste operation failed:', error);
		}
	}

	/**
	 * Handle file download
	 */
	function handleDownload(items: FileInfo[]) {
		for (const item of items) {
			if (!item.isDir) {
				const downloadUrl = getDownloadUrl(item.path);
				window.open(downloadUrl, '_blank');
			}
		}
	}

	/**
	 * Handle rename confirmation
	 */
	async function handleRenameConfirm() {
		if (!renameDialog.file || !renameDialog.newName.trim()) return;

		const oldPath = renameDialog.file.path;
		const parentPath = oldPath.substring(0, oldPath.lastIndexOf('/'));
		const newPath = parentPath ? `${parentPath}/${renameDialog.newName}` : renameDialog.newName;

		try {
			await rename(oldPath, newPath);
			renameDialog = { open: false, file: null, newName: '' };
			directoryQuery.refetch();
		} catch (error) {
			console.error('Rename failed:', error);
		}
	}

	/**
	 * Handle delete confirmation
	 */
	async function handleDeleteConfirm() {
		if (deleteDialog.items.length === 0) return;

		try {
			for (const item of deleteDialog.items) {
				if (item.isDir) {
					// Use job for directory deletion
					await createDeleteJob(item.path);
				} else {
					await deleteFile(item.path);
				}
			}
			deleteDialog = { open: false, items: [] };
			selectedPaths = new Set();
			directoryQuery.refetch();
		} catch (error) {
			console.error('Delete failed:', error);
		}
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
					{cutPaths}
					{canPaste}
					onItemClick={handleFileClick}
					onSortChange={handleSortChange}
					onSelectionChange={handleSelectionChange}
					onContextMenuAction={handleContextMenuAction}
				/>
			{/if}
		</div>

		<!-- Status bar -->
		<StatusBar {itemCount} {selectedCount} {viewMode} onViewModeChange={handleViewModeChange} />
	</div>
</div>

<!-- File Preview Modal -->
<FilePreview file={previewFile} onClose={handleClosePreview} />

<!-- Rename Dialog -->
<Modal
	open={renameDialog.open}
	title="Rename"
	onclose={() => (renameDialog = { open: false, file: null, newName: '' })}
>
	<div class="flex flex-col gap-4">
		<p class="text-text-secondary text-sm">Enter a new name:</p>
		<Input
			bind:value={renameDialog.newName}
			placeholder="New name"
			onkeydown={(e) => e.key === 'Enter' && handleRenameConfirm()}
		/>
	</div>
	{#snippet footer()}
		<Button variant="secondary" onclick={() => (renameDialog = { open: false, file: null, newName: '' })}>
			Cancel
		</Button>
		<Button variant="primary" onclick={handleRenameConfirm}>
			Rename
		</Button>
	{/snippet}
</Modal>

<!-- Delete Confirmation Dialog -->
<Modal
	open={deleteDialog.open}
	title="Delete"
	onclose={() => (deleteDialog = { open: false, items: [] })}
>
	<div class="flex flex-col gap-2">
		<p class="text-text-primary">
			Are you sure you want to delete {deleteDialog.items.length === 1
				? `"${deleteDialog.items[0]?.name}"`
				: `${deleteDialog.items.length} items`}?
		</p>
		<p class="text-text-secondary text-sm">This action cannot be undone.</p>
	</div>
	{#snippet footer()}
		<Button variant="secondary" onclick={() => (deleteDialog = { open: false, items: [] })}>
			Cancel
		</Button>
		<Button variant="danger" onclick={handleDeleteConfirm}>
			Delete
		</Button>
	{/snippet}
</Modal>

<!-- Properties Dialog -->
<Modal
	open={propertiesDialog.open}
	title="Properties"
	onclose={() => (propertiesDialog = { open: false, file: null })}
>
	{#if propertiesDialog.file}
		{@const file = propertiesDialog.file}
		<div class="flex flex-col gap-3 text-sm">
			<div class="flex justify-between">
				<span class="text-text-secondary">Name:</span>
				<span class="text-text-primary font-medium">{file.name}</span>
			</div>
			<div class="flex justify-between">
				<span class="text-text-secondary">Type:</span>
				<span class="text-text-primary">{file.isDir ? 'Folder' : file.mimeType || 'File'}</span>
			</div>
			<div class="flex justify-between">
				<span class="text-text-secondary">Path:</span>
				<span class="text-text-primary break-all">{file.path}</span>
			</div>
			{#if !file.isDir}
				<div class="flex justify-between">
					<span class="text-text-secondary">Size:</span>
					<span class="text-text-primary">{formatFileSize(file.size)}</span>
				</div>
			{/if}
			<div class="flex justify-between">
				<span class="text-text-secondary">Modified:</span>
				<span class="text-text-primary">{formatFileDate(file.modTime)}</span>
			</div>
			<div class="flex justify-between">
				<span class="text-text-secondary">Permissions:</span>
				<span class="text-text-primary font-mono">{file.permissions}</span>
			</div>
		</div>
	{/if}
	{#snippet footer()}
		<Button variant="secondary" onclick={() => (propertiesDialog = { open: false, file: null })}>
			Close
		</Button>
	{/snippet}
</Modal>
