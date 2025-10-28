<script lang="ts">
	/* Note to self: The filepath nomenclature is as follows:
   - user_id-note_id-~#~
    - Example: 3-15-~#~
*/

	// Props
	let { note, select_note } = $props();

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
	<button onclick={handle_click}>
		<div class="flex flex-col">
			<div class="flex flex-col">
				<h3 class="text-lg font-semibold text-slate-100">{note?.title}</h3>
				<!-- HTML markup: -->
				<p class="text-sm text-slate-300">{@html note?.content}</p>
			</div>
		</div>
	</button>
</li>
