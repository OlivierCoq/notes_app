<script lang="ts">
	// types
	import type { Folder } from '$lib/types/Folder';
	// Props
	let { folder, select_note } = $props();

	// Components
	import NoteSelector from '$lib/components/NoteSelector.svelte';
	import NoteFolder from '$lib/components/NoteFolder.svelte';

	// console.log('Folder in NoteFolder:', folder);

	// imports
	import { AccordionItem, Accordion } from 'flowbite-svelte';
	// icons:
	import { FolderSolid } from 'flowbite-svelte-icons';
	//   Svelte
	import { fade } from 'svelte/transition';
</script>

<AccordionItem
	id={`folder-${folder.id}`}
	class="align-start mb-2 flex flex-row rounded bg-slate-700 text-slate-200"
	classes={{ inactive: 'text-slate-200' }}
>
	{#snippet header()}
		<div class="flex flex-row items-center gap-2">
			<FolderSolid class="h-6 w-6 shrink-0" />
			<p class="text-slate-200">{folder.title}</p>
		</div>
	{/snippet}
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
					<NoteFolder folder={subfolder} {select_note} />
				{/each}
			</div>
		{/if}
	</div>
</AccordionItem>
