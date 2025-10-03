import { createServerClient } from '@supabase/ssr'
import { type Handle, redirect } from '@sveltejs/kit'
import { sequence } from '@sveltejs/kit/hooks'
import { PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY } from '$env/static/public'

// Environment variable validation
function validateEnvironment(): void {
  const requiredEnvVars = {
    PUBLIC_SUPABASE_URL,
    PUBLIC_SUPABASE_ANON_KEY
  };

  const missingVars = Object.entries(requiredEnvVars)
    .filter(([name, value]) => !value)
    .map(([name]) => name);

  if (missingVars.length > 0) {
    throw new Error(
      `Missing required environment variables: ${missingVars.join(', ')}`
    );
  }

  // Validate Supabase URL format
  if (!PUBLIC_SUPABASE_URL.startsWith('https://') && !PUBLIC_SUPABASE_URL.includes('localhost')) {
    throw new Error('PUBLIC_SUPABASE_URL must be a valid HTTPS URL or localhost');
  }

  console.log('✅ Environment variables validated successfully');
}

// Validate environment on startup
validateEnvironment();

const supabase: Handle = async ({ event, resolve }) => {
  event.locals.supabase = createServerClient(PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY, {
    cookies: {
      getAll: () => event.cookies.getAll(),
      setAll: (cookiesToSet) => {
        cookiesToSet.forEach(({ name, value, options }) => {
          event.cookies.set(name, value, { ...options, path: '/' })
        })
      },
    },
  })

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

  return resolve(event, {
    filterSerializedResponseHeaders(name) {
      return name === 'content-range' || name === 'x-supabase-api-version'
    },
  })
}

const authGuard: Handle = async ({ event, resolve }) => {
  const { session } = await event.locals.safeGetSession()
  event.locals.session = session
  event.locals.user = session?.user ?? null

  if (!event.locals.session && event.url.pathname.startsWith('/account')) {
    redirect(303, '/auth/login')
  }

  if (!event.locals.session && event.url.pathname.startsWith('/upload')) {
    redirect(303, '/auth/login')
  }

  if (event.locals.session && event.url.pathname === '/auth/login') {
    redirect(303, '/account')
  }

  return resolve(event)
}

export const handle: Handle = sequence(supabase, authGuard)