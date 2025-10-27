<script lang="ts">
	import { onMount } from 'svelte';
	import { user as userStore } from '$lib/stores/AppUser';
	import type { User } from '$lib/types/User';
	let user: User | null = null;

	userStore.subscribe((value) => {
		user = value;
	});
</script>

<div>
	{#if user}
		<h1 class="mb-4 text-2xl font-bold">Dashboard</h1>
		<p>Welcome, {user.username}!</p>
	{:else}
		<p>Loading user data...</p>
	{/if}
</div>
