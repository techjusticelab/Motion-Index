package models

import (
	"mime/multipart"
	"time"
)

// ProcessDocumentRequest represents a document processing request
type ProcessDocumentRequest struct {
	File        *multipart.FileHeader `form:"file" validate:"required"`
	Category    string                `form:"category" validate:"omitempty,oneof=motion order contract brief memo other"`
	Description string                `form:"description" validate:"omitempty,max=500"`
	CaseName    string                `form:"case_name" validate:"omitempty,max=200"`
	CaseNumber  string                `form:"case_number" validate:"omitempty,max=50"`
	Author      string                `form:"author" validate:"omitempty,max=100"`
	Judge       string                `form:"judge" validate:"omitempty,max=100"`
	Court       string                `form:"court" validate:"omitempty,max=200"`
	LegalTags   []string              `form:"legal_tags" validate:"omitempty,dive,max=50"`
	Options     *ProcessOptions       `json:"options,omitempty"`
}

// ProcessOptions defines processing options
type ProcessOptions struct {
	ExtractText    bool `json:"extract_text" validate:"omitempty"`
	ClassifyDoc    bool `json:"classify_document" validate:"omitempty"`
	IndexDocument  bool `json:"index_document" validate:"omitempty"`
	StoreDocument  bool `json:"store_document" validate:"omitempty"`
	TimeoutSeconds int  `json:"timeout_seconds" validate:"omitempty,min=1,max=300"`
	RetryCount     int  `json:"retry_count" validate:"omitempty,min=0,max=3"`
}

// BatchProcessRequest represents a batch document processing request
type BatchProcessRequest struct {
	Files       []*multipart.FileHeader `form:"files" validate:"required,min=1,max=10"`
	Category    string                  `form:"category" validate:"omitempty,oneof=motion order contract brief memo other"`
	Description string                  `form:"description" validate:"omitempty,max=500"`
	CaseName    string                  `form:"case_name" validate:"omitempty,max=200"`
	CaseNumber  string                  `form:"case_number" validate:"omitempty,max=50"`
	Options     *ProcessOptions         `json:"options,omitempty"`
}

// SearchDocumentsRequest represents a document search request
type SearchDocumentsRequest struct {
	Query     string            `json:"query" validate:"required,min=1,max=500"`
	Filters   map[string]string `json:"filters" validate:"omitempty"`
	Page      int               `json:"page" validate:"omitempty,min=1"`
	Size      int               `json:"size" validate:"omitempty,min=1,max=100"`
	SortBy    string            `json:"sort_by" validate:"omitempty,oneof=relevance date size name"`
	SortOrder string            `json:"sort_order" validate:"omitempty,oneof=asc desc"`
	Highlight bool              `json:"highlight" validate:"omitempty"`
}

// UpdateMetadataRequest represents a request to update document metadata
type UpdateMetadataRequest struct {
	DocumentID string            `json:"document_id" validate:"required"`
	Metadata   map[string]string `json:"metadata" validate:"required"`
	CaseName   string            `json:"case_name" validate:"omitempty,max=200"`
	CaseNumber string            `json:"case_number" validate:"omitempty,max=50"`
	Author     string            `json:"author" validate:"omitempty,max=100"`
	Judge      string            `json:"judge" validate:"omitempty,max=100"`
	Court      string            `json:"court" validate:"omitempty,max=200"`
	LegalTags  []string          `json:"legal_tags" validate:"omitempty,dive,max=50"`
	Status     string            `json:"status" validate:"omitempty,oneof=draft review approved published archived"`
}

// DeleteDocumentRequest represents a request to delete a document
type DeleteDocumentRequest struct {
	DocumentID string `json:"document_id" validate:"required"`
	Reason     string `json:"reason" validate:"omitempty,max=200"`
}

// AnalyzeRedactionsRequest represents a request to analyze document redactions
type AnalyzeRedactionsRequest struct {
	File        *multipart.FileHeader `form:"file" validate:"required"`
	Sensitivity string                `form:"sensitivity" validate:"omitempty,oneof=low medium high"`
}

// BulkUploadRequest represents a bulk upload request
type BulkUploadRequest struct {
	Files       []*multipart.FileHeader `form:"files" validate:"required,min=1,max=50"`
	Category    string                  `form:"category" validate:"omitempty,oneof=motion order contract brief memo other"`
	CaseName    string                  `form:"case_name" validate:"omitempty,max=200"`
	CaseNumber  string                  `form:"case_number" validate:"omitempty,max=50"`
	AutoProcess bool                    `form:"auto_process" validate:"omitempty"`
}

// GetDocumentRequest represents a request to retrieve a document
type GetDocumentRequest struct {
	DocumentID string `json:"document_id" validate:"required"`
	Format     string `json:"format" validate:"omitempty,oneof=json full metadata"`
}

// GetDocumentStatsRequest represents a request for document statistics
type GetDocumentStatsRequest struct {
	DateRange   *DateRange `json:"date_range,omitempty"`
	GroupBy     string     `json:"group_by" validate:"omitempty,oneof=date category court judge author"`
	Granularity string     `json:"granularity" validate:"omitempty,oneof=day week month year"`
}

// DateRange represents a date range filter
type DateRange struct {
	Start time.Time `json:"start" validate:"omitempty"`
	End   time.Time `json:"end" validate:"omitempty,gtfield=Start"`
}

// HealthCheckRequest represents a health check request
type HealthCheckRequest struct {
	Component string `json:"component" validate:"omitempty,oneof=storage search pipeline classifier extractor"`
	Deep      bool   `json:"deep" validate:"omitempty"`
}

// Default values for process options
func DefaultProcessOptions() *ProcessOptions {
	return &ProcessOptions{
		ExtractText:    true,
		ClassifyDoc:    true,
		IndexDocument:  true,
		StoreDocument:  true,
		TimeoutSeconds: 120,
		RetryCount:     1,
	}
}

// Validate checks if ProcessOptions are valid
func (opts *ProcessOptions) Validate() error {
	if opts == nil {
		return nil
	}

	if opts.TimeoutSeconds < 1 || opts.TimeoutSeconds > 300 {
		opts.TimeoutSeconds = 120
	}

	if opts.RetryCount < 0 || opts.RetryCount > 3 {
		opts.RetryCount = 1
	}

	return nil
}

// ApplyDefaults applies default values to process options
func (opts *ProcessOptions) ApplyDefaults() {
	if opts == nil {
		return
	}

	if opts.TimeoutSeconds == 0 {
		opts.TimeoutSeconds = 120
	}

	if opts.RetryCount == 0 {
		opts.RetryCount = 1
	}
}
