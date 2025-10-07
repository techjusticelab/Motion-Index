// API Configuration for production deployment
import { PUBLIC_API_URL } from '$env/static/public';

// Default to Go Fiber backend port - the Fiber backend requires PORT environment variable
// For local development, the Fiber backend should be configured with PORT=8003
export const API_URL = PUBLIC_API_URL || 'http://localhost:8003';

/**
 * Get authentication headers for API requests
 */
export async function getAuthHeaders(session?: any): Promise<Record<string, string>> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    // Add headers to bypass ad blockers
    'X-Requested-With': 'XMLHttpRequest',
    'Cache-Control': 'no-cache'
  };
  
  if (session?.access_token) {
    headers['Authorization'] = `Bearer ${session.access_token}`;
  }
  
  return headers;
}

/**
 * Handle API errors consistently for Fiber backend
 */
export function handleApiError(error: any, operation: string): never {
  console.error(`API Error in ${operation}:`, error);
  
  // Handle Fiber ApiResponse error format
  if (error.response?.data?.error?.message) {
    throw new Error(error.response.data.error.message);
  }
  
  // Handle direct ApiResponse error format
  if (error.error?.message) {
    throw new Error(error.error.message);
  }
  
  // Handle legacy error formats
  if (error.response?.data?.message) {
    throw new Error(error.response.data.message);
  }
  
  if (error.message) {
    throw new Error(error.message);
  }
  
  throw new Error(`Failed to ${operation}`);
}