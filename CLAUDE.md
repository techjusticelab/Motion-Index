# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Motion-Index is a legal document processing system designed for serverless cloud deployment on DigitalOcean. The system consists of four main components:

1. **API** - Python/FastAPI backend with unified storage (cloud/local) and Elasticsearch integration
2. **Fiber** - Go/Fiber backend for high-performance document processing and classification
3. **Web** - SvelteKit frontend for document search and management
4. **Database** - Managed OpenSearch with 23,500+ migrated legal documents

## Architecture

### Production Implementation (DigitalOcean)
- **Document Storage**: DigitalOcean Spaces with global CDN
- **Search Engine**: DigitalOcean Managed OpenSearch cluster
- **Authentication**: Supabase for user management
- **File Processing**: Cloud-based text extraction and multi-model AI classification with enhanced date extraction
- **Deployment**: DigitalOcean App Platform with auto-scaling

### Local Development
- **Document Storage**: Local filesystem (`API/data/documents/`)
- **Search Engine**: DigitalOcean Managed OpenSearch (same as production)
- **Authentication**: Supabase (same as production)
- **File Processing**: Local text extraction (textract) and multi-model AI classification (OpenAI, Claude, Ollama)

### API Backend Structure

The API uses a modular handler-based architecture:

**Active Handlers** (`src/handlers/`):
- `elasticsearch_handler.py`: Search operations and document indexing (OpenSearch client)
- `storage_handler.py`: Unified storage (DigitalOcean Spaces + local filesystem)
- `file_processor.py`: Multi-format text extraction
- `document_classifier.py`: AI-powered document classification
- `redaction_handler.py`: PDF redaction detection and processing

**Core Components**:
- `src/core/document_processor.py`: Processing pipeline orchestration
- `src/middleware/auth.py`: JWT authentication with Supabase
- `src/models/`: Data models and schemas
- `server.py`: Main FastAPI application

### Fiber Backend Structure (Go)

The Fiber backend provides high-performance document processing:

**Command-line Tools** (`cmd/`):
- `cmd/api-classifier/main.go`: Single-threaded document classification script
- `cmd/api-batch-classifier/main.go`: Multi-threaded batch document classification
- `cmd/setup-index/main.go`: OpenSearch index setup and management
- `cmd/inspect-index/main.go`: Index inspection and debugging tools
- `cmd/server/main.go`: Main Fiber web server

**Processing Components** (`pkg/processing/`):
- `classifier/`: Multi-model AI classification with unified prompts and enhanced date extraction
  - `prompts.go`: Centralized prompt templates for all AI models
  - `date_extraction.go`: Comprehensive date parsing and validation utilities
  - `openai.go`, `claude.go`, `ollama.go`: Model-specific implementations
- `extractor/`: Text extraction from PDFs, DOCX, OCR
- `pipeline/`: Document processing pipeline with date field integration
- `queue/`: Priority queue and rate limiting for batch processing

**Storage & Search** (`pkg/`):
- `cloud/digitalocean/spaces/`: DigitalOcean Spaces integration
- `search/`: OpenSearch client and query building
- `storage/`: Unified storage abstraction layer

### Web Frontend Structure

SvelteKit application with:
- **Routes** (`src/routes/`): File-based routing
  - `/`: Document search interface
  - `/auth/*`: Authentication flows
  - `/upload`: Document upload
  - `/account`: User and case management
- **Components** (`src/routes/lib/components/`): Reusable UI components
- **Styling**: Tailwind CSS v4

## Development Commands

### DigitalOcean Deployment (Production)
```bash
# Deploy via GitHub Actions (automatic on push to main)
git push origin main

# Manual deployment with doctl
doctl apps create --spec app.yaml
doctl apps update <app-id> --spec app.yaml

# Migrate PDFs to Spaces
python scripts/migrate_to_spaces.py --base-path ./API/data
```

### Local Development
```bash
# 1. Configure environment (OpenSearch connection)
cd API && cp ../.env.local.example .env  # Edit with OpenSearch credentials

# 2. No local Elasticsearch needed - uses DigitalOcean OpenSearch

# 2. Start API
python server.py  # Port 8000

# 3. Start Web frontend  
cd ../Web && npm install && npm run dev  # Port 5173
```

### Web Frontend Commands
```bash
cd Web
npm install
npm run dev       # Development server (port 5173)
npm run build     # Production build
npm run preview   # Preview production build
npm run check     # Type checking
npm run lint      # Linting
npm run format    # Code formatting
```

### API Backend Commands
```bash
cd API
pip install -r requirements.txt
python server.py       # Run server (port 8000)

# Alternative with uvicorn:
uvicorn server:app --reload --host 0.0.0.0 --port 8000
```

### Fiber Backend Commands (Go)
```bash
cd API  # Go Fiber code is in API directory (motion-index-fiber)

# Run main Fiber server
go run cmd/server/main.go

# Single-threaded document classification
go run cmd/api-classifier/main.go test-connection      # Test API connectivity
go run cmd/api-classifier/main.go classify-count 50   # Process 50 documents
go run cmd/api-classifier/main.go classify-all        # Process all documents

# Batch document classification (multi-threaded)
go run cmd/api-batch-classifier/main.go classify-all --limit=100

# Index management
go run cmd/setup-index/main.go     # Setup OpenSearch index
go run cmd/inspect-index/main.go   # Inspect index contents

# Build commands
go build -o bin/api-classifier cmd/api-classifier/main.go
go build -o bin/server cmd/server/main.go
```

Access local development at http://localhost:5173

## Complete Data Model

### Core Document Structure

The Motion-Index system uses a comprehensive legal document data model optimized for legal workflows and search:

#### Document Model (`pkg/models/document.go`)
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
    UpdatedAt   time.Time         `json:"updated_at"`
    Metadata    *DocumentMetadata `json:"metadata"`
    Size        int64             `json:"size,omitempty"`
    ContentType string            `json:"content_type,omitempty"`
}
```

#### Enhanced Legal Metadata (`pkg/models/document.go`)
The DocumentMetadata struct contains comprehensive legal-specific fields:

```go
type DocumentMetadata struct {
    // Basic Information
    DocumentName string       `json:"document_name"`
    Subject      string       `json:"subject"`
    Summary      string       `json:"summary,omitempty"`
    DocumentType DocumentType `json:"document_type"`

    // Case Information
    Case *CaseInfo `json:"case,omitempty"`

    // Court Information
    Court *CourtInfo `json:"court,omitempty"`

    // People & Parties
    Parties   []Party    `json:"parties,omitempty"`
    Attorneys []Attorney `json:"attorneys,omitempty"`
    Judge     *Judge     `json:"judge,omitempty"`

    // Enhanced Date Fields for Legal Documents
    FilingDate   *time.Time `json:"filing_date,omitempty"`   // When document was filed with court
    EventDate    *time.Time `json:"event_date,omitempty"`    // Key event or action date
    HearingDate  *time.Time `json:"hearing_date,omitempty"`  // Scheduled court hearing date
    DecisionDate *time.Time `json:"decision_date,omitempty"` // When court decision was made
    ServedDate   *time.Time `json:"served_date,omitempty"`   // When documents were served
    Status       string     `json:"status,omitempty"`

    // Document Properties
    Language  string `json:"language,omitempty"`
    Pages     int    `json:"pages,omitempty"`
    WordCount int    `json:"word_count,omitempty"`

    // Legal Classification
    LegalTags   []string    `json:"legal_tags,omitempty"`
    Charges     []Charge    `json:"charges,omitempty"`
    Authorities []Authority `json:"authorities,omitempty"`

    // Processing Metadata
    ProcessedAt  time.Time `json:"processed_at"`
    Confidence   float64   `json:"confidence,omitempty"`
    AIClassified bool      `json:"ai_classified"`
}
```

### Legal Entity Models (`pkg/models/legal.go`)

#### Case Information
```go
type CaseInfo struct {
    CaseNumber   string `json:"case_number"`
    CaseName     string `json:"case_name"`
    CaseType     string `json:"case_type,omitempty"`     // "criminal", "civil", "traffic"
    Chapter      string `json:"chapter,omitempty"`       // Bankruptcy chapter
    Docket       string `json:"docket,omitempty"`        // Full docket number
    NatureOfSuit string `json:"nature_of_suit,omitempty"`
}
```

#### Court Information
```go
type CourtInfo struct {
    CourtID      string `json:"court_id"`
    CourtName    string `json:"court_name"`
    Jurisdiction string `json:"jurisdiction"`            // "federal", "state", "local"
    Level        string `json:"level"`                   // "trial", "appellate", "supreme"
    District     string `json:"district,omitempty"`
    Division     string `json:"division,omitempty"`
    County       string `json:"county,omitempty"`
}
```

#### Parties and Legal Professionals
```go
type Party struct {
    Name      string     `json:"name"`
    Role      string     `json:"role"`                    // "defendant", "plaintiff", "appellant"
    PartyType string     `json:"party_type,omitempty"`    // "individual", "corporation", "government"
    Date      *time.Time `json:"date,omitempty"`
}

type Attorney struct {
    Name         string `json:"name"`
    BarNumber    string `json:"bar_number,omitempty"`
    Role         string `json:"role"`                    // "defense", "prosecution", "counsel"
    Organization string `json:"organization,omitempty"`
    ContactInfo  string `json:"contact_info,omitempty"`
}

type Judge struct {
    Name    string `json:"name"`
    Title   string `json:"title,omitempty"`
    JudgeID string `json:"judge_id,omitempty"`
}
```

#### Legal Content
```go
type Charge struct {
    Statute     string `json:"statute"`
    Description string `json:"description"`
    Grade       string `json:"grade,omitempty"`         // "felony", "misdemeanor"
    Class       string `json:"class,omitempty"`         // "A", "B", "C"
    Count       int    `json:"count,omitempty"`
}

type Authority struct {
    Citation  string `json:"citation"`
    CaseTitle string `json:"case_title,omitempty"`
    Type      string `json:"type"`                    // "case_law", "statute", "regulation"
    Precedent bool   `json:"precedent"`
    Page      string `json:"page,omitempty"`
}
```

#### Document Type Enumeration
The system supports 25+ specific legal document types with semantic categorization:

```go
type DocumentType string

// Motion Types
const (
    DocTypeMotionToSuppress         DocumentType = "motion_to_suppress"
    DocTypeMotionToDismiss          DocumentType = "motion_to_dismiss"
    DocTypeMotionToCompel           DocumentType = "motion_to_compel"
    DocTypeMotionInLimine           DocumentType = "motion_in_limine"
    DocTypeMotionForSummaryJudgment DocumentType = "motion_summary_judgment"
    // ... additional motion types
)

// Orders and Rulings
const (
    DocTypeOrder      DocumentType = "order"
    DocTypeRuling     DocumentType = "ruling"
    DocTypeJudgment   DocumentType = "judgment"
    DocTypeSentence   DocumentType = "sentence"
    DocTypeInjunction DocumentType = "injunction"
)

// Briefs and Pleadings
const (
    DocTypeBrief     DocumentType = "brief"
    DocTypeComplaint DocumentType = "complaint"
    DocTypeAnswer    DocumentType = "answer"
    DocTypePlea      DocumentType = "plea"
    DocTypeReply     DocumentType = "reply"
)

// Administrative
const (
    DocTypeDocketEntry    DocumentType = "docket_entry"
    DocTypeNotice         DocumentType = "notice"
    DocTypeStipulation    DocumentType = "stipulation"
    DocTypeCorrespondence DocumentType = "correspondence"
    DocTypeTranscript     DocumentType = "transcript"
    DocTypeEvidence       DocumentType = "evidence"
)
```

### Search and Query Models (`pkg/models/search.go`)

#### Search Request
```go
type SearchRequest struct {
    Query             string             `json:"query,omitempty"`
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
    Size              int                `json:"size"`
    From              int                `json:"from"`
    SortBy            string             `json:"sort_by,omitempty"`
    SortOrder         string             `json:"sort_order,omitempty"`
    IncludeHighlights bool               `json:"include_highlights"`
    FuzzySearch       bool               `json:"fuzzy_search"`
}
```

#### Date Range Filtering
```go
type DateRange struct {
    From *time.Time `json:"from,omitempty"`
    To   *time.Time `json:"to,omitempty"`
}
```

### API Response Models (`pkg/models/api.go`)

#### Standard API Response
```go
type APIResponse struct {
    Success   bool        `json:"success"`
    Message   string      `json:"message,omitempty"`
    Data      interface{} `json:"data,omitempty"`
    Error     *APIError   `json:"error,omitempty"`
    RequestID string      `json:"request_id,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
}

type APIError struct {
    Code    string                 `json:"code"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
    Field   string                 `json:"field,omitempty"`
}
```

## API Endpoints

### Core API Routes (`cmd/server/main.go`)

#### Health and System
- `GET /` - Root health check
- `GET /health` - Detailed system health status

#### Document Processing (`/api/v1`)
- `POST /categorise` - Upload and process documents with AI classification
- `POST /analyze-redactions` - Analyze PDF redactions and privacy data
- `POST /redact-document` - Apply redactions to PDF documents

#### Search and Discovery
- `POST /search` - Full-text and metadata search with legal filters
- `GET /legal-tags` - Available legal classification tags
- `GET /document-types` - Supported document type enumeration
- `GET /document-stats` - Index statistics and document counts
- `GET /field-options` - Available filter values for search fields
- `GET /metadata-fields` - Available metadata field names
- `GET /metadata-fields/:field` - Values for specific metadata field
- `POST /metadata-field-values` - Custom metadata field value queries

#### Document Management
- `GET /documents/:id` - Get document metadata by ID
- `GET /documents/:id/redactions` - Get redaction analysis for document
- `DELETE /documents/:id` - Delete document (currently unprotected)
- `POST /update-metadata` - Update document metadata (currently unprotected)

#### File Serving and Storage
- `GET /files/*` - Serve document files (PDF, DOCX, etc.) with embedding support
- `GET /files/search` - Find documents by filename pattern
- `GET /storage/documents` - List all documents in storage
- `GET /storage/documents/count` - Get total document count

#### Batch Processing
- `POST /batch/classify` - Start batch document classification job
- `GET /batch/:job_id/status` - Get batch job processing status
- `GET /batch/:job_id/results` - Get batch job results
- `DELETE /batch/:job_id` - Cancel batch processing job

#### Search Index Management
- `POST /index/document` - Index document directly to OpenSearch

## Environment Configuration

### Production (DigitalOcean)
Set via App Platform dashboard or GitHub Secrets:
```env
# DigitalOcean Spaces
DO_SPACES_KEY=your-spaces-key
DO_SPACES_SECRET=your-spaces-secret
DO_SPACES_BUCKET=motion-index-docs
DO_SPACES_REGION=nyc3
USE_CLOUD_STORAGE=true
STORAGE_BACKEND=spaces

# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key

# DigitalOcean Managed OpenSearch
ES_HOST=your-opensearch-host.k.db.ondigitalocean.com
ES_PORT=25060
ES_USERNAME=doadmin
ES_PASSWORD=your-opensearch-password
ES_USE_SSL=true
ES_INDEX=documents
```

### Local Development (.env)
```env
# Local storage
USE_CLOUD_STORAGE=false
STORAGE_BACKEND=local
STORAGE_PATH=./data

# DigitalOcean Managed OpenSearch (same as production)
ES_HOST=your-opensearch-host.k.db.ondigitalocean.com
ES_PORT=25060
ES_USERNAME=doadmin
ES_PASSWORD=your-opensearch-password
ES_USE_SSL=true
ES_INDEX=documents

# Supabase (same as production)
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key
```

### Web Frontend (.env.local)
```env
PUBLIC_API_URL=http://localhost:8000  # Local dev
# PUBLIC_API_URL=https://your-app.ondigitalocean.app  # Production

PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-anon-key
```

## Key Implementation Notes

- **Unified Storage**: `storage_handler.py` automatically switches between DigitalOcean Spaces (production) and local filesystem (development)
- **Document Processing**: Supports PDF, DOCX, TXT, RTF via textract with cloud/local processing
- **Authentication**: Flexible JWT verification supporting multiple algorithms via Supabase
- **Search**: Full-text search with metadata filtering via DigitalOcean Managed OpenSearch
- **File Serving**: CDN redirects for cloud storage, direct serving for local development
- **Auto-scaling**: DigitalOcean App Platform automatically scales based on traffic
- **Migration**: `migrate_to_spaces.py` handles bulk PDF upload and OpenSearch URL updates

## Testing & Quality

Currently no automated tests configured. For development:
- API: Check endpoints at http://localhost:8000/docs (Swagger)
- Frontend: Manual testing with `npm run check` for type safety
- OpenSearch: Use `python test_opensearch.py` script or OpenSearch Dashboard

## Single-Threaded Classification Script

The `cmd/api-classifier/main.go` script provides sequential document processing for debugging and controlled operations:

### Features
- **Sequential Processing**: Documents processed one at a time for easier debugging
- **Document Discovery**: Automatic discovery from DigitalOcean Spaces via storage API
- **Duplicate Detection**: Checks if documents already indexed before processing
- **Complete Pipeline**: Downloads → Classifies → Indexes in single workflow
- **Progress Tracking**: Real-time progress with detailed statistics
- **Error Handling**: Retry logic with detailed error reporting

### Commands
```bash
# Test API connectivity
go run cmd/api-classifier/main.go test-connection

# Classify first N documents (default: 10)
go run cmd/api-classifier/main.go classify-count [N]

# Classify ALL documents in storage (sequential)
go run cmd/api-classifier/main.go classify-all
```

### Configuration (Environment Variables)
```env
API_BASE_URL=http://localhost:8003          # API base URL
REQUEST_TIMEOUT=120                         # Request timeout in seconds
RETRY_ATTEMPTS=3                           # Number of retry attempts
PROCESSING_DELAY_MS=100                    # Delay between documents in milliseconds
```

### Processing Workflow
For each document:
1. **Check Existence**: Verifies if document already indexed in OpenSearch
2. **Download**: Downloads document content from DigitalOcean Spaces
3. **Process**: Sends document to processing API for text extraction and classification
4. **Index**: Directly indexes processed document to OpenSearch
5. **Progress**: Reports success/failure and updates statistics

### Key Differences from Batch Classifier
| Feature | Batch Classifier | Single-Threaded Classifier |
|---------|-----------------|---------------------------|
| Concurrency | Multi-threaded workers | Single-threaded sequential |
| Processing | Batch job API | Individual processing API |
| Queue | Background job queue | Direct processing |
| Monitoring | Job status polling | Real-time progress |
| Debugging | Complex multi-worker logs | Simple sequential logs |

### Use Cases
- **Development**: Test classification changes with controlled processing
- **Debugging**: Isolate issues with individual documents
- **Maintenance**: Process specific document sets with detailed monitoring
- **Recovery**: Re-process failed documents from batch operations
- **Testing**: Validate processing pipeline with small document sets

## Enhanced Date Extraction System

The Fiber backend now features a comprehensive date extraction system designed specifically for legal documents, enabling precise date-based search capabilities in the frontend.

### Supported Date Types

The system extracts five critical date types from legal documents:

1. **Filing Date** (`filing_date`): When the document was filed with the court
2. **Event Date** (`event_date`): Key event or action date referenced in the document  
3. **Hearing Date** (`hearing_date`): Scheduled court hearing or proceeding date
4. **Decision Date** (`decision_date`): When a court decision, ruling, or order was made
5. **Served Date** (`served_date`): When documents were served to parties

### Multi-Model AI Classification

The system supports three AI models with unified prompt architecture:

- **OpenAI GPT-4**: Production-grade classification with comprehensive legal analysis
- **Claude 3.5 Sonnet**: Advanced legal reasoning and structured extraction  
- **Ollama (Local)**: Privacy-focused local model support (Llama3, etc.)

### Key Features

- **Unified Prompts**: Centralized prompt templates ensure consistent results across all AI models
- **ISO Date Standardization**: All dates converted to YYYY-MM-DD format for reliable indexing
- **Context-Aware Validation**: Legal document date range validation (1950-present)
- **Multiple Format Support**: Handles MM/DD/YYYY, Month DD YYYY, relative dates
- **Frontend Integration**: All extracted dates are indexed in OpenSearch as searchable fields

### Classification Commands

```bash
# Single-threaded classification (debugging)
go run cmd/api-classifier/main.go test-connection       # Test API connectivity
go run cmd/api-classifier/main.go classify-count 10    # Classify 10 documents
go run cmd/api-classifier/main.go classify-all         # Classify all documents

# Batch classification (production)
go run cmd/api-batch-classifier/main.go classify-all --limit=100

# Index management  
go run cmd/setup-index/main.go                         # Setup OpenSearch index
go run cmd/inspect-index/main.go                       # Inspect index structure
```

### Date Field Integration

Enhanced date fields are now available in document metadata for frontend search:

```json
{
  "metadata": {
    "filing_date": "2024-03-15",
    "event_date": "2024-03-10", 
    "hearing_date": "2024-04-20",
    "decision_date": "2024-03-25",
    "served_date": "2024-03-12"
  }
}
```

## Performance Specs

- **Database**: 23,500+ documents indexed  
- **Search Response**: <100ms typical (DigitalOcean Managed OpenSearch)
- **Migration Speed**: 50-100 docs/sec (Python batch processing)
- **Single-threaded Processing**: 20-40 docs/min (Go sequential processing)
- **Storage**: 7.6GB for complete document set
- **CDN**: Global document delivery with edge caching
- **Auto-scaling**: Handles traffic spikes automatically