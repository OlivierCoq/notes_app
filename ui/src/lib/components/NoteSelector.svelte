<script lang="ts">
	// Types
	import type { Note } from '$lib/types/Note';

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
</script>

<li
	in:scale
	out:fade={{ duration: 400 }}
	class="note-selector mb-2 cursor-pointer rounded-md p-2 hover:bg-slate-600 dark:hover:bg-slate-500"
>
	<button onclick={handle_click} class="cursor-pointer">
		<div class="flex flex-col">
			<div class="flex max-w-[300px] flex-col p-2 text-start">
				<h3 class="text-lg font-semibold text-slate-100">{note?.title}</h3>
				<!-- HTML markup, preview: -->
				<div class="prose line-clamp text-sm text-slate-300">{@html note?.content}</div>
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
