import { API_URL, getAuthHeaders, handleApiError } from './config';
import type { Document, ApiResponse } from './types';

/**
 * Upload and categorize a document
 */
export async function categoriseDocument(file: File, session?: any): Promise<any> {
  try {
    const formData = new FormData();
    formData.append("file", file);
    const authHeaders = await getAuthHeaders(session);
    
    const response = await fetch(`${API_URL}/api/v1/categorise`, {
      method: 'POST',
      headers: {
        ...authHeaders
        // Don't set Content-Type for FormData - let browser set it with boundary
      },
      body: formData
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    console.log("Categorise response:", apiResponse);
    if (apiResponse.status === 'success' && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error('Document categorization failed - invalid response format');
    }
  } catch (error) {
    return handleApiError(error, 'categorize document');
  }
}

/**
 * Update document metadata
 */
export async function updateDocumentMetadata(
  documentId: string, 
  metadata: any, 
  session?: any
): Promise<Document> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await fetch(`${API_URL}/api/v1/update-metadata`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authHeaders
      },
      body: JSON.stringify({
        document_id: documentId,
        metadata
      })
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    console.log("Metadata update response:", apiResponse);
    if (apiResponse.status === 'success' && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error('Metadata update failed - invalid response format');
    }
  } catch (error) {
    return handleApiError(error, 'update document metadata');
  }
}

/**
 * Get document URL for viewing/downloading using new file serving endpoints
 */
export function getDocumentUrl(document: Document): string | null {
  if (!document) {
    console.warn('No document provided to getDocumentUrl');
    return null;
  }

  // For local development, use the file_url if available and valid
  if (document.file_url && document.file_url.startsWith('http')) {
    console.log('Using local file URL:', document.file_url);
    return document.file_url;
  }

  // Try to construct file serving URL from file_path using new /api/v1/files/ endpoint
  if (document.file_path) {
    // Clean the file path - remove any leading slashes and normalize
    let cleanPath = document.file_path;
    
    // Remove leading slash if present
    if (cleanPath.startsWith('/')) {
      cleanPath = cleanPath.substring(1);
    }
    
    // The new file serving endpoint expects the full path after /api/v1/files/
    // For example: /api/v1/files/data/1385.pdf
    const fileUrl = `${API_URL}/api/v1/files/${cleanPath}`;
    console.log('Using new file serving URL:', fileUrl);
    console.log('Original file_path:', document.file_path);
    console.log('Cleaned path:', cleanPath);
    return fileUrl;
  }

  // Try using file_name to search if no file_path available
  if (document.file_name) {
    // Use the search endpoint to find file by name, then serve it
    const searchUrl = `${API_URL}/api/v1/files/search?name=${encodeURIComponent(document.file_name)}`;
    console.log('File path not available, would need to search by name:', searchUrl);
    // Note: This would require an async call to search first, then serve
    // For now, we'll try a direct file serving approach with just the filename
    const directUrl = `${API_URL}/api/v1/files/data/${encodeURIComponent(document.file_name)}`;
    console.log('Trying direct file URL with filename:', directUrl);
    return directUrl;
  }

  // Try using document ID if available (legacy fallback)
  if (document.id) {
    const idUrl = `${API_URL}/api/v1/files/data/${encodeURIComponent(document.id)}`;
    console.log('Using document ID with file endpoint:', idUrl);
    return idUrl;
  }

  // Fallback to S3 URI if available (for backward compatibility)
  if (document.s3_uri) {
    console.log('Using S3 URI:', document.s3_uri);
    return document.s3_uri;
  }

  console.warn('No valid URL found for document:', document);
  return null;
}

/**
 * Search for files by name using the new file search endpoint
 */
export async function searchFilesByName(fileName: string, session?: any): Promise<string[]> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await fetch(`${API_URL}/api/v1/files/search?name=${encodeURIComponent(fileName)}`, {
      headers: authHeaders
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    console.log('File search response:', apiResponse);
    
    if (apiResponse.status === 'success' && apiResponse.data) {
      // Return array of file paths found
      return Array.isArray(apiResponse.data) ? apiResponse.data : [apiResponse.data];
    } else {
      return [];
    }
  } catch (error) {
    console.error('Error searching files by name:', error);
    return [];
  }
}

/**
 * Get document URL with fallback to file search if needed
 */
export async function getDocumentUrlWithSearch(document: Document, session?: any): Promise<string | null> {
  // Try the regular URL generation first
  const url = getDocumentUrl(document);
  if (url) {
    return url;
  }

  // If no URL could be generated and we have a file name, try searching
  if (document.file_name) {
    console.log('Attempting to search for file by name:', document.file_name);
    const searchResults = await searchFilesByName(document.file_name, session);
    
    if (searchResults.length > 0) {
      // Use the first result found
      const foundPath = searchResults[0];
      console.log('File found via search:', foundPath);
      return `${API_URL}/api/v1/files/${foundPath}`;
    }
  }

  return null;
}

/**
 * Download document as blob
 */
export async function downloadDocument(document: Document, session?: any): Promise<Blob> {
  const url = getDocumentUrl(document);
  if (!url) {
    throw new Error('No valid URL found for document');
  }

  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await fetch(url, {
      headers: authHeaders
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    return await response.blob();
  } catch (error) {
    return handleApiError(error, 'download document');
  }
}

/**
 * Get document by ID
 */
export async function getDocument(documentId: string, session?: any): Promise<Document> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await fetch(`${API_URL}/api/v1/documents/${documentId}`, {
      headers: authHeaders
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    if (apiResponse.status === 'success' && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error('Failed to get document - invalid response format');
    }
  } catch (error) {
    return handleApiError(error, 'get document');
  }
}