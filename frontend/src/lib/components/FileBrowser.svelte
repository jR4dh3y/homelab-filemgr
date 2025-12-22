<script lang="ts">
	/**
	 * FileBrowser component - main container for file browsing
	 * Composes Breadcrumb, SearchBar, FileList, and UploadDropzone
	 * Requirements: 1.2, 1.4
	 */

	import type { FileInfo, FileList as FileListType } from '$lib/api/files';
	import type { SortField, SortDir } from '$lib/types/files';
	import Breadcrumb from './Breadcrumb.svelte';
	import SearchBar from './SearchBar.svelte';
	import FileList from './FileList.svelte';

	interface Props {
		/** Current path */
		currentPath?: string;
		/** Path segments for breadcrumb */
		pathSegments?: string[];
		/** File list data */
		fileList?: FileListType | null;
		/** Whether files are loading */
		isLoading?: boolean;
		/** Search query */
		searchQuery?: string;
		/** Whether search is in progress */
		isSearching?: boolean;
		/** Search results */
		searchResults?: FileInfo[];
		/** Current sort field */
		sortBy?: SortField;
		/** Current sort direction */
		sortDir?: SortDir;
		/** Selected file paths */
		selectedPaths?: Set<string>;
		/** Callback when navigating to a path */
		onNavigate?: (path: string) => void;
		/** Callback when a file is clicked */
		onFileClick?: (file: FileInfo) => void;
		/** Callback when search is submitted */
		onSearch?: (query: string) => void;
		/** Callback when search is cleared */
		onSearchClear?: () => void;
		/** Callback when sort changes */
		onSortChange?: (field: SortField, dir: SortDir) => void;
		/** Callback when selection changes */
		onSelectionChange?: (paths: Set<string>) => void;
		/** Callback when files are dropped for upload */
		onFilesDropped?: (files: File[]) => void;
	}

	let {
		currentPath = '',
		pathSegments = [],
		fileList = null,
		isLoading = false,
		searchQuery = '',
		isSearching = false,
		searchResults = [],
		sortBy = 'name',
		sortDir = 'asc',
		selectedPaths = new Set<string>(),
		onNavigate,
		onFileClick,
		onSearch,
		onSearchClear,
		onSortChange,
		onSelectionChange,
		onFilesDropped
	}: Props = $props();

	let isDragOver = $state(false);

	// Determine which items to display (search results or file list)
	const displayItems = $derived(
		searchQuery && searchResults.length > 0 
			? searchResults 
			: (fileList?.items ?? [])
	);

	const isShowingSearchResults = $derived(searchQuery.length > 0 && searchResults.length > 0);

	/**
	 * Handle file/folder click - navigate to folders, trigger callback for files
	 */
	function handleItemClick(item: FileInfo) {
		if (item.isDir) {
			onNavigate?.(item.path);
		} else {
			onFileClick?.(item);
		}
	}

	/**
	 * Handle drag over event
	 */
	function handleDragOver(event: DragEvent) {
		event.preventDefault();
		isDragOver = true;
	}

	/**
	 * Handle drag leave event
	 */
	function handleDragLeave(event: DragEvent) {
		event.preventDefault();
		isDragOver = false;
	}

	/**
	 * Handle drop event
	 */
	function handleDrop(event: DragEvent) {
		event.preventDefault();
		isDragOver = false;
		
		const files = event.dataTransfer?.files;
		if (files && files.length > 0) {
			onFilesDropped?.(Array.from(files));
		}
	}
</script>

<div 
	class="file-browser"
	class:drag-over={isDragOver}
	ondragover={handleDragOver}
	ondragleave={handleDragLeave}
	ondrop={handleDrop}
	role="application"
	aria-label="File browser"
>
	<!-- Header with breadcrumb and search -->
	<header class="browser-header">
		<div class="breadcrumb-container">
			<Breadcrumb segments={pathSegments} onNavigate={onNavigate} />
		</div>
		<div class="search-container">
			<SearchBar 
				value={searchQuery}
				isLoading={isSearching}
				onSearch={onSearch}
				onClear={onSearchClear}
				placeholder="Search in current folder..."
			/>
		</div>
	</header>

	<!-- Status bar -->
	{#if fileList || isShowingSearchResults}
		<div class="status-bar">
			{#if isShowingSearchResults}
				<span class="status-text">
					{searchResults.length} result{searchResults.length !== 1 ? 's' : ''} for "{searchQuery}"
				</span>
				<button 
					type="button" 
					class="clear-search-btn"
					onclick={onSearchClear}
				>
					Clear search
				</button>
			{:else if fileList}
				<span class="status-text">
					{fileList.totalCount} item{fileList.totalCount !== 1 ? 's' : ''}
				</span>
				{#if fileList.totalCount > fileList.pageSize}
					<span class="page-info">
						Page {fileList.page} of {Math.ceil(fileList.totalCount / fileList.pageSize)}
					</span>
				{/if}
			{/if}
		</div>
	{/if}

	<!-- File list -->
	<main class="browser-content">
		<FileList 
			items={displayItems}
			{sortBy}
			{sortDir}
			{selectedPaths}
			isLoading={isLoading || isSearching}
			onItemClick={handleItemClick}
			onSortChange={onSortChange}
			onSelectionChange={onSelectionChange}
		/>
	</main>

	<!-- Drag overlay -->
	{#if isDragOver}
		<div class="drag-overlay">
			<div class="drag-content">
				<span class="drag-icon">📤</span>
				<span class="drag-text">Drop files to upload</span>
			</div>
		</div>
	{/if}
</div>

<style>
	.file-browser {
		display: flex;
		flex-direction: column;
		height: 100%;
		position: relative;
		background: white;
		border-radius: 0.5rem;
		border: 1px solid #e5e7eb;
		overflow: hidden;
	}

	.file-browser.drag-over {
		border-color: #3b82f6;
		border-style: dashed;
	}

	.browser-header {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		padding: 1rem;
		border-bottom: 1px solid #e5e7eb;
		background: #f9fafb;
	}

	.breadcrumb-container {
		flex: 1;
		min-width: 200px;
	}

	.search-container {
		flex-shrink: 0;
		width: 100%;
		max-width: 320px;
	}

	@media (max-width: 640px) {
		.browser-header {
			flex-direction: column;
			align-items: stretch;
		}
		.search-container {
			max-width: none;
		}
	}

	.status-bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.5rem 1rem;
		background: #f3f4f6;
		border-bottom: 1px solid #e5e7eb;
		font-size: 0.75rem;
		color: #6b7280;
	}

	.status-text {
		font-weight: 500;
	}

	.page-info {
		color: #9ca3af;
	}

	.clear-search-btn {
		padding: 0.25rem 0.5rem;
		font-size: 0.75rem;
		color: #3b82f6;
		background: transparent;
		border: none;
		cursor: pointer;
		border-radius: 0.25rem;
		transition: background-color 0.15s;
	}

	.clear-search-btn:hover {
		background: #dbeafe;
	}

	.browser-content {
		flex: 1;
		overflow: auto;
	}

	.drag-overlay {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(59, 130, 246, 0.1);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 20;
		pointer-events: none;
	}

	.drag-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		padding: 2rem;
		background: white;
		border-radius: 0.75rem;
		border: 2px dashed #3b82f6;
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
	}

	.drag-icon {
		font-size: 3rem;
	}

	.drag-text {
		font-size: 1rem;
		font-weight: 500;
		color: #3b82f6;
	}

	/* Dark mode */
	@media (prefers-color-scheme: dark) {
		.file-browser {
			background: #111827;
			border-color: #374151;
		}

		.file-browser.drag-over {
			border-color: #60a5fa;
		}

		.browser-header {
			background: #1f2937;
			border-bottom-color: #374151;
		}

		.status-bar {
			background: #1f2937;
			border-bottom-color: #374151;
			color: #9ca3af;
		}

		.status-text {
			color: #d1d5db;
		}

		.page-info {
			color: #6b7280;
		}

		.clear-search-btn {
			color: #60a5fa;
		}

		.clear-search-btn:hover {
			background: #1e3a5f;
		}

		.drag-overlay {
			background: rgba(59, 130, 246, 0.15);
		}

		.drag-content {
			background: #1f2937;
			border-color: #60a5fa;
		}

		.drag-text {
			color: #60a5fa;
		}
	}
</style>
