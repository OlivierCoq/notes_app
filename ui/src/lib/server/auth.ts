import type { JwtPayload } from 'jwt-decode';
import { jwtDecode } from 'jwt-decode';

export type User = {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
}

export function getUserFromToken(token: string): User | null {
  try {
    const decoded = jwtDecode<JwtPayload & { sub?: string; email?: string; first_name?: string; last_name?: string }>(token);
    if (decoded && decoded.sub && decoded.email) {
      return {
        id: decoded.sub.toString(),
        email: decoded.email,
        first_name: decoded.first_name || '',
        last_name: decoded.last_name || ''
      };
    }
  } catch (error) {
    console.error('Error decoding token:', error);
  }
  return null;
}