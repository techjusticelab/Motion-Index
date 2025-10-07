// API types and interfaces - Updated to match Fiber backend

export interface SearchParams {
  query?: string;
  doc_type?: string;
  case_number?: string;
  case_name?: string;
  judge?: string[];
  court?: string[];
  author?: string;
  status?: string;
  legal_tags?: string[];
  legal_tags_match_all?: boolean;
  date_range?: DateRange;
  size?: number;
  from?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
  include_highlights?: boolean;
  fuzzy_search?: boolean;
  filters?: SearchFilters;
  sort?: SortOptions;
  highlight?: HighlightOptions;
  // Legacy support
  page?: number;
  limit?: number;
  tags?: string[];
  use_fuzzy?: boolean;
}

export interface SearchFilters {
  doc_type?: string[];
  court?: string[];
  judge?: string[];
  author?: string[];
  status?: string[];
  legal_tags?: string[];
  date_range?: DateRange;
  custom_filters?: Record<string, any>;
}

export interface SortOptions {
  field: string;
  order: 'asc' | 'desc';
}

export interface HighlightOptions {
  fields: string[];
}

export interface DateRange {
  start?: string;
  end?: string;
}

export interface Document {
  id: string;
  file_name: string;
  file_path: string;
  file_url?: string;
  s3_uri?: string;
  text: string;
  doc_type: string;
  category?: string;
  hash?: string;
  created_at: string;
  updated_at?: string;
  metadata: DocumentMetadata;
  size?: number;
  content_type?: string;
  // Legacy fields for backward compatibility
  title?: string;
  content?: string;
  highlight?: {
    text?: string[];
  };
}

export interface DocumentMetadata {
  // Basic Information
  document_name: string;
  subject: string;
  summary?: string;
  document_type?: DocumentType;
  
  // Case Information
  case?: CaseInfo;
  
  // Court Information
  court?: CourtInfo;
  
  // People & Parties
  parties?: Party[];
  attorneys?: Attorney[];
  judge?: Judge;
  
  // Enhanced Date Fields for Legal Documents
  filing_date?: string;   // When document was filed with court
  event_date?: string;    // Key event or action date
  hearing_date?: string;  // Scheduled court hearing date
  decision_date?: string; // When court decision was made
  served_date?: string;   // When documents were served
  timestamp?: string;     // Legacy field
  status?: string;
  
  // Document Properties
  language?: string;
  pages?: number;
  word_count?: number;
  
  // Legal Classification
  legal_tags?: string[];
  charges?: Charge[];
  authorities?: Authority[];
  
  // Processing Metadata
  processed_at: string;
  confidence?: number;
  ai_classified: boolean;
  classification_confidence?: number;
  
  // Document Analysis
  sensitive_terms?: string[];
  has_redactions?: boolean;
  redaction_score?: number;
  extraction_method?: string;
  file_type?: string;
  
  // Legacy fields for backward compatibility
  case_name?: string;
  case_number?: string;
  author?: string;
}

export interface CaseInfo {
  case_number: string;
  case_name: string;
  case_type?: string;
  chapter?: string;
  docket?: string;
  nature_of_suit?: string;
}

export interface CourtInfo {
  court_id: string;
  court_name: string;
  jurisdiction: string;
  level: string;
  district?: string;
  division?: string;
  county?: string;
}

export interface Party {
  name: string;
  role: string;
  party_type?: string;
  date?: string;
}

export interface Attorney {
  name: string;
  bar_number?: string;
  role: string;
  organization?: string;
  contact_info?: string;
}

export interface Judge {
  name: string;
  title?: string;
  judge_id?: string;
}

export interface Charge {
  statute: string;
  description: string;
  grade?: string;
  class?: string;
  count?: number;
}

export interface Authority {
  citation: string;
  case_title?: string;
  type: string;
  precedent: boolean;
  page?: string;
}

export type DocumentType = 
  | 'motion_to_suppress'
  | 'motion_to_dismiss'
  | 'motion_to_compel'
  | 'motion_in_limine'
  | 'motion_summary_judgment'
  | 'motion_to_strike'
  | 'motion_for_reconsideration'
  | 'motion_to_amend'
  | 'motion_for_continuance'
  | 'order'
  | 'ruling'
  | 'judgment'
  | 'sentence'
  | 'injunction'
  | 'brief'
  | 'complaint'
  | 'answer'
  | 'plea'
  | 'reply'
  | 'docket_entry'
  | 'notice'
  | 'stipulation'
  | 'correspondence'
  | 'transcript'
  | 'evidence'
  | 'other'
  | 'unknown';

// Fiber backend search response format
export interface SearchResponse {
  total: number; // Frontend expects this format
  hits: Document[]; // Frontend expects this format
  aggregations?: Record<string, any>;
}

// Fiber backend search result format (what the API actually returns)
export interface FiberSearchResult {
  total_hits: number;
  max_score?: number;
  documents: SearchDocument[];
  aggregations?: Record<string, any>;
  took_ms: number;
  timed_out: boolean;
}

// Fiber backend search document format
export interface SearchDocument {
  id: string;
  score?: number;
  document: Record<string, any>;
  highlights?: Record<string, string[]>;
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
  index_size?: string;
  type_counts?: TypeCount[];
  tag_counts?: TagCount[];
  last_updated?: string;
  field_stats?: Record<string, FieldStat>;
}

// New Fiber backend types
export interface TypeCount {
  type: string;
  count: number;
}

export interface TagCount {
  tag: string;
  count: number;
}

export interface FieldStat {
  unique_values: number;
  total_values: number;
}

export interface FieldValue {
  value: string;
  count: number;
  last_seen?: string;
}

export interface LegalTagsResponse {
  document_types: Array<{
    type: string;
    label: string;
    count: number;
    subcategories?: Array<{
      type: string;
      label: string;
      count: number;
    }>;
  }>;
  practice_areas: Array<{
    area: string;
    label: string;
    count: number;
  }>;
}

export interface DocumentTypesResponse {
  types: Array<{
    id: string;
    name: string;
    description: string;
    count: number;
    common_subtypes?: string[];
    typical_length?: string;
    filing_requirements?: string;
  }>;
}

export interface FieldOptionsResponse {
  filterable_fields: Array<{
    field: string;
    type: string;
    values?: string[];
    searchable?: boolean;
    description: string;
  }>;
  sortable_fields: Array<{
    field: string;
    default?: boolean;
    description: string;
  }>;
}

export interface MetadataFieldValuesResponse {
  field: string;
  total: number;
  values: FieldValue[];
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

// Fiber backend API response format
export interface ApiResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: ApiError;
  request_id?: string;
  timestamp: string;
}

// Fiber backend error structure
export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, any>;
  field?: string;
}

// Legacy wrapper for backward compatibility
export interface LegacyApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

// Batch Processing Types
export interface BatchJob {
  id: string;
  type: string;
  status: 'queued' | 'running' | 'completed' | 'failed' | 'cancelled';
  progress: BatchProgress;
  results?: BatchResult[];
  error?: string;
  created_at: string;
  updated_at: string;
  completed_at?: string;
  options: Record<string, any>;
}

export interface BatchProgress {
  total_documents: number;
  processed_count: number;
  success_count: number;
  error_count: number;
  skipped_count: number;
  indexed_count: number;
  index_error_count: number;
  percent_complete: number;
  estimated_duration?: string;
}

export interface BatchResult {
  document_id: string;
  document_path: string;
  status: 'success' | 'error' | 'skipped';
  classification_result?: ClassificationResult;
  error?: string;
  indexed: boolean;
  index_error?: string;
  index_id?: string;
  processed_at: string;
}

export interface ClassificationResult {
  document_type: string;
  legal_category?: string;
  subject?: string;
  summary?: string;
  confidence: number;
  success: boolean;
  error?: string;
  // Enhanced date extraction fields
  filing_date?: string;
  event_date?: string;
  hearing_date?: string;
  decision_date?: string;
  served_date?: string;
  // Additional metadata
  legal_tags?: string[];
  case_number?: string;
  case_name?: string;
  court?: string;
  judge?: string;
  parties?: Party[];
  attorneys?: Attorney[];
}

export interface BatchClassifyRequest {
  documents: BatchDocumentInput[];
  options?: Record<string, any>;
}

export interface BatchDocumentInput {
  document_id: string;
  document_path?: string;
  text?: string;
}

// Storage Management Types
export interface StorageDocument {
  path: string;
  name: string;
  size: number;
  modified: string;
  url?: string;
}

export interface StorageStats {
  total_documents: number;
  total_size: number;
  storage_backend: string;
  last_updated: string;
}