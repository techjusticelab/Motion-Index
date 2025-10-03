package processing

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"sync/atomic"

	"motion-index-fiber/pkg/processing/extractor"
	"motion-index-fiber/pkg/processing/queue"
	"motion-index-fiber/pkg/models"
)

// createExtractionProcessor creates a processor for document text extraction
func (c *DocumentCoordinator) createExtractionProcessor() queue.ProcessorFunc {
	return func(ctx context.Context, job interface{}) error {
		processingJob, ok := job.(*ProcessingJob)
		if !ok {
			return fmt.Errorf("invalid job type for extraction processor")
		}
		
		// Download document content
		content, err := c.downloadDocument(ctx, processingJob.DocumentKey)
		if err != nil {
			return fmt.Errorf("failed to download document: %w", err)
		}
		
		// Extract text using enhanced extraction service
		metadata := &extractor.DocumentMetadata{
			Filename: processingJob.DocumentKey,
			Size:     int64(len(content)),
		}
		
		result, err := c.extractor.ExtractText(ctx, content, metadata)
		if err != nil {
			return fmt.Errorf("text extraction failed: %w", err)
		}
		
		// Check if we got meaningful text
		if result.Text == "" || len(strings.TrimSpace(result.Text)) < 10 {
			log.Printf("⚠️ No meaningful text extracted from %s, skipping", processingJob.DocumentKey)
			atomic.AddInt64(c.skipped, 1)
			return nil
		}
		
		// Update job with extracted text
		processingJob.ExtractedText = result.Text
		processingJob.Metadata = result.Metadata
		
		// Queue for classification
		classificationQueue, err := c.queueManager.GetQueue("classification")
		if err != nil {
			return fmt.Errorf("failed to get classification queue: %w", err)
		}
		
		return classificationQueue.Enqueue(ctx, processingJob)
	}
}

// createClassificationProcessor creates a processor for AI document classification
func (c *DocumentCoordinator) createClassificationProcessor() queue.ProcessorFunc {
	return func(ctx context.Context, job interface{}) error {
		processingJob, ok := job.(*ProcessingJob)
		if !ok {
			return fmt.Errorf("invalid job type for classification processor")
		}
		
		// Prepare classification metadata
		classifierMetadata := &classifier.DocumentMetadata{
			Filename: processingJob.DocumentKey,
			Text:     processingJob.ExtractedText,
		}
		
		// Classify document
		result, err := c.classifier.ClassifyDocument(ctx, classifierMetadata)
		if err != nil {
			return fmt.Errorf("classification failed: %w", err)
		}
		
		// Update job with classification results
		processingJob.Classification = result.Category
		if processingJob.Metadata == nil {
			processingJob.Metadata = make(map[string]interface{})
		}
		processingJob.Metadata["classification"] = result.Category
		processingJob.Metadata["confidence"] = result.Confidence
		processingJob.Metadata["reasoning"] = result.Reasoning
		
		// Queue for indexing
		indexingQueue, err := c.queueManager.GetQueue("indexing")
		if err != nil {
			return fmt.Errorf("failed to get indexing queue: %w", err)
		}
		
		return indexingQueue.Enqueue(ctx, processingJob)
	}
}

// createIndexingProcessor creates a processor for search index updates
func (c *DocumentCoordinator) createIndexingProcessor() queue.ProcessorFunc {
	return func(ctx context.Context, job interface{}) error {
		processingJob, ok := job.(*ProcessingJob)
		if !ok {
			return fmt.Errorf("invalid job type for indexing processor")
		}
		
		// Create search document
		doc := &models.Document{
			ID:       processingJob.DocumentKey,
			Title:    extractTitle(processingJob.DocumentKey),
			Content:  processingJob.ExtractedText,
			Category: processingJob.Classification,
			Metadata: processingJob.Metadata,
		}
		
		// Index document (assuming search service has Index method)
		// This would need to be adapted based on your actual search interface
		if searchService, ok := c.search.(interface{ IndexDocument(context.Context, *models.Document) error }); ok {
			if err := searchService.IndexDocument(ctx, doc); err != nil {
				return fmt.Errorf("indexing failed: %w", err)
			}
		} else {
			log.Printf("⚠️ Search service does not support indexing, skipping %s", processingJob.DocumentKey)
		}
		
		// Mark as successfully processed
		atomic.AddInt64(c.processed, 1)
		
		log.Printf("✅ Successfully processed: %s (%s)", 
			processingJob.DocumentKey, processingJob.Classification)
		
		return nil
	}
}

// downloadDocument downloads document content from storage
func (c *DocumentCoordinator) downloadDocument(ctx context.Context, key string) (io.Reader, error) {
	reader, err := c.storage.GetObject(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from storage: %w", err)
	}
	
	return reader, nil
}

// extractTitle extracts a title from the document key/filename
func extractTitle(key string) string {
	// Remove path prefix and file extension
	filename := key
	if lastSlash := strings.LastIndex(filename, "/"); lastSlash != -1 {
		filename = filename[lastSlash+1:]
	}
	if lastDot := strings.LastIndex(filename, "."); lastDot != -1 {
		filename = filename[:lastDot]
	}
	
	// Replace underscores and hyphens with spaces
	filename = strings.ReplaceAll(filename, "_", " ")
	filename = strings.ReplaceAll(filename, "-", " ")
	
	// Capitalize first letter of each word
	words := strings.Fields(filename)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	
	return strings.Join(words, " ")
}

// Helper function to safely convert interface{} to string
func toString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}

// Helper function to safely convert interface{} to float64
func toFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0.0
	}
}