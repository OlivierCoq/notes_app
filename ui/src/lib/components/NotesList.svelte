<script lang="ts">
	// Props
	let { notes, select_note, user, folders, onPostMove } = $props<{
		notes: Note[];
		select_note: (note: Note) => void;
		user: User | null;
		folders: Folder[] | null;
		onPostMove: () => Promise<void>;
	}>();

	// Imports
	import { AccordionItem, Accordion } from 'flowbite-svelte';
	//   Svelte
	//   Components
	import NoteSelector from '$lib/components/NoteSelector.svelte';
	import NoteFolder from '$lib/components/NoteFolder.svelte';
	//   Icons
	import { FolderPlusOutline, FolderSolid } from 'flowbite-svelte-icons';

	// Types
	import type { Note } from '$lib/types/Note';
	import type { Folder } from '$lib/types/Folder';
	import type { User } from '$lib/types/User';

	// Stores
	import { folders_store } from '../../stores/Folders';
	import { notes_store } from '../../stores/Notes';

	// console.log('Foldsers store in NotesList:', $folders_store);

	// State
	let notes_list_state = $state({
		folders: [] as Folder[],
		selected_folder: null as Folder | null,
		single_notes: [] as Note[],
		adding_new_folder: false,
		new_folder: {
			user_id: user?.id || null,
			title: '',
			is_favorite: false,
			parent_folder_id: null
		}
	});

	// Lifecycle
	import { onMount } from 'svelte';

	// Immediately organize notes into folders and single notes
	$effect(() => {
		if (folders && notes) {
			notes_list_state.folders = attachNotesToFolders(folders, notes);

			// Identify single notes (not in any folder)
			notes_list_state.single_notes = notes.filter((n: Note) => n.folder_id == null);
		}
	});

	// Functions
	const addFolder = async (parent_folder: Folder | null) => {
		try {
			const response = await fetch('/api/folders/add', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					user_id: user?.id,
					title: notes_list_state.new_folder.title,
					is_favorite: false,
					parent_folder_id: parent_folder ? parent_folder.id : null
				})
			});

			if (!response.ok) {
				throw new Error('Failed to add folder');
			}

			const newFolder = await response.json();
			console.log('Added folder:', newFolder);
			notes_list_state.folders = [...notes_list_state.folders, newFolder];
			// add folder to Folders Store:
			folders_store.update((folders) => [...folders, newFolder?.folder]);
		} catch (error) {
			console.error('Error adding folder:', error);
		}
	};

	const attachNotesToFolders = (rawFolders: Folder[], allNotes: Note[]): Folder[] => {
		const byFolder = new Map<number, Note[]>();
		for (const n of allNotes) {
			if (n.folder_id != null) {
				const arr = byFolder.get(n.folder_id) ?? [];
				arr.push(n);
				byFolder.set(n.folder_id, arr);
			}
		}

		// Recursive
		const clone = (f: Folder): Folder => ({
			...f,
			notes: byFolder.get(f.id) ?? [],
			subfolders: f.subfolders?.map(clone)
		});
		// Add subfolders to parents
		attachFoldersToFolderParents(rawFolders);

		// Identify folders without parents:
		const topLevelFolders = rawFolders.filter((f) => f.parent_folder_id?.Int64 == 0);
		return topLevelFolders.map(clone);
	};
	const attachFoldersToFolderParents = (rawFolders: Folder[]): Folder[] => {
		const byId = new Map<number, Folder>();
		for (const f of rawFolders) {
			byId.set(f?.id, f);
		}

		const roots: Folder[] = [];
		for (const f of rawFolders) {
			if (f.parent_folder_id != null) {
				const parent = byId?.get(f.parent_folder_id?.Int64);
				if (parent) {
					if (!parent?.subfolders) {
						parent.subfolders = [];
					}
					parent?.subfolders?.push(f);
				}
			} else {
				roots.push(f);
			}
		}
		return roots;
	};
</script>

<div
	id="notes-list"
	class="notes-list flex h-[90vh] w-1/5 flex-col overflow-scroll border-r border-slate-600 p-4"
>
	<div class="mb-4 flex w-full flex-row items-center justify-between">
		<h2 class="text-xl font-bold text-slate-200">Your Notes</h2>
		<button
			class="cursor-pointer rounded bg-slate-700 p-2 hover:bg-slate-600"
			onclick={() => (notes_list_state.adding_new_folder = !notes_list_state.adding_new_folder)}
		>
			<FolderPlusOutline class="h-6 w-6 shrink-0 text-slate-200" />
		</button>
	</div>
	{#if notes_list_state.adding_new_folder}
		<div class="mb-4 flex w-full flex-row items-center">
			<input
				name="title"
				type="text"
				bind:value={notes_list_state.new_folder.title}
				placeholder="New Folder Title"
				class="me-2 flex-1 rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
				onkeydown={(e) => {
					if (e.key === 'Enter') {
						addFolder(null);
						notes_list_state.adding_new_folder = false;
					}
				}}
			/>
			<button
				onclick={() => addFolder(null)}
				class="flex h-8 w-8 cursor-pointer items-center justify-center rounded bg-sky-500 p-2 text-white hover:bg-sky-600"
				disabled={notes_list_state.new_folder.title.trim() === ''}
			>
				+
			</button>
		</div>
	{/if}
	{#if folders && folders?.length > 0}
		<div class="folder mb-4 text-slate-200">
			<h3 class="mb-2 text-lg font-semibold text-slate-300">Folders</h3>
			<Accordion flush multiple>
				{#each notes_list_state.folders as folder}
					<NoteFolder {folder} {select_note} {onPostMove} />
				{/each}
			</Accordion>
		</div>
	{/if}
	{#if notes?.length > 0}
		<ul>
			{#each notes_list_state.single_notes as note}
				<NoteSelector {note} {select_note} />
			{/each}
		</ul>
	{:else}
		<p class="text-slate-400">Hmm. No notes. Get writing!</p>
	{/if}
</div>

<style>
	/* target svg child of .folder: */
	.folder :global(svg) {
		color: #94a3b8; /* Tailwind's text-slate-400 color */
	}
</style>
