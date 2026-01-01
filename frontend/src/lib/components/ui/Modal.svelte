<script lang="ts">
	import type { Snippet } from 'svelte';
	import { X } from 'lucide-svelte';

	interface Props {
		open?: boolean;
		title?: string;
		children: Snippet;
		footer?: Snippet;
		onclose?: () => void;
	}

	let { open = false, title, children, footer, onclose }: Props = $props();

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			onclose?.();
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			onclose?.();
		}
	}
</script>

{#if open}
	<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		role="dialog"
		aria-modal="true"
		tabindex="-1"
		onclick={handleBackdropClick}
		onkeydown={handleKeydown}
	>
		<div class="bg-surface-primary border border-border-primary rounded-lg shadow-xl max-w-md w-full mx-4 max-h-[90vh] overflow-hidden flex flex-col">
			{#if title}
				<div class="flex items-center justify-between px-4 py-3 border-b border-border-secondary">
					<h2 class="text-lg font-medium text-text-primary">{title}</h2>
					<button
						type="button"
						class="p-1 text-text-secondary hover:text-text-primary rounded transition-colors"
						onclick={onclose}
						aria-label="Close"
					>
						<X size={18} />
					</button>
				</div>
			{/if}
			<div class="p-4 overflow-y-auto">
				{@render children()}
			</div>
			{#if footer}
				<div class="px-4 py-3 border-t border-border-secondary flex justify-end gap-2">
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
{/if}
