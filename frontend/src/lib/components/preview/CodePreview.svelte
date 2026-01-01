<script lang="ts">
	/**
	 * CodePreview - Code/text viewer with syntax highlighting using Monaco Editor
	 */
	import { onMount, onDestroy } from 'svelte';
	import { getMonacoLanguage } from '$lib/utils/fileTypes';
	import { getFileContent } from '$lib/api/files';
	import { Spinner } from '$lib/components/ui';

	interface Props {
		url: string;
		filename: string;
	}

	let { url, filename }: Props = $props();

	let containerElement: HTMLDivElement | null = $state(null);
	let editor: any = $state(null);
	let monaco: any = $state(null);
	let content = $state<string | null>(null);
	let error = $state<string | null>(null);
	let loading = $state(true);

	const language = $derived(getMonacoLanguage(filename));

	onMount(async () => {
		// Load file content
		try {
			content = await getFileContent(url);
		} catch (e) {
			error = 'Failed to load file content.';
			loading = false;
			return;
		}

		// Dynamically import Monaco Editor
		try {
			const monacoModule = await import('monaco-editor');
			monaco = monacoModule;

			// Configure Monaco environment for web workers
			self.MonacoEnvironment = {
				getWorker: function (_moduleId: string, label: string) {
					return new Worker(
						URL.createObjectURL(new Blob([`self.onmessage = function() {}`], { type: 'text/javascript' }))
					);
				},
			};

			if (containerElement && content !== null) {
				// Define dark theme
				monaco.editor.defineTheme('filemanager-dark', {
					base: 'vs-dark',
					inherit: true,
					rules: [],
					colors: {
						'editor.background': '#1e1e1e',
						'editor.foreground': '#d4d4d4',
						'editorLineNumber.foreground': '#5a5a5a',
						'editorLineNumber.activeForeground': '#c6c6c6',
						'editor.selectionBackground': '#264f78',
						'editor.lineHighlightBackground': '#2a2a2a',
					},
				});

				editor = monaco.editor.create(containerElement, {
					value: content,
					language: language,
					theme: 'filemanager-dark',
					readOnly: true,
					minimap: { enabled: true },
					scrollBeyondLastLine: false,
					fontSize: 13,
					fontFamily: "'Fira Code', 'Cascadia Code', 'JetBrains Mono', Consolas, monospace",
					lineNumbers: 'on',
					renderLineHighlight: 'line',
					automaticLayout: true,
					wordWrap: 'on',
					scrollbar: {
						vertical: 'auto',
						horizontal: 'auto',
						verticalScrollbarSize: 10,
						horizontalScrollbarSize: 10,
					},
				});
			}
		} catch (e) {
			console.error('Failed to load Monaco Editor:', e);
			// Fallback to plain text display
		}

		loading = false;
	});

	onDestroy(() => {
		if (editor) {
			editor.dispose();
		}
	});

	// Update editor content when URL changes
	$effect(() => {
		if (editor && content !== null) {
			editor.setValue(content);
			if (monaco) {
				monaco.editor.setModelLanguage(editor.getModel(), language);
			}
		}
	});
</script>

<div class="w-full h-full flex flex-col bg-surface-primary">
	{#if loading}
		<div class="flex items-center justify-center h-full">
			<Spinner />
		</div>
	{:else if error}
		<div class="flex items-center justify-center h-full text-danger text-sm p-5">{error}</div>
	{:else if !editor && content !== null}
		<!-- Fallback: plain text display if Monaco fails to load -->
		<pre class="flex-1 m-0 p-4 overflow-auto bg-surface-primary text-text-primary font-mono text-[13px] leading-relaxed whitespace-pre-wrap break-words">{content}</pre>
	{/if}
	<div bind:this={containerElement} class="flex-1 w-full min-h-0 {loading || error || (!editor && content !== null) ? 'hidden' : ''}"></div>
</div>
