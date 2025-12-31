<script lang="ts">
	/**
	 * CodePreview - Code/text viewer with syntax highlighting using Monaco Editor
	 */
	import { onMount, onDestroy } from 'svelte';
	import { getMonacoLanguage } from '$lib/utils/fileTypes';
	import { getFileContent } from '$lib/api/files';

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
					// Return a basic worker - Monaco will work without specialized workers
					// but won't have full language features
					return new Worker(
						URL.createObjectURL(
							new Blob(
								[`self.onmessage = function() {}`],
								{ type: 'text/javascript' }
							)
						)
					);
				}
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
						'editor.lineHighlightBackground': '#2a2a2a'
					}
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
						horizontalScrollbarSize: 10
					}
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

<div class="code-preview">
	{#if loading}
		<div class="loading">Loading...</div>
	{:else if error}
		<div class="error-message">{error}</div>
	{:else if !editor && content !== null}
		<!-- Fallback: plain text display if Monaco fails to load -->
		<pre class="fallback-code">{content}</pre>
	{/if}
	<div 
		bind:this={containerElement} 
		class="editor-container"
		class:hidden={loading || error || (!editor && content !== null)}
	></div>
</div>

<style>
	.code-preview {
		width: 100%;
		height: 100%;
		display: flex;
		flex-direction: column;
		background: #1e1e1e;
	}

	.editor-container {
		flex: 1;
		width: 100%;
		min-height: 0;
	}

	.editor-container.hidden {
		display: none;
	}

	.loading {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 100%;
		color: #888;
		font-size: 14px;
	}

	.error-message {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 100%;
		color: #f87171;
		font-size: 14px;
		padding: 20px;
	}

	.fallback-code {
		flex: 1;
		margin: 0;
		padding: 16px;
		overflow: auto;
		background: #1e1e1e;
		color: #d4d4d4;
		font-family: 'Fira Code', 'Cascadia Code', Consolas, monospace;
		font-size: 13px;
		line-height: 1.5;
		white-space: pre-wrap;
		word-wrap: break-word;
	}
</style>
