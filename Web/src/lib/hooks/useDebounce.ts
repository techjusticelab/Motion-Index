import { writable, derived } from 'svelte/store';

/**
 * Composable hook for debouncing values
 * Useful for search inputs and API calls that should be delayed
 */
export function useDebounce<T>(initialValue: T, delay: number = 300) {
  const value = writable<T>(initialValue);
  const debouncedValue = writable<T>(initialValue);

  let timeoutId: ReturnType<typeof setTimeout> | null = null;

  // Subscribe to value changes and debounce them
  value.subscribe((newValue) => {
    if (timeoutId) {
      clearTimeout(timeoutId);
    }

    timeoutId = setTimeout(() => {
      debouncedValue.set(newValue);
    }, delay);
  });

  return {
    // Immediate value (changes instantly)
    value,
    
    // Debounced value (changes after delay)
    debouncedValue,
    
    // Set the value
    setValue: (newValue: T) => value.set(newValue),
    
    // Cancel pending debounce
    cancel: () => {
      if (timeoutId) {
        clearTimeout(timeoutId);
        timeoutId = null;
      }
    },
    
    // Trigger immediate update to debounced value
    flush: () => {
      if (timeoutId) {
        clearTimeout(timeoutId);
        timeoutId = null;
      }
      value.update(val => {
        debouncedValue.set(val);
        return val;
      });
    }
  };
}

/**
 * Simple debounce function for standalone use
 */
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout> | null = null;

  return (...args: Parameters<T>) => {
    if (timeoutId) {
      clearTimeout(timeoutId);
    }

    timeoutId = setTimeout(() => {
      func(...args);
    }, delay);
  };
}