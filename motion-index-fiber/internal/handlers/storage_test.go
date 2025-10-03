package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"motion-index-fiber/internal/config"
)

func TestNewStorageHandler(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "6000",
		},
	}

	handler := NewStorageHandler(cfg)

	assert.NotNil(t, handler)
	assert.Equal(t, cfg, handler.cfg)
}

func TestStorageHandler_ServeDocument(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "6000",
		},
	}

	handler := NewStorageHandler(cfg)
	app := fiber.New()
	app.Get("/documents/:path", handler.ServeDocument)

	// Test serving a document (should return not implemented for now)
	req := httptest.NewRequest("GET", "/documents/test-doc.pdf", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotImplemented, resp.StatusCode)

	// Verify response content
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Document serving not yet implemented", response["message"])
	assert.Equal(t, "not_implemented", response["status"])
}

func TestStorageHandler_ServeDocument_EmptyPath(t *testing.T) {
	cfg := &config.Config{}
	handler := NewStorageHandler(cfg)
	app := fiber.New()
	app.Get("/documents/:path?", handler.ServeDocument)

	// Test with empty path
	req := httptest.NewRequest("GET", "/documents/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	// Should still return not implemented since the method isn't implemented yet
	assert.Equal(t, fiber.StatusNotImplemented, resp.StatusCode)
}

func TestStorageHandler_ServeDocument_WithParams(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "6000",
		},
	}

	handler := NewStorageHandler(cfg)
	app := fiber.New()
	app.Get("/documents/*", handler.ServeDocument)

	// Test with nested path
	req := httptest.NewRequest("GET", "/documents/subfolder/test-doc.pdf", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotImplemented, resp.StatusCode)
}

func TestStorageHandler_ServeDocument_NilConfig(t *testing.T) {
	// Test behavior with nil config
	handler := NewStorageHandler(nil)
	app := fiber.New()
	app.Get("/documents/:path", handler.ServeDocument)

	req := httptest.NewRequest("GET", "/documents/test.pdf", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	// Should still return not implemented
	assert.Equal(t, fiber.StatusNotImplemented, resp.StatusCode)
}
