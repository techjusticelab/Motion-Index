import { writable } from 'svelte/store';
import type { User } from '@supabase/supabase-js';

export const isLoading = writable(true);

export const user = writable<User | null>(null);

export const isAuthenticated = writable<boolean>(false);

export function updateUserStore(userData: User | null) {
    user.set(userData);
    isAuthenticated.set(!!userData);
}

export function initAuthStore(initialUser: User | null) {
    updateUserStore(initialUser);
}