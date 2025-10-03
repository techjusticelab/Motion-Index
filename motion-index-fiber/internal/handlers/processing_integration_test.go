package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"motion-index-fiber/internal/config"
	"motion-index-fiber/internal/models"
	"motion-index-fiber/pkg/processing/classifier"
	"motion-index-fiber/pkg/processing/extractor"
	"motion-index-fiber/pkg/processing/pipeline"
	"motion-index-fiber/pkg/search"
	"motion-index-fiber/pkg/storage"
)

func TestProcessingHandler_Integration(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Environment: "local",
		Server: config.ServerConfig{
			Port: "6000",
		},
		Processing: config.ProcessingConfig{
			MaxWorkers:     2,
			BatchSize:      10,
			ProcessTimeout: 30 * time.Second,
		},
		OpenAI: config.OpenAIConfig{
			APIKey: "", // Empty key will trigger mock service
			Model:  "gpt-4",
		},
	}

	// Initialize services
	extractorService := extractor.NewService()
	classifierService := &mockClassifierService{}
	storageService := storage.NewMockService()
	searchService := search.NewMockService()

	// Initialize processing pipeline
	pipelineConfig := &pipeline.Config{
		MaxWorkers:     cfg.Processing.MaxWorkers,
		QueueSize:      cfg.Processing.BatchSize,
		ProcessTimeout: cfg.Processing.ProcessTimeout,
		RetryAttempts:  3,
		RetryDelay:     1 * time.Second,
		EnableMetrics:  true,
	}

	processingPipeline, err := pipeline.NewPipeline(
		extractorService,
		classifierService,
		searchService,
		storageService,
		pipelineConfig,
	)
	assert.NoError(t, err)

	// Pipeline doesn't need explicit start/stop for this test

	// Create handler
	handler := NewProcessingHandler(processingPipeline, storageService, searchService)

	// Create Fiber app
	app := fiber.New()
	app.Post("/process", handler.ProcessDocument)

	t.Run("ProcessDocument_WithPipeline", func(t *testing.T) {
		// Create test text content for easier processing
		testContent := "This is a test legal document for motion processing. It contains sample text that can be classified and processed by the pipeline."

		// Create multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add file (using .txt extension for easier processing)
		fileWriter, err := writer.CreateFormFile("file", "test.txt")
		assert.NoError(t, err)
		_, err = fileWriter.Write([]byte(testContent))
		assert.NoError(t, err)

		// Add form fields
		writer.WriteField("category", "motion")
		writer.WriteField("case_name", "Test Case")
		writer.WriteField("case_number", "TC-001")

		writer.Close()

		// Create request
		req := httptest.NewRequest("POST", "/process", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Execute request
		resp, err := app.Test(req, 30000) // 30 second timeout
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Read response body
		responseBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		// Debug: print response if not successful
		if resp.StatusCode != 200 {
			t.Logf("Response status: %d, body: %s", resp.StatusCode, string(responseBody))
		}

		// Parse response
		var apiResponse models.APIResponse
		err = json.Unmarshal(responseBody, &apiResponse)
		assert.NoError(t, err)

		// Check response
		assert.Equal(t, 200, resp.StatusCode)
		assert.True(t, apiResponse.Success)

		// Check if data contains expected processing response
		data, ok := apiResponse.Data.(map[string]interface{})
		assert.True(t, ok)

		// Basic validation of response structure
		assert.Contains(t, data, "document_id")
		assert.Contains(t, data, "file_name")
		assert.Equal(t, "test.txt", data["file_name"])
		assert.Contains(t, data, "status")
	})

	t.Run("GetPipelineStatus", func(t *testing.T) {
		// Test pipeline status
		status := processingPipeline.GetStatus()
		assert.NotNil(t, status)
		assert.True(t, processingPipeline.IsHealthy())
	})
}

func TestMockClassifierService(t *testing.T) {
	mock := &mockClassifierService{}

	t.Run("ClassifyDocument", func(t *testing.T) {
		result, err := mock.ClassifyDocument(context.Background(), "This is a legal motion for summary judgment", &classifier.DocumentMetadata{
			FileName: "motion.pdf",
			FileType: "application/pdf",
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, classifier.DocumentTypeMotion, result.DocumentType)
		assert.Equal(t, classifier.LegalCategoryCivil, result.LegalCategory)
		assert.Greater(t, result.Confidence, 0.0)
	})

	t.Run("IsHealthy", func(t *testing.T) {
		assert.True(t, mock.IsHealthy())
	})

	t.Run("GetAvailableCategories", func(t *testing.T) {
		categories := mock.GetAvailableCategories()
		assert.NotEmpty(t, categories)
		assert.Contains(t, categories, classifier.LegalCategoryCivil)
	})
}

// parseJSONResponse parses JSON response body into the provided interface
func parseJSONResponse(resp *http.Response, v interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
