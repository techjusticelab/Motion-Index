import axios from 'axios';

const API_URL = 'http://localhost:8000';

// Define types
export interface SearchParams {
  query?: string;
  doc_type?: string;
  case_number?: string;
  case_name?: string;
  judge?: string;
  court?: string;
  author?: string;
  status?: string;
  date_range?: {
    start?: string;
    end?: string;
  };
  size?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
  page?: number;
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

// API client functions
export async function searchDocuments(params: SearchParams): Promise<SearchResponse> {
  try {
    const response = await axios.post(`${API_URL}/search`, params);
    return response.data;
  } catch (error) {
    console.error('Error searching documents:', error);
    throw error;
  }
}

export async function getDocumentTypes(): Promise<Record<string, number>> {
  try {
    const response = await axios.get(`${API_URL}/document-types`);
    return response.data;
  } catch (error) {
    console.error('Error getting document types:', error);
    throw error;
  }
}

export async function getMetadataFieldValues(field: string, prefix?: string, size: number = 20): Promise<string[]> {
  try {
    const response = await axios.post(`${API_URL}/metadata-field-values`, {
      field,
      prefix,
      size
    });
    return response.data;
  } catch (error) {
    console.error(`Error getting metadata field values for ${field}:`, error);
    throw error;
  }
}

export async function getAllFieldOptions(): Promise<Record<string, string[]>> {
  try {
    const response = await axios.get(`${API_URL}/all-field-options`);
    return response.data;
  } catch (error) {
    console.error('Error getting all field options:', error);
    throw error;
  }
}

export async function getDocumentStats(): Promise<any> {
  try {
    const response = await axios.get(`${API_URL}/document-stats`);
    return response.data;
  } catch (error) {
    console.error('Error getting document stats:', error);
    throw error;
  }
}

export async function getMetadataFields(): Promise<{ fields: MetadataField[] }> {
  try {
    const response = await axios.get(`${API_URL}/metadata-fields`);
    return response.data;
  } catch (error) {
    console.error('Error getting metadata fields:', error);
    throw error;
  }
}
