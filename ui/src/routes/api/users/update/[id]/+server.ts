import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { PUBLIC_API_URL } from '$env/static/public';

const getAuthToken = (cookieValue: string | undefined): string => {
  return cookieValue || '';
};

export const PATCH: RequestHandler = async ({ params, request, cookies }) => {
  try {
    const userId = params.id;
    const formData = await request.json();
    console.log('from update endpoint:', formData);


    const authToken = getAuthToken(cookies.get('auth_token'));

    const response = await fetch(`${PUBLIC_API_URL}/users/${userId}`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken}`
      },
      body: JSON.stringify(formData)
    });

    if (!response.ok) {
      const errorData = await response.json();
      return json({ error: errorData }, { status: response.status });
    }
    const data = await response.json();
    console.log('User updated successfully:', data);

    return json(data);
  } catch (error) {
    console.error('Error updating user:', error);
    return json({ error: 'Failed to update user' }, { status: 500 });
  }
}