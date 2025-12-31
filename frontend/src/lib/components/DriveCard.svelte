<script lang="ts">
	/**
	 * DriveCard component - displays drive/mount point with usage stats
	 */
	import type { DriveStats } from '$lib/api/files';
	import { HardDrive } from 'lucide-svelte';
	import { formatFileSize } from '$lib/utils/format';

	interface Props {
		drive: DriveStats;
		onClick?: () => void;
	}

	let { drive, onClick }: Props = $props();

	const usedFormatted = $derived(formatFileSize(drive.usedBytes));
	const totalFormatted = $derived(formatFileSize(drive.totalBytes));
	const freeFormatted = $derived(formatFileSize(drive.freeBytes));

	// Color based on usage percentage
	const barColor = $derived.by(() => {
		if (drive.usedPct >= 90) return '#e74c3c';
		if (drive.usedPct >= 75) return '#f39c12';
		return '#3498db';
	});
</script>

<button type="button" class="drive-card" onclick={onClick}>
	<div class="drive-header">
		<div class="drive-icon">
			<HardDrive size={24} />
		</div>
		<div class="drive-info">
			<span class="drive-name">{drive.name}</span>
			{#if drive.readOnly}
				<span class="drive-badge">Read-only</span>
			{/if}
		</div>
	</div>

	<div class="drive-stats">
		<div class="usage-bar-container">
			<div 
				class="usage-bar" 
				style="width: {drive.usedPct}%; background-color: {barColor};"
			></div>
		</div>
		<div class="usage-text">
			<span class="usage-used">{usedFormatted} used</span>
			<span class="usage-free">{freeFormatted} free</span>
		</div>
		<div class="usage-total">
			{totalFormatted} total
		</div>
	</div>
</button>

<style>
	.drive-card {
		display: flex;
		flex-direction: column;
		gap: 12px;
		padding: 16px;
		background: #252525;
		border: 1px solid #333;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.15s ease;
		text-align: left;
		width: 100%;
	}

	.drive-card:hover {
		background: #2a2a2a;
		border-color: #444;
	}

	.drive-header {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.drive-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 40px;
		height: 40px;
		background: #333;
		border-radius: 6px;
		color: #888;
	}

	.drive-info {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.drive-name {
		font-size: 14px;
		font-weight: 500;
		color: #e0e0e0;
	}

	.drive-badge {
		font-size: 10px;
		color: #888;
		background: #333;
		padding: 1px 6px;
		border-radius: 3px;
		width: fit-content;
	}

	.drive-stats {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.usage-bar-container {
		height: 6px;
		background: #333;
		border-radius: 3px;
		overflow: hidden;
	}

	.usage-bar {
		height: 100%;
		border-radius: 3px;
		transition: width 0.3s ease;
	}

	.usage-text {
		display: flex;
		justify-content: space-between;
		font-size: 12px;
	}

	.usage-used {
		color: #aaa;
	}

	.usage-free {
		color: #888;
	}

	.usage-total {
		font-size: 11px;
		color: #666;
	}
</style>
