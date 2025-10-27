import { writable } from 'svelte/store';
import type { User } from '$lib/types/User';

// Store to hold the authenticated user
export const user = writable<User | null>(null);

// Function to update the user store
export const setUser = (newUser: User | null) => {
  user.set(newUser);
}; 