/**
 * File type detection utilities for preview
 */

export type PreviewType = 'video' | 'audio' | 'image' | 'pdf' | 'code' | 'text' | 'unsupported';

const VIDEO_EXTENSIONS = ['mp4', 'webm', 'mkv', 'avi', 'mov', 'wmv', 'flv', 'm4v', 'ogv'];
const AUDIO_EXTENSIONS = ['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a', 'wma', 'opus'];
const IMAGE_EXTENSIONS = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp', 'ico', 'avif'];
const PDF_EXTENSIONS = ['pdf'];
const CODE_EXTENSIONS = [
	// JavaScript/TypeScript
	'js', 'jsx', 'ts', 'tsx', 'mjs', 'cjs',
	// Web
	'html', 'htm', 'css', 'scss', 'sass', 'less',
	// Data formats
	'json', 'yaml', 'yml', 'toml', 'xml',
	// Backend languages
	'py', 'go', 'rs', 'java', 'kt', 'scala', 'rb', 'php', 'cs', 'fs',
	// Systems languages
	'c', 'cpp', 'cc', 'cxx', 'h', 'hpp', 'hxx',
	// Shell/Scripts
	'sh', 'bash', 'zsh', 'fish', 'ps1', 'bat', 'cmd',
	// Config files
	'ini', 'conf', 'cfg', 'env',
	// Documentation
	'md', 'mdx', 'rst', 'tex',
	// Database
	'sql', 'graphql', 'gql',
	// Other
	'dockerfile', 'makefile', 'cmake', 'gradle', 'swift', 'r', 'lua', 'vim', 'asm'
];
const TEXT_EXTENSIONS = ['txt', 'log', 'csv', 'tsv', 'rtf'];

/**
 * Get the file extension from a filename
 */
export function getExtension(filename: string): string {
	const lastDot = filename.lastIndexOf('.');
	if (lastDot === -1) return '';
	return filename.slice(lastDot + 1).toLowerCase();
}

/**
 * Determine the preview type for a file
 */
export function getPreviewType(filename: string): PreviewType {
	const ext = getExtension(filename);
	
	// Handle special filenames without extensions
	const lowerName = filename.toLowerCase();
	if (['dockerfile', 'makefile', 'cmakelists.txt', 'gemfile', 'rakefile'].includes(lowerName)) {
		return 'code';
	}
	if (lowerName.startsWith('.') && !ext) {
		// Dotfiles like .gitignore, .env, etc.
		return 'code';
	}

	if (!ext) return 'unsupported';

	if (VIDEO_EXTENSIONS.includes(ext)) return 'video';
	if (AUDIO_EXTENSIONS.includes(ext)) return 'audio';
	if (IMAGE_EXTENSIONS.includes(ext)) return 'image';
	if (PDF_EXTENSIONS.includes(ext)) return 'pdf';
	if (CODE_EXTENSIONS.includes(ext)) return 'code';
	if (TEXT_EXTENSIONS.includes(ext)) return 'text';

	return 'unsupported';
}

/**
 * Get Monaco Editor language ID from file extension
 */
export function getMonacoLanguage(filename: string): string {
	const ext = getExtension(filename);
	const lowerName = filename.toLowerCase();

	// Special filenames
	if (lowerName === 'dockerfile') return 'dockerfile';
	if (lowerName === 'makefile' || lowerName === 'gnumakefile') return 'makefile';
	if (lowerName.endsWith('.gitignore') || lowerName.endsWith('.dockerignore')) return 'ignore';

	const languageMap: Record<string, string> = {
		// JavaScript/TypeScript
		js: 'javascript',
		jsx: 'javascript',
		mjs: 'javascript',
		cjs: 'javascript',
		ts: 'typescript',
		tsx: 'typescript',
		// Web
		html: 'html',
		htm: 'html',
		css: 'css',
		scss: 'scss',
		sass: 'scss',
		less: 'less',
		// Data formats
		json: 'json',
		yaml: 'yaml',
		yml: 'yaml',
		toml: 'ini',
		xml: 'xml',
		// Backend languages
		py: 'python',
		go: 'go',
		rs: 'rust',
		java: 'java',
		kt: 'kotlin',
		scala: 'scala',
		rb: 'ruby',
		php: 'php',
		cs: 'csharp',
		fs: 'fsharp',
		// Systems languages
		c: 'c',
		cpp: 'cpp',
		cc: 'cpp',
		cxx: 'cpp',
		h: 'c',
		hpp: 'cpp',
		hxx: 'cpp',
		// Shell/Scripts
		sh: 'shell',
		bash: 'shell',
		zsh: 'shell',
		fish: 'shell',
		ps1: 'powershell',
		bat: 'bat',
		cmd: 'bat',
		// Config files
		ini: 'ini',
		conf: 'ini',
		cfg: 'ini',
		env: 'dotenv',
		// Documentation
		md: 'markdown',
		mdx: 'markdown',
		rst: 'restructuredtext',
		tex: 'latex',
		// Database
		sql: 'sql',
		graphql: 'graphql',
		gql: 'graphql',
		// Other
		swift: 'swift',
		r: 'r',
		lua: 'lua',
		vim: 'vim',
		// Plain text
		txt: 'plaintext',
		log: 'plaintext',
		csv: 'plaintext',
		tsv: 'plaintext'
	};

	return languageMap[ext] || 'plaintext';
}

/**
 * Check if a file can be previewed
 */
export function canPreview(filename: string): boolean {
	return getPreviewType(filename) !== 'unsupported';
}

/**
 * Format file size for display
 */
export function formatSize(bytes: number): string {
	if (bytes === 0) return '0 B';
	const k = 1024;
	const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}
