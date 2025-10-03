package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockS3Client for testing
type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) Upload(ctx context.Context, key string, content io.Reader, size int64, contentType string) (*UploadResult, error) {
	args := m.Called(ctx, key, content, size, contentType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UploadResult), args.Error(1)
}

func (m *MockS3Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockS3Client) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockS3Client) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockS3Client) GetMetadata(ctx context.Context, key string) (*FileMetadata, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*FileMetadata), args.Error(1)
}

func (m *MockS3Client) GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	args := m.Called(ctx, key, expiration)
	return args.String(0), args.Error(1)
}

func (m *MockS3Client) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	args := m.Called(ctx, prefix)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockS3Client) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestNewSpacesService(t *testing.T) {
	tests := []struct {
		name        string
		endpoint    string
		region      string
		bucket      string
		accessKey   string
		secretKey   string
		expectError bool
	}{
		{
			name:        "valid configuration",
			endpoint:    "https://nyc3.digitaloceanspaces.com",
			region:      "nyc3",
			bucket:      "test-bucket",
			accessKey:   "test-access-key",
			secretKey:   "test-secret-key",
			expectError: false,
		},
		{
			name:        "empty endpoint",
			endpoint:    "",
			region:      "nyc3",
			bucket:      "test-bucket",
			accessKey:   "test-access-key",
			secretKey:   "test-secret-key",
			expectError: true,
		},
		{
			name:        "empty region",
			endpoint:    "https://nyc3.digitaloceanspaces.com",
			region:      "",
			bucket:      "test-bucket",
			accessKey:   "test-access-key",
			secretKey:   "test-secret-key",
			expectError: true,
		},
		{
			name:        "empty bucket",
			endpoint:    "https://nyc3.digitaloceanspaces.com",
			region:      "nyc3",
			bucket:      "",
			accessKey:   "test-access-key",
			secretKey:   "test-secret-key",
			expectError: true,
		},
		{
			name:        "empty access key",
			endpoint:    "https://nyc3.digitaloceanspaces.com",
			region:      "nyc3",
			bucket:      "test-bucket",
			accessKey:   "",
			secretKey:   "test-secret-key",
			expectError: true,
		},
		{
			name:        "empty secret key",
			endpoint:    "https://nyc3.digitaloceanspaces.com",
			region:      "nyc3",
			bucket:      "test-bucket",
			accessKey:   "test-access-key",
			secretKey:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &SpacesConfig{
				Endpoint:  tt.endpoint,
				Region:    tt.region,
				Bucket:    tt.bucket,
				AccessKey: tt.accessKey,
				SecretKey: tt.secretKey,
			}

			service, err := NewSpacesService(config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
			}
		})
	}
}

func TestSpacesService_Upload(t *testing.T) {
	mockClient := &MockS3Client{}
	service := &SpacesService{
		client: mockClient,
		config: &SpacesConfig{
			Bucket: "test-bucket",
		},
	}

	tests := []struct {
		name          string
		key           string
		content       []byte
		contentType   string
		mockResult    *UploadResult
		mockError     error
		expectedError bool
	}{
		{
			name:        "successful upload",
			key:         "documents/test.pdf",
			content:     []byte("test content"),
			contentType: "application/pdf",
			mockResult: &UploadResult{
				URL:        "https://test-bucket.nyc3.digitaloceanspaces.com/documents/test.pdf",
				Path:       "documents/test.pdf",
				Size:       12,
				Success:    true,
				UploadedAt: time.Now(),
			},
			expectedError: false,
		},
		{
			name:          "upload failure",
			key:           "documents/test.pdf",
			content:       []byte("test content"),
			contentType:   "application/pdf",
			mockError:     assert.AnError,
			expectedError: true,
		},
		{
			name:          "empty key",
			key:           "",
			content:       []byte("test content"),
			contentType:   "application/pdf",
			expectedError: true,
		},
		{
			name:          "empty content",
			key:           "documents/test.pdf",
			content:       []byte{},
			contentType:   "application/pdf",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.key != "" && len(tt.content) > 0 {
				mockClient.On("Upload", mock.Anything, tt.key, mock.Anything, int64(len(tt.content)), tt.contentType).Return(tt.mockResult, tt.mockError)
			}

			ctx := context.Background()
			result, err := service.Upload(ctx, tt.key, bytes.NewReader(tt.content), int64(len(tt.content)), tt.contentType)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockResult, result)
			}

			if tt.key != "" && len(tt.content) > 0 {
				mockClient.AssertExpectations(t)
			}
			mockClient.ExpectedCalls = nil // Reset for next test
		})
	}
}

func TestSpacesService_Download(t *testing.T) {
	mockClient := &MockS3Client{}
	service := &SpacesService{
		client: mockClient,
		config: &SpacesConfig{
			Bucket: "test-bucket",
		},
	}

	tests := []struct {
		name          string
		key           string
		mockContent   io.ReadCloser
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful download",
			key:           "documents/test.pdf",
			mockContent:   io.NopCloser(bytes.NewReader([]byte("test content"))),
			expectedError: false,
		},
		{
			name:          "download failure",
			key:           "documents/nonexistent.pdf",
			mockError:     assert.AnError,
			expectedError: true,
		},
		{
			name:          "empty key",
			key:           "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.key != "" {
				mockClient.On("Download", mock.Anything, tt.key).Return(tt.mockContent, tt.mockError)
			}

			ctx := context.Background()
			result, err := service.Download(ctx, tt.key)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				// Read content to verify
				content, readErr := io.ReadAll(result)
				assert.NoError(t, readErr)
				assert.Equal(t, []byte("test content"), content)
				result.Close()
			}

			if tt.key != "" {
				mockClient.AssertExpectations(t)
			}
			mockClient.ExpectedCalls = nil // Reset for next test
		})
	}
}

func TestSpacesService_Delete(t *testing.T) {
	mockClient := &MockS3Client{}
	service := &SpacesService{
		client: mockClient,
		config: &SpacesConfig{
			Bucket: "test-bucket",
		},
	}

	tests := []struct {
		name          string
		key           string
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful deletion",
			key:           "documents/test.pdf",
			expectedError: false,
		},
		{
			name:          "deletion failure",
			key:           "documents/protected.pdf",
			mockError:     assert.AnError,
			expectedError: true,
		},
		{
			name:          "empty key",
			key:           "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.key != "" {
				mockClient.On("Delete", mock.Anything, tt.key).Return(tt.mockError)
			}

			ctx := context.Background()
			err := service.Delete(ctx, tt.key)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.key != "" {
				mockClient.AssertExpectations(t)
			}
			mockClient.ExpectedCalls = nil // Reset for next test
		})
	}
}

func TestSpacesService_Exists(t *testing.T) {
	mockClient := &MockS3Client{}
	service := &SpacesService{
		client: mockClient,
		config: &SpacesConfig{
			Bucket: "test-bucket",
		},
	}

	tests := []struct {
		name           string
		key            string
		mockExists     bool
		mockError      error
		expectedError  bool
		expectedExists bool
	}{
		{
			name:           "file exists",
			key:            "documents/test.pdf",
			mockExists:     true,
			expectedError:  false,
			expectedExists: true,
		},
		{
			name:           "file does not exist",
			key:            "documents/nonexistent.pdf",
			mockExists:     false,
			expectedError:  false,
			expectedExists: false,
		},
		{
			name:          "check failure",
			key:           "documents/error.pdf",
			mockError:     assert.AnError,
			expectedError: true,
		},
		{
			name:          "empty key",
			key:           "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.key != "" {
				mockClient.On("Exists", mock.Anything, tt.key).Return(tt.mockExists, tt.mockError)
			}

			ctx := context.Background()
			exists, err := service.Exists(ctx, tt.key)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExists, exists)
			}

			if tt.key != "" {
				mockClient.AssertExpectations(t)
			}
			mockClient.ExpectedCalls = nil // Reset for next test
		})
	}
}

func TestSpacesService_GeneratePresignedURL(t *testing.T) {
	mockClient := &MockS3Client{}
	service := &SpacesService{
		client: mockClient,
		config: &SpacesConfig{
			Bucket: "test-bucket",
		},
	}

	tests := []struct {
		name          string
		key           string
		expiration    time.Duration
		mockURL       string
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful URL generation",
			key:           "documents/test.pdf",
			expiration:    1 * time.Hour,
			mockURL:       "https://test-bucket.nyc3.digitaloceanspaces.com/documents/test.pdf?presigned=true",
			expectedError: false,
		},
		{
			name:          "URL generation failure",
			key:           "documents/error.pdf",
			expiration:    1 * time.Hour,
			mockError:     assert.AnError,
			expectedError: true,
		},
		{
			name:          "empty key",
			key:           "",
			expiration:    1 * time.Hour,
			expectedError: true,
		},
		{
			name:          "zero expiration",
			key:           "documents/test.pdf",
			expiration:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.key != "" && tt.expiration > 0 {
				mockClient.On("GeneratePresignedURL", mock.Anything, tt.key, tt.expiration).Return(tt.mockURL, tt.mockError)
			}

			ctx := context.Background()
			url, err := service.GeneratePresignedURL(ctx, tt.key, tt.expiration)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockURL, url)
			}

			if tt.key != "" && tt.expiration > 0 {
				mockClient.AssertExpectations(t)
			}
			mockClient.ExpectedCalls = nil // Reset for next test
		})
	}
}

func TestSpacesService_IsHealthy(t *testing.T) {
	tests := []struct {
		name            string
		mockHealthy     bool
		expectedHealthy bool
	}{
		{
			name:            "healthy client",
			mockHealthy:     true,
			expectedHealthy: true,
		},
		{
			name:            "unhealthy client",
			mockHealthy:     false,
			expectedHealthy: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockS3Client{}
			service := &SpacesService{
				client: mockClient,
				config: &SpacesConfig{
					Bucket: "test-bucket",
				},
			}

			mockClient.On("IsHealthy").Return(tt.mockHealthy)

			result := service.IsHealthy()
			assert.Equal(t, tt.expectedHealthy, result)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestValidateSpacesConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *SpacesConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: &SpacesConfig{
				Endpoint:  "https://nyc3.digitaloceanspaces.com",
				Region:    "nyc3",
				Bucket:    "test-bucket",
				AccessKey: "access-key",
				SecretKey: "secret-key",
			},
			valid: true,
		},
		{
			name:   "nil config",
			config: nil,
			valid:  false,
		},
		{
			name: "empty endpoint",
			config: &SpacesConfig{
				Endpoint:  "",
				Region:    "nyc3",
				Bucket:    "test-bucket",
				AccessKey: "access-key",
				SecretKey: "secret-key",
			},
			valid: false,
		},
		{
			name: "empty region",
			config: &SpacesConfig{
				Endpoint:  "https://nyc3.digitaloceanspaces.com",
				Region:    "",
				Bucket:    "test-bucket",
				AccessKey: "access-key",
				SecretKey: "secret-key",
			},
			valid: false,
		},
		{
			name: "empty bucket",
			config: &SpacesConfig{
				Endpoint:  "https://nyc3.digitaloceanspaces.com",
				Region:    "nyc3",
				Bucket:    "",
				AccessKey: "access-key",
				SecretKey: "secret-key",
			},
			valid: false,
		},
		{
			name: "empty access key",
			config: &SpacesConfig{
				Endpoint:  "https://nyc3.digitaloceanspaces.com",
				Region:    "nyc3",
				Bucket:    "test-bucket",
				AccessKey: "",
				SecretKey: "secret-key",
			},
			valid: false,
		},
		{
			name: "empty secret key",
			config: &SpacesConfig{
				Endpoint:  "https://nyc3.digitaloceanspaces.com",
				Region:    "nyc3",
				Bucket:    "test-bucket",
				AccessKey: "access-key",
				SecretKey: "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSpacesConfig(tt.config)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestSpacesService_BatchOperations(t *testing.T) {
	mockClient := &MockS3Client{}
	service := &SpacesService{
		client: mockClient,
		config: &SpacesConfig{
			Bucket: "test-bucket",
		},
	}

	t.Run("batch upload", func(t *testing.T) {
		files := []struct {
			key     string
			content []byte
		}{
			{"documents/file1.pdf", []byte("content1")},
			{"documents/file2.pdf", []byte("content2")},
			{"documents/file3.pdf", []byte("content3")},
		}

		// Mock successful uploads
		for _, file := range files {
			mockClient.On("Upload", mock.Anything, file.key, mock.Anything, int64(len(file.content)), "application/pdf").Return(&UploadResult{
				URL:     "https://test-bucket.nyc3.digitaloceanspaces.com/" + file.key,
				Path:    file.key,
				Size:    int64(len(file.content)),
				Success: true,
			}, nil)
		}

		ctx := context.Background()
		results := make([]*UploadResult, 0, len(files))

		for _, file := range files {
			result, err := service.Upload(ctx, file.key, bytes.NewReader(file.content), int64(len(file.content)), "application/pdf")
			assert.NoError(t, err)
			assert.NotNil(t, result)
			results = append(results, result)
		}

		assert.Len(t, results, 3)
		mockClient.AssertExpectations(t)
	})
}

// Helper function that might be in the actual implementation
func validateSpacesConfig(config *SpacesConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if config.Endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}
	if config.Region == "" {
		return fmt.Errorf("region is required")
	}
	if config.Bucket == "" {
		return fmt.Errorf("bucket is required")
	}
	if config.AccessKey == "" {
		return fmt.Errorf("access key is required")
	}
	if config.SecretKey == "" {
		return fmt.Errorf("secret key is required")
	}
	return nil
}
