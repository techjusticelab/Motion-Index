import { SupabaseClient, Session } from '@supabase/supabase-js';

declare global {
	namespace App {
		interface Locals {
			supabase: SupabaseClient;
			getSession(): Promise<Session | null>;
		}
		interface PageData {
			session: Session | null;
			supabase: SupabaseClient;
		}
		// interface Error {}
		// interface Platform {}
	}

	// Add environment variable typing
	interface ImportMetaEnv {
		VITE_SUPABASE_URL: string;
		VITE_SUPABASE_ANON_KEY: string;
		// Add any other environment variables you need
	}
}

export { };