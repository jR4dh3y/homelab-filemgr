<script lang="ts">
	/**
	 * Main page - redirects to browse
	 * Requirements: 1.2
	 */
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		goto('/browse', { replaceState: true });
	});
</script>

<div class="redirect-container">
	<div class="redirect-content">
		<div class="spinner"></div>
		<p>Redirecting to file browser...</p>
	</div>
</div>

<style>
	.redirect-container {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		background: #f3f4f6;
	}

	.redirect-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		color: #6b7280;
	}

	.spinner {
		width: 32px;
		height: 32px;
		border: 3px solid #e5e7eb;
		border-top-color: #3b82f6;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	/* Dark mode */
	@media (prefers-color-scheme: dark) {
		.redirect-container {
			background: #0f172a;
		}

		.redirect-content {
			color: #9ca3af;
		}

		.spinner {
			border-color: #374151;
			border-top-color: #60a5fa;
		}
	}
</style>
