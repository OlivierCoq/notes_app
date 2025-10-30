
import type { PageServerLoad } from './$types';
import { PUBLIC_API_URL } from '$env/static/public';

//

export const load: PageServerLoad = async ({ locals, fetch }) => {
  // Access the user from locals (set by hooks.server.ts)
  const user = locals.user;

  if (!user) {
    // This should be handled by the layout, but just in case
    return {
      notes: [],
      error: 'User not authenticated'
    };
  }


};