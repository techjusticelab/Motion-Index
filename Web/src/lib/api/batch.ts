import { API_URL, getAuthHeaders, handleApiError } from './config';
import type { 
  BatchJob, 
  BatchClassifyRequest, 
  ApiResponse 
} from './types';

/**
 * Start a batch classification job
 */
export async function startBatchClassification(
  request: BatchClassifyRequest,
  session?: any
): Promise<{
  job_id: string;
  status: string;
  total_documents: number;
  created_at: string;
}> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await fetch(`${API_URL}/api/v1/batch/classify`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authHeaders
      },
      body: JSON.stringify(request)
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    console.log('Batch classification start response:', apiResponse);
    
    if (apiResponse.success === true && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error('Failed to start batch classification - invalid response format');
    }
  } catch (error) {
    return handleApiError(error, 'start batch classification');
  }
}

/**
 * Get batch job status and progress
 */
export async function getBatchJobStatus(
  jobId: string,
  session?: any
): Promise<BatchJob> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await fetch(`${API_URL}/api/v1/batch/${jobId}/status`, {
      headers: authHeaders
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    
    if (apiResponse.success === true && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error('Failed to get batch job status - invalid response format');
    }
  } catch (error) {
    return handleApiError(error, 'get batch job status');
  }
}

/**
 * Get completed batch job results
 */
export async function getBatchJobResults(
  jobId: string,
  session?: any
): Promise<{
  job_id: string;
  status: string;
  progress: any;
  results: any[];
  completed_at: string;
}> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await fetch(`${API_URL}/api/v1/batch/${jobId}/results`, {
      headers: authHeaders
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    
    if (apiResponse.success === true && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error('Failed to get batch job results - invalid response format');
    }
  } catch (error) {
    return handleApiError(error, 'get batch job results');
  }
}

/**
 * Cancel a running batch job
 */
export async function cancelBatchJob(
  jobId: string,
  session?: any
): Promise<{
  job_id: string;
  status: string;
}> {
  try {
    const authHeaders = await getAuthHeaders(session);
    const response = await fetch(`${API_URL}/api/v1/batch/${jobId}`, {
      method: 'DELETE',
      headers: authHeaders
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const apiResponse = await response.json();
    
    if (apiResponse.success === true && apiResponse.data) {
      return apiResponse.data;
    } else {
      throw new Error('Failed to cancel batch job - invalid response format');
    }
  } catch (error) {
    return handleApiError(error, 'cancel batch job');
  }
}

/**
 * Helper to create batch classification request from file list
 */
export function createBatchClassifyRequest(
  documents: Array<{ id: string; path?: string; text?: string }>,
  options?: Record<string, any>
): BatchClassifyRequest {
  return {
    documents: documents.map(doc => ({
      document_id: doc.id,
      document_path: doc.path,
      text: doc.text
    })),
    options: options || {}
  };
}

/**
 * Poll batch job status until completion
 */
export async function pollBatchJob(
  jobId: string,
  onProgress?: (job: BatchJob) => void,
  pollInterval = 2000,
  session?: any
): Promise<BatchJob> {
  return new Promise((resolve, reject) => {
    const poll = async () => {
      try {
        const job = await getBatchJobStatus(jobId, session);
        
        if (onProgress) {
          onProgress(job);
        }

        if (job.status === 'completed' || job.status === 'failed' || job.status === 'cancelled') {
          resolve(job);
        } else {
          setTimeout(poll, pollInterval);
        }
      } catch (error) {
        reject(error);
      }
    };

    poll();
  });
}