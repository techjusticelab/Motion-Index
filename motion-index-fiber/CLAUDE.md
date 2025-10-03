# CLAUDE.md

This file provides comprehensive guidance to Claude Code (claude.ai/code) when working with the Motion-Index Fiber codebase.

## Project Overview

Motion-Index Fiber is a high-performance legal document processing API built with Go and Fiber, designed for California public defenders. The system features production-ready DigitalOcean integration with direct REST API calls and AWS S3 SDK, following strict TDD and UNIX philosophy principles.

### Core Technologies
- **Runtime**: Go 1.21+ with Fiber v2 web framework
- **Cloud Platform**: DigitalOcean (Spaces storage, Managed OpenSearch, App Platform)
- **Integration Pattern**: Direct DigitalOcean REST API + AWS S3 SDK for Spaces
- **Authentication**: JWT with Supabase integration
- **AI Processing**: OpenAI GPT-4 for document classification
- **Testing**: Comprehensive TDD with testify framework

## Architecture

### High-Level Design
```
┌─────────────────┐
│   Web Client    │
└─────────┬───────┘
          │
    ┌─────▼─────┐
    │Motion Index│
    │    API     │
    │ (Fiber v2) │
    └─────┬─────┘
          │
┌─────────┼─────────┐
│    DigitalOcean   │    │   Supabase    │    │   OpenAI    │
│ (REST API + S3)   │    │     (JWT)     │    │   (GPT-4)   │
└─────────┬─────────┘    └───────────────┘    └─────────────┘
          │
    ┌─────▼─────┐        ┌─────────────┐
    │  Spaces   │        │ OpenSearch  │
    │Storage+CDN│        │   Cluster   │
    └───────────┘        └─────────────┘
```

### Package Structure
- **`cmd/server/`**: Application entry point with graceful shutdown
- **`pkg/cloud/digitalocean/`**: Direct DigitalOcean API integration and service factory
- **`pkg/processing/`**: Document processing pipeline (extract, classify, process)
- **`pkg/search/`**: Search interfaces and OpenSearch client
- **`pkg/storage/`**: Storage interfaces and utilities
- **`internal/config/`**: Application configuration with validation
- **`internal/handlers/`**: HTTP request handlers
- **`internal/middleware/`**: Custom middleware (auth, error handling)
- **`internal/models/`**: Data models and validation
- **`internal/testutil/`**: Test utilities following UNIX principles

## Development Commands

### Core Development Workflow
```bash
# Initial setup
go mod tidy
cp .env.example .env  # Edit with your configuration

# Development server
go run cmd/server/main.go

# Or with auto-reload (install air first)
go install github.com/cosmtrek/air@latest
air

# Build for development
go build -o bin/server cmd/server/main.go

# Build for production
go build -ldflags="-s -w" -o bin/server cmd/server/main.go
```

### Testing Commands (Critical)
This project follows strict TDD principles. ALWAYS run tests before and after code changes:

```bash
# Run all tests with coverage (primary command)
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out

# Unit tests only (fast feedback loop)
go test ./... -short -v

# Integration tests (requires real DigitalOcean credentials)
RUN_INTEGRATION_TESTS=true go test ./... -v -tags=integration

# Test specific packages
go test ./internal/config/... -v
go test ./pkg/cloud/digitalocean/... -v
go test ./internal/handlers/... -v

# Benchmark tests
go test ./... -bench=. -benchmem

# Race condition detection
go test ./... -race -v

# Test coverage by package
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Code Quality Commands
```bash
# Format code
go fmt ./...

# Vet for common issues
go vet ./...

# Lint (if golangci-lint is installed)
golangci-lint run

# Dependency verification
go mod verify
go mod download
```

## Configuration Management

### Environment Files
- **`.env.example`**: Template with all required variables
- **`.env.local.example`**: Local development configuration
- **`.env.production.example`**: Production deployment configuration

### Key Configuration Sections

#### Core Application
```bash
PORT=6000
ENVIRONMENT=local  # local, staging, production
JWT_SECRET=your-jwt-secret
```

#### DigitalOcean Services (Production)
```bash
# DigitalOcean API Token for CDN/Management operations
DO_API_TOKEN=your-digitalocean-api-token

# Spaces Storage
DO_SPACES_ACCESS_KEY=your-access-key
DO_SPACES_SECRET_KEY=your-secret-key
DO_SPACES_BUCKET=motion-index-docs
DO_SPACES_REGION=nyc3

# Managed OpenSearch
DO_OPENSEARCH_HOST=your-cluster.k.db.ondigitalocean.com
DO_OPENSEARCH_PORT=25060
DO_OPENSEARCH_USERNAME=doadmin
DO_OPENSEARCH_PASSWORD=your-password
DO_OPENSEARCH_USE_SSL=true
```

#### External Services
```bash
# Supabase Authentication
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key

# OpenAI Classification
OPENAI_API_KEY=your-api-key
OPENAI_MODEL=gpt-4
```

## UNIX Philosophy Adherence

### Core Principles
1. **Do One Thing Well**: Each component has a single, focused responsibility
2. **Composable**: Services can be combined and used independently
3. **Testable**: Every component has comprehensive test coverage
4. **Observable**: Health checks, metrics, and logging throughout

### Implementation Guidelines

#### Package Design
- Each package should have a clear, single purpose
- Interfaces should be minimal and focused
- Dependencies should be explicit and injectable
- Avoid global state and singleton patterns

#### Function Design
- Functions should be pure when possible
- Single responsibility and clear naming
- Comprehensive error handling
- Input validation at boundaries

#### Testing Strategy
- Unit tests for isolated logic
- Integration tests for service interactions
- Benchmark tests for performance validation
- Mock external dependencies consistently

## Test-Driven Development (TDD)

### TDD Workflow
1. **Red**: Write a failing test first
2. **Green**: Write minimal code to make test pass
3. **Refactor**: Improve code while keeping tests green
4. **Repeat**: Continue with next requirement

### Test Categories

#### Unit Tests (`*_test.go`)
- Fast, isolated tests with no external dependencies
- Use mocks for external services
- Test edge cases and error conditions
- Example: `internal/config/config_test.go`

#### Integration Tests (`integration_test.go`)
- Test interactions between components
- Use real services when possible (with proper setup)
- Mark with `// +build integration` tag
- Example: `pkg/cloud/digitalocean/integration_test.go`

#### Benchmark Tests
- Performance validation and optimization
- Memory allocation tracking
- Comparison between implementations
- Run with `go test -bench=. -benchmem`

### Test Utilities (`internal/testutil/`)
Following UNIX principles, test utilities are composable and single-purpose:

```go
// Environment setup
cleanup := testutil.SetEnv(t, map[string]string{
    "PORT": "8080",
    "JWT_SECRET": "test-secret",
})
defer cleanup()

// Configuration testing
cfg := testutil.TestConfig()
cfg.DigitalOcean = testutil.TestDigitalOceanConfig()

// Temporary files and directories
tmpDir, cleanup := testutil.TempDir(t)
defer cleanup()
```

## API Development Patterns

### Handler Structure
```go
type SomeHandler struct {
    service SomeService
    logger  Logger
}

func NewSomeHandler(service SomeService) *SomeHandler {
    return &SomeHandler{service: service}
}

func (h *SomeHandler) HandleRequest(c *fiber.Ctx) error {
    // 1. Parse and validate input
    // 2. Call business logic (service layer)
    // 3. Handle errors appropriately
    // 4. Return structured response
}
```

### Error Handling
```go
// Use structured errors with context
if err != nil {
    return fiber.NewError(fiber.StatusBadRequest, 
        fmt.Sprintf("failed to process document: %v", err))
}

// Or use custom error middleware
return models.NewErrorResponse("VALIDATION_ERROR", 
    "Invalid document format", details)
```

### Response Structure
```go
// Success responses
return c.JSON(models.NewSuccessResponse(data, "Operation completed"))

// Error responses  
return c.JSON(models.NewErrorResponse("ERROR_CODE", "Message", details))
```

## Service Integration Patterns

### DigitalOcean Services
```go
// Use service factory pattern
provider, err := digitalocean.NewProviderFromEnvironment()
if err != nil {
    return fmt.Errorf("failed to create provider: %w", err)
}

// Initialize services
err = provider.Initialize()
if err != nil {
    return fmt.Errorf("failed to initialize: %w", err)
}

// Get specific services
services := provider.GetServices()
storage := services.Storage
search := services.Search
```

### Health Checks
```go
// Implement health check interface
func (s *SomeService) IsHealthy() bool {
    // Implement actual health check logic
    return s.client != nil && s.ping()
}

// Use in handlers
func (h *HealthHandler) DetailedStatus(c *fiber.Ctx) error {
    status := &models.SystemStatus{
        Service:   "motion-index-fiber",
        Status:    "healthy",
        Storage:   h.getStorageStatus(),
        Search:    h.getSearchStatus(),
    }
    return c.JSON(models.NewSuccessResponse(status, "System status"))
}
```

## Deployment & Production

### Local Development
1. Configure `.env` with local/staging credentials
2. Use `ENVIRONMENT=local` for mock services
3. Run with `go run cmd/server/main.go`
4. Access at `http://localhost:6000`

### DigitalOcean App Platform
1. Push to GitHub repository
2. Configure environment variables in dashboard
3. Deploy with `doctl apps create --spec deployments/digitalocean/app.yaml`
4. Monitor with health checks and logs

### Docker Deployment
```bash
# Build image
docker build -f deployments/docker/Dockerfile.prod -t motion-index:latest .

# Run container
docker run -d --name motion-index -p 6000:6000 --env-file .env.production motion-index:latest
```

## Code Style & Conventions

### Go Conventions
- Follow standard Go formatting (`go fmt`)
- Use meaningful package and variable names
- Implement interfaces explicitly
- Handle errors explicitly (no silent failures)
- Use context.Context for cancellation and timeouts

### Project Conventions
- Configuration through environment variables
- Structured logging with context
- Comprehensive error handling with wrapped errors
- Health checks for all external dependencies
- Metrics collection for observability

## Troubleshooting Common Issues

### Configuration Issues
- Verify all required environment variables are set
- Check configuration validation in `internal/config/config_test.go`
- Use `go test ./internal/config/... -v` to test configuration

### Service Connection Issues
- Check DigitalOcean service status and credentials
- Verify network connectivity to external services
- Use health check endpoints to diagnose issues
- Check logs for specific error messages

### Test Failures
- Ensure test environment is clean (no conflicting env vars)
- For integration tests, verify real service credentials
- Use `-v` flag for detailed test output
- Check test coverage for missing test cases

### Performance Issues
- Use benchmark tests to identify bottlenecks
- Check memory allocation with `-benchmem`
- Profile with `go tool pprof` if needed
- Monitor health check metrics

## Important Implementation Notes

- **Never commit secrets**: Use environment variables for all sensitive data
- **Always test**: Follow TDD principles and maintain high test coverage
- **Handle errors**: Wrap errors with context and return meaningful messages
- **Document APIs**: Use clear endpoint documentation and examples
- **Monitor health**: Implement comprehensive health checks for all services
- **Follow UNIX**: Keep components small, focused, and composable
- **Validate input**: Sanitize and validate all user input at API boundaries
- **Use interfaces**: Design with interfaces for testability and flexibility

## Quick Reference Commands

```bash
# Start development server
go run cmd/server/main.go

# Run tests with coverage
go test ./... -v -coverprofile=coverage.out

# Build production binary
go build -ldflags="-s -w" -o bin/server cmd/server/main.go

# Test specific package
go test ./internal/handlers/... -v

# Check configuration
go test ./internal/config/... -v

# Benchmark performance
go test ./... -bench=. -benchmem

# Integration tests (requires credentials)
RUN_INTEGRATION_TESTS=true go test ./... -v -tags=integration
```