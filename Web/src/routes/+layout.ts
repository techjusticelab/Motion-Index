import { createSupabaseLoadClient } from '$lib/supabase-client'
import { browser } from '$app/environment'
import type { LayoutLoad } from './$types'

export const load: LayoutLoad = async ({ data, depends, fetch }) => {
  depends('supabase:auth')

  const supabase = createSupabaseLoadClient()

  /**
   * It's fine to use `getSession` here, because on the client, `getSession` is
   * safe, and on the server, it reads `session` from the `LayoutData`, which
   * safely checked the session using `safeGetSession`.
   */
  const {
    data: { session },
  } = await supabase.auth.getSession()

  /**
   * Use getUser() to get authenticated user data, which validates the JWT
   * by contacting the Supabase Auth server, making it more secure than
   * relying on the user object from getSession().
   */
  const {
    data: { user },
  } = await supabase.auth.getUser()

  return {
    ...data,
    supabase,
    session: session ?? data.session,
    user: user ?? data.user
  }
}