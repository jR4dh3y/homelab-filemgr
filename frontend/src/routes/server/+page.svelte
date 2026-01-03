<script lang="ts">
	/**
	 * Server page - displays system information including all attached drives
	 */
	import { createQuery } from '@tanstack/svelte-query';
	import { goto } from '$app/navigation';
	import { getSystemDrives, type SystemDrivesResponse, type SystemDrive } from '$lib/api/system';
	import { formatFileSize } from '$lib/utils/format';
	import { Button, Card, Spinner } from '$lib/components/ui';
	import SystemDriveCard from '$lib/components/SystemDriveCard.svelte';
	import { Server, ArrowLeft, HardDrive, RefreshCw, AlertCircle } from 'lucide-svelte';

	// Query for system drives
	const drivesQuery = createQuery<SystemDrivesResponse>(() => ({
		queryKey: ['system', 'drives'],
		queryFn: () => getSystemDrives(),
		refetchInterval: 30000, // Refresh every 30 seconds
	}));

	const drives = $derived(drivesQuery.data?.drives ?? []);
	const isLoading = $derived(drivesQuery.isLoading);
	const isError = $derived(drivesQuery.isError);
	const isRefetching = $derived(drivesQuery.isFetching && !drivesQuery.isLoading);

	// Calculate totals
	const totalStorage = $derived(drives.reduce((sum: number, d: SystemDrive) => sum + d.totalBytes, 0));
	const totalUsed = $derived(drives.reduce((sum: number, d: SystemDrive) => sum + d.usedBytes, 0));
	const totalFree = $derived(drives.reduce((sum: number, d: SystemDrive) => sum + d.freeBytes, 0));
	const overallUsedPct = $derived(totalStorage > 0 ? (totalUsed / totalStorage) * 100 : 0);

	function goBack() {
		goto('/browse');
	}

	function handleRefresh() {
		drivesQuery.refetch();
	}
</script>

<svelte:head>
	<title>Server - File Manager</title>
</svelte:head>

<div class="min-h-screen bg-surface-primary p-6">
	<div class="max-w-4xl mx-auto">
		<!-- Header -->
		<header class="flex items-center justify-between mb-8">
			<div class="flex items-center gap-4">
				<Button variant="secondary" size="icon" onclick={goBack}>
					<ArrowLeft size={20} />
				</Button>
				<div class="flex items-center gap-3 text-text-primary">
					<Server size={24} />
					<h1 class="text-2xl font-semibold m-0">Server</h1>
				</div>
			</div>
			<Button variant="secondary" onclick={handleRefresh} disabled={isLoading || isRefetching}>
				<RefreshCw size={16} class={isRefetching ? 'animate-spin' : ''} />
				Refresh
			</Button>
		</header>

		<!-- Summary Card -->
		<div class="mb-6">
			<Card>
				<div class="flex items-center gap-3 mb-4">
					<HardDrive size={20} class="text-text-secondary" />
					<h2 class="text-lg font-semibold text-text-primary m-0">Storage Overview</h2>
				</div>
			
			{#if isLoading}
				<div class="flex items-center justify-center py-8">
					<Spinner size="md" />
				</div>
			{:else if isError}
				<div class="flex items-center gap-2 text-danger py-4">
					<AlertCircle size={18} />
					<span>Failed to load system drives</span>
				</div>
			{:else}
				<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
					<div class="text-center p-3 bg-surface-elevated rounded-lg">
						<div class="text-2xl font-bold text-text-primary">{drives.length}</div>
						<div class="text-xs text-text-muted">Filesystems</div>
					</div>
					<div class="text-center p-3 bg-surface-elevated rounded-lg">
						<div class="text-2xl font-bold text-text-primary">{formatFileSize(totalStorage)}</div>
						<div class="text-xs text-text-muted">Total Storage</div>
					</div>
					<div class="text-center p-3 bg-surface-elevated rounded-lg">
						<div class="text-2xl font-bold text-text-primary">{formatFileSize(totalUsed)}</div>
						<div class="text-xs text-text-muted">Used</div>
					</div>
					<div class="text-center p-3 bg-surface-elevated rounded-lg">
						<div class="text-2xl font-bold text-text-primary">{formatFileSize(totalFree)}</div>
						<div class="text-xs text-text-muted">Free</div>
					</div>
				</div>
				
				{#if drives.length > 0}
					<div class="mt-4 flex items-center gap-3">
						<div class="flex-1 h-3 bg-surface-elevated rounded-full overflow-hidden">
							<div
								class="h-full transition-all duration-300 {overallUsedPct >= 90 ? 'bg-danger' : overallUsedPct >= 75 ? 'bg-warning' : 'bg-accent-primary'}"
								style="width: {overallUsedPct}%"
							></div>
						</div>
						<span class="text-sm font-medium text-text-secondary shrink-0">
							{overallUsedPct.toFixed(1)}% used
						</span>
					</div>
				{/if}
			{/if}
			</Card>
		</div>

		<!-- Drives List -->
		<div class="mb-4">
			<h2 class="text-sm font-semibold text-text-secondary uppercase tracking-wide mb-3">
				All Filesystems
			</h2>
		</div>

		{#if isLoading}
			<div class="flex items-center justify-center py-12">
				<Spinner size="lg" />
			</div>
		{:else if isError}
			<Card>
				<div class="flex flex-col items-center justify-center py-8 gap-3">
					<AlertCircle size={32} class="text-danger" />
					<p class="text-text-muted">Failed to load system drives. Please try again.</p>
					<Button variant="secondary" onclick={handleRefresh}>
						<RefreshCw size={16} />
						Retry
					</Button>
				</div>
			</Card>
		{:else if drives.length === 0}
			<Card>
				<div class="flex flex-col items-center justify-center py-8 gap-2">
					<HardDrive size={32} class="text-text-muted" />
					<p class="text-text-muted">No filesystems found</p>
				</div>
			</Card>
		{:else}
			<div class="grid gap-3">
				{#each drives as drive (drive.device + drive.mountPoint)}
					<SystemDriveCard {drive} />
				{/each}
			</div>
		{/if}
	</div>
</div>
