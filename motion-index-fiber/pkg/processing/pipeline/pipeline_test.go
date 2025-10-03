package pipeline

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"motion-index-fiber/pkg/processing/classifier"
	"motion-index-fiber/pkg/processing/extractor"
	"motion-index-fiber/pkg/search"
	"motion-index-fiber/pkg/models"
	"motion-index-fiber/pkg/storage"
)

// Mock services for testing
type mockExtractorService struct {
	mock.Mock
}

func (m *mockExtractorService) ExtractText(ctx context.Context, reader io.Reader, metadata *extractor.DocumentMetadata) (*extractor.ExtractionResult, error) {
	args := m.Called(ctx, reader, metadata)
	return args.Get(0).(*extractor.ExtractionResult), args.Error(1)
}

func (m *mockExtractorService) GetExtractor(format string) (extractor.Extractor, error) {
	args := m.Called(format)
	return args.Get(0).(extractor.Extractor), args.Error(1)
}

func (m *mockExtractorService) SupportedFormats() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

type mockClassifierService struct {
	mock.Mock
}

func (m *mockClassifierService) ClassifyDocument(ctx context.Context, text string, metadata *classifier.DocumentMetadata) (*classifier.ClassificationResult, error) {
	args := m.Called(ctx, text, metadata)
	return args.Get(0).(*classifier.ClassificationResult), args.Error(1)
}

func (m *mockClassifierService) GetAvailableCategories() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *mockClassifierService) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockClassifierService) ValidateResult(result *classifier.ClassificationResult) error {
	args := m.Called(result)
	return args.Error(0)
}

type mockSearchService struct {
	mock.Mock
}

func (m *mockSearchService) IndexDocument(ctx context.Context, doc *models.Document) (string, error) {
	args := m.Called(ctx, doc)
	return args.String(0), args.Error(1)
}

func (m *mockSearchService) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockSearchService) SearchDocuments(ctx context.Context, req *models.SearchRequest) (*models.SearchResult, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.SearchResult), args.Error(1)
}

func (m *mockSearchService) GetLegalTags(ctx context.Context) ([]*models.TagCount, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.TagCount), args.Error(1)
}

func (m *mockSearchService) GetDocumentTypes(ctx context.Context) ([]*models.TypeCount, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.TypeCount), args.Error(1)
}

func (m *mockSearchService) GetDocumentStats(ctx context.Context) (*models.DocumentStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.DocumentStats), args.Error(1)
}

func (m *mockSearchService) GetAllFieldOptions(ctx context.Context) (*models.FieldOptions, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.FieldOptions), args.Error(1)
}

func (m *mockSearchService) GetMetadataFieldValues(ctx context.Context, field string, prefix string, size int) ([]*models.FieldValue, error) {
	args := m.Called(ctx, field, prefix, size)
	return args.Get(0).([]*models.FieldValue), args.Error(1)
}

func (m *mockSearchService) BulkIndexDocuments(ctx context.Context, docs []*models.Document) (*models.BulkResult, error) {
	args := m.Called(ctx, docs)
	return args.Get(0).(*models.BulkResult), args.Error(1)
}

func (m *mockSearchService) UpdateDocumentMetadata(ctx context.Context, docID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, docID, metadata)
	return args.Error(0)
}

func (m *mockSearchService) DeleteDocument(ctx context.Context, docID string) error {
	args := m.Called(ctx, docID)
	return args.Error(0)
}

func (m *mockSearchService) GetDocument(ctx context.Context, docID string) (*models.Document, error) {
	args := m.Called(ctx, docID)
	return args.Get(0).(*models.Document), args.Error(1)
}

func (m *mockSearchService) DocumentExists(ctx context.Context, docID string) (bool, error) {
	args := m.Called(ctx, docID)
	return args.Bool(0), args.Error(1)
}

func (m *mockSearchService) Health(ctx context.Context) (*search.HealthStatus, error) {
	args := m.Called(ctx)
	return args.Get(0).(*search.HealthStatus), args.Error(1)
}

func TestNewPipeline(t *testing.T) {
	extractorSvc := &mockExtractorService{}
	classifierSvc := &mockClassifierService{}
	searchSvc := &mockSearchService{}
	storageSvc := storage.NewMockService()

	pipeline, err := NewPipeline(extractorSvc, classifierSvc, searchSvc, storageSvc, nil)

	assert.NoError(t, err)
	assert.NotNil(t, pipeline)
}

func TestPipeline_ProcessDocument(t *testing.T) {
	// Setup mocks
	extractorSvc := &mockExtractorService{}
	classifierSvc := &mockClassifierService{}
	searchSvc := &mockSearchService{}
	storageSvc := storage.NewMockService()

	// Setup expectations
	extractorSvc.On("ExtractText", mock.Anything, mock.Anything, mock.Anything).Return(
		&extractor.ExtractionResult{
			Text:      "Sample extracted text",
			WordCount: 3,
			CharCount: 21,
			Success:   true,
		}, nil)

	classifierSvc.On("ClassifyDocument", mock.Anything, mock.Anything, mock.Anything).Return(
		&classifier.ClassificationResult{
			DocumentType:  "Motion",
			LegalCategory: "Civil Law",
			Confidence:    0.8,
			Success:       true,
		}, nil)

	classifierSvc.On("IsHealthy").Return(true)

	searchSvc.On("IndexDocument", mock.Anything, mock.Anything).Return("doc123", nil)
	searchSvc.On("IsHealthy").Return(true)

	// Create pipeline
	pipeline, err := NewPipeline(extractorSvc, classifierSvc, searchSvc, storageSvc, DefaultConfig())
	assert.NoError(t, err)

	// Create process request
	req := &ProcessRequest{
		ID:          "test-doc-1",
		FileName:    "test.pdf",
		ContentType: "application/pdf",
		Size:        1024,
		Content:     strings.NewReader("Sample document content"),
		Options:     DefaultProcessOptions(),
		Timestamp:   time.Now(),
	}

	// Process document
	ctx := context.Background()
	result, err := pipeline.ProcessDocument(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "test-doc-1", result.ID)
	assert.NotEmpty(t, result.Steps)
	assert.NotNil(t, result.ExtractionResult)
	assert.NotNil(t, result.ClassificationResult)
	assert.NotNil(t, result.IndexResult)
	assert.NotNil(t, result.StorageResult)
	assert.GreaterOrEqual(t, result.ProcessingTime, int64(0))

	// Verify all mocks were called
	extractorSvc.AssertExpectations(t)
	classifierSvc.AssertExpectations(t)
	searchSvc.AssertExpectations(t)
}

func TestPipeline_ProcessBatch(t *testing.T) {
	// Setup mocks
	extractorSvc := &mockExtractorService{}
	classifierSvc := &mockClassifierService{}
	searchSvc := &mockSearchService{}
	storageSvc := storage.NewMockService()

	// Setup expectations for multiple documents
	extractorSvc.On("ExtractText", mock.Anything, mock.Anything, mock.Anything).Return(
		&extractor.ExtractionResult{
			Text:      "Sample extracted text",
			WordCount: 3,
			CharCount: 21,
			Success:   true,
		}, nil).Times(2)

	classifierSvc.On("ClassifyDocument", mock.Anything, mock.Anything, mock.Anything).Return(
		&classifier.ClassificationResult{
			DocumentType:  "Motion",
			LegalCategory: "Civil Law",
			Confidence:    0.8,
			Success:       true,
		}, nil).Times(2)

	classifierSvc.On("IsHealthy").Return(true)

	searchSvc.On("IndexDocument", mock.Anything, mock.Anything).Return("doc123", nil).Times(2)
	searchSvc.On("IsHealthy").Return(true)

	// Create pipeline
	pipeline, err := NewPipeline(extractorSvc, classifierSvc, searchSvc, storageSvc, DefaultConfig())
	assert.NoError(t, err)

	// Create batch requests
	requests := []*ProcessRequest{
		{
			ID:          "batch-doc-1",
			FileName:    "doc1.pdf",
			ContentType: "application/pdf",
			Size:        1024,
			Content:     strings.NewReader("Document 1 content"),
			Options:     DefaultProcessOptions(),
		},
		{
			ID:          "batch-doc-2",
			FileName:    "doc2.pdf",
			ContentType: "application/pdf",
			Size:        2048,
			Content:     strings.NewReader("Document 2 content"),
			Options:     DefaultProcessOptions(),
		},
	}

	// Process batch
	ctx := context.Background()
	batchResult, err := pipeline.ProcessBatch(ctx, requests)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, batchResult)
	assert.Equal(t, 2, batchResult.TotalCount)
	assert.Equal(t, 2, batchResult.SuccessCount)
	assert.Equal(t, 0, batchResult.FailureCount)
	assert.Len(t, batchResult.Results, 2)

	// Verify all results are successful
	for _, result := range batchResult.Results {
		assert.True(t, result.Success)
		assert.NotEmpty(t, result.Steps)
	}

	// Verify all mocks were called
	extractorSvc.AssertExpectations(t)
	classifierSvc.AssertExpectations(t)
	searchSvc.AssertExpectations(t)
}

func TestPipeline_GetStatus(t *testing.T) {
	extractorSvc := &mockExtractorService{}
	classifierSvc := &mockClassifierService{}
	searchSvc := &mockSearchService{}
	storageSvc := storage.NewMockService()

	// Setup health check expectations
	classifierSvc.On("IsHealthy").Return(true)
	searchSvc.On("IsHealthy").Return(true)

	pipeline, err := NewPipeline(extractorSvc, classifierSvc, searchSvc, storageSvc, DefaultConfig())
	assert.NoError(t, err)

	status := pipeline.GetStatus()

	assert.NotNil(t, status)
	assert.NotEmpty(t, status.ProcessorStatus)
	assert.GreaterOrEqual(t, status.CompletedJobs, int64(0))
	assert.GreaterOrEqual(t, status.FailedJobs, int64(0))
}

func TestPipeline_IsHealthy(t *testing.T) {
	extractorSvc := &mockExtractorService{}
	classifierSvc := &mockClassifierService{}
	searchSvc := &mockSearchService{}
	storageSvc := storage.NewMockService()

	// Setup health check expectations
	classifierSvc.On("IsHealthy").Return(true)
	searchSvc.On("IsHealthy").Return(true)

	pipeline, err := NewPipeline(extractorSvc, classifierSvc, searchSvc, storageSvc, DefaultConfig())
	assert.NoError(t, err)

	healthy := pipeline.IsHealthy()
	assert.True(t, healthy)
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.Greater(t, config.MaxWorkers, 0)
	assert.Greater(t, config.QueueSize, 0)
	assert.Greater(t, config.ProcessTimeout, time.Duration(0))
	assert.Greater(t, config.RetryAttempts, 0)
	assert.True(t, config.EnableMetrics)
}

func TestDefaultProcessOptions(t *testing.T) {
	options := DefaultProcessOptions()

	assert.NotNil(t, options)
	assert.True(t, options.ExtractText)
	assert.True(t, options.ClassifyDoc)
	assert.True(t, options.IndexDocument)
	assert.True(t, options.StoreDocument)
	assert.Greater(t, options.TimeoutSeconds, 0)
	assert.Greater(t, options.RetryCount, 0)
}
