<script lang="ts">
	/**
	 * Settings page - user preferences and account management
	 */
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { settingsStore, type UserSettings } from '$lib/stores/settings';
	import { Settings, LogOut, Eye, ArrowLeft, RotateCcw, Layout, MousePointer } from 'lucide-svelte';
	import { Button, Toggle, Select, Card } from '$lib/components/ui';

	let settings = $state<UserSettings>({ ...$settingsStore });
	const hasChanges = $derived(JSON.stringify(settings) !== JSON.stringify($settingsStore));

	function handleSave() {
		settingsStore.set(settings);
	}

	function handleReset() {
		settingsStore.reset();
		settings = { ...$settingsStore };
	}

	function handleCancel() {
		settings = { ...$settingsStore };
	}

	async function handleLogout() {
		await authStore.logout();
		goto('/login');
	}

	function goBack() {
		goto('/browse');
	}

	const sortByOptions = [
		{ value: 'name', label: 'Name' },
		{ value: 'size', label: 'Size' },
		{ value: 'modTime', label: 'Date modified' },
		{ value: 'type', label: 'Type' },
	];

	const sortDirOptions = [
		{ value: 'asc', label: 'Ascending' },
		{ value: 'desc', label: 'Descending' },
	];

	const viewModeOptions = [
		{ value: 'list', label: 'List' },
		{ value: 'grid', label: 'Grid' },
	];
</script>

<svelte:head>
	<title>Settings - File Manager</title>
</svelte:head>

<div class="min-h-screen bg-surface-primary p-6">
	<div class="max-w-[640px] mx-auto">
		<!-- Header -->
		<header class="flex items-center gap-4 mb-8">
			<Button variant="secondary" size="icon" onclick={goBack}>
				<ArrowLeft size={20} />
			</Button>
			<div class="flex items-center gap-3 text-text-primary">
				<Settings size={24} />
				<h1 class="text-2xl font-semibold m-0">Settings</h1>
			</div>
		</header>

		<!-- Settings sections -->
		<div class="flex flex-col gap-6">
			<!-- File Display -->
			<Card>
				<h2 class="flex items-center gap-2.5 text-sm font-semibold text-text-primary m-0 mb-4 pb-3 border-b border-border-secondary">
					<Eye size={18} />
					File Display
				</h2>

				<div class="flex flex-col gap-4">
					<div class="flex items-center justify-between py-2">
						<div class="flex flex-col gap-0.5">
							<span class="text-sm text-text-primary">Show hidden files</span>
							<span class="text-xs text-text-muted">Display files and folders starting with a dot (.)</span>
						</div>
						<Toggle bind:checked={settings.showHiddenFiles} />
					</div>

					<div class="flex items-center justify-between py-2">
						<div class="flex flex-col gap-0.5">
							<span class="text-sm text-text-primary">Show file extensions</span>
							<span class="text-xs text-text-muted">Display file extensions in the file list</span>
						</div>
						<Toggle bind:checked={settings.showFileExtensions} />
					</div>

					<div class="flex items-center justify-between py-2">
						<div class="flex flex-col gap-0.5">
							<span class="text-sm text-text-primary">Compact mode</span>
							<span class="text-xs text-text-muted">Reduce spacing in file list for more items</span>
						</div>
						<Toggle bind:checked={settings.compactMode} />
					</div>
				</div>
			</Card>

			<!-- Behavior -->
			<Card>
				<h2 class="flex items-center gap-2.5 text-sm font-semibold text-text-primary m-0 mb-4 pb-3 border-b border-border-secondary">
					<MousePointer size={18} />
					Behavior
				</h2>

				<div class="flex flex-col gap-4">
					<div class="flex items-center justify-between py-2">
						<div class="flex flex-col gap-0.5">
							<span class="text-sm text-text-primary">Confirm before delete</span>
							<span class="text-xs text-text-muted">Show confirmation dialog before deleting files</span>
						</div>
						<Toggle bind:checked={settings.confirmDelete} />
					</div>

					<div class="flex items-center justify-between py-2">
						<div class="flex flex-col gap-0.5">
							<span class="text-sm text-text-primary">Preview on single click</span>
							<span class="text-xs text-text-muted">Open file preview with single click instead of double click</span>
						</div>
						<Toggle bind:checked={settings.previewOnSingleClick} />
					</div>
				</div>
			</Card>

			<!-- Default View -->
			<Card>
				<h2 class="flex items-center gap-2.5 text-sm font-semibold text-text-primary m-0 mb-4 pb-3 border-b border-border-secondary">
					<Layout size={18} />
					Default View
				</h2>

				<div class="flex flex-col gap-4">
					<div class="flex items-center justify-between py-2">
						<div class="flex flex-col gap-0.5">
							<span class="text-sm text-text-primary">Default sort by</span>
							<span class="text-xs text-text-muted">How files are sorted when opening a folder</span>
						</div>
						<div class="w-36">
							<Select options={sortByOptions} bind:value={settings.defaultSortBy} />
						</div>
					</div>

					<div class="flex items-center justify-between py-2">
						<div class="flex flex-col gap-0.5">
							<span class="text-sm text-text-primary">Sort direction</span>
							<span class="text-xs text-text-muted">Ascending or descending order</span>
						</div>
						<div class="w-36">
							<Select options={sortDirOptions} bind:value={settings.defaultSortDir} />
						</div>
					</div>

					<div class="flex items-center justify-between py-2">
						<div class="flex flex-col gap-0.5">
							<span class="text-sm text-text-primary">Default view mode</span>
							<span class="text-xs text-text-muted">List or grid view</span>
						</div>
						<div class="w-36">
							<Select options={viewModeOptions} bind:value={settings.defaultViewMode} />
						</div>
					</div>
				</div>
			</Card>

			<!-- Actions -->
			{#if hasChanges}
				<div class="flex justify-end gap-3 p-4 bg-surface-primary border border-border-primary rounded-lg">
					<Button variant="secondary" onclick={handleCancel}>Cancel</Button>
					<Button variant="primary" onclick={handleSave}>Save Changes</Button>
				</div>
			{/if}

			<!-- Account section -->
			<Card>
				<h2 class="flex items-center gap-2.5 text-sm font-semibold text-text-primary m-0 mb-4 pb-3 border-b border-border-secondary">
					<LogOut size={18} />
					Account
				</h2>

				<div class="flex gap-3">
					<Button variant="secondary" onclick={handleReset}>
						<RotateCcw size={16} />
						Reset to Defaults
					</Button>
					<Button variant="danger" onclick={handleLogout}>
						<LogOut size={16} />
						Logout
					</Button>
				</div>
			</Card>
		</div>
	</div>
</div>
