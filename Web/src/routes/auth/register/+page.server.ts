// /src/routes/auth/register/+page.server.ts
import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, locals }) => {
        const formData = await request.formData();
        const email = formData.get('email')?.toString();
        const password = formData.get('password')?.toString();
        const confirmPassword = formData.get('confirm-password')?.toString();

        if (!email || !password || !confirmPassword) {
            return fail(400, {
                error: 'All fields are required',
                email
            });
        }

        if (password !== confirmPassword) {
            return fail(400, {
                error: 'Passwords do not match',
                email
            });
        }

        if (password.length < 6) {
            return fail(400, {
                error: 'Password must be at least 6 characters',
                email
            });
        }

        try {
            const { data, error } = await locals.supabase.auth.signUp({
                email,
                password,
                options: {
                    emailRedirectTo: `${process.env.PUBLIC_SITE_URL || 'http://localhost:5173'}/auth/confirm`
                }
            });

            if (error) {
                return fail(400, {
                    error: error.message,
                    email
                });
            }

            // Redirect to login page with confirmation message
            return { 
                success: true,
                message: 'Please check your email to confirm your account'
            };
        } catch (err) {
            console.error('Registration error:', err);
            return fail(500, {
                error: 'An unexpected error occurred',
                email
            });
        }
    }
};
