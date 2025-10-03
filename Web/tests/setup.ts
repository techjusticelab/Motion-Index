import { beforeAll, afterEach } from 'vitest';
import { cleanup } from '@testing-library/svelte';

// Clean up after each test
afterEach(() => {
	cleanup();
});

// Mock environment variables for tests
beforeAll(() => {
	Object.defineProperty(globalThis, '__SVELTEKIT__', {
		value: true,
		writable: true
	});
});

// Mock SvelteKit modules
vi.mock('$app/environment', () => ({
	browser: false,
	dev: true,
	building: false,
	version: '1.0.0'
}));

vi.mock('$app/stores', () => ({
	page: {
		subscribe: vi.fn(() => vi.fn())
	},
	navigating: {
		subscribe: vi.fn(() => vi.fn())
	},
	updated: {
		subscribe: vi.fn(() => vi.fn())
	}
}));

vi.mock('$env/static/public', () => ({
	PUBLIC_SUPABASE_URL: 'https://test.supabase.co',
	PUBLIC_SUPABASE_ANON_KEY: 'test-anon-key',
	PUBLIC_API_URL: 'http://localhost:8003'
}));