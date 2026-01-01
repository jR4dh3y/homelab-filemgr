<script lang="ts">
	/**
	 * Settings page - user preferences and account management
	 */
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { settingsStore, type UserSettings } from '$lib/stores/settings';
	import {
		Settings,
		LogOut,
		Eye,
		EyeOff,
		FileText,
		Trash2,
		ArrowLeft,
		RotateCcw,
		Layout,
		MousePointer
	} from 'lucide-svelte';

	// Local copy of settings for the form
	let settings = $state<UserSettings>({ ...$settingsStore });

	// Track if settings have changed
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
</script>

<svelte:head>
	<title>Settings - File Manager</title>
</svelte:head>

<div class="settings-page">
	<div class="settings-container">
		<!-- Header -->
		<header class="settings-header">
			<button type="button" class="back-btn" onclick={goBack}>
				<ArrowLeft size={20} />
			</button>
			<div class="header-title">
				<Settings size={24} />
				<h1>Settings</h1>
			</div>
		</header>

		<!-- Settings sections -->
		<div class="settings-content">
			<!-- File Display -->
			<section class="settings-section">
				<h2 class="section-title">
					<Eye size={18} />
					File Display
				</h2>
				
				<div class="setting-item">
					<div class="setting-info">
						<label for="showHidden">Show hidden files</label>
						<span class="setting-desc">Display files and folders starting with a dot (.)</span>
					</div>
					<label class="toggle">
						<input 
							type="checkbox" 
							id="showHidden"
							bind:checked={settings.showHiddenFiles}
						/>
						<span class="toggle-slider"></span>
					</label>
				</div>

				<div class="setting-item">
					<div class="setting-info">
						<label for="showExtensions">Show file extensions</label>
						<span class="setting-desc">Display file extensions in the file list</span>
					</div>
					<label class="toggle">
						<input 
							type="checkbox" 
							id="showExtensions"
							bind:checked={settings.showFileExtensions}
						/>
						<span class="toggle-slider"></span>
					</label>
				</div>

				<div class="setting-item">
					<div class="setting-info">
						<label for="compactMode">Compact mode</label>
						<span class="setting-desc">Reduce spacing in file list for more items</span>
					</div>
					<label class="toggle">
						<input 
							type="checkbox" 
							id="compactMode"
							bind:checked={settings.compactMode}
						/>
						<span class="toggle-slider"></span>
					</label>
				</div>
			</section>

			<!-- Behavior -->
			<section class="settings-section">
				<h2 class="section-title">
					<MousePointer size={18} />
					Behavior
				</h2>

				<div class="setting-item">
					<div class="setting-info">
						<label for="confirmDelete">Confirm before delete</label>
						<span class="setting-desc">Show confirmation dialog before deleting files</span>
					</div>
					<label class="toggle">
						<input 
							type="checkbox" 
							id="confirmDelete"
							bind:checked={settings.confirmDelete}
						/>
						<span class="toggle-slider"></span>
					</label>
				</div>

				<div class="setting-item">
					<div class="setting-info">
						<label for="previewOnClick">Preview on single click</label>
						<span class="setting-desc">Open file preview with single click instead of double click</span>
					</div>
					<label class="toggle">
						<input 
							type="checkbox" 
							id="previewOnClick"
							bind:checked={settings.previewOnSingleClick}
						/>
						<span class="toggle-slider"></span>
					</label>
				</div>
			</section>

			<!-- Default View -->
			<section class="settings-section">
				<h2 class="section-title">
					<Layout size={18} />
					Default View
				</h2>

				<div class="setting-item">
					<div class="setting-info">
						<label for="defaultSort">Default sort by</label>
						<span class="setting-desc">How files are sorted when opening a folder</span>
					</div>
					<select id="defaultSort" bind:value={settings.defaultSortBy}>
						<option value="name">Name</option>
						<option value="size">Size</option>
						<option value="modTime">Date modified</option>
						<option value="type">Type</option>
					</select>
				</div>

				<div class="setting-item">
					<div class="setting-info">
						<label for="defaultSortDir">Sort direction</label>
						<span class="setting-desc">Ascending or descending order</span>
					</div>
					<select id="defaultSortDir" bind:value={settings.defaultSortDir}>
						<option value="asc">Ascending</option>
						<option value="desc">Descending</option>
					</select>
				</div>

				<div class="setting-item">
					<div class="setting-info">
						<label for="defaultView">Default view mode</label>
						<span class="setting-desc">List or grid view</span>
					</div>
					<select id="defaultView" bind:value={settings.defaultViewMode}>
						<option value="list">List</option>
						<option value="grid">Grid</option>
					</select>
				</div>
			</section>

			<!-- Actions -->
			{#if hasChanges}
				<div class="settings-actions">
					<button type="button" class="btn btn-secondary" onclick={handleCancel}>
						Cancel
					</button>
					<button type="button" class="btn btn-primary" onclick={handleSave}>
						Save Changes
					</button>
				</div>
			{/if}

			<!-- Account section -->
			<section class="settings-section account-section">
				<h2 class="section-title">
					<LogOut size={18} />
					Account
				</h2>

				<div class="account-actions">
					<button type="button" class="btn btn-outline" onclick={handleReset}>
						<RotateCcw size={16} />
						Reset to Defaults
					</button>
					<button type="button" class="btn btn-danger" onclick={handleLogout}>
						<LogOut size={16} />
						Logout
					</button>
				</div>
			</section>
		</div>
	</div>
</div>

<style>
	.settings-page {
		min-height: 100vh;
		background: #141414;
		padding: 24px;
	}

	.settings-container {
		max-width: 640px;
		margin: 0 auto;
	}

	.settings-header {
		display: flex;
		align-items: center;
		gap: 16px;
		margin-bottom: 32px;
	}

	.back-btn {
		width: 40px;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #252525;
		border: 1px solid #333;
		border-radius: 8px;
		color: #888;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.back-btn:hover {
		background: #2a2a2a;
		color: #ccc;
	}

	.header-title {
		display: flex;
		align-items: center;
		gap: 12px;
		color: #e0e0e0;
	}

	.header-title h1 {
		font-size: 24px;
		font-weight: 600;
		margin: 0;
	}

	.settings-content {
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	.settings-section {
		background: #1e1e1e;
		border: 1px solid #2a2a2a;
		border-radius: 8px;
		padding: 20px;
	}

	.section-title {
		display: flex;
		align-items: center;
		gap: 10px;
		font-size: 14px;
		font-weight: 600;
		color: #e0e0e0;
		margin: 0 0 16px 0;
		padding-bottom: 12px;
		border-bottom: 1px solid #2a2a2a;
	}

	.setting-item {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 12px 0;
		border-bottom: 1px solid #252525;
	}

	.setting-item:last-child {
		border-bottom: none;
		padding-bottom: 0;
	}

	.setting-info {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.setting-info label {
		font-size: 14px;
		color: #ccc;
		cursor: pointer;
	}

	.setting-desc {
		font-size: 12px;
		color: #666;
	}

	/* Toggle switch */
	.toggle {
		position: relative;
		display: inline-block;
		width: 44px;
		height: 24px;
		flex-shrink: 0;
	}

	.toggle input {
		opacity: 0;
		width: 0;
		height: 0;
	}

	.toggle-slider {
		position: absolute;
		cursor: pointer;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: #333;
		border-radius: 24px;
		transition: 0.2s;
	}

	.toggle-slider:before {
		position: absolute;
		content: "";
		height: 18px;
		width: 18px;
		left: 3px;
		bottom: 3px;
		background: #888;
		border-radius: 50%;
		transition: 0.2s;
	}

	.toggle input:checked + .toggle-slider {
		background: #2d4a6f;
	}

	.toggle input:checked + .toggle-slider:before {
		transform: translateX(20px);
		background: #fff;
	}

	/* Select dropdown */
	select {
		padding: 8px 12px;
		background: #252525;
		border: 1px solid #333;
		border-radius: 6px;
		color: #ccc;
		font-size: 13px;
		cursor: pointer;
		min-width: 140px;
	}

	select:hover {
		border-color: #444;
	}

	select:focus {
		outline: none;
		border-color: #4a9eff;
	}

	/* Buttons */
	.btn {
		display: inline-flex;
		align-items: center;
		gap: 8px;
		padding: 10px 16px;
		font-size: 14px;
		font-weight: 500;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.15s ease;
		border: none;
	}

	.btn-primary {
		background: #2d4a6f;
		color: #fff;
	}

	.btn-primary:hover {
		background: #345580;
	}

	.btn-secondary {
		background: #333;
		color: #ccc;
	}

	.btn-secondary:hover {
		background: #3a3a3a;
	}

	.btn-outline {
		background: transparent;
		border: 1px solid #333;
		color: #888;
	}

	.btn-outline:hover {
		background: #252525;
		color: #ccc;
	}

	.btn-danger {
		background: #dc3545;
		color: #fff;
	}

	.btn-danger:hover {
		background: #c82333;
	}

	.settings-actions {
		display: flex;
		justify-content: flex-end;
		gap: 12px;
		padding: 16px 20px;
		background: #1a1a1a;
		border: 1px solid #2a2a2a;
		border-radius: 8px;
	}

	.account-section {
		margin-top: 16px;
	}

	.account-actions {
		display: flex;
		gap: 12px;
	}
</style>
