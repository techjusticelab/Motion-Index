# Cloud Storage Management

## Overview
This feature provides DigitalOcean Spaces integration for document storage with CDN support, removing all local filesystem dependencies for a pure cloud-native architecture.

## Current Python Implementation Analysis

### Key Components (from API analysis):
- **`src/handlers/storage_handler.py`**: Unified storage with cloud/local switching
- **`src/handlers/s3_handler.py`**: S3-compatible operations (currently used for DigitalOcean Spaces)

### Endpoints from `server.py`:
- `GET /api/documents/{file_path:path}` - Serve documents (redirects to CDN for cloud storage)

### Current Logic:
- Checks `USE_CLOUD_STORAGE` environment variable
- For cloud: Redirects to CDN URLs with signed access
- For local: Direct file serving with path traversal protection

## Go Package Design (Cloud-Only)

### Package Structure:
```
pkg/
├── storage/
│   ├── spaces/              # DigitalOcean Spaces integration
│   │   ├── client.go        # S3-compatible client for Spaces
│   │   ├── upload.go        # File upload operations
│   │   ├── download.go      # File access and URL generation
│   │   ├── delete.go        # File deletion operations
│   │   └── config.go        # Spaces configuration
│   ├── cdn/                 # CDN integration
│   │   ├── urls.go          # CDN URL generation
│   │   ├── cache.go         # Cache control headers
│   │   └── signed.go        # Signed URL generation
│   ├── models/              # Storage data models
│   │   ├── file.go          # File metadata
│   │   ├── upload.go        # Upload request/response
│   │   └── access.go        # Access control models
│   └── interface.go         # Storage service interface
```

### Core Interfaces:

```go
// StorageService interface for cloud storage operations
type StorageService interface {
    // File Operations
    UploadFile(ctx context.Context, req *UploadRequest) (*UploadResult, error)
    DeleteFile(ctx context.Context, key string) error
    FileExists(ctx context.Context, key string) (bool, error)
    GetFileMetadata(ctx context.Context, key string) (*FileMetadata, error)
    
    // Access Operations
    GenerateAccessURL(ctx context.Context, key string, expiry time.Duration) (string, error)
    GenerateCDNURL(ctx context.Context, key string) (string, error)
    GenerateSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
    
    // Batch Operations
    UploadFiles(ctx context.Context, requests []*UploadRequest) ([]*UploadResult, error)
    DeleteFiles(ctx context.Context, keys []string) error
    
    // Storage Information
    GetStorageInfo(ctx context.Context) (*StorageInfo, error)
}

// CDNService interface for content delivery network operations
type CDNService interface {
    GetCDNURL(key string) string
    InvalidateCache(ctx context.Context, keys []string) error
    SetCacheHeaders(key string, maxAge time.Duration) map[string]string
}
```

### Data Models:

```go
type UploadRequest struct {
    Key         string            `json:"key"`         // Storage path/key
    File        io.Reader         `json:"-"`           // File content
    ContentType string            `json:"content_type"`
    Size        int64             `json:"size"`
    Metadata    map[string]string `json:"metadata,omitempty"`
    CacheControl string           `json:"cache_control,omitempty"`
    ACL         string            `json:"acl,omitempty"` // Access control
}

type UploadResult struct {
    Key       string            `json:"key"`
    URL       string            `json:"url"`        // Direct Spaces URL
    CDNURL    string            `json:"cdn_url"`    // CDN URL
    ETag      string            `json:"etag"`
    Size      int64             `json:"size"`
    Metadata  map[string]string `json:"metadata,omitempty"`
    UploadedAt time.Time        `json:"uploaded_at"`
}

type FileMetadata struct {
    Key          string            `json:"key"`
    Size         int64             `json:"size"`
    ContentType  string            `json:"content_type"`
    ETag         string            `json:"etag"`
    LastModified time.Time         `json:"last_modified"`
    Metadata     map[string]string `json:"metadata,omitempty"`
    ACL          string            `json:"acl,omitempty"`
}

type StorageInfo struct {
    Provider    string `json:"provider"`     // "digitalocean-spaces"
    Region      string `json:"region"`       // e.g., "nyc3"
    Bucket      string `json:"bucket"`       // Bucket name
    CDNEnabled  bool   `json:"cdn_enabled"`  // CDN availability
    CDNEndpoint string `json:"cdn_endpoint"` // CDN base URL
}

type AccessURL struct {
    URL       string    `json:"url"`
    ExpiresAt time.Time `json:"expires_at"`
    Type      string    `json:"type"` // "cdn", "direct", "signed"
}
```

## Fiber Handlers

### Document Serving Handler:
```go
func (h *StorageHandler) ServeDocument(c *fiber.Ctx) error {
    filePath := c.Params("file_path")
    if filePath == "" {
        return fiber.NewError(fiber.StatusBadRequest, "File path required")
    }
    
    // URL decode the file path (handle double encoding)
    decodedPath, err := url.QueryUnescape(filePath)
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid file path")
    }
    
    // Convert S3 URI to storage key if needed
    storageKey := h.convertS3URIToKey(decodedPath)
    
    // Check if file exists
    exists, err := h.storageService.FileExists(c.Context(), storageKey)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    if !exists {
        return fiber.NewError(fiber.StatusNotFound, "Document not found")
    }
    
    // Generate CDN URL and redirect
    cdnURL, err := h.storageService.GenerateCDNURL(c.Context(), storageKey)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    // Redirect to CDN with caching headers
    c.Set("Cache-Control", "public, max-age=3600")
    c.Set("X-Content-Type-Options", "nosniff")
    
    return c.Redirect(cdnURL, fiber.StatusFound)
}
```

### Upload Handler:
```go
func (h *StorageHandler) UploadDocument(c *fiber.Ctx) error {
    file, err := c.FormFile("document")
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "No file provided")
    }
    
    // Open file for reading
    src, err := file.Open()
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Failed to open file")
    }
    defer src.Close()
    
    // Generate storage key with folder structure
    storageKey := h.generateStorageKey(file.Filename)
    
    // Prepare upload request
    uploadReq := &UploadRequest{
        Key:         storageKey,
        File:        src,
        ContentType: file.Header.Get("Content-Type"),
        Size:        file.Size,
        Metadata: map[string]string{
            "original-filename": file.Filename,
            "uploaded-by":       c.Get("user-id", "anonymous"),
            "upload-time":       time.Now().UTC().Format(time.RFC3339),
        },
        ACL: "private", // Documents are private by default
    }
    
    // Upload to DigitalOcean Spaces
    result, err := h.storageService.UploadFile(c.Context(), uploadReq)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Document uploaded successfully",
        "file":    result,
    })
}
```

### Storage Info Handler:
```go
func (h *StorageHandler) GetStorageInfo(c *fiber.Ctx) error {
    info, err := h.storageService.GetStorageInfo(c.Context())
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    return c.JSON(info)
}
```

## DigitalOcean Spaces Implementation

### Client Configuration:
```go
type SpacesClient struct {
    s3Client   *s3.Client
    bucket     string
    region     string
    cdnURL     string
    endpoint   string
}

func NewSpacesClient(config *SpacesConfig) (*SpacesClient, error) {
    // Configure AWS SDK for DigitalOcean Spaces
    cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
        awsconfig.WithRegion(config.Region),
        awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            config.AccessKey,
            config.SecretKey,
            "",
        )),
        awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
            func(service, region string, options ...interface{}) (aws.Endpoint, error) {
                return aws.Endpoint{
                    URL:           fmt.Sprintf("https://%s.digitaloceanspaces.com", config.Region),
                    SigningRegion: config.Region,
                }, nil
            }),
        ),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to load AWS config: %w", err)
    }
    
    s3Client := s3.NewFromConfig(cfg)
    
    return &SpacesClient{
        s3Client: s3Client,
        bucket:   config.Bucket,
        region:   config.Region,
        cdnURL:   config.CDNURL,
        endpoint: fmt.Sprintf("https://%s.digitaloceanspaces.com", config.Region),
    }, nil
}
```

### Upload Implementation:
```go
func (sc *SpacesClient) UploadFile(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
    // Prepare S3 upload input
    input := &s3.PutObjectInput{
        Bucket:      &sc.bucket,
        Key:         &req.Key,
        Body:        req.File,
        ContentType: &req.ContentType,
        ACL:         types.ObjectCannedACL(req.ACL),
    }
    
    // Add metadata
    if len(req.Metadata) > 0 {
        input.Metadata = req.Metadata
    }
    
    // Add cache control
    if req.CacheControl != "" {
        input.CacheControl = &req.CacheControl
    }
    
    // Upload to Spaces
    result, err := sc.s3Client.PutObject(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("failed to upload to Spaces: %w", err)
    }
    
    // Generate URLs
    directURL := fmt.Sprintf("%s/%s/%s", sc.endpoint, sc.bucket, req.Key)
    cdnURL := fmt.Sprintf("%s/%s", sc.cdnURL, req.Key)
    
    return &UploadResult{
        Key:        req.Key,
        URL:        directURL,
        CDNURL:     cdnURL,
        ETag:       strings.Trim(*result.ETag, "\""),
        Size:       req.Size,
        Metadata:   req.Metadata,
        UploadedAt: time.Now().UTC(),
    }, nil
}
```

### CDN URL Generation:
```go
func (sc *SpacesClient) GenerateCDNURL(ctx context.Context, key string) (string, error) {
    if sc.cdnURL == "" {
        // Fallback to direct Spaces URL
        return fmt.Sprintf("%s/%s/%s", sc.endpoint, sc.bucket, key), nil
    }
    
    return fmt.Sprintf("%s/%s", sc.cdnURL, key), nil
}

func (sc *SpacesClient) GenerateSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
    presignClient := s3.NewPresignClient(sc.s3Client)
    
    request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
        Bucket: &sc.bucket,
        Key:    &key,
    }, func(opts *s3.PresignOptions) {
        opts.Expires = expiry
    })
    
    if err != nil {
        return "", fmt.Errorf("failed to generate signed URL: %w", err)
    }
    
    return request.URL, nil
}
```

## Helper Functions

### Storage Key Generation:
```go
func (h *StorageHandler) generateStorageKey(filename string) string {
    // Generate folder structure: documents/YYYY/MM/DD/filename_hash.ext
    now := time.Now().UTC()
    
    // Generate short hash for uniqueness
    hash := sha256.Sum256([]byte(filename + now.Format(time.RFC3339Nano)))
    shortHash := hex.EncodeToString(hash[:4]) // 8 character hash
    
    // Extract file extension
    ext := filepath.Ext(filename)
    nameWithoutExt := strings.TrimSuffix(filename, ext)
    
    // Clean filename
    cleanName := h.sanitizeFilename(nameWithoutExt)
    
    // Construct key
    return fmt.Sprintf("documents/%04d/%02d/%02d/%s_%s%s",
        now.Year(), now.Month(), now.Day(),
        cleanName, shortHash, ext)
}

func (h *StorageHandler) sanitizeFilename(filename string) string {
    // Remove invalid characters and limit length
    reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
    clean := reg.ReplaceAllString(filename, "_")
    
    if len(clean) > 100 {
        clean = clean[:100]
    }
    
    return clean
}

func (h *StorageHandler) convertS3URIToKey(uri string) string {
    // Convert s3://bucket/key format to just key
    if strings.HasPrefix(uri, "s3://") {
        parts := strings.SplitN(strings.TrimPrefix(uri, "s3://"), "/", 2)
        if len(parts) > 1 {
            return parts[1] // Return everything after bucket name
        }
    }
    
    return uri
}
```

## Test Strategy

### Unit Tests:
```go
func TestSpacesClient_UploadFile(t *testing.T) {
    tests := []struct {
        name    string
        request *UploadRequest
        want    *UploadResult
        wantErr bool
    }{
        {
            name: "successful upload",
            request: &UploadRequest{
                Key:         "test/document.pdf",
                File:        strings.NewReader("test content"),
                ContentType: "application/pdf",
                Size:        12,
            },
            want: &UploadResult{
                Key:  "test/document.pdf",
                Size: 12,
            },
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Mock Spaces client and test
        })
    }
}
```

### Integration Tests:
- Real DigitalOcean Spaces upload/download
- CDN URL generation and access
- Signed URL functionality
- File existence checks
- Metadata handling

## Implementation Priority

1. **Spaces Client** - Basic S3-compatible operations
2. **File Upload** - Document upload with metadata
3. **CDN Integration** - URL generation and redirects
4. **File Serving** - Document access via redirects
5. **Advanced Features** - Signed URLs, batch operations
6. **Monitoring** - Storage metrics and health checks

## Dependencies

### External Libraries:
- `github.com/aws/aws-sdk-go-v2` - AWS SDK for S3-compatible operations
- `github.com/aws/aws-sdk-go-v2/service/s3` - S3 service client
- `github.com/aws/aws-sdk-go-v2/config` - AWS configuration

### Configuration:
```go
type SpacesConfig struct {
    AccessKey   string `env:"DO_SPACES_KEY" required:"true"`
    SecretKey   string `env:"DO_SPACES_SECRET" required:"true"`
    Bucket      string `env:"DO_SPACES_BUCKET" required:"true"`
    Region      string `env:"DO_SPACES_REGION" default:"nyc3"`
    CDNURL      string `env:"DO_SPACES_CDN_URL"` // Optional CDN endpoint
    
    // Upload Configuration
    MaxFileSize int64         `env:"MAX_FILE_SIZE" default:"104857600"` // 100MB
    Timeout     time.Duration `env:"UPLOAD_TIMEOUT" default:"5m"`
    
    // URL Configuration
    SignedURLExpiry time.Duration `env:"SIGNED_URL_EXPIRY" default:"1h"`
    CDNCacheMaxAge  time.Duration `env:"CDN_CACHE_MAX_AGE" default:"24h"`
}
```

## Performance Considerations

- **Multipart Upload**: Use for large files (>100MB)
- **Connection Pooling**: Efficient HTTP client configuration
- **CDN Optimization**: Proper cache headers and invalidation
- **Parallel Operations**: Concurrent uploads for batch processing
- **Streaming**: Memory-efficient file handling with io.Reader

## Security Considerations

- **Access Control**: Private documents with signed URLs
- **File Validation**: Content type and size validation
- **Path Sanitization**: Prevent directory traversal in keys
- **URL Security**: Signed URLs with proper expiration
- **CORS Configuration**: Secure cross-origin access