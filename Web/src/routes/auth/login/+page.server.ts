import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, locals, url }) => {
        const formData = await request.formData();
        const email = formData.get('email')?.toString();
        const password = formData.get('password')?.toString();
        const redirectTo = url.searchParams.get('redirectTo') || '/';

        // Validate form data
        if (!email || !password) {
            return fail(400, {
                error: 'Email and password are required',
                email
            });
        }

        // Basic email validation
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email)) {
            return fail(400, {
                error: 'Please enter a valid email address',
                email
            });
        }

        try {
            const { data, error } = await locals.supabase.auth.signInWithPassword({
                email,
                password
            });

            if (error) {
                console.error('Supabase auth error:', error);
                
                // Provide user-friendly error messages based on error type
                let userMessage = error.message;
                
                // Authentication errors
                if (error.message.includes('Invalid login credentials')) {
                    userMessage = 'Invalid email or password. Please check your credentials and try again.';
                } else if (error.message.includes('Email not confirmed')) {
                    userMessage = 'Please check your email and click the confirmation link before signing in.';
                } else if (error.message.includes('Too many requests')) {
                    userMessage = 'Too many login attempts. Please wait a few minutes before trying again.';
                } 
                // Network and connection errors
                else if (error.message.includes('Network request failed') || 
                         error.message.includes('fetch') || 
                         error.message.includes('ERR_BLOCKED_BY_CLIENT') ||
                         error.message.includes('ERR_NETWORK') ||
                         error.message.includes('ERR_INTERNET_DISCONNECTED')) {
                    userMessage = 'Connection failed. This may be caused by:\n• Ad blockers blocking Supabase\n• Firewall or network restrictions\n• Poor internet connection\n\nTry disabling ad blockers or check your network.';
                }
                // Timeout errors
                else if (error.message.includes('timeout') || error.message.includes('TIMEOUT')) {
                    userMessage = 'Request timed out. Please check your internet connection and try again.';
                }
                // CORS or domain errors
                else if (error.message.includes('CORS') || error.message.includes('origin')) {
                    userMessage = 'Authentication service blocked. Please disable ad blockers and privacy extensions temporarily.';
                }
                
                return fail(400, {
                    error: userMessage,
                    email
                });
            }

            if (!data?.session) {
                return fail(400, {
                    error: 'Login failed - no session created. Please try again or contact support.',
                    email
                });
            }

            // Successful login - redirect to requested page or home
            redirect(303, redirectTo);
        } catch (err) {
            // Handle redirect responses (these are normal)
            if (err instanceof Response) {
                throw err;
            }

            console.error('Unexpected login error:', err);
            
            // Enhanced error detection for different types of failures
            if (err instanceof Error) {
                const errorMessage = err.message.toLowerCase();
                
                // Network connectivity issues
                if (errorMessage.includes('fetch') || 
                    errorMessage.includes('network') || 
                    errorMessage.includes('connection') ||
                    errorMessage.includes('timeout') ||
                    errorMessage.includes('blocked') ||
                    errorMessage.includes('cors')) {
                    
                    return fail(500, {
                        error: 'Connection to authentication service failed. Please:\n• Check your internet connection\n• Disable ad blockers (uBlock Origin, AdBlock, etc.)\n• Disable privacy extensions temporarily\n• Try a different browser or incognito mode',
                        email
                    });
                }
                
                // Browser security restrictions
                if (errorMessage.includes('security') || 
                    errorMessage.includes('same-origin') ||
                    errorMessage.includes('mixed content')) {
                    
                    return fail(500, {
                        error: 'Browser security settings are blocking authentication. Try using incognito mode or a different browser.',
                        email
                    });
                }
            }
            
            // Generic fallback error
            return fail(500, {
                error: 'An unexpected error occurred during login. Please try again or contact support if the problem persists.',
                email
            });
        }
    }
};