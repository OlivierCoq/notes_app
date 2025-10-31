import { fail, redirect } from '@sveltejs/kit';

import { PUBLIC_API_URL } from '$env/static/public'

export const actions = {


  default: async ({ request, fetch, cookies, url }) => {
    const form = await request.formData();
    const username = String(form.get('username') ?? '')
    const password = String(form.get('password') ?? '')


    // console.log('Logging in with', { username, password });


    const res = await fetch(`${PUBLIC_API_URL}/tokens/authentication`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username, password })
    });
    if (!res.ok) {
      // locals error interface
      let errorMessage = 'Incorrect username or password. bummer.';
      try {
        const errorData = await res.json();
        console.error('Login failed:', errorData);
        // errorMessage += errorData.message || errorData.error || 'Login failed';
      } catch (e) {
        console.error('Failed to parse error response:', e);
        if (res.status === 401) {
          errorMessage = 'Invalid username or password';
        } else if (res.status === 500) {
          errorMessage = 'Server error. Please try again later.';
        }
      }
      return fail(res.status, {
        error: errorMessage,
        username: username
      });
    }
    const { auth_token } = await res.json();
    console.log('Received auth token:', auth_token)
    cookies.set('auth_token', auth_token, {
      httpOnly: true,
      path: '/',
      sameSite: 'lax',
      secure: true
    });


    const redirectTo = url.searchParams.get('redirectTo') || '/dashboard';
    throw redirect(303, redirectTo);
  }
}