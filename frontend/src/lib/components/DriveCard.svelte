<script lang="ts">
	import type { DriveStats } from '$lib/api/files';
	import { HardDrive, Pencil, X } from 'lucide-svelte';
	import { formatFileSize } from '$lib/utils/format';
	import { Badge, ProgressBar, ContextMenu, InlineRename } from '$lib/components/ui';
	import { settingsStore } from '$lib/stores/settings';

	interface Props {
		drive: DriveStats;
		onClick?: () => void;
	}

	let { drive, onClick }: Props = $props();

	let renaming = $state(false);
	let contextMenuOpen = $state(false);
	let contextMenuPosition = $state({ x: 0, y: 0 });

	const usedFormatted = $derived(formatFileSize(drive.usedBytes));
	const totalFormatted = $derived(formatFileSize(drive.totalBytes));
	const freeFormatted = $derived(formatFileSize(drive.freeBytes));

	// Subscribe to the store reactively so UI updates when drive names change
	const customName = $derived($settingsStore.driveNameOverrides[drive.name]);
	const displayName = $derived(customName || drive.name);

	const progressVariant = $derived.by(() => {
		if (drive.usedPct >= 90) return 'danger' as const;
		if (drive.usedPct >= 75) return 'warning' as const;
		return 'default' as const;
	});

	// Use the already-derived customName for reactivity
	const hasCustomName = $derived(!!customName);

	function startRenaming(): void {
		renaming = true;
		contextMenuOpen = false;
	}

	function handleSaveRename(newValue: string): void {
		if (newValue) {
			settingsStore.setDriveName(drive.name, newValue);
		} else {
			settingsStore.removeDriveName(drive.name);
		}
		renaming = false;
	}

	function handleCancelRename(): void {
		renaming = false;
	}

	function resetName(): void {
		settingsStore.removeDriveName(drive.name);
		contextMenuOpen = false;
	}

	function handleCardClick(): void {
		if (!renaming) {
			onClick?.();
		}
	}

	function handleContextMenu(e: MouseEvent): void {
		e.preventDefault();
		contextMenuPosition = { x: e.clientX, y: e.clientY };
		contextMenuOpen = true;
	}

	function handleMenuSelect(id: string): void {
		if (id === 'rename') {
			startRenaming();
		} else if (id === 'reset') {
			resetName();
		}
	}

	function handleMenuClose(): void {
		contextMenuOpen = false;
	}

	function handleClickOutside(e: MouseEvent): void {
		if (contextMenuOpen && e.target !== null) {
			contextMenuOpen = false;
		}
	}
</script>

<svelte:window onclick={handleClickOutside} />

<button
	type="button"
	class="flex items-stretch gap-3 p-4 bg-surface-secondary border border-border-primary rounded-lg cursor-pointer transition-all duration-150 text-left w-full hover:bg-surface-tertiary hover:border-border-focus relative"
	onclick={handleCardClick}
	oncontextmenu={handleContextMenu}
>
	<div class="shrink-0 w-16 flex items-center justify-center rounded bg-surface-elevated text-text-secondary">
		<HardDrive size={24} />
	</div>

	<div class="flex-1 min-w-0 flex flex-col gap-1 py-0.5">
		{#if renaming}
			<InlineRename
				value={displayName}
				onSave={handleSaveRename}
				onCancel={handleCancelRename}
			/>
		{:else}
			<div class="flex items-center justify-between gap-2 h-5">
				<span class="text-sm font-medium text-text-primary">{displayName}</span>
				<Badge variant="default">{totalFormatted}</Badge>
			</div>
		{/if}
		<div class="text-xs text-text-muted">
			{usedFormatted} used · {freeFormatted} free
			{#if drive.readOnly}
				<span class="text-warning ml-1">· Read-only</span>
			{/if}
		</div>
		{#if drive.device}
			<div class="text-[10px] text-text-muted font-mono truncate" title="{drive.device} ({drive.fsType || 'unknown'})">
				{drive.device}{#if drive.fsType} · {drive.fsType}{/if}
			</div>
		{/if}

		<div class="flex items-center gap-3">
			<div class="flex-1">
				<ProgressBar value={drive.usedPct} size="sm" variant={progressVariant} />
			</div>
			<span class="text-[11px] shrink-0 {drive.usedPct >= 90 ? 'text-danger' : drive.usedPct >= 75 ? 'text-warning' : 'text-text-muted'}">
				{drive.usedPct.toFixed(1)}%
			</span>
		</div>
	</div>

	{#if contextMenuOpen && !renaming}
		<ContextMenu
			x={contextMenuPosition.x}
			y={contextMenuPosition.y}
			items={[
				{ id: 'rename', label: 'Rename', icon: Pencil },
				...(hasCustomName ? [{ id: 'reset', label: 'Reset Name', icon: X }] : [])
			]}
			onSelect={handleMenuSelect}
			onClose={handleMenuClose}
		/>
	{/if}
</button>
