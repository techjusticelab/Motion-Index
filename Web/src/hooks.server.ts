import { PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY } from '$env/static/public';
import { redirect, type Handle } from '@sveltejs/kit';
import { createClient } from '@supabase/supabase-js';

export const handle: Handle = async ({ event, resolve }) => {
    event.locals.supabase = createClient(
        PUBLIC_SUPABASE_URL,
        PUBLIC_SUPABASE_ANON_KEY,
        {
            auth: {
                autoRefreshToken: true,
                persistSession: true,
                detectSessionInUrl: true
            }
        }
    );

    event.locals.getSession = async () => {
        const {
            data: { session }
        } = await event.locals.supabase.auth.getSession();
        return session;
    };

    // Protected routes - redirect to login if not authenticated
    const session = await event.locals.getSession();
    const protectedRoutes = ['/account', '/upload', '/documents'];
    const isProtectedRoute = protectedRoutes.some(route =>
        event.url.pathname.startsWith(route)
    );

    if (isProtectedRoute && !session) {
        throw redirect(303, '/auth/login?redirectTo=' + event.url.pathname);
    }

    return resolve(event, {
        filterSerializedResponseHeaders(name) {
            return name === 'content-range';
        }
    });
};