<script lang="ts">
	// Note: User data now comes from server load function instead of client store

	// Components
	import { Modal, DarkMode } from 'flowbite-svelte';
	import { Editor } from '@tadashi/svelte-editor-quill';
	import NotesList from '$lib/components/NotesList.svelte';
	import NoteViewer from '$lib/components/NoteViewer.svelte';

	// Types
	import type { Note } from '$lib/types/Note';
	import type { User } from '$lib/types/User';

	// Stores
	import { folders_store, refreshFolders } from '../../stores/Folders';
	import { notes_store, refreshNotes } from '../../stores/Notes';
	import { onMount } from 'svelte';

	// Data and State
	//     props from server load function:
	let { data } = $props();

	// Extract notes and user from the server data
	/* 
    Using $derived to create reactive variables for notes, user, and error
    based on the data provided by the server load function.
  */
	let notes = $derived(data.notes || []);
	let folders = $derived(data.folders || []);
	let user = $derived(data.user);
	let error = $derived(data.error);

	// Set folders in the store
	$effect(() => {
		folders_store.set(folders);
		notes_store.set(notes);
		// console.log
		// refreshData();
	});

	// onMount(() => {
	// 	refreshData();
	// });

	let dashboard_state = $state({
		new_note: {
			title: ''
		},
		new_note_modal: false,
		quill: {
			options: {
				placeholder: 'Compose an epic...',
				plainclipboard: true
			},
			content: {
				html: '',
				text: ''
			}
		},
		selected_id: null as number | null
	});

	let selected_note = $derived(
		() => notes.find((note: Note) => note.id === dashboard_state.selected_id) || null
	);

	const onTextChange = (markup: any, plaintext: any) => {
		dashboard_state.quill.content.html = markup;
		dashboard_state.quill.content.text = plaintext;
	};

	// Methods
	const refreshData = async () => {
		let user = data.user;
		if (user && user.id) {
			await Promise.all([refreshNotes(user?.id), refreshFolders(user?.id)]);
		}
	};

	const submitNewNote = async () => {
		// Handle new note submission logic here using server-provided user data
		if (!user?.id) {
			console.error('No user ID available');
			return;
		}

		let post_obj = {
			user_id: user.id,
			title: dashboard_state.new_note.title,
			content: dashboard_state.quill.content.html,
			is_favorite: false,
			folder_id: null
		};
		// console.log('submitting: ', post_obj);

		const response = await fetch(`/api/notes/add/`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(post_obj)
		});

		// console.log('Response:', response);

		if (response.ok) {
			const data = await response.json();
			console.log('Note added successfully:', data);
			// update notes store:
			notes_store.update((notes) => [...notes, data.note]);
			// Update local notes array
			notes = [...notes, data.note];
			// Auto select the newly created note
			dashboard_state.selected_id = data.note.id;
		} else {
			console.error('Error adding note:', response.statusText);
		}

		// Close modal and reset fields
		dashboard_state.new_note_modal = false;
		dashboard_state.new_note.title = '';
		dashboard_state.quill.content.html = '';
		dashboard_state.quill.content.text = '';
	};

	const select_note = (note: Note) => {
		dashboard_state.selected_id = note.id;
	};

	const deleted_note = () => {
		// console.log('Note deleted callback received');
		dashboard_state.selected_id = null;
	};
</script>

<svelte:head>
	<link
		rel="stylesheet"
		href="https://unpkg.com/quill@2.0.3/dist/quill.snow.css"
		crossorigin="anonymous"
	/>
</svelte:head>

<div id="dashboard" class="flex min-h-[95vh] w-full flex-col">
	<!-- Error Display -->
	{#if error}
		<div class="mb-4 rounded border border-red-400 bg-red-100 px-4 py-3 text-red-700">
			<p><strong>Error:</strong> {error}</p>
		</div>
	{/if}

	<div id="notes-interface" class="flex w-full flex-1 flex-row">
		<!-- Notes List -->
		{#if user && notes}
			<!-- Populate notesList from store: -->
			<NotesList
				notes={$notes_store}
				folders={$folders_store}
				{select_note}
				{user}
				onPostMove={refreshData}
			/>
			{#if selected_note()}
				<NoteViewer selected_note={selected_note()} {deleted_note} />
			{:else}
				<div
					class="note-viewer flex w-full flex-1 flex-col border-l border-slate-700 p-4 text-slate-400"
				>
					<p>Select a note to view its content.</p>
				</div>
			{/if}
		{/if}
	</div>
	<div id="user-actions-interface" class="flex w-full flex-row justify-end p-4">
		<button
			class="cursor-pointer rounded-full bg-sky-400 px-4 py-2 text-white shadow-xl hover:bg-sky-500"
			onclick={() => (dashboard_state.new_note_modal = true)}
			><span class="text-4xl">+</span></button
		>
		<Modal
			title="New note"
			class="dark"
			bind:open={dashboard_state.new_note_modal}
			onaction={({ action }) => alert(`Handle "${action}"`)}
		>
			<div>
				<input
					type="text"
					name="title"
					placeholder="Note Title"
					class="mb-2 w-full rounded-md border border-slate-100 p-2 focus:border-slate-500 focus:outline-none dark:border-slate-600
          dark:bg-slate-800 dark:text-slate-200 dark:focus:border-slate-400"
					bind:value={dashboard_state.new_note.title}
				/>
			</div>
			<div class="h-full w-full overflow-hidden rounded-sm bg-slate-100">
				<Editor class="h-96" options={dashboard_state.quill.options} {onTextChange}
					>{@html dashboard_state.quill.content.html}</Editor
				>
			</div>
			<div class="w-full p-2">
				<button
					class="btn btn-block w-full cursor-pointer rounded-md bg-sky-400 px-4 py-2 text-white shadow-md hover:bg-sky-500"
					disabled={dashboard_state.new_note.title === '' ||
						dashboard_state.quill.content.text === ''}
					onclick={submitNewNote}>Submit</button
				>
			</div>
		</Modal>
	</div>
</div>
