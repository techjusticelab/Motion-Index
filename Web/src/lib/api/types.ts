// API types and interfaces

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
  file_url?: string;  // New field for local file access
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
  s3_uri?: string;  // Keep as optional for backward compatibility
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

export interface DocumentStats {
  total_documents: number;
  document_types: Record<string, number>;
  recent_uploads: number;
  storage_size?: string;
}

export interface RedactionAnalysis {
  redactions_found: number;
  sensitive_terms: string[];
  confidence_scores: Record<string, number>;
  redaction_areas: Array<{
    page: number;
    x: number;
    y: number;
    width: number;
    height: number;
  }>;
}

export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: Record<string, any>;
  };
  message?: string;
  timestamp?: string;
}

// Legacy wrapper for backward compatibility
export interface LegacyApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}