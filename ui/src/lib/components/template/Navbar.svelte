<script lang="ts">
	// Components
	import { DarkMode, Tooltip } from 'flowbite-svelte';
	import { onMount } from 'svelte';

	// Stores
	import { user as userStore } from '$lib/stores/AppUser';

	// Types
	import type { User } from '$lib/types/User';

	// Data and State
	let user: User | null = null;

	userStore.subscribe((value) => {
		user = value;
	});

	// Functions
	const logout = () => {
		// Implement logout functionality here
		// POST request to logout endpoint:
		fetch('/logout', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			}
		})
			.then((response) => {
				if (response.ok) {
					// Handle successful logout
					console.log('Logout successful');
					userStore.set(null);
					// Redirect to home:
					window.location.href = '/';
				} else {
					// Handle logout error
					console.error('Logout failed');
				}
			})
			.catch((error) => {
				console.error('Error:', error);
			});
		console.log('Logout clicked');
	};
</script>

<div id="navbar" class="fixed flex w-full flex-row justify-between bg-slate-800 p-3 shadow-md">
	<div class="flex w-1/2 flex-row justify-start text-start">
		<a id="dashboard-link" href="/dashboard">
			<h1 class="text-xl font-bold text-slate-200">n o t e z</h1>
		</a>
		<Tooltip triggeredBy="#dashboard-link">Your notes</Tooltip>
	</div>
	<div class="align-right flex w-1/2 flex-row items-end justify-end gap-2 space-y-1 text-end">
		{#if user}
			<div class="flex flex-col items-end">
				<p class="m-0 p-0 text-slate-200">Logged in: <strong>{user.username}</strong></p>
				<button class="cursor-pointer text-xs text-slate-200" onclick={logout}>Logout</button>
			</div>
			<a id="account-link" href="/account">
				<img src={user?.pfp_url} alt="User Profile" class="h-10 w-10 rounded-full" /></a
			>
			<Tooltip triggeredBy="#account-link">Edit account</Tooltip>
		{/if}
	</div>
</div>
