<script lang="ts">
	import '../app.css';
	// import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<!-- <link rel="icon" href={favicon} /> -->
</svelte:head>

<main class="align-center flex h-[100vh] w-full flex-col justify-center bg-slate-800">
	{@render children?.()}
</main>
