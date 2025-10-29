<script lang="ts">
	// Types
	import type { Note } from '$lib/types/Note';
	// Props
	let { selected_note } = $props<{
		selected_note: Note;
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
	let note_viewer_state = $state({
		editing: false,
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

	// Lifecycle
	import { onMount } from 'svelte';
	// Prepopulate quill content when selected_note changes
	$effect(() => {
		if (selected_note) {
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

	const saveNote = async () => {
		try {
			const response = await fetch(`/api/notes/${selected_note.id}/update`, {
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
			console.log('Note updated successfully:', data);
			// Optionally, update the selected_note content
			selected_note.content = note_viewer_state.quill.content.html;
			note_viewer_state.editing = false;
		} catch (error) {
			console.error('Error saving note:', error);
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
	class="note-viewer flex w-4/5 flex-col border-l border-slate-700 p-4"
>
	<div
		id={`note-actions-${selected_note.id}`}
		class="mx-auto my-2 flex w-[98%] flex-row items-end justify-end p-2"
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
			onclick={() => (note_viewer_state.editing = !note_viewer_state.editing)}
		>
			<EditSolid class="h-8 w-8 shrink-0" />
		</button>
		<button class="btn btn-lg cursor-pointer text-red-400 hover:text-red-500">
			<TrashBinSolid class="h-8 w-8 shrink-0" />
		</button>
	</div>
	{#if note_viewer_state.editing}
		<div
			class="h-full w-full overflow-hidden rounded-sm bg-slate-100"
			in:fade
			out:fade={{ duration: 400 }}
		>
			<Editor class="h-96" options={note_viewer_state.quill.options} {onTextChange}
				>{@html note_viewer_state.quill.content.html}</Editor
			>
		</div>
	{:else}
		<div class="note-body">
			<h2 class="text-xl font-bold text-slate-200">{selected_note.title}</h2>
			<div
				class="note-body prose prose-slate mt-2 flex max-w-none flex-wrap text-wrap text-slate-300"
			>
				{@html selected_note.content}
			</div>
		</div>
	{/if}
</div>

<style>
</style>
