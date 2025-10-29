<script lang="ts">
	// Types
	import type { Note } from '$lib/types/Note';
	// Props
	let { selected_note } = $props<{
		selected_note: Note;
	}>();

	// imports
	//   Svelte
	import { fade } from 'svelte/transition';
</script>

<div
	id={`note-view-${selected_note.id}`}
	in:fade
	out:fade={{ duration: 400 }}
	class="note-viewer flex w-4/5 flex-col border-l border-slate-700 p-4"
>
	<h2 class="text-xl font-bold text-slate-200">{selected_note.title}</h2>
	<div class="note-body prose prose-slate mt-2 flex max-w-none flex-wrap text-wrap text-slate-300">
		{@html selected_note.content}
	</div>
</div>

<style>
</style>
