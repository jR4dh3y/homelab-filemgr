/**
 * Auth store for managing authentication state
 * Requirements: 7.1, 7.4
 */

import { writable, derived, get } from 'svelte/store';
import {
	login as apiLogin,
	logout as apiLogout,
	refresh as apiRefresh,
	isAuthenticated as checkAuth
} from '$lib/api/auth';

/**
 * Auth state interface
 */
export interface AuthState {
	isAuthenticated: boolean;
	isLoading: boolean;
	error: string | null;
	username: string | null;
}

/**
 * Initial auth state
 */
const initialState: AuthState = {
	isAuthenticated: false,
	isLoading: false,
	error: null,
	username: null
};

/**
 * Create the auth store
 */
function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(initialState);

	// Token refresh interval (refresh 1 minute before expiry, assuming 15 min tokens)
	let refreshInterval: ReturnType<typeof setInterval> | null = null;
	const REFRESH_INTERVAL_MS = 14 * 60 * 1000; // 14 minutes

	/**
	 * Initialize auth state from stored tokens
	 */
	function initialize(): void {
		const isAuth = checkAuth();
		update((state) => ({
			...state,
			isAuthenticated: isAuth
		}));

		if (isAuth) {
			startTokenRefresh();
		}
	}

	/**
	 * Start automatic token refresh
	 */
	function startTokenRefresh(): void {
		stopTokenRefresh();
		refreshInterval = setInterval(async () => {
			try {
				await apiRefresh();
			} catch {
				// Refresh failed - user will be logged out on next API call
				stopTokenRefresh();
				set({
					...initialState,
					error: 'Session expired. Please log in again.'
				});
			}
		}, REFRESH_INTERVAL_MS);
	}

	/**
	 * Stop automatic token refresh
	 */
	function stopTokenRefresh(): void {
		if (refreshInterval) {
			clearInterval(refreshInterval);
			refreshInterval = null;
		}
	}

	/**
	 * Login with username and password
	 */
	async function login(username: string, password: string): Promise<boolean> {
		update((state) => ({
			...state,
			isLoading: true,
			error: null
		}));

		try {
			await apiLogin(username, password);
			update((state) => ({
				...state,
				isAuthenticated: true,
				isLoading: false,
				username
			}));
			startTokenRefresh();
			return true;
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Login failed';
			update((state) => ({
				...state,
				isLoading: false,
				error: message
			}));
			return false;
		}
	}

	/**
	 * Logout and clear tokens
	 */
	async function logout(): Promise<void> {
		stopTokenRefresh();
		try {
			await apiLogout();
		} catch {
			// Ignore logout errors
		}
		set(initialState);
	}

	/**
	 * Manually refresh the access token
	 */
	async function refresh(): Promise<boolean> {
		try {
			await apiRefresh();
			return true;
		} catch {
			stopTokenRefresh();
			set({
				...initialState,
				error: 'Session expired. Please log in again.'
			});
			return false;
		}
	}

	/**
	 * Clear any error message
	 */
	function clearError(): void {
		update((state) => ({
			...state,
			error: null
		}));
	}

	/**
	 * Check if currently authenticated
	 */
	function isAuthenticated(): boolean {
		return get({ subscribe }).isAuthenticated;
	}

	return {
		subscribe,
		initialize,
		login,
		logout,
		refresh,
		clearError,
		isAuthenticated
	};
}

/**
 * Auth store singleton
 */
export const authStore = createAuthStore();

/**
 * Derived store for just the authentication status
 */
export const isAuthenticated = derived(authStore, ($auth) => $auth.isAuthenticated);

/**
 * Derived store for loading state
 */
export const isAuthLoading = derived(authStore, ($auth) => $auth.isLoading);

/**
 * Derived store for auth error
 */
export const authError = derived(authStore, ($auth) => $auth.error);
