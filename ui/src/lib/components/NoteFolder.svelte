<script lang="ts">
	// types
	import type { Folder } from '$lib/types/Folder';
	// Props
	let { folder, select_note, onPostMove } = $props();

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

	// Drag + Drop:
	let isOver = false;
	const allowDrop = (event: DragEvent) => {
		event.preventDefault();
		console.log('Drag over folder:', folder.title);
		if (event.dataTransfer?.types?.includes('application/x-note-id')) {
			event.dataTransfer.dropEffect = 'move';
		}
	};
	const setHover = (hover: boolean) => {
		isOver = hover;
	};
	const handleDrop = async (event: DragEvent) => {
		event.preventDefault();
		setHover(false);
		const data_transfer = event.dataTransfer;
		console.log('Data transfer on drop:', data_transfer);
		if (!data_transfer) return;
		const id_str = data_transfer.getData('application/x-note-id');
		console.log('Data transfer ID string:', id_str);
		const note_id = Number(id_str);
		if (!note_id) return;
		console.log('Dropped note ID:', note_id, 'on folder ID:', folder.id);
		// Update the note's folder_id via API call
		try {
			const response = await fetch(`/api/notes/update/${note_id}`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ folder_id: folder.id })
			});
			if (!response.ok) {
				throw new Error(`Error updating note: ${response.statusText}`);
			}
			console.log('Note moved successfully');
			if (onPostMove) {
				await onPostMove();
			}
		} catch (error) {
			console.error('Failed to move note:', error);
		}
	};
</script>

<AccordionItem
	id={`folder-${folder.id}`}
	class="align-start mb-2 flex flex-row rounded bg-slate-700 text-slate-200"
	classes={{ inactive: 'text-slate-200' }}
>
	{#snippet header()}
		<div
			class="flex flex-row items-center gap-2"
			ondragenter={() => setHover(true)}
			ondragleave={() => setHover(false)}
			ondragover={allowDrop}
			ondrop={handleDrop}
			role="group"
			aria-label={`Folder ${folder.title}`}
		>
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
					<NoteFolder folder={subfolder} {select_note} {onPostMove} />
				{/each}
			</div>
		{/if}
	</div>
</AccordionItem>
