import { json } from '@sveltejs/kit';

export async function GET() {
  // Return an empty JSON response for Chrome DevTools requests
  // This prevents the 404 error from appearing in logs
  return json({}, { status: 200 });
}