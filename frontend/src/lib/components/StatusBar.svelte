<script lang="ts">
	/**
	 * StatusBar component - bottom status bar with item count and view options
	 */

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
				<svg viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
				</svg>
			</button>
			<button
				type="button"
				class="view-btn"
				class:active={viewMode === 'grid'}
				onclick={() => onViewModeChange?.('grid')}
				title="Grid view"
			>
				<svg viewBox="0 0 20 20" fill="currentColor">
					<path d="M5 3a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2V5a2 2 0 00-2-2H5zM5 11a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2v-2a2 2 0 00-2-2H5zM11 5a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V5zM11 13a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
				</svg>
			</button>
		</div>
		<span class="details-label">Details</span>
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

	.view-btn svg {
		width: 14px;
		height: 14px;
	}

	.details-label {
		color: #666;
	}
</style>
