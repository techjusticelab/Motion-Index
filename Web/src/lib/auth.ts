// Simple, reliable authentication management
import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import Cookies from 'js-cookie';

// Define the user type
export interface User {
  id: string;
  email: string;
  name?: string;
}

// Create a store for the authentication state
export const isAuthenticated = writable<boolean>(false);
export const currentUser = writable<User | null>(null);
export const authToken = writable<string | null>(null);

// Cookie settings
const COOKIE_OPTIONS = {
  expires: 7, // 7 days
  path: '/',
  secure: browser && window.location.protocol === 'https:',
  sameSite: 'strict' as const
};

// Initialize auth state from cookies
export function initAuth() {
  if (browser) {
    try {
      // Clear any existing Supabase auth tokens from localStorage
      const keys = Object.keys(localStorage);
      keys.forEach(key => {
        if (key.startsWith('sb-')) {
          localStorage.removeItem(key);
        }
      });
      
      // Check for our auth cookie
      const authCookie = Cookies.get('motion-index-auth');
      const userCookie = Cookies.get('motion-index-user');
      
      if (authCookie && userCookie) {
        try {
          const user = JSON.parse(userCookie);
          isAuthenticated.set(true);
          currentUser.set(user);
          authToken.set(authCookie);
          console.log('User authenticated from cookies');
        } catch (e) {
          console.error('Failed to parse user data:', e);
          clearAuth();
        }
      } else {
        clearAuth();
      }
    } catch (e) {
      console.error('Error initializing auth:', e);
      clearAuth();
    }
  }
}

// Set authentication state
export function setAuth(user: User, token: string) {
  if (browser) {
    try {
      // Store in cookies
      Cookies.set('motion-index-auth', token, COOKIE_OPTIONS);
      Cookies.set('motion-index-user', JSON.stringify(user), COOKIE_OPTIONS);
      
      // Update stores
      isAuthenticated.set(true);
      currentUser.set(user);
      authToken.set(token);
      
      console.log('User authenticated and cookies set');
    } catch (e) {
      console.error('Error setting auth:', e);
    }
  }
}

// Clear authentication state
export function clearAuth() {
  if (browser) {
    // Clear cookies
    Cookies.remove('motion-index-auth', { path: '/' });
    Cookies.remove('motion-index-user', { path: '/' });
    
    // Clear any Supabase auth tokens from localStorage
    const keys = Object.keys(localStorage);
    keys.forEach(key => {
      if (key.startsWith('sb-')) {
        localStorage.removeItem(key);
      }
    });
    
    // Update stores
    isAuthenticated.set(false);
    currentUser.set(null);
    authToken.set(null);
    
    console.log('User logged out and cookies cleared');
  }
}

// Get current auth token
export function getAuthToken(): string | null {
  if (browser) {
    return Cookies.get('motion-index-auth') || null;
  }
  return null;
}

// Initialize auth on import
if (browser) {
  initAuth();
}
