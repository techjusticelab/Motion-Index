# Document Redaction

## Overview
This feature provides California legal compliance for document redaction, including pattern-based detection, AI-enhanced analysis, and PDF manipulation for creating redacted copies.

## Current Python Implementation Analysis

### Key Components (from API analysis):
- **`src/handlers/redaction_handler.py`**: Main redaction processing with California legal compliance
- **California Legal Codes**: CCP, WIC, PC, etc. compliance patterns

### Endpoints from `server.py`:
- `POST /analyze-redactions` - PDF redaction analysis without requiring Elasticsearch
- `POST /redact-document` - Create redacted version from stored document
- Redaction analysis in `POST /categorise` - Optional redaction during upload

### Current Features:
- Pattern-based sensitive content detection
- AI-enhanced analysis using OpenAI
- PDF manipulation with PyMuPDF
- Redacted copy generation
- California-specific legal code compliance

## Go Package Design

### Package Structure:
```
pkg/
├── redaction/
│   ├── analyzer/            # Redaction analysis engine
│   │   ├── patterns.go      # Pattern-based detection
│   │   ├── ai.go           # AI-enhanced analysis
│   │   ├── legal.go        # California legal compliance
│   │   └── interface.go    # Analyzer interface
│   ├── processor/          # PDF processing and redaction
│   │   ├── pdf.go          # PDF manipulation
│   │   ├── redactor.go     # Redaction application
│   │   ├── renderer.go     # Redacted PDF generation
│   │   └── interface.go    # Processor interface
│   ├── rules/              # Redaction rules and patterns
│   │   ├── california.go   # California-specific rules
│   │   ├── patterns.go     # Regex patterns for detection
│   │   ├── sensitive.go    # Sensitive information types
│   │   └── legal_codes.go  # Legal code references
│   ├── models/             # Redaction data models
│   │   ├── analysis.go     # Analysis result models
│   │   ├── pattern.go      # Pattern match models
│   │   ├── redaction.go    # Redaction operation models
│   │   └── report.go       # Redaction report models
│   └── service/            # Redaction service
│       ├── service.go      # Main redaction service
│       ├── batch.go        # Batch redaction operations
│       └── config.go       # Configuration
```

### Core Interfaces:

```go
// RedactionAnalyzer interface for detecting sensitive content
type RedactionAnalyzer interface {
    AnalyzeDocument(ctx context.Context, doc *Document) (*AnalysisResult, error)
    AnalyzePDF(ctx context.Context, pdfPath string) (*AnalysisResult, error)
    FindPatterns(ctx context.Context, text string) ([]*PatternMatch, error)
    FindSensitiveInfo(ctx context.Context, text string) ([]*SensitiveMatch, error)
}

// RedactionProcessor interface for applying redactions
type RedactionProcessor interface {
    RedactDocument(ctx context.Context, req *RedactionRequest) (*RedactionResult, error)
    CreateRedactedCopy(ctx context.Context, sourcePath string, redactions []*Redaction) (string, error)
    ApplyRedactions(ctx context.Context, pdfPath string, redactions []*Redaction) error
    GenerateRedactionReport(ctx context.Context, redactions []*Redaction) (*RedactionReport, error)
}

// RedactionService interface for complete redaction workflow
type RedactionService interface {
    AnalyzeForRedactions(ctx context.Context, req *AnalysisRequest) (*AnalysisResult, error)
    ProcessRedactions(ctx context.Context, req *RedactionRequest) (*RedactionResult, error)
    GetRedactionReport(ctx context.Context, documentID string) (*RedactionReport, error)
    ValidateRedactions(ctx context.Context, redactions []*Redaction) error
}
```

### Data Models:

```go
type AnalysisRequest struct {
    DocumentID       string `json:"document_id,omitempty"`
    FilePath         string `json:"file_path,omitempty"`
    Content          []byte `json:"-"`              // Raw file content
    AnalysisType     string `json:"analysis_type"`  // "pattern", "ai", "both"
    CaliforniaRules  bool   `json:"california_rules" default:"true"`
    CustomPatterns   []string `json:"custom_patterns,omitempty"`
}

type AnalysisResult struct {
    DocumentID        string            `json:"document_id,omitempty"`
    FileName          string            `json:"file_name"`
    RedactionsFound   int               `json:"redactions_found"`
    PatternsDetected  []*PatternMatch   `json:"patterns_detected"`
    SensitiveContent  []*SensitiveMatch `json:"sensitive_content"`
    Recommendations   []*Recommendation `json:"recommendations"`
    ComplianceStatus  *ComplianceStatus `json:"compliance_status"`
    AnalysisTime      time.Duration     `json:"analysis_time"`
    CreatedAt         time.Time         `json:"created_at"`
}

type PatternMatch struct {
    Type        string    `json:"type"`         // "ssn", "phone", "email", "case_number"
    Pattern     string    `json:"pattern"`      // Regex pattern used
    Match       string    `json:"match"`        // Actual matched text
    Location    *Location `json:"location"`     // Position in document
    Confidence  float64   `json:"confidence"`   // Match confidence (0-1)
    Sensitive   bool      `json:"sensitive"`    // Requires redaction
    LegalCode   string    `json:"legal_code,omitempty"` // Relevant legal code
}

type SensitiveMatch struct {
    Type         string    `json:"type"`          // "personal_info", "medical", "financial"
    Content      string    `json:"content"`       // Detected content
    Context      string    `json:"context"`       // Surrounding context
    Location     *Location `json:"location"`      // Position in document
    Severity     string    `json:"severity"`      // "low", "medium", "high"
    Justification string   `json:"justification"` // Why it's sensitive
    LegalBasis   string    `json:"legal_basis"`   // Legal requirement
}

type Location struct {
    Page       int     `json:"page"`
    X          float64 `json:"x"`
    Y          float64 `json:"y"`
    Width      float64 `json:"width"`
    Height     float64 `json:"height"`
    TextOffset int     `json:"text_offset"` // Character offset in text
    TextLength int     `json:"text_length"` // Length of matched text
}

type Redaction struct {
    ID        string    `json:"id"`
    Type      string    `json:"type"`      // Type of redaction
    Location  *Location `json:"location"`  // Where to redact
    Method    string    `json:"method"`    // "blackout", "blur", "replace"
    Color     string    `json:"color"`     // Redaction color (default: black)
    Applied   bool      `json:"applied"`   // Whether redaction was applied
    Reason    string    `json:"reason"`    // Reason for redaction
}

type RedactionRequest struct {
    DocumentID       string       `json:"document_id"`
    ApplyRedactions  bool         `json:"apply_redactions"`
    Redactions       []*Redaction `json:"redactions,omitempty"`
    OutputPath       string       `json:"output_path,omitempty"`
    CreateCopy       bool         `json:"create_copy" default:"true"`
    PreserveCitations bool        `json:"preserve_citations" default:"true"`
}

type RedactionResult struct {
    DocumentID       string            `json:"document_id"`
    OriginalPath     string            `json:"original_path"`
    RedactedPath     string            `json:"redacted_path,omitempty"`
    RedactedURL      string            `json:"redacted_url,omitempty"`
    RedactionsApplied int              `json:"redactions_applied"`
    FileSize         int64             `json:"file_size"`
    ProcessingTime   time.Duration     `json:"processing_time"`
    Report           *RedactionReport  `json:"report,omitempty"`
}

type ComplianceStatus struct {
    Compliant        bool     `json:"compliant"`
    ViolatedCodes    []string `json:"violated_codes,omitempty"`
    RequiredActions  []string `json:"required_actions,omitempty"`
    ComplianceLevel  string   `json:"compliance_level"` // "full", "partial", "none"
}
```

## Fiber Handlers

### Analyze Redactions Handler:
```go
func (h *RedactionHandler) AnalyzeRedactions(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "No file provided")
    }
    
    // Validate PDF file type
    if !strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
        return fiber.NewError(fiber.StatusBadRequest, "Only PDF files are supported for redaction analysis")
    }
    
    // Save file temporarily
    tempPath, err := h.saveTemporaryFile(file)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Failed to save file")
    }
    defer os.Remove(tempPath)
    
    // Analyze document for redactions
    result, err := h.redactionService.AnalyzeForRedactions(c.Context(), &AnalysisRequest{
        FilePath:        tempPath,
        AnalysisType:    "both", // Pattern + AI analysis
        CaliforniaRules: true,
    })
    
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message":           fmt.Sprintf("Redaction analysis completed. Found %d potential redactions.", result.RedactionsFound),
        "filename":          file.Filename,
        "redaction_analysis": result,
        "status":            "success",
    })
}
```

### Create Redacted Document Handler:
```go
func (h *RedactionHandler) CreateRedactedDocument(c *fiber.Ctx) error {
    var req RedactionRequest
    if err := c.BodyParser(&req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid redaction request")
    }
    
    // Validate request
    if req.DocumentID == "" {
        return fiber.NewError(fiber.StatusBadRequest, "Document ID required")
    }
    
    // Get document from storage service
    doc, err := h.documentService.GetDocument(c.Context(), req.DocumentID)
    if err != nil {
        return fiber.NewError(fiber.StatusNotFound, "Document not found")
    }
    
    // Validate PDF file
    if !strings.HasSuffix(strings.ToLower(doc.FileName), ".pdf") {
        return fiber.NewError(fiber.StatusBadRequest, "Only PDF files can be redacted")
    }
    
    // Download file from storage
    filePath, err := h.downloadDocumentFile(c.Context(), doc)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Failed to access document file")
    }
    defer os.Remove(filePath)
    
    // Process redactions
    result, err := h.redactionService.ProcessRedactions(c.Context(), &RedactionRequest{
        DocumentID:      req.DocumentID,
        ApplyRedactions: req.ApplyRedactions,
        Redactions:      req.Redactions,
        CreateCopy:      true,
    })
    
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message":    "Document redaction completed",
        "document_id": req.DocumentID,
        "redaction_result": result,
    })
}
```

## Pattern Detection Implementation

### California Legal Patterns:
```go
var CaliforniaPatterns = []RedactionPattern{
    {
        Name:        "SSN",
        Pattern:     `\b\d{3}-?\d{2}-?\d{4}\b`,
        Type:        "personal_identifier",
        Sensitivity: "high",
        LegalCode:   "CCP §1985.3",
        Description: "Social Security Number",
    },
    {
        Name:        "Case Number",
        Pattern:     `(?i)case\s*(?:no|number|#)[\s:]*([A-Z0-9-]+)`,
        Type:        "court_identifier",
        Sensitivity: "medium",
        LegalCode:   "CCP §367.75",
        Description: "Court case number",
    },
    {
        Name:        "Birth Date",
        Pattern:     `\b\d{1,2}[\/\-]\d{1,2}[\/\-]\d{4}\b`,
        Type:        "personal_info",
        Sensitivity: "high",
        LegalCode:   "WIC §827",
        Description: "Date of birth",
    },
    {
        Name:        "Phone Number",
        Pattern:     `\b\d{3}[-.]?\d{3}[-.]?\d{4}\b`,
        Type:        "contact_info",
        Sensitivity: "medium",
        LegalCode:   "CCP §1985.3",
        Description: "Phone number",
    },
    {
        Name:        "Address",
        Pattern:     `\b\d+\s+[A-Za-z\s]+(?:Street|St|Avenue|Ave|Road|Rd|Drive|Dr|Lane|Ln|Boulevard|Blvd)\b`,
        Type:        "address",
        Sensitivity: "medium",
        LegalCode:   "CCP §1985.3",
        Description: "Street address",
    },
}

type RedactionPattern struct {
    Name        string  `json:"name"`
    Pattern     string  `json:"pattern"`
    Type        string  `json:"type"`
    Sensitivity string  `json:"sensitivity"`
    LegalCode   string  `json:"legal_code"`
    Description string  `json:"description"`
    Enabled     bool    `json:"enabled"`
    Weight      float64 `json:"weight"`
}
```

### Pattern Matcher:
```go
func (a *PatternAnalyzer) FindPatterns(ctx context.Context, text string) ([]*PatternMatch, error) {
    var matches []*PatternMatch
    
    for _, pattern := range a.patterns {
        if !pattern.Enabled {
            continue
        }
        
        regex, err := regexp.Compile(pattern.Pattern)
        if err != nil {
            continue // Skip invalid patterns
        }
        
        // Find all matches
        regexMatches := regex.FindAllStringSubmatch(text, -1)
        regexIndexes := regex.FindAllStringIndex(text, -1)
        
        for i, match := range regexMatches {
            if len(match) == 0 {
                continue
            }
            
            matchText := match[0]
            startIndex := regexIndexes[i][0]
            
            // Calculate confidence based on pattern specificity
            confidence := a.calculateConfidence(pattern, matchText)
            
            matches = append(matches, &PatternMatch{
                Type:       pattern.Type,
                Pattern:    pattern.Pattern,
                Match:      matchText,
                Location: &Location{
                    TextOffset: startIndex,
                    TextLength: len(matchText),
                },
                Confidence: confidence,
                Sensitive:  pattern.Sensitivity == "high" || pattern.Sensitivity == "medium",
                LegalCode:  pattern.LegalCode,
            })
        }
    }
    
    return matches, nil
}

func (a *PatternAnalyzer) calculateConfidence(pattern RedactionPattern, match string) float64 {
    confidence := 0.5 // Base confidence
    
    // Adjust based on pattern specificity
    switch pattern.Type {
    case "personal_identifier":
        confidence = 0.9
    case "contact_info":
        confidence = 0.8
    case "address":
        confidence = 0.7
    case "court_identifier":
        confidence = 0.85
    }
    
    // Adjust based on match characteristics
    if len(match) > 20 {
        confidence -= 0.1 // Longer matches might be false positives
    }
    
    // Apply pattern weight
    confidence *= pattern.Weight
    
    // Ensure confidence is within bounds
    if confidence > 1.0 {
        confidence = 1.0
    }
    if confidence < 0.0 {
        confidence = 0.0
    }
    
    return confidence
}
```

## AI-Enhanced Analysis

### OpenAI Integration:
```go
type AIAnalyzer struct {
    client  *openai.Client
    config  *AIConfig
}

func (a *AIAnalyzer) FindSensitiveInfo(ctx context.Context, text string) ([]*SensitiveMatch, error) {
    prompt := fmt.Sprintf(`
Analyze the following legal document text and identify sensitive information that should be redacted according to California privacy laws and court rules.

Focus on:
1. Personal identifiers (SSN, driver's license, etc.)
2. Medical information
3. Financial information
4. Minor's information
5. Victim information in criminal cases
6. Confidential attorney-client communications

For each sensitive item found, provide:
- Type of sensitive information
- Exact text to redact
- Reason for redaction
- Relevant legal code or rule

Text to analyze:
%s

Respond in JSON format with an array of sensitive items.
`, text)
    
    resp, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT4,
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    openai.ChatMessageRoleSystem,
                Content: "You are an expert in California legal document redaction and privacy law.",
            },
            {
                Role:    openai.ChatMessageRoleUser,
                Content: prompt,
            },
        },
        MaxTokens:   2000,
        Temperature: 0.1, // Low temperature for consistent results
    })
    
    if err != nil {
        return nil, fmt.Errorf("AI analysis failed: %w", err)
    }
    
    // Parse JSON response
    var aiResults []struct {
        Type          string `json:"type"`
        Content       string `json:"content"`
        Reason        string `json:"reason"`
        LegalBasis    string `json:"legal_basis"`
        Severity      string `json:"severity"`
    }
    
    if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &aiResults); err != nil {
        return nil, fmt.Errorf("failed to parse AI response: %w", err)
    }
    
    // Convert to SensitiveMatch objects
    var matches []*SensitiveMatch
    for _, result := range aiResults {
        // Find location of content in text
        index := strings.Index(text, result.Content)
        if index == -1 {
            continue // Content not found
        }
        
        matches = append(matches, &SensitiveMatch{
            Type:         result.Type,
            Content:      result.Content,
            Context:      a.extractContext(text, index, 50),
            Location: &Location{
                TextOffset: index,
                TextLength: len(result.Content),
            },
            Severity:      result.Severity,
            Justification: result.Reason,
            LegalBasis:    result.LegalBasis,
        })
    }
    
    return matches, nil
}
```

## PDF Processing Implementation

### PDF Redaction:
```go
type PDFProcessor struct {
    tempDir string
}

func (p *PDFProcessor) CreateRedactedCopy(ctx context.Context, sourcePath string, redactions []*Redaction) (string, error) {
    // Create output path
    outputPath := filepath.Join(p.tempDir, fmt.Sprintf("redacted_%d.pdf", time.Now().Unix()))
    
    // Copy original file
    if err := p.copyFile(sourcePath, outputPath); err != nil {
        return "", fmt.Errorf("failed to copy source file: %w", err)
    }
    
    // Open PDF for editing
    doc, err := pdf.Open(outputPath)
    if err != nil {
        return "", fmt.Errorf("failed to open PDF: %w", err)
    }
    defer doc.Close()
    
    // Apply redactions
    for _, redaction := range redactions {
        if err := p.applyRedaction(doc, redaction); err != nil {
            return "", fmt.Errorf("failed to apply redaction %s: %w", redaction.ID, err)
        }
    }
    
    // Save redacted PDF
    if err := doc.Save(); err != nil {
        return "", fmt.Errorf("failed to save redacted PDF: %w", err)
    }
    
    return outputPath, nil
}

func (p *PDFProcessor) applyRedaction(doc *pdf.Document, redaction *Redaction) error {
    page := doc.Page(redaction.Location.Page)
    if page == nil {
        return fmt.Errorf("invalid page number: %d", redaction.Location.Page)
    }
    
    // Create redaction rectangle
    rect := pdf.Rect{
        X:      redaction.Location.X,
        Y:      redaction.Location.Y,
        Width:  redaction.Location.Width,
        Height: redaction.Location.Height,
    }
    
    // Apply redaction based on method
    switch redaction.Method {
    case "blackout":
        return page.AddBlackoutRect(rect, redaction.Color)
    case "blur":
        return page.AddBlurRect(rect)
    case "replace":
        return page.ReplaceTextInRect(rect, "[REDACTED]")
    default:
        return page.AddBlackoutRect(rect, "#000000") // Default to black
    }
}
```

## Test Strategy

### Unit Tests:
```go
func TestPatternAnalyzer_FindPatterns(t *testing.T) {
    tests := []struct {
        name     string
        text     string
        patterns []RedactionPattern
        want     int
        wantErr  bool
    }{
        {
            name: "find SSN",
            text: "John Doe, SSN: 123-45-6789, was born on 01/01/1990",
            patterns: []RedactionPattern{
                {
                    Name:    "SSN",
                    Pattern: `\b\d{3}-?\d{2}-?\d{4}\b`,
                    Type:    "personal_identifier",
                    Enabled: true,
                    Weight:  1.0,
                },
            },
            want:    1,
            wantErr: false,
        },
    }
    
    analyzer := NewPatternAnalyzer(tests[0].patterns)
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            matches, err := analyzer.FindPatterns(context.Background(), tt.text)
            if (err != nil) != tt.wantErr {
                t.Errorf("FindPatterns() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if len(matches) != tt.want {
                t.Errorf("FindPatterns() found %d matches, want %d", len(matches), tt.want)
            }
        })
    }
}
```

### Integration Tests:
- PDF processing and redaction application
- AI analysis accuracy and performance
- California legal compliance validation
- End-to-end redaction workflow

## Implementation Priority

1. **Pattern Detection** - Core regex-based sensitive content detection
2. **PDF Processing** - Basic PDF manipulation and redaction
3. **California Rules** - Legal compliance patterns and validation
4. **AI Integration** - OpenAI-enhanced analysis
5. **Batch Processing** - Multi-document redaction workflows
6. **Reporting** - Detailed redaction reports and compliance tracking

## Dependencies

### External Libraries:
- `github.com/unidoc/unipdf/v3` - PDF processing and manipulation
- `github.com/sashabaranov/go-openai` - OpenAI API integration
- `github.com/go-playground/validator/v10` - Input validation

### Configuration:
```go
type RedactionConfig struct {
    // AI Configuration
    OpenAIAPIKey    string `env:"OPENAI_API_KEY"`
    OpenAIModel     string `env:"OPENAI_MODEL" default:"gpt-4"`
    
    // Processing Configuration
    TempDirectory   string        `env:"TEMP_DIR" default:"/tmp"`
    MaxFileSize     int64         `env:"MAX_REDACTION_FILE_SIZE" default:"52428800"` // 50MB
    ProcessTimeout  time.Duration `env:"REDACTION_TIMEOUT" default:"5m"`
    
    // Legal Configuration
    CaliforniaRules bool     `env:"CALIFORNIA_RULES" default:"true"`
    CustomPatterns  []string `env:"CUSTOM_PATTERNS"`
    
    // Output Configuration
    DefaultMethod   string `env:"DEFAULT_REDACTION_METHOD" default:"blackout"`
    DefaultColor    string `env:"DEFAULT_REDACTION_COLOR" default:"#000000"`
    PreserveCites   bool   `env:"PRESERVE_CITATIONS" default:"true"`
}
```

## Performance Considerations

- **Memory Management**: Efficient PDF processing for large documents
- **Concurrent Analysis**: Parallel pattern matching and AI analysis
- **Caching**: Cache AI analysis results for similar content
- **Streaming**: Process large PDFs without loading entirely into memory

## Security Considerations

- **Secure Deletion**: Ensure temporary files are securely deleted
- **Content Validation**: Validate PDF structure before processing
- **AI Privacy**: Ensure sensitive content is not logged by AI services
- **Access Control**: Restrict redaction operations to authorized users