// Centralized type exports

// Re-export API types
export type {
  SearchParams,
  Document,
  SearchResponse,
  MetadataField,
  DocumentStats,
  RedactionAnalysis,
  ApiResponse,
  BatchJob,
  BatchProgress,
  BatchResult,
  BatchClassifyRequest,
  BatchDocumentInput,
  ClassificationResult,
  StorageDocument,
  StorageStats
} from '../api/types';

// Re-export search types (for backward compatibility)
export * from './search_types';

// Re-export Supabase types
export * from './supabase_types';

// Common UI types
export interface LoadingState {
  isLoading: boolean;
  error: string | null;
  success: boolean;
}

export interface ModalState {
  isOpen: boolean;
  title?: string;
  content?: string;
}

export interface FormValidation {
  isValid: boolean;
  errors: Record<string, string>;
  touched: Record<string, boolean>;
}

export interface PaginationOptions {
  page: number;
  perPage: number;
  total: number;
  hasNext: boolean;
  hasPrevious: boolean;
}

// Event types for custom dispatchers  
export interface CustomEvents {
  search: CustomEvent<import('../api/types').SearchParams>;
  reset: CustomEvent<void>;
  select: CustomEvent<import('../api/types').Document>;
  delete: CustomEvent<string>;
  update: CustomEvent<{ id: string; data: any }>;
}

// Component prop interfaces
export interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'danger' | 'warning' | 'outline';
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  loading?: boolean;
  type?: 'button' | 'submit' | 'reset';
}

export interface InputProps {
  value?: string;
  placeholder?: string;
  disabled?: boolean;
  required?: boolean;
  error?: string;
  type?: 'text' | 'email' | 'password' | 'number' | 'tel' | 'url';
}