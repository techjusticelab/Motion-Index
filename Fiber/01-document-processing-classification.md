# Document Processing & Classification

## Overview
This feature handles multi-format document processing and AI-powered classification for legal documents in the Motion-Index system, built with Go Fiber.

## Current Python Implementation Analysis

### Key Components (from API analysis):
- **`src/core/document_processor.py`**: Main orchestration with 28-core optimization
- **`src/handlers/file_processor.py`**: Multi-format text extraction (textract, pytesseract, pdf2image, tabula, pdfminer, PyPDF2)
- **`src/handlers/document_classifier.py`**: AI classification using OpenAI ChatGPT
- **`src/handlers/high_performance_classifier.py`**: Concurrent processing engine

### Endpoints from `server.py`:
- `POST /categorise` - Upload and process documents with optional redaction
- `POST /analyze-redactions` - PDF redaction analysis only

## Go Package Design

### Package Structure:
```
pkg/
├── processing/
│   ├── extractor/          # Text extraction from various formats  
│   │   ├── pdf.go         # PDF text extraction
│   │   ├── docx.go        # DOCX handling
│   │   ├── ocr.go         # OCR for scanned documents
│   │   └── interface.go   # TextExtractor interface
│   ├── classifier/         # AI-powered classification
│   │   ├── openai.go      # OpenAI ChatGPT integration
│   │   ├── legal.go       # Legal tag validation
│   │   └── interface.go   # Classifier interface
│   ├── pipeline/          # Document processing pipeline
│   │   ├── processor.go   # Main processing orchestration
│   │   ├── worker.go      # Worker pool implementation
│   │   └── batch.go       # Batch processing logic
│   └── models/            # Data models
│       ├── document.go    # Document struct
│       ├── metadata.go    # Metadata struct
│       └── result.go      # Processing result types
```

### Core Interfaces:

```go
// TextExtractor interface for different file formats
type TextExtractor interface {
    Extract(ctx context.Context, file io.Reader, filename string) (*ExtractionResult, error)
    CanProcess(filename string) bool
    GetMetadata(file io.Reader) (*FileMetadata, error)
}

// Classifier interface for AI-powered classification
type Classifier interface {
    ClassifyDocument(ctx context.Context, text string, metadata *FileMetadata) (*ClassificationResult, error)
    ValidateLegalTags(tags []string) []string
}

// Processor interface for document processing pipeline
type Processor interface {
    ProcessDocument(ctx context.Context, file *FileUpload) (*ProcessingResult, error)
    ProcessBatch(ctx context.Context, files []*FileUpload) ([]*ProcessingResult, error)
}
```

### Data Models:

```go
type Document struct {
    ID          string            `json:"id"`
    FileName    string            `json:"file_name"`
    FilePath    string            `json:"file_path"`
    FileURL     string            `json:"file_url,omitempty"`
    S3URI       string            `json:"s3_uri,omitempty"`
    Text        string            `json:"text"`
    DocType     string            `json:"doc_type"`
    Category    string            `json:"category,omitempty"`
    Hash        string            `json:"hash"`
    CreatedAt   time.Time         `json:"created_at"`
    Metadata    *Metadata         `json:"metadata"`
}

type Metadata struct {
    DocumentName string    `json:"document_name"`
    Subject      string    `json:"subject"`
    Status       string    `json:"status,omitempty"`
    Timestamp    time.Time `json:"timestamp,omitempty"`
    CaseName     string    `json:"case_name,omitempty"`
    CaseNumber   string    `json:"case_number,omitempty"`
    Author       string    `json:"author,omitempty"`
    Judge        string    `json:"judge,omitempty"`
    Court        string    `json:"court,omitempty"`
    LegalTags    []string  `json:"legal_tags,omitempty"`
}

type ProcessingResult struct {
    Document    *Document            `json:"document"`
    Success     bool                 `json:"success"`
    Error       string               `json:"error,omitempty"`
    ProcessTime time.Duration        `json:"process_time"`
    Redactions  *RedactionAnalysis   `json:"redaction_analysis,omitempty"`
}
```

## Fiber Handlers

### Document Upload Handler:
```go
func (h *ProcessingHandler) UploadDocument(c *fiber.Ctx) error {
    // Parse multipart form
    file, err := c.FormFile("file")
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "No file provided")
    }
    
    // Get optional redaction flag
    applyRedaction := c.FormValue("apply_redaction") == "true"
    
    // Process document
    result, err := h.processor.ProcessDocument(c.Context(), &FileUpload{
        File:           file,
        ApplyRedaction: applyRedaction,
    })
    
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Document categorised successfully",
        "document": result.Document,
        "redaction_analysis": result.Redactions,
    })
}
```

### Redaction Analysis Handler:
```go
func (h *ProcessingHandler) AnalyzeRedactions(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "No file provided")
    }
    
    // Validate PDF file type
    if !strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
        return fiber.NewError(fiber.StatusBadRequest, "Only PDF files are supported")
    }
    
    // Analyze redactions only
    analysis, err := h.redactionService.AnalyzeDocument(c.Context(), file)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message":           "Redaction analysis completed",
        "filename":          file.Filename,
        "redaction_analysis": analysis,
        "status":           "success",
    })
}
```

## Test Strategy

### Unit Tests:
```go
func TestPDFExtractor_Extract(t *testing.T) {
    tests := []struct {
        name     string
        filename string
        want     string
        wantErr  bool
    }{
        {
            name:     "valid PDF",
            filename: "test.pdf",
            want:     "extracted text content",
            wantErr:  false,
        },
        {
            name:     "corrupted PDF",
            filename: "corrupted.pdf",
            want:     "",
            wantErr:  true,
        },
    }
    
    extractor := NewPDFExtractor()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests:
- End-to-end document processing pipeline
- Claude API integration tests
- Error handling and recovery tests
- Performance tests with concurrent processing

## Implementation Priority

1. **Text Extraction** - Core functionality for PDF, DOCX, TXT
2. **Basic Classification** - Simple rule-based classification  
3. **AI Integration** - Claude/Ollama API integration
4. **Worker Pools** - Concurrent processing implementation
5. **Advanced Features** - OCR, table extraction, batch processing

## Dependencies

### External Libraries:
- `github.com/unidoc/unipdf/v3` - PDF text extraction
- `github.com/nguyenthenguyen/docx` - DOCX processing  
- `github.com/sashabaranov/go-openai` - OpenAI ChatGPT API
- `github.com/aws/aws-sdk-go-v2` - S3 integration for storage
- `gocv.io/x/gocv` - OpenCV for OCR (optional)

### Configuration:
```go
type ProcessingConfig struct {
    MaxFileSize     int64         `env:"MAX_FILE_SIZE" default:"104857600"` // 100MB
    MaxWorkers      int           `env:"MAX_WORKERS" default:"10"`
    BatchSize       int           `env:"BATCH_SIZE" default:"50"`
    ProcessTimeout  time.Duration `env:"PROCESS_TIMEOUT" default:"5m"`
    
    // AI Configuration
    OpenAIAPIKey    string `env:"OPENAI_API_KEY" required:"true"`
    OpenAIModel     string `env:"OPENAI_MODEL" default:"gpt-4"`
}
```

## Performance Considerations

- **Streaming Processing**: Use `io.Reader` interfaces for memory efficiency
- **Worker Pools**: Limit concurrent processing to prevent resource exhaustion
- **Context Cancellation**: Implement proper context handling for timeouts
- **Memory Management**: Use object pools for frequent allocations
- **Error Recovery**: Graceful degradation when AI services are unavailable

## Security Considerations

- **File Validation**: Strict file type and size validation
- **Sanitization**: Clean extracted text before AI processing
- **Rate Limiting**: Prevent abuse of AI APIs
- **Virus Scanning**: Optional malware detection integration