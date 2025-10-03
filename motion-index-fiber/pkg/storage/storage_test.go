package storage

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMockService(t *testing.T) {
	service := NewMockService()

	assert.NotNil(t, service)
	assert.True(t, service.IsHealthy())
}

func TestMockService_Upload(t *testing.T) {
	service := NewMockService()
	ctx := context.Background()

	tests := []struct {
		name     string
		path     string
		content  string
		metadata *UploadMetadata
		wantErr  bool
	}{
		{
			name:    "successful upload",
			path:    "documents/test.pdf",
			content: "test content",
			metadata: &UploadMetadata{
				ContentType: "application/pdf",
				Size:        12,
			},
			wantErr: false,
		},
		{
			name:     "upload without metadata",
			path:     "documents/test2.txt",
			content:  "another test",
			metadata: nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Upload(ctx, tt.path, strings.NewReader(tt.content), tt.metadata)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.path, result.Path)
				assert.Equal(t, int64(len(tt.content)), result.Size)
				assert.True(t, result.Success)
				assert.NotEmpty(t, result.URL)
				assert.NotEmpty(t, result.ETag)
				assert.False(t, result.UploadedAt.IsZero())
			}
		})
	}
}

func TestMockService_Download(t *testing.T) {
	service := NewMockService()
	ctx := context.Background()

	// First upload a document
	path := "documents/test-download.txt"
	content := "test download content"

	_, err := service.Upload(ctx, path, strings.NewReader(content), nil)
	require.NoError(t, err)

	// Test download
	t.Run("successful download", func(t *testing.T) {
		reader, err := service.Download(ctx, path)
		assert.NoError(t, err)
		assert.NotNil(t, reader)

		// Read the content
		buf := make([]byte, len(content))
		n, err := reader.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, len(content), n)
		assert.Equal(t, content, string(buf))

		reader.Close()
	})

	t.Run("download non-existent file", func(t *testing.T) {
		reader, err := service.Download(ctx, "non-existent.txt")
		assert.Error(t, err)
		assert.Nil(t, reader)

		storageErr, ok := err.(*StorageError)
		assert.True(t, ok)
		assert.Equal(t, "not_found", storageErr.Type)
	})
}

func TestMockService_Delete(t *testing.T) {
	service := NewMockService()
	ctx := context.Background()

	// First upload a document
	path := "documents/test-delete.txt"
	content := "test delete content"

	_, err := service.Upload(ctx, path, strings.NewReader(content), nil)
	require.NoError(t, err)

	// Test delete
	t.Run("successful delete", func(t *testing.T) {
		err := service.Delete(ctx, path)
		assert.NoError(t, err)

		// Verify it's deleted
		reader, err := service.Download(ctx, path)
		assert.Error(t, err)
		assert.Nil(t, reader)
	})

	t.Run("delete non-existent file", func(t *testing.T) {
		err := service.Delete(ctx, "non-existent.txt")
		assert.NoError(t, err) // Mock service doesn't error on non-existent delete
	})
}

func TestMockService_GetURL(t *testing.T) {
	service := NewMockService()

	path := "documents/test.pdf"
	url := service.GetURL(path)

	assert.NotEmpty(t, url)
	assert.Contains(t, url, path)
	assert.Contains(t, url, "mock-storage.example.com")
}

func TestMockService_GetSignedURL(t *testing.T) {
	service := NewMockService()

	path := "documents/test.pdf"
	expiration := time.Hour

	url, err := service.GetSignedURL(path, expiration)
	assert.NoError(t, err)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, path)
	assert.Contains(t, url, "signed=true")
	assert.Contains(t, url, "expires=")
}

func TestMockService_Exists(t *testing.T) {
	service := NewMockService()
	ctx := context.Background()

	// Test non-existent file
	exists, err := service.Exists(ctx, "non-existent.txt")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Upload a file and test existence
	path := "documents/test-exists.txt"
	_, err = service.Upload(ctx, path, strings.NewReader("test content"), nil)
	require.NoError(t, err)

	exists, err = service.Exists(ctx, path)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestMockService_List(t *testing.T) {
	service := NewMockService()
	ctx := context.Background()

	// Upload several files
	files := []string{
		"documents/test1.txt",
		"documents/test2.pdf",
		"documents/subdir/test3.txt",
		"other/test4.txt",
	}

	for _, file := range files {
		_, err := service.Upload(ctx, file, strings.NewReader("content"), nil)
		require.NoError(t, err)
	}

	// Test listing with prefix
	t.Run("list with prefix", func(t *testing.T) {
		objects, err := service.List(ctx, "documents/")
		assert.NoError(t, err)
		assert.Len(t, objects, 3) // Should exclude "other/test4.txt"

		for _, obj := range objects {
			assert.True(t, strings.HasPrefix(obj.Path, "documents/"))
			assert.Greater(t, obj.Size, int64(0))
			assert.False(t, obj.LastModified.IsZero())
			assert.NotEmpty(t, obj.ETag)
		}
	})

	t.Run("list all", func(t *testing.T) {
		objects, err := service.List(ctx, "")
		assert.NoError(t, err)
		assert.Len(t, objects, 4) // All files
	})
}

func TestMockService_HealthStatus(t *testing.T) {
	service := NewMockService().(*mockService)

	// Test healthy service
	assert.True(t, service.IsHealthy())

	// Test unhealthy service
	service.SetHealthy(false)
	assert.False(t, service.IsHealthy())

	// Test operations on unhealthy service
	ctx := context.Background()

	_, err := service.Upload(ctx, "test.txt", strings.NewReader("content"), nil)
	assert.Error(t, err)

	_, err = service.Download(ctx, "test.txt")
	assert.Error(t, err)

	err = service.Delete(ctx, "test.txt")
	assert.Error(t, err)

	_, err = service.GetSignedURL("test.txt", time.Hour)
	assert.Error(t, err)

	_, err = service.Exists(ctx, "test.txt")
	assert.Error(t, err)

	_, err = service.List(ctx, "")
	assert.Error(t, err)
}

func TestStorageError(t *testing.T) {
	t.Run("with cause", func(t *testing.T) {
		cause := assert.AnError
		err := NewStorageError("test_type", "test message", "test/path", cause)

		assert.Equal(t, "test_type", err.Type)
		assert.Equal(t, "test message", err.Message)
		assert.Equal(t, "test/path", err.Path)
		assert.Equal(t, cause, err.Cause)
		assert.Equal(t, "test message: assert.AnError general error for testing", err.Error())
	})

	t.Run("without cause", func(t *testing.T) {
		err := NewStorageError("test_type", "test message", "test/path", nil)

		assert.Equal(t, "test_type", err.Type)
		assert.Equal(t, "test message", err.Message)
		assert.Equal(t, "test/path", err.Path)
		assert.Nil(t, err.Cause)
		assert.Equal(t, "test message", err.Error())
	})
}

func TestUploadMetadata(t *testing.T) {
	metadata := &UploadMetadata{
		ContentType: "application/pdf",
		Size:        1024,
		Tags:        map[string]string{"category": "legal", "type": "motion"},
	}

	assert.Equal(t, "application/pdf", metadata.ContentType)
	assert.Equal(t, int64(1024), metadata.Size)
	assert.Equal(t, "legal", metadata.Tags["category"])
	assert.Equal(t, "motion", metadata.Tags["type"])
}

func TestUploadResult(t *testing.T) {
	result := &UploadResult{
		Path:       "documents/test.pdf",
		URL:        "https://example.com/documents/test.pdf",
		Size:       1024,
		ETag:       "abc123",
		Success:    true,
		UploadedAt: time.Now(),
	}

	assert.Equal(t, "documents/test.pdf", result.Path)
	assert.Equal(t, "https://example.com/documents/test.pdf", result.URL)
	assert.Equal(t, int64(1024), result.Size)
	assert.Equal(t, "abc123", result.ETag)
	assert.True(t, result.Success)
	assert.False(t, result.UploadedAt.IsZero())
}

func TestStorageObject(t *testing.T) {
	obj := &StorageObject{
		Path:         "documents/test.pdf",
		Size:         1024,
		LastModified: time.Now(),
		ETag:         "abc123",
		ContentType:  "application/pdf",
	}

	assert.Equal(t, "documents/test.pdf", obj.Path)
	assert.Equal(t, int64(1024), obj.Size)
	assert.False(t, obj.LastModified.IsZero())
	assert.Equal(t, "abc123", obj.ETag)
	assert.Equal(t, "application/pdf", obj.ContentType)
}
