import { PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY } from '$env/static/public';
import { createClient } from '@supabase/supabase-js';

export const load = async ({ fetch, data }) => {
    const supabase = createClient(
        PUBLIC_SUPABASE_URL,
        PUBLIC_SUPABASE_ANON_KEY,
        {
            global: {
                fetch: fetch
            },
            auth: {
                persistSession: true
            }
        }
    );

    return {
        supabase,
        session: data.session
    };
};