<script lang="ts">
	/**
	 * FileList component with sortable columns
	 * Requirements: 1.1, 1.2
	 */

	import type { FileInfo } from '$lib/api/files';
	import type { SortField, SortDir } from '$lib/types/files';
	import { formatFileSize, formatFileDate, getFileTypeDescription } from '$lib/utils/format';

	let {
		items = [],
		sortBy = 'name',
		sortDir = 'asc',
		selectedPaths = new Set<string>(),
		isLoading = false,
		onItemClick,
		onSortChange,
		onSelectionChange
	}: {
		items?: FileInfo[];
		sortBy?: SortField;
		sortDir?: SortDir;
		selectedPaths?: Set<string>;
		isLoading?: boolean;
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
			const newSelection = new Set<string>(selectedPaths);
			if (newSelection.has(item.path)) {
				newSelection.delete(item.path);
			} else {
				newSelection.add(item.path);
			}
			onSelectionChange?.(newSelection);
		} else {
			onItemClick?.(item);
		}
	}

	function handleKeyDown(item: FileInfo, event: KeyboardEvent) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			onItemClick?.(item);
		}
	}

	function getSortIndicator(field: SortField): string {
		if (sortBy !== field) return '';
		return sortDir === 'asc' ? '↑' : '↓';
	}

	function getFileIcon(item: FileInfo): string {
		if (item.isDir) return '📁';
		const ext = item.name.includes('.') 
			? item.name.slice(item.name.lastIndexOf('.') + 1).toLowerCase()
			: '';
		const iconMap: Record<string, string> = {
			jpg: '🖼️', jpeg: '🖼️', png: '🖼️', gif: '🖼️', svg: '🖼️', webp: '🖼️',
			mp4: '🎬', mkv: '🎬', avi: '🎬', mov: '🎬', webm: '🎬',
			mp3: '🎵', wav: '🎵', flac: '🎵', aac: '🎵', ogg: '🎵',
			pdf: '📄', doc: '📝', docx: '📝', txt: '📝', md: '📝',
			xls: '📊', xlsx: '📊', csv: '📊',
			zip: '📦', rar: '📦', '7z': '📦', tar: '📦', gz: '📦',
			js: '💻', ts: '💻', py: '💻', go: '💻', rs: '💻', java: '💻',
			html: '🌐', css: '🎨', json: '📋', xml: '📋', yaml: '📋', yml: '📋'
		};
		return iconMap[ext] || '📄';
	}

	function isSelected(path: string): boolean {
		return selectedPaths.has(path);
	}
</script>

<div class="file-list">
	{#if isLoading}
		<div class="loading-overlay">
			<div class="spinner"></div>
			<span>Loading...</span>
		</div>
	{/if}

	<table class="file-table" role="grid" aria-busy={isLoading}>
		<thead>
			<tr>
				<th 
					class="sortable" 
					class:sorted={sortBy === 'name'}
					onclick={() => handleSort('name')}
					onkeydown={(e) => e.key === 'Enter' && handleSort('name')}
					tabindex="0"
					role="columnheader"
					aria-sort={sortBy === 'name' ? (sortDir === 'asc' ? 'ascending' : 'descending') : 'none'}
				>
					Name {getSortIndicator('name')}
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
					Size {getSortIndicator('size')}
				</th>
				<th 
					class="sortable date-col" 
					class:sorted={sortBy === 'modTime'}
					onclick={() => handleSort('modTime')}
					onkeydown={(e) => e.key === 'Enter' && handleSort('modTime')}
					tabindex="0"
					role="columnheader"
					aria-sort={sortBy === 'modTime' ? (sortDir === 'asc' ? 'ascending' : 'descending') : 'none'}
				>
					Modified {getSortIndicator('modTime')}
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
					Type {getSortIndicator('type')}
				</th>
			</tr>
		</thead>
		<tbody>
			{#if items.length === 0 && !isLoading}
				<tr class="empty-row">
					<td colspan="4">
						<div class="empty-state">
							<span class="empty-icon">📂</span>
							<span>This folder is empty</span>
						</div>
					</td>
				</tr>
			{:else}
				{#each items as item (item.path)}
					<tr 
						class="file-row"
						class:selected={isSelected(item.path)}
						class:directory={item.isDir}
						onclick={(e) => handleRowClick(item, e)}
						onkeydown={(e) => handleKeyDown(item, e)}
						ondblclick={() => onItemClick?.(item)}
						tabindex="0"
						aria-selected={isSelected(item.path)}
					>
						<td class="name-cell">
							<span class="file-icon">{getFileIcon(item)}</span>
							<span class="file-name" title={item.name}>{item.name}</span>
						</td>
						<td class="size-cell">
							{item.isDir ? '-' : formatFileSize(item.size)}
						</td>
						<td class="date-cell">
							{formatFileDate(item.modTime)}
						</td>
						<td class="type-cell">
							{item.isDir ? 'Folder' : getFileTypeDescription(item.name)}
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
		overflow: auto;
	}

	.loading-overlay {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(255, 255, 255, 0.8);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		z-index: 10;
	}

	.spinner {
		width: 24px;
		height: 24px;
		border: 3px solid #e5e7eb;
		border-top-color: #3b82f6;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.file-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.875rem;
	}

	.file-table th {
		text-align: left;
		padding: 0.75rem 1rem;
		background: #f9fafb;
		border-bottom: 2px solid #e5e7eb;
		font-weight: 600;
		color: #374151;
		white-space: nowrap;
		user-select: none;
	}

	.file-table th.sortable {
		cursor: pointer;
		transition: background-color 0.15s;
	}

	.file-table th.sortable:hover {
		background: #f3f4f6;
	}

	.file-table th.sortable:focus {
		outline: 2px solid #3b82f6;
		outline-offset: -2px;
	}

	.file-table th.sorted {
		background: #eff6ff;
		color: #1d4ed8;
	}

	.size-col { width: 100px; text-align: right; }
	.date-col { width: 140px; }
	.type-col { width: 160px; }

	.file-table td {
		padding: 0.625rem 1rem;
		border-bottom: 1px solid #e5e7eb;
		color: #4b5563;
	}

	.file-row {
		cursor: pointer;
		transition: background-color 0.15s;
	}

	.file-row:hover { background: #f9fafb; }
	.file-row:focus { outline: 2px solid #3b82f6; outline-offset: -2px; background: #eff6ff; }
	.file-row.selected { background: #dbeafe; }
	.file-row.selected:hover { background: #bfdbfe; }
	.file-row.directory .file-name { font-weight: 500; }

	.name-cell {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		min-width: 200px;
	}

	.file-icon { font-size: 1.125rem; flex-shrink: 0; }
	.file-name { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
	.size-cell { text-align: right; font-variant-numeric: tabular-nums; }
	.date-cell { white-space: nowrap; }
	.type-cell { color: #6b7280; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

	.empty-row td { padding: 3rem 1rem; }
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		color: #9ca3af;
	}
	.empty-icon { font-size: 2.5rem; }

	@media (prefers-color-scheme: dark) {
		.loading-overlay { background: rgba(17, 24, 39, 0.8); }
		.spinner { border-color: #374151; border-top-color: #60a5fa; }
		.file-table th { background: #1f2937; border-bottom-color: #374151; color: #e5e7eb; }
		.file-table th.sortable:hover { background: #374151; }
		.file-table th.sorted { background: #1e3a5f; color: #93c5fd; }
		.file-table td { border-bottom-color: #374151; color: #d1d5db; }
		.file-row:hover { background: #1f2937; }
		.file-row:focus { background: #1e3a5f; }
		.file-row.selected { background: #1e3a5f; }
		.file-row.selected:hover { background: #1e40af; }
		.type-cell { color: #9ca3af; }
		.empty-state { color: #6b7280; }
	}
</style>
