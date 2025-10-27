// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		interface Locals {
			user: {
				id: number;
				username: string;
				email: string;
				bio: string | null;
				first_name: string;
				last_name: string;
				pfp_url: string | null;
				address_line1: string | null;
				address_line2: string | null;
				city: string | null;
				state: string | null;
				zip: string | null;
				country: string | null;
				created_at: string;
				updated_at: string;
			} | null;
			token: string | null;
		}


		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export { };
