<script lang="ts">
	// Types
	import type { Note } from '$lib/types/Note';
	// Props
	let { selected_note, deleted_note } = $props<{
		selected_note: Note;
		deleted_note: () => void; // This is a function sent up to notify parent on deletion
	}>();

	// imports
	//   Components
	import { Modal, DarkMode } from 'flowbite-svelte';
	import { Editor } from '@tadashi/svelte-editor-quill';

	//   Icons
	import { EditSolid, TrashBinSolid, FloppyDiskSolid } from 'flowbite-svelte-icons';

	//   Svelte
	import { fade } from 'svelte/transition';

	// State

	//.    Stores
	import { notes_store } from '../../stores/Notes';
	//.    Local state
	let note_viewer_state = $state({
		editing: false,
		editing_title: false,
		deleting: false,
		quill: {
			options: {
				placeholder: 'Compose an epic...',
				plainclipboard: true
			},
			content: {
				html: '',
				text: ''
			}
		}
	});

	// Auto format onMount
	import { onMount, tick } from 'svelte';
	onMount(() => {
		autoFormat();
	});

	// Prepopulate quill content when selected_note changes
	$effect(() => {
		if (selected_note) {
			autoFormat();
			note_viewer_state.quill.content.html = selected_note.content;
			// Optionally, you can extract plain text if needed
			const tempDiv = document.createElement('div');
			tempDiv.innerHTML = selected_note.content;
			note_viewer_state.quill.content.text = tempDiv.textContent || tempDiv.innerText || '';
		}
	});

	// Functions
	const onTextChange = (markup: any, plaintext: any) => {
		note_viewer_state.quill.content.html = markup;
		note_viewer_state.quill.content.text = plaintext;
	};

	const editTitle = (newTitle: string) => {
		note_viewer_state.editing_title = true;
		selected_note.title = newTitle;
		saveNote();
	};

	const saveNote = async () => {
		try {
			const response = await fetch(`/api/notes/update/${selected_note.id}`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					title: selected_note.title,
					content: note_viewer_state.quill.content.html
				})
			});

			if (!response.ok) {
				throw new Error(`Error updating note: ${response.statusText}`);
			}

			const data = await response.json();
			console.log('Note updated successfully!');
			// Optionally, update the selected_note content
			selected_note.content = note_viewer_state.quill.content.html;
			note_viewer_state.editing = false;

			// Update notes store:
			notes_store.update((notes) => {
				return notes.map((note) =>
					note.id === selected_note.id
						? { ...note, title: selected_note.title, content: selected_note.content }
						: note
				);
			});
			await tick();
			autoFormat();
		} catch (error) {
			console.error('Error saving note:', error);
		}
	};

	const autoFormat = () => {
		// find all <p> tags and add wrapping class:
		const noteBody = document.querySelector('.note-body.prose');
		if (noteBody) {
			const paragraphs = noteBody.querySelectorAll('p');
			paragraphs.forEach((p) => p.classList.add('wrap-anywhere'));
		}
	};

	const toggleEditing = async () => {
		autoFormat();
		note_viewer_state.editing = !note_viewer_state.editing;
		await tick();
		autoFormat();
		// if (note_viewer_state.editing) {
		//   // Prepopulate quill content when entering edit mode
		//   note_viewer_state.quill.content.html = selected_note.content;
		//   const tempDiv = document.createElement('div');
		//   tempDiv.innerHTML = selected_note.content;
		//   note_viewer_state.quill.content.text = tempDiv.textContent || tempDiv.innerText || '';
		// }
	};

	const deleteNote = async () => {
		try {
			const response = await fetch(`/api/notes/delete/${selected_note.id}`, {
				method: 'DELETE'
			});
			if (!response.ok) {
				throw new Error(`Error deleting note: ${response.statusText}`);
			}
			console.log('Note deleted successfully');
			// Update notes store:
			notes_store.update((notes) => notes.filter((note) => note.id !== selected_note.id));
			deleted_note();
		} catch (error) {
			console.error('Error deleting note:', error);
		}
	};
</script>

<svelte:head>
	<link
		rel="stylesheet"
		href="https://unpkg.com/quill@2.0.3/dist/quill.snow.css"
		crossorigin="anonymous"
	/>
</svelte:head>

<div
	id={`note-view-${selected_note.id}`}
	in:fade
	out:fade={{ duration: 400 }}
	class="note-viewer mt-6 flex w-4/5 flex-col border-l border-slate-700 px-4 pb-4 pt-6"
>
	<div
		id={`note-actions-${selected_note.id}`}
		class="mx-auto mb-2 mt-6 flex w-[98%] flex-row items-end justify-end p-2 pt-4"
	>
		{#if note_viewer_state.editing}
			<button
				class="btn btn-lg mr-4 cursor-pointer text-green-400 hover:text-green-500"
				onclick={saveNote}
			>
				<FloppyDiskSolid class="h-8 w-8 shrink-0" />
			</button>
		{/if}
		<button
			class="btn btn-lg mr-4 cursor-pointer text-slate-400 hover:text-slate-500"
			onclick={toggleEditing}
		>
			<EditSolid class="h-8 w-8 shrink-0" />
		</button>
		<button
			class="btn btn-lg cursor-pointer text-red-400 hover:text-red-500"
			onclick={() => (note_viewer_state.deleting = true)}
		>
			<TrashBinSolid class="h-8 w-8 shrink-0" />
		</button>
		<!-- Confirm delete modal: -->
		<Modal
			title={`Delete: ${selected_note.title}`}
			class="dark"
			bind:open={note_viewer_state.deleting}
			onaction={({ action }) => alert(`Handle "${action}"`)}
		>
			<div>
				<p>Are you sure you want to delete this note?</p>
				<div class="mt-4 flex flex-row justify-end">
					<button
						class="btn btn-secondary mr-4 cursor-pointer"
						onclick={() => (note_viewer_state.deleting = false)}
					>
						Cancel
					</button>
					<button
						class="btn btn-danger cursor-pointer text-red-500"
						onclick={() => {
							deleteNote();
							note_viewer_state.deleting = false;
						}}
					>
						Delete
					</button>
				</div>
			</div>
		</Modal>
	</div>
	{#if note_viewer_state.editing}
		<div
			class="h-full w-full overflow-hidden rounded-sm bg-slate-100"
			in:fade={{ duration: 400 }}
			out:fade={{ duration: 400 }}
		>
			<Editor class="h-96" options={note_viewer_state.quill.options} {onTextChange}
				>{@html note_viewer_state.quill.content.html}</Editor
			>
		</div>
	{:else}
		<div
			in:fade={{ duration: 400 }}
			out:fade={{ duration: 400 }}
			class="flex-column flex w-full flex-col"
		>
			<div class="flex w-full flex-row items-start">
				{#if note_viewer_state.editing_title}
					<input
						type="text"
						bind:value={selected_note.title}
						onblur={() => (
							editTitle(selected_note.title),
							(note_viewer_state.editing_title = false)
						)}
						onkeydown={(e) => {
							if (e.key === 'Enter') {
								editTitle(selected_note.title);
							}
						}}
						class="rounded-sm border-b border-slate-500 bg-transparent text-xl font-bold text-slate-200 focus:outline-none"
					/>
				{:else}
					<h2 class="text-2xl font-bold text-slate-100">
						{selected_note.title}
					</h2>
				{/if}
				<button
					class="btn btn-sm cursor-pointer text-slate-400 hover:text-slate-500"
					onclick={() => (note_viewer_state.editing_title = true)}
				>
					<EditSolid class="h-5 w-5 shrink-0" />
				</button>
			</div>
			<!-- Wrap text:  -->
			<div class="note-body prose mt-2 flex w-4/5 flex-col text-start text-slate-300">
				{@html selected_note.content}
			</div>
		</div>
	{/if}
</div>

<style>
	.prose {
		max-width: 150ch !important;
	}
</style>
