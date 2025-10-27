<script lang="ts">
	// Components
	import Modal from '$lib/components/ui/Modal.svelte';

	// Data and State
	let dashboard_state = $state({});

	// Stores

	// Methods
</script>

<div class="flex min-h-[95vh] w-full flex-col">
	<div class="flex w-full flex-1"></div>
	<div class="flex w-full flex-row justify-end p-4">
		<button
			class="cursor-pointer rounded-full bg-sky-400 px-4 py-2 text-white shadow-xl hover:bg-sky-500"
			><span class="text-4xl">+</span></button
		>
	</div>
</div>
