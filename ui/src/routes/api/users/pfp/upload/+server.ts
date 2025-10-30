import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Import environment variables for Cloudinary
import { PUBLIC_API_URL } from '$env/static/public';
import { CLOUDINARY_API_KEY, CLOUDINARY_API_SECRET, CLOUDINARY_CLOUD_NAME } from '$env/static/private';



const getAuthToken = (cookieValue: string | undefined): string => {
  return cookieValue || '';
};

// Get file from frontend request, upload to Cloudinary, return URL
export const POST: RequestHandler = async ({ request, cookies }) => {
  try {
    const formData = await request.formData();
    const file = formData.get('pfp');

    console.log('Received file for upload:', file);

    if (!file || !(file instanceof Blob)) {
      return json({ error: 'No file uploaded' }, { status: 400 });
    }

    // Upload file to Cloudinary using signed upload POST request:
    const uploadUrl = `https://api.cloudinary.com/v1_1/${CLOUDINARY_CLOUD_NAME}/image/upload`;
    let signature = ``

    return json({ message: 'File received, upload logic not implemented yet.' }, { status: 200 });
  } catch (error) {
    console.error('Error uploading profile picture:', error);
    return json({ error: 'Failed to upload profile picture' }, { status: 500 });
  }
}
