<script lang="ts">
	// imports
	import { PUBLIC_API_URL } from '$env/static/public';
	import { page } from '$app/stores';
	import { enhance } from '$app/forms';

	// Svelte
	import { Alert } from 'flowbite-svelte';
	import { fade } from 'svelte/transition';

	// form data from server
	let { form } = $props<{ form?: any }>();

	// error handlign and stuff
	$effect(() => {
		if (form?.username) {
			login_state.user_creds.username = form.username;
		}
	});

	// State
	const login_state = $state({
		user_creds: {
			username: '',
			password: '',
			email: '',
			pfp_url:
				'https://res.cloudinary.com/dxsjva9e0/image/upload/v1761835316/user_avatar_ry4fdr.png'
		},
		logged_in: false,
		signup_mode: false,
		posting: false,
		success: '',
		error: ''
	});

	// Functions
	const signUp = async () => {
		login_state.posting = true;
		try {
			const response = await fetch(`/api/users/register`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(login_state.user_creds)
			});
			if (response.ok) {
				login_state.logged_in = true;
				// Success message:
				login_state.success = 'Signup successful! Redirecting to dashboard...';
				window.location.href = '/dashboard';
			} else {
				login_state.error = 'Signup failed. Please try again.';
				console.error('Signup failed:', response.statusText);
			}
		} catch (error) {
			login_state.error = `Signup failed: ${error} `;
			console.error('Error during signup:', error);
		} finally {
			login_state.posting = false;
		}
	};
</script>

<div class="flex h-[100vh] w-full flex-col justify-center">
	<div
		class="align-center mx-auto flex w-96 flex-col justify-center rounded-md bg-slate-100 p-8 shadow-md dark:bg-slate-700"
	>
		<h1 class="mb-6 text-center text-2xl font-bold text-slate-800 dark:text-slate-200">notez</h1>
		{#if login_state.signup_mode}
			<div class="flex flex-col">
				<input
					name="username"
					bind:value={login_state.user_creds.username}
					type="text"
					placeholder="Username"
					class="mb-2 rounded-md border border-slate-300 p-2 focus:border-slate-500 focus:outline-none dark:border-slate-600 dark:bg-slate-800 dark:text-slate-200 dark:focus:border-slate-400"
				/>
				<input
					name="password"
					bind:value={login_state.user_creds.password}
					type="password"
					placeholder="Password"
					class="mb-2 rounded-md border border-slate-300 p-2 focus:border-slate-500 focus:outline-none dark:border-slate-600 dark:bg-slate-800 dark:text-slate-200 dark:focus:border-slate-400"
				/>
				<input
					name="email"
					bind:value={login_state.user_creds.email}
					type="email"
					placeholder="Email"
					class="mb-2 rounded-md border border-slate-300 p-2 focus:border-slate-500 focus:outline-none dark:border-slate-600 dark:bg-slate-800 dark:text-slate-200 dark:focus:border-slate-400"
				/>
				<button
					class="mt-4 cursor-pointer rounded-md bg-slate-600 p-2 text-white hover:bg-slate-700 dark:bg-slate-500 dark:hover:bg-slate-600"
					onclick={signUp}
					disabled={login_state.posting}
				>
					Sign Up
				</button>
				<!-- Success/Error Messages -->
				{#if login_state.success}
					<div in:fade out:fade={{ duration: 400 }}>
						<Alert color="green" class="mt-2">
							<span>{login_state.success}</span>
						</Alert>
					</div>
				{/if}
				{#if login_state.error}
					<div in:fade out:fade={{ duration: 400 }}>
						<Alert color="red" class="mt-2">
							<span>{login_state.error}</span>
						</Alert>
					</div>
				{/if}
			</div>
			<div class="flex w-full flex-col text-center">
				<p class="mt-4 text-slate-600 dark:text-slate-300">
					Already have an account?
					<button
						class="cursor-pointer text-slate-800 underline hover:text-slate-600 dark:text-slate-200 dark:hover:text-slate-400"
						onclick={() => (login_state.signup_mode = false)}
					>
						Sign in
					</button>
				</p>
			</div>
		{:else}
			<form method="POST" use:enhance class="flex flex-col">
				<input
					name="username"
					bind:value={login_state.user_creds.username}
					type="text"
					placeholder="Username"
					class="mb-2 rounded-md border border-slate-300 p-2 focus:border-slate-500 focus:outline-none dark:border-slate-600 dark:bg-slate-800 dark:text-slate-200 dark:focus:border-slate-400"
				/>
				<input
					name="password"
					bind:value={login_state.user_creds.password}
					type="password"
					placeholder="Password"
					class="cursor-pointer rounded-md border border-slate-300 p-2 focus:border-slate-500 focus:outline-none dark:border-slate-600 dark:bg-slate-800 dark:text-slate-200 dark:focus:border-slate-400"
				/>
				<button
					type="submit"
					class="mt-4 cursor-pointer rounded-md bg-slate-600 p-2 text-white hover:bg-slate-700 dark:bg-slate-500 dark:hover:bg-slate-600"
				>
					Login
				</button>
				<!-- Server-side Error Messages from form action -->
				{#if form?.error}
					<div in:fade out:fade={{ duration: 400 }}>
						<Alert color="red" class="mt-2">
							<span>{form.error}</span>
						</Alert>
					</div>
				{/if}
				<!-- Client-side Success/Error Messages -->
				{#if login_state.success}
					<div in:fade out:fade={{ duration: 400 }}>
						<Alert color="green" class="mt-2">
							<span>{login_state.success}</span>
						</Alert>
					</div>
				{/if}
				{#if login_state.error}
					<div in:fade out:fade={{ duration: 400 }}>
						<Alert color="red" class="mt-2">
							<span>{login_state.error}</span>
						</Alert>
					</div>
				{/if}
			</form>
			<div class="flex w-full flex-col text-center">
				<p class="mt-4 text-slate-600 dark:text-slate-300">
					Don't have an account?
					<button
						class="cursor-pointer text-slate-800 underline hover:text-slate-600 dark:text-slate-200 dark:hover:text-slate-400"
						onclick={() => (login_state.signup_mode = true)}
					>
						Sign up here.
					</button>
				</p>
			</div>
		{/if}
	</div>
</div>
