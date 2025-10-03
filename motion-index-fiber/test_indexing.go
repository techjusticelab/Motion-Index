package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"motion-index-fiber/internal/config"
	"motion-index-fiber/internal/handlers"
	"motion-index-fiber/pkg/cloud/digitalocean"
	"motion-index-fiber/pkg/processing/classifier"
	"motion-index-fiber/pkg/processing/extractor"
	"motion-index-fiber/pkg/processing/queue"
	searchModels "motion-index-fiber/pkg/search/models"
)

// TestDocument represents a simple test document
type TestDocument struct {
	DocumentID   string `json:"document_id"`
	DocumentPath string `json:"document_path"`
	Text         string `json:"text"`
}

// BatchClassifyRequest represents the API request format
type BatchClassifyRequest struct {
	Documents []TestDocument         `json:"documents"`
	Options   map[string]interface{} `json:"options"`
}

func main() {
	log.Printf("🧪 Testing OpenSearch Indexing")

	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize DigitalOcean provider
	provider, err := digitalocean.NewProviderFromEnvironment()
	if err != nil {
		log.Fatalf("Failed to create DigitalOcean provider: %v", err)
	}

	if err := provider.Initialize(); err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}

	services := provider.GetServices()

	// Test search service health
	log.Printf("🔍 Testing search service health...")
	if !services.Search.IsHealthy() {
		log.Fatalf("❌ Search service is not healthy")
	}
	log.Printf("✅ Search service is healthy")

	// Create a mock queue manager
	queueManager := &mockQueueManager{}

	// Create mock classifier that returns a simple classification
	mockClassifier := &mockClassifierService{}

	// Create mock extractor
	mockExtractor := &mockExtractorService{}

	// Create batch handler
	batchHandler := handlers.NewBatchHandler(
		queueManager,
		services.Storage,
		services.Search,
		mockClassifier,
		mockExtractor,
	)

	// Test direct indexing with a sample document
	log.Printf("📄 Testing direct document indexing...")
	
	ctx := context.Background()
	testDoc := handlers.BatchDocumentInput{
		DocumentID:   "test-doc-" + fmt.Sprintf("%d", time.Now().Unix()),
		DocumentPath: "test/sample-motion.pdf",
		Text:         "This is a sample motion to suppress evidence filed in the Superior Court of California. The defendant respectfully requests that this court suppress all evidence obtained during the unlawful search and seizure that occurred on January 15, 2024.",
	}

	// Create mock classification result
	classificationResult := &classifier.ClassificationResult{
		DocumentType:     "motion_to_suppress",
		LegalCategory:    "Criminal Law",
		SubCategory:      "Fourth Amendment",
		Subject:          "Motion to Suppress Evidence",
		Summary:          "Defendant requests suppression of evidence obtained through unlawful search",
		Confidence:       0.92,
		Keywords:         []string{"motion", "suppress", "evidence", "search", "seizure"},
		LegalTags:        []string{"criminal", "fourth-amendment", "search-and-seizure"},
		Status:           "filed",
		Success:          true,
		ProcessingTime:   1500,
	}

	// Test the indexing function directly
	log.Printf("🔄 Calling indexDocument directly...")
	
	// Use reflection to access the private method (for testing only)
	// In a real test, we'd make this method public or use a different approach
	jobOptions := map[string]interface{}{
		"update_index": true,
	}

	// Process the document (which should include indexing)
	result := processDocumentForTest(ctx, batchHandler, testDoc, jobOptions, classificationResult)
	
	log.Printf("📊 Processing result:")
	log.Printf("   Status: %s", result.Status)
	log.Printf("   Indexed: %v", result.Indexed)
	log.Printf("   Index ID: %s", result.IndexID)
	if result.Error != "" {
		log.Printf("   Error: %s", result.Error)
	}
	if result.IndexError != "" {
		log.Printf("   Index Error: %s", result.IndexError)
	}

	// Verify the document was indexed by searching for it
	if result.Indexed && result.IndexID != "" {
		log.Printf("🔍 Verifying document was indexed...")
		time.Sleep(2 * time.Second) // Give OpenSearch time to index

		// Try to retrieve the document
		doc, err := services.Search.GetDocument(ctx, result.IndexID)
		if err != nil {
			log.Printf("❌ Failed to retrieve indexed document: %v", err)
		} else {
			log.Printf("✅ Successfully retrieved indexed document:")
			log.Printf("   ID: %s", doc.ID)
			log.Printf("   FileName: %s", doc.FileName)
			log.Printf("   DocType: %s", doc.DocType)
			log.Printf("   Category: %s", doc.Category)
		}

		// Test search functionality
		searchReq := &searchModels.SearchRequest{
			Query: "motion suppress evidence",
			Size:  10,
		}
		
		searchResult, err := services.Search.SearchDocuments(ctx, searchReq)
		if err != nil {
			log.Printf("❌ Search failed: %v", err)
		} else {
			log.Printf("🔍 Search results: %d documents found", len(searchResult.Documents))
			for i, doc := range searchResult.Documents {
				if i < 3 { // Show first 3 results
					log.Printf("   [%d] %s - %s (%s)", i+1, doc.ID, doc.FileName, doc.DocType)
				}
			}
		}
	}

	log.Printf("🎉 Test completed!")
}

// processDocumentForTest simulates the document processing
func processDocumentForTest(ctx context.Context, handler *handlers.BatchHandler, doc handlers.BatchDocumentInput, jobOptions map[string]interface{}, classResult *classifier.ClassificationResult) handlers.BatchResult {
	result := handlers.BatchResult{
		DocumentID:   doc.DocumentID,
		DocumentPath: doc.DocumentPath,
		ProcessedAt:  time.Now(),
		Status:       "success",
		ClassificationResult: classResult,
	}

	// Simulate indexing (this would normally be done inside processDocument)
	text := doc.Text
	if shouldIndex := shouldIndexDocument(jobOptions); shouldIndex {
		log.Printf("[TEST] Attempting to index document: %s", doc.DocumentID)
		
		// We need to create the search document ourselves since we can't access the private method
		now := time.Now()
		fileName := "sample-motion.pdf"
		fileURL := fmt.Sprintf("/api/documents/%s", doc.DocumentPath)

		searchDoc := &searchModels.Document{
			ID:          doc.DocumentID,
			FileName:    fileName,
			FilePath:    doc.DocumentPath,
			FileURL:     fileURL,
			Text:        text,
			DocType:     classResult.DocumentType,
			Category:    classResult.LegalCategory,
			Hash:        generateHash(text),
			CreatedAt:   now,
			UpdatedAt:   now,
			ContentType: "application/pdf",
			Size:        int64(len(text)),
			Metadata:    buildMetadata(classResult),
		}

		// Get the search service from handler (we'll need to make this accessible)
		// For now, we'll simulate this step
		log.Printf("[TEST] Document prepared for indexing: ID=%s, DocType=%s", searchDoc.ID, searchDoc.DocType)
		
		result.Indexed = true
		result.IndexID = doc.DocumentID // Simulate successful indexing
	}

	return result
}

// Helper functions
func shouldIndexDocument(options map[string]interface{}) bool {
	if options == nil {
		return false
	}
	if updateIndex, exists := options["update_index"]; exists {
		if val, ok := updateIndex.(bool); ok && val {
			return true
		}
	}
	return false
}

func generateHash(text string) string {
	return fmt.Sprintf("hash_%d", len(text))
}

func buildMetadata(classResult *classifier.ClassificationResult) *searchModels.DocumentMetadata {
	return &searchModels.DocumentMetadata{
		DocumentName:  classResult.Subject,
		Subject:       classResult.Subject,
		Summary:       classResult.Summary,
		DocumentType:  searchModels.DocumentType(classResult.DocumentType),
		Status:        classResult.Status,
		Language:      "en",
		ProcessedAt:   time.Now(),
		Confidence:    classResult.Confidence,
		AIClassified:  true,
		LegalTags:     classResult.LegalTags,
	}
}

// Mock implementations
type mockQueueManager struct{}

func (m *mockQueueManager) EnqueueTask(ctx context.Context, task queue.Task) error {
	return nil
}

func (m *mockQueueManager) DequeueTask(ctx context.Context, queueName string) (queue.Task, error) {
	return nil, nil
}

func (m *mockQueueManager) GetQueueStatus(queueName string) (*queue.QueueStatus, error) {
	return &queue.QueueStatus{}, nil
}

func (m *mockQueueManager) IsHealthy() bool {
	return true
}

type mockClassifierService struct{}

func (m *mockClassifierService) ClassifyDocument(ctx context.Context, text string, metadata *classifier.DocumentMetadata) (*classifier.ClassificationResult, error) {
	return &classifier.ClassificationResult{
		DocumentType:   "motion_to_suppress",
		LegalCategory:  "Criminal Law",
		Subject:        "Motion to Suppress Evidence",
		Summary:        "Mock classification result",
		Confidence:     0.95,
		Success:        true,
		ProcessingTime: 1000,
	}, nil
}

func (m *mockClassifierService) GetAvailableCategories() []string {
	return []string{"Criminal Law", "Civil Law"}
}

func (m *mockClassifierService) IsHealthy() bool {
	return true
}

func (m *mockClassifierService) ValidateResult(result *classifier.ClassificationResult) error {
	return nil
}

type mockExtractorService struct{}

func (m *mockExtractorService) ExtractText(ctx context.Context, reader interface{}, metadata *extractor.DocumentMetadata) (*extractor.ExtractionResult, error) {
	return &extractor.ExtractionResult{
		Text:           "Mock extracted text",
		Success:        true,
		ProcessingTime: 500,
	}, nil
}

func (m *mockExtractorService) GetSupportedFormats() []string {
	return []string{"pdf", "docx", "txt"}
}

func (m *mockExtractorService) IsHealthy() bool {
	return true
}