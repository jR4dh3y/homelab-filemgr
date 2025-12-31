<script lang="ts">
	/**
	 * Login page component
	 * Requirements: 7.1
	 */
	import { authStore, authError, isAuthLoading } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

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

<div class="login-container">
	<div class="login-card">
		<div class="login-header">
			<span class="login-icon">📁</span>
			<h1 class="login-title">File Manager</h1>
			<p class="login-subtitle">Sign in to access your files</p>
		</div>

		<form class="login-form" onsubmit={handleSubmit}>
			{#if $authError}
				<div class="error-message" role="alert">
					<span class="error-icon">⚠️</span>
					<span>{$authError}</span>
					<button
						type="button"
						class="error-dismiss"
						onclick={clearError}
						aria-label="Dismiss error"
					>
						×
					</button>
				</div>
			{/if}

			<div class="form-group">
				<label for="username" class="form-label">Username</label>
				<input
					type="text"
					id="username"
					bind:value={username}
					class="form-input"
					placeholder="Enter your username"
					autocomplete="username"
					required
					disabled={$isAuthLoading}
				/>
			</div>

			<div class="form-group">
				<label for="password" class="form-label">Password</label>
				<input
					type="password"
					id="password"
					bind:value={password}
					class="form-input"
					placeholder="Enter your password"
					autocomplete="current-password"
					required
					disabled={$isAuthLoading}
				/>
			</div>

			<button
				type="submit"
				class="submit-btn"
				disabled={$isAuthLoading || !username.trim() || !password}
			>
				{#if $isAuthLoading}
					<span class="btn-spinner"></span>
					<span>Signing in...</span>
				{:else}
					<span>Sign In</span>
				{/if}
			</button>
		</form>
	</div>
</div>

<style>
	.login-container {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
		background: #141414;
	}

	.login-card {
		width: 100%;
		max-width: 400px;
		background: #1e1e1e;
		border: 1px solid #2a2a2a;
		border-radius: 8px;
		padding: 2rem;
	}

	.login-header {
		text-align: center;
		margin-bottom: 2rem;
	}

	.login-icon {
		font-size: 3rem;
		display: block;
		margin-bottom: 0.5rem;
	}

	.login-title {
		font-size: 1.5rem;
		font-weight: 600;
		color: #e0e0e0;
		margin: 0 0 0.5rem;
	}

	.login-subtitle {
		font-size: 0.875rem;
		color: #888;
		margin: 0;
	}

	.login-form {
		display: flex;
		flex-direction: column;
		gap: 1.25rem;
	}

	.error-message {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		background: #3d1f1f;
		border: 1px solid #5c2a2a;
		border-radius: 4px;
		color: #f87171;
		font-size: 0.875rem;
	}

	.error-icon {
		flex-shrink: 0;
	}

	.error-dismiss {
		margin-left: auto;
		padding: 0;
		width: 1.5rem;
		height: 1.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: none;
		font-size: 1.25rem;
		color: #f87171;
		cursor: pointer;
		border-radius: 4px;
		transition: background-color 0.15s;
	}

	.error-dismiss:hover {
		background: #5c2a2a;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #aaa;
	}

	.form-input {
		padding: 0.75rem 1rem;
		font-size: 1rem;
		border: 1px solid #333;
		border-radius: 4px;
		background: #252525;
		color: #e0e0e0;
		transition: border-color 0.15s;
	}

	.form-input:focus {
		outline: none;
		border-color: #4a9eff;
	}

	.form-input:disabled {
		background: #1a1a1a;
		cursor: not-allowed;
		color: #666;
	}

	.form-input::placeholder {
		color: #555;
	}

	.submit-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.875rem 1.5rem;
		font-size: 1rem;
		font-weight: 600;
		color: white;
		background: #2d4a6f;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		transition: background-color 0.15s;
	}

	.submit-btn:hover:not(:disabled) {
		background: #345580;
	}

	.submit-btn:active:not(:disabled) {
		background: #2a4060;
	}

	.submit-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-spinner {
		width: 1rem;
		height: 1rem;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-top-color: white;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
