import type { LayoutServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import type { RequestEvent } from '@sveltejs/kit';

export const load: LayoutServerLoad = async ({ locals, url }: { locals: { getSession: () => Promise<any> }, url: URL }) => {
    const { getSession } = locals;
    const session = await getSession();

    // Redirect to login if not authenticated and not on auth pages


    return {
        session
    };
};