import { redirect } from '@sveltejs/kit';

export const POST = async ({ cookies }) => {
  cookies.set('auth_token', '', { path: '/', httpOnly: true, maxAge: 0, sameSite: 'lax', secure: true });
  throw redirect(303, '/');
};