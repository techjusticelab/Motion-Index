# API Management & Performance

## Overview
This feature provides comprehensive API management including health monitoring, error handling, performance optimization, and system diagnostics for the Motion-Index Go Fiber application.

## Current Python Implementation Analysis

### Key Components (from API analysis):
- **Health endpoint** in `server.py`: Basic OpenSearch connectivity check
- **Error handling**: Try-catch with HTTPException responses
- **CORS configuration**: Multiple domain support
- **Logging**: Basic Python logging with configurable levels
- **Demo mode**: Graceful degradation when OpenSearch unavailable

### Endpoints from `server.py`:
- `GET /` - Root endpoint status check
- `GET /health` - Health check with OpenSearch status
- Error responses with proper HTTP status codes
- CORS preflight handling

## Go Package Design

### Package Structure:
```
pkg/
├── api/
│   ├── health/              # Health monitoring and diagnostics
│   │   ├── checker.go       # Health check implementation
│   │   ├── monitor.go       # System monitoring
│   │   ├── metrics.go       # Performance metrics
│   │   └── status.go        # Service status tracking
│   ├── middleware/          # API middleware
│   │   ├── logging.go       # Request/response logging
│   │   ├── metrics.go       # Performance metrics collection
│   │   ├── recovery.go      # Panic recovery
│   │   ├── timeout.go       # Request timeout handling
│   │   └── compression.go   # Response compression
│   ├── errors/              # Error handling
│   │   ├── handler.go       # Global error handler
│   │   ├── types.go         # Error types and codes
│   │   ├── formatter.go     # Error response formatting
│   │   └── recovery.go      # Error recovery strategies
│   ├── response/            # Response formatting
│   │   ├── formatter.go     # Standard response formats
│   │   ├── pagination.go    # Paginated response handling
│   │   ├── headers.go       # Response headers management
│   │   └── cache.go         # Response caching
│   └── server/              # Server management
│       ├── server.go        # Main server setup
│       ├── routes.go        # Route registration
│       ├── config.go        # Server configuration
│       ├── graceful.go      # Graceful shutdown
│       └── telemetry.go     # Telemetry and monitoring
```

### Core Interfaces:

```go
// HealthChecker interface for service health monitoring
type HealthChecker interface {
    CheckHealth(ctx context.Context) (*HealthStatus, error)
    CheckDependencies(ctx context.Context) ([]*DependencyStatus, error)
    GetSystemMetrics(ctx context.Context) (*SystemMetrics, error)
    IsHealthy(ctx context.Context) bool
}

// MetricsCollector interface for performance metrics
type MetricsCollector interface {
    RecordRequest(method, path string, statusCode int, duration time.Duration)
    RecordError(err error, context string)
    GetMetrics(ctx context.Context) (*Metrics, error)
    IncrementCounter(name string, tags map[string]string)
    RecordHistogram(name string, value float64, tags map[string]string)
}

// ErrorHandler interface for error management
type ErrorHandler interface {
    HandleError(c *fiber.Ctx, err error) error
    FormatError(err error) *ErrorResponse
    RecordError(err error, context map[string]interface{})
    ShouldRetry(err error) bool
}

// ResponseFormatter interface for consistent API responses
type ResponseFormatter interface {
    Success(data interface{}) *APIResponse
    Error(code int, message string, details interface{}) *ErrorResponse
    Paginated(data interface{}, pagination *PaginationInfo) *PaginatedResponse
    WithMetadata(data interface{}, metadata map[string]interface{}) *APIResponse
}
```

### Data Models:

```go
type HealthStatus struct {
    Status      string                 `json:"status"`      // "healthy", "degraded", "unhealthy"
    Timestamp   time.Time              `json:"timestamp"`
    Uptime      time.Duration          `json:"uptime"`
    Version     string                 `json:"version"`
    Environment string                 `json:"environment"`
    Services    map[string]ServiceStatus `json:"services"`
    System      *SystemInfo            `json:"system"`
}

type ServiceStatus struct {
    Name      string        `json:"name"`
    Status    string        `json:"status"`    // "up", "down", "degraded"
    Latency   time.Duration `json:"latency"`
    LastCheck time.Time     `json:"last_check"`
    Error     string        `json:"error,omitempty"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type SystemInfo struct {
    OS           string  `json:"os"`
    Architecture string  `json:"architecture"`
    CPUs         int     `json:"cpus"`
    Memory       *Memory `json:"memory"`
    Disk         *Disk   `json:"disk"`
    GoVersion    string  `json:"go_version"`
    Goroutines   int     `json:"goroutines"`
}

type Memory struct {
    Total     uint64  `json:"total"`
    Available uint64  `json:"available"`
    Used      uint64  `json:"used"`
    UsedMB    float64 `json:"used_mb"`
    Percent   float64 `json:"percent"`
}

type Disk struct {
    Total     uint64  `json:"total"`
    Free      uint64  `json:"free"`
    Used      uint64  `json:"used"`
    UsedGB    float64 `json:"used_gb"`
    Percent   float64 `json:"percent"`
}

type APIResponse struct {
    Success   bool                   `json:"success"`
    Data      interface{}            `json:"data,omitempty"`
    Message   string                 `json:"message,omitempty"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
}

type ErrorResponse struct {
    Success   bool                   `json:"success"`
    Error     *ErrorInfo             `json:"error"`
    RequestID string                 `json:"request_id,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
}

type ErrorInfo struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
    Type    string      `json:"type"`
}

type PaginatedResponse struct {
    Success    bool            `json:"success"`
    Data       interface{}     `json:"data"`
    Pagination *PaginationInfo `json:"pagination"`
    Timestamp  time.Time       `json:"timestamp"`
}

type PaginationInfo struct {
    Page       int `json:"page"`
    PerPage    int `json:"per_page"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
    HasNext    bool `json:"has_next"`
    HasPrev    bool `json:"has_prev"`
}

type Metrics struct {
    Requests    *RequestMetrics    `json:"requests"`
    Performance *PerformanceMetrics `json:"performance"`
    Errors      *ErrorMetrics      `json:"errors"`
    System      *SystemMetrics     `json:"system"`
    Timestamp   time.Time          `json:"timestamp"`
}

type RequestMetrics struct {
    Total            int64                      `json:"total"`
    TotalToday       int64                      `json:"total_today"`
    Rate             float64                    `json:"rate_per_second"`
    ByMethod         map[string]int64           `json:"by_method"`
    ByStatus         map[string]int64           `json:"by_status"`
    ByEndpoint       map[string]int64           `json:"by_endpoint"`
    AverageLatency   time.Duration              `json:"average_latency"`
}
```

## Fiber Server Implementation

### Main Server Setup:
```go
func NewServer(config *ServerConfig) *Server {
    // Create Fiber app with optimized config
    app := fiber.New(fiber.Config{
        AppName:               "Motion-Index API",
        ServerHeader:          "Motion-Index",
        DisableStartupMessage: config.Production,
        ErrorHandler:          NewGlobalErrorHandler(),
        BodyLimit:             100 * 1024 * 1024, // 100MB for large PDFs
        ReadTimeout:           30 * time.Second,
        WriteTimeout:          30 * time.Second,
        IdleTimeout:           60 * time.Second,
        Immutable:             true, // Enable for zero allocation
        UnescapePath:          true,
        CaseSensitive:         false,
        StrictRouting:         false,
        EnableTrustedProxyCheck: true,
        TrustedProxies:        config.TrustedProxies,
    })

    server := &Server{
        app:    app,
        config: config,
        logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
    }

    // Setup middleware
    server.setupMiddleware()
    
    // Setup routes
    server.setupRoutes()
    
    return server
}

func (s *Server) setupMiddleware() {
    // Recovery middleware (first)
    s.app.Use(recover.New(recover.Config{
        EnableStackTrace: !s.config.Production,
    }))
    
    // Request ID middleware
    s.app.Use(requestid.New())
    
    // Logging middleware
    s.app.Use(logger.New(logger.Config{
        Format: "[${time}] ${status} - ${latency} ${method} ${path} ${error}\n",
        TimeFormat: "2006-01-02 15:04:05",
        TimeZone:   "UTC",
    }))
    
    // CORS middleware
    s.app.Use(NewCORSMiddleware())
    
    // Security headers middleware
    s.app.Use(NewSecurityMiddleware())
    
    // Compression middleware
    s.app.Use(compress.New(compress.Config{
        Level: compress.LevelBestSpeed,
    }))
    
    // Rate limiting middleware
    s.app.Use(NewRateLimitMiddleware())
    
    // Metrics collection middleware
    s.app.Use(NewMetricsMiddleware())
}
```

### Health Check Handlers:
```go
func (h *HealthHandler) GetHealth(c *fiber.Ctx) error {
    status, err := h.healthChecker.CheckHealth(c.Context())
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    // Determine HTTP status based on health
    httpStatus := fiber.StatusOK
    switch status.Status {
    case "degraded":
        httpStatus = fiber.StatusPartialContent
    case "unhealthy":
        httpStatus = fiber.StatusServiceUnavailable
    }
    
    return c.Status(httpStatus).JSON(status)
}

func (h *HealthHandler) GetReadiness(c *fiber.Ctx) error {
    // Check if all critical services are available
    deps, err := h.healthChecker.CheckDependencies(c.Context())
    if err != nil {
        return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
            "ready": false,
            "error": err.Error(),
        })
    }
    
    ready := true
    for _, dep := range deps {
        if dep.Status != "up" && dep.Required {
            ready = false
            break
        }
    }
    
    status := fiber.StatusOK
    if !ready {
        status = fiber.StatusServiceUnavailable
    }
    
    return c.Status(status).JSON(fiber.Map{
        "ready":        ready,
        "dependencies": deps,
        "timestamp":    time.Now(),
    })
}

func (h *HealthHandler) GetLiveness(c *fiber.Ctx) error {
    // Simple liveness check - server is running
    return c.JSON(fiber.Map{
        "alive":     true,
        "timestamp": time.Now(),
        "uptime":    time.Since(h.startTime),
    })
}

func (h *HealthHandler) GetMetrics(c *fiber.Ctx) error {
    metrics, err := h.metricsCollector.GetMetrics(c.Context())
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    return c.JSON(metrics)
}
```

### Error Handler Implementation:
```go
func NewGlobalErrorHandler() fiber.ErrorHandler {
    return func(c *fiber.Ctx, err error) error {
        // Default error response
        code := fiber.StatusInternalServerError
        message := "Internal Server Error"
        var details interface{}
        
        // Handle Fiber errors
        if e, ok := err.(*fiber.Error); ok {
            code = e.Code
            message = e.Message
        }
        
        // Handle validation errors
        if ve, ok := err.(validator.ValidationErrors); ok {
            code = fiber.StatusBadRequest
            message = "Validation failed"
            details = formatValidationErrors(ve)
        }
        
        // Handle context errors
        if err == context.DeadlineExceeded {
            code = fiber.StatusRequestTimeout
            message = "Request timeout"
        }
        
        // Handle database/service errors
        if isServiceUnavailable(err) {
            code = fiber.StatusServiceUnavailable
            message = "Service temporarily unavailable"
        }
        
        // Generate request ID if not present
        requestID := c.Locals("requestid")
        if requestID == nil {
            requestID = generateRequestID()
        }
        
        // Log error
        logError(c, err, code, requestID)
        
        // Return JSON error response
        return c.Status(code).JSON(&ErrorResponse{
            Success:   false,
            Error: &ErrorInfo{
                Code:    getErrorCode(code),
                Message: message,
                Details: details,
                Type:    getErrorType(err),
            },
            RequestID: requestID.(string),
            Timestamp: time.Now(),
        })
    }
}

func formatValidationErrors(errs validator.ValidationErrors) map[string]string {
    result := make(map[string]string)
    for _, err := range errs {
        result[err.Field()] = fmt.Sprintf("Field validation failed on '%s' tag", err.Tag())
    }
    return result
}

func getErrorCode(status int) string {
    switch status {
    case fiber.StatusBadRequest:
        return "BAD_REQUEST"
    case fiber.StatusUnauthorized:
        return "UNAUTHORIZED"
    case fiber.StatusForbidden:
        return "FORBIDDEN"
    case fiber.StatusNotFound:
        return "NOT_FOUND"
    case fiber.StatusTooManyRequests:
        return "RATE_LIMITED"
    case fiber.StatusInternalServerError:
        return "INTERNAL_ERROR"
    case fiber.StatusServiceUnavailable:
        return "SERVICE_UNAVAILABLE"
    default:
        return "UNKNOWN_ERROR"
    }
}
```

### Metrics Collection Middleware:
```go
func NewMetricsMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        // Process request
        err := c.Next()
        
        // Calculate duration
        duration := time.Since(start)
        
        // Record metrics
        method := c.Method()
        path := c.Route().Path
        statusCode := c.Response().StatusCode()
        
        // Record request metrics
        recordRequestMetrics(method, path, statusCode, duration)
        
        // Record error metrics if needed
        if err != nil {
            recordErrorMetrics(err, method, path)
        }
        
        return err
    }
}

func recordRequestMetrics(method, path string, statusCode int, duration time.Duration) {
    // Update request counters
    requestTotal.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
    
    // Update latency histogram
    requestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
    
    // Update active requests gauge
    activeRequests.Inc()
    defer activeRequests.Dec()
}
```

## Performance Optimization

### Connection Management:
```go
type ServerConfig struct {
    Port            int           `env:"PORT" default:"8000"`
    Host            string        `env:"HOST" default:"0.0.0.0"`
    Production      bool          `env:"PRODUCTION" default:"false"`
    MaxConnections  int           `env:"MAX_CONNECTIONS" default:"1000"`
    ReadTimeout     time.Duration `env:"READ_TIMEOUT" default:"30s"`
    WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" default:"30s"`
    IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" default:"60s"`
    BodyLimit       int           `env:"BODY_LIMIT" default:"104857600"` // 100MB
    TrustedProxies  []string      `env:"TRUSTED_PROXIES"`
    
    // Performance settings
    Prefork         bool `env:"PREFORK" default:"false"`
    WorkerCount     int  `env:"WORKER_COUNT" default:"0"` // 0 = auto
    EnableProfiling bool `env:"ENABLE_PROFILING" default:"false"`
}
```

### Graceful Shutdown:
```go
func (s *Server) Start() error {
    // Start server in goroutine
    go func() {
        addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
        s.logger.Info("Starting server", "address", addr)
        
        if err := s.app.Listen(addr); err != nil && err != http.ErrServerClosed {
            s.logger.Error("Server error", "error", err)
        }
    }()
    
    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    s.logger.Info("Shutting down server...")
    
    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    return s.app.ShutdownWithContext(ctx)
}
```

### Health Check Implementation:
```go
type HealthService struct {
    openSearchClient *opensearch.Client
    storageService   storage.StorageService
    startTime        time.Time
}

func (h *HealthService) CheckHealth(ctx context.Context) (*HealthStatus, error) {
    status := &HealthStatus{
        Timestamp:   time.Now(),
        Uptime:      time.Since(h.startTime),
        Version:     getVersion(),
        Environment: getEnvironment(),
        Services:    make(map[string]ServiceStatus),
        System:      h.getSystemInfo(),
    }
    
    // Check OpenSearch
    osStatus := h.checkOpenSearch(ctx)
    status.Services["opensearch"] = osStatus
    
    // Check Storage
    storageStatus := h.checkStorage(ctx)
    status.Services["storage"] = storageStatus
    
    // Determine overall status
    status.Status = h.calculateOverallStatus(status.Services)
    
    return status, nil
}

func (h *HealthService) checkOpenSearch(ctx context.Context) ServiceStatus {
    start := time.Now()
    
    if h.openSearchClient == nil {
        return ServiceStatus{
            Name:      "opensearch",
            Status:    "down",
            Latency:   0,
            LastCheck: start,
            Error:     "client not initialized",
        }
    }
    
    // Ping OpenSearch
    res, err := h.openSearchClient.Ping()
    latency := time.Since(start)
    
    if err != nil {
        return ServiceStatus{
            Name:      "opensearch",
            Status:    "down",
            Latency:   latency,
            LastCheck: start,
            Error:     err.Error(),
        }
    }
    
    if res.IsError() {
        return ServiceStatus{
            Name:      "opensearch",
            Status:    "down",
            Latency:   latency,
            LastCheck: start,
            Error:     fmt.Sprintf("HTTP %d", res.StatusCode),
        }
    }
    
    return ServiceStatus{
        Name:      "opensearch",
        Status:    "up",
        Latency:   latency,
        LastCheck: start,
    }
}
```

## Test Strategy

### Unit Tests:
```go
func TestHealthHandler_GetHealth(t *testing.T) {
    tests := []struct {
        name           string
        healthStatus   *HealthStatus
        expectedStatus int
    }{
        {
            name: "healthy system",
            healthStatus: &HealthStatus{
                Status: "healthy",
            },
            expectedStatus: 200,
        },
        {
            name: "degraded system",
            healthStatus: &HealthStatus{
                Status: "degraded",
            },
            expectedStatus: 206,
        },
        {
            name: "unhealthy system",
            healthStatus: &HealthStatus{
                Status: "unhealthy",
            },
            expectedStatus: 503,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests:
- Health endpoint functionality
- Error handling across different scenarios
- Performance metrics collection
- Graceful shutdown behavior

## Implementation Priority

1. **Basic Server Setup** - Fiber app with middleware
2. **Health Monitoring** - Health checks and system metrics
3. **Error Handling** - Global error handler and response formatting
4. **Performance Monitoring** - Metrics collection and reporting
5. **Graceful Shutdown** - Clean server shutdown handling
6. **Advanced Features** - Performance optimization and monitoring

## Dependencies

### External Libraries:
- `github.com/gofiber/fiber/v2` - Web framework
- `github.com/gofiber/fiber/v2/middleware/*` - Fiber middleware
- `github.com/prometheus/client_golang` - Metrics collection
- `github.com/shirou/gopsutil/v3` - System metrics

### Configuration:
```go
type ServerConfig struct {
    // Server Configuration
    Port           int           `env:"PORT" default:"8000"`
    Host           string        `env:"HOST" default:"0.0.0.0"`
    Production     bool          `env:"PRODUCTION" default:"false"`
    
    // Performance Configuration
    MaxConnections int           `env:"MAX_CONNECTIONS" default:"1000"`
    ReadTimeout    time.Duration `env:"READ_TIMEOUT" default:"30s"`
    WriteTimeout   time.Duration `env:"WRITE_TIMEOUT" default:"30s"`
    IdleTimeout    time.Duration `env:"IDLE_TIMEOUT" default:"60s"`
    BodyLimit      int           `env:"BODY_LIMIT" default:"104857600"`
    
    // Monitoring Configuration
    MetricsEnabled bool   `env:"METRICS_ENABLED" default:"true"`
    HealthPath     string `env:"HEALTH_PATH" default:"/health"`
    MetricsPath    string `env:"METRICS_PATH" default:"/metrics"`
    
    // Security Configuration
    TrustedProxies []string `env:"TRUSTED_PROXIES"`
    EnableProfiling bool    `env:"ENABLE_PROFILING" default:"false"`
}
```

## Performance Considerations

- **Zero Allocation**: Use Fiber's immutable mode for performance
- **Connection Pooling**: Efficient HTTP client management
- **Compression**: Response compression for large payloads
- **Caching**: Response caching for static data
- **Monitoring**: Low-overhead metrics collection

## Security Considerations

- **Error Information**: Limit error details in production
- **Rate Limiting**: Protect against abuse
- **Request Validation**: Comprehensive input validation
- **Security Headers**: Complete security header implementation
- **Logging**: Secure logging without sensitive data exposure