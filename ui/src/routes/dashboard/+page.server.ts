
import type { PageServerLoad } from './$types';
import { PUBLIC_API_URL } from '$env/static/public';

//

export const load: PageServerLoad = async ({ locals, fetch }) => {
  // Access the user from locals (set by hooks.server.ts)
  const user = locals.user;
  
  if (!user) {
    // This should be handled by the layout, but just in case
    return {
      notes: [],
      error: 'User not authenticated'
    };
  }

  try {
    console.log('Fetching notes for user ID:', user.id);
    
    // Use the internal fetch to automatically include auth headers
    const res = await fetch(`${PUBLIC_API_URL}/user-notes/${user.id}`, {
      headers: {
        'Authorization': `Bearer ${locals.token}`
      }
    });

    const folder_res = await fetch(`${PUBLIC_API_URL}/user-folders/${user.id}`, {
      headers: {
        'Authorization': `Bearer ${locals.token}`
      }
    });

    const folder_data = await folder_res.json();
    if( !folder_res.ok ) {
      console.error('Failed to fetch folders:', folder_res.status, folder_res.statusText);
    } else {
      console.log('Fetched folders data:', folder_data);
    }

    
    if (!res.ok) {
      console.error('Failed to fetch notes:', res.status, res.statusText);
      return {
        notes: [],
        error: 'Failed to fetch notes'
      };
    }
    
    const data = await res.json();
    if (res.ok) {
      // console.log('Fetched notes data:', data);
      return {
        notes: data.notes,
        folders: folder_data.folders,
        user: user // Pass the user data to the page
      };
    }
  } catch (error) {
    console.error('Error fetching notes:', error);
    return {
      notes: [],
      error: 'Failed to fetch notes'
    };
  }
};