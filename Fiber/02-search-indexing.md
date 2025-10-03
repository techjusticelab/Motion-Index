# Search & Indexing

## Overview
This feature provides full-text search and document indexing capabilities using OpenSearch/Elasticsearch, optimized for legal document retrieval with advanced filtering.

## Current Python Implementation Analysis

### Key Components (from API analysis):
- **`src/handlers/elasticsearch_handler.py`**: Main OpenSearch client wrapper
- **`src/utils/text_normalizer.py`**: Court name normalization for consistent searching

### Endpoints from `server.py`:
- `POST /search` - Advanced document search with filtering
- `GET /legal-tags` - Get all legal types and counts
- `GET /document-types` - Get document types and counts  
- `POST /metadata-field-values` - Get unique values for metadata fields
- `GET /document-stats` - Document index statistics
- `GET /metadata-fields` - Available metadata fields for filtering
- `GET /all-field-options` - All filter options for UI dropdowns
- `POST /update-metadata` - Update document metadata (requires auth)

## Go Package Design

### Package Structure:
```
pkg/
├── search/
│   ├── client/            # OpenSearch client management
│   │   ├── opensearch.go  # Client configuration and connection
│   │   ├── health.go      # Health checks and monitoring
│   │   └── config.go      # Connection configuration
│   ├── indexing/          # Document indexing operations
│   │   ├── indexer.go     # Document indexing logic
│   │   ├── mapping.go     # Index mapping definitions
│   │   ├── bulk.go        # Bulk indexing operations
│   │   └── metadata.go    # Metadata update operations
│   ├── query/             # Search query building
│   │   ├── builder.go     # Query builder interface
│   │   ├── legal.go       # Legal-specific search logic
│   │   ├── filters.go     # Metadata filtering
│   │   └── aggregation.go # Aggregation queries
│   ├── normalizer/        # Text normalization
│   │   ├── court.go       # Court name standardization
│   │   ├── legal.go       # Legal tag normalization
│   │   └── text.go        # General text normalization
│   └── models/            # Search data models
│       ├── query.go       # Search request models
│       ├── result.go      # Search result models
│       └── aggregation.go # Aggregation result models
```

### Core Interfaces:

```go
// SearchService interface for document search operations
type SearchService interface {
    SearchDocuments(ctx context.Context, req *SearchRequest) (*SearchResult, error)
    IndexDocument(ctx context.Context, doc *Document) (string, error)
    BulkIndexDocuments(ctx context.Context, docs []*Document) (*BulkResult, error)
    UpdateDocumentMetadata(ctx context.Context, docID string, metadata map[string]interface{}) error
    DeleteDocument(ctx context.Context, docID string) error
    GetDocument(ctx context.Context, docID string) (*Document, error)
    DocumentExists(ctx context.Context, docID string) (bool, error)
}

// AggregationService interface for metadata aggregations
type AggregationService interface {
    GetLegalTags(ctx context.Context) ([]*TagCount, error)
    GetDocumentTypes(ctx context.Context) ([]*TypeCount, error)
    GetMetadataFieldValues(ctx context.Context, field string, prefix string, size int) ([]*FieldValue, error)
    GetDocumentStats(ctx context.Context) (*DocumentStats, error)
    GetAllFieldOptions(ctx context.Context) (*FieldOptions, error)
}

// QueryBuilder interface for search query construction
type QueryBuilder interface {
    BuildQuery(req *SearchRequest) (map[string]interface{}, error)
    AddTextQuery(query string, fuzzy bool) QueryBuilder
    AddMetadataFilters(filters map[string]interface{}, matchAll bool) QueryBuilder
    AddDateRange(field string, from, to time.Time) QueryBuilder
    AddSorting(field string, order SortOrder) QueryBuilder
    AddPagination(from, size int) QueryBuilder
}
```

### Data Models:

```go
type SearchRequest struct {
    Query              string                 `json:"query,omitempty"`
    DocType            string                 `json:"doc_type,omitempty"`
    CaseNumber         string                 `json:"case_number,omitempty"`
    CaseName           string                 `json:"case_name,omitempty"`
    Judge              []string               `json:"judge,omitempty"`
    Court              []string               `json:"court,omitempty"`
    Author             string                 `json:"author,omitempty"`
    Status             string                 `json:"status,omitempty"`
    LegalTags          []string               `json:"legal_tags,omitempty"`
    LegalTagsMatchAll  bool                   `json:"legal_tags_match_all"`
    DateRange          *DateRange             `json:"date_range,omitempty"`
    Size               int                    `json:"size" validate:"min=1,max=100"`
    Page               int                    `json:"page" validate:"min=1"`
    SortBy             string                 `json:"sort_by,omitempty"`
    SortOrder          SortOrder              `json:"sort_order" validate:"oneof=asc desc"`
    UseFuzzy           bool                   `json:"use_fuzzy"`
}

type SearchResult struct {
    Total        int64                  `json:"total"`
    Hits         []*DocumentHit         `json:"hits"`
    Aggregations map[string]interface{} `json:"aggregations"`
    ProcessTime  time.Duration          `json:"process_time"`
}

type DocumentHit struct {
    ID        string                 `json:"id"`
    Score     float64                `json:"score"`
    Document  *Document              `json:"document"`
    Highlight map[string][]string    `json:"highlight,omitempty"`
}

type DocumentStats struct {
    TotalDocuments int64            `json:"total_documents"`
    IndexSize      string           `json:"index_size"`
    LastUpdated    time.Time        `json:"last_updated"`
    TypeCounts     map[string]int64 `json:"type_counts"`
}

type FieldOptions struct {
    DocTypes     []*FieldValue `json:"doc_types"`
    Categories   []*FieldValue `json:"categories"`
    CaseNumbers  []*FieldValue `json:"case_numbers"`
    Judges       []*FieldValue `json:"judges"`
    Courts       []*FieldValue `json:"courts"`
    LegalTags    []*FieldValue `json:"legal_tags"`
    Statuses     []*FieldValue `json:"statuses"`
}

type FieldValue struct {
    Value string `json:"value"`
    Count int64  `json:"count"`
}
```

## Fiber Handlers

### Search Documents Handler:
```go
func (h *SearchHandler) SearchDocuments(c *fiber.Ctx) error {
    var req SearchRequest
    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid search request")
    }
    
    // Validate request
    if err := h.validator.Validate(&req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, err.Error())
    }
    
    // Execute search
    result, err := h.searchService.SearchDocuments(c.Context(), &req)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    return c.JSON(result)
}
```

### Metadata Fields Handler:
```go
func (h *SearchHandler) GetMetadataFields(c *fiber.Ctx) error {
    fields := []MetadataField{
        {ID: "doc_type", Name: "Document Type", Type: "string"},
        {ID: "category", Name: "Category", Type: "string"},
        {ID: "metadata.case_number", Name: "Case Number", Type: "string"},
        {ID: "metadata.case_name", Name: "Case Name", Type: "string"},
        {ID: "metadata.judge", Name: "Judge", Type: "string"},
        {ID: "metadata.court", Name: "Court", Type: "string"},
        {ID: "metadata.legal_tags", Name: "Legal Tags", Type: "string"},
        {ID: "metadata.author", Name: "Author", Type: "string"},
        {ID: "metadata.status", Name: "Status", Type: "string"},
        {ID: "created_at", Name: "Date", Type: "date"},
    }
    
    return c.JSON(fiber.Map{"fields": fields})
}
```

### Update Metadata Handler (Authenticated):
```go
func (h *SearchHandler) UpdateDocumentMetadata(c *fiber.Ctx) error {
    var req MetadataUpdateRequest
    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid update request")
    }
    
    // Check if document exists
    exists, err := h.searchService.DocumentExists(c.Context(), req.DocumentID)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    if !exists {
        return fiber.NewError(fiber.StatusNotFound, "Document not found")
    }
    
    // Update metadata
    err = h.searchService.UpdateDocumentMetadata(c.Context(), req.DocumentID, req.Metadata)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    // Return updated document
    doc, err := h.searchService.GetDocument(c.Context(), req.DocumentID)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    return c.JSON(fiber.Map{
        "message":  "Document metadata updated successfully",
        "document": doc,
    })
}
```

### Field Options Handler:
```go
func (h *SearchHandler) GetAllFieldOptions(c *fiber.Ctx) error {
    options, err := h.aggregationService.GetAllFieldOptions(c.Context())
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    return c.JSON(options)
}
```

## OpenSearch Configuration

### Index Mapping:
```go
const DocumentMapping = `{
    "mappings": {
        "properties": {
            "file_name": {"type": "keyword"},
            "file_path": {"type": "keyword"},
            "s3_uri": {"type": "keyword"},
            "text": {
                "type": "text",
                "analyzer": "standard"
            },
            "doc_type": {"type": "keyword"},
            "category": {"type": "keyword"},
            "hash": {"type": "keyword"},
            "created_at": {"type": "date"},
            "metadata": {
                "properties": {
                    "document_name": {"type": "text"},
                    "subject": {"type": "text"},
                    "status": {"type": "keyword"},
                    "timestamp": {"type": "date"},
                    "case_name": {"type": "text"},
                    "case_number": {"type": "keyword"},
                    "author": {"type": "keyword"},
                    "judge": {"type": "keyword"},
                    "court": {"type": "keyword"},
                    "legal_tags": {"type": "keyword"}
                }
            }
        }
    },
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0,
        "analysis": {
            "analyzer": {
                "legal_analyzer": {
                    "type": "custom",
                    "tokenizer": "standard",
                    "filter": ["lowercase", "stop"]
                }
            }
        }
    }
}`
```

## Test Strategy

### Unit Tests:
```go
func TestQueryBuilder_BuildQuery(t *testing.T) {
    tests := []struct {
        name    string
        request *SearchRequest
        want    map[string]interface{}
        wantErr bool
    }{
        {
            name: "simple text query",
            request: &SearchRequest{
                Query: "motion to dismiss",
                Size:  10,
                Page:  1,
            },
            want: map[string]interface{}{
                "query": map[string]interface{}{
                    "match": map[string]interface{}{
                        "text": "motion to dismiss",
                    },
                },
                "size": 10,
                "from": 0,
            },
            wantErr: false,
        },
    }
    
    builder := NewQueryBuilder()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := builder.BuildQuery(tt.request)
            if (err != nil) != tt.wantErr {
                t.Errorf("BuildQuery() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Compare results
        })
    }
}
```

### Integration Tests:
- OpenSearch connection and health checks
- Document indexing and retrieval
- Complex search queries with filters
- Aggregation results accuracy
- Performance tests with large datasets

## Implementation Priority

1. **OpenSearch Client** - Connection and health monitoring
2. **Basic Indexing** - Document CRUD operations
3. **Simple Search** - Text queries with pagination
4. **Advanced Filtering** - Metadata filters and date ranges
5. **Aggregations** - Field options and statistics
6. **Bulk Operations** - High-performance batch processing

## Dependencies

### External Libraries:
- `github.com/opensearch-project/opensearch-go/v2` - OpenSearch client
- `github.com/go-playground/validator/v10` - Request validation
- `github.com/elastic/go-elasticsearch/v8` - Alternative Elasticsearch client

### Configuration:
```go
type SearchConfig struct {
    Host        string `env:"ES_HOST" required:"true"`
    Port        int    `env:"ES_PORT" default:"9200"`
    Username    string `env:"ES_USERNAME"`
    Password    string `env:"ES_PASSWORD"`
    APIKey      string `env:"ES_API_KEY"`
    CloudID     string `env:"ES_CLOUD_ID"`
    UseSSL      bool   `env:"ES_USE_SSL" default:"true"`
    Index       string `env:"ES_INDEX" default:"documents"`
    
    // Search Configuration
    MaxResults  int           `env:"MAX_SEARCH_RESULTS" default:"10000"`
    Timeout     time.Duration `env:"SEARCH_TIMEOUT" default:"30s"`
    
    // Bulk Configuration
    BulkSize    int           `env:"BULK_INDEX_SIZE" default:"100"`
    FlushBytes  int           `env:"BULK_FLUSH_BYTES" default:"5242880"` // 5MB
}
```

## Performance Considerations

- **Connection Pooling**: Efficient OpenSearch client management
- **Bulk Operations**: Batch document indexing for high throughput
- **Query Optimization**: Efficient query structure for large datasets
- **Caching**: Aggregate results caching for UI dropdowns
- **Pagination**: Cursor-based pagination for large result sets

## Security Considerations

- **Authentication**: Secure OpenSearch authentication
- **Input Validation**: Prevent injection attacks in search queries
- **Access Control**: User-based document filtering
- **Rate Limiting**: Prevent search abuse