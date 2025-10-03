package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"motion-index-fiber/pkg/search"
	"motion-index-fiber/pkg/search/models"
)

// SearchHandler handles search-related HTTP requests
type SearchHandler struct {
	searchService search.Service
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService search.Service) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// SearchDocuments handles POST /search
func (h *SearchHandler) SearchDocuments(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	var req models.SearchRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	// Parse query parameters if body is empty
	if req.Query == "" {
		req.Query = c.Query("q", "")
	}
	if req.Size == 0 {
		if sizeStr := c.Query("size"); sizeStr != "" {
			if size, err := strconv.Atoi(sizeStr); err == nil {
				req.Size = size
			}
		}
	}
	if req.From == 0 {
		if fromStr := c.Query("from"); fromStr != "" {
			if from, err := strconv.Atoi(fromStr); err == nil {
				req.From = from
			}
		}
	}

	// Validate request
	if err := validateSearchRequest(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Execute search
	result, err := h.searchService.SearchDocuments(ctx, &req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Search failed: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"data":    result,
		"message": "Search completed successfully",
	})
}

// GetLegalTags handles GET /legal-tags
func (h *SearchHandler) GetLegalTags(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	tags, err := h.searchService.GetLegalTags(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve legal tags: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   tags,
	})
}

// GetDocumentTypes handles GET /document-types
func (h *SearchHandler) GetDocumentTypes(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	types, err := h.searchService.GetDocumentTypes(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve document types: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   types,
	})
}

// GetDocumentStats handles GET /document-stats
func (h *SearchHandler) GetDocumentStats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	stats, err := h.searchService.GetDocumentStats(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve document stats: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   stats,
	})
}

// GetFieldOptions handles GET /field-options
func (h *SearchHandler) GetFieldOptions(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	options, err := h.searchService.GetAllFieldOptions(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve field options: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   options,
	})
}

// GetMetadataFieldValues handles GET /metadata-fields/{field}
func (h *SearchHandler) GetMetadataFieldValues(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	field := c.Params("field")
	if field == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Field parameter is required")
	}

	// Parse query parameters
	prefix := c.Query("prefix", "")
	size := 50
	if sizeStr := c.Query("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			size = s
		}
	}

	values, err := h.searchService.GetMetadataFieldValues(ctx, field, prefix, size)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve field values: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   values,
	})
}

// GetDocument handles GET /documents/{id}
func (h *SearchHandler) GetDocument(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	docID := c.Params("id")
	if docID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Document ID is required")
	}

	document, err := h.searchService.GetDocument(ctx, docID)
	if err != nil {
		if err.Error() == "document not found" {
			return fiber.NewError(fiber.StatusNotFound, "Document not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve document: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   document,
	})
}

// DeleteDocument handles DELETE /documents/{id} (protected endpoint)
func (h *SearchHandler) DeleteDocument(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	docID := c.Params("id")
	if docID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Document ID is required")
	}

	err := h.searchService.DeleteDocument(ctx, docID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete document: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document deleted successfully",
	})
}

// validateSearchRequest validates a search request
func validateSearchRequest(req *models.SearchRequest) error {
	if req.Size > models.MaxSearchSize {
		req.Size = models.MaxSearchSize
	}
	if req.Size <= 0 {
		req.Size = models.DefaultSearchSize
	}
	if req.From < 0 {
		req.From = 0
	}

	// Validate sort order
	if req.SortOrder != "" && req.SortOrder != "asc" && req.SortOrder != "desc" {
		req.SortOrder = "desc"
	}

	return nil
}
