<script lang="ts">
	/**
	 * Breadcrumb component for path navigation
	 * Requirements: 1.4
	 */

	interface Props {
		/** Current path segments */
		segments: string[];
		/** Callback when a segment is clicked */
		onNavigate?: (path: string) => void;
	}

	let { segments = [], onNavigate }: Props = $props();

	/**
	 * Build the full path up to a given segment index
	 */
	function buildPath(index: number): string {
		return segments.slice(0, index + 1).join('/');
	}

	/**
	 * Handle segment click
	 */
	function handleClick(index: number): void {
		if (onNavigate) {
			onNavigate(buildPath(index));
		}
	}

	/**
	 * Handle root click
	 */
	function handleRootClick(): void {
		if (onNavigate) {
			onNavigate('');
		}
	}
</script>

<nav aria-label="Breadcrumb" class="flex items-center space-x-1 text-sm">
	<!-- Home/Root -->
	<button
		type="button"
		onclick={handleRootClick}
		class="flex items-center text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
		aria-label="Go to root"
	>
		<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
			/>
		</svg>
	</button>

	{#if segments.length > 0}
		{#each segments as segment, index (index)}
			<!-- Separator -->
			<svg
				class="h-4 w-4 flex-shrink-0 text-gray-400 dark:text-gray-500"
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
				aria-hidden="true"
			>
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
			</svg>

			{#if index === segments.length - 1}
				<!-- Current segment (not clickable) -->
				<span
					class="max-w-[200px] truncate font-medium text-gray-900 dark:text-gray-100"
					title={segment}
					aria-current="page"
				>
					{segment}
				</span>
			{:else}
				<!-- Clickable segment -->
				<button
					type="button"
					onclick={() => handleClick(index)}
					class="max-w-[150px] truncate text-gray-500 transition-colors hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
					title={segment}
				>
					{segment}
				</button>
			{/if}
		{/each}
	{/if}
</nav>
