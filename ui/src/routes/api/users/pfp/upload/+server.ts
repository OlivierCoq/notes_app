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

    // Create params to sign (alphabetical order is required for Cloudinary)
    const params = {
      eager: eager,
      public_id: publicId,
      timestamp: timestamp.toString()
    };

    // Sort parameters alphabetically and create string to sign
    const sortedParams = Object.keys(params)
      .sort()
      .map(key => `${key}=${params[key as keyof typeof params]}`)
      .join('&');

    const stringToSign = `${sortedParams}${CLOUDINARY_API_SECRET}`;
    console.log('String to sign:', stringToSign);

    // Create SHA-1 signature
    const signature = await createSHA1Hash(stringToSign);

    // Create FormData for file upload
    const uploadFormData = new FormData();
    uploadFormData.append('file', file);
    uploadFormData.append('api_key', CLOUDINARY_API_KEY);
    uploadFormData.append('timestamp', timestamp.toString());
    uploadFormData.append('signature', signature);
    uploadFormData.append('eager', eager);
    uploadFormData.append('public_id', publicId);

    const response = await fetch(uploadUrl, {
      method: 'POST',
      body: uploadFormData
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error('Cloudinary upload failed:', errorText);
      return json({ error: `Cloudinary upload failed: ${errorText}` }, { status: response.status });
    }

    const responseData = await response.json();
    console.log('Cloudinary upload response:', responseData);

    // Return the secure URL from Cloudinary response
    return json({
      message: 'File uploaded successfully',
      url: responseData.secure_url,
      public_id: responseData.public_id
    }, { status: 200 });
  } catch (error) {
    console.error('Error uploading profile picture:', error);
    return json({ error: 'Failed to upload profile picture' }, { status: 500 });
  }
}
