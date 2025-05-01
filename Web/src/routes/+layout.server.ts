import type { LayoutServerLoad } from './$types'

export const load: LayoutServerLoad = async({ locals: { safeGetSession }: { safeGetSession: () => Promise<{ session: any }> }, cookies }) => {
    const { session } = await safeGetSession()
    return {
        session,
        cookies: cookies.getAll(),
    }
}