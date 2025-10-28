<script lang="ts">
	// Imports
	// environment
	import { PUBLIC_API_URL } from '$env/static/public';
	// Svelte
	import { fade, scale } from 'svelte/transition';

	// Note: User data now comes from server load function instead of client store

	// Components
	import { Modal, DarkMode } from 'flowbite-svelte';
	import { Editor } from '@tadashi/svelte-editor-quill';

	// Types
	import type { Note } from '$lib/types/Note';

	// Data and State
	// props from server load function:
	let { data } = $props();

	// Extract notes and user from the server data
	/* 
  
    Using $derived to create reactive variables for notes, user, and error
    based on the data provided by the server load function.
  */
	let notes = $derived(data.notes || []);
	let user = $derived(data.user);
	let error = $derived(data.error);

	// Update dashboard state when notes change
	$effect(() => {
		// Waits for notes to be available
		if (notes) {
			dashboard_state.all_notes = notes;
		}
	});

	let dashboard_state = $state({
		all_notes: [] as Note[],
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
		}
	});

	const onTextChange = (markup: any, plaintext: any) => {
		dashboard_state.quill.content.html = markup;
		dashboard_state.quill.content.text = plaintext;
	};

	const loadNotes = async () => {
		// Fetch notes logic here using server-provided user data
		if (!user?.id) {
			console.error('No user ID available');
			return;
		}

		const response = await fetch(`/api/notes/all/${user.id}`);
		if (response.ok) {
			const data = await response.json();
			dashboard_state.all_notes = data.notes;
			// console.log('Loaded notes:', dashboard_state.all_notes);
		} else {
			console.error('Error loading notes:', response.statusText);
		}
	};

	// Stores

	// Methods
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
			is_favorite: false
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
			// console.log('Note added successfully:', data);
			// Optionally, refresh the notes list or provide user feedback here
			await loadNotes();
		} else {
			console.error('Error adding note:', response.statusText);
		}

		// Close modal and reset fields
		dashboard_state.new_note_modal = false;
		dashboard_state.new_note.title = '';
		dashboard_state.quill.content.html = '';
		dashboard_state.quill.content.text = '';
	};
</script>

<svelte:head>
	<link
		rel="stylesheet"
		href="https://unpkg.com/quill@2.0.3/dist/quill.snow.css"
		crossorigin="anonymous"
	/>
</svelte:head>

<div class="flex min-h-[95vh] w-full flex-col">
	<!-- Error Display -->
	{#if error}
		<div class="mb-4 rounded border border-red-400 bg-red-100 px-4 py-3 text-red-700">
			<p><strong>Error:</strong> {error}</p>
		</div>
	{/if}

	<div class="flex w-full flex-1">
		<!-- Notes List -->
		<div class="flex h-[90vh] w-1/5 flex-col overflow-scroll border-r border-slate-600 p-4">
			<h2 class="mb-4 text-xl font-bold text-slate-200">Your Notes</h2>
			<!-- If dashboard_state.all_notes.length > 0 -->
			{#if dashboard_state.all_notes.length > 0}
				<ul>
					{#each dashboard_state.all_notes as note}
						<li
							in:scale
							out:fade={{ duration: 400 }}
							class="mb-2 cursor-pointer rounded-md p-2 hover:bg-slate-600 dark:hover:bg-slate-500"
						>
							<h3 class="text-lg font-semibold text-slate-100">{note?.title}</h3>
							<!-- <p class="text-sm text-slate-300">{note.content.replace(/<[^>]+>/g, '').slice(0, 50)}{note.content.length > 50 ? '...' : ''}</p> -->
							<!-- HTML markup: -->
							<p class="text-sm text-slate-300">{@html note?.content}</p>
						</li>
					{/each}
				</ul>
			{:else}
				<p class="text-slate-400">Hmm. No notes. Get writing!</p>
			{/if}
		</div>
	</div>
	<div class="flex w-full flex-row justify-end p-4">
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
			<!-- {#snippet footer()}
        <button type="submit" value="success">I accept</button>
        <button type="submit" value="decline" color="alternative">Decline</button>
      {/snippet} -->
		</Modal>
	</div>
</div>

<!-- Modals -->
