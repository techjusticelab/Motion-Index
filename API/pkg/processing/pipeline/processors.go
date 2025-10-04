package pipeline

import (
	"context"
	"fmt"
	"strings"
	"time"

	"motion-index-fiber/pkg/processing/classifier"
	"motion-index-fiber/pkg/processing/extractor"
	"motion-index-fiber/pkg/search"
	"motion-index-fiber/pkg/models"
	"motion-index-fiber/pkg/storage"
)

// extractionProcessor handles text extraction
type extractionProcessor struct {
	service extractor.Service
}

// NewExtractionProcessor creates a new extraction processor
func NewExtractionProcessor(service extractor.Service) Processor {
	return &extractionProcessor{
		service: service,
	}
}

// Process executes text extraction
func (p *extractionProcessor) Process(ctx context.Context, req *ProcessRequest) (*ProcessResult, error) {
	if p.service == nil {
		return nil, fmt.Errorf("extraction service not available")
	}

	// Create extraction metadata
	metadata := &extractor.DocumentMetadata{
		FileName: req.FileName,
		MimeType: req.ContentType,
		Size:     req.Size,
	}

	// Extract text
	result, err := p.service.ExtractText(ctx, req.Content, metadata)
	if err != nil {
		return nil, fmt.Errorf("text extraction failed: %w", err)
	}

	return &ProcessResult{
		ID:               req.ID,
		ExtractionResult: result,
	}, nil
}

// GetType returns the processor type
func (p *extractionProcessor) GetType() ProcessorType {
	return ProcessorTypeExtraction
}

// IsHealthy returns true if the processor is healthy
func (p *extractionProcessor) IsHealthy() bool {
	return p.service != nil
}

// classificationProcessor handles document classification
type classificationProcessor struct {
	service classifier.Service
}

// NewClassificationProcessor creates a new classification processor
func NewClassificationProcessor(service classifier.Service) Processor {
	return &classificationProcessor{
		service: service,
	}
}

// Process executes document classification
func (p *classificationProcessor) Process(ctx context.Context, req *ProcessRequest) (*ProcessResult, error) {
	if p.service == nil {
		return nil, fmt.Errorf("classification service not available")
	}

	// Get extracted text from request context or extract it ourselves
	var text string
	var wordCount, pageCount int

	// Check if we have extracted text in the metadata
	if extractedText, exists := req.Metadata["extracted_text"]; exists {
		text = extractedText
		if wc, exists := req.Metadata["word_count"]; exists {
			fmt.Sscanf(wc, "%d", &wordCount)
		}
		if pc, exists := req.Metadata["page_count"]; exists {
			fmt.Sscanf(pc, "%d", &pageCount)
		}
	} else {
		// Extract text first using basic extractor
		extractorSvc := extractor.NewService()
		extractorMetadata := &extractor.DocumentMetadata{
			FileName: req.FileName,
			MimeType: req.ContentType,
			Size:     req.Size,
		}

		extractionResult, err := extractorSvc.ExtractText(ctx, req.Content, extractorMetadata)
		if err != nil {
			return nil, fmt.Errorf("failed to extract text for classification: %w", err)
		}

		text = extractionResult.Text
		wordCount = extractionResult.WordCount
		pageCount = extractionResult.PageCount
	}

	// Create classification metadata
	metadata := &classifier.DocumentMetadata{
		FileName:  req.FileName,
		FileType:  req.ContentType,
		Size:      req.Size,
		WordCount: wordCount,
		PageCount: pageCount,
	}

	// Classify document
	result, err := p.service.ClassifyDocument(ctx, text, metadata)
	if err != nil {
		return nil, fmt.Errorf("document classification failed: %w", err)
	}

	return &ProcessResult{
		ID:                   req.ID,
		ClassificationResult: result,
	}, nil
}

// GetType returns the processor type
func (p *classificationProcessor) GetType() ProcessorType {
	return ProcessorTypeClassification
}

// IsHealthy returns true if the processor is healthy
func (p *classificationProcessor) IsHealthy() bool {
	return p.service != nil && p.service.IsHealthy()
}

// indexingProcessor handles document indexing
type indexingProcessor struct {
	service search.Service
}

// NewIndexingProcessor creates a new indexing processor
func NewIndexingProcessor(service search.Service) Processor {
	return &indexingProcessor{
		service: service,
	}
}

// Process executes document indexing
func (p *indexingProcessor) Process(ctx context.Context, req *ProcessRequest) (*ProcessResult, error) {
	if p.service == nil {
		return nil, fmt.Errorf("search service not available")
	}

	// Extract data from previous processing steps
	extractedText := req.Metadata["extracted_text"]
	if extractedText == "" {
		extractedText = "No text extracted"
	}

	// Create document for indexing with all collected data
	// Only include fields that exist in the actual index schema
	doc := &models.Document{
		ID:          req.ID,
		FileName:    req.FileName,
		FilePath:    req.ID, // Use ID as file path since documents are already stored
		ContentType: req.ContentType,
		Size:        req.Size,
		Text:        extractedText,
		Hash:        fmt.Sprintf("hash_%s", req.ID), // Generate a basic hash
		Metadata:    &models.DocumentMetadata{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Populate metadata from processing results
	doc.Metadata.DocumentName = req.FileName

	// Add extraction metadata (only fields that exist in the actual index schema)
	// Note: word_count, pages, and language fields don't exist in the current index
	// So we'll skip setting them to avoid indexing errors

	// Add classification metadata (map to existing fields)
	if documentType, exists := req.Metadata["document_type"]; exists {
		doc.DocType = documentType // Map to Document.DocType
	} else {
		doc.DocType = "Other" // Default document type
	}
	if legalCategory, exists := req.Metadata["legal_category"]; exists {
		doc.Category = legalCategory // Map to Document.Category
	} else {
		doc.Category = "Civil" // Default legal category
	}
	if subCategory, exists := req.Metadata["sub_category"]; exists {
		doc.Metadata.Subject = subCategory // Use Subject field for sub-category
	}
	// Note: Confidence is not in the existing schema, so we'll skip it for now
	if summary, exists := req.Metadata["summary"]; exists {
		doc.Metadata.Subject = summary // Use Subject field if no sub-category
	}

	// Add storage metadata
	if storagePath, exists := req.Metadata["storage_path"]; exists {
		doc.FilePath = storagePath // Map to Document.FilePath
	}
	if storageURL, exists := req.Metadata["storage_url"]; exists {
		doc.FileURL = storageURL // Map to Document.FileURL
	}

	// Set processing timestamp
	now := time.Now()
	doc.Metadata.Timestamp = &now

	// Index document
	docID, err := p.service.IndexDocument(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("document indexing failed: %w", err)
	}

	return &ProcessResult{
		ID: req.ID,
		IndexResult: &IndexResult{
			DocumentID: docID,
			Success:    true,
		},
		Document: doc,
	}, nil
}

// GetType returns the processor type
func (p *indexingProcessor) GetType() ProcessorType {
	return ProcessorTypeIndexing
}

// IsHealthy returns true if the processor is healthy
func (p *indexingProcessor) IsHealthy() bool {
	return p.service != nil && p.service.IsHealthy()
}

// storageProcessor handles document storage
type storageProcessor struct {
	service storage.Service
}

// NewStorageProcessor creates a new storage processor
func NewStorageProcessor(service storage.Service) Processor {
	return &storageProcessor{
		service: service,
	}
}

// Process executes document storage
func (p *storageProcessor) Process(ctx context.Context, req *ProcessRequest) (*ProcessResult, error) {
	if p.service == nil {
		return nil, fmt.Errorf("storage service not available")
	}

	// Generate storage path
	storagePath := p.generateStoragePath(req.FileName, req.ID)

	// Store document (this is a placeholder implementation)
	// In a real implementation, you would upload the document to storage
	url := fmt.Sprintf("https://storage.example.com/%s", storagePath)

	return &ProcessResult{
		ID: req.ID,
		StorageResult: &StorageResult{
			StoragePath: storagePath,
			URL:         url,
			Success:     true,
		},
	}, nil
}

// generateStoragePath generates a storage path for the document
func (p *storageProcessor) generateStoragePath(fileName, docID string) string {
	// Create a path based on document ID and filename
	// Format: documents/{year}/{month}/{docID}/{filename}
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")

	// Clean filename
	cleanName := strings.ReplaceAll(fileName, " ", "_")
	cleanName = strings.ReplaceAll(cleanName, "/", "_")

	return fmt.Sprintf("documents/%s/%s/%s/%s", year, month, docID, cleanName)
}

// GetType returns the processor type
func (p *storageProcessor) GetType() ProcessorType {
	return ProcessorTypeStorage
}

// IsHealthy returns true if the processor is healthy
func (p *storageProcessor) IsHealthy() bool {
	return p.service != nil
}

// validationProcessor handles document validation (optional)
type validationProcessor struct{}

// NewValidationProcessor creates a new validation processor
func NewValidationProcessor() Processor {
	return &validationProcessor{}
}

// Process executes document validation
func (p *validationProcessor) Process(ctx context.Context, req *ProcessRequest) (*ProcessResult, error) {
	// Validate file size
	if req.Size > 100*1024*1024 { // 100MB limit
		return nil, fmt.Errorf("file size too large: %d bytes", req.Size)
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"application/pdf":  true,
		"application/docx": true,
		"text/plain":       true,
	}

	if !allowedTypes[req.ContentType] {
		return nil, fmt.Errorf("unsupported file type: %s", req.ContentType)
	}

	// Validate filename
	if req.FileName == "" {
		return nil, fmt.Errorf("filename is required")
	}

	return &ProcessResult{
		ID: req.ID,
	}, nil
}

// GetType returns the processor type
func (p *validationProcessor) GetType() ProcessorType {
	return ProcessorTypeValidation
}

// IsHealthy returns true if the processor is healthy
func (p *validationProcessor) IsHealthy() bool {
	return true
}
