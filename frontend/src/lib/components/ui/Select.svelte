<script lang="ts">
	interface Option {
		value: string;
		label: string;
	}

	interface Props {
		value?: string;
		options: Option[];
		disabled?: boolean;
		id?: string;
		name?: string;
		onchange?: (value: string) => void;
	}

	let { value = $bindable(''), options, disabled = false, id, name, onchange }: Props = $props();

	function handleChange(e: Event) {
		const target = e.target as HTMLSelectElement;
		value = target.value;
		onchange?.(value);
	}
</script>

<select
	{value}
	{disabled}
	{id}
	{name}
	onchange={handleChange}
	class="w-full px-3 py-2 bg-surface-secondary border border-border-primary rounded text-text-primary text-sm transition-colors duration-150 focus:outline-none focus:border-border-focus disabled:opacity-50 disabled:cursor-not-allowed appearance-none cursor-pointer"
>
	{#each options as option (option.value)}
		<option value={option.value}>{option.label}</option>
	{/each}
</select>
