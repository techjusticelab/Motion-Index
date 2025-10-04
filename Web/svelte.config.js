import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const config = {
	preprocess: vitePreprocess(),
	kit: { 
		adapter: adapter({
			out: 'build',
			precompress: true,
			envPrefix: ''
		}),
		// Temporarily disable CSRF protection for demo
		csrf: {
			checkOrigin: false
		}
	}
};

export default config;
