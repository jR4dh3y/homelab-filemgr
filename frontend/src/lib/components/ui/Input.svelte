<script lang="ts">
	interface Props {
		type?: 'text' | 'password' | 'email' | 'search' | 'number';
		value?: string;
		placeholder?: string;
		disabled?: boolean;
		required?: boolean;
		id?: string;
		name?: string;
		autocomplete?: AutoFill;
		oninput?: (e: Event) => void;
		onkeydown?: (e: KeyboardEvent) => void;
	}

	let {
		type = 'text',
		value = $bindable(''),
		placeholder,
		disabled = false,
		required = false,
		id,
		name,
		autocomplete,
		oninput,
		onkeydown,
	}: Props = $props();
</script>

<input
	{type}
	bind:value
	{placeholder}
	{disabled}
	{required}
	{id}
	{name}
	autocomplete={autocomplete}
	{oninput}
	{onkeydown}
	class="w-full px-3 py-2 bg-surface-secondary border border-border-primary rounded text-text-primary placeholder:text-text-muted text-sm transition-colors duration-150 focus:outline-none focus:border-border-focus disabled:opacity-50 disabled:cursor-not-allowed"
/>
