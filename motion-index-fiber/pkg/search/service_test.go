package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"motion-index-fiber/pkg/search/models"
)

// MockSearchClient implements SearchClient interface for testing
type MockSearchClient struct {
	mock.Mock
}

func (m *MockSearchClient) Search(ctx context.Context, req *models.SearchRequest) (*models.SearchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SearchResponse), args.Error(1)
}

func (m *MockSearchClient) Index(ctx context.Context, doc *models.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockSearchClient) Update(ctx context.Context, id string, doc *models.Document) error {
	args := m.Called(ctx, id, doc)
	return args.Error(0)
}

func (m *MockSearchClient) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSearchClient) GetDocument(ctx context.Context, id string) (*models.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Document), args.Error(1)
}

func (m *MockSearchClient) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSearchClient) GetStats(ctx context.Context) (*models.IndexStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IndexStats), args.Error(1)
}

func TestNewService(t *testing.T) {
	client := &MockSearchClient{}
	service := NewService(client)

	assert.NotNil(t, service)
	assert.Equal(t, client, service.client)
}

func TestService_Search(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.SearchRequest
		mockResponse  *models.SearchResponse
		mockError     error
		expectedError bool
	}{
		{
			name: "successful search",
			request: &models.SearchRequest{
				Query: "test query",
				Limit: 10,
			},
			mockResponse: &models.SearchResponse{
				Documents: []*models.Document{
					{
						ID:       "doc1",
						Title:    "Test Document",
						Content:  "Test content",
						Category: "motion",
					},
				},
				Total: 1,
			},
			expectedError: false,
		},
		{
			name: "empty query",
			request: &models.SearchRequest{
				Query: "",
				Limit: 10,
			},
			expectedError: true,
		},
		{
			name: "search client error",
			request: &models.SearchRequest{
				Query: "test query",
				Limit: 10,
			},
			mockError:     assert.AnError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockSearchClient{}
			service := NewService(client)

			if tt.mockResponse != nil || tt.mockError != nil {
				client.On("Search", mock.Anything, tt.request).Return(tt.mockResponse, tt.mockError)
			}

			ctx := context.Background()
			result, err := service.Search(ctx, tt.request)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockResponse, result)
			}

			client.AssertExpectations(t)
		})
	}
}

func TestService_IndexDocument(t *testing.T) {
	tests := []struct {
		name          string
		document      *models.Document
		mockError     error
		expectedError bool
	}{
		{
			name: "successful indexing",
			document: &models.Document{
				ID:       "doc1",
				Title:    "Test Document",
				Content:  "Test content",
				Category: "motion",
			},
			expectedError: false,
		},
		{
			name:          "nil document",
			document:      nil,
			expectedError: true,
		},
		{
			name: "document without ID",
			document: &models.Document{
				Title:    "Test Document",
				Content:  "Test content",
				Category: "motion",
			},
			expectedError: true,
		},
		{
			name: "client error",
			document: &models.Document{
				ID:       "doc1",
				Title:    "Test Document",
				Content:  "Test content",
				Category: "motion",
			},
			mockError:     assert.AnError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockSearchClient{}
			service := NewService(client)

			if tt.document != nil && tt.document.ID != "" {
				client.On("Index", mock.Anything, tt.document).Return(tt.mockError)
			}

			ctx := context.Background()
			err := service.IndexDocument(ctx, tt.document)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.document != nil && tt.document.ID != "" {
				client.AssertExpectations(t)
			}
		})
	}
}

func TestService_UpdateDocument(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		document      *models.Document
		mockError     error
		expectedError bool
	}{
		{
			name: "successful update",
			id:   "doc1",
			document: &models.Document{
				ID:       "doc1",
				Title:    "Updated Document",
				Content:  "Updated content",
				Category: "motion",
			},
			expectedError: false,
		},
		{
			name:          "empty ID",
			id:            "",
			document:      &models.Document{},
			expectedError: true,
		},
		{
			name:          "nil document",
			id:            "doc1",
			document:      nil,
			expectedError: true,
		},
		{
			name: "client error",
			id:   "doc1",
			document: &models.Document{
				ID:       "doc1",
				Title:    "Updated Document",
				Content:  "Updated content",
				Category: "motion",
			},
			mockError:     assert.AnError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockSearchClient{}
			service := NewService(client)

			if tt.id != "" && tt.document != nil {
				client.On("Update", mock.Anything, tt.id, tt.document).Return(tt.mockError)
			}

			ctx := context.Background()
			err := service.UpdateDocument(ctx, tt.id, tt.document)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.id != "" && tt.document != nil {
				client.AssertExpectations(t)
			}
		})
	}
}

func TestService_DeleteDocument(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful deletion",
			id:            "doc1",
			expectedError: false,
		},
		{
			name:          "empty ID",
			id:            "",
			expectedError: true,
		},
		{
			name:          "client error",
			id:            "doc1",
			mockError:     assert.AnError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockSearchClient{}
			service := NewService(client)

			if tt.id != "" {
				client.On("Delete", mock.Anything, tt.id).Return(tt.mockError)
			}

			ctx := context.Background()
			err := service.DeleteDocument(ctx, tt.id)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.id != "" {
				client.AssertExpectations(t)
			}
		})
	}
}

func TestService_GetDocument(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockDocument  *models.Document
		mockError     error
		expectedError bool
	}{
		{
			name: "successful get",
			id:   "doc1",
			mockDocument: &models.Document{
				ID:       "doc1",
				Title:    "Test Document",
				Content:  "Test content",
				Category: "motion",
			},
			expectedError: false,
		},
		{
			name:          "empty ID",
			id:            "",
			expectedError: true,
		},
		{
			name:          "document not found",
			id:            "nonexistent",
			mockError:     assert.AnError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockSearchClient{}
			service := NewService(client)

			if tt.id != "" {
				client.On("GetDocument", mock.Anything, tt.id).Return(tt.mockDocument, tt.mockError)
			}

			ctx := context.Background()
			result, err := service.GetDocument(ctx, tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockDocument, result)
			}

			if tt.id != "" {
				client.AssertExpectations(t)
			}
		})
	}
}

func TestService_IsHealthy(t *testing.T) {
	tests := []struct {
		name           string
		clientHealthy  bool
		expectedHealth bool
	}{
		{
			name:           "healthy client",
			clientHealthy:  true,
			expectedHealth: true,
		},
		{
			name:           "unhealthy client",
			clientHealthy:  false,
			expectedHealth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockSearchClient{}
			service := NewService(client)

			client.On("IsHealthy").Return(tt.clientHealthy)

			result := service.IsHealthy()
			assert.Equal(t, tt.expectedHealth, result)

			client.AssertExpectations(t)
		})
	}
}

func TestService_GetStats(t *testing.T) {
	tests := []struct {
		name          string
		mockStats     *models.IndexStats
		mockError     error
		expectedError bool
	}{
		{
			name: "successful stats retrieval",
			mockStats: &models.IndexStats{
				DocumentCount: 100,
				IndexSize:     "50MB",
			},
			expectedError: false,
		},
		{
			name:          "client error",
			mockError:     assert.AnError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockSearchClient{}
			service := NewService(client)

			client.On("GetStats", mock.Anything).Return(tt.mockStats, tt.mockError)

			ctx := context.Background()
			result, err := service.GetStats(ctx)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockStats, result)
			}

			client.AssertExpectations(t)
		})
	}
}
