import { createServerClient } from '@supabase/ssr'
import { type Handle, redirect } from '@sveltejs/kit'
import { sequence } from '@sveltejs/kit/hooks'

// Use environment variables directly for now to avoid import errors
const PUBLIC_SUPABASE_URL = process.env.PUBLIC_SUPABASE_URL || ''
const PUBLIC_SUPABASE_ANON_KEY = process.env.PUBLIC_SUPABASE_ANON_KEY || ''

const supabase: Handle = async ({ event, resolve }) => {
    /**
     * Creates a Supabase client specific to this server request.
     * The Supabase client gets the Auth token from the request cookies.
     */
    event.locals.supabase = createServerClient(PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY, {
        cookies: {
            getAll: () => event.cookies.getAll(),
            /**
             * SvelteKit's cookies API requires `path` to be explicitly set in
             * the cookie options. Setting `path` to `/` replicates previous/
             * standard behavior.
             */
            setAll: (cookiesToSet) => {
                cookiesToSet.forEach(({ name, value, options }) => {
                    event.cookies.set(name, value, { ...options, path: '/' })
                })
            },
        },
        global: {
            headers: {
                'X-Client-Info': `sveltekit-supabase@${new Date().toISOString()}`
            }
        }
    })

    /**
     * Unlike `supabase.auth.getSession()`, which returns the session _without_
     * validating the JWT, this function also calls `getUser()` to validate the
     * JWT before returning the session.
     */
    event.locals.safeGetSession = async () => {
        const {
            data: { session },
        } = await event.locals.supabase.auth.getSession()
        if (!session) {
            return { session: null, user: null }
        }

        const {
            data: { user },
            error,
        } = await event.locals.supabase.auth.getUser()
        if (error) {
            // JWT validation has failed
            return { session: null, user: null }
        }

        return { session, user }
    }

    // For backward compatibility with existing code
    event.locals.getSession = async () => {
        const { session } = await event.locals.safeGetSession();
        return session;
    };

    return resolve(event, {
        filterSerializedResponseHeaders(name) {
            /**
             * Supabase libraries use the `content-range` and `x-supabase-api-version`
             * headers, so we need to tell SvelteKit to pass it through.
             */
            return name === 'content-range' || name === 'x-supabase-api-version'
        },
    })
}

const authGuard: Handle = async ({ event, resolve }) => {
    // Get the session and user from the event
    const { session, user } = await event.locals.safeGetSession()
    
    // Store them in event.locals for use in routes
    event.locals.session = session
    event.locals.user = user
    
    // Define all protected routes
    const protectedRoutes = [
        '/account',
        '/upload',
        '/documents',
        '/cases',
        '/private'
    ];

    // Check if current path is a protected route (exact match or starts with the route)
    const isProtectedRoute = protectedRoutes.some(route => {
        // Check if the pathname exactly matches the route or starts with the route followed by a slash
        return event.url.pathname === route || 
               event.url.pathname.startsWith(`${route}/`);
    });

    // For API routes, return 401 instead of redirecting
    const isApiRoute = event.url.pathname.startsWith('/api/');

    // If user is trying to access a protected route, check auth cookie
    if (isProtectedRoute) {
        // Check for our custom auth cookie
        const authCookie = event.cookies.get('motion-index-auth');
        
        // If no auth cookie and no session, redirect to login
        if (!authCookie && !session) {
            console.log(`Redirecting unauthenticated user from ${event.url.pathname} to login`);
            
            if (isApiRoute) {
                // For API routes, return a 401 response
                return new Response(JSON.stringify({ error: 'Unauthorized' }), {
                    status: 401,
                    headers: { 'Content-Type': 'application/json' }
                });
            }

            // For regular routes, redirect to login
            const redirectUrl = encodeURIComponent(event.url.pathname + event.url.search);
            throw redirect(303, `/auth/login?redirectTo=${redirectUrl}`);
        }
    }

    // Redirect away from auth pages if already logged in (except logout)
    const isAuthRoute = event.url.pathname.startsWith('/auth');
    if (session && isAuthRoute && !event.url.pathname.includes('logout')) {
        throw redirect(303, '/account');
    }

    const response = await resolve(event);

    // Add security headers
    const securityHeaders = {
        'X-Frame-Options': 'SAMEORIGIN',
        'X-Content-Type-Options': 'nosniff',
        'Referrer-Policy': 'strict-origin-when-cross-origin'
    };

    // Add security headers to the response
    Object.entries(securityHeaders).forEach(([header, value]) => {
        response.headers.set(header, value);
    });

    return response;
}

export const handle: Handle = sequence(supabase, authGuard)