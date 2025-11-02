import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import crypto from 'node:crypto'; // Node.js 16+ style import

// Alternative: Using Node.js crypto directly
function createSHA1HashNode(message: string): string {
  return crypto.createHash('sha1').update(message).digest('hex');
}

// Example usage:
// const signature = createSHA1HashNode(paramsToSign);