<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { authStore, isAuthenticated } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
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

	onMount(() => {
		authStore.initialize();
		initialized = true;
	});

	// Reactive navigation based on auth state
	$effect(() => {
		if (!initialized) return;

		const currentPath = $page.url.pathname;
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
	{:else}
		<div class="app-container">
			{#if $isAuthenticated && !$page.url.pathname.startsWith('/login')}
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
				class:with-header={$isAuthenticated && !$page.url.pathname.startsWith('/login')}
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
		background: #f9fafb;
	}

	.loading-spinner {
		width: 40px;
		height: 40px;
		border: 3px solid #e5e7eb;
		border-top-color: #3b82f6;
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
		background: #f3f4f6;
	}

	.app-header {
		background: white;
		border-bottom: 1px solid #e5e7eb;
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
		color: #111827;
		font-weight: 600;
		font-size: 1.125rem;
	}

	.app-logo:hover {
		color: #3b82f6;
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
		color: #6b7280;
		background: transparent;
		border: 1px solid #e5e7eb;
		border-radius: 0.375rem;
		cursor: pointer;
		transition: all 0.15s;
	}

	.logout-btn:hover {
		color: #111827;
		border-color: #d1d5db;
		background: #f9fafb;
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

	/* Dark mode */
	@media (prefers-color-scheme: dark) {
		.loading-screen {
			background: #111827;
		}

		.loading-spinner {
			border-color: #374151;
			border-top-color: #60a5fa;
		}

		.app-container {
			background: #0f172a;
		}

		.app-header {
			background: #1f2937;
			border-bottom-color: #374151;
		}

		.app-logo {
			color: #f9fafb;
		}

		.app-logo:hover {
			color: #60a5fa;
		}

		.logout-btn {
			color: #9ca3af;
			border-color: #374151;
		}

		.logout-btn:hover {
			color: #f9fafb;
			border-color: #4b5563;
			background: #374151;
		}
	}
</style>
