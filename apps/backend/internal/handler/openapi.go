// Package handler provides HTTP handlers for the ARK API.
// This file contains handlers for serving OpenAPI documentation.
package handler

import (
	"fmt"
	"net/http"
	"os"

	"ark/internal/server"

	"github.com/labstack/echo/v4"
)

// OpenAPIHandler handles requests for OpenAPI documentation.
// It serves both the interactive UI and the raw JSON specification.
type OpenAPIHandler struct {
	Handler
}

// NewOpenAPIHandler creates a new OpenAPIHandler with the given server.
func NewOpenAPIHandler(s *server.Server) *OpenAPIHandler {
	return &OpenAPIHandler{
		Handler: NewHandler(s),
	}
}

// ServeOpenAPIUI serves the interactive OpenAPI documentation UI.
// This provides a user-friendly interface for exploring the ARK API.
//
// Route: GET /docs
func (h *OpenAPIHandler) ServeOpenAPIUI(c echo.Context) error {
	templateBytes, err := os.ReadFile("static/openapi.html")
	c.Response().Header().Set("Cache-Control", "no-cache")
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI UI template: %w", err)
	}

	templateString := string(templateBytes)

	err = c.HTML(http.StatusOK, templateString)
	if err != nil {
		return fmt.Errorf("failed to write HTML response: %w", err)
	}

	return nil
}

// ServeOpenAPISpec serves the raw OpenAPI JSON specification.
// This is primarily used for testing and programmatic access.
//
// Note: The spec is also available via the static file route at /static/openapi.json
func (h *OpenAPIHandler) ServeOpenAPISpec(c echo.Context) error {
	specBytes, err := os.ReadFile("static/openapi.json")
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI spec: %w", err)
	}

	return c.JSONBlob(http.StatusOK, specBytes)
}
