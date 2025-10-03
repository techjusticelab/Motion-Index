package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type BatchClassifyRequest struct {
	Documents []BatchDocumentInput   `json:"documents"`
	Options   map[string]interface{} `json:"options"`
}

type BatchDocumentInput struct {
	DocumentID   string `json:"document_id"`
	DocumentPath string `json:"document_path,omitempty"`
	Text         string `json:"text,omitempty"`
}

type BatchResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
}

type BatchJob struct {
	ID       string      `json:"id"`
	Status   string      `json:"status"`
	Progress BatchProgress `json:"progress"`
	Results  []BatchResult `json:"results"`
}

type BatchProgress struct {
	TotalDocuments  int     `json:"total_documents"`
	ProcessedCount  int     `json:"processed_count"`
	SuccessCount    int     `json:"success_count"`
	ErrorCount      int     `json:"error_count"`
	SkippedCount    int     `json:"skipped_count"`
	IndexedCount    int     `json:"indexed_count"`
	IndexErrorCount int     `json:"index_error_count"`
	PercentComplete float64 `json:"percent_complete"`
}

type BatchResult struct {
	DocumentID   string `json:"document_id"`
	DocumentPath string `json:"document_path"`
	Status       string `json:"status"`
	Indexed      bool   `json:"indexed"`
	IndexError   string `json:"index_error,omitempty"`
	IndexID      string `json:"index_id,omitempty"`
	Error        string `json:"error,omitempty"`
}

func main() {
	log.Printf("🧪 Testing Batch API with Enhanced Indexing")

	// Create a test batch request
	batchReq := BatchClassifyRequest{
		Documents: []BatchDocumentInput{
			{
				DocumentID: fmt.Sprintf("test-batch-%d", time.Now().Unix()),
				Text:       "This is a test motion to suppress evidence in a criminal case. The defendant argues that the search violated Fourth Amendment rights.",
			},
			{
				DocumentID: fmt.Sprintf("test-batch-%d-2", time.Now().Unix()),
				Text:       "This is a civil complaint regarding breach of contract. The plaintiff seeks damages for non-performance of contractual obligations.",
			},
		},
		Options: map[string]interface{}{
			"update_index": true, // This should trigger indexing
		},
	}

	// Submit the batch job
	log.Printf("📤 Submitting batch job...")
	jobID, err := submitBatchJob(batchReq)
	if err != nil {
		log.Fatalf("❌ Failed to submit batch job: %v", err)
	}

	log.Printf("✅ Batch job submitted: %s", jobID)

	// Monitor progress
	log.Printf("👀 Monitoring progress...")
	for i := 0; i < 30; i++ { // Check for up to 30 seconds
		time.Sleep(2 * time.Second)
		
		job, err := getBatchJobStatus(jobID)
		if err != nil {
			log.Printf("❌ Failed to get job status: %v", err)
			continue
		}

		log.Printf("📊 Progress: %.1f%% (%d/%d processed, %d indexed, %d index errors)", 
			job.Progress.PercentComplete, 
			job.Progress.ProcessedCount, 
			job.Progress.TotalDocuments,
			job.Progress.IndexedCount,
			job.Progress.IndexErrorCount)

		if job.Status == "completed" || job.Status == "failed" {
			log.Printf("🎯 Job %s: %s", job.Status, jobID)
			
			// Get full results
			results, err := getBatchJobResults(jobID)
			if err != nil {
				log.Printf("❌ Failed to get results: %v", err)
			} else {
				log.Printf("📋 Final Results:")
				for i, result := range results {
					log.Printf("  [%d] %s:", i+1, result.DocumentID)
					log.Printf("      Status: %s", result.Status)
					log.Printf("      Indexed: %v", result.Indexed)
					if result.IndexID != "" {
						log.Printf("      Index ID: %s", result.IndexID)
					}
					if result.IndexError != "" {
						log.Printf("      Index Error: %s", result.IndexError)
					}
					if result.Error != "" {
						log.Printf("      Error: %s", result.Error)
					}
				}
			}
			break
		}
	}

	// Search for our test documents
	log.Printf("🔍 Searching for indexed documents...")
	time.Sleep(3 * time.Second) // Give OpenSearch time to index

	searchResp, err := searchDocuments("test motion suppress")
	if err != nil {
		log.Printf("❌ Search failed: %v", err)
	} else {
		log.Printf("🔍 Search found %d documents", len(searchResp.Data["documents"].([]interface{})))
		
		// Check if our test documents appear
		docs := searchResp.Data["documents"].([]interface{})
		for i, doc := range docs {
			docMap := doc.(map[string]interface{})
			docData := docMap["document"].(map[string]interface{})
			docID := docData["id"].(string)
			
			if i < 3 { // Show first 3 results
				log.Printf("  [%d] %s", i+1, docID)
			}
			
			// Check if it's one of our test documents
			for _, testDoc := range batchReq.Documents {
				if docID == testDoc.DocumentID {
					log.Printf("  ⭐ Found our test document: %s", docID)
				}
			}
		}
	}

	log.Printf("🎉 Test completed!")
}

func submitBatchJob(req BatchClassifyRequest) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:6000/api/v1/batch/classify", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var batchResp BatchResponse
	if err := json.Unmarshal(body, &batchResp); err != nil {
		return "", err
	}

	if !batchResp.Success {
		return "", fmt.Errorf("batch submission failed: %s", batchResp.Message)
	}

	jobID := batchResp.Data["job_id"].(string)
	return jobID, nil
}

func getBatchJobStatus(jobID string) (*BatchJob, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:6000/api/v1/batch/%s/status", jobID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var batchResp struct {
		Success bool     `json:"success"`
		Data    BatchJob `json:"data"`
	}
	
	if err := json.Unmarshal(body, &batchResp); err != nil {
		return nil, err
	}

	return &batchResp.Data, nil
}

func getBatchJobResults(jobID string) ([]BatchResult, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:6000/api/v1/batch/%s/results", jobID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var batchResp struct {
		Success bool `json:"success"`
		Data    struct {
			Results []BatchResult `json:"results"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &batchResp); err != nil {
		return nil, err
	}

	return batchResp.Data.Results, nil
}

func searchDocuments(query string) (*BatchResponse, error) {
	searchReq := map[string]interface{}{
		"query": query,
		"size":  10,
	}

	jsonData, err := json.Marshal(searchReq)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:6000/api/v1/search", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResp BatchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, err
	}

	return &searchResp, nil
}