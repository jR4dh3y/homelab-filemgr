/**
 * Settings store - persists user preferences to localStorage
 */
import { writable, derived, get } from 'svelte/store';
import { settingsStorage } from '$lib/utils/storage';

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
	const stored = settingsStorage.get<UserSettings>();
	return stored ? { ...defaultSettings, ...stored } : defaultSettings;
}

function saveSettings(settings: UserSettings): void {
	settingsStorage.set(settings);
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
