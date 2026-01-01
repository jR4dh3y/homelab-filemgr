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
	} from 'lucide-svelte';
	import { Badge } from '$lib/components/ui';

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
		{ name: 'Videos', path: 'videos', icon: Video },
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

	const navItemClass =
		'w-full flex items-center gap-2.5 py-1.5 px-3 pl-5 bg-transparent border-none text-text-primary text-[13px] cursor-pointer text-left transition-colors duration-100 hover:bg-surface-secondary';
	const navItemActiveClass = 'bg-selection text-white hover:bg-selection-hover';
</script>

<aside class="w-[220px] min-w-[220px] bg-surface-primary border-r border-border-secondary flex flex-col overflow-y-auto overflow-x-hidden">
	<!-- Storage Section -->
	<div class="border-b border-border-secondary">
		<button
			type="button"
			class="w-full flex items-center gap-1.5 px-3 py-2.5 bg-transparent border-none text-text-secondary text-[11px] font-medium uppercase tracking-wide cursor-pointer text-left hover:text-text-primary"
			onclick={() => (storageCollapsed = !storageCollapsed)}
		>
			<ChevronDown size={14} class="shrink-0 transition-transform duration-150 {storageCollapsed ? '-rotate-90' : ''}" />
			<span>Storage</span>
		</button>
		{#if !storageCollapsed}
			<div class="pb-2">
				{#each roots as root (root.name)}
					<button
						type="button"
						class="{navItemClass} {isActive(root.name) ? navItemActiveClass : ''}"
						onclick={() => handleNavigate(root.name)}
					>
						<HardDrive size={16} class="shrink-0 opacity-80" />
						<span class="flex-1 overflow-hidden text-ellipsis whitespace-nowrap">{root.name}</span>
						{#if root.readOnly}
							<Badge>RO</Badge>
						{/if}
					</button>
				{/each}
			</div>
		{/if}
	</div>

	<!-- Places Section -->
	<div class="border-b border-border-secondary">
		<button
			type="button"
			class="w-full flex items-center gap-1.5 px-3 py-2.5 bg-transparent border-none text-text-secondary text-[11px] font-medium uppercase tracking-wide cursor-pointer text-left hover:text-text-primary"
			onclick={() => (placesCollapsed = !placesCollapsed)}
		>
			<ChevronDown size={14} class="shrink-0 transition-transform duration-150 {placesCollapsed ? '-rotate-90' : ''}" />
			<span>Places</span>
		</button>
		{#if !placesCollapsed}
			<div class="pb-2">
				{#each places as place (place.path)}
					<button
						type="button"
						class="{navItemClass} {isActive(place.path) ? navItemActiveClass : ''}"
						onclick={() => handleNavigate(place.path)}
					>
						<place.icon size={16} class="shrink-0 opacity-80" />
						<span class="flex-1 overflow-hidden text-ellipsis whitespace-nowrap">{place.name}</span>
					</button>
				{/each}
			</div>
		{/if}
	</div>

	<!-- Favorites Section -->
	<div class="border-b border-border-secondary">
		<button
			type="button"
			class="w-full flex items-center gap-1.5 px-3 py-2.5 bg-transparent border-none text-text-secondary text-[11px] font-medium uppercase tracking-wide cursor-pointer text-left hover:text-text-primary"
			onclick={() => (favoritesCollapsed = !favoritesCollapsed)}
		>
			<ChevronDown size={14} class="shrink-0 transition-transform duration-150 {favoritesCollapsed ? '-rotate-90' : ''}" />
			<span>Favorites</span>
		</button>
		{#if !favoritesCollapsed}
			<div class="pb-2">
				{#if favorites.length === 0}
					<div class="py-2 px-5 text-text-muted text-xs italic">No favorites yet</div>
				{:else}
					{#each favorites as fav (fav.path)}
						<button
							type="button"
							class="{navItemClass} {isActive(fav.path) ? navItemActiveClass : ''}"
							onclick={() => handleNavigate(fav.path)}
						>
							<Star size={16} class="shrink-0 opacity-80" />
							<span class="flex-1 overflow-hidden text-ellipsis whitespace-nowrap">{fav.name}</span>
						</button>
					{/each}
				{/if}
			</div>
		{/if}
	</div>
</aside>
