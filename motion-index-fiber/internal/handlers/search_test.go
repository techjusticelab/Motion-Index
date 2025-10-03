package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"motion-index-fiber/pkg/search/models"
)

func TestSearchHandler_SearchDocuments(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		queryParams    map[string]string
		mockResponse   *models.SearchResult
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful search with body",
			requestBody: models.SearchRequest{
				Query: "test query",
				Size:  10,
			},
			mockResponse: &models.SearchResult{
				TotalHits: 5,
				Documents: []*models.SearchDocument{
					{
						ID:    "doc1",
						Score: 0.95,
						Document: map[string]interface{}{
							"title": "Test Document",
						},
					},
				},
			},
			expectedStatus: 200,
		},
		{
			name:        "successful search with query params",
			requestBody: map[string]interface{}{},
			queryParams: map[string]string{
				"q":    "test query",
				"size": "5",
				"from": "10",
			},
			mockResponse: &models.SearchResult{
				TotalHits: 15,
				Documents: []*models.SearchDocument{},
			},
			expectedStatus: 200,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockSearchService)
			handler := NewSearchHandler(mockService)
			app := fiber.New()
			app.Post("/search", handler.SearchDocuments)

			if tt.mockResponse != nil {
				mockService.On("SearchDocuments", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*models.SearchRequest")).Return(tt.mockResponse, tt.mockError)
			}

			// Prepare request
			var reqBody []byte
			var err error
			if tt.requestBody != nil {
				reqBody, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/search", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Add query parameters
			if tt.queryParams != nil {
				q := req.URL.Query()
				for k, v := range tt.queryParams {
					q.Add(k, v)
				}
				req.URL.RawQuery = q.Encode()
			}

			// Execute
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Assert
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.mockResponse != nil {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestSearchHandler_GetLegalTags(t *testing.T) {
	mockService := new(MockSearchService)
	handler := NewSearchHandler(mockService)
	app := fiber.New()
	app.Get("/legal-tags", handler.GetLegalTags)

	expectedTags := []*models.TagCount{
		{Tag: "Contract Law", Count: 150},
		{Tag: "Criminal Law", Count: 89},
	}

	mockService.On("GetLegalTags", mock.AnythingOfType("*context.timerCtx")).Return(expectedTags, nil)

	req := httptest.NewRequest("GET", "/legal-tags", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestSearchHandler_GetDocumentStats(t *testing.T) {
	mockService := new(MockSearchService)
	handler := NewSearchHandler(mockService)
	app := fiber.New()
	app.Get("/document-stats", handler.GetDocumentStats)

	expectedStats := &models.DocumentStats{
		TotalDocuments: 1000,
		IndexSize:      "50MB",
		LastUpdated:    time.Now(),
		TypeCounts: []*models.TypeCount{
			{Type: "Motion", Count: 500},
			{Type: "Order", Count: 300},
		},
	}

	mockService.On("GetDocumentStats", mock.AnythingOfType("*context.timerCtx")).Return(expectedStats, nil)

	req := httptest.NewRequest("GET", "/document-stats", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestSearchHandler_GetDocument(t *testing.T) {
	tests := []struct {
		name           string
		docID          string
		mockDocument   *models.Document
		mockError      error
		expectedStatus int
	}{
		{
			name:  "successful document retrieval",
			docID: "doc123",
			mockDocument: &models.Document{
				ID:       "doc123",
				FileName: "test.pdf",
				Text:     "Document content",
			},
			expectedStatus: 200,
		},
		{
			name:           "document not found",
			docID:          "nonexistent",
			mockDocument:   nil,
			mockError:      assert.AnError,
			expectedStatus: 500,
		},
		{
			name:           "missing document ID",
			docID:          "",
			expectedStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSearchService)
			handler := NewSearchHandler(mockService)
			app := fiber.New()
			app.Get("/documents/:id", handler.GetDocument)

			if tt.docID != "" && tt.expectedStatus != 404 {
				mockService.On("GetDocument", mock.AnythingOfType("*context.timerCtx"), tt.docID).Return(tt.mockDocument, tt.mockError)
			}

			req := httptest.NewRequest("GET", "/documents/"+tt.docID, nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus != 404 {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestValidateSearchRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *models.SearchRequest
		want *models.SearchRequest
	}{
		{
			name: "valid request unchanged",
			req: &models.SearchRequest{
				Query:     "test",
				Size:      20,
				From:      10,
				SortOrder: "desc",
			},
			want: &models.SearchRequest{
				Query:     "test",
				Size:      20,
				From:      10,
				SortOrder: "desc",
			},
		},
		{
			name: "size too large gets capped",
			req: &models.SearchRequest{
				Size: 2000,
			},
			want: &models.SearchRequest{
				Size: models.MaxSearchSize,
			},
		},
		{
			name: "invalid sort order gets defaulted",
			req: &models.SearchRequest{
				SortOrder: "invalid",
			},
			want: &models.SearchRequest{
				Size:      models.DefaultSearchSize,
				SortOrder: "desc",
			},
		},
		{
			name: "negative from gets reset",
			req: &models.SearchRequest{
				From: -5,
			},
			want: &models.SearchRequest{
				Size: models.DefaultSearchSize,
				From: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSearchRequest(tt.req)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Size, tt.req.Size)
			assert.Equal(t, tt.want.From, tt.req.From)
			assert.Equal(t, tt.want.SortOrder, tt.req.SortOrder)
		})
	}
}
