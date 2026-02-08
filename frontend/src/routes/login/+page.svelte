<script lang="ts">
	/**
	 * Login page component
	 */
	import { authStore, authError, isAuthLoading } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { Button, Input, Spinner } from '$lib/components/ui';
	import { X, FolderOpen, AlertTriangle } from 'lucide-svelte';

	let username = $state('');
	let password = $state('');

	async function handleSubmit(event: Event) {
		event.preventDefault();

		if (!username.trim() || !password) {
			return;
		}

		const success = await authStore.login(username.trim(), password);
		if (success) {
			goto('/browse');
		}
	}

	function clearError() {
		authStore.clearError();
	}
</script>

<svelte:head>
	<title>Login - File Manager</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center p-4 bg-surface-primary">
	<div class="w-full max-w-[400px] bg-surface-secondary border border-border-primary rounded-lg p-8">
		<div class="flex flex-col items-center mb-8">
			<div class="flex items-center gap-3 mb-2">
				<span class="text-accent"><FolderOpen size={32} /></span>
				<h1 class="text-2xl font-semibold text-text-primary m-0">File Manager</h1>
			</div>
			<p class="text-sm text-text-secondary m-0">Sign in to access your files</p>
		</div>

		<form class="flex flex-col gap-5" onsubmit={handleSubmit}>
			{#if $authError}
				<div class="flex items-center gap-2 px-4 py-3 bg-danger/20 border border-danger/30 rounded text-danger text-sm" role="alert">
					<span class="shrink-0"><AlertTriangle size={16} /></span>
					<span class="flex-1">{$authError}</span>
					<button
						type="button"
						class="ml-auto p-0 w-6 h-6 flex items-center justify-center bg-transparent border-none text-xl text-danger cursor-pointer rounded transition-colors hover:bg-danger/30"
						onclick={clearError}
						aria-label="Dismiss error"
					>
						<X size={16} />
					</button>
				</div>
			{/if}

			<div class="flex flex-col gap-2">
				<label for="username" class="text-sm font-medium text-text-secondary">Username</label>
				<Input
					type="text"
					id="username"
					bind:value={username}
					placeholder="Enter your username"
					autocomplete="username"
					required
					disabled={$isAuthLoading}
				/>
			</div>

			<div class="flex flex-col gap-2">
				<label for="password" class="text-sm font-medium text-text-secondary">Password</label>
				<Input
					type="password"
					id="password"
					bind:value={password}
					placeholder="Enter your password"
					autocomplete="current-password"
					required
					disabled={$isAuthLoading}
				/>
			</div>

			<Button
				type="submit"
				variant="primary"
				disabled={$isAuthLoading || !username.trim() || !password}
			>
				{#if $isAuthLoading}
					<Spinner size="sm" />
					<span>Signing in...</span>
				{:else}
					<span>Sign In</span>
				{/if}
			</Button>
		</form>
	</div>
</div>
