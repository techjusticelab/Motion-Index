import axios from 'axios';
import { API_URL, getAuthHeaders, handleApiError } from './config';
import type { RedactionAnalysis, ApiResponse } from './types';

/**
 * Analyze document for redactions only (bypasses Elasticsearch requirement)
 */
export async function analyzeRedactionsOnly(file: File, session?: any): Promise<{
  redaction_analysis?: RedactionAnalysis;
  message?: string;
}> {
  try {
    const formData = new FormData();
    formData.append("file", file);
    const authHeaders = await getAuthHeaders(session);
    
    const response = await axios.post(`${API_URL}/api/v1/analyze-redactions`, formData, {
      headers: {
        "Content-Type": "multipart/form-data",
        ...authHeaders
      },
    });
    const apiResponse: ApiResponse<{redaction_analysis?: RedactionAnalysis; message?: string}> = response.data;
    console.log("Redaction analysis response:", apiResponse);
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error(apiResponse.error?.message || 'Redaction analysis failed');
    }
  } catch (error) {
    return handleApiError(error, 'analyze redactions');
  }
}

/**
 * Create redacted version of a document
 */
export async function createRedactedDocument(
  documentId: string, 
  applyRedactions: boolean = true, 
  session?: any
): Promise<{
  document_id: string;
  redacted_url?: string;
  message: string;
}> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await axios.post(`${API_URL}/api/v1/redact-document`, {
      document_id: documentId,
      apply_redactions: applyRedactions
    }, {
      headers: {
        "Content-Type": "application/json",
        ...authHeaders
      },
    });
    const apiResponse: ApiResponse<{document_id: string; redacted_url?: string; message: string}> = response.data;
    console.log("Redaction response:", apiResponse);
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error(apiResponse.error?.message || 'Document redaction failed');
    }
  } catch (error) {
    return handleApiError(error, 'create redacted document');
  }
}

/**
 * Get redaction analysis for an existing document
 */
export async function getDocumentRedactionAnalysis(
  documentId: string, 
  session?: any
): Promise<RedactionAnalysis | null> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await axios.get(`${API_URL}/api/v1/documents/${documentId}/redactions`, {
      headers: authHeaders
    });
    const apiResponse: ApiResponse<{redaction_analysis?: RedactionAnalysis}> = response.data;
    if (apiResponse.success && apiResponse.data) {
      return apiResponse.data.redaction_analysis || null;
    } else {
      return null;
    }
  } catch (error) {
    console.warn('No redaction analysis found for document:', documentId);
    return null;
  }
}