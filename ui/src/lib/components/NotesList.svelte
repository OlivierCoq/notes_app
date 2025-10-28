<script lang="ts">
	// Props
	let { notes, select_note } = $props();

	// Imports
	//   Svelte
	//   Components
	import NoteSelector from '$lib/components/NoteSelector.svelte';

	// Types
	import type { Note } from '$lib/types/Note';
	interface Folder {
		id: number;
		name: string;
		notes: Note[];
		subfolders?: Folder[];
	}

	// State
	let notes_list_state = $state({
		folders: [] as Folder[],
		selected_folder: null as Folder | null,
		single_notes: [] as Note[]
	});

	// Lifecycle
	import { onMount } from 'svelte';

	// Functions
</script>

<div class="notes-list flex h-[90vh] w-1/5 flex-col overflow-scroll border-r border-slate-600 p-4">
	<h2 class="mb-4 text-xl font-bold text-slate-200">Your Notes</h2>
	{#if notes.length > 0}
		<ul>
			{#each notes as note}
				<NoteSelector {note} {select_note} />
			{/each}
		</ul>
	{:else}
		<p class="text-slate-400">Hmm. No notes. Get writing!</p>
	{/if}
</div>
