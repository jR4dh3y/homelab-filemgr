<script lang="ts">
	/**
	 * SearchBar component with search input and loading state
	 * Requirements: 9.1, 9.4
	 */

	interface Props {
		/** Current search query */
		value?: string;
		/** Placeholder text */
		placeholder?: string;
		/** Whether search is in progress */
		isLoading?: boolean;
		/** Callback when search is submitted */
		onSearch?: (query: string) => void;
		/** Callback when input changes */
		onInput?: (query: string) => void;
		/** Callback when search is cleared */
		onClear?: () => void;
	}

	let {
		value = '',
		placeholder = 'Search files...',
		isLoading = false,
		onSearch,
		onInput,
		onClear
	}: Props = $props();

	let inputValue = $state('');

	// Sync external value changes (including initial value)
	$effect(() => {
		inputValue = value;
	});

	/**
	 * Handle form submission
	 */
	function handleSubmit(event: Event): void {
		event.preventDefault();
		const trimmed = inputValue.trim();
		if (trimmed && onSearch) {
			onSearch(trimmed);
		}
	}

	/**
	 * Handle input change
	 */
	function handleInput(event: Event): void {
		const target = event.target as HTMLInputElement;
		inputValue = target.value;
		if (onInput) {
			onInput(inputValue);
		}
	}

	/**
	 * Handle clear button click
	 */
	function handleClear(): void {
		inputValue = '';
		if (onClear) {
			onClear();
		}
		if (onInput) {
			onInput('');
		}
	}

	/**
	 * Handle keyboard events
	 */
	function handleKeydown(event: KeyboardEvent): void {
		if (event.key === 'Escape') {
			handleClear();
		}
	}
</script>

<form onsubmit={handleSubmit} class="relative w-full max-w-md">
	<div class="relative">
		<!-- Search icon -->
		<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
			{#if isLoading}
				<!-- Loading spinner -->
				<svg
					class="h-5 w-5 animate-spin text-gray-400"
					fill="none"
					viewBox="0 0 24 24"
					aria-hidden="true"
				>
					<circle
						class="opacity-25"
						cx="12"
						cy="12"
						r="10"
						stroke="currentColor"
						stroke-width="4"
					></circle>
					<path
						class="opacity-75"
						fill="currentColor"
						d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
					></path>
				</svg>
			{:else}
				<!-- Search icon -->
				<svg
					class="h-5 w-5 text-gray-400"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
					aria-hidden="true"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
					/>
				</svg>
			{/if}
		</div>

		<!-- Input field -->
		<input
			type="search"
			{placeholder}
			value={inputValue}
			oninput={handleInput}
			onkeydown={handleKeydown}
			disabled={isLoading}
			class="block w-full rounded-lg border border-gray-300 bg-white py-2 pl-10 pr-10 text-sm text-gray-900 placeholder-gray-500 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100 dark:placeholder-gray-400 dark:focus:border-blue-400 dark:focus:ring-blue-400"
			aria-label="Search"
		/>

		<!-- Clear button -->
		{#if inputValue}
			<button
				type="button"
				onclick={handleClear}
				class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
				aria-label="Clear search"
			>
				<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M6 18L18 6M6 6l12 12"
					/>
				</svg>
			</button>
		{/if}
	</div>

	<!-- Loading indicator text -->
	{#if isLoading}
		<p class="mt-1 text-xs text-gray-500 dark:text-gray-400">Searching...</p>
	{/if}
</form>
