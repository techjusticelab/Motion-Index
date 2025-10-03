import { writable } from 'svelte/store';

export interface LoadingState {
  isLoading: boolean;
  error: string | null;
  success: boolean;
}

/**
 * Composable hook for managing loading states
 * Provides consistent loading, error, and success state management
 */
export function useLoading(initialState: Partial<LoadingState> = {}) {
  const { subscribe, set, update } = writable<LoadingState>({
    isLoading: false,
    error: null,
    success: false,
    ...initialState
  });

  return {
    subscribe,
    
    // Start loading
    start: () => update(state => ({
      ...state,
      isLoading: true,
      error: null,
      success: false
    })),
    
    // Complete successfully
    success: (successMessage?: string) => update(state => ({
      ...state,
      isLoading: false,
      error: null,
      success: true
    })),
    
    // Complete with error
    error: (errorMessage: string) => update(state => ({
      ...state,
      isLoading: false,
      error: errorMessage,
      success: false
    })),
    
    // Reset to initial state
    reset: () => set({
      isLoading: false,
      error: null,
      success: false
    }),
    
    // Async wrapper function
    async withLoading<T>(asyncFn: () => Promise<T>): Promise<T | null> {
      try {
        this.start();
        const result = await asyncFn();
        this.success();
        return result;
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'An error occurred';
        this.error(errorMessage);
        return null;
      }
    }
  };
}