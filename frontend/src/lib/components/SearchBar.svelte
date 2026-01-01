<script lang="ts">
	/**
	 * SearchBar component with search input and loading state
	 */
	import { Search, X } from 'lucide-svelte';
	import { Spinner } from '$lib/components/ui';

	interface Props {
		value?: string;
		placeholder?: string;
		isLoading?: boolean;
		onSearch?: (query: string) => void;
		onInput?: (query: string) => void;
		onClear?: () => void;
	}

	let {
		value = '',
		placeholder = 'Search files...',
		isLoading = false,
		onSearch,
		onInput,
		onClear,
	}: Props = $props();

	let inputValue = $state(value);

	$effect(() => {
		if (value !== inputValue) {
			inputValue = value;
		}
	});

	function handleSubmit(event: Event): void {
		event.preventDefault();
		const trimmed = inputValue.trim();
		if (trimmed && onSearch) {
			onSearch(trimmed);
		}
	}

	function handleInput(event: Event): void {
		const target = event.target as HTMLInputElement;
		inputValue = target.value;
		onInput?.(inputValue);
	}

	function handleClear(): void {
		inputValue = '';
		onClear?.();
		onInput?.('');
	}

	function handleKeydown(event: KeyboardEvent): void {
		if (event.key === 'Escape') {
			handleClear();
		}
	}
</script>

<form onsubmit={handleSubmit} class="relative w-full max-w-md">
	<div class="relative">
		<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
			{#if isLoading}
				<Spinner size="sm" />
			{:else}
				<Search size={18} class="text-text-muted" />
			{/if}
		</div>

		<input
			type="search"
			{placeholder}
			value={inputValue}
			oninput={handleInput}
			onkeydown={handleKeydown}
			disabled={isLoading}
			class="block w-full rounded-lg border border-border-primary bg-surface-secondary py-2 pr-10 pl-10 text-sm text-text-primary placeholder-text-muted focus:border-border-focus focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
			aria-label="Search"
		/>

		{#if inputValue}
			<button
				type="button"
				onclick={handleClear}
				class="absolute inset-y-0 right-0 flex items-center pr-3 text-text-muted hover:text-text-secondary"
				aria-label="Clear search"
			>
				<X size={18} />
			</button>
		{/if}
	</div>

	{#if isLoading}
		<p class="mt-1 text-xs text-text-muted">Searching...</p>
	{/if}
</form>
