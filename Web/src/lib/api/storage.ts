import { API_URL, getAuthHeaders, handleApiError } from './config';
import type { StorageDocument, StorageStats, ApiResponse } from './types';

/**
 * List all documents in storage
 */
export async function listStorageDocuments(
  session?: any
): Promise<StorageDocument[]> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await fetch(`${API_URL}/api/v1/storage/documents`, {
      headers: authHeaders
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    console.log('Storage documents response:', apiResponse);
    
    if (apiResponse.success === true && apiResponse.data) {
      // Handle array response or object with documents array
      if (Array.isArray(apiResponse.data)) {
        return apiResponse.data;
      } else if (apiResponse.data.documents) {
        return apiResponse.data.documents;
      } else {
        return [];
      }
    } else {
      throw new Error('Failed to list storage documents - invalid response format');
    }
  } catch (error) {
    return handleApiError(error, 'list storage documents');
  }
}

/**
 * Get total document count in storage
 */
export async function getStorageDocumentsCount(
  session?: any
): Promise<{
  total_documents: number;
  storage_backend: string;
}> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await fetch(`${API_URL}/api/v1/storage/documents/count`, {
      headers: authHeaders
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    console.log('Storage count response:', apiResponse);
    
    if (apiResponse.success === true && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error('Failed to get storage documents count - invalid response format');
    }
  } catch (error) {
    return handleApiError(error, 'get storage documents count');
  }
}

/**
 * Get storage statistics and overview
 */
export async function getStorageStats(session?: any): Promise<StorageStats> {
  try {
    // Get both count and list to compile stats
    const [countData, documents] = await Promise.all([
      getStorageDocumentsCount(session),
      listStorageDocuments(session)
    ]);

    // Calculate total size from documents if available
    const totalSize = documents.reduce((sum, doc) => sum + (doc.size || 0), 0);

    return {
      total_documents: countData.total_documents,
      total_size: totalSize,
      storage_backend: countData.storage_backend,
      last_updated: new Date().toISOString()
    };
  } catch (error) {
    return handleApiError(error, 'get storage statistics');
  }
}

/**
 * Search for documents by name pattern
 */
export async function searchStorageDocumentsByName(
  namePattern: string,
  session?: any
): Promise<StorageDocument[]> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const params = new URLSearchParams({ name: namePattern });
    const response = await fetch(`${API_URL}/api/v1/files/search?${params}`, {
      headers: authHeaders
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    console.log('Storage search response:', apiResponse);
    
    if (apiResponse.success === true && apiResponse.data) {
      // Convert file paths to storage document format
      const filePaths = Array.isArray(apiResponse.data) ? apiResponse.data : [apiResponse.data];
      return filePaths.map((path: string) => ({
        path,
        name: path.split('/').pop() || path,
        size: 0, // Size not available from search
        modified: new Date().toISOString(),
        url: `${API_URL}/api/v1/files/${path}`
      }));
    } else {
      return [];
    }
  } catch (error) {
    console.error('Error searching storage documents:', error);
    return [];
  }
}

/**
 * Get storage document URL for download/viewing
 */
export function getStorageDocumentUrl(documentPath: string): string {
  // Clean the document path
  let cleanPath = documentPath;
  if (cleanPath.startsWith('/')) {
    cleanPath = cleanPath.substring(1);
  }
  
  return `${API_URL}/api/v1/files/${cleanPath}`;
}

/**
 * Download storage document as blob
 */
export async function downloadStorageDocument(
  documentPath: string,
  session?: any
): Promise<Blob> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const url = getStorageDocumentUrl(documentPath);
    
    const response = await fetch(url, {
      headers: authHeaders
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    return await response.blob();
  } catch (error) {
    return handleApiError(error, 'download storage document');
  }
}

/**
 * Check if document exists in storage
 */
export async function checkStorageDocumentExists(
  documentPath: string,
  session?: any
): Promise<boolean> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const url = getStorageDocumentUrl(documentPath);
    
    const response = await fetch(url, {
      method: 'HEAD',
      headers: authHeaders
    });

    return response.ok;
  } catch (error) {
    console.warn('Error checking document existence:', error);
    return false;
  }
}