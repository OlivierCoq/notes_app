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
        const noteData = await request.json();

        // Here you would typically call your createNote function
        // const newNote = await createNote(noteData);
        const response = await fetch(`${PUBLIC_API_URL}/notes`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getAuthToken(cookies.get('auth_token'))}`
            },
            body: JSON.stringify(noteData)
        });

        if (!response.ok) {
            const errorData = await response.json();
            return json({ error: errorData }, { status: response.status });
        }

        const data = await response.json();
        return json(data);
    } catch (error) {
        console.error('Error adding note:', error);
        return json({ error: 'Failed to add note' }, { status: 500 });
    }
}