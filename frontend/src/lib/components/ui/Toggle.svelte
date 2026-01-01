<script lang="ts">
	interface Props {
		checked?: boolean;
		disabled?: boolean;
		id?: string;
		label?: string;
		onchange?: (checked: boolean) => void;
	}

	let { checked = $bindable(false), disabled = false, id, label, onchange }: Props = $props();

	function handleChange() {
		checked = !checked;
		onchange?.(checked);
	}
</script>

<label class="inline-flex items-center gap-2 cursor-pointer" class:opacity-50={disabled}>
	<button
		type="button"
		role="switch"
		aria-checked={checked}
		aria-label={label || 'Toggle'}
		{disabled}
		{id}
		onclick={handleChange}
		class="relative w-10 h-5 rounded-full transition-colors duration-200 {checked
			? 'bg-accent'
			: 'bg-surface-elevated'}"
	>
		<span
			class="absolute top-0.5 left-0.5 w-4 h-4 bg-white rounded-full transition-transform duration-200 {checked
				? 'translate-x-5'
				: 'translate-x-0'}"
		></span>
	</button>
	{#if label}
		<span class="text-sm text-text-primary">{label}</span>
	{/if}
</label>
