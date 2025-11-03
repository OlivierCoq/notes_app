<script lang="ts">
	// Stores

	// Imports
	import USStates from '$lib/assets/data/USStates.json';
	import WorldCountries from '$lib/assets/data/WorldCountries.json';

	//   Svelte
	import Dropzone from 'svelte-file-dropzone';
	import { Alert, Spinner } from 'flowbite-svelte';
	import { fade } from 'svelte/transition';
	//  Components
	import { Editor } from '@tadashi/svelte-editor-quill';

	// Types
	import type { User } from '$lib/types/User';
	import { onMount, tick } from 'svelte';

	// State
	let { data } = $props();
	let user = $derived(data.user as User);
	let error = $derived(data.error);

	let user_account_state = $state({
		posting: false,
		user: {
			id: data?.user?.id || null,
			username: data?.user?.username || '',
			email: data?.user?.email || '',
			first_name: data?.user?.first_name || '',
			last_name: data?.user?.last_name || '',
			bio: data?.user?.bio || '',
			pfp_url: data?.user?.pfp_url || '',
			address_line_1: data?.user?.address_line1 || '',
			address_line_2: data?.user?.address_line2 || '',
			address_city: data?.user?.city || '',
			address_state: data?.user?.state || '',
			address_zip: data?.user?.zip || '',
			address_country: data?.user?.country || ''
		},
		password_update: {
			current_password: '',
			new_password: '',
			confirm_password: '',
			error: '',
			success: ''
		},
		pfp_file: null as File | null,
		saved: false,
		quill: {
			options: {
				placeholder: '',
				plainclipboard: true
			},
			content: {
				html: '',
				text: ''
			}
		},
		success_message: '',
		error_message: ''
	});

	// Lifecycle
	onMount(() => {
		// Prepopulate quill content
		user_account_state.quill.content.html = user_account_state.user.bio;
		user_account_state.quill.content.text = user_account_state.user.bio;
		user_account_state.quill.options.placeholder = user_account_state.user.bio
			? ''
			: 'Tell us about yourself...';
	});

	// Functions
	const handlePfpDrop = async (event: CustomEvent) => {
		user_account_state.posting = true;
		let files = {
			accepted: [],
			rejected: []
		};

		const { acceptedFiles, rejectedFiles } = event.detail;
		files.accepted = acceptedFiles;
		files.rejected = rejectedFiles;

		if (files.accepted.length > 0) {
			user_account_state.pfp_file = files.accepted[0];
		}
		await tick();
		console.log('PFP file set to:', user_account_state.pfp_file);
		if (user_account_state.pfp_file) {
			const formData = new FormData();
			formData.append('pfp', user_account_state?.pfp_file);
			try {
				const response = await fetch(`/api/users/pfp/upload`, {
					method: 'POST',
					body: formData
				});
				if (!response.ok) {
					throw new Error('Failed to upload profile picture');
				}
				const result = await response.json();
				console.log('Profile picture uploaded successfully:', result);
				user_account_state.success_message = 'Profile picture updated successfully!';
				// Update user state with new pfp_url
				user_account_state.user.pfp_url = result.url;
				await tick();
				// Now update the rest of the account info
				await updateAccount();
			} catch (error) {
				console.error('Error uploading profile picture:', error);
			}
		}
	};

	const updateAccount = async () => {
		user_account_state.posting = true;
		try {
			const response = await fetch(`/api/users/update/${data?.user?.id}`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(user_account_state.user)
			});

			if (!response.ok) {
				throw new Error('Failed to update account');
			}

			const result = await response.json();
			console.log('Account updated successfully:', result);
		} catch (error) {
			console.error('Error updating account:', error);
		} finally {
			user_account_state.posting = false;
		}
	};
	const onTextChange = (markup: any, plaintext: any) => {
		user_account_state.user.bio = markup;
		user_account_state.quill.content.html = markup;
		user_account_state.quill.content.text = plaintext;
	};
	const handlePasswordUpdate = async () => {
		// Implement password update logic here
		user_account_state.posting = true;
		user_account_state.password_update.error = '';
		try {
			const response = await fetch(`/api/users/password/${data?.user?.id}`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					id: data?.user?.id,
					new_password: user_account_state.password_update.new_password
				})
			});
			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.message || 'Failed to update password');
				user_account_state.password_update.error = errorData.message || 'Failed to update password';
			}
			const result = await response.json();
			console.log('Password updated successfully:', result);
			user_account_state.password_update.success = 'Password updated successfully!';
			// Clear password fields
			user_account_state.password_update.current_password = '';
			user_account_state.password_update.new_password = '';
			user_account_state.password_update.confirm_password = '';
		} catch (error: any) {
			console.error('Error updating password:', error);
			user_account_state.password_update.error = error.message;
		} finally {
			user_account_state.posting = false;
		}
	};
</script>

<svelte:head>
	<link
		rel="stylesheet"
		href="https://unpkg.com/quill@2.0.3/dist/quill.snow.css"
		crossorigin="anonymous"
	/>
</svelte:head>

<div id="account" class="flex min-h-[95vh] w-full flex-col">
	<div
		class="align-center container mx-auto mt-6 flex h-full flex-col items-center justify-center pt-6"
	>
		<!-- Error Display -->
		{#if error}
			<div class="mb-4 rounded border border-red-400 bg-red-100 px-4 py-3 text-red-700">
				<p><strong>Error:</strong> {error}</p>
			</div>
		{/if}
		<div class="mx-auto mt-6 flex h-full w-full flex-col">
			<h1 class="my-6 text-center text-3xl font-bold text-slate-200">Account Settings</h1>
			<div class="mt-6 flex w-full flex-row">
				<!-- Profile pic -->
				<div class="flex w-full flex-col p-6 pe-6 pt-6 md:w-1/2">
					<div
						id="profile-picture"
						class="mx-6 mb-4 flex h-[250px] w-[250px] flex-col items-center justify-center rounded-full bg-white"
						style="background-image: url({user_account_state.user.pfp_url ||
							'https://res.cloudinary.com/dxsjva9e0/image/upload/v1761835316/user_avatar_ry4fdr.png'}); background-size: cover; background-position: center;"
					>
						<Dropzone
							on:drop={handlePfpDrop}
							class="flex h-full w-full cursor-pointer items-center justify-center rounded-full bg-black/40 hover:bg-black/30"
						>
							{#if user_account_state.posting}
								<Spinner color="blue" size="8" />
							{:else}
								<p class="text-center text-sm text-white">Click or Drag & Drop to Change</p>
							{/if}
						</Dropzone>
					</div>
					<!-- Edit bio -->
					<div
						class="me-6 mt-6 w-full overflow-hidden overflow-scroll rounded-md bg-slate-200 px-6"
					>
						<label class="block">
							<span class="text-slate-200">Bio</span>
							<Editor class="h-96" options={user_account_state.quill.options} {onTextChange}
								>{@html user_account_state.quill.content.html}</Editor
							>
						</label>
					</div>
				</div>
				<!-- Account info -->
				<div class="w-full md:w-1/2">
					<h3 class="mt-6 text-lg font-semibold text-slate-200">Personal info</h3>
					<div class="mt-4 flex w-5/6 flex-col space-y-4 pb-5">
						<label class="block">
							<span class="text-slate-200">Username</span>
							<input
								type="text"
								bind:value={user_account_state.user.username}
								class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
							/>
						</label>
						<label class="block">
							<span class="text-slate-200">Email</span>
							<input
								type="email"
								bind:value={user_account_state.user.email}
								class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
							/>
						</label>

						<div class="flex w-full flex-row">
							<div class="w-1/2">
								<label class="block">
									<span class="text-slate-200">First Name</span>
									<input
										type="text"
										bind:value={user_account_state.user.first_name}
										class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
									/>
								</label>
							</div>
							<div class="w-1/2 ps-4">
								<label class="block">
									<span class="text-slate-200">Last Name</span>
									<input
										type="text"
										bind:value={user_account_state.user.last_name}
										class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
									/>
								</label>
							</div>
						</div>
						<h3 class="text-md mt-6 font-semibold text-slate-200">Address</h3>
						<div class="flex w-full flex-row">
							<div class="w-1/2">
								<label class="block">
									<span class="text-slate-200">Line 1</span>
									<input
										type="text"
										bind:value={user_account_state.user.address_line_1}
										class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
									/>
								</label>
							</div>
							<div class="w-1/2 ps-4">
								<label class="block">
									<span class="text-slate-200">Line 2</span>
									<input
										type="text"
										bind:value={user_account_state.user.address_line_2}
										class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
									/>
								</label>
							</div>
						</div>
						<div class="flex w-full flex-row">
							<div class="w-1/3">
								<label class="block">
									<span class="text-slate-200">City</span>
									<input
										type="text"
										bind:value={user_account_state.user.address_city}
										class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
									/>
								</label>
							</div>
							{#if user_account_state.user.address_country === 'US'}
								<div class="w-1/3 px-2">
									<label class="block">
										<span class="text-slate-200">State</span>
										<select
											bind:value={user_account_state.user.address_state}
											class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
										>
											<option value="state.abbreviation" disabled>Select State</option>
											{#each USStates.data as state}
												<option value={state.abbreviation}>{state.name}</option>
											{/each}
										</select>
									</label>
								</div>
							{:else}
								<div class="w-1/3 px-2">
									<label class="block">
										<span class="text-slate-200">State/Province</span>
										<input
											type="text"
											bind:value={user_account_state.user.address_state}
											class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
										/>
									</label>
								</div>
							{/if}
							<div class="w-1/3 ps-4">
								<label class="block">
									<span class="text-slate-200">ZIP</span>
									<input
										type="text"
										bind:value={user_account_state.user.address_zip}
										class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
									/>
								</label>
							</div>
						</div>
						<div class="flex w-full flex-row">
							<div class="w-3/5">
								<label class="block" for="country-select">
									<span class="text-slate-200">Country</span>
								</label>
								<select
									id="country-select"
									bind:value={user_account_state.user.address_country}
									class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
								>
									<option value="" disabled>Select Country</option>
									{#each WorldCountries.data as country}
										<option value={country.code}>{country.name}</option>
									{/each}
								</select>
							</div>
						</div>
						<h3 class="text-md mt-6 font-semibold text-slate-200">Change password</h3>
						<div class="flex w-full flex-col">
							<div class="w-1/2">
								<label class="block">
									<span class="text-slate-200">New Password</span>
									<input
										type="password"
										bind:value={user_account_state.password_update.new_password}
										class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
									/>
								</label>
							</div>
							<div class="w-1/2">
								<label class="block">
									<span class="text-slate-200">Confirm New Password</span>
									<input
										type="password"
										bind:value={user_account_state.password_update.confirm_password}
										class="mt-1 block w-full rounded border border-slate-600 bg-slate-800 p-2 text-slate-200 focus:border-sky-400 focus:outline-none"
									/>
								</label>
							</div>
							<div class="w-1/2">
								<button
									onclick={handlePasswordUpdate}
									class="btn btn-sm mt-2 cursor-pointer rounded-md bg-sky-500 px-4 py-2 text-white hover:bg-sky-600"
									disabled={user_account_state.posting ||
										user_account_state.password_update.new_password === '' ||
										user_account_state.password_update.new_password !==
											user_account_state.password_update.confirm_password}
								>
									Update Password
								</button>
								<!-- Alert: -->
								{#if user_account_state.password_update.error}
									<div in:fade out:fade={{ duration: 400 }}>
										<Alert color="red" class="mt-2">
											<span>{user_account_state.password_update.error}</span>
										</Alert>
									</div>
								{/if}
								{#if user_account_state.password_update.success}
									<div in:fade out:fade={{ duration: 400 }}>
										<Alert color="green" class="mt-2">
											<span>{user_account_state.password_update.success}</span>
										</Alert>
									</div>
								{/if}
							</div>
						</div>

						<button
							onclick={updateAccount}
							class="mb-5 mt-4 cursor-pointer rounded bg-sky-500 px-4 py-2 text-white hover:bg-sky-600"
							disabled={user_account_state.posting}
						>
							{user_account_state.posting ? 'Saving...' : 'Save Changes'}
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
