<script lang="ts">
	// Types
	import type { Note } from '$lib/types/Note';
	import { tick } from 'svelte';

	// Props
	let { note, select_note } = $props<{
		note: Note;
		select_note: (note: Note) => void;
	}>();

	// imports
	//   Svelte
	import { fade, scale } from 'svelte/transition';

	// functions
	const handle_click = () => {
		if (select_note) {
			select_note(note);
		}
	};
	// Dragging:
	const handle_dragstart = async (event: DragEvent) => {
		if (!event.dataTransfer) return;
		event.dataTransfer?.setData('text/plain', JSON.stringify(note));
		event.dataTransfer?.setData('application/note-id', note.id.toString());
		// push into dataTransfer items:
		event.dataTransfer?.setData('application/x-note-id', note.id.toString());
		event.dataTransfer.effectAllowed = 'move';
		// await tick();
		// console.log(event.dataTransfer);
	};
</script>

<li
	in:scale
	out:fade={{ duration: 400 }}
	class="note-selector mb-2 max-h-[140px] cursor-pointer overflow-hidden rounded-md p-2 hover:bg-slate-600 dark:hover:bg-slate-500"
>
	<button
		onclick={handle_click}
		draggable="true"
		ondragstart={handle_dragstart}
		class="cursor-pointer overflow-hidden"
	>
		<div class="flex flex-col overflow-hidden rounded-md">
			<div class="flex max-w-[300px] flex-col overflow-hidden rounded-md p-2 text-start">
				<h3 class="text-lg font-semibold text-slate-100">{note?.title}</h3>
				<!-- HTML markup, preview: -->
				<div class="prose line-clamp overflow-hidden pe-6 text-sm text-slate-300">
					{@html note?.content}
				</div>
			</div>
		</div>
	</button>
</li>

<style>
	.line-clamp {
		display: -webkit-box;
		-webkit-line-clamp: 3; /* number of lines you want */
		line-clamp: 3; /* standard property for compatibility */
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
