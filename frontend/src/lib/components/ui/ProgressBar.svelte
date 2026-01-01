<script lang="ts">
	interface Props {
		value: number;
		max?: number;
		variant?: 'default' | 'success' | 'warning' | 'danger';
		size?: 'sm' | 'md';
		showLabel?: boolean;
	}

	let { value, max = 100, variant = 'default', size = 'md', showLabel = false }: Props = $props();

	const percentage = $derived(Math.min(100, Math.max(0, (value / max) * 100)));

	const variantClasses: Record<string, string> = {
		default: 'bg-accent',
		success: 'bg-success',
		warning: 'bg-warning',
		danger: 'bg-danger',
	};

	const sizeClasses: Record<string, string> = {
		sm: 'h-1',
		md: 'h-2',
	};
</script>

<div class="w-full">
	<div class="w-full bg-surface-elevated rounded-full overflow-hidden {sizeClasses[size]}">
		<div
			class="{variantClasses[variant]} h-full rounded-full transition-all duration-300"
			style="width: {percentage}%"
			role="progressbar"
			aria-valuenow={value}
			aria-valuemin={0}
			aria-valuemax={max}
		></div>
	</div>
	{#if showLabel}
		<span class="text-xs text-text-secondary mt-1">{percentage.toFixed(0)}%</span>
	{/if}
</div>
