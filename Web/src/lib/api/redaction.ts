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
    
    const response = await fetch(`${API_URL}/api/v1/analyze-redactions`, {
      method: 'POST',
      headers: {
        ...authHeaders
        // Don't set Content-Type for FormData - let browser set it with boundary
      },
      body: formData
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    console.log("Redaction analysis response:", apiResponse);
    if (apiResponse.success === true && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error('Redaction analysis failed - invalid response format');
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
    const response = await fetch(`${API_URL}/api/v1/redact-document`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authHeaders
      },
      body: JSON.stringify({
        document_id: documentId,
        apply_redactions: applyRedactions
      })
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    console.log("Redaction response:", apiResponse);
    if (apiResponse.success === true && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error('Document redaction failed - invalid response format');
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
    const response = await fetch(`${API_URL}/api/v1/documents/${documentId}/redactions`, {
      headers: authHeaders
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    if (apiResponse.success === true && apiResponse.data) {
      return apiResponse.data.redaction_analysis || null;
    } else {
      return null;
    }
  } catch (error) {
    console.warn('No redaction analysis found for document:', documentId);
    return null;
  }
}