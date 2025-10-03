import { writable } from 'svelte/store';

export interface ModalState {
  isOpen: boolean;
  title?: string;
  content?: any;
  data?: any;
}

/**
 * Composable hook for managing modal state
 * Provides consistent modal open/close functionality with data passing
 */
export function useModal(initialState: Partial<ModalState> = {}) {
  const { subscribe, set, update } = writable<ModalState>({
    isOpen: false,
    title: '',
    content: null,
    data: null,
    ...initialState
  });

  return {
    subscribe,
    
    // Open modal with optional data
    open: (options: Partial<ModalState> = {}) => update(state => ({
      ...state,
      isOpen: true,
      ...options
    })),
    
    // Close modal and optionally clear data
    close: (clearData = true) => update(state => ({
      ...state,
      isOpen: false,
      ...(clearData && { title: '', content: null, data: null })
    })),
    
    // Toggle modal state
    toggle: () => update(state => ({
      ...state,
      isOpen: !state.isOpen
    })),
    
    // Update modal data without changing open state
    updateData: (data: any) => update(state => ({
      ...state,
      data
    })),
    
    // Update modal title
    updateTitle: (title: string) => update(state => ({
      ...state,
      title
    })),
    
    // Reset to initial state
    reset: () => set({
      isOpen: false,
      title: '',
      content: null,
      data: null
    })
  };
}