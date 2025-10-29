import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { PUBLIC_API_URL } from '$env/static/public';
// import { createNote } from '$lib/server/api';

// Helper function to get auth token from cookies
const getAuthToken = (cookieValue: string | undefined): string => {
  return cookieValue || '';
};

// Get slug from params to update note for that user

export const PATCH: RequestHandler = async ({ params, cookies, request }) => {
  try {
    // Get the id parameter from the route
    const { id } = params;

    // Validate that id exists
    if (!id) {
      return json({ error: 'Note ID is required' }, { status: 400 });
    }

    const body = await request.json();
    console.log('Request body for updating note:', body);

    // Here you would typically call your updateNote function
    const response = await fetch(`${PUBLIC_API_URL}/notes/${id}`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken(cookies.get('auth_token'))}`
      },
      body: JSON.stringify(body)
    });

    if (!response.ok) {
      const errorData = await response.json();
      return json({ error: errorData }, { status: response.status });
    }
    const data = await response.json();
    return json(data);
  } catch (error) {
    console.error('Error updating note:', error);
    return json({ error: 'Failed to update note' }, { status: 500 });
  }
} 