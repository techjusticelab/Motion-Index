import axios from 'axios';
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
    
    const response = await axios.post(`${API_URL}/api/v1/categorise`, formData, {
      headers: {
        "Content-Type": "multipart/form-data",
        ...authHeaders
      },
    });
    const apiResponse: ApiResponse<any> = response.data;
    console.log("Categorise response:", apiResponse);
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error(apiResponse.error?.message || 'Document categorization failed');
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
    const response = await axios.post(`${API_URL}/api/v1/update-metadata`, {
      document_id: documentId,
      metadata
    }, {
      headers: {
        "Content-Type": "application/json",
        ...authHeaders
      },
    });
    const apiResponse: ApiResponse<Document> = response.data;
    console.log("Metadata update response:", apiResponse);
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error(apiResponse.error?.message || 'Metadata update failed');
    }
  } catch (error) {
    return handleApiError(error, 'update document metadata');
  }
}

/**
 * Get document URL for viewing/downloading
 */
export function getDocumentUrl(document: Document): string | null {
  if (!document) {
    console.warn('No document provided to getDocumentUrl');
    return null;
  }

  // For local development, use the file_url if available
  if (document.file_url) {
    console.log('Using local file URL:', document.file_url);
    return document.file_url;
  }

  // Try to construct local API URL from file_path
  if (document.file_path) {
    const localUrl = `${API_URL}/api/v1/documents/${encodeURIComponent(document.file_path)}`;
    console.log('Using constructed local URL:', localUrl);
    return localUrl;
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
 * Download document as blob
 */
export async function downloadDocument(document: Document, session?: any): Promise<Blob> {
  const url = getDocumentUrl(document);
  if (!url) {
    throw new Error('No valid URL found for document');
  }

  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await axios.get(url, {
      responseType: 'blob',
      headers: authHeaders
    });
    return response.data;
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
    const response = await axios.get(`${API_URL}/api/v1/documents/${documentId}`, {
      headers: authHeaders
    });
    const apiResponse: ApiResponse<Document> = response.data;
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error(apiResponse.error?.message || 'Failed to get document');
    }
  } catch (error) {
    return handleApiError(error, 'get document');
  }
}