// /src/routes/auth/login/+page.server.ts
import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ request, locals }) => {
        const formData = await request.formData();
        const email = formData.get('email')?.toString();
        const password = formData.get('password')?.toString();

        if (!email || !password) {
            return fail(400, {
                error: 'Email and password are required',
                email
            });
        }

        try {
            const { data, error } = await locals.supabase.auth.signInWithPassword({
                email,
                password
            });

            if (error) {
                return fail(400, {
                    error: error.message,
                    email
                });
            }

            // Redirect after successful login
            throw redirect(303, '/');
        } catch (err) {
            if (err instanceof Response) throw err;

            console.error('Login error:', err);
            return fail(500, {
                error: 'An unexpected error occurred',
                email
            });
        }
    }
};