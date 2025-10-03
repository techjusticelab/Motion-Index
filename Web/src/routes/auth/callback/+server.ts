import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ url, locals: { supabase } }) => {
	const code = url.searchParams.get('code');
	const next = url.searchParams.get('next') ?? '/';
	const error = url.searchParams.get('error');
	const error_description = url.searchParams.get('error_description');

	// Handle OAuth errors
	if (error) {
		console.error('OAuth error:', error, error_description);
		const errorMessage = encodeURIComponent(error_description || 'Authentication failed');
		redirect(303, `/auth/login?error=${errorMessage}`);
	}

	// Handle PKCE code exchange
	if (code) {
		try {
			const { data, error: exchangeError } = await supabase.auth.exchangeCodeForSession(code);
			
			if (exchangeError) {
				console.error('Code exchange error:', exchangeError);
				const errorMessage = encodeURIComponent(exchangeError.message || 'Failed to complete authentication');
				redirect(303, `/auth/login?error=${errorMessage}`);
			}

			if (data?.session) {
				console.log('Authentication successful, redirecting to:', next);
				redirect(303, next);
			}
		} catch (err) {
			console.error('Unexpected error during code exchange:', err);
			redirect(303, '/auth/login?error=Authentication%20failed');
		}
	}

	// Handle email confirmation
	const token_hash = url.searchParams.get('token_hash');
	const type = url.searchParams.get('type');

	if (token_hash && type) {
		try {
			const { data, error: verifyError } = await supabase.auth.verifyOtp({
				type: type as any,
				token_hash
			});

			if (verifyError) {
				console.error('OTP verification error:', verifyError);
				const errorMessage = encodeURIComponent(verifyError.message || 'Failed to verify email');
				redirect(303, `/auth/login?error=${errorMessage}`);
			}

			if (data?.session) {
				console.log('Email verification successful');
				redirect(303, '/auth/login?message=Email%20verified%20successfully');
			}
		} catch (err) {
			console.error('Unexpected error during email verification:', err);
			redirect(303, '/auth/login?error=Email%20verification%20failed');
		}
	}

	// Fallback redirect
	console.warn('No valid authentication parameters found in callback');
	redirect(303, '/auth/login?error=Invalid%20authentication%20callback');
};