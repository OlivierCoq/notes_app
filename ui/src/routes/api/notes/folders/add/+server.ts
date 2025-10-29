import { json} from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { PUBLIC_API_URL } from '$env/static/public';
// import { createNote } from '$lib/server/api';

// Helper function to get auth token from cookies
const getAuthToken = (cookieValue: string | undefined): string => {
    return cookieValue || '';
};

export const POST: RequestHandler = async ({ request, cookies }) => {
    try {
        const folderData = await request.json();

        // Here you would typically call your createNote function
        // const newNote = await createNote(folderData);
        const response = await fetch(`${PUBLIC_API_URL}/folders`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getAuthToken(cookies.get('auth_token'))}`
            },
            body: JSON.stringify(folderData)
        });

        if (!response.ok) {
            const errorData = await response.json();
            return json({ error: errorData }, { status: response.status });
        }

        const data = await response.json();
        return json(data);
    } catch (error) {
        console.error('Error adding folder:', error);
        return json({ error: 'Failed to add folder' }, { status: 500 });
    }
}