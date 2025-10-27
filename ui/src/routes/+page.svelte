<script lang="ts">
	// imports
	import { PUBLIC_API_URL } from '$env/static/public';

	// State
	const login_state = $state({
		user_creds: {
			username: '',
			password: ''
		},
		logged_in: false,
		signup_mode: false
	});

	// methods
	const login = async (event: Event) => {
		event.preventDefault();
		try {
			const response = await fetch(`${PUBLIC_API_URL}/tokens/authentication`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Accept: 'application/json'
				},
				body: JSON.stringify(login_state.user_creds)
			});

			// console.log('Login response:', response);

			if (response.ok) {
				login_state.logged_in = true;
				const token = await response.json();
				console.log('response data', token);
				// Handle successful login (e.g., redirect to dashboard)
			} else {
				const errorData = await response.json();
				console.error('Login failed:', errorData.message);
				// Handle login failure (e.g., show error message)
			}
		} catch (error) {
			console.error('Error during login:', error);
			// Handle network or other errors
		}
	};
</script>

<div
	class="align-center mx-auto flex w-96 flex-col justify-center rounded-md bg-slate-100 p-8 shadow-md dark:bg-slate-700"
>
	<h1 class="mb-6 text-center text-2xl font-bold text-slate-800 dark:text-slate-200">notez</h1>
	<form class="flex flex-col">
		<input
			bind:value={login_state.user_creds.username}
			type="text"
			placeholder="Username"
			class="mb-2 rounded-md border border-slate-300 p-2 focus:border-slate-500 focus:outline-none dark:border-slate-600 dark:bg-slate-800 dark:text-slate-200 dark:focus:border-slate-400"
		/>
		<input
			bind:value={login_state.user_creds.password}
			type="password"
			placeholder="Password"
			class="cursor-pointer rounded-md border border-slate-300 p-2 focus:border-slate-500 focus:outline-none dark:border-slate-600 dark:bg-slate-800 dark:text-slate-200 dark:focus:border-slate-400"
		/>
		<button
			onclick={login}
			class="mt-4 cursor-pointer rounded-md bg-slate-600 p-2 text-white hover:bg-slate-700 dark:bg-slate-500 dark:hover:bg-slate-600"
		>
			Login
		</button>
	</form>
</div>
