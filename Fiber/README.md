# Motion-Index Go Fiber Rewrite Plan

## Overview

This document outlines the complete plan for rebuilding the Motion-Index API from Python/FastAPI to Go/Fiber. The system is a legal document processing platform designed for California public defenders, featuring document classification, search, and compliance tools.

## Architecture Goals

- **Cloud-Native**: DigitalOcean Spaces only (no local filesystem)
- **High Performance**: Go concurrency with Fiber's zero-allocation design
- **Test-Driven**: Simple unit tests for individual functions
- **UNIX Philosophy**: Modular, single-purpose packages
- **Security-First**: JWT authentication, input validation, rate limiting

## Core Features

### 1. [[Document Processing & Classification]]
**File**: `01-document-processing-classification.md`

**Purpose**: Multi-format document processing with AI-powered legal classification
- Text extraction from PDF, DOCX, TXT, RTF with OCR support
- OpenAI ChatGPT integration for legal document classification  
- California-specific legal metadata extraction and validation
- Concurrent processing with Go worker pools

**Key Endpoints**:
- `POST /categorise` - Upload and process documents
- `POST /analyze-redactions` - PDF redaction analysis

**Dependencies**: `unidoc/unipdf`, `sashabaranov/go-openai`, `aws-sdk-go-v2`

---

### 2. [[Search & Indexing]]
**File**: `02-search-indexing.md`

**Purpose**: OpenSearch integration for full-text search with legal-specific filtering
- Advanced search with metadata filtering (court, judge, legal tags)
- Document statistics and field aggregations
- Bulk indexing operations for high throughput
- Court name normalization and legal tag validation

**Key Endpoints**:
- `POST /search` - Advanced document search
- `GET /legal-tags` - Get legal types and counts
- `GET /document-stats` - Index statistics
- `POST /update-metadata` - Update document metadata (auth required)

**Dependencies**: `opensearch-go`, `go-playground/validator`

---

### 3. [[Cloud Storage Management]]  
**File**: `03-cloud-storage-management.md`

**Purpose**: DigitalOcean Spaces integration with CDN support (no local storage)
- S3-compatible operations using AWS SDK v2
- CDN integration for global document delivery
- Secure file serving with redirect responses
- Signed URL generation for private access

**Key Endpoints**:
- `GET /api/documents/{file_path:path}` - Serve documents (CDN redirects)
- File upload handling in document processing

**Dependencies**: `aws-sdk-go-v2/service/s3`, `aws-sdk-go-v2/config`

---

### 4. [[Authentication & Security]]
**File**: `04-authentication-security.md`

**Purpose**: JWT authentication with Supabase integration and comprehensive security
- Multi-algorithm JWT validation (HS256, RS256)
- Supabase user management integration
- CORS, security headers, rate limiting
- Input validation and sanitization

**Key Middleware**:
- JWT authentication middleware
- CORS configuration for multiple domains
- Rate limiting with Redis backend
- Security headers (HSTS, CSP, XSS protection)

**Dependencies**: `golang-jwt/jwt`, `gofiber/contrib/jwt`, `supabase-go`

---

### 5. [[Document Redaction]]
**File**: `05-document-redaction.md`

**Purpose**: California legal compliance for document redaction
- Pattern-based sensitive content detection
- AI-enhanced analysis using OpenAI
- PDF manipulation for redacted copy generation
- California legal code compliance (CCP, WIC, PC)

**Key Endpoints**:
- `POST /analyze-redactions` - PDF redaction analysis 
- `POST /redact-document` - Create redacted document copy

**Dependencies**: `unidoc/unipdf`, `sashabaranov/go-openai`

---

### 6. [[API Management & Performance]]
**File**: `06-api-management-performance.md`

**Purpose**: Comprehensive API management with monitoring and performance optimization
- Health monitoring with dependency checks
- Global error handling and response formatting
- Performance metrics collection
- Graceful shutdown and recovery

**Key Endpoints**:
- `GET /` - Root status check
- `GET /health` - Health check with service status
- `GET /metrics` - Performance metrics

**Dependencies**: `gofiber/fiber/v2`, `prometheus/client_golang`, `shirou/gopsutil`

## Implementation Strategy

### Phase 1: Foundation (Week 1-2)
1. **Server Setup**: Basic Fiber application with middleware
2. **Health Monitoring**: Health checks and system metrics
3. **Cloud Storage**: DigitalOcean Spaces integration
4. **Authentication**: JWT middleware and Supabase integration

### Phase 2: Core Features (Week 3-4)
1. **Document Processing**: Text extraction and basic classification
2. **Search Engine**: OpenSearch client and basic search
3. **File Upload**: Document upload pipeline
4. **Error Handling**: Comprehensive error management

### Phase 3: Advanced Features (Week 5-6)
1. **AI Classification**: OpenAI ChatGPT integration
2. **Advanced Search**: Complex filtering and aggregations
3. **Redaction System**: Pattern detection and PDF processing
4. **Performance Optimization**: Caching, compression, monitoring

### Phase 4: Production Ready (Week 7-8)
1. **Security Hardening**: Rate limiting, validation, headers
2. **Monitoring**: Metrics, logging, alerting
3. **Testing**: Comprehensive test suite
4. **Documentation**: API documentation and deployment guides

## Project Structure

```
motion-index-fiber/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── pkg/
│   ├── processing/              # Document processing & classification
│   ├── search/                  # Search & indexing
│   ├── storage/                 # Cloud storage management  
│   ├── auth/                    # Authentication & security
│   ├── redaction/               # Document redaction
│   └── api/                     # API management & performance
├── internal/
│   ├── config/                  # Configuration management
│   ├── handlers/                # HTTP handlers
│   ├── middleware/              # Custom middleware
│   └── models/                  # Shared data models
├── test/
│   ├── unit/                    # Unit tests
│   ├── integration/             # Integration tests
│   └── testdata/                # Test fixtures
├── scripts/
│   ├── build.sh                 # Build scripts
│   ├── deploy.sh                # Deployment scripts
│   └── test.sh                  # Test scripts
├── deployments/
│   ├── docker/                  # Docker configurations
│   ├── k8s/                     # Kubernetes manifests
│   └── digitalocean/            # DigitalOcean App Platform specs
├── docs/
│   ├── api/                     # API documentation
│   ├── deployment/              # Deployment guides
│   └── development/             # Development setup
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Technology Stack

### Core Technologies
- **Language**: Go 1.21+
- **Framework**: Fiber v2 (Express-like, built on Fasthttp)
- **Database**: OpenSearch/Elasticsearch (managed)
- **Storage**: DigitalOcean Spaces (S3-compatible)
- **Authentication**: Supabase + JWT

### Key Libraries
- **Web Framework**: `github.com/gofiber/fiber/v2`
- **PDF Processing**: `github.com/unidoc/unipdf/v3`
- **Cloud Storage**: `github.com/aws/aws-sdk-go-v2`
- **Search**: `github.com/opensearch-project/opensearch-go/v2`
- **Authentication**: `github.com/golang-jwt/jwt/v5`
- **AI Integration**: `github.com/sashabaranov/go-openai`
- **Validation**: `github.com/go-playground/validator/v10`

### Development Tools
- **Testing**: Go standard library + testify
- **Linting**: golangci-lint
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker + docker-compose
- **CI/CD**: GitHub Actions

## Environment Configuration

### Required Environment Variables
```bash
# DigitalOcean Spaces
DO_SPACES_KEY=your-spaces-key
DO_SPACES_SECRET=your-spaces-secret  
DO_SPACES_BUCKET=motion-index-docs
DO_SPACES_REGION=nyc3

# OpenSearch
ES_HOST=your-opensearch-host.db.ondigitalocean.com
ES_PORT=25060
ES_USERNAME=doadmin
ES_PASSWORD=your-opensearch-password
ES_USE_SSL=true
ES_INDEX=documents

# Supabase Authentication
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key
SUPABASE_JWT_SECRET=your-jwt-secret

# AI Services
CLAUDE_API_KEY=your-claude-api-key
OPENAI_API_KEY=your-openai-api-key

# Server Configuration
PORT=8000
PRODUCTION=false
MAX_WORKERS=10
BATCH_SIZE=50
```

## Testing Strategy

### Unit Tests
- **Coverage**: >80% code coverage
- **Focus**: Individual function testing
- **Tools**: Go testing package + testify assertions
- **Patterns**: Table-driven tests for multiple scenarios

### Integration Tests  
- **Scope**: End-to-end feature workflows
- **Services**: Real OpenSearch and DigitalOcean Spaces
- **Authentication**: Test JWT validation and user management
- **Performance**: Load testing with realistic data

### Test Data
- **Legal Documents**: Sample PDFs for processing tests
- **Search Queries**: Complex legal search scenarios
- **Redaction Cases**: Documents with sensitive information
- **Error Cases**: Invalid inputs and edge cases

## Performance Targets

### Response Times
- **Document Upload**: <5 seconds for 50MB PDFs
- **Search Queries**: <100ms typical, <500ms complex
- **Health Checks**: <50ms
- **File Serving**: <200ms (CDN redirects)

### Throughput
- **Concurrent Requests**: 1000+ simultaneous
- **Document Processing**: 10+ documents/second
- **Search Operations**: 100+ queries/second
- **File Operations**: 50+ uploads/second

### Resource Usage
- **Memory**: <1GB typical, <2GB peak
- **CPU**: <50% utilization under normal load
- **Storage**: Efficient temporary file management
- **Network**: Optimized for CDN delivery

## Security Considerations

### Authentication & Authorization
- JWT tokens with proper expiration
- Role-based access control
- Supabase user management integration
- Rate limiting per user/IP

### Data Protection
- TLS encryption for all communications
- Secure file uploads with validation
- PDF redaction for sensitive content
- Audit logging for document access

### Input Validation
- Strict file type and size validation
- SQL injection prevention (parameterized queries)
- XSS protection through proper encoding
- Path traversal prevention

### Infrastructure Security
- Private network configuration
- Firewall rules for database access
- Environment variable management
- Secure container deployment

## Deployment Strategy

### DigitalOcean App Platform
- **Primary Deployment**: Auto-scaling cloud platform
- **Benefits**: Managed infrastructure, auto-scaling, zero-downtime deployments
- **Configuration**: App spec with environment variables
- **Monitoring**: Built-in metrics and logging

### Container Deployment
- **Alternative**: Docker containers on Droplets
- **Orchestration**: Docker Compose or Kubernetes
- **Scaling**: Manual scaling with load balancer
- **Monitoring**: Custom metrics collection

### CI/CD Pipeline
- **Source Control**: GitHub integration
- **Build Process**: Automated Go builds
- **Testing**: Automated test suite execution
- **Deployment**: Automatic deployment on merge to main

## Migration Plan

### Data Migration
1. **Document Storage**: Migrate PDFs to DigitalOcean Spaces
2. **Search Index**: Re-index documents in OpenSearch
3. **User Data**: Validate Supabase user management
4. **Configuration**: Environment variable setup

### Rollout Strategy
1. **Parallel Deployment**: Run new system alongside existing
2. **Feature Toggle**: Gradual feature migration
3. **User Migration**: Migrate users in batches
4. **Monitoring**: Compare performance metrics
5. **Cutover**: Complete migration after validation

### Rollback Plan
- Maintain Python system during transition
- Database snapshots before migration
- Quick rollback procedures documented
- Health monitoring for early issue detection

## Success Metrics

### Performance Improvements
- **Response Time**: 50% faster than Python version
- **Memory Usage**: 60% less memory consumption
- **CPU Efficiency**: Better multi-core utilization
- **Concurrent Users**: 5x more concurrent capacity

### Reliability Improvements
- **Uptime**: 99.9% availability target
- **Error Rate**: <0.1% error rate
- **Recovery Time**: <30 seconds for service recovery
- **Data Integrity**: Zero data loss during operations

### Developer Experience
- **Build Time**: <30 seconds for full build
- **Test Suite**: <60 seconds for full test run
- **Hot Reload**: <3 seconds for development changes
- **Deployment**: <5 minutes for production deployment

This comprehensive plan provides the foundation for successfully rebuilding the Motion-Index API using Go and Fiber, with a focus on performance, security, and maintainability while leveraging cloud-native technologies.