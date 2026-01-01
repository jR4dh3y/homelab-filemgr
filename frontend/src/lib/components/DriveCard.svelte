<script lang="ts">
	/**
	 * DriveCard component - displays drive/mount point with usage stats
	 */
	import type { DriveStats } from '$lib/api/files';
	import { HardDrive } from 'lucide-svelte';
	import { formatFileSize } from '$lib/utils/format';
	import { Badge } from '$lib/components/ui';

	interface Props {
		drive: DriveStats;
		onClick?: () => void;
	}

	let { drive, onClick }: Props = $props();

	const usedFormatted = $derived(formatFileSize(drive.usedBytes));
	const totalFormatted = $derived(formatFileSize(drive.totalBytes));
	const freeFormatted = $derived(formatFileSize(drive.freeBytes));

	// Color based on usage percentage
	const barColorClass = $derived.by(() => {
		if (drive.usedPct >= 90) return 'bg-danger';
		if (drive.usedPct >= 75) return 'bg-warning';
		return 'bg-accent';
	});
</script>

<button
	type="button"
	class="flex flex-col gap-3 p-4 bg-surface-secondary border border-border-primary rounded-md cursor-pointer transition-all duration-150 text-left w-full hover:bg-surface-tertiary hover:border-border-focus"
	onclick={onClick}
>
	<div class="flex items-center gap-3">
		<div class="flex items-center justify-center w-10 h-10 bg-surface-elevated rounded-md text-text-secondary">
			<HardDrive size={24} />
		</div>
		<div class="flex flex-col gap-0.5">
			<span class="text-sm font-medium text-text-primary">{drive.name}</span>
			{#if drive.readOnly}
				<Badge>Read-only</Badge>
			{/if}
		</div>
	</div>

	<div class="flex flex-col gap-1.5">
		<div class="h-1.5 bg-surface-elevated rounded-full overflow-hidden">
			<div class="{barColorClass} h-full rounded-full transition-all duration-300" style="width: {drive.usedPct}%"></div>
		</div>
		<div class="flex justify-between text-xs">
			<span class="text-text-secondary">{usedFormatted} used</span>
			<span class="text-text-muted">{freeFormatted} free</span>
		</div>
		<div class="text-[11px] text-text-muted">
			{totalFormatted} total
		</div>
	</div>
</button>
