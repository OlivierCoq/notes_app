import { writable } from "svelte/store";

// Types
import type { Folder } from "../lib/types/Folder";

export const folders_store = writable<Array<Folder>>([]);

export const setFolders = (folders: Array<Folder>) => {
    folders_store.set(folders);
}