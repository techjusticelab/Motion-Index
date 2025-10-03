import { writable, derived } from 'svelte/store';
import type { SupabaseClient, Session, User } from '@supabase/supabase-js';

export interface AuthState {
  session: Session | null;
  user: User | null;
  loading: boolean;
  error: string | null;
}

/**
 * Composable hook for managing authentication state
 * Provides consistent auth state management across components
 */
export function useAuth(supabase: SupabaseClient) {
  const authStore = writable<AuthState>({
    session: null,
    user: null,
    loading: true,
    error: null
  });

  // Derived stores for convenient access
  const isAuthenticated = derived(authStore, $auth => !!$auth.user);
  const user = derived(authStore, $auth => $auth.user);
  const session = derived(authStore, $auth => $auth.session);

  // Initialize auth state
  async function initialize(): Promise<void> {
    try {
      authStore.update(state => ({ ...state, loading: true, error: null }));
      
      const { data: { session }, error } = await supabase.auth.getSession();
      
      if (error) {
        authStore.update(state => ({ ...state, error: error.message, loading: false }));
        return;
      }

      authStore.update(state => ({
        ...state,
        session,
        user: session?.user || null,
        loading: false,
        error: null
      }));
    } catch (err) {
      authStore.update(state => ({
        ...state,
        error: err instanceof Error ? err.message : 'Authentication error',
        loading: false
      }));
    }
  }

  // Listen for auth state changes
  supabase.auth.onAuthStateChange((event, session) => {
    authStore.update(state => ({
      ...state,
      session,
      user: session?.user || null,
      loading: false,
      error: null
    }));
  });

  return {
    subscribe: authStore.subscribe,
    isAuthenticated,
    user,
    session,

    // Initialize auth state
    initialize,

    // Sign in with email and password
    signIn: async (email: string, password: string): Promise<{ success: boolean; error?: string; data?: any }> => {
      try {
        authStore.update(state => ({ ...state, loading: true, error: null }));
        
        const { data, error } = await supabase.auth.signInWithPassword({
          email,
          password
        });

        if (error) {
          authStore.update(state => ({ ...state, error: error.message, loading: false }));
          return { success: false, error: error.message };
        }

        return { success: true, data };
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Sign in failed';
        authStore.update(state => ({ ...state, error: errorMessage, loading: false }));
        return { success: false, error: errorMessage };
      }
    },

    // Sign up with email and password
    signUp: async (email: string, password: string, options?: { data?: object }): Promise<{ success: boolean; error?: string; data?: any }> => {
      try {
        authStore.update(state => ({ ...state, loading: true, error: null }));
        
        const { data, error } = await supabase.auth.signUp({
          email,
          password,
          options
        });

        if (error) {
          authStore.update(state => ({ ...state, error: error.message, loading: false }));
          return { success: false, error: error.message };
        }

        return { success: true, data };
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Sign up failed';
        authStore.update(state => ({ ...state, error: errorMessage, loading: false }));
        return { success: false, error: errorMessage };
      }
    },

    // Sign out
    signOut: async (): Promise<{ success: boolean; error?: string }> => {
      try {
        authStore.update(state => ({ ...state, loading: true, error: null }));
        
        const { error } = await supabase.auth.signOut();

        if (error) {
          authStore.update(state => ({ ...state, error: error.message, loading: false }));
          return { success: false, error: error.message };
        }

        authStore.update(state => ({
          session: null,
          user: null,
          loading: false,
          error: null
        }));

        return { success: true };
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Sign out failed';
        authStore.update(state => ({ ...state, error: errorMessage, loading: false }));
        return { success: false, error: errorMessage };
      }
    },

    // Reset password
    resetPassword: async (email: string): Promise<{ success: boolean; error?: string; data?: any }> => {
      try {
        const { data, error } = await supabase.auth.resetPasswordForEmail(email);

        if (error) {
          return { success: false, error: error.message };
        }

        return { success: true, data };
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Password reset failed';
        return { success: false, error: errorMessage };
      }
    },

    // Update password
    updatePassword: async (password: string): Promise<{ success: boolean; error?: string; data?: any }> => {
      try {
        authStore.update(state => ({ ...state, loading: true, error: null }));
        
        const { data, error } = await supabase.auth.updateUser({ password });

        if (error) {
          authStore.update(state => ({ ...state, error: error.message, loading: false }));
          return { success: false, error: error.message };
        }

        authStore.update(state => ({ ...state, loading: false }));
        return { success: true, data };
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Password update failed';
        authStore.update(state => ({ ...state, error: errorMessage, loading: false }));
        return { success: false, error: errorMessage };
      }
    }
  };
}