import { writable } from "svelte/store";

// Types
import type { Folder } from "../lib/types/Folder";

export const folders_store = writable<Array<Folder>>([]);

export const setFolders = (folders: Array<Folder>) => {
    folders_store.set(folders);
}

export const refreshFolders = async (user_id: number) => {
    try {
        const response = await fetch(`/api/folders/all/${user_id}`);
        if (!response.ok) {
            throw new Error("Failed to fetch folders");
        }
        const data: Array<Folder> = await response.json();
        folders_store.set(data);
    } catch (error) {
        console.error("Error refreshing folders:", error);
    }
}