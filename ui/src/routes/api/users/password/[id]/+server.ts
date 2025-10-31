import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { PUBLIC_API_URL } from '$env/static/public';

const getAuthToken = (cookieValue: string | undefined): string => {
  return cookieValue || '';
};

export const PATCH: RequestHandler = async ({ params, request, cookies }) => {
  const authToken = getAuthToken(cookies.get('auth_token'));
  const userId = params.id;

  const body = await request.json();
  console.log('Request body:', body);
  const postObj = {
    id: userId,
    new_password: body.new_password
  }

  const response = await fetch(`${PUBLIC_API_URL}/users/password/${userId}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${authToken}`
    },
    body: JSON.stringify(postObj)
  });

  if (!response.ok) {
    return json({ error: 'Failed to update password' }, { status: 500 });
  }

  return json({ success: true });
};
