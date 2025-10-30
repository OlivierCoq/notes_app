<script lang="ts">
	// types
	import type { Folder } from '$lib/types/Folder';
	// Props
	let { folder, select_note, user } = $props();

	// Stores:
	import { notes_store } from '../../stores/Notes';
	import { folders_store } from '../../stores/Folders';

	// Components
	import NoteSelector from '$lib/components/NoteSelector.svelte';
	import NoteFolder from '$lib/components/NoteFolder.svelte';

	// console.log('Folder in NoteFolder:', folder);

	// imports
	import { AccordionItem, Accordion } from 'flowbite-svelte';
	// icons:
	import { FolderPlusOutline, FolderSolid, DotsVerticalOutline } from 'flowbite-svelte-icons';
	//   Svelte
	import { fade } from 'svelte/transition';
	import { tick } from 'svelte';

	// State
	let note_folder_state = $state({
		adding_subfolder: false,
		new_subfolder: {
			user_id: user?.id || null,
			title: '',
			is_favorite: false,
			parent_folder_id: {
				Int64: folder.id,
				Valid: true
			}
		}
	});

	// Functions
	const addSubfolder = async () => {
		// Implementation for adding a subfolder
		try {
			const response = await fetch('/api/folders/add', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(note_folder_state.new_subfolder)
			});
			if (!response.ok) {
				throw new Error(`Error adding subfolder: ${response.statusText}`);
			}
			const newFolder = await response.json();
			// update folders_store:

			console.log('New subfolder added:', newFolder);
			folders_store.update((folders) => [...folders, newFolder]);
			await tick();
			folder.subfolders.push(newFolder?.folder);
		} catch (error) {
			console.error('Failed to add subfolder:', error);
		}
	};

	//     Drag + Drop:
	let isOver = false;
	const allowDrop = (event: DragEvent) => {
		event.preventDefault();
		console.log('Drag over folder:', folder.title);
		if (event.dataTransfer?.types?.includes('application/x-note-id')) {
			event.dataTransfer.dropEffect = 'move';
		}
	};
	const setHover = (hover: boolean) => {
		isOver = hover;
	};
	const handleDrop = async (event: DragEvent) => {
		event.preventDefault();
		setHover(false);
		const data_transfer = event.dataTransfer;
		// console.log('Data transfer on drop:', data_transfer);
		if (!data_transfer) return;
		const id_str = data_transfer.getData('application/x-note-id');
		// console.log('Data transfer ID string:', id_str);
		const note_id = Number(id_str);
		if (!note_id) return;
		// console.log('Dropped note ID:', note_id, 'on folder ID:', folder.id);
		// Update the note's folder_id via API call
		try {
			const response = await fetch(`/api/notes/update/${note_id}`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ folder_id: folder.id })
			});
			if (!response.ok) {
				throw new Error(`Error updating note: ${response.statusText}`);
			}
			console.log('Note moved successfully');
			// Update notes_store to reflect the change
			notes_store.update((notes) => {
				return notes.map((note) => {
					if (note.id === note_id) {
						return { ...note, folder_id: folder.id };
					}
					return note;
				});
			});
		} catch (error) {
			console.error('Failed to move note:', error);
		}
	};
	/// Dragging Folders
	const handle_dragstart = (event: DragEvent) => {
		if (!event.dataTransfer) return;
		event.dataTransfer?.setData('application/x-folder-id', folder.id.toString());
		event.dataTransfer.effectAllowed = 'move';
		event.dataTransfer?.setData('text/plain', JSON.stringify(folder));
		// console.log('Dragging folder:', folder);
	};
	// Receiving folders
	const handle_folder_drop = async (event: DragEvent) => {
		event.preventDefault();
		setHover(false);
		const data_transfer = event.dataTransfer;
		if (!data_transfer) return;
		const id_str = data_transfer.getData('application/x-folder-id');
		const folder_id = Number(id_str);
		if (!folder_id) return;
		// console.log('Dropped folder ID:', folder_id, 'on folder ID:', folder.id);
		// Update the folder's parent_folder_id via API call
		try {
			const response = await fetch(`/api/folders/update/${folder_id}`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ parent_folder_id: folder.id })
			});
			if (!response.ok) {
				throw new Error(`Error updating folder: ${response.statusText}`);
			}
			console.log('Folder moved successfully');
			// Update folders_store to reflect the change
			folders_store.update((folders) => {
				return folders.map((f) => {
					if (f.id === folder_id) {
						return { ...f, parent_folder_id: folder.id };
					}
					return f;
				});
			});
		} catch (error) {
			console.error('Failed to move folder:', error);
		}
	};
</script>

<AccordionItem
	id={`folder-${folder.id}`}
	class="align-start mb-2 flex flex-row rounded bg-slate-700 text-slate-200"
	classes={{ inactive: 'text-slate-200' }}
>
	{#snippet header()}
		<div
			class="flex flex-row items-center gap-2"
			ondragenter={() => setHover(true)}
			ondragleave={() => setHover(false)}
			ondragover={allowDrop}
			ondrop={handleDrop}
			role="group"
			aria-label={`Folder ${folder.title}`}
			draggable="true"
			ondragstart={handle_dragstart}
			class:bg-slate-600={isOver}
		>
			<FolderSolid class="h-6 w-6 shrink-0" />
			<p class="text-slate-200">{folder.title}</p>
			<div class="flex-1"></div>
		</div>
	{/snippet}
	<div class="flex flex-col">
		<div class="flex flex-row items-center gap-2">
			<button
				class="cursor-pointer p-1 hover:visible hover:bg-slate-600 group-hover:visible"
				onclick={() => (note_folder_state.adding_subfolder = !note_folder_state.adding_subfolder)}
				aria-label="Add Subfolder"
			>
				<FolderPlusOutline class="h-6 w-6 shrink-0" />
			</button>
		</div>
		{#if note_folder_state.adding_subfolder}
			<div class="flex flex-row">
				<input
					in:fade
					out:fade={{ duration: 400 }}
					type="text"
					placeholder="New Folder"
					bind:value={note_folder_state.new_subfolder.title}
					class="mb-2 ms-4 mt-1 w-3/4 rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
				/>
				<button
					class="mb-2 ms-2 mt-1 h-10 w-10 cursor-pointer rounded bg-sky-400 p-2 text-white hover:bg-sky-500"
					onclick={() => {
						addSubfolder();
						note_folder_state.adding_subfolder = false;
					}}
				>
					+
				</button>
			</div>
		{/if}
	</div>
	<div>
		<ul class="mt-2">
			{#each folder.notes as note}
				<NoteSelector {note} {select_note} />
			{/each}
		</ul>
		{#if folder.subfolders && folder.subfolders.length > 0}
			<div class="ml-4 mt-4 border-l border-slate-600 pl-4">
				{#each folder.subfolders as subfolder}
					<!-- Insane self-referencing loop!!! :)  -->
					<NoteFolder folder={subfolder} {select_note} {user} />
				{/each}
			</div>
		{/if}
	</div>
</AccordionItem>
