package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"motion-index-fiber/internal/models"
	"motion-index-fiber/pkg/search"
	searchModels "motion-index-fiber/pkg/search/models"
	"motion-index-fiber/pkg/storage"
)

// MockStorage is a mock implementation of the storage interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Upload(ctx context.Context, path string, content io.Reader, metadata *storage.UploadMetadata) (*storage.UploadResult, error) {
	args := m.Called(ctx, path, content, metadata)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*storage.UploadResult), args.Error(1)
}

func (m *MockStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorage) Delete(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *MockStorage) GetURL(path string) string {
	args := m.Called(path)
	return args.String(0)
}

func (m *MockStorage) GetSignedURL(path string, expiration time.Duration) (string, error) {
	args := m.Called(path, expiration)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) Exists(ctx context.Context, path string) (bool, error) {
	args := m.Called(ctx, path)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) List(ctx context.Context, prefix string) ([]*storage.StorageObject, error) {
	args := m.Called(ctx, prefix)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*storage.StorageObject), args.Error(1)
}

func (m *MockStorage) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockStorage) GetMetrics() map[string]interface{} {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]interface{})
}

// MockSearchService is a mock implementation of the search service interface
type MockSearchService struct {
	mock.Mock
}

// SearchService methods
func (m *MockSearchService) SearchDocuments(ctx context.Context, req *searchModels.SearchRequest) (*searchModels.SearchResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*searchModels.SearchResult), args.Error(1)
}

func (m *MockSearchService) IndexDocument(ctx context.Context, doc *searchModels.Document) (string, error) {
	args := m.Called(ctx, doc)
	return args.String(0), args.Error(1)
}

func (m *MockSearchService) BulkIndexDocuments(ctx context.Context, docs []*searchModels.Document) (*searchModels.BulkResult, error) {
	args := m.Called(ctx, docs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*searchModels.BulkResult), args.Error(1)
}

func (m *MockSearchService) UpdateDocumentMetadata(ctx context.Context, docID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, docID, metadata)
	return args.Error(0)
}

func (m *MockSearchService) DeleteDocument(ctx context.Context, docID string) error {
	args := m.Called(ctx, docID)
	return args.Error(0)
}

func (m *MockSearchService) GetDocument(ctx context.Context, docID string) (*searchModels.Document, error) {
	args := m.Called(ctx, docID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*searchModels.Document), args.Error(1)
}

func (m *MockSearchService) DocumentExists(ctx context.Context, docID string) (bool, error) {
	args := m.Called(ctx, docID)
	return args.Bool(0), args.Error(1)
}

// AggregationService methods
func (m *MockSearchService) GetLegalTags(ctx context.Context) ([]*searchModels.TagCount, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*searchModels.TagCount), args.Error(1)
}

func (m *MockSearchService) GetDocumentTypes(ctx context.Context) ([]*searchModels.TypeCount, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*searchModels.TypeCount), args.Error(1)
}

func (m *MockSearchService) GetMetadataFieldValues(ctx context.Context, field string, prefix string, size int) ([]*searchModels.FieldValue, error) {
	args := m.Called(ctx, field, prefix, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*searchModels.FieldValue), args.Error(1)
}

func (m *MockSearchService) GetDocumentStats(ctx context.Context) (*searchModels.DocumentStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*searchModels.DocumentStats), args.Error(1)
}

func (m *MockSearchService) GetAllFieldOptions(ctx context.Context) (*searchModels.FieldOptions, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*searchModels.FieldOptions), args.Error(1)
}

// HealthChecker methods
func (m *MockSearchService) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSearchService) Health(ctx context.Context) (*search.HealthStatus, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*search.HealthStatus), args.Error(1)
}

func TestNewHealthHandler(t *testing.T) {
	mockStorage := &MockStorage{}
	mockSearchSvc := &MockSearchService{}

	handler := NewHealthHandler(mockStorage, mockSearchSvc)

	assert.NotNil(t, handler)
	assert.Equal(t, mockStorage, handler.storage)
	assert.Equal(t, mockSearchSvc, handler.searchSvc)
}

func TestHealthHandler_HealthCheck(t *testing.T) {
	mockStorage := &MockStorage{}
	mockSearchSvc := &MockSearchService{}
	handler := NewHealthHandler(mockStorage, mockSearchSvc)

	app := fiber.New()
	app.Get("/health", handler.HealthCheck)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Service is healthy", response.Message)

	// Check the health response data
	healthData := response.Data.(map[string]interface{})
	assert.Equal(t, "healthy", healthData["status"])
	assert.Equal(t, "1.0.0", healthData["version"])
	assert.Equal(t, "motion-index-fiber", healthData["service"])
	assert.NotNil(t, healthData["timestamp"])
}

func TestHealthHandler_DetailedStatus_Healthy(t *testing.T) {
	mockStorage := &MockStorage{}
	mockSearchSvc := &MockSearchService{}

	// Mock healthy responses
	mockStorage.On("IsHealthy").Return(true)
	mockSearchSvc.On("IsHealthy").Return(true)

	handler := NewHealthHandler(mockStorage, mockSearchSvc)

	app := fiber.New()
	app.Get("/status", handler.DetailedStatus)

	req := httptest.NewRequest("GET", "/status", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check the system status data
	statusData := response.Data.(map[string]interface{})
	assert.Equal(t, "healthy", statusData["status"])
	assert.Equal(t, "motion-index-fiber", statusData["service"])
	assert.NotNil(t, statusData["system"])
	assert.NotNil(t, statusData["storage"])
	assert.NotNil(t, statusData["indexer"])

	mockStorage.AssertExpectations(t)
	mockSearchSvc.AssertExpectations(t)
}

func TestHealthHandler_DetailedStatus_Degraded(t *testing.T) {
	mockStorage := &MockStorage{}
	mockSearchSvc := &MockSearchService{}

	// Mock unhealthy storage
	mockStorage.On("IsHealthy").Return(false)
	mockSearchSvc.On("IsHealthy").Return(true)

	handler := NewHealthHandler(mockStorage, mockSearchSvc)

	app := fiber.New()
	app.Get("/status", handler.DetailedStatus)

	req := httptest.NewRequest("GET", "/status", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check the system status data
	statusData := response.Data.(map[string]interface{})
	assert.Equal(t, "degraded", statusData["status"])

	mockStorage.AssertExpectations(t)
	mockSearchSvc.AssertExpectations(t)
}

func TestHealthHandler_ReadinessCheck_Ready(t *testing.T) {
	mockStorage := &MockStorage{}
	mockSearchSvc := &MockSearchService{}

	// Mock healthy responses
	mockStorage.On("IsHealthy").Return(true)
	mockSearchSvc.On("IsHealthy").Return(true)

	handler := NewHealthHandler(mockStorage, mockSearchSvc)

	app := fiber.New()
	app.Get("/ready", handler.ReadinessCheck)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check the readiness response data
	readinessData := response.Data.(map[string]interface{})
	assert.Equal(t, true, readinessData["ready"])
	assert.NotNil(t, readinessData["checks"])

	checks := readinessData["checks"].(map[string]interface{})
	assert.Equal(t, true, checks["storage"])
	assert.Equal(t, true, checks["search"])

	mockStorage.AssertExpectations(t)
	mockSearchSvc.AssertExpectations(t)
}

func TestHealthHandler_ReadinessCheck_NotReady(t *testing.T) {
	mockStorage := &MockStorage{}
	mockSearchSvc := &MockSearchService{}

	// Mock unhealthy search service
	mockStorage.On("IsHealthy").Return(true)
	mockSearchSvc.On("IsHealthy").Return(false)

	handler := NewHealthHandler(mockStorage, mockSearchSvc)

	app := fiber.New()
	app.Get("/ready", handler.ReadinessCheck)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Check the readiness response data
	readinessData := response.Data.(map[string]interface{})
	assert.Equal(t, false, readinessData["ready"])

	checks := readinessData["checks"].(map[string]interface{})
	assert.Equal(t, true, checks["storage"])
	assert.Equal(t, false, checks["search"])

	mockStorage.AssertExpectations(t)
	mockSearchSvc.AssertExpectations(t)
}

func TestHealthHandler_LivenessCheck(t *testing.T) {
	mockStorage := &MockStorage{}
	mockSearchSvc := &MockSearchService{}
	handler := NewHealthHandler(mockStorage, mockSearchSvc)

	app := fiber.New()
	app.Get("/live", handler.LivenessCheck)

	req := httptest.NewRequest("GET", "/live", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Service is alive", response.Message)

	// Check the liveness response data
	livenessData := response.Data.(map[string]interface{})
	assert.Equal(t, true, livenessData["alive"])
	assert.NotNil(t, livenessData["timestamp"])
	assert.NotNil(t, livenessData["pid"])
}

func TestHealthHandler_Metrics(t *testing.T) {
	mockStorage := &MockStorage{}
	mockSearchSvc := &MockSearchService{}

	// Mock metrics responses
	storageMetrics := map[string]interface{}{
		"uploads_count": 100,
		"storage_used":  "1GB",
	}

	healthStatus := &search.HealthStatus{
		Status:        "green",
		ClusterName:   "test-cluster",
		NumberOfNodes: 3,
		ActiveShards:  10,
		IndexExists:   true,
		IndexHealth:   "green",
	}

	mockStorage.On("GetMetrics").Return(storageMetrics)
	mockSearchSvc.On("Health", mock.AnythingOfType("*context.timerCtx")).Return(healthStatus, nil)

	handler := NewHealthHandler(mockStorage, mockSearchSvc)

	app := fiber.New()
	app.Get("/metrics", handler.Metrics)

	req := httptest.NewRequest("GET", "/metrics", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Application metrics", response.Message)

	// Check the metrics response data
	metricsData := response.Data.(map[string]interface{})
	assert.NotNil(t, metricsData["timestamp"])
	assert.NotNil(t, metricsData["memory"])
	assert.NotNil(t, metricsData["goroutines"])
	assert.NotNil(t, metricsData["gc"])
	assert.NotNil(t, metricsData["storage"])
	assert.NotNil(t, metricsData["indexer"])

	mockStorage.AssertExpectations(t)
	mockSearchSvc.AssertExpectations(t)
}

func TestHealthHandler_GetSystemInfo(t *testing.T) {
	systemInfo := getSystemInfo()

	assert.NotNil(t, systemInfo)
	assert.NotEmpty(t, systemInfo.OS)
	assert.NotEmpty(t, systemInfo.Architecture)
	assert.NotEmpty(t, systemInfo.GoVersion)
	assert.Greater(t, systemInfo.NumCPU, 0)
	assert.GreaterOrEqual(t, systemInfo.Goroutines, 1)
	assert.NotNil(t, systemInfo.Memory)
	assert.Greater(t, systemInfo.Memory.Alloc, uint64(0))
}

func TestHealthHandler_GetUptime(t *testing.T) {
	// Reset start time for testing
	originalStartTime := startTime
	startTime = time.Now().Add(-5 * time.Second)
	defer func() {
		startTime = originalStartTime
	}()

	uptime := getUptime()
	assert.GreaterOrEqual(t, uptime, 4*time.Second)
	assert.Less(t, uptime, 6*time.Second)
}

func TestHealthHandler_GetMemoryStats(t *testing.T) {
	memStats := getMemoryStats()

	assert.NotNil(t, memStats)
	assert.Greater(t, memStats.Alloc, uint64(0))
	assert.Greater(t, memStats.TotalAlloc, uint64(0))
	assert.Greater(t, memStats.Sys, uint64(0))
	assert.GreaterOrEqual(t, memStats.NumGC, uint32(0))
}

func TestHealthHandler_GetGCStats(t *testing.T) {
	gcStats := getGCStats()

	assert.NotNil(t, gcStats)
	assert.GreaterOrEqual(t, gcStats.NumGC, uint32(0))
	assert.GreaterOrEqual(t, gcStats.PauseTotal, time.Duration(0))
	assert.Greater(t, gcStats.NextGC, uint64(0))
}

func TestHealthHandler_GetStorageStatus_Healthy(t *testing.T) {
	mockStorage := &MockStorage{}
	mockStorage.On("IsHealthy").Return(true)

	handler := NewHealthHandler(mockStorage, nil)
	status := handler.getStorageStatus()

	assert.Equal(t, "storage", status.Name)
	assert.Equal(t, "healthy", status.Status)
	assert.Empty(t, status.Error)
	assert.True(t, status.LastError.IsZero())

	mockStorage.AssertExpectations(t)
}

func TestHealthHandler_GetStorageStatus_Unhealthy(t *testing.T) {
	mockStorage := &MockStorage{}
	mockStorage.On("IsHealthy").Return(false)

	handler := NewHealthHandler(mockStorage, nil)
	status := handler.getStorageStatus()

	assert.Equal(t, "storage", status.Name)
	assert.Equal(t, "unhealthy", status.Status)
	assert.Equal(t, "storage service is not healthy", status.Error)
	assert.False(t, status.LastError.IsZero())

	mockStorage.AssertExpectations(t)
}

func TestHealthHandler_GetStorageStatus_Nil(t *testing.T) {
	handler := NewHealthHandler(nil, nil)
	status := handler.getStorageStatus()

	assert.Equal(t, "storage", status.Name)
	assert.Equal(t, "unhealthy", status.Status)
	assert.Equal(t, "storage not initialized", status.Error)
	assert.False(t, status.LastError.IsZero())
}

func TestHealthHandler_GetSearchStatus_Healthy(t *testing.T) {
	mockSearchSvc := &MockSearchService{}
	mockSearchSvc.On("IsHealthy").Return(true)

	handler := NewHealthHandler(nil, mockSearchSvc)
	status := handler.getSearchStatus()

	assert.Equal(t, "search", status.Name)
	assert.Equal(t, "healthy", status.Status)
	assert.Empty(t, status.Error)
	assert.True(t, status.LastError.IsZero())

	mockSearchSvc.AssertExpectations(t)
}

func TestHealthHandler_GetSearchStatus_Unhealthy(t *testing.T) {
	mockSearchSvc := &MockSearchService{}
	mockSearchSvc.On("IsHealthy").Return(false)

	handler := NewHealthHandler(nil, mockSearchSvc)
	status := handler.getSearchStatus()

	assert.Equal(t, "search", status.Name)
	assert.Equal(t, "unhealthy", status.Status)
	assert.Equal(t, "search service is not healthy", status.Error)
	assert.False(t, status.LastError.IsZero())

	mockSearchSvc.AssertExpectations(t)
}

func TestHealthHandler_GetStorageMetrics(t *testing.T) {
	mockStorage := &MockStorage{}
	expectedMetrics := map[string]interface{}{
		"uploads": 100,
		"size":    "1GB",
	}
	mockStorage.On("GetMetrics").Return(expectedMetrics)

	handler := NewHealthHandler(mockStorage, nil)
	metrics := handler.getStorageMetrics()

	assert.Equal(t, expectedMetrics, metrics)
	mockStorage.AssertExpectations(t)
}

func TestHealthHandler_GetStorageMetrics_Nil(t *testing.T) {
	handler := NewHealthHandler(nil, nil)
	metrics := handler.getStorageMetrics()

	assert.NotNil(t, metrics)
	assert.Empty(t, metrics)
}

func TestHealthHandler_GetSearchMetrics(t *testing.T) {
	mockSearchSvc := &MockSearchService{}
	healthStatus := &search.HealthStatus{
		Status:        "green",
		ClusterName:   "test-cluster",
		NumberOfNodes: 3,
		ActiveShards:  10,
		IndexExists:   true,
		IndexHealth:   "green",
	}
	mockSearchSvc.On("Health", mock.AnythingOfType("*context.timerCtx")).Return(healthStatus, nil)

	handler := NewHealthHandler(nil, mockSearchSvc)
	metrics := handler.getSearchMetrics()

	assert.Equal(t, "test-cluster", metrics["cluster_name"])
	assert.Equal(t, 3, metrics["number_of_nodes"])
	assert.Equal(t, 10, metrics["active_shards"])
	assert.Equal(t, true, metrics["index_exists"])
	assert.Equal(t, "green", metrics["index_health"])
	mockSearchSvc.AssertExpectations(t)
}

func TestHealthHandler_GetSearchMetrics_Nil(t *testing.T) {
	handler := NewHealthHandler(nil, nil)
	metrics := handler.getSearchMetrics()

	assert.NotNil(t, metrics)
	assert.Empty(t, metrics)
}
