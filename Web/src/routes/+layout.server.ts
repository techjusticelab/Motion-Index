import type { LayoutServerLoad } from './$types'

export const load: LayoutServerLoad = async ({ locals, cookies }) => {
    // With Supabase, the session is usually accessed from locals.supabase.auth.getSession()
    // or directly available as locals.getSession or locals.session

    // Try this approach:
    const session = locals.getSession ? await locals.getSession() : null;
    // Or this:
    // const session = locals.supabase?.auth.getSession ? await locals.supabase.auth.getSession() : null;

    return {
        session,
        cookies: cookies.getAll(),
    };
}