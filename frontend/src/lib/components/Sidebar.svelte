<script lang="ts">
	/**
	 * Sidebar component - navigation panel with storage, places, favorites
	 */
	import type { MountPoint } from '$lib/api/files';
	import {
		ChevronDown,
		HardDrive,
		Server,
		Monitor,
		Download,
		FileText,
		Music,
		Image,
		Video,
		Star,
		Folder
	} from 'lucide-svelte';

	interface Props {
		roots?: MountPoint[];
		currentPath?: string;
		onNavigate?: (path: string) => void;
	}

	let { roots = [], currentPath = '', onNavigate }: Props = $props();

	// Quick access places
	const places = [
		{ name: 'This Server', path: '', icon: Server },
		{ name: 'Desktop', path: 'desktop', icon: Monitor },
		{ name: 'Downloads', path: 'downloads', icon: Download },
		{ name: 'Documents', path: 'documents', icon: FileText },
		{ name: 'Music', path: 'music', icon: Music },
		{ name: 'Pictures', path: 'pictures', icon: Image },
		{ name: 'Videos', path: 'videos', icon: Video }
	];

	// Favorites (could be stored in localStorage later)
	let favorites = $state<{ name: string; path: string }[]>([]);

	function isActive(path: string): boolean {
		if (path === '' && currentPath === '') return true;
		return currentPath.startsWith(path) && path !== '';
	}

	function handleNavigate(path: string) {
		onNavigate?.(path);
	}

	// Collapsed sections state
	let storageCollapsed = $state(false);
	let placesCollapsed = $state(false);
	let favoritesCollapsed = $state(false);
</script>

<aside class="sidebar">
	<!-- Storage Section -->
	<div class="sidebar-section">
		<button 
			type="button" 
			class="section-header"
			onclick={() => storageCollapsed = !storageCollapsed}
		>
			<ChevronDown size={14} class="collapse-icon {storageCollapsed ? 'collapsed' : ''}" />
			<span>Storage</span>
		</button>
		{#if !storageCollapsed}
			<div class="section-content">
				{#each roots as root (root.name)}
					<button
						type="button"
						class="nav-item"
						class:active={isActive(root.name)}
						onclick={() => handleNavigate(root.name)}
					>
						<HardDrive size={16} class="nav-icon" />
						<span class="nav-label">{root.name}</span>
						{#if root.readOnly}
							<span class="badge">RO</span>
						{/if}
					</button>
				{/each}
			</div>
		{/if}
	</div>

	<!-- Places Section -->
	<div class="sidebar-section">
		<button 
			type="button" 
			class="section-header"
			onclick={() => placesCollapsed = !placesCollapsed}
		>
			<ChevronDown size={14} class="collapse-icon {placesCollapsed ? 'collapsed' : ''}" />
			<span>Places</span>
		</button>
		{#if !placesCollapsed}
			<div class="section-content">
				{#each places as place (place.path)}
					<button
						type="button"
						class="nav-item"
						class:active={isActive(place.path)}
						onclick={() => handleNavigate(place.path)}
					>
						<place.icon size={16} class="nav-icon" />
						<span class="nav-label">{place.name}</span>
					</button>
				{/each}
			</div>
		{/if}
	</div>

	<!-- Favorites Section -->
	<div class="sidebar-section">
		<button 
			type="button" 
			class="section-header"
			onclick={() => favoritesCollapsed = !favoritesCollapsed}
		>
			<ChevronDown size={14} class="collapse-icon {favoritesCollapsed ? 'collapsed' : ''}" />
			<span>Favorites</span>
		</button>
		{#if !favoritesCollapsed}
			<div class="section-content">
				{#if favorites.length === 0}
					<div class="empty-favorites">No favorites yet</div>
				{:else}
					{#each favorites as fav (fav.path)}
						<button
							type="button"
							class="nav-item"
							class:active={isActive(fav.path)}
							onclick={() => handleNavigate(fav.path)}
						>
							<Star size={16} class="nav-icon" />
							<span class="nav-label">{fav.name}</span>
						</button>
					{/each}
				{/if}
			</div>
		{/if}
	</div>
</aside>

<style>
	.sidebar {
		width: 220px;
		min-width: 220px;
		background: #1a1a1a;
		border-right: 1px solid #2a2a2a;
		display: flex;
		flex-direction: column;
		overflow-y: auto;
		overflow-x: hidden;
	}

	.sidebar-section {
		border-bottom: 1px solid #2a2a2a;
	}

	.section-header {
		width: 100%;
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 10px 12px;
		background: transparent;
		border: none;
		color: #888;
		font-size: 11px;
		font-weight: 500;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		cursor: pointer;
		text-align: left;
	}

	.section-header:hover {
		color: #aaa;
	}

	:global(.collapse-icon) {
		transition: transform 0.15s ease;
		flex-shrink: 0;
	}

	:global(.collapse-icon.collapsed) {
		transform: rotate(-90deg);
	}

	.section-content {
		padding-bottom: 8px;
	}

	.nav-item {
		width: 100%;
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 7px 12px 7px 20px;
		background: transparent;
		border: none;
		color: #ccc;
		font-size: 13px;
		cursor: pointer;
		text-align: left;
		transition: background-color 0.1s ease;
	}

	.nav-item:hover {
		background: #252525;
	}

	.nav-item.active {
		background: #2d4a6f;
		color: #fff;
	}

	:global(.nav-icon) {
		flex-shrink: 0;
		opacity: 0.8;
	}

	.nav-item.active :global(.nav-icon) {
		opacity: 1;
	}

	.nav-label {
		flex: 1;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.badge {
		font-size: 9px;
		padding: 1px 4px;
		background: #444;
		color: #999;
		border-radius: 3px;
		flex-shrink: 0;
	}

	.empty-favorites {
		padding: 8px 20px;
		color: #555;
		font-size: 12px;
		font-style: italic;
	}
</style>
