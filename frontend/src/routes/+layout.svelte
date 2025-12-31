<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { authStore, isAuthenticated } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';

	let { children } = $props();
	let initialized = $state(false);

	// Create QueryClient for TanStack Query
	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 1000 * 60, // 1 minute
				retry: 1
			}
		}
	});

	// Public routes that don't require authentication
	const publicRoutes = ['/login'];

	// Check if we're on the browse page (full-screen file manager)
	const isBrowsePage = $derived(page.url.pathname.startsWith('/browse'));
	const isLoginPage = $derived(page.url.pathname.startsWith('/login'));

	onMount(() => {
		authStore.initialize();
		initialized = true;
	});

	// Reactive navigation based on auth state
	$effect(() => {
		if (!initialized) return;

		const currentPath = page.url.pathname;
		const isPublicRoute = publicRoutes.some((route) => currentPath.startsWith(route));

		if (!$isAuthenticated && !isPublicRoute) {
			goto('/login');
		} else if ($isAuthenticated && currentPath.startsWith('/login')) {
			goto('/browse');
		}
	});

	async function handleLogout() {
		await authStore.logout();
		goto('/login');
	}
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<QueryClientProvider client={queryClient}>
	{#if !initialized}
		<div class="loading-screen">
			<div class="loading-spinner"></div>
		</div>
	{:else if isBrowsePage}
		<!-- Full-screen file manager mode -->
		{@render children()}
	{:else}
		<div class="app-container">
			{#if $isAuthenticated && !isLoginPage}
				<header class="app-header">
					<div class="header-content">
						<a href="/browse" class="app-logo">
							<span class="logo-icon">📁</span>
							<span class="logo-text">File Manager</span>
						</a>
						<nav class="header-nav">
							<button type="button" class="logout-btn" onclick={handleLogout}> Logout </button>
						</nav>
					</div>
				</header>
			{/if}
			<main
				class="app-main"
				class:with-header={$isAuthenticated && !isLoginPage}
			>
				{@render children()}
			</main>
		</div>
	{/if}
</QueryClientProvider>

<style>
	.loading-screen {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		background: #1e1e1e;
	}

	.loading-spinner {
		width: 40px;
		height: 40px;
		border: 3px solid #333;
		border-top-color: #4a9eff;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.app-container {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		background: #1e1e1e;
	}

	.app-header {
		background: #1a1a1a;
		border-bottom: 1px solid #2a2a2a;
		padding: 0 1rem;
		position: sticky;
		top: 0;
		z-index: 50;
	}

	.header-content {
		max-width: 1400px;
		margin: 0 auto;
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: 56px;
	}

	.app-logo {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		text-decoration: none;
		color: #ccc;
		font-weight: 600;
		font-size: 1.125rem;
	}

	.app-logo:hover {
		color: #4a9eff;
	}

	.logo-icon {
		font-size: 1.5rem;
	}

	.header-nav {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.logout-btn {
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
		font-weight: 500;
		color: #888;
		background: transparent;
		border: 1px solid #333;
		border-radius: 0.375rem;
		cursor: pointer;
		transition: all 0.15s;
	}

	.logout-btn:hover {
		color: #ccc;
		border-color: #444;
		background: #252525;
	}

	.app-main {
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.app-main.with-header {
		padding: 1.5rem;
		max-width: 1400px;
		width: 100%;
		margin: 0 auto;
	}
</style>
