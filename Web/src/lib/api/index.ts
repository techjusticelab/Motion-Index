// Main API exports - organized by functionality

// Configuration and utilities
export { API_URL, getAuthHeaders, handleApiError } from './config';

// Type definitions
export type {
  SearchParams,
  Document,
  SearchResponse,
  MetadataField,
  DocumentStats,
  RedactionAnalysis,
  ApiResponse,
  LegacyApiResponse
} from './types';

// Search API
export {
  searchDocuments,
  getDocumentTypes,
  getLegalTags,
  getMetadataFieldValues,
  getAllFieldOptions,
  getDocumentStats,
  getMetadataFields
} from './search';

// Documents API
export {
  categoriseDocument,
  updateDocumentMetadata,
  getDocumentUrl,
  downloadDocument,
  getDocument
} from './documents';

// Redaction API
export {
  analyzeRedactionsOnly,
  createRedactedDocument,
  getDocumentRedactionAnalysis
} from './redaction';

// Legacy re-exports for backward compatibility
// These will be removed in a future version
export * as api from './client';