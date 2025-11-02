<script lang="ts">
	import '../app.css';
	import '../app.scss';
	// import favicon from '$lib/assets/favicon.svg';

	// Data and State

	// Stores
	import { user as userStore, setUser } from '$lib/stores/AppUser';
	// Access props using $props() in runes mode
	let { data, children } = $props();
	// console.log('data.user:', data?.user);

	// Populate userStore with data.user
	$effect(() => {
		// console.log('Updating userStore with:', data?.user);
		setUser(data?.user ?? null);
	});
	// Components
	import Navbar from '$lib/components/template/Navbar.svelte';
</script>

<svelte:head>
	<!-- <link rel="icon" href={favicon} /> -->
</svelte:head>

<main class="align-center min-h-[100vh] w-full bg-slate-700">
	{#if $userStore}
		<Navbar />
	{/if}
	{@render children?.()}
</main>
