import { SupabaseClient, Session } from '@supabase/supabase-js';
import { Database } from '$lib/database.types';

declare global {
	namespace App {
		interface Locals {
			supabase: SupabaseClient<Database>;
			getSession(): Promise<Session | null>;
		}
		interface PageData {
			session: Session | null;
			supabase: SupabaseClient<Database>;
		}
		// interface Error {}
		// interface Platform {}
	}
}

export { };