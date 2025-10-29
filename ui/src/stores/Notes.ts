import { writable } from "svelte/store";

// Types
import type { Note } from "../lib/types/Note";


export const notes_store = writable<Array<Note>>([]);


export const setNotes = (notes: Array<Note>) => {
    notes_store.set(notes);
}

export const refreshNotes = async (user_id: number) => {
    try {
        const response = await fetch(`/api/notes/all/${user_id}`);
        if (!response.ok) {
            throw new Error("Failed to fetch notes");
        }
        const data: Array<Note> = await response.json();
        notes_store.set(data);
    } catch (error) {
        console.error("Error refreshing notes:", error);
    }
}