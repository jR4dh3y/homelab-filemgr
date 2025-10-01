/**
 * Files store for managing file listing state
 * Requirements: 1.1, 1.2
 */

import { writable, derived, get } from 'svelte/store';
import {
	listRoots,
	listDirectory,
	createDirectory,
	rename,
	deleteFile,
	search,
	type FileInfo,
	type FileList,
	type ListOptions,
	type MountPoint,
	type RootsResponse,
	type SearchResponse
} from '$lib/api/files';

/**
 * Current path state
 */
export interface PathState {
	currentPath: string;
	pathSegments: string[];
}

/**
 * Create the path store
 */
function createPathStore() {
	const { subscribe, set, update } = writable<PathState>({
		currentPath: '',
		pathSegments: []
	});

	/**
	 * Navigate to a path
	 */
	function navigateTo(path: string): void {
		const normalizedPath = path.replace(/^\/+|\/+$/g, '');
		const segments = normalizedPath ? normalizedPath.split('/') : [];
		set({
			currentPath: normalizedPath,
			pathSegments: segments
		});
	}

	/**
	 * Navigate up one level
	 */
	function navigateUp(): void {
		update((state) => {
			if (state.pathSegments.length === 0) return state;
			const newSegments = state.pathSegments.slice(0, -1);
			return {
				currentPath: newSegments.join('/'),
				pathSegments: newSegments
			};
		});
	}

	/**
	 * Navigate to root
	 */
	function navigateToRoot(): void {
		set({
			currentPath: '',
			pathSegments: []
		});
	}

	/**
	 * Get current path
	 */
	function getCurrentPath(): string {
		return get({ subscribe }).currentPath;
	}

	return {
		subscribe,
		navigateTo,
		navigateUp,
		navigateToRoot,
		getCurrentPath
	};
}

/**
 * Path store singleton
 */
export const pathStore = createPathStore();

/**
 * Derived store for current path string
 */
export const currentPath = derived(pathStore, ($path) => $path.currentPath);

/**
 * Derived store for path segments (for breadcrumbs)
 */
export const pathSegments = derived(pathStore, ($path) => $path.pathSegments);

/**
 * List options state
 */
export interface ListOptionsState extends ListOptions {
	page: number;
	pageSize: number;
	sortBy: 'name' | 'size' | 'modTime' | 'type';
	sortDir: 'asc' | 'desc';
	filter: string;
}

/**
 * Default list options
 */
const defaultListOptions: ListOptionsState = {
	page: 1,
	pageSize: 50,
	sortBy: 'name',
	sortDir: 'asc',
	filter: ''
};

/**
 * Create the list options store
 */
function createListOptionsStore() {
	const { subscribe, set, update } = writable<ListOptionsState>(defaultListOptions);

	function setPage(page: number): void {
		update((state) => ({ ...state, page }));
	}

	function setPageSize(pageSize: number): void {
		update((state) => ({ ...state, pageSize, page: 1 }));
	}

	function setSortBy(sortBy: ListOptionsState['sortBy']): void {
		update((state) => ({ ...state, sortBy, page: 1 }));
	}

	function setSortDir(sortDir: ListOptionsState['sortDir']): void {
		update((state) => ({ ...state, sortDir, page: 1 }));
	}

	function toggleSortDir(): void {
		update((state) => ({
			...state,
			sortDir: state.sortDir === 'asc' ? 'desc' : 'asc',
			page: 1
		}));
	}

	function setFilter(filter: string): void {
		update((state) => ({ ...state, filter, page: 1 }));
	}

	function reset(): void {
		set(defaultListOptions);
	}

	function getOptions(): ListOptionsState {
		return get({ subscribe });
	}

	return {
		subscribe,
		setPage,
		setPageSize,
		setSortBy,
		setSortDir,
		toggleSortDir,
		setFilter,
		reset,
		getOptions
	};
}

/**
 * List options store singleton
 */
export const listOptionsStore = createListOptionsStore();

/**
 * Query key factory for files
 */
export const fileQueryKeys = {
	all: ['files'] as const,
	roots: () => [...fileQueryKeys.all, 'roots'] as const,
	list: (path: string, options: ListOptions) =>
		[...fileQueryKeys.all, 'list', path, options] as const,
	search: (path: string, query: string) =>
		[...fileQueryKeys.all, 'search', path, query] as const
};

/**
 * Query options factory for listing mount points (roots)
 */
export function rootsQueryOptions() {
	return {
		queryKey: fileQueryKeys.roots(),
		queryFn: () => listRoots()
	};
}

/**
 * Query options factory for listing directory contents
 */
export function directoryQueryOptions(path: string, options: ListOptions) {
	return {
		queryKey: fileQueryKeys.list(path, options),
		queryFn: () => listDirectory(path, options),
		enabled: path !== ''
	};
}

/**
 * Query options factory for searching files
 */
export function searchQueryOptions(path: string, query: string) {
	return {
		queryKey: fileQueryKeys.search(path, query),
		queryFn: () => search(path, query),
		enabled: query.length > 0
	};
}

/**
 * Mutation options for creating directories
 */
export function createDirectoryMutationOptions() {
	return {
		mutationFn: ({ basePath, name }: { basePath: string; name: string }) =>
			createDirectory(basePath, name)
	};
}

/**
 * Mutation options for renaming files/directories
 */
export function renameMutationOptions() {
	return {
		mutationFn: ({ oldPath, newPath }: { oldPath: string; newPath: string }) =>
			rename(oldPath, newPath)
	};
}

/**
 * Mutation options for deleting files/directories
 */
export function deleteMutationOptions() {
	return {
		mutationFn: ({ path, confirm }: { path: string; confirm?: boolean }) =>
			deleteFile(path, confirm)
	};
}

/**
 * Selection state for multi-select operations
 */
export interface SelectionState {
	selectedItems: Set<string>;
	lastSelectedItem: string | null;
}

/**
 * Create the selection store
 */
function createSelectionStore() {
	const { subscribe, set, update } = writable<SelectionState>({
		selectedItems: new Set(),
		lastSelectedItem: null
	});

	function select(path: string): void {
		update((state) => {
			const newSelected = new Set(state.selectedItems);
			newSelected.add(path);
			return {
				selectedItems: newSelected,
				lastSelectedItem: path
			};
		});
	}

	function deselect(path: string): void {
		update((state) => {
			const newSelected = new Set(state.selectedItems);
			newSelected.delete(path);
			return {
				...state,
				selectedItems: newSelected
			};
		});
	}

	function toggle(path: string): void {
		update((state) => {
			const newSelected = new Set(state.selectedItems);
			if (newSelected.has(path)) {
				newSelected.delete(path);
			} else {
				newSelected.add(path);
			}
			return {
				selectedItems: newSelected,
				lastSelectedItem: path
			};
		});
	}

	function selectOnly(path: string): void {
		set({
			selectedItems: new Set([path]),
			lastSelectedItem: path
		});
	}

	function selectAll(paths: string[]): void {
		set({
			selectedItems: new Set(paths),
			lastSelectedItem: paths[paths.length - 1] || null
		});
	}

	function clearSelection(): void {
		set({
			selectedItems: new Set(),
			lastSelectedItem: null
		});
	}

	function isSelected(path: string): boolean {
		return get({ subscribe }).selectedItems.has(path);
	}

	function getSelectedItems(): string[] {
		return Array.from(get({ subscribe }).selectedItems);
	}

	return {
		subscribe,
		select,
		deselect,
		toggle,
		selectOnly,
		selectAll,
		clearSelection,
		isSelected,
		getSelectedItems
	};
}

/**
 * Selection store singleton
 */
export const selectionStore = createSelectionStore();

/**
 * Derived store for selected items count
 */
export const selectedCount = derived(
	selectionStore,
	($selection) => $selection.selectedItems.size
);

/**
 * Derived store for whether any items are selected
 */
export const hasSelection = derived(selectionStore, ($selection) => $selection.selectedItems.size > 0);
