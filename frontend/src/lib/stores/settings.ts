/**
 * Settings store - persists user preferences to localStorage
 */
import { writable, derived, get } from 'svelte/store';

export interface UserSettings {
	showHiddenFiles: boolean;
	showFileExtensions: boolean;
	confirmDelete: boolean;
	defaultSortBy: 'name' | 'size' | 'modTime' | 'type';
	defaultSortDir: 'asc' | 'desc';
	defaultViewMode: 'list' | 'grid';
	previewOnSingleClick: boolean;
	compactMode: boolean;
}

const SETTINGS_KEY = 'filemanager_settings';

const defaultSettings: UserSettings = {
	showHiddenFiles: false,
	showFileExtensions: true,
	confirmDelete: true,
	defaultSortBy: 'name',
	defaultSortDir: 'asc',
	defaultViewMode: 'list',
	previewOnSingleClick: false,
	compactMode: false
};

function loadSettings(): UserSettings {
	if (typeof window === 'undefined') return defaultSettings;
	
	try {
		const stored = localStorage.getItem(SETTINGS_KEY);
		if (stored) {
			return { ...defaultSettings, ...JSON.parse(stored) };
		}
	} catch (e) {
		console.error('Failed to load settings:', e);
	}
	return defaultSettings;
}

function saveSettings(settings: UserSettings): void {
	if (typeof window === 'undefined') return;
	
	try {
		localStorage.setItem(SETTINGS_KEY, JSON.stringify(settings));
	} catch (e) {
		console.error('Failed to save settings:', e);
	}
}

function createSettingsStore() {
	const { subscribe, set, update } = writable<UserSettings>(loadSettings());

	return {
		subscribe,
		
		set(settings: UserSettings) {
			saveSettings(settings);
			set(settings);
		},

		update(updater: (settings: UserSettings) => UserSettings) {
			update(current => {
				const updated = updater(current);
				saveSettings(updated);
				return updated;
			});
		},

		reset() {
			saveSettings(defaultSettings);
			set(defaultSettings);
		},

		setSetting<K extends keyof UserSettings>(key: K, value: UserSettings[K]) {
			update(current => {
				const updated = { ...current, [key]: value };
				saveSettings(updated);
				return updated;
			});
		},

		getSetting<K extends keyof UserSettings>(key: K): UserSettings[K] {
			return get({ subscribe })[key];
		}
	};
}

export const settingsStore = createSettingsStore();

// Derived stores for individual settings
export const showHiddenFiles = derived(settingsStore, $s => $s.showHiddenFiles);
export const showFileExtensions = derived(settingsStore, $s => $s.showFileExtensions);
export const confirmDelete = derived(settingsStore, $s => $s.confirmDelete);
export const compactMode = derived(settingsStore, $s => $s.compactMode);
