<script lang="ts">
	/**
	 * Toolbar component - navigation buttons and path bar
	 */

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
		onRefresh
	}: Props = $props();

	const currentPathDisplay = $derived(
		pathSegments.length > 0 ? pathSegments.join(' / ') : 'This PC'
	);

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
			<svg viewBox="0 0 20 20" fill="currentColor">
				<path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
			</svg>
		</button>
		<button
			type="button"
			class="nav-btn"
			disabled={!canGoForward}
			onclick={onForward}
			title="Forward"
		>
			<svg viewBox="0 0 20 20" fill="currentColor">
				<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
			</svg>
		</button>
		<button
			type="button"
			class="nav-btn"
			disabled={!canGoUp}
			onclick={onUp}
			title="Up"
		>
			<svg viewBox="0 0 20 20" fill="currentColor">
				<path fill-rule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clip-rule="evenodd" />
			</svg>
		</button>
	</div>

	<!-- Path bar -->
	<div class="path-bar">
		<button type="button" class="path-icon" onclick={handleRootClick} title="Go to root">
			<svg viewBox="0 0 20 20" fill="currentColor">
				<path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z" />
			</svg>
		</button>
		<div class="path-segments">
			{#if pathSegments.length === 0}
				<span class="path-segment current">This PC</span>
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
			<svg viewBox="0 0 20 20" fill="currentColor">
				<path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd" />
			</svg>
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

	.nav-btn svg {
		width: 16px;
		height: 16px;
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

	.path-icon svg {
		width: 14px;
		height: 14px;
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

	.action-btn svg {
		width: 16px;
		height: 16px;
	}
</style>
