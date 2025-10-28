<script lang="ts">

  import { DarkMode } from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import { user as userStore } from '$lib/stores/AppUser';
	import type { User } from '$lib/types/User';
	let user: User | null = null;

	userStore.subscribe((value) => {
		user = value;
	});

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

<div class="flex w-full flex-row justify-between bg-slate-800 p-3 shadow-md">
	<div class="flex w-1/2 flex-row justify-start text-start">
		<a href="/dashboard">
			<h1 class="text-xl font-bold text-slate-200">notez</h1>
		</a>
	</div>
	<div class="align-right flex w-1/2 flex-row items-end justify-end gap-2 space-y-1 text-end">
		{#if user}
			<div class="flex flex-col items-end">
				<p class="m-0 p-0 text-slate-200">Logged in: <strong>{user.username}</strong></p>
				<button class="cursor-pointer text-xs text-slate-200" onclick={logout}>Logout</button>
			</div>
			<a href="/account">
				<img src={user?.pfp_url} alt="User Profile" class="h-10 w-10 rounded-full" /></a
			>
		{/if}
	</div>
</div>
