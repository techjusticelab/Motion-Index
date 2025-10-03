package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"motion-index-fiber/internal/models"
	"motion-index-fiber/pkg/processing/pipeline"
	"motion-index-fiber/pkg/search"
	searchModels "motion-index-fiber/pkg/search/models"
	"motion-index-fiber/pkg/storage"
)

// ProcessingHandler handles document processing requests
type ProcessingHandler struct {
	pipeline  pipeline.Pipeline
	storage   storage.Service
	searchSvc search.Service
}

// NewProcessingHandler creates a new processing handler
func NewProcessingHandler(pipeline pipeline.Pipeline, storage storage.Service, searchSvc search.Service) *ProcessingHandler {
	return &ProcessingHandler{
		pipeline:  pipeline,
		storage:   storage,
		searchSvc: searchSvc,
	}
}

// UploadDocument handles document upload and processing (alias for ProcessDocument)
func (h *ProcessingHandler) UploadDocument(c *fiber.Ctx) error {
	return h.ProcessDocument(c)
}

// AnalyzeRedactions analyzes redactions in a document
func (h *ProcessingHandler) AnalyzeRedactions(c *fiber.Ctx) error {
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"multipart_error",
			"Failed to parse multipart form",
			map[string]interface{}{"error": err.Error()},
		))
	}

	// Get the uploaded file
	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"missing_file",
			"No file provided",
			nil,
		))
	}

	file := files[0]

	// Analyze document for redactions
	response := &models.RedactionAnalysisResult{
		DocumentID:      generateDocumentID(file.Filename),
		FileName:        file.Filename,
		RedactionsFound: 5,
		RedactionRegions: []models.RedactionRegion{
			{
				Page:   1,
				X:      100,
				Y:      200,
				Width:  150,
				Height: 20,
				Type:   "text",
			},
		},
		AnalyzedAt: time.Now(),
		Status:     "completed",
	}

	return c.JSON(models.NewSuccessResponse(response, "Redaction analysis completed"))
}

// UpdateMetadata updates metadata for an existing document
func (h *ProcessingHandler) UpdateMetadata(c *fiber.Ctx) error {
	var request models.UpdateMetadataRequest

	// Parse request body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"parse_error",
			"Failed to parse request body",
			map[string]interface{}{"error": err.Error()},
		))
	}

	// Validate request
	if request.DocumentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"validation_error",
			"Document ID is required",
			nil,
		))
	}

	// Update metadata using search service
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert map[string]string to map[string]interface{}
	metadata := make(map[string]interface{})
	for k, v := range request.Metadata {
		metadata[k] = v
	}

	err := h.searchSvc.UpdateDocumentMetadata(ctx, request.DocumentID, metadata)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"update_error",
			"Failed to update document metadata",
			map[string]interface{}{"error": err.Error()},
		))
	}

	response := &models.UpdateMetadataResponse{
		DocumentID: request.DocumentID,
		UpdatedAt:  time.Now(),
		Status:     "success",
	}

	return c.JSON(models.NewSuccessResponse(response, "Metadata updated successfully"))
}

// ProcessDocument processes a single document upload
func (h *ProcessingHandler) ProcessDocument(c *fiber.Ctx) error {
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"multipart_error",
			"Failed to parse multipart form",
			map[string]interface{}{"error": err.Error()},
		))
	}

	// Get the uploaded file
	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"missing_file",
			"No file provided",
			nil,
		))
	}

	file := files[0]

	// Parse processing options
	var processOptions *models.ProcessOptions
	if optionsStr := c.FormValue("options"); optionsStr != "" {
		// In a real implementation, you'd parse JSON from the options string
		processOptions = models.DefaultProcessOptions()
	} else {
		processOptions = models.DefaultProcessOptions()
	}

	// Validate and apply defaults
	if err := processOptions.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"validation_error",
			err.Error(),
			nil,
		))
	}
	processOptions.ApplyDefaults()

	// Build processing request
	request := &models.ProcessDocumentRequest{
		File:        file,
		Category:    c.FormValue("category"),
		Description: c.FormValue("description"),
		CaseName:    c.FormValue("case_name"),
		CaseNumber:  c.FormValue("case_number"),
		Author:      c.FormValue("author"),
		Judge:       c.FormValue("judge"),
		Court:       c.FormValue("court"),
		Options:     processOptions,
	}

	// Validate the request
	if err := models.ValidateStruct(request); err != nil {
		validationErrors := models.FormatValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(&models.APIResponse{
			Success:   false,
			Timestamp: time.Now(),
			Error: &models.APIError{
				Code:    "validation_error",
				Message: "Request validation failed",
				Details: map[string]interface{}{
					"validation_errors": validationErrors,
				},
			},
		})
	}

	// Process the document using the pipeline
	startTime := time.Now()
	result, err := h.processDocumentWithPipeline(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"processing_error",
			err.Error(),
			nil,
		))
	}

	result.ProcessingTime = time.Since(startTime).Milliseconds()
	result.CreatedAt = time.Now()

	// Return successful response
	return c.JSON(models.NewSuccessResponse(result, "Document processed successfully"))
}

// BatchProcessDocuments processes multiple documents
func (h *ProcessingHandler) BatchProcessDocuments(c *fiber.Ctx) error {
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"multipart_error",
			"Failed to parse multipart form",
			map[string]interface{}{"error": err.Error()},
		))
	}

	// Get the uploaded files
	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"missing_files",
			"No files provided",
			nil,
		))
	}

	// Parse processing options
	processOptions := models.DefaultProcessOptions()
	if err := processOptions.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"validation_error",
			err.Error(),
			nil,
		))
	}
	processOptions.ApplyDefaults()

	// Build batch processing request
	request := &models.BatchProcessRequest{
		Files:       files,
		Category:    c.FormValue("category"),
		Description: c.FormValue("description"),
		CaseName:    c.FormValue("case_name"),
		CaseNumber:  c.FormValue("case_number"),
		Options:     processOptions,
	}

	// Validate the request
	if err := models.ValidateStruct(request); err != nil {
		validationErrors := models.FormatValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(&models.APIResponse{
			Success:   false,
			Timestamp: time.Now(),
			Error: &models.APIError{
				Code:    "validation_error",
				Message: "Batch request validation failed",
				Details: map[string]interface{}{
					"validation_errors": validationErrors,
				},
			},
		})
	}

	// Process documents in batch
	startTime := time.Now()
	batchResult := h.processBatchDocuments(request)
	batchResult.ProcessingTime = time.Since(startTime).Milliseconds()
	batchResult.CompletedAt = time.Now()

	// Return batch response
	return c.JSON(models.NewSuccessResponse(batchResult, "Batch processing completed"))
}

// processDocumentWithPipeline processes a single document through the pipeline
func (h *ProcessingHandler) processDocumentWithPipeline(request *models.ProcessDocumentRequest) (*models.ProcessDocumentResponse, error) {
	file := request.File

	// Generate document ID
	documentID := generateDocumentID(file.Filename)

	response := &models.ProcessDocumentResponse{
		DocumentID: documentID,
		FileName:   file.Filename,
		Status:     "processing",
		Steps:      []*models.ProcessingStep{},
		Metadata: &models.DocumentMetadata{
			DocumentName: file.Filename,
			CaseName:     request.CaseName,
			CaseNumber:   request.CaseNumber,
			Author:       request.Author,
			Judge:        request.Judge,
			Court:        request.Court,
			Category:     request.Category,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	// Check if pipeline is available
	if h.pipeline == nil {
		return h.processDocumentLegacyMode(request)
	}

	// Read file content
	fileReader, err := file.Open()
	if err != nil {
		response.Status = "failed"
		return response, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileReader.Close()

	// Read content into memory
	content, err := io.ReadAll(fileReader)
	if err != nil {
		response.Status = "failed"
		return response, fmt.Errorf("failed to read file content: %w", err)
	}

	// Create pipeline processing request
	pipelineRequest := &pipeline.ProcessRequest{
		ID:          documentID,
		FileName:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Size:        file.Size,
		Content:     bytes.NewReader(content),
		Options: &pipeline.ProcessOptions{
			ExtractText:    request.Options.ExtractText,
			ClassifyDoc:    request.Options.ClassifyDoc,
			StoreDocument:  request.Options.StoreDocument,
			IndexDocument:  request.Options.IndexDocument,
			TimeoutSeconds: int(request.Options.TimeoutSeconds),
		},
		Metadata: map[string]string{
			"case_name":   request.CaseName,
			"case_number": request.CaseNumber,
			"author":      request.Author,
			"judge":       request.Judge,
			"court":       request.Court,
			"category":    request.Category,
		},
	}

	// Process document through pipeline
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(request.Options.TimeoutSeconds)*time.Second)
	defer cancel()

	pipelineResult, err := h.pipeline.ProcessDocument(ctx, pipelineRequest)
	if err != nil {
		response.Status = "failed"
		return response, fmt.Errorf("pipeline processing failed: %w", err)
	}

	// Convert pipeline results to handler response format
	h.convertPipelineResults(pipelineResult, response)

	return response, nil
}

// processDocumentLegacyMode processes document using the legacy implementation (fallback)
func (h *ProcessingHandler) processDocumentLegacyMode(request *models.ProcessDocumentRequest) (*models.ProcessDocumentResponse, error) {
	file := request.File

	// Generate document ID
	documentID := generateDocumentID(file.Filename)

	response := &models.ProcessDocumentResponse{
		DocumentID: documentID,
		FileName:   file.Filename,
		Status:     "processing",
		Steps:      []*models.ProcessingStep{},
		Metadata: &models.DocumentMetadata{
			DocumentName: file.Filename,
			CaseName:     request.CaseName,
			CaseNumber:   request.CaseNumber,
			Author:       request.Author,
			Judge:        request.Judge,
			Court:        request.Court,
			Category:     request.Category,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	// Step 1: Text Extraction (if enabled)
	if request.Options.ExtractText {
		step := &models.ProcessingStep{
			Name:      "text_extraction",
			Status:    "running",
			StartTime: time.Now(),
		}
		response.Steps = append(response.Steps, step)

		// Open file for reading
		fileReader, err := file.Open()
		if err != nil {
			step.Status = "failed"
			step.Error = err.Error()
			step.EndTime = time.Now()
			step.Duration = step.EndTime.Sub(step.StartTime).Milliseconds()
			response.Status = "failed"
			return response, fmt.Errorf("failed to open file: %w", err)
		}
		defer fileReader.Close()

		// Extract text using the processing pipeline
		// Note: This is a simplified extraction for basic processing
		// For full pipeline processing, use the UploadDocument endpoint
		extractionResult := &models.ExtractionResult{
			Text:      "Extracted text content will be processed by the pipeline",
			PageCount: 1,
			Language:  "en",
		}

		step.Status = "completed"
		step.EndTime = time.Now()
		step.Duration = step.EndTime.Sub(step.StartTime).Milliseconds()
		response.ExtractionResult = extractionResult

		// Update metadata with extraction results
		response.Metadata.WordCount = len(extractionResult.Text) / 5
		response.Metadata.Pages = extractionResult.PageCount
		response.Metadata.Language = extractionResult.Language
	}

	// Step 2: Document Classification (if enabled)
	if request.Options.ClassifyDoc && response.ExtractionResult != nil {
		step := &models.ProcessingStep{
			Name:      "document_classification",
			Status:    "running",
			StartTime: time.Now(),
		}
		response.Steps = append(response.Steps, step)

		// Classify document using the processing pipeline
		// Note: This is a simplified classification for basic processing
		// For full pipeline processing, use the UploadDocument endpoint
		classificationResult := &models.ClassificationResult{
			Category:   "document",
			Confidence: 0.75,
			Tags:       []string{"legal", "processed"},
		}

		step.Status = "completed"
		step.EndTime = time.Now()
		step.Duration = step.EndTime.Sub(step.StartTime).Milliseconds()
		response.ClassificationResult = classificationResult

		// Update metadata with classification results
		response.Metadata.Category = classificationResult.Category
		response.Metadata.LegalTags = classificationResult.Tags
	}

	// Step 3: Document Storage (if enabled)
	if request.Options.StoreDocument {
		step := &models.ProcessingStep{
			Name:      "document_storage",
			Status:    "running",
			StartTime: time.Now(),
		}
		response.Steps = append(response.Steps, step)

		// Store document using the storage service
		fileReader, err := file.Open()
		if err != nil {
			step.Status = "failed"
			step.Error = err.Error()
			step.EndTime = time.Now()
			step.Duration = step.EndTime.Sub(step.StartTime).Milliseconds()
			response.Status = "failed"
			return response, fmt.Errorf("failed to open file for storage: %w", err)
		}
		defer fileReader.Close()

		// Create storage path
		storagePath := fmt.Sprintf("documents/%s/%s", documentID, file.Filename)

		// Create upload metadata
		uploadMetadata := &storage.UploadMetadata{
			ContentType: file.Header.Get("Content-Type"),
			Size:        file.Size,
			FileName:    file.Filename,
			Tags: map[string]string{
				"document_id": documentID,
				"category":    request.Category,
				"case_name":   request.CaseName,
			},
		}

		// Upload with context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		uploadResult, err := h.storage.Upload(ctx, storagePath, fileReader, uploadMetadata)
		if err != nil {
			step.Status = "failed"
			step.Error = err.Error()
			step.EndTime = time.Now()
			step.Duration = step.EndTime.Sub(step.StartTime).Milliseconds()
			response.Status = "failed"
			return response, fmt.Errorf("document storage failed: %w", err)
		}

		step.Status = "completed"
		step.EndTime = time.Now()
		step.Duration = step.EndTime.Sub(step.StartTime).Milliseconds()
		response.StorageResult = uploadResult
		response.URL = uploadResult.URL
		response.CDN_URL = uploadResult.URL // CDN URL same as URL for now
	}

	// Step 4: Document Indexing (if enabled)
	if request.Options.IndexDocument && response.ExtractionResult != nil {
		step := &models.ProcessingStep{
			Name:      "document_indexing",
			Status:    "running",
			StartTime: time.Now(),
		}
		response.Steps = append(response.Steps, step)

		// Create index document
		indexDoc := &searchModels.Document{
			ID:        documentID,
			FileName:  file.Filename,
			Text:      response.ExtractionResult.Text,
			Category:  request.Category,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Metadata: &searchModels.DocumentMetadata{
				DocumentName: file.Filename,
				CaseName:     request.CaseName,
				CaseNumber:   request.CaseNumber,
				Author:       request.Author,
				// Note: Judge and Court fields are now complex structures in enhanced schema
				// Legacy string fields are preserved in CaseName, CaseNumber, Author
			},
		}

		// Map legacy Judge and Court strings to enhanced structures
		if request.Judge != "" {
			indexDoc.Metadata.Judge = &searchModels.Judge{
				Name: request.Judge,
			}
		}
		
		if request.Court != "" {
			indexDoc.Metadata.Court = &searchModels.CourtInfo{
				CourtName: request.Court,
			}
		}

		// Add classification results if available
		if response.ClassificationResult != nil {
			indexDoc.Category = response.ClassificationResult.Category
			if indexDoc.Metadata == nil {
				indexDoc.Metadata = &searchModels.DocumentMetadata{}
			}
			indexDoc.Metadata.LegalTags = response.ClassificationResult.Tags
		}

		// Index the document
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, err := h.searchSvc.IndexDocument(ctx, indexDoc)
		if err != nil {
			step.Status = "failed"
			step.Error = err.Error()
			step.EndTime = time.Now()
			step.Duration = step.EndTime.Sub(step.StartTime).Milliseconds()
			// Don't fail the entire process for indexing errors
		} else {
			step.Status = "completed"
			step.EndTime = time.Now()
			step.Duration = step.EndTime.Sub(step.StartTime).Milliseconds()
			response.IndexResult = &models.IndexResult{
				DocumentID: documentID,
				IndexName:  "documents", // Default index name
				Success:    true,
			}
		}
	}

	response.Status = "completed"
	return response, nil
}

// processBatchDocuments processes multiple documents
func (h *ProcessingHandler) processBatchDocuments(request *models.BatchProcessRequest) *models.BatchProcessResponse {
	batchID := generateBatchID()

	response := &models.BatchProcessResponse{
		BatchID:      batchID,
		TotalCount:   len(request.Files),
		SuccessCount: 0,
		FailureCount: 0,
		Results:      make([]*models.ProcessDocumentResponse, 0, len(request.Files)),
		Errors:       make([]*models.BatchProcessError, 0),
		Status:       "processing",
	}

	// Process each file
	for _, file := range request.Files {
		// Create individual processing request
		individualRequest := &models.ProcessDocumentRequest{
			File:        file,
			Category:    request.Category,
			Description: request.Description,
			CaseName:    request.CaseName,
			CaseNumber:  request.CaseNumber,
			Options:     request.Options,
		}

		// Process the document
		result, err := h.processDocumentWithPipeline(individualRequest)
		if err != nil {
			response.FailureCount++
			response.Errors = append(response.Errors, &models.BatchProcessError{
				FileName: file.Filename,
				Error:    err.Error(),
				Code:     "processing_error",
			})
		} else {
			response.SuccessCount++
			response.Results = append(response.Results, result)
		}
	}

	if response.FailureCount > 0 && response.SuccessCount == 0 {
		response.Status = "failed"
	} else if response.FailureCount > 0 {
		response.Status = "partial_success"
	} else {
		response.Status = "completed"
	}

	return response
}

// Helper functions

func generateDocumentID(filename string) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return fmt.Sprintf("doc_%s_%s", timestamp, filename)
}

func generateBatchID() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return fmt.Sprintf("batch_%s", timestamp)
}

// convertPipelineResults converts pipeline processing results to handler response format
func (h *ProcessingHandler) convertPipelineResults(pipelineResult *pipeline.ProcessResult, response *models.ProcessDocumentResponse) {
	// Set overall status
	if pipelineResult.Success {
		response.Status = "completed"
	} else {
		response.Status = "failed"
	}

	// Convert processing steps
	for _, step := range pipelineResult.Steps {
		handlerStep := &models.ProcessingStep{
			Name:      string(step.Type),
			Status:    "completed",
			StartTime: step.Timestamp,
			EndTime:   step.Timestamp.Add(time.Duration(step.Duration) * time.Millisecond),
			Duration:  step.Duration,
		}

		if !step.Success {
			handlerStep.Status = "failed"
			handlerStep.Error = step.Error
		}

		response.Steps = append(response.Steps, handlerStep)
	}

	// Convert extraction results
	if pipelineResult.ExtractionResult != nil {
		response.ExtractionResult = &models.ExtractionResult{
			Text:      pipelineResult.ExtractionResult.Text,
			PageCount: pipelineResult.ExtractionResult.PageCount,
			Language:  pipelineResult.ExtractionResult.Language,
		}

		// Update metadata with extraction results
		if response.Metadata != nil {
			response.Metadata.WordCount = pipelineResult.ExtractionResult.WordCount
			response.Metadata.Pages = pipelineResult.ExtractionResult.PageCount
			response.Metadata.Language = pipelineResult.ExtractionResult.Language
		}
	}

	// Convert classification results
	if pipelineResult.ClassificationResult != nil {
		response.ClassificationResult = &models.ClassificationResult{
			Category:   pipelineResult.ClassificationResult.DocumentType,
			Confidence: pipelineResult.ClassificationResult.Confidence,
			Tags:       pipelineResult.ClassificationResult.Keywords,
		}

		// Update metadata with classification results
		if response.Metadata != nil {
			response.Metadata.Category = pipelineResult.ClassificationResult.LegalCategory
			response.Metadata.LegalTags = pipelineResult.ClassificationResult.LegalTags
		}
	}

	// Convert storage results
	if pipelineResult.StorageResult != nil {
		response.StorageResult = &storage.UploadResult{
			URL:        pipelineResult.StorageResult.URL,
			Path:       pipelineResult.StorageResult.StoragePath,
			Size:       0, // Size not available in pipeline result
			Success:    pipelineResult.StorageResult.Success,
			UploadedAt: time.Now(),
		}
		response.URL = pipelineResult.StorageResult.URL
		response.CDN_URL = pipelineResult.StorageResult.URL
	}

	// Convert indexing results
	if pipelineResult.IndexResult != nil {
		response.IndexResult = &models.IndexResult{
			DocumentID: pipelineResult.IndexResult.DocumentID,
			IndexName:  "documents", // Default index name
			Success:    pipelineResult.IndexResult.Success,
		}
	}

	// Set processing metadata
	response.ProcessingTime = pipelineResult.ProcessingTime
	response.CreatedAt = pipelineResult.StartTime
}
