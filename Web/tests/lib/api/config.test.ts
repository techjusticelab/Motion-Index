import { describe, it, expect, vi } from 'vitest';
import { getAuthHeaders, handleApiError } from '$lib/api/config';

describe('API Config', () => {
	describe('getAuthHeaders', () => {
		it('returns headers with auth token when session provided', async () => {
			const mockSession = {
				access_token: 'test-token'
			};

			const headers = await getAuthHeaders(mockSession);
			
			expect(headers).toEqual({
				'Authorization': 'Bearer test-token',
				'Content-Type': 'application/json'
			});
		});

		it('returns default headers when no session', async () => {
			const headers = await getAuthHeaders();
			
			expect(headers).toEqual({
				'Content-Type': 'application/json'
			});
		});
	});

	describe('handleApiError', () => {
		it('throws error with response message', () => {
			const error = {
				response: {
					data: {
						message: 'API Error Message'
					}
				}
			};

			expect(() => {
				handleApiError(error, 'test operation');
			}).toThrow('API Error Message');
		});

		it('throws error with error message', () => {
			const error = {
				message: 'Network Error'
			};

			expect(() => {
				handleApiError(error, 'test operation');
			}).toThrow('Network Error');
		});

		it('throws generic error when no specific message', () => {
			const error = {};

			expect(() => {
				handleApiError(error, 'test operation');
			}).toThrow('Failed to test operation');
		});
	});
});