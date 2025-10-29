import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { PUBLIC_API_URL } from '$env/static/public';
// import { createNote } from '$lib/server/api';

// Helper function to get auth token from cookies
const getAuthToken = (cookieValue: string | undefined): string => {
  return cookieValue || '';
};

// Get slug from params and fetch folders for that user


export const GET: RequestHandler = async ({ params, cookies }) => {
  try {
    // Get the id parameter from the route
    const { id } = params;

    // Validate that id exists
    if (!id) {
      return json({ error: 'User ID is required' }, { status: 400 });
    }

    // Here you would typically call your getAllFolders function
    const response = await fetch(`${PUBLIC_API_URL}/user-folders/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken(cookies.get('auth_token'))}`
      }
    });

    if (!response.ok) {
      const errorData = await response.json();
      return json({ error: errorData }, { status: response.status });
    }
    const data = await response.json();
    return json(data);
  } catch (error) {
    console.error('Error fetching folders:', error);
    return json({ error: 'Failed to fetch folders' }, { status: 500 });
  }
}
