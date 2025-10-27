<script lang="ts">
	import '../app.css';
	// import favicon from '$lib/assets/favicon.svg';

	// Data and State

	// Stores
	import { user as userStore, setUser } from '$lib/stores/AppUser';
	// Access props using $props() in runes mode
	let { data, children } = $props(); // Uncomment and adjust if you need to use 'data' prop

	// console.log('Layout data:', data);
	// console.log('data.user:', data?.user);

	// Populate userStore with data.user
	$effect(() => {
		console.log('Updating userStore with:', data?.user);
		setUser(data?.user ?? null);
	});
	// Components
	import Navbar from '$lib/components/Navbar.svelte';
</script>

<svelte:head>
	<!-- <link rel="icon" href={favicon} /> -->
</svelte:head>

<main class="align-center h-[100vh] w-full bg-slate-800">
	<Navbar />
	{@render children?.()}
</main>
