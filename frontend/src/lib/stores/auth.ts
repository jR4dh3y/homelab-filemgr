/**
 * Auth store for managing authentication state
 */

import { writable } from 'svelte/store';

export interface AuthState {
	isAuthenticated: boolean;
	isLoading: boolean;
	error: string | null;
}

const initialState: AuthState = {
	isAuthenticated: false,
	isLoading: false,
	error: null
};

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(initialState);

	async function login(username: string, password: string): Promise<boolean> {
		update((state) => ({ ...state, isLoading: true, error: null }));

		try {
			const response = await fetch('/api/v1/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username, password })
			});

			if (!response.ok) {
				throw new Error('Login failed');
			}

			const data = await response.json();
			localStorage.setItem('accessToken', data.accessToken);
			
			update((state) => ({ ...state, isAuthenticated: true, isLoading: false }));
			return true;
		} catch (err) {
			update((state) => ({
				...state,
				isLoading: false,
				error: err instanceof Error ? err.message : 'Login failed'
			}));
			return false;
		}
	}

	function logout(): void {
		localStorage.removeItem('accessToken');
		set(initialState);
	}

	return { subscribe, login, logout };
}

export const authStore = createAuthStore();