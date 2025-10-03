import { describe, it, expect } from 'vitest';
import { 
	formatDate, 
	formatFileSize, 
	truncateText, 
	debounce, 
	generateId, 
	filterOptions 
} from '$lib/utils';

describe('Utils', () => {
	describe('formatDate', () => {
		it('formats valid date string', () => {
			const result = formatDate('2023-12-25');
			expect(result).toMatch(/Dec 25, 2023/);
		});

		it('returns N/A for invalid input', () => {
			expect(formatDate('')).toBe('N/A');
			expect(formatDate(null as any)).toBe('N/A');
		});
	});

	describe('formatFileSize', () => {
		it('formats bytes correctly', () => {
			expect(formatFileSize(0)).toBe('0 B');
			expect(formatFileSize(1024)).toBe('1 KB');
			expect(formatFileSize(1048576)).toBe('1 MB');
			expect(formatFileSize(1073741824)).toBe('1 GB');
		});
	});

	describe('truncateText', () => {
		it('truncates long text', () => {
			const text = 'This is a very long text that should be truncated';
			const result = truncateText(text, 20);
			expect(result).toBe('This is a very long ...');
		});

		it('returns original text if shorter than max length', () => {
			const text = 'Short text';
			const result = truncateText(text, 20);
			expect(result).toBe('Short text');
		});
	});

	describe('generateId', () => {
		it('generates unique IDs', () => {
			const id1 = generateId();
			const id2 = generateId();
			expect(id1).not.toBe(id2);
			expect(typeof id1).toBe('string');
			expect(id1.length).toBeGreaterThan(0);
		});
	});

	describe('filterOptions', () => {
		const options = ['Apple', 'Banana', 'Cherry', 'Date'];

		it('filters options by search input', () => {
			const result = filterOptions(options, 'a');
			expect(result).toEqual(['Apple', 'Banana', 'Date']);
		});

		it('returns all options for empty search', () => {
			const result = filterOptions(options, '');
			expect(result).toEqual(options);
		});

		it('is case insensitive', () => {
			const result = filterOptions(options, 'APPLE');
			expect(result).toEqual(['Apple']);
		});
	});

	describe('debounce', () => {
		it('debounces function calls', async () => {
			let callCount = 0;
			const fn = () => callCount++;
			const debouncedFn = debounce(fn, 100);

			debouncedFn();
			debouncedFn();
			debouncedFn();

			expect(callCount).toBe(0);

			await new Promise(resolve => setTimeout(resolve, 150));
			expect(callCount).toBe(1);
		});
	});
});