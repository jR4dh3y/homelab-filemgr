<script lang="ts">
	/**
	 * Settings page - user preferences and account management
	 * Design follows the same visual language as the file browser sidebar
	 */
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { settingsStore, type UserSettings } from '$lib/stores/settings';
	import {
		ChevronDown,
		ChevronLeft,
		Eye,
		RotateCcw,
		Layout,
		MousePointer,
		User,
		LogOut,
		Save,
		X,
	} from 'lucide-svelte';
	import { Button, Toggle, Select } from '$lib/components/ui';

	let settings = $state<UserSettings>({ ...$settingsStore });
	const hasChanges = $derived(JSON.stringify(settings) !== JSON.stringify($settingsStore));

	// Section collapse state
	let displayCollapsed = $state(false);
	let behaviorCollapsed = $state(false);
	let viewCollapsed = $state(false);
	let accountCollapsed = $state(false);

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

	// Shared styles matching the sidebar design
	const sectionHeaderClass =
		'w-full flex items-center gap-1.5 px-3 py-2.5 bg-transparent border-none text-text-secondary text-[11px] font-medium uppercase tracking-wide cursor-pointer text-left hover:text-text-primary';
	const settingRowClass = 'flex items-center justify-between py-2 px-3';
</script>

<svelte:head>
	<title>Settings - File Manager</title>
</svelte:head>

<div class="flex h-screen w-full bg-surface-primary overflow-hidden">
	<!-- Settings Sidebar (matching main sidebar width) -->
	<aside class="w-[220px] min-w-[220px] bg-surface-primary border-r border-border-secondary flex flex-col overflow-y-auto overflow-x-hidden">
		<!-- Header -->
		<div class="flex items-center gap-2 px-3 py-3 border-b border-border-secondary">
			<button
				type="button"
				class="w-7 h-7 flex items-center justify-center bg-transparent border-none rounded text-text-secondary cursor-pointer transition-all duration-100 hover:bg-surface-elevated hover:text-text-primary"
				onclick={goBack}
				title="Back to Files"
			>
				<ChevronLeft size={18} />
			</button>
			<span class="text-text-primary text-[13px] font-medium">Settings</span>
		</div>

		<!-- Navigation within settings -->
		<nav class="flex-1 py-2">
			<button
				type="button"
				class="w-full flex items-center gap-2.5 py-1.5 px-3 pl-5 bg-transparent border-none text-text-primary text-[13px] cursor-pointer text-left transition-colors duration-100 hover:bg-surface-secondary bg-selection"
			>
				<Eye size={16} class="shrink-0 opacity-80" />
				<span>Preferences</span>
			</button>
			<button
				type="button"
				class="w-full flex items-center gap-2.5 py-1.5 px-3 pl-5 bg-transparent border-none text-text-secondary text-[13px] cursor-pointer text-left transition-colors duration-100 hover:bg-surface-secondary"
				onclick={handleLogout}
			>
				<LogOut size={16} class="shrink-0 opacity-80" />
				<span>Logout</span>
			</button>
		</nav>
	</aside>

	<!-- Main content area -->
	<div class="flex-1 flex flex-col min-w-0">
		<!-- Toolbar matching browse toolbar -->
		<div class="flex items-center gap-2 px-3 py-1.5 bg-surface-primary border-b border-border-secondary">
			<div class="flex items-center gap-2">
				<button
					type="button"
					class="w-7 h-7 flex items-center justify-center bg-transparent border-none rounded text-text-secondary cursor-pointer transition-all duration-100 hover:bg-surface-elevated hover:text-text-primary"
					onclick={goBack}
					title="Back"
				>
					<ChevronLeft size={18} />
				</button>
			</div>

			<!-- Path bar style breadcrumb -->
			<div class="flex-1 flex items-center gap-1.5 bg-surface-secondary border border-border-primary rounded px-2 py-1 min-w-0">
				<span class="text-text-secondary text-[13px]">Settings</span>
				<span class="text-text-muted text-xs">/</span>
				<span class="text-text-primary text-[13px]">Preferences</span>
			</div>

			<!-- Action buttons -->
			{#if hasChanges}
				<div class="flex gap-1">
					<Button variant="ghost" size="sm" onclick={handleCancel} title="Discard changes">
						<X size={16} />
						<span class="hidden sm:inline">Cancel</span>
					</Button>
					<Button variant="primary" size="sm" onclick={handleSave} title="Save changes">
						<Save size={16} />
						<span class="hidden sm:inline">Save</span>
					</Button>
				</div>
			{/if}
		</div>

		<!-- Settings content (styled like file list area) -->
		<div class="flex-1 overflow-auto">
			<div class="max-w-[600px]">
				<!-- File Display Section -->
				<div class="border-b border-border-secondary">
					<button
						type="button"
						class={sectionHeaderClass}
						onclick={() => (displayCollapsed = !displayCollapsed)}
					>
						<ChevronDown
							size={14}
							class="shrink-0 transition-transform duration-150 {displayCollapsed ? '-rotate-90' : ''}"
						/>
						<Eye size={14} class="shrink-0 opacity-60" />
						<span>File Display</span>
					</button>
					{#if !displayCollapsed}
						<div class="pb-3">
							<div class={settingRowClass}>
								<div class="flex flex-col gap-0.5">
									<span class="text-[13px] text-text-primary">Show hidden files</span>
									<span class="text-[11px] text-text-muted">Display files starting with a dot</span>
								</div>
								<Toggle bind:checked={settings.showHiddenFiles} />
							</div>

							<div class={settingRowClass}>
								<div class="flex flex-col gap-0.5">
									<span class="text-[13px] text-text-primary">Show file extensions</span>
									<span class="text-[11px] text-text-muted">Display extensions in file list</span>
								</div>
								<Toggle bind:checked={settings.showFileExtensions} />
							</div>

							<div class={settingRowClass}>
								<div class="flex flex-col gap-0.5">
									<span class="text-[13px] text-text-primary">Compact mode</span>
									<span class="text-[11px] text-text-muted">Reduce spacing for more items</span>
								</div>
								<Toggle bind:checked={settings.compactMode} />
							</div>
						</div>
					{/if}
				</div>

				<!-- Behavior Section -->
				<div class="border-b border-border-secondary">
					<button
						type="button"
						class={sectionHeaderClass}
						onclick={() => (behaviorCollapsed = !behaviorCollapsed)}
					>
						<ChevronDown
							size={14}
							class="shrink-0 transition-transform duration-150 {behaviorCollapsed ? '-rotate-90' : ''}"
						/>
						<MousePointer size={14} class="shrink-0 opacity-60" />
						<span>Behavior</span>
					</button>
					{#if !behaviorCollapsed}
						<div class="pb-3">
							<div class={settingRowClass}>
								<div class="flex flex-col gap-0.5">
									<span class="text-[13px] text-text-primary">Confirm before delete</span>
									<span class="text-[11px] text-text-muted">Show confirmation dialog</span>
								</div>
								<Toggle bind:checked={settings.confirmDelete} />
							</div>

							<div class={settingRowClass}>
								<div class="flex flex-col gap-0.5">
									<span class="text-[13px] text-text-primary">Preview on single click</span>
									<span class="text-[11px] text-text-muted">Open preview with single click</span>
								</div>
								<Toggle bind:checked={settings.previewOnSingleClick} />
							</div>
						</div>
					{/if}
				</div>

				<!-- Default View Section -->
				<div class="border-b border-border-secondary">
					<button
						type="button"
						class={sectionHeaderClass}
						onclick={() => (viewCollapsed = !viewCollapsed)}
					>
						<ChevronDown
							size={14}
							class="shrink-0 transition-transform duration-150 {viewCollapsed ? '-rotate-90' : ''}"
						/>
						<Layout size={14} class="shrink-0 opacity-60" />
						<span>Default View</span>
					</button>
					{#if !viewCollapsed}
						<div class="pb-3">
							<div class={settingRowClass}>
								<div class="flex flex-col gap-0.5">
									<span class="text-[13px] text-text-primary">Sort by</span>
									<span class="text-[11px] text-text-muted">Default sort field</span>
								</div>
								<div class="w-32">
									<Select options={sortByOptions} bind:value={settings.defaultSortBy} />
								</div>
							</div>

							<div class={settingRowClass}>
								<div class="flex flex-col gap-0.5">
									<span class="text-[13px] text-text-primary">Sort direction</span>
									<span class="text-[11px] text-text-muted">Ascending or descending</span>
								</div>
								<div class="w-32">
									<Select options={sortDirOptions} bind:value={settings.defaultSortDir} />
								</div>
							</div>

							<div class={settingRowClass}>
								<div class="flex flex-col gap-0.5">
									<span class="text-[13px] text-text-primary">View mode</span>
									<span class="text-[11px] text-text-muted">List or grid view</span>
								</div>
								<div class="w-32">
									<Select options={viewModeOptions} bind:value={settings.defaultViewMode} />
								</div>
							</div>
						</div>
					{/if}
				</div>

				<!-- Account Section -->
				<div class="border-b border-border-secondary">
					<button
						type="button"
						class={sectionHeaderClass}
						onclick={() => (accountCollapsed = !accountCollapsed)}
					>
						<ChevronDown
							size={14}
							class="shrink-0 transition-transform duration-150 {accountCollapsed ? '-rotate-90' : ''}"
						/>
						<User size={14} class="shrink-0 opacity-60" />
						<span>Account</span>
					</button>
					{#if !accountCollapsed}
						<div class="pb-3 px-3">
							<div class="flex gap-2 pt-2">
								<Button variant="secondary" size="sm" onclick={handleReset}>
									<RotateCcw size={14} />
									Reset to Defaults
								</Button>
								<Button variant="danger" size="sm" onclick={handleLogout}>
									<LogOut size={14} />
									Logout
								</Button>
							</div>
						</div>
					{/if}
				</div>
			</div>
		</div>

		<!-- Status bar matching browse status bar -->
		<div class="flex items-center justify-between px-3 py-1.5 bg-surface-primary border-t border-border-secondary text-[11px] text-text-muted">
			<span>
				{#if hasChanges}
					<span class="text-warning">● Unsaved changes</span>
				{:else}
					All changes saved
				{/if}
			</span>
			<span>Preferences</span>
		</div>
	</div>
</div>
