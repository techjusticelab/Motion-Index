package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"motion-index-fiber/internal/config"
	"motion-index-fiber/internal/models"
	"motion-index-fiber/pkg/storage"
)

type StorageHandler struct {
	cfg     *config.Config
	storage storage.Service
}

func NewStorageHandler(cfg *config.Config, storage storage.Service) *StorageHandler {
	return &StorageHandler{
		cfg:     cfg,
		storage: storage,
	}
}

// ListDocuments handles GET /api/storage/documents - List documents with pagination
func (h *StorageHandler) ListDocuments(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 60*time.Second)
	defer cancel()

	// Parse query parameters
	prefix := c.Query("prefix", "documents/")
	limitStr := c.Query("limit", "50")
	cursor := c.Query("cursor", "")
	fileType := c.Query("file_type", "")
	minSizeStr := c.Query("min_size", "0")
	maxSizeStr := c.Query("max_size", "")

	// Parse and validate limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 500 {
		limit = 50
	}

	// Parse size filters
	minSize, _ := strconv.ParseInt(minSizeStr, 10, 64)
	var maxSize int64 = -1
	if maxSizeStr != "" {
		maxSize, _ = strconv.ParseInt(maxSizeStr, 10, 64)
	}

	// List documents from storage
	objects, err := h.storage.List(ctx, prefix)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"storage_error",
			"Failed to list documents",
			map[string]interface{}{"error": err.Error()},
		))
	}

	// Apply filters and pagination
	filtered := h.filterDocuments(objects, fileType, minSize, maxSize)
	paginatedResult := h.paginateDocuments(filtered, cursor, limit)

	return c.JSON(models.NewSuccessResponse(paginatedResult, "Documents listed successfully"))
}

// GetDocumentsCount handles GET /api/storage/documents/count - Get total document count
func (h *StorageHandler) GetDocumentsCount(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 60*time.Second)
	defer cancel()

	// Parse query parameters for filtering
	prefix := c.Query("prefix", "documents/")
	fileType := c.Query("file_type", "")
	minSizeStr := c.Query("min_size", "0")
	maxSizeStr := c.Query("max_size", "")

	// Parse size filters
	minSize, _ := strconv.ParseInt(minSizeStr, 10, 64)
	var maxSize int64 = -1
	if maxSizeStr != "" {
		maxSize, _ = strconv.ParseInt(maxSizeStr, 10, 64)
	}

	// List and filter documents
	objects, err := h.storage.List(ctx, prefix)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"storage_error",
			"Failed to count documents",
			map[string]interface{}{"error": err.Error()},
		))
	}

	filtered := h.filterDocuments(objects, fileType, minSize, maxSize)

	response := map[string]interface{}{
		"total_count":    len(filtered),
		"prefix":         prefix,
		"applied_filters": map[string]interface{}{
			"file_type": fileType,
			"min_size":  minSize,
			"max_size":  maxSize,
		},
	}

	return c.JSON(models.NewSuccessResponse(response, "Document count retrieved successfully"))
}

func (h *StorageHandler) ServeDocument(c *fiber.Ctx) error {
	// Get document path from URL parameters
	documentPath := c.Params("*")
	if documentPath == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Document path is required",
		})
	}

	// For production, this should redirect to CDN URL
	// For now, return a not implemented response with the path
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"message": "Document serving via CDN redirect not yet implemented",
		"path":    documentPath,
		"status":  "not_implemented",
	})
}

// filterDocuments applies file type and size filters to the document list
func (h *StorageHandler) filterDocuments(objects []*storage.StorageObject, fileType string, minSize, maxSize int64) []*storage.StorageObject {
	var filtered []*storage.StorageObject

	for _, obj := range objects {
		// Skip directories
		if strings.HasSuffix(obj.Path, "/") {
			continue
		}

		// Skip very small files (likely empty or corrupt)
		if obj.Size < 100 {
			continue
		}

		// Skip system files
		filename := filepath.Base(obj.Path)
		if strings.Contains(filename, "__MACOSX") ||
			strings.Contains(filename, ".DS_Store") ||
			strings.HasSuffix(filename, ".tmp") ||
			strings.HasSuffix(filename, ".log") {
			continue
		}

		// Apply file type filter
		if fileType != "" {
			ext := strings.ToLower(filepath.Ext(obj.Path))
			if !strings.HasSuffix(ext, strings.ToLower(fileType)) {
				continue
			}
		}

		// Apply size filters
		if obj.Size < minSize {
			continue
		}
		if maxSize > 0 && obj.Size > maxSize {
			continue
		}

		filtered = append(filtered, obj)
	}

	return filtered
}

// paginateDocuments implements cursor-based pagination
func (h *StorageHandler) paginateDocuments(objects []*storage.StorageObject, cursor string, limit int) map[string]interface{} {
	startIndex := 0

	// Decode cursor if provided
	if cursor != "" {
		if decoded, err := base64.URLEncoding.DecodeString(cursor); err == nil {
			var cursorData map[string]interface{}
			if json.Unmarshal(decoded, &cursorData) == nil {
				if idx, ok := cursorData["index"].(float64); ok {
					startIndex = int(idx)
				}
			}
		}
	}

	// Apply pagination
	endIndex := startIndex + limit
	if endIndex > len(objects) {
		endIndex = len(objects)
	}

	var paginatedObjects []*storage.StorageObject
	if startIndex < len(objects) {
		paginatedObjects = objects[startIndex:endIndex]
	}

	// Generate next cursor
	var nextCursor string
	hasMore := endIndex < len(objects)
	if hasMore {
		cursorData := map[string]interface{}{
			"index": endIndex,
		}
		if cursorBytes, err := json.Marshal(cursorData); err == nil {
			nextCursor = base64.URLEncoding.EncodeToString(cursorBytes)
		}
	}

	// Convert storage objects to response format
	var documents []map[string]interface{}
	for _, obj := range paginatedObjects {
		documents = append(documents, map[string]interface{}{
			"path":          obj.Path,
			"size":          obj.Size,
			"last_modified": obj.LastModified,
			"file_type":     strings.ToLower(filepath.Ext(obj.Path)),
			"filename":      filepath.Base(obj.Path),
		})
	}

	return map[string]interface{}{
		"documents":        documents,
		"next_cursor":      nextCursor,
		"has_more":         hasMore,
		"total_returned":   len(documents),
		"total_estimated":  len(objects),
	}
}
