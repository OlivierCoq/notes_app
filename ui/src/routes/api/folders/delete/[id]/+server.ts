// Deleting folder:

import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';
import { PUBLIC_API_URL } from '$env/static/public';

// Helper function to get auth token from cookies
const getAuthToken = (cookieValue: string | undefined): string => {
  return cookieValue || '';
};

// Delete folder by ID
export const DELETE: RequestHandler = async ({ params, cookies }) => {
  try {
    // Get the id parameter from the route
    const id = params.id;

    // Validate that id exists
    if (!id) {
      return json({ error: 'Folder ID is required' }, { status: 400 });
    }

    // Here you would typically call your deleteFolder function
    const response = await fetch(`${PUBLIC_API_URL}/folders/${id}`, {
      method: 'DELETE',
      headers: {
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
    console.error('Error deleting folder:', error);
    return json({ error: 'Failed to delete folder' }, { status: 500 });
  }
}