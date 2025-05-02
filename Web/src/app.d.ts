import type { Session, SupabaseClient, User } from '@supabase/supabase-js'
import type { Database } from './database.types' // import generated types

declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			supabase: SupabaseClient<Database>
			safeGetSession: () => Promise<{ session: Session | null; user: User | null }>
			getSession: () => Promise<Session | null> // Added for backward compatibility
			session: Session | null
			user: User | null
		}
		interface PageData {
			session: Session | null
			supabase: SupabaseClient<Database> // Added for components that need supabase client
		}
		// interface PageState {}
		// interface Platform {}
	}
}

export { }