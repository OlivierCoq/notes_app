import { writable } from "svelte/store";

// Types
import type { Note } from "../lib/types/Note";


export const notes_store = writable<Array<Note>>([]);


export const setNotes = (notes: Array<Note>) => {
    notes_store.set(notes);
}
