<script lang="ts">
	/**
	 * Toolbar component - navigation buttons and path bar
	 */
	import { ChevronLeft, ChevronRight, ChevronUp, Home, RefreshCw, Settings } from 'lucide-svelte';

	interface Props {
		pathSegments?: string[];
		canGoBack?: boolean;
		canGoForward?: boolean;
		canGoUp?: boolean;
		onBack?: () => void;
		onForward?: () => void;
		onUp?: () => void;
		onNavigate?: (path: string) => void;
		onRefresh?: () => void;
		onSettings?: () => void;
	}

	let {
		pathSegments = [],
		canGoBack = false,
		canGoForward = false,
		canGoUp = false,
		onBack,
		onForward,
		onUp,
		onNavigate,
		onRefresh,
		onSettings
	}: Props = $props();

	function buildPath(index: number): string {
		return pathSegments.slice(0, index + 1).join('/');
	}

	function handleSegmentClick(index: number) {
		onNavigate?.(buildPath(index));
	}

	function handleRootClick() {
		onNavigate?.('');
	}
</script>

<div class="toolbar">
	<!-- Navigation buttons -->
	<div class="nav-buttons">
		<button
			type="button"
			class="nav-btn"
			disabled={!canGoBack}
			onclick={onBack}
			title="Back"
		>
			<ChevronLeft size={18} />
		</button>
		<button
			type="button"
			class="nav-btn"
			disabled={!canGoForward}
			onclick={onForward}
			title="Forward"
		>
			<ChevronRight size={18} />
		</button>
		<button
			type="button"
			class="nav-btn"
			disabled={!canGoUp}
			onclick={onUp}
			title="Up"
		>
			<ChevronUp size={18} />
		</button>
	</div>

	<!-- Path bar -->
	<div class="path-bar">
		<button type="button" class="path-icon" onclick={handleRootClick} title="Go to root">
			<Home size={14} />
		</button>
		<div class="path-segments">
			{#if pathSegments.length === 0}
				<span class="path-segment current">This Server</span>
			{:else}
				{#each pathSegments as segment, index (index)}
					{#if index > 0}
						<span class="path-separator">/</span>
					{/if}
					{#if index === pathSegments.length - 1}
						<span class="path-segment current">{segment}</span>
					{:else}
						<button
							type="button"
							class="path-segment clickable"
							onclick={() => handleSegmentClick(index)}
						>
							{segment}
						</button>
					{/if}
				{/each}
			{/if}
		</div>
	</div>

	<!-- Action buttons -->
	<div class="action-buttons">
		<button type="button" class="action-btn" onclick={onRefresh} title="Refresh">
			<RefreshCw size={16} />
		</button>
		<button type="button" class="action-btn" onclick={onSettings} title="Settings">
			<Settings size={16} />
		</button>
	</div>
</div>

<style>
	.toolbar {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 6px 12px;
		background: #1e1e1e;
		border-bottom: 1px solid #2a2a2a;
	}

	.nav-buttons {
		display: flex;
		gap: 2px;
	}

	.nav-btn {
		width: 28px;
		height: 28px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: none;
		border-radius: 4px;
		color: #888;
		cursor: pointer;
		transition: all 0.1s ease;
	}

	.nav-btn:hover:not(:disabled) {
		background: #333;
		color: #ccc;
	}

	.nav-btn:disabled {
		color: #444;
		cursor: not-allowed;
	}

	.path-bar {
		flex: 1;
		display: flex;
		align-items: center;
		gap: 6px;
		background: #252525;
		border: 1px solid #333;
		border-radius: 4px;
		padding: 4px 8px;
		min-width: 0;
	}

	.path-icon {
		width: 18px;
		height: 18px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: none;
		color: #888;
		cursor: pointer;
		flex-shrink: 0;
	}

	.path-icon:hover {
		color: #ccc;
	}

	.path-segments {
		display: flex;
		align-items: center;
		gap: 4px;
		overflow: hidden;
		flex: 1;
	}

	.path-segment {
		color: #888;
		font-size: 13px;
		white-space: nowrap;
	}

	.path-segment.current {
		color: #ccc;
	}

	.path-segment.clickable {
		background: transparent;
		border: none;
		padding: 0;
		cursor: pointer;
		font-size: 13px;
	}

	.path-segment.clickable:hover {
		color: #fff;
		text-decoration: underline;
	}

	.path-separator {
		color: #555;
		font-size: 12px;
	}

	.action-buttons {
		display: flex;
		gap: 4px;
	}

	.action-btn {
		width: 28px;
		height: 28px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: none;
		border-radius: 4px;
		color: #888;
		cursor: pointer;
		transition: all 0.1s ease;
	}

	.action-btn:hover {
		background: #333;
		color: #ccc;
	}
</style>
