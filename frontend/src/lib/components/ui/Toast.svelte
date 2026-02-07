<script lang="ts">
	/**
	 * Toast notification container component
	 * Displays toast notifications from the toastStore
	 */
	import { toastStore, type Toast } from '$lib/stores/toast.svelte';
	import { CheckCircle, XCircle, Info, AlertTriangle, X } from 'lucide-svelte';
	import { fly, fade } from 'svelte/transition';

	const iconMap = {
		success: CheckCircle,
		error: XCircle,
		info: Info,
		warning: AlertTriangle
	};

	const colorMap = {
		success: 'bg-success/10 border-success/30 text-success',
		error: 'bg-danger/10 border-danger/30 text-danger',
		info: 'bg-accent/10 border-accent/30 text-accent',
		warning: 'bg-warning/10 border-warning/30 text-warning'
	};

	const iconColorMap = {
		success: 'text-success',
		error: 'text-danger',
		info: 'text-accent',
		warning: 'text-warning'
	};

	function handleDismiss(id: string) {
		toastStore.remove(id);
	}
</script>

{#if toastStore.toasts.length > 0}
	<div class="fixed bottom-4 right-4 z-[100] flex flex-col gap-2 max-w-sm pointer-events-none">
		{#each toastStore.toasts as toast (toast.id)}
			<div
				class="pointer-events-auto flex items-start gap-3 px-4 py-3 rounded-lg border shadow-lg backdrop-blur-sm {colorMap[toast.type]}"
				in:fly={{ x: 100, duration: 200 }}
				out:fade={{ duration: 150 }}
				role="alert"
			>
				<!-- Icon -->
				<div class="shrink-0 mt-0.5 {iconColorMap[toast.type]}">
					<svelte:component this={iconMap[toast.type]} size={18} />
				</div>

				<!-- Message -->
				<p class="flex-1 text-sm text-text-primary m-0 pr-2">
					{toast.message}
				</p>

				<!-- Dismiss button -->
				<button
					type="button"
					class="shrink-0 p-0.5 rounded text-text-muted hover:text-text-primary hover:bg-surface-elevated transition-colors"
					onclick={() => handleDismiss(toast.id)}
					aria-label="Dismiss notification"
				>
					<X size={14} />
				</button>
			</div>
		{/each}
	</div>
{/if}
