import { redirect } from '@sveltejs/kit';

export const load = async ({ locals, url }) => {
  const isProtected = !['/login', '/register', '/'].includes(url.pathname);

  if (isProtected && !locals.user) {
    // Redirect to login if not authenticated
    throw redirect(303, `/login?redirectTo=${url.pathname}`);
  }

  return {
    user: locals.user
  };
}