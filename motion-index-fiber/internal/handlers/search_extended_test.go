package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"motion-index-fiber/pkg/models"
)

func TestSearchHandler_SearchDocuments_Extended(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		queryParams    map[string]string
		mockResponse   *models.SearchResult
		mockError      error
		expectedStatus int
		expectedData   bool
	}{
		{
			name: "complex search with filters",
			requestBody: models.SearchRequest{
				Query:             "contract dispute",
				Size:              20,
				From:              10,
				DocType:           "Motion",
				LegalTags:         []string{"Contract Law", "Dispute Resolution"},
				CaseNumber:        "CV-2023-001",
				SortBy:            "created_at",
				SortOrder:         "desc",
				FuzzySearch:       true,
				IncludeHighlights: true,
			},
			mockResponse: &models.SearchResult{
				TotalHits: 150,
				MaxScore:  0.95,
				Documents: []*models.SearchDocument{
					{
						ID:    "doc1",
						Score: 0.95,
						Document: map[string]interface{}{
							"title":       "Contract Dispute Motion",
							"doc_type":    "Motion",
							"case_number": "CV-2023-001",
						},
						Highlights: map[string][]string{
							"text": {"<mark>contract</mark> <mark>dispute</mark>"},
						},
					},
				},
				Took:     25,
				TimedOut: false,
			},
			expectedStatus: 200,
			expectedData:   true,
		},
		{
			name: "search with date range",
			requestBody: models.SearchRequest{
				Query: "motion",
				DateRange: &models.DateRange{
					From: timePtr(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
					To:   timePtr(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
				},
			},
			mockResponse: &models.SearchResult{
				TotalHits: 50,
				Documents: []*models.SearchDocument{
					{
						ID:    "doc2",
						Score: 0.85,
						Document: map[string]interface{}{
							"title":      "Motion to Dismiss",
							"created_at": "2023-06-15T14:30:00Z",
						},
					},
				},
			},
			expectedStatus: 200,
			expectedData:   true,
		},
		{
			name: "empty search query",
			requestBody: models.SearchRequest{
				Query: "",
				Size:  10,
			},
			mockResponse: &models.SearchResult{
				TotalHits: 1000,
				Documents: []*models.SearchDocument{},
			},
			expectedStatus: 200,
			expectedData:   true,
		},
		{
			name:           "search service timeout",
			requestBody:    models.SearchRequest{Query: "test"},
			mockError:      errors.New("context deadline exceeded"),
			expectedStatus: 500,
			expectedData:   false,
		},
		{
			name: "search with pagination beyond limits",
			requestBody: models.SearchRequest{
				Query: "test",
				Size:  2000, // Beyond max size
				From:  -10,  // Negative from
			},
			mockResponse: &models.SearchResult{
				TotalHits: 100,
				Documents: []*models.SearchDocument{},
			},
			expectedStatus: 200,
			expectedData:   true,
		},
		{
			name:           "malformed JSON request",
			requestBody:    `{"query": "test", "size": "invalid"}`,
			expectedStatus: 400,
			expectedData:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSearchService)
			handler := NewSearchHandler(mockService)
			app := fiber.New()
			app.Post("/search", handler.SearchDocuments)

			if tt.mockResponse != nil || tt.mockError != nil {
				mockService.On("SearchDocuments", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*models.SearchRequest")).Return(tt.mockResponse, tt.mockError)
			}

			var reqBody []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/search", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedData {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
				assert.NotNil(t, response["data"])
			}

			if tt.mockResponse != nil || tt.mockError != nil {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestSearchHandler_GetLegalTags_Extended(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   []*models.TagCount
		mockError      error
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "successful with multiple tags",
			mockResponse: []*models.TagCount{
				{Tag: "Contract Law", Count: 150},
				{Tag: "Criminal Law", Count: 89},
				{Tag: "Family Law", Count: 45},
				{Tag: "Corporate Law", Count: 123},
				{Tag: "Immigration Law", Count: 67},
			},
			expectedStatus: 200,
			expectedCount:  5,
		},
		{
			name:           "empty tags list",
			mockResponse:   []*models.TagCount{},
			expectedStatus: 200,
			expectedCount:  0,
		},
		{
			name:           "service error",
			mockError:      errors.New("aggregation failed"),
			expectedStatus: 500,
		},
		{
			name:           "service timeout",
			mockError:      errors.New("context deadline exceeded"),
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSearchService)
			handler := NewSearchHandler(mockService)
			app := fiber.New()
			app.Get("/legal-tags", handler.GetLegalTags)

			mockService.On("GetLegalTags", mock.AnythingOfType("*context.timerCtx")).Return(tt.mockResponse, tt.mockError)

			req := httptest.NewRequest("GET", "/legal-tags", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == 200 {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])

				data, ok := response["data"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, data, tt.expectedCount)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSearchHandler_GetDocumentTypes_Extended(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   []*models.TypeCount
		mockError      error
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "successful with multiple types",
			mockResponse: []*models.TypeCount{
				{Type: "Motion", Count: 500},
				{Type: "Order", Count: 300},
				{Type: "Brief", Count: 150},
				{Type: "Complaint", Count: 200},
				{Type: "Answer", Count: 180},
			},
			expectedStatus: 200,
			expectedCount:  5,
		},
		{
			name:           "empty types list",
			mockResponse:   []*models.TypeCount{},
			expectedStatus: 200,
			expectedCount:  0,
		},
		{
			name:           "aggregation service error",
			mockError:      errors.New("failed to retrieve document types"),
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSearchService)
			handler := NewSearchHandler(mockService)
			app := fiber.New()
			app.Get("/document-types", handler.GetDocumentTypes)

			mockService.On("GetDocumentTypes", mock.AnythingOfType("*context.timerCtx")).Return(tt.mockResponse, tt.mockError)

			req := httptest.NewRequest("GET", "/document-types", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == 200 {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])

				data, ok := response["data"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, data, tt.expectedCount)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSearchHandler_GetDocumentStats_Extended(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   *models.DocumentStats
		mockError      error
		expectedStatus int
	}{
		{
			name: "comprehensive stats",
			mockResponse: &models.DocumentStats{
				TotalDocuments: 1500,
				IndexSize:      "75MB",
				LastUpdated:    time.Now(),
				TypeCounts: []*models.TypeCount{
					{Type: "Motion", Count: 600},
					{Type: "Order", Count: 400},
					{Type: "Brief", Count: 300},
					{Type: "Complaint", Count: 200},
				},
				TagCounts: []*models.TagCount{
					{Tag: "Contract Law", Count: 400},
					{Tag: "Criminal Law", Count: 350},
					{Tag: "Family Law", Count: 300},
				},
				FieldStats: map[string]models.FieldStat{
					"court": {
						UniqueValues: 25,
						TotalValues:  1500,
					},
					"judge": {
						UniqueValues: 45,
						TotalValues:  1500,
					},
				},
			},
			expectedStatus: 200,
		},
		{
			name: "minimal stats",
			mockResponse: &models.DocumentStats{
				TotalDocuments: 0,
				IndexSize:      "0B",
				LastUpdated:    time.Now(),
				TypeCounts:     []*models.TypeCount{},
				TagCounts:      []*models.TagCount{},
				FieldStats:     map[string]models.FieldStat{},
			},
			expectedStatus: 200,
		},
		{
			name:           "stats service error",
			mockError:      errors.New("stats aggregation failed"),
			expectedStatus: 500,
		},
		{
			name:           "opensearch unavailable",
			mockError:      errors.New("connection refused"),
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSearchService)
			handler := NewSearchHandler(mockService)
			app := fiber.New()
			app.Get("/document-stats", handler.GetDocumentStats)

			mockService.On("GetDocumentStats", mock.AnythingOfType("*context.timerCtx")).Return(tt.mockResponse, tt.mockError)

			req := httptest.NewRequest("GET", "/document-stats", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == 200 {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
				assert.NotNil(t, response["data"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSearchHandler_GetFieldOptions_Extended(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   *models.FieldOptions
		mockError      error
		expectedStatus int
	}{
		{
			name: "comprehensive field options",
			mockResponse: &models.FieldOptions{
				Courts: []*models.FieldValue{
					{Value: "Superior Court", Count: 500},
					{Value: "District Court", Count: 300},
					{Value: "Appeals Court", Count: 150},
				},
				Judges: []*models.FieldValue{
					{Value: "Judge Smith", Count: 200},
					{Value: "Judge Johnson", Count: 180},
					{Value: "Judge Williams", Count: 150},
				},
				DocTypes: []*models.FieldValue{
					{Value: "Motion", Count: 600},
					{Value: "Order", Count: 400},
				},
				LegalTags: []*models.FieldValue{
					{Value: "Contract Law", Count: 400},
					{Value: "Criminal Law", Count: 350},
				},
				Statuses: []*models.FieldValue{
					{Value: "Active", Count: 800},
					{Value: "Closed", Count: 700},
				},
				Authors: []*models.FieldValue{
					{Value: "Attorney A", Count: 100},
					{Value: "Attorney B", Count: 90},
				},
			},
			expectedStatus: 200,
		},
		{
			name: "empty field options",
			mockResponse: &models.FieldOptions{
				Courts:    []*models.FieldValue{},
				Judges:    []*models.FieldValue{},
				DocTypes:  []*models.FieldValue{},
				LegalTags: []*models.FieldValue{},
				Statuses:  []*models.FieldValue{},
				Authors:   []*models.FieldValue{},
			},
			expectedStatus: 200,
		},
		{
			name:           "field options service error",
			mockError:      errors.New("failed to retrieve field options"),
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSearchService)
			handler := NewSearchHandler(mockService)
			app := fiber.New()
			app.Get("/field-options", handler.GetFieldOptions)

			mockService.On("GetAllFieldOptions", mock.AnythingOfType("*context.timerCtx")).Return(tt.mockResponse, tt.mockError)

			req := httptest.NewRequest("GET", "/field-options", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == 200 {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
				assert.NotNil(t, response["data"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSearchHandler_GetMetadataFieldValues_Extended(t *testing.T) {
	tests := []struct {
		name           string
		field          string
		queryParams    map[string]string
		mockResponse   []*models.FieldValue
		mockError      error
		expectedStatus int
		expectedField  string
		expectedPrefix string
		expectedSize   int
	}{
		{
			name:  "court field with prefix",
			field: "metadata.court",
			queryParams: map[string]string{
				"prefix": "Superior",
				"size":   "10",
			},
			mockResponse: []*models.FieldValue{
				{Value: "Superior Court of California", Count: 150},
				{Value: "Superior Court of New York", Count: 120},
			},
			expectedStatus: 200,
			expectedField:  "metadata.court",
			expectedPrefix: "Superior",
			expectedSize:   10,
		},
		{
			name:  "judge field without prefix",
			field: "metadata.judge",
			mockResponse: []*models.FieldValue{
				{Value: "Judge Smith", Count: 200},
				{Value: "Judge Johnson", Count: 180},
				{Value: "Judge Williams", Count: 150},
			},
			expectedStatus: 200,
			expectedField:  "metadata.judge",
			expectedPrefix: "",
			expectedSize:   50, // default
		},
		{
			name:  "large size parameter",
			field: "metadata.legal_tags",
			queryParams: map[string]string{
				"size": "200",
			},
			mockResponse: []*models.FieldValue{
				{Value: "Contract Law", Count: 400},
				{Value: "Criminal Law", Count: 350},
			},
			expectedStatus: 200,
			expectedField:  "metadata.legal_tags",
			expectedSize:   200,
		},
		{
			name:           "missing field parameter",
			field:          "",
			expectedStatus: 404, // Fiber route doesn't match empty param
		},
		{
			name:           "field values service error",
			field:          "metadata.court",
			mockError:      errors.New("failed to retrieve field values"),
			expectedStatus: 500,
			expectedField:  "metadata.court",
			expectedSize:   50,
		},
		{
			name:  "invalid size parameter gets defaulted",
			field: "metadata.status",
			queryParams: map[string]string{
				"size": "invalid",
			},
			mockResponse: []*models.FieldValue{
				{Value: "Active", Count: 800},
			},
			expectedStatus: 200,
			expectedField:  "metadata.status",
			expectedSize:   50, // default due to invalid size
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSearchService)
			handler := NewSearchHandler(mockService)
			app := fiber.New()
			app.Get("/metadata-fields/:field", handler.GetMetadataFieldValues)

			if tt.expectedStatus != 404 {
				mockService.On("GetMetadataFieldValues",
					mock.AnythingOfType("*context.timerCtx"),
					tt.expectedField,
					tt.expectedPrefix,
					tt.expectedSize).Return(tt.mockResponse, tt.mockError)
			}

			url := "/metadata-fields/" + tt.field
			if len(tt.queryParams) > 0 {
				url += "?"
				for k, v := range tt.queryParams {
					url += k + "=" + v + "&"
				}
				url = url[:len(url)-1] // remove trailing &
			}

			req := httptest.NewRequest("GET", url, nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == 200 {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
				assert.NotNil(t, response["data"])
			}

			if tt.expectedStatus != 404 {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestSearchHandler_GetDocument_Extended(t *testing.T) {
	tests := []struct {
		name           string
		docID          string
		mockDocument   *models.Document
		mockError      error
		expectedStatus int
	}{
		{
			name:  "successful document retrieval with full metadata",
			docID: "doc123",
			mockDocument: &models.Document{
				ID:       "doc123",
				FileName: "contract_motion.pdf",
				Text:     "This is a motion regarding contract dispute...",
				DocType:  "Motion",
				Metadata: &models.DocumentMetadata{
					Subject:    "Contract Dispute Motion",
					CaseNumber: "CV-2023-001",
					CaseName:   "Smith vs Jones",
					Court:      &models.CourtInfo{CourtName: "Superior Court"},
					Judge:      &models.Judge{Name: "Judge Johnson"},
					Author:     "Attorney Smith",
					Status:     "Active",
					LegalTags:  []string{"Contract Law", "Dispute Resolution"},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedStatus: 200,
		},
		{
			name:  "document with minimal metadata",
			docID: "doc456",
			mockDocument: &models.Document{
				ID:       "doc456",
				FileName: "simple_order.pdf",
				Text:     "Court order text...",
				DocType:  "Order",
				Metadata: &models.DocumentMetadata{
					Subject: "Simple Order",
				},
			},
			expectedStatus: 200,
		},
		{
			name:           "document not found",
			docID:          "nonexistent",
			mockDocument:   nil,
			mockError:      errors.New("document not found"),
			expectedStatus: 404,
		},
		{
			name:           "service error",
			docID:          "error_doc",
			mockDocument:   nil,
			mockError:      errors.New("database connection failed"),
			expectedStatus: 500,
		},
		{
			name:           "empty document ID in URL",
			docID:          "",
			expectedStatus: 404, // Fiber route doesn't match
		},
		{
			name:           "invalid document ID format",
			docID:          "invalid_doc_id", // Use valid URL-safe format
			mockDocument:   nil,
			mockError:      errors.New("document not found"),
			expectedStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSearchService)
			handler := NewSearchHandler(mockService)
			app := fiber.New()
			app.Get("/documents/:id", handler.GetDocument)

			if tt.docID != "" && tt.expectedStatus == 200 {
				mockService.On("GetDocument", mock.AnythingOfType("*context.timerCtx"), tt.docID).Return(tt.mockDocument, tt.mockError)
			} else if tt.docID != "" && tt.expectedStatus == 404 && tt.mockError != nil {
				mockService.On("GetDocument", mock.AnythingOfType("*context.timerCtx"), tt.docID).Return(tt.mockDocument, tt.mockError)
			} else if tt.docID != "" && tt.expectedStatus == 500 {
				mockService.On("GetDocument", mock.AnythingOfType("*context.timerCtx"), tt.docID).Return(tt.mockDocument, tt.mockError)
			}

			req := httptest.NewRequest("GET", "/documents/"+tt.docID, nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == 200 {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
				assert.NotNil(t, response["data"])
			}

			if tt.docID != "" && (tt.expectedStatus == 200 || tt.expectedStatus == 500 || (tt.expectedStatus == 404 && tt.mockError != nil)) {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestSearchHandler_DeleteDocument_Extended(t *testing.T) {
	tests := []struct {
		name           string
		docID          string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful document deletion",
			docID:          "doc123",
			mockError:      nil,
			expectedStatus: 200,
		},
		{
			name:           "delete nonexistent document",
			docID:          "nonexistent",
			mockError:      nil, // Service handles this gracefully
			expectedStatus: 200,
		},
		{
			name:           "service error during deletion",
			docID:          "error_doc",
			mockError:      errors.New("failed to delete document"),
			expectedStatus: 500,
		},
		{
			name:           "empty document ID",
			docID:          "",
			expectedStatus: 404, // Route doesn't match
		},
		{
			name:           "opensearch connection error",
			docID:          "doc456",
			mockError:      errors.New("connection refused"),
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSearchService)
			handler := NewSearchHandler(mockService)
			app := fiber.New()
			app.Delete("/documents/:id", handler.DeleteDocument)

			if tt.docID != "" && tt.expectedStatus != 404 {
				mockService.On("DeleteDocument", mock.AnythingOfType("*context.timerCtx"), tt.docID).Return(tt.mockError)
			}

			req := httptest.NewRequest("DELETE", "/documents/"+tt.docID, nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == 200 {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "success", response["status"])
				assert.Contains(t, response["message"], "deleted successfully")
			}

			if tt.docID != "" && tt.expectedStatus != 404 {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestValidateSearchRequest_Extended(t *testing.T) {
	tests := []struct {
		name string
		req  *models.SearchRequest
		want *models.SearchRequest
	}{
		{
			name: "complex valid request unchanged",
			req: &models.SearchRequest{
				Query:             "contract dispute",
				Size:              50,
				From:              20,
				SortOrder:         "asc",
				DocType:           "Motion",
				LegalTags:         []string{"Contract Law"},
				CaseNumber:        "CV-2023-001",
				FuzzySearch:       true,
				IncludeHighlights: true,
			},
			want: &models.SearchRequest{
				Query:             "contract dispute",
				Size:              50,
				From:              20,
				SortOrder:         "asc",
				DocType:           "Motion",
				LegalTags:         []string{"Contract Law"},
				CaseNumber:        "CV-2023-001",
				FuzzySearch:       true,
				IncludeHighlights: true,
			},
		},
		{
			name: "zero size gets defaulted",
			req: &models.SearchRequest{
				Size: 0,
			},
			want: &models.SearchRequest{
				Size: models.DefaultSearchSize,
			},
		},
		{
			name: "very large from value accepted",
			req: &models.SearchRequest{
				From: 10000,
				Size: 20,
			},
			want: &models.SearchRequest{
				From: 10000,
				Size: 20,
			},
		},
		{
			name: "invalid sort order variations",
			req: &models.SearchRequest{
				SortOrder: "ASCENDING",
			},
			want: &models.SearchRequest{
				Size:      models.DefaultSearchSize,
				SortOrder: "desc",
			},
		},
		{
			name: "edge case max size",
			req: &models.SearchRequest{
				Size: models.MaxSearchSize,
			},
			want: &models.SearchRequest{
				Size: models.MaxSearchSize,
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

// Helper function to create time pointers
func timePtr(t time.Time) *time.Time {
	return &t
}
