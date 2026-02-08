<script lang="ts">
	/**
	 * InlineRename component - reusable inline text editing with save/cancel
	 * Extracts common renaming logic used across DriveCard, SystemDriveCard, Sidebar
	 */
	import { Check, X } from 'lucide-svelte';

	interface Props {
		value: string;
		onSave: (newValue: string) => void;
		onCancel: () => void;
		placeholder?: string;
		class?: string;
	}

	let { value, onSave, onCancel, placeholder = '', class: className = '' }: Props = $props();

	let inputValue = $state(value);

	function handleSave(): void {
		onSave(inputValue.trim());
	}

	function handleKeydown(e: KeyboardEvent): void {
		if (e.key === 'Enter') handleSave();
		if (e.key === 'Escape') onCancel();
	}

	function handleFocus(e: FocusEvent): void {
		(e.target as HTMLInputElement).select();
	}

	function handleButtonClick(e: MouseEvent, action: () => void): void {
		e.stopPropagation();
		action();
	}
</script>

<div class="flex items-center gap-1 h-5 {className}">
	<input
		type="text"
		bind:value={inputValue}
		onkeydown={handleKeydown}
		onfocus={handleFocus}
		{placeholder}
		class="flex-1 min-w-0 h-5 box-border bg-surface-primary border border-border-focus rounded px-2 text-xs text-text-primary outline-none"
	/>
	<button
		type="button"
		onclick={(e: MouseEvent) => handleButtonClick(e, handleSave)}
		class="shrink-0 w-5 h-5 flex items-center justify-center text-success hover:text-green-400 hover:bg-success/20 rounded transition-colors"
		title="Save"
	>
		<Check size={12} />
	</button>
	<button
		type="button"
		onclick={(e: MouseEvent) => handleButtonClick(e, onCancel)}
		class="shrink-0 w-5 h-5 flex items-center justify-center text-text-muted hover:text-danger hover:bg-danger/20 rounded transition-colors"
		title="Cancel"
	>
		<X size={12} />
	</button>
</div>
