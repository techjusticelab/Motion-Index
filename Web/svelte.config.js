import adapter from '@sveltejs/adapter-auto';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const config = {
	preprocess: vitePreprocess(),
	kit: { 
		adapter: adapter(),
		// Temporarily disable CSRF protection for demo
		csrf: {
			checkOrigin: false
		}
	}
};

export default config;
