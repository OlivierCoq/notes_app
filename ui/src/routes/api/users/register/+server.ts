import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { PUBLIC_API_URL } from '$env/static/public';

const getAuthToken = (cookieValue: string | undefined): string => {
  return cookieValue || '';
};

export const POST: RequestHandler = async ({ request, cookies }) => {
  try {
    const formData = await request.json();
    console.log('from register endpoint:', formData);
    const username = String(formData.username ?? '');
    const password = String(formData.password ?? '');
    const email = String(formData.email ?? '');
    const pfp_url = String(formData.pfp_url ?? '');

    const response = await fetch(`${PUBLIC_API_URL}/users/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username, email, password, pfp_url })
    });

    if (!response.ok) {
      const errorData = await response.json();
      return json({ error: errorData }, { status: response.status });
    }
    const data = await response.json();
    console.log('User registered successfully:', data);
    // Log user in:
    const loginResponse = await fetch(`${PUBLIC_API_URL}/tokens/authentication`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username, password })
    });

    if (!loginResponse.ok) {
      const errorData = await loginResponse.json();
      return json({ error: errorData }, { status: loginResponse.status });
    }

    const { auth_token } = await loginResponse.json();
    console.log('Received auth token:', auth_token);
    cookies.set('auth_token', auth_token, {
      httpOnly: true,
      path: '/',
      sameSite: 'lax',
      secure: true
    });

    return json(data);
  } catch (error) {
    console.error('Error registering user:', error);
    return json({ error: 'Failed to register user' }, { status: 500 });
  }
}
