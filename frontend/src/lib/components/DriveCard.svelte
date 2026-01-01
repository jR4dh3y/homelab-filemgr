<script lang="ts">
	/**
	 * DriveCard component - displays drive/mount point with usage stats
	 */
	import type { DriveStats } from '$lib/api/files';
	import { HardDrive } from 'lucide-svelte';
	import { formatFileSize } from '$lib/utils/format';
	import { Badge, ProgressBar } from '$lib/components/ui';

	interface Props {
		drive: DriveStats;
		onClick?: () => void;
	}

	let { drive, onClick }: Props = $props();

	const usedFormatted = $derived(formatFileSize(drive.usedBytes));
	const totalFormatted = $derived(formatFileSize(drive.totalBytes));
	const freeFormatted = $derived(formatFileSize(drive.freeBytes));

	// Variant based on usage percentage
	const progressVariant = $derived.by(() => {
		if (drive.usedPct >= 90) return 'danger' as const;
		if (drive.usedPct >= 75) return 'warning' as const;
		return 'default' as const;
	});
</script>

<button
	type="button"
	class="flex items-stretch gap-3 p-4 bg-surface-secondary border border-border-primary rounded-lg cursor-pointer transition-all duration-150 text-left w-full hover:bg-surface-tertiary hover:border-border-focus"
	onclick={onClick}
>
	<!-- Icon -->
	<div class="shrink-0 w-16 flex items-center justify-center rounded bg-surface-elevated text-text-secondary">
		<HardDrive size={24} />
	</div>

	<!-- Content -->
	<div class="flex-1 min-w-0 flex flex-col gap-1 py-0.5">
		<div class="flex items-center justify-between gap-2">
			<span class="text-sm font-medium text-text-primary">{drive.name}</span>
			<Badge variant="default">{totalFormatted}</Badge>
		</div>
		<div class="text-xs text-text-muted">
			{usedFormatted} used · {freeFormatted} free
			{#if drive.readOnly}
				<span class="text-warning ml-1">· Read-only</span>
			{/if}
		</div>
		
		<!-- Progress bar with percentage -->
		<div class="flex items-center gap-3">
			<div class="flex-1">
				<ProgressBar value={drive.usedPct} size="sm" variant={progressVariant} />
			</div>
			<span class="text-[11px] shrink-0 {drive.usedPct >= 90 ? 'text-danger' : drive.usedPct >= 75 ? 'text-warning' : 'text-text-muted'}">
				{drive.usedPct.toFixed(1)}%
			</span>
		</div>
	</div>
</button>
