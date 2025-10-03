import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit(),
	],
	server: {
		allowedHosts: ['motionindex.techjusticelab.org'],
		cors: true,
		host: '0.0.0.0'  // Allow external access
	}
});