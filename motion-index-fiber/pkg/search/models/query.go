package models

import "time"

// SearchRequest represents a search query with legal-specific filters
type SearchRequest struct {
	Query             string             `json:"query,omitempty" validate:"max=1000"`
	DocType           string             `json:"doc_type,omitempty"`
	CaseNumber        string             `json:"case_number,omitempty"`
	CaseName          string             `json:"case_name,omitempty"`
	Judge             []string           `json:"judge,omitempty"`
	Court             []string           `json:"court,omitempty"`
	Author            string             `json:"author,omitempty"`
	Status            string             `json:"status,omitempty"`
	LegalTags         []string           `json:"legal_tags,omitempty"`
	LegalTagsMatchAll bool               `json:"legal_tags_match_all"`
	DateRange         *DateRange         `json:"date_range,omitempty"`
	Size              int                `json:"size" validate:"min=1,max=100"`
	From              int                `json:"from" validate:"min=0"`
	SortBy            string             `json:"sort_by,omitempty"`
	SortOrder         string             `json:"sort_order,omitempty"`
	IncludeHighlights bool               `json:"include_highlights"`
	FuzzySearch       bool               `json:"fuzzy_search"`
	Filters           interface{}        `json:"filters,omitempty"` // Can be *Filters or map[string]interface{}
	Sort              *SortOptions       `json:"sort,omitempty"`
	Pagination        *PaginationOptions `json:"pagination,omitempty"`
	Highlight         *HighlightOptions  `json:"highlight,omitempty"`
	Limit             int                `json:"limit,omitempty"` // For backward compatibility with tests
}

// DateRange represents a date range filter
type DateRange struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

// SearchResult represents the response from a search query
type SearchResult struct {
	TotalHits    int64                  `json:"total_hits"`
	MaxScore     float64                `json:"max_score,omitempty"`
	Documents    []*SearchDocument      `json:"documents"`
	Aggregations map[string]interface{} `json:"aggregations,omitempty"`
	Took         int64                  `json:"took_ms"`
	TimedOut     bool                   `json:"timed_out"`
}

// SearchDocument represents a document in search results
type SearchDocument struct {
	ID         string                 `json:"id"`
	Score      float64                `json:"score,omitempty"`
	Document   map[string]interface{} `json:"document"`
	Highlights map[string][]string    `json:"highlights,omitempty"`
}

// TagCount represents a legal tag with its document count
type TagCount struct {
	Tag   string `json:"tag"`
	Count int64  `json:"count"`
}

// TypeCount represents a document type with its count
type TypeCount struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

// FieldValue represents a metadata field value
type FieldValue struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// DocumentStats represents overall document statistics
type DocumentStats struct {
	TotalDocuments int64                `json:"total_documents"`
	IndexSize      string               `json:"index_size"`
	TypeCounts     []*TypeCount         `json:"type_counts"`
	TagCounts      []*TagCount          `json:"tag_counts"`
	LastUpdated    time.Time            `json:"last_updated"`
	FieldStats     map[string]FieldStat `json:"field_stats"`
}

// FieldStat represents statistics for a specific field
type FieldStat struct {
	UniqueValues int64 `json:"unique_values"`
	TotalValues  int64 `json:"total_values"`
}

// FieldOptions represents all available filter options
type FieldOptions struct {
	Courts    []*FieldValue `json:"courts"`
	Judges    []*FieldValue `json:"judges"`
	DocTypes  []*FieldValue `json:"doc_types"`
	LegalTags []*FieldValue `json:"legal_tags"`
	Statuses  []*FieldValue `json:"statuses"`
	Authors   []*FieldValue `json:"authors"`
}

// BulkResult represents the result of a bulk operation
type BulkResult struct {
	Took       int64            `json:"took"`
	Errors     bool             `json:"errors"`
	Items      []BulkResultItem `json:"items"`
	Indexed    int              `json:"indexed"`
	Failed     int              `json:"failed"`
	FailedDocs []*BulkFailedDoc `json:"failed_docs,omitempty"`
}

// BulkResultItem represents a single item in bulk operation result
type BulkResultItem struct {
	Index  *BulkItemResult `json:"index,omitempty"`
	Create *BulkItemResult `json:"create,omitempty"`
	Update *BulkItemResult `json:"update,omitempty"`
	Delete *BulkItemResult `json:"delete,omitempty"`
}

// BulkItemResult represents the result of a single bulk item
type BulkItemResult struct {
	ID     string     `json:"_id"`
	Index  string     `json:"_index"`
	Type   string     `json:"_type"`
	Status int        `json:"status"`
	Error  *BulkError `json:"error,omitempty"`
}

// BulkError represents an error in bulk operation
type BulkError struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// BulkFailedDoc represents a document that failed to be indexed
type BulkFailedDoc struct {
	ID     string `json:"id"`
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// SortOrder represents search result ordering
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// DefaultSearchSize is the default number of results to return
const DefaultSearchSize = 20

// MaxSearchSize is the maximum number of results allowed
const MaxSearchSize = 100

// SearchResponse represents the top-level search response
type SearchResponse struct {
	Success    bool         `json:"success"`
	Message    string       `json:"message,omitempty"`
	Data       *SearchResult `json:"data,omitempty"`
	Error      *SearchError  `json:"error,omitempty"`
	RequestID  string       `json:"request_id,omitempty"`
	Timestamp  string       `json:"timestamp"`
	Documents  []*Document  `json:"documents,omitempty"` // For backward compatibility with tests
	Total      int64        `json:"total,omitempty"`     // For backward compatibility with tests
}

// SearchError represents an error in search operations
type SearchError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// IndexStats represents statistics about the search index
type IndexStats struct {
	TotalDocuments int64              `json:"total_documents"`
	IndexSize      string             `json:"index_size"`
	IndexHealth    string             `json:"index_health"`
	ShardInfo      *ShardInfo         `json:"shard_info,omitempty"`
	FieldCounts    map[string]int64   `json:"field_counts"`
	LastUpdated    string             `json:"last_updated"`
	Performance    *PerformanceStats  `json:"performance,omitempty"`
}

// ShardInfo represents shard information for the index
type ShardInfo struct {
	Primary   int `json:"primary"`
	Replicas  int `json:"replicas"`
	Total     int `json:"total"`
	Active    int `json:"active"`
	Failed    int `json:"failed"`
}

// PerformanceStats represents performance metrics
type PerformanceStats struct {
	AvgQueryTime    string `json:"avg_query_time"`
	TotalQueries    int64  `json:"total_queries"`
	IndexingRate    string `json:"indexing_rate"`
	CacheHitRatio   string `json:"cache_hit_ratio"`
}

// AggregationResponse represents the response from aggregation queries
type AggregationResponse struct {
	Success       bool                    `json:"success"`
	Message       string                  `json:"message,omitempty"`
	Aggregations  map[string]interface{}  `json:"aggregations"`
	TotalHits     int64                   `json:"total_hits"`
	RequestID     string                  `json:"request_id,omitempty"`
	Timestamp     string                  `json:"timestamp"`
	Error         *SearchError            `json:"error,omitempty"`
	DocumentTypes []AggregationBucket     `json:"document_types,omitempty"`
	Categories    []AggregationBucket     `json:"categories,omitempty"`
	DateRanges    []AggregationBucket     `json:"date_ranges,omitempty"`
	Courts        []AggregationBucket     `json:"courts,omitempty"`
	Judges        []AggregationBucket     `json:"judges,omitempty"`
}

// SortOptions represents sorting configuration for search queries
type SortOptions struct {
	Field     string    `json:"field"`
	Order     SortOrder `json:"order"`
	Ascending bool      `json:"ascending"` // For backward compatibility with tests
}

// Filters represents search filters that can be applied
type Filters struct {
	DocType     []string            `json:"doc_type,omitempty"`
	Court       []string            `json:"court,omitempty"`
	Judge       []string            `json:"judge,omitempty"`
	Author      []string            `json:"author,omitempty"`
	Status      []string            `json:"status,omitempty"`
	LegalTags   []string            `json:"legal_tags,omitempty"`
	DateRange   *DateRange          `json:"date_range,omitempty"`
	CustomFilters map[string]interface{} `json:"custom_filters,omitempty"`
}

// PaginationOptions represents pagination configuration
type PaginationOptions struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// HighlightOptions represents highlighting configuration
type HighlightOptions struct {
	Fields []string `json:"fields"`
}

// AggregationBucket represents a single aggregation bucket
type AggregationBucket struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
}

