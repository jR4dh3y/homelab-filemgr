<script lang="ts">
	/**
	 * FileList component with sortable columns - FilePilot style
	 * Requirements: 1.1, 1.2
	 */

	import type { FileInfo } from '$lib/api/files';
	import type { SortField, SortDir } from '$lib/types/files';
	import { formatFileSize, formatFileDate, getFileTypeDescription } from '$lib/utils/format';
	import { SvelteSet } from 'svelte/reactivity';
	import {
		Folder,
		File,
		FileImage,
		FileVideo,
		FileAudio,
		FileText,
		FileCode,
		FileArchive,
		FileSpreadsheet,
		Globe,
		Palette,
		FileJson
	} from 'lucide-svelte';

	let {
		items = [],
		sortBy = 'name',
		sortDir = 'asc',
		selectedPaths = new SvelteSet<string>(),
		isLoading = false,
		compactMode = false,
		onItemClick,
		onSortChange,
		onSelectionChange
	}: {
		items?: FileInfo[];
		sortBy?: SortField;
		sortDir?: SortDir;
		selectedPaths?: Set<string>;
		isLoading?: boolean;
		compactMode?: boolean;
		onItemClick?: (item: FileInfo) => void;
		onSortChange?: (field: SortField, dir: SortDir) => void;
		onSelectionChange?: (paths: Set<string>) => void;
	} = $props();

	function handleSort(field: SortField) {
		if (sortBy === field) {
			const newDir = sortDir === 'asc' ? 'desc' : 'asc';
			onSortChange?.(field, newDir);
		} else {
			onSortChange?.(field, 'asc');
		}
	}

	function handleRowClick(item: FileInfo, event: MouseEvent) {
		if (event.ctrlKey || event.metaKey) {
			const newSelection = new SvelteSet<string>(selectedPaths);
			if (newSelection.has(item.path)) {
				newSelection.delete(item.path);
			} else {
				newSelection.add(item.path);
			}
			onSelectionChange?.(newSelection);
		} else if (event.shiftKey && selectedPaths.size > 0) {
			// Shift-click for range selection
			const newSelection = new SvelteSet<string>(selectedPaths);
			newSelection.add(item.path);
			onSelectionChange?.(newSelection);
		} else {
			// Single click selects, double click opens
			const newSelection = new SvelteSet<string>([item.path]);
			onSelectionChange?.(newSelection);
		}
	}

	function handleDoubleClick(item: FileInfo) {
		onItemClick?.(item);
	}

	function handleKeyDown(item: FileInfo, event: KeyboardEvent) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			onItemClick?.(item);
		}
	}

	function getSortIndicator(field: SortField): string {
		if (sortBy !== field) return '';
		return sortDir === 'asc' ? '▲' : '▼';
	}

	function getFileIcon(item: FileInfo) {
		if (item.isDir) return Folder;
		const ext = item.name.includes('.')
			? item.name.slice(item.name.lastIndexOf('.') + 1).toLowerCase()
			: '';
		
		const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'svg', 'webp', 'bmp', 'ico'];
		const videoExts = ['mp4', 'mkv', 'avi', 'mov', 'webm', 'wmv', 'flv'];
		const audioExts = ['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a', 'wma'];
		const codeExts = ['js', 'ts', 'py', 'go', 'rs', 'java', 'c', 'cpp', 'h', 'rb', 'php', 'swift', 'kt'];
		const archiveExts = ['zip', 'rar', '7z', 'tar', 'gz', 'bz2', 'xz'];
		const spreadsheetExts = ['xls', 'xlsx', 'csv', 'ods'];
		const docExts = ['pdf', 'doc', 'docx', 'txt', 'md', 'rtf', 'odt'];
		const webExts = ['html', 'htm', 'xml'];
		const styleExts = ['css', 'scss', 'sass', 'less'];
		const dataExts = ['json', 'yaml', 'yml', 'toml'];

		if (imageExts.includes(ext)) return FileImage;
		if (videoExts.includes(ext)) return FileVideo;
		if (audioExts.includes(ext)) return FileAudio;
		if (codeExts.includes(ext)) return FileCode;
		if (archiveExts.includes(ext)) return FileArchive;
		if (spreadsheetExts.includes(ext)) return FileSpreadsheet;
		if (docExts.includes(ext)) return FileText;
		if (webExts.includes(ext)) return Globe;
		if (styleExts.includes(ext)) return Palette;
		if (dataExts.includes(ext)) return FileJson;
		
		return File;
	}

	function isSelected(path: string): boolean {
		return selectedPaths.has(path);
	}
</script>

<div class="file-list" class:compact={compactMode}>
	{#if isLoading}
		<div class="loading-overlay">
			<div class="spinner"></div>
		</div>
	{/if}

	<table class="file-table" role="grid" aria-busy={isLoading}>
		<thead>
			<tr>
				<th
					class="sortable name-col"
					class:sorted={sortBy === 'name'}
					onclick={() => handleSort('name')}
					onkeydown={(e) => e.key === 'Enter' && handleSort('name')}
					tabindex="0"
					role="columnheader"
					aria-sort={sortBy === 'name' ? (sortDir === 'asc' ? 'ascending' : 'descending') : 'none'}
				>
					<span class="col-label">Name</span>
					<span class="sort-indicator">{getSortIndicator('name')}</span>
				</th>
				<th
					class="sortable type-col"
					class:sorted={sortBy === 'type'}
					onclick={() => handleSort('type')}
					onkeydown={(e) => e.key === 'Enter' && handleSort('type')}
					tabindex="0"
					role="columnheader"
					aria-sort={sortBy === 'type' ? (sortDir === 'asc' ? 'ascending' : 'descending') : 'none'}
				>
					<span class="col-label">Type</span>
					<span class="sort-indicator">{getSortIndicator('type')}</span>
				</th>
				<th
					class="sortable size-col"
					class:sorted={sortBy === 'size'}
					onclick={() => handleSort('size')}
					onkeydown={(e) => e.key === 'Enter' && handleSort('size')}
					tabindex="0"
					role="columnheader"
					aria-sort={sortBy === 'size' ? (sortDir === 'asc' ? 'ascending' : 'descending') : 'none'}
				>
					<span class="col-label">Size</span>
					<span class="sort-indicator">{getSortIndicator('size')}</span>
				</th>
				<th
					class="sortable date-col"
					class:sorted={sortBy === 'modTime'}
					onclick={() => handleSort('modTime')}
					onkeydown={(e) => e.key === 'Enter' && handleSort('modTime')}
					tabindex="0"
					role="columnheader"
					aria-sort={sortBy === 'modTime'
						? sortDir === 'asc'
							? 'ascending'
							: 'descending'
						: 'none'}
				>
					<span class="col-label">Modified</span>
					<span class="sort-indicator">{getSortIndicator('modTime')}</span>
				</th>
			</tr>
		</thead>
		<tbody>
			{#if items.length === 0 && !isLoading}
				<tr class="empty-row">
					<td colspan="4">
						<div class="empty-state">
							<span class="empty-icon">📂</span>
							<span class="empty-text">This folder is empty</span>
						</div>
					</td>
				</tr>
			{:else}
				{#each items as item (item.path)}
					{@const IconComponent = getFileIcon(item)}
					<tr
						class="file-row"
						class:selected={isSelected(item.path)}
						class:directory={item.isDir}
						onclick={(e) => handleRowClick(item, e)}
						onkeydown={(e) => handleKeyDown(item, e)}
						ondblclick={() => handleDoubleClick(item)}
						tabindex="0"
						aria-selected={isSelected(item.path)}
					>
						<td class="name-cell">
							<span class="file-icon" class:folder={item.isDir}>
								<IconComponent size={16} />
							</span>
							<span class="file-name" title={item.name}>{item.name}</span>
						</td>
						<td class="type-cell">
							{item.isDir ? 'Folder' : getFileTypeDescription(item.name)}
						</td>
						<td class="size-cell">
							{item.isDir ? '' : formatFileSize(item.size)}
						</td>
						<td class="date-cell">
							{formatFileDate(item.modTime)}
						</td>
					</tr>
				{/each}
			{/if}
		</tbody>
	</table>
</div>

<style>
	.file-list {
		position: relative;
		width: 100%;
		height: 100%;
		overflow: auto;
		background: #1e1e1e;
	}

	.loading-overlay {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(30, 30, 30, 0.8);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 10;
	}

	.spinner {
		width: 24px;
		height: 24px;
		border: 2px solid #333;
		border-top-color: #4a9eff;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.file-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 13px;
	}

	.file-table th {
		text-align: left;
		padding: 8px 12px;
		background: #252525;
		border-bottom: 1px solid #333;
		font-weight: 500;
		color: #888;
		white-space: nowrap;
		user-select: none;
		position: sticky;
		top: 0;
		z-index: 5;
	}

	.file-table th.sortable {
		cursor: pointer;
		transition: background-color 0.1s ease;
	}

	.file-table th.sortable:hover {
		background: #2a2a2a;
		color: #aaa;
	}

	.file-table th.sortable:focus {
		outline: 1px solid #4a9eff;
		outline-offset: -1px;
	}

	.file-table th.sorted {
		color: #4a9eff;
	}

	.col-label {
		margin-right: 4px;
	}

	.sort-indicator {
		font-size: 10px;
		opacity: 0.8;
	}

	.name-col {
		min-width: 200px;
	}

	.type-col {
		width: 120px;
	}

	.size-col {
		width: 100px;
		text-align: right;
	}

	.date-col {
		width: 150px;
	}

	.file-table td {
		padding: 6px 12px;
		border-bottom: 1px solid #2a2a2a;
		color: #ccc;
	}

	.file-row {
		cursor: default;
		transition: background-color 0.05s ease;
	}

	.file-row:hover {
		background: #252525;
	}

	.file-row:focus {
		outline: none;
		background: #2d4a6f;
	}

	.file-row.selected {
		background: #2d4a6f;
	}

	.file-row.selected:hover {
		background: #345580;
	}

	.file-row.directory .file-name {
		color: #e8c36a;
	}

	.name-cell {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.file-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		width: 20px;
		color: #888;
	}

	.file-icon.folder {
		color: #e8c36a;
	}

	.file-name {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.type-cell {
		color: #888;
	}

	.size-cell {
		text-align: right;
		font-variant-numeric: tabular-nums;
		color: #888;
	}

	.date-cell {
		white-space: nowrap;
		color: #888;
	}

	.empty-row td {
		padding: 48px 12px;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8px;
		color: #555;
	}

	.empty-icon {
		font-size: 32px;
		opacity: 0.5;
	}

	.empty-text {
		font-size: 13px;
	}

	/* Compact mode */
	.file-list.compact .file-table th {
		padding: 6px 12px;
	}

	.file-list.compact .file-table td {
		padding: 4px 12px;
	}

	.file-list.compact .file-icon {
		width: 16px;
	}

	.file-list.compact .name-cell {
		gap: 6px;
	}
</style>
