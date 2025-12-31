<script lang="ts">
	/**
	 * Sidebar component - navigation panel with storage, places, favorites
	 */
	import type { MountPoint } from '$lib/api/files';

	interface Props {
		roots?: MountPoint[];
		currentPath?: string;
		onNavigate?: (path: string) => void;
	}

	let { roots = [], currentPath = '', onNavigate }: Props = $props();

	// Quick access places
	const places = [
		{ name: 'This PC', path: '', icon: '💻' },
		{ name: 'Desktop', path: 'desktop', icon: '🖥️' },
		{ name: 'Downloads', path: 'downloads', icon: '⬇️' },
		{ name: 'Documents', path: 'documents', icon: '📄' },
		{ name: 'Music', path: 'music', icon: '🎵' },
		{ name: 'Pictures', path: 'pictures', icon: '🖼️' },
		{ name: 'Videos', path: 'videos', icon: '🎬' }
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
	<!-- Search/Filter -->
	<div class="sidebar-search">
		<svg class="search-icon" viewBox="0 0 20 20" fill="currentColor">
			<path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
		</svg>
		<input type="text" placeholder="Filter / New..." class="search-input" />
	</div>

	<!-- Storage Section -->
	<div class="sidebar-section">
		<button 
			type="button" 
			class="section-header"
			onclick={() => storageCollapsed = !storageCollapsed}
		>
			<svg class="collapse-icon" class:collapsed={storageCollapsed} viewBox="0 0 20 20" fill="currentColor">
				<path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
			</svg>
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
						<span class="nav-icon">💾</span>
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
			<svg class="collapse-icon" class:collapsed={placesCollapsed} viewBox="0 0 20 20" fill="currentColor">
				<path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
			</svg>
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
						<span class="nav-icon">{place.icon}</span>
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
			<svg class="collapse-icon" class:collapsed={favoritesCollapsed} viewBox="0 0 20 20" fill="currentColor">
				<path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
			</svg>
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
							<span class="nav-icon">⭐</span>
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

	.sidebar-search {
		padding: 8px 12px;
		border-bottom: 1px solid #2a2a2a;
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.search-icon {
		width: 14px;
		height: 14px;
		color: #666;
		flex-shrink: 0;
	}

	.search-input {
		flex: 1;
		background: transparent;
		border: none;
		color: #888;
		font-size: 12px;
		outline: none;
	}

	.search-input::placeholder {
		color: #555;
	}

	.sidebar-section {
		border-bottom: 1px solid #2a2a2a;
	}

	.section-header {
		width: 100%;
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 8px 12px;
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

	.collapse-icon {
		width: 14px;
		height: 14px;
		transition: transform 0.15s ease;
	}

	.collapse-icon.collapsed {
		transform: rotate(-90deg);
	}

	.section-content {
		padding-bottom: 8px;
	}

	.nav-item {
		width: 100%;
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 6px 12px 6px 24px;
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

	.nav-icon {
		font-size: 14px;
		width: 18px;
		text-align: center;
		flex-shrink: 0;
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
		padding: 8px 24px;
		color: #555;
		font-size: 12px;
		font-style: italic;
	}
</style>
