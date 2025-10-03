// API Configuration for local development
export const API_URL = 'http://localhost:8003';

/**
 * Get authentication headers for API requests
 */
export async function getAuthHeaders(session?: any) {
  if (session?.access_token) {
    return {
      'Authorization': `Bearer ${session.access_token}`,
      'Content-Type': 'application/json'
    };
  }
  
  return {
    'Content-Type': 'application/json'
  };
}

/**
 * Handle API errors consistently
 */
export function handleApiError(error: any, operation: string): never {
  console.error(`API Error in ${operation}:`, error);
  
  if (error.response?.data?.message) {
    throw new Error(error.response.data.message);
  }
  
  if (error.message) {
    throw new Error(error.message);
  }
  
  throw new Error(`Failed to ${operation}`);
}