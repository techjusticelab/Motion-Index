import axios from 'axios';
import { browser } from '$app/environment';
import { getAuthToken, isAuthenticated } from '$lib/auth';
import { get } from 'svelte/store';

// Using a separate API deployment on Vercel
// Using ngrok for HTTPS tunneling to the API server
const API_URL = 'https://rational-evolving-joey.ngrok-free.app';
// const API_URL = 'https://3.88.135.105:8000';
// const API_URL = 'https://172.20.0.2:8000';
//const API_URL = 'https://0.0.0.0:8000'

// Define types
export interface SearchParams {
  query?: string;
  doc_type?: string;
  case_number?: string;
  case_name?: string;
  judge?: string | string[];
  court?: string | string[];
  author?: string;
  status?: string;
  date_range?: {
    start?: string;
    end?: string;
  };
  tags?: string[];
  size?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
  page?: number;
  use_fuzzy?: boolean;
}

export interface Document {
  id: string;
  file_name: string;
  file_path: string;
  text: string;
  doc_type: string;
  category?: string;
  metadata: {
    document_name: string;
    subject: string;
    status?: string;
    timestamp?: string;
    case_name?: string;
    case_number?: string;
    author?: string;
    judge?: string;
    legal_tags?: string[];
    court?: string;
  };
  created_at: string;
  s3_uri?: string;
  highlight?: {
    text?: string[];
  };
}

export interface SearchResponse {
  total: number;
  hits: Document[];
  aggregations?: any;
}

export interface MetadataField {
  id: string;
  name: string;
  type: string;
}

/**
 * Get authentication headers for API requests using our new auth system
 */
async function getAuthHeaders() {
  // Default headers to include the ngrok-skip-browser-warning header
  const defaultHeaders = {
    'ngrok-skip-browser-warning': 'true'
  };
  
  // Check if user is authenticated
  const authenticated = get(isAuthenticated);
  
  if (authenticated && browser) {
    const token = getAuthToken();
    
    if (token) {
      console.log('Using auth token for API request');
      return {
        ...defaultHeaders,
        Authorization: `Bearer ${token}`
      };
    }
  }
  
  console.log('No authenticated user, proceeding without auth token');
  return defaultHeaders;
}

// API client functions
export async function searchDocuments(params: SearchParams): Promise<SearchResponse> {
  try {
    console.log('Searching documents with params:', params);
    
    // Get auth headers since search is now protected
    const headers = await getAuthHeaders();
    
    // Make the request with auth headers
    const response = await axios.post(`${API_URL}/search`, params, { headers });
    return response.data;
  } catch (error) {
    console.error('Error searching documents:', error);
    throw error;
  }
}

export async function getDocumentTypes(): Promise<Record<string, number>> {
  try {
    console.log('Getting document types');
    
    // Make the request without auth headers to ensure it works for all users
    const response = await axios.get(`${API_URL}/document-types`);
    return response.data;
  } catch (error) {
    console.error('Error getting document types:', error);
    throw error;
  }
}

export async function getLegalTags(): Promise<string[]> {
  try {
    const headers = await getAuthHeaders();
    const response = await axios.get(`${API_URL}/legal-tags`, { headers });
    return response.data;
  } catch (error) {
    console.error('Error getting legal tags:', error);
    throw error;
  }
}
export async function getMetadataFieldValues(field: string, prefix?: string, size: number = 20): Promise<string[]> {
  try {
    const headers = await getAuthHeaders();
    const response = await axios.post(`${API_URL}/metadata-field-values`, {
      field,
      prefix,
      size
    }, { headers });
    return response.data;
  } catch (error) {
    console.error(`Error getting metadata field values for ${field}:`, error);
    throw error;
  }
}

export async function getAllFieldOptions(): Promise<Record<string, string[]>> {
  try {
    const headers = await getAuthHeaders();
    const response = await axios.get(`${API_URL}/all-field-options`, { headers });
    return response.data;
  } catch (error) {
    console.error('Error getting all field options:', error);
    throw error;
  }
}

export async function getDocumentStats(): Promise<any> {
  try {
    const headers = await getAuthHeaders();
    const response = await axios.get(`${API_URL}/document-stats`, { headers });
    return response.data;
  } catch (error) {
    console.error('Error getting document stats:', error);
    throw error;
  }
}

export async function getMetadataFields(): Promise<{ fields: MetadataField[] }> {
  try {
    const headers = await getAuthHeaders();
    const response = await axios.get(`${API_URL}/metadata-fields`, { headers });
    return response.data;
  } catch (error) {
    console.error('Error getting metadata fields:', error);
    throw error;
  }
}

export async function categoriseDocument(file: File): Promise<any> {
  try {
    const formData = new FormData();
    formData.append("file", file);
    
    const authHeaders = await getAuthHeaders();
    const response = await axios.post(`${API_URL}/categorise`, formData, {
      headers: {
        "Content-Type": "multipart/form-data",
        ...authHeaders
      },
    });
    console.log("Categorise response:", response.data);
    return response.data;
  } catch (error) {
    console.error("Error categorising document:", error);
    throw error;
  }
}

export async function updateDocumentMetadata(documentId: string, metadata: any): Promise<any> {
  try {
    const headers = await getAuthHeaders();
    const response = await axios.post(`${API_URL}/update-metadata`, {
      document_id: documentId,
      metadata: metadata
    }, { headers });
    return response.data;
  } catch (error) {
    console.error('Error updating document metadata:', error);
    throw error;
  }
}
