<script lang="ts">
	/**
	 * SystemDriveCard component - displays system drive/filesystem with usage stats
	 */
	import type { SystemDrive } from '$lib/api/system';
	import { HardDrive, Database } from 'lucide-svelte';
	import { formatFileSize } from '$lib/utils/format';
	import { Badge, ProgressBar } from '$lib/components/ui';

	interface Props {
		drive: SystemDrive;
	}

	let { drive }: Props = $props();

	const usedFormatted = $derived(formatFileSize(drive.usedBytes));
	const totalFormatted = $derived(formatFileSize(drive.totalBytes));
	const freeFormatted = $derived(formatFileSize(drive.freeBytes));

	// Variant based on usage percentage
	const progressVariant = $derived.by(() => {
		if (drive.usedPct >= 90) return 'danger' as const;
		if (drive.usedPct >= 75) return 'warning' as const;
		return 'default' as const;
	});

	// Determine icon based on mount point or device
	const isRootDrive = $derived(drive.mountPoint === '/' || drive.mountPoint.match(/^[A-Z]:\\$/));
</script>

<div
	class="flex items-stretch gap-3 p-4 bg-surface-secondary border border-border-primary rounded-lg w-full"
>
	<!-- Icon -->
	<div class="shrink-0 w-14 flex items-center justify-center rounded bg-surface-elevated text-text-secondary">
		{#if isRootDrive}
			<Database size={22} />
		{:else}
			<HardDrive size={22} />
		{/if}
	</div>

	<!-- Content -->
	<div class="flex-1 min-w-0 flex flex-col gap-1.5">
		<div class="flex items-center justify-between gap-2">
			<span class="text-sm font-medium text-text-primary truncate" title={drive.mountPoint}>
				{drive.mountPoint}
			</span>
			<Badge variant="default">{totalFormatted}</Badge>
		</div>
		
		<div class="text-xs text-text-muted">
			{usedFormatted} used · {freeFormatted} free
		</div>
		
		<div class="text-[11px] text-text-muted font-mono truncate" title="{drive.device}">
			{drive.device}
			{#if drive.fsType}
				<span class="text-text-secondary"> · {drive.fsType}</span>
			{/if}
		</div>
		
		<!-- Progress bar with percentage -->
		<div class="flex items-center gap-3 mt-0.5">
			<div class="flex-1">
				<ProgressBar value={drive.usedPct} size="sm" variant={progressVariant} />
			</div>
			<span class="text-[11px] shrink-0 font-medium {drive.usedPct >= 90 ? 'text-danger' : drive.usedPct >= 75 ? 'text-warning' : 'text-text-muted'}">
				{drive.usedPct.toFixed(1)}%
			</span>
		</div>
	</div>
</div>
