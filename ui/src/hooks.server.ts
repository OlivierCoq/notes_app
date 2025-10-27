import type { Handle } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';

// API base
import { PUBLIC_API_URL } from '$env/static/public';

const handleAuth: Handle = async ({ event, resolve }) => {
  const token = event.cookies.get('auth_token');

  event.locals.token = token || null;
  event.locals.user = null;

  if (token) {
    try {
      // Call your Go API to hydrate the user for this request
      const res = await fetch(`${PUBLIC_API_URL}/users/me`, {
        headers: { Authorization: `Bearer ${token}` }
      });

      if (res.ok) {
        const user = await res.json();
        event.locals.user = user?.user; // get user data from GO api. The reason it's user.user is because of the way the response is structured in Go
        // console.log('Full user: ', event.locals.user);
      } else {
        // invalid/expired token
        event.locals.token = null;
      }
    } catch {
      event.locals.token = null;
    }
  }

  return resolve(event);
}

export const handleFetch = async ({ event, request, fetch }) => {
  const token = event.locals.token;
  const url = new URL(request.url);

  const isGoAPI =
    url.hostname === PUBLIC_API_URL || // direct calls
    url.pathname.startsWith('/api/');       // your proxy

  if (token && isGoAPI) {
    request = new Request(request, {
      headers: new Headers({
        ...Object.fromEntries(request.headers),
        Authorization: `Bearer ${token}`
      })
    });
  }

  return fetch(request);
};

export const handle = sequence(handleAuth)