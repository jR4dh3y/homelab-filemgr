<script lang="ts">
	import type { SystemDrive } from '$lib/api/system';
	import { HardDrive, Database, Pencil, X } from 'lucide-svelte';
	import { formatFileSize } from '$lib/utils/format';
	import { Badge, ProgressBar, ContextMenu, InlineRename } from '$lib/components/ui';
	import { settingsStore } from '$lib/stores/settings';

	interface Props {
		drive: SystemDrive;
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
	const customName = $derived($settingsStore.driveNameOverrides[drive.mountPoint]);
	const displayName = $derived(customName || drive.mountPoint);

	const progressVariant = $derived.by(() => {
		if (drive.usedPct >= 90) return 'danger' as const;
		if (drive.usedPct >= 75) return 'warning' as const;
		return 'default' as const;
	});

	const isRootDrive = $derived(drive.mountPoint === '/' || drive.mountPoint.match(/^[A-Z]:\\$/));
	const hasCustomName = $derived(!!customName);

	function startRenaming(): void {
		renaming = true;
		contextMenuOpen = false;
	}

	function handleSaveRename(newValue: string): void {
		if (newValue) {
			settingsStore.setDriveName(drive.mountPoint, newValue);
		} else {
			settingsStore.removeDriveName(drive.mountPoint);
		}
		renaming = false;
	}

	function handleCancelRename(): void {
		renaming = false;
	}

	function resetName(): void {
		settingsStore.removeDriveName(drive.mountPoint);
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
	class="flex items-stretch gap-3 p-4 bg-surface-secondary border border-border-primary rounded-lg w-full text-left cursor-pointer transition-all duration-150 hover:bg-surface-tertiary hover:border-border-focus relative"
	onclick={handleCardClick}
	oncontextmenu={handleContextMenu}
>
	<div class="shrink-0 w-22 flex items-center justify-center rounded bg-surface-elevated text-text-secondary">
		{#if isRootDrive}
			<Database size={22} />
		{:else}
			<HardDrive size={22} />
		{/if}
	</div>

	<div class="flex-1 min-w-0 flex flex-col gap-1.5">
		{#if renaming}
			<InlineRename
				value={displayName}
				onSave={handleSaveRename}
				onCancel={handleCancelRename}
			/>
		{:else}
			<div class="flex items-center justify-between gap-2 h-5">
				<span class="text-sm font-medium text-text-primary truncate" title={drive.mountPoint}>
					{displayName}
				</span>
				<Badge variant="default">{totalFormatted}</Badge>
			</div>
		{/if}
		
		<div class="text-xs text-text-muted">
			{usedFormatted} used · {freeFormatted} free
		</div>
		
		<div class="text-[11px] text-text-muted font-mono truncate" title="{drive.device}">
			{drive.device}
			{#if drive.fsType}
				<span class="text-text-secondary"> · {drive.fsType}</span>
			{/if}
		</div>
		
		<div class="flex items-center gap-3 mt-0.5">
			<div class="flex-1">
				<ProgressBar value={drive.usedPct} size="sm" variant={progressVariant} />
			</div>
			<span class="text-[11px] shrink-0 font-medium {drive.usedPct >= 90 ? 'text-danger' : drive.usedPct >= 75 ? 'text-warning' : 'text-text-muted'}">
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
