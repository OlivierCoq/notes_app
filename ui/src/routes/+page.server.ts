import { fail, redirect } from '@sveltejs/kit';

import { PUBLIC_API_URL } from '$env/static/public'

export const actions = {

  // Nothing is coming in request formData:

  default: async ({ request, fetch, cookies, url }) => {
    const form = await request.formData();
    const username = String(form.get('username') ?? '')
    const password = String(form.get('password') ?? '')


    console.log('Logging in with', { username, password });


    const res = await fetch(`${PUBLIC_API_URL}/tokens/authentication`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username, password })
    });
    if (!res.ok) {
      return fail(400, { error: 'Invalid credentials' });
    }
    const { auth_token } = await res.json();
    console.log('Received auth token:', auth_token)
    cookies.set('auth_token', auth_token, {
      httpOnly: true,
      path: '/',
      sameSite: 'lax',
      secure: true
    });

    // Redirect to the home page or intended page after login
    const redirectTo = url.searchParams.get('redirectTo') || '/dashboard';
    throw redirect(303, redirectTo);
  }
}