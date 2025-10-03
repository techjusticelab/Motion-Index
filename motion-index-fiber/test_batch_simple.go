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

func main() {
	log.Printf("🧪 Simple Batch Test - No AI Classification")

	// Create a simple test without AI classification to isolate indexing issue
	batchReq := BatchClassifyRequest{
		Documents: []BatchDocumentInput{
			{
				DocumentID: fmt.Sprintf("simple-test-%d", time.Now().Unix()),
				Text:       "This is a simple test document for indexing only.",
			},
		},
		Options: map[string]interface{}{
			"update_index":    true,  // Enable indexing
			"skip_ai":         true,  // Skip AI classification to avoid quota issues
			"force_index":     true,  // Force indexing even without classification
		},
	}

	// Submit the batch job
	log.Printf("📤 Submitting simple batch job...")
	jobID, err := submitBatchJob(batchReq)
	if err != nil {
		log.Fatalf("❌ Failed to submit batch job: %v", err)
	}

	log.Printf("✅ Batch job submitted: %s", jobID)

	// Monitor progress for a shorter time
	log.Printf("👀 Monitoring progress...")
	for i := 0; i < 15; i++ { // Check for up to 30 seconds
		time.Sleep(2 * time.Second)
		
		job, err := getBatchJobStatus(jobID)
		if err != nil {
			log.Printf("❌ Failed to get job status: %v", err)
			continue
		}

		log.Printf("📊 Progress: %.1f%% (%d/%d processed, %d indexed, %d index errors)", 
			job["progress"].(map[string]interface{})["percent_complete"].(float64), 
			int(job["progress"].(map[string]interface{})["processed_count"].(float64)), 
			int(job["progress"].(map[string]interface{})["total_documents"].(float64)),
			int(job["progress"].(map[string]interface{})["indexed_count"].(float64)),
			int(job["progress"].(map[string]interface{})["index_error_count"].(float64)))

		status := job["status"].(string)
		if status == "completed" || status == "failed" {
			log.Printf("🎯 Job %s: %s", status, jobID)
			
			// Get full results
			results, err := getBatchJobResults(jobID)
			if err != nil {
				log.Printf("❌ Failed to get results: %v", err)
			} else {
				log.Printf("📋 Final Results:")
				for i, result := range results {
					log.Printf("  [%d] %s:", i+1, result["document_id"].(string))
					log.Printf("      Status: %s", result["status"].(string))
					if indexed, ok := result["indexed"]; ok {
						log.Printf("      Indexed: %v", indexed.(bool))
					}
					if indexID, ok := result["index_id"]; ok && indexID != nil {
						log.Printf("      Index ID: %s", indexID.(string))
					}
					if indexError, ok := result["index_error"]; ok && indexError != nil && indexError.(string) != "" {
						log.Printf("      Index Error: %s", indexError.(string))
					}
					if resultError, ok := result["error"]; ok && resultError != nil && resultError.(string) != "" {
						log.Printf("      Error: %s", resultError.(string))
					}
				}
			}
			break
		}
	}

	log.Printf("🎉 Simple test completed!")
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

func getBatchJobStatus(jobID string) (map[string]interface{}, error) {
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
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &batchResp); err != nil {
		return nil, err
	}

	return batchResp.Data, nil
}

func getBatchJobResults(jobID string) ([]map[string]interface{}, error) {
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
			Results []map[string]interface{} `json:"results"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &batchResp); err != nil {
		return nil, err
	}

	return batchResp.Data.Results, nil
}