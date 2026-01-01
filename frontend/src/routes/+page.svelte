<script lang="ts">
	/**
	 * Main page - redirects to browse
	 */
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Spinner } from '$lib/components/ui';

	onMount(() => {
		goto('/browse', { replaceState: true });
	});
</script>

<div class="flex items-center justify-center min-h-screen bg-surface-primary">
	<div class="flex flex-col items-center gap-4 text-text-secondary">
		<Spinner />
		<p>Redirecting to file browser...</p>
	</div>
</div>
