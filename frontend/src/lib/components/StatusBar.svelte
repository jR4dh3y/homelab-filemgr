<script lang="ts">
	/**
	 * StatusBar component - bottom status bar with item count and view options
	 */
	import { List, LayoutGrid } from 'lucide-svelte';

	interface Props {
		itemCount?: number;
		selectedCount?: number;
		viewMode?: 'list' | 'grid';
		onViewModeChange?: (mode: 'list' | 'grid') => void;
	}

	let {
		itemCount = 0,
		selectedCount = 0,
		viewMode = 'list',
		onViewModeChange
	}: Props = $props();

	const statusText = $derived.by(() => {
		if (selectedCount > 0) {
			return `${selectedCount} of ${itemCount} selected`;
		}
		return `${itemCount} item${itemCount !== 1 ? 's' : ''}`;
	});
</script>

<footer class="status-bar">
	<div class="status-left">
		<span class="item-count">{statusText}</span>
	</div>
	<div class="status-right">
		<div class="view-toggle">
			<button
				type="button"
				class="view-btn"
				class:active={viewMode === 'list'}
				onclick={() => onViewModeChange?.('list')}
				title="List view"
			>
				<List size={14} />
			</button>
			<button
				type="button"
				class="view-btn"
				class:active={viewMode === 'grid'}
				onclick={() => onViewModeChange?.('grid')}
				title="Grid view"
			>
				<LayoutGrid size={14} />
			</button>
		</div>
	</div>
</footer>

<style>
	.status-bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 4px 12px;
		background: #1a1a1a;
		border-top: 1px solid #2a2a2a;
		font-size: 12px;
		color: #888;
	}

	.status-left {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.item-count {
		color: #888;
	}

	.status-right {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.view-toggle {
		display: flex;
		gap: 2px;
		background: #252525;
		border-radius: 3px;
		padding: 2px;
	}

	.view-btn {
		width: 22px;
		height: 22px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: none;
		border-radius: 2px;
		color: #666;
		cursor: pointer;
		transition: all 0.1s ease;
	}

	.view-btn:hover {
		color: #888;
	}

	.view-btn.active {
		background: #333;
		color: #ccc;
	}
</style>
