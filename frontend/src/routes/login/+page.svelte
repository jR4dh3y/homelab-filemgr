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
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	}

	.login-card {
		width: 100%;
		max-width: 400px;
		background: white;
		border-radius: 1rem;
		box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
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
		font-weight: 700;
		color: #111827;
		margin: 0 0 0.5rem;
	}

	.login-subtitle {
		font-size: 0.875rem;
		color: #6b7280;
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
		background: #fef2f2;
		border: 1px solid #fecaca;
		border-radius: 0.5rem;
		color: #dc2626;
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
		color: #dc2626;
		cursor: pointer;
		border-radius: 0.25rem;
		transition: background-color 0.15s;
	}

	.error-dismiss:hover {
		background: #fee2e2;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #374151;
	}

	.form-input {
		padding: 0.75rem 1rem;
		font-size: 1rem;
		border: 1px solid #d1d5db;
		border-radius: 0.5rem;
		background: white;
		color: #111827;
		transition:
			border-color 0.15s,
			box-shadow 0.15s;
	}

	.form-input:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.form-input:disabled {
		background: #f9fafb;
		cursor: not-allowed;
	}

	.form-input::placeholder {
		color: #9ca3af;
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
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		border: none;
		border-radius: 0.5rem;
		cursor: pointer;
		transition:
			opacity 0.15s,
			transform 0.15s;
	}

	.submit-btn:hover:not(:disabled) {
		opacity: 0.9;
		transform: translateY(-1px);
	}

	.submit-btn:active:not(:disabled) {
		transform: translateY(0);
	}

	.submit-btn:disabled {
		opacity: 0.6;
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

	/* Dark mode */
	@media (prefers-color-scheme: dark) {
		.login-container {
			background: linear-gradient(135deg, #1e3a5f 0%, #312e81 100%);
		}

		.login-card {
			background: #1f2937;
		}

		.login-title {
			color: #f9fafb;
		}

		.login-subtitle {
			color: #9ca3af;
		}

		.error-message {
			background: #450a0a;
			border-color: #7f1d1d;
			color: #fca5a5;
		}

		.error-dismiss {
			color: #fca5a5;
		}

		.error-dismiss:hover {
			background: #7f1d1d;
		}

		.form-label {
			color: #d1d5db;
		}

		.form-input {
			background: #374151;
			border-color: #4b5563;
			color: #f9fafb;
		}

		.form-input:focus {
			border-color: #60a5fa;
			box-shadow: 0 0 0 3px rgba(96, 165, 250, 0.1);
		}

		.form-input:disabled {
			background: #1f2937;
		}

		.form-input::placeholder {
			color: #6b7280;
		}
	}
</style>
