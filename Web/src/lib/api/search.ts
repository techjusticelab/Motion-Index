import axios from 'axios';
import { API_URL, getAuthHeaders, handleApiError } from './config';
import type { SearchParams, SearchResponse, MetadataField, DocumentStats, ApiResponse } from './types';

/**
 * Search documents with given parameters
 */
export async function searchDocuments(params: SearchParams, session?: any): Promise<SearchResponse> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await axios.post(`${API_URL}/api/v1/search`, params, {
      headers: {
        'Content-Type': 'application/json',
        ...authHeaders
      },
    });
    const apiResponse: ApiResponse<SearchResponse> = response.data;
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error(apiResponse.error?.message || 'Search request failed');
    }
  } catch (error) {
    return handleApiError(error, 'search documents');
  }
}

/**
 * Get document type statistics
 */
export async function getDocumentTypes(): Promise<Record<string, number>> {
  try {
    console.log('Fetching document types from:', `${API_URL}/api/v1/document-types`);
    
    // Use fetch instead of axios to avoid network blocking issues
    const response = await fetch(`${API_URL}/api/v1/document-types`);
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
    
    const data = await response.json();
    console.log('Raw document types response:', data);
    
    // Handle the actual API response format: {data: [{type: string, count: number}], status: "success"}
    if (data && data.status === 'success' && data.data) {
      const docTypesArray = data.data;
      const docTypesMap: Record<string, number> = {};
      
      // Convert array format to object format
      docTypesArray.forEach((item: {type: string, count: number}) => {
        docTypesMap[item.type] = item.count;
      });
      
      console.log('Converted document types:', docTypesMap);
      return docTypesMap;
    } else {
      throw new Error('Invalid response format from document types API');
    }
  } catch (error) {
    console.error('Error fetching document types:', error);
    return handleApiError(error, 'get document types');
  }
}

/**
 * Get available legal tags
 */
export async function getLegalTags(session?: any): Promise<string[]> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await axios.get(`${API_URL}/api/v1/legal-tags`, { headers: authHeaders });
    const apiResponse: ApiResponse<{tags: string[]}> = response.data;
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data.tags || [];
    } else {
      throw new Error(apiResponse.error?.message || 'Failed to get legal tags');
    }
  } catch (error) {
    return handleApiError(error, 'get legal tags');
  }
}

/**
 * Get metadata field values with optional prefix filtering
 */
export async function getMetadataFieldValues(
  field: string, 
  prefix?: string, 
  size: number = 20, 
  session?: any
): Promise<string[]> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const params = new URLSearchParams({ field, size: size.toString() });
    if (prefix) {
      params.append('prefix', prefix);
    }
    
    const response = await axios.get(`${API_URL}/api/v1/metadata-fields/${field}?${params}`, { 
      headers: authHeaders 
    });
    const apiResponse: ApiResponse<{values: string[]}> = response.data;
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data.values || [];
    } else {
      throw new Error(apiResponse.error?.message || 'Failed to get metadata field values');
    }
  } catch (error) {
    return handleApiError(error, 'get metadata field values');
  }
}

/**
 * Get all field options for search filters
 */
export async function getAllFieldOptions(session?: any): Promise<Record<string, string[]>> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await axios.get(`${API_URL}/api/v1/field-options`, { headers: authHeaders });
    const apiResponse: ApiResponse<Record<string, string[]>> = response.data;
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error(apiResponse.error?.message || 'Failed to get field options');
    }
  } catch (error) {
    return handleApiError(error, 'get field options');
  }
}

/**
 * Get document statistics
 */
export async function getDocumentStats(session?: any): Promise<DocumentStats> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await axios.get(`${API_URL}/api/v1/document-stats`, { headers: authHeaders });
    const apiResponse: ApiResponse<DocumentStats> = response.data;
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error(apiResponse.error?.message || 'Failed to get document statistics');
    }
  } catch (error) {
    return handleApiError(error, 'get document statistics');
  }
}

/**
 * Get available metadata fields
 */
export async function getMetadataFields(session?: any): Promise<{ fields: MetadataField[] }> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await axios.get(`${API_URL}/api/v1/metadata-fields`, { headers: authHeaders });
    const apiResponse: ApiResponse<{ fields: MetadataField[] }> = response.data;
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error(apiResponse.error?.message || 'Failed to get metadata fields');
    }
  } catch (error) {
    return handleApiError(error, 'get metadata fields');
  }
}