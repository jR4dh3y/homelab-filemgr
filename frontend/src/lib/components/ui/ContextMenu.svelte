<script lang="ts">
	/**
	 * Context Menu component - reusable right-click menu
	 * Follows UI component patterns from contributing guidelines
	 */
	import type { Snippet } from 'svelte';

	export interface ContextMenuItem {
		id: string;
		label: string;
		icon?: typeof import('lucide-svelte').Copy;
		disabled?: boolean;
		separator?: boolean;
		shortcut?: string;
	}

	interface Props {
		items: ContextMenuItem[];
		x: number;
		y: number;
		onSelect: (id: string) => void;
		onClose: () => void;
	}

	let { items, x, y, onSelect, onClose }: Props = $props();

	let menuRef: HTMLDivElement | undefined = $state();

	// Adjust position to keep menu within viewport
	let adjustedPosition = $derived.by(() => {
		if (!menuRef) return { x, y };
		
		const menuWidth = 200; // approximate width
		const menuHeight = items.length * 36; // approximate height
		const viewportWidth = typeof window !== 'undefined' ? window.innerWidth : 1920;
		const viewportHeight = typeof window !== 'undefined' ? window.innerHeight : 1080;
		
		let adjustedX = x;
		let adjustedY = y;
		
		if (x + menuWidth > viewportWidth) {
			adjustedX = viewportWidth - menuWidth - 8;
		}
		if (y + menuHeight > viewportHeight) {
			adjustedY = viewportHeight - menuHeight - 8;
		}
		
		return { x: adjustedX, y: adjustedY };
	});

	function handleItemClick(item: ContextMenuItem) {
		if (item.disabled || item.separator) return;
		onSelect(item.id);
		onClose();
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			onClose();
		}
	}

	function handleBackdropClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			onClose();
		}
	}
</script>

<svelte:window onkeydown={handleKeyDown} />

<!-- svelte-ignore a11y_no_static_element_interactions -->
<!-- svelte-ignore a11y_click_events_have_key_events -->
<div
	class="fixed inset-0 z-50"
	onclick={handleBackdropClick}
	oncontextmenu={(e) => { e.preventDefault(); onClose(); }}
>
	<div
		bind:this={menuRef}
		class="fixed bg-surface-primary border border-border-primary rounded-lg shadow-xl py-1 min-w-45 max-w-70"
		style="left: {adjustedPosition.x}px; top: {adjustedPosition.y}px;"
		role="menu"
	>
		{#each items as item (item.id)}
			{#if item.separator}
				<div class="h-px bg-border-secondary my-1 mx-2"></div>
			{:else}
				<button
					type="button"
					class="w-full flex items-center gap-3 px-3 py-2 text-left text-[13px] transition-colors
						{item.disabled 
							? 'text-text-disabled cursor-not-allowed' 
							: 'text-text-primary hover:bg-surface-secondary cursor-pointer'}"
					disabled={item.disabled}
					onclick={() => handleItemClick(item)}
					role="menuitem"
				>
					{#if item.icon}
					{@const IconComponent = item.icon}
					<span class="flex items-center justify-center w-4 h-4 text-text-secondary">
						<IconComponent size={14} />
						</span>
					{:else}
						<span class="w-4"></span>
					{/if}
					<span class="flex-1">{item.label}</span>
					{#if item.shortcut}
						<span class="text-text-muted text-[11px]">{item.shortcut}</span>
					{/if}
				</button>
			{/if}
		{/each}
	</div>
</div>
