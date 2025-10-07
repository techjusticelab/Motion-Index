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
  LegacyApiResponse,
  BatchJob,
  BatchProgress,
  BatchResult,
  BatchClassifyRequest,
  BatchDocumentInput,
  ClassificationResult,
  StorageDocument,
  StorageStats
} from './types';

// Search API
export {
  searchDocuments,
  getDocumentTypes,
  getLegalTags,
  getMetadataFieldValues,
  getAllFieldOptions,
  getDocumentStats
} from './search';

// Documents API
export {
  categoriseDocument,
  updateDocumentMetadata,
  getDocumentUrl,
  getDocumentUrlWithSearch,
  searchFilesByName,
  downloadDocument,
  getDocument
} from './documents';

// Redaction API
export {
  analyzeRedactionsOnly,
  createRedactedDocument,
  getDocumentRedactionAnalysis
} from './redaction';

// Batch Processing API
export {
  startBatchClassification,
  getBatchJobStatus,
  getBatchJobResults,
  cancelBatchJob,
  createBatchClassifyRequest,
  pollBatchJob
} from './batch';

// Storage Management API
export {
  listStorageDocuments,
  getStorageDocumentsCount,
  getStorageStats,
  searchStorageDocumentsByName,
  getStorageDocumentUrl,
  downloadStorageDocument,
  checkStorageDocumentExists
} from './storage';

// Legacy re-exports for backward compatibility
// These will be removed in a future version
export * as api from './client';