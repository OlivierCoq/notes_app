import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Import environment variables for Cloudinary
import { PUBLIC_API_URL } from '$env/static/public';
import { CLOUDINARY_API_KEY, CLOUDINARY_API_SECRET, CLOUDINARY_CLOUD_NAME } from '$env/static/private';

// SHA-1 hash function using Web Crypto API
async function createSHA1Hash(message: string): Promise<string> {
  const encoder = new TextEncoder();
  const data = encoder.encode(message);
  const hashBuffer = await crypto.subtle.digest('SHA-1', data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  return hashHex;
}



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
    const timestamp = Math.floor(Date.now() / 1000);
    const eager = `w_400,h_400,c_pad`;
    const publicId = file.name.replace(/\.[^/.]+$/, ""); // Remove file extension
    
    // Create params to sign (order matters for Cloudinary)
    const paramsToSign = `eager=${eager}&public_id=${publicId}&timestamp=${timestamp}${CLOUDINARY_API_SECRET}`;
    
    // Create SHA-1 signature
    const signature = await createSHA1Hash(paramsToSign);

    return json({ message: 'File received, upload logic not implemented yet.' }, { status: 200 });
  } catch (error) {
    console.error('Error uploading profile picture:', error);
    return json({ error: 'Failed to upload profile picture' }, { status: 500 });
  }
}
