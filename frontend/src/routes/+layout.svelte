<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { authStore, isAuthenticated } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { CONFIG } from '$lib/config';
	import { Spinner, Button } from '$lib/components/ui';
	import { FolderOpen } from 'lucide-svelte';

	let { children } = $props();
	let initialized = $state(false);

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: CONFIG.query.staleTimeMs,
				retry: 1,
			},
		},
	});

	// Public routes that don't require authentication
	const publicRoutes = ['/login', '/test'];
	const isBrowsePage = $derived(page.url.pathname.startsWith('/browse'));
	const isLoginPage = $derived(page.url.pathname.startsWith('/login'));

	onMount(() => {
		authStore.initialize();
		initialized = true;
	});

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
		<div class="flex items-center justify-center min-h-screen bg-surface-primary">
			<Spinner size="lg" />
		</div>
	{:else if isBrowsePage}
		{@render children()}
	{:else}
		<div class="min-h-screen flex flex-col bg-surface-primary">
			{#if $isAuthenticated && !isLoginPage}
				<header class="bg-surface-primary border-b border-border-secondary px-4 sticky top-0 z-50">
					<div class="max-w-[1400px] mx-auto flex items-center justify-between h-14">
						<a href="/browse" class="flex items-center gap-2 no-underline text-text-primary font-semibold text-lg hover:text-accent">
							<FolderOpen size={24} class="text-accent" />
							<span>File Manager</span>
						</a>
						<nav class="flex items-center gap-4">
							<Button variant="secondary" size="sm" onclick={handleLogout}>Logout</Button>
						</nav>
					</div>
				</header>
			{/if}
			<main class="flex-1 flex flex-col {$isAuthenticated && !isLoginPage ? 'p-6 max-w-[1400px] w-full mx-auto' : ''}">
				{@render children()}
			</main>
		</div>
	{/if}
</QueryClientProvider>
