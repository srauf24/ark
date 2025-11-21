// Package handler provides HTTP request handlers for the Ark API.
// This file contains handlers for asset log operations.
package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"ark/internal/middleware"
	"ark/internal/model"
	"ark/internal/service"
)

// LogHandler handles HTTP requests for asset log operations.
// It provides CRUD operations for configuration logs and troubleshooting notes
// associated with homelab assets.
//
// Routes:
//   - POST   /api/v1/assets/:id/logs  - Create log for asset (nested)
//   - GET    /api/v1/assets/:id/logs  - List logs for asset (nested)
//   - GET    /api/v1/logs/:id          - Get log by ID (flat)
//   - PATCH  /api/v1/logs/:id          - Update log (flat)
//   - DELETE /api/v1/logs/:id          - Delete log (flat)
//
// All endpoints require authentication via the auth middleware.
type LogHandler struct {
	service *service.LogService
}

// NewLogHandler creates a new LogHandler with the given LogService.
// The service handles business logic, validation, and data persistence.
func NewLogHandler(service *service.LogService) *LogHandler {
	return &LogHandler{
		service: service,
	}
}

// ListByAsset handles GET /api/v1/assets/:id/logs
//
// Returns a paginated list of logs for the specified asset. This is a nested route
// that requires the asset ID in the URL path.
//
// Authentication: Required (user_id from context)
//
// URL Parameters:
//   - id: UUID of the asset (required)
//
// Query Parameters:
//   - limit:  Maximum number of logs to return (default: 50, max: 200)
//   - offset: Number of logs to skip for pagination (default: 0)
//   - tags:   Filter by tags (optional, can specify multiple)
//   - search: Search in log content (optional)
//   - start_date: Filter logs created after this date (optional)
//   - end_date: Filter logs created before this date (optional)
//   - sort_by: Field to sort by (default: "created_at")
//   - sort_order: Sort direction "asc" or "desc" (default: "desc")
//
// Response:
//   - 200 OK: Returns LogListResponse with logs array and pagination metadata
//   - 400 Bad Request: Invalid asset ID format or query parameters
//   - 401 Unauthorized: Missing or invalid authentication
//   - 404 Not Found: Asset doesn't exist or belongs to another user
//
// Security:
//   - Service layer verifies asset ownership before returning logs
//   - Returns 404 (not 403) to prevent information leakage
//
// Example Response:
//
//	{
//	  "logs": [
//	    {
//	      "id": "660e8400-e29b-41d4-a716-446655440000",
//	      "asset_id": "550e8400-e29b-41d4-a716-446655440000",
//	      "content": "Fixed nginx by restarting service",
//	      "tags": ["nginx", "fix"],
//	      "created_at": "2024-03-15T14:30:00Z",
//	      "updated_at": "2024-03-15T14:30:00Z"
//	    }
//	  ],
//	  "total": 1,
//	  "limit": 50,
//	  "offset": 0
//	}
func (h *LogHandler) ListByAsset(c echo.Context) error {
	// Extract user_id from context (set by auth middleware)
	userID, err := middleware.GetUserIDOrError(c)
	if err != nil {
		return err
	}

	// Parse and validate asset ID from URL parameter
	idParam := c.Param("id")
	assetID, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid asset id")
	}

	// Parse query parameters
	var params model.LogQueryParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	// Set defaults for pagination (logs use 50 as default)
	params.SetDefaults()

	// Call service (service verifies asset ownership)
	response, err := h.service.ListByAsset(c.Request().Context(), userID, assetID, &params)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, response)
}

// Create handles POST /api/v1/assets/:id/logs
//
// Creates a new log entry for the specified asset. This is a nested route that
// requires the asset ID in the URL path.
//
// Authentication: Required (user_id from context)
//
// URL Parameters:
//   - id: UUID of the asset (required)
//
// Request Body (JSON):
//   - content: Log content/description (required, 2-10000 chars)
//   - tags: Array of tag strings (optional, max 20 tags, each max 50 chars)
//
// Tags Handling:
//   - Service layer processes tags: trim, lowercase, deduplicate
//   - Tags can be null, empty array, or populated array
//   - Empty/null tags result in log with no tags
//
// Response:
//   - 201 Created: Returns LogResponse with created log including ID and timestamps
//   - 400 Bad Request: Invalid asset ID, JSON format, or validation error
//   - 401 Unauthorized: Missing or invalid authentication
//   - 404 Not Found: Asset doesn't exist or belongs to another user
//
// Security:
//   - Service layer verifies asset ownership before creating log
//   - User can only create logs for their own assets
//
// Example Request:
//
//	{
//	  "content": "Updated nginx config and restarted service",
//	  "tags": ["nginx", "config", "fix"]
//	}
//
// Example Response:
//
//	{
//	  "id": "660e8400-e29b-41d4-a716-446655440000",
//	  "asset_id": "550e8400-e29b-41d4-a716-446655440000",
//	  "user_id": "user-123",
//	  "content": "Updated nginx config and restarted service",
//	  "tags": ["nginx", "config", "fix"],
//	  "created_at": "2024-03-15T14:30:00Z",
//	  "updated_at": "2024-03-15T14:30:00Z"
//	}
func (h *LogHandler) Create(c echo.Context) error {
	// Extract user_id from context
	userID, err := middleware.GetUserIDOrError(c)
	if err != nil {
		return err
	}

	// Parse and validate asset ID from URL parameter
	idParam := c.Param("id")
	assetID, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid asset id")
	}

	// Parse request body
	var req model.CreateLogRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Call service (service verifies asset ownership and processes tags)
	response, err := h.service.Create(c.Request().Context(), userID, assetID, &req)
	if err != nil {
		return err
	}

	// Return response with 201 Created
	return c.JSON(http.StatusCreated, response)
}

// GetByID handles GET /api/v1/logs/:id
//
// Returns a single log by its ID. This is a flat route that accesses logs directly
// without requiring the asset ID in the URL.
//
// Authentication: Required (user_id from context)
//
// URL Parameters:
//   - id: UUID of the log (required)
//
// Response:
//   - 200 OK: Returns LogResponse with the requested log
//   - 400 Bad Request: Invalid log ID format
//   - 401 Unauthorized: Missing or invalid authentication
//   - 404 Not Found: Log doesn't exist or belongs to another user
//
// Security:
//   - Service layer filters by user_id to ensure users can only access their own logs
//   - Returns 404 (not 403) to prevent information leakage
func (h *LogHandler) GetByID(c echo.Context) error {
	// Extract user_id from context
	userID, err := middleware.GetUserIDOrError(c)
	if err != nil {
		return err
	}

	// Parse and validate log ID from URL parameter
	idParam := c.Param("id")
	logID, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid log id")
	}

	// Call service
	response, err := h.service.GetByID(c.Request().Context(), userID, logID)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, response)
}

// Update handles PATCH /api/v1/logs/:id
//
// Updates an existing log entry. This is a flat route with PATCH semantics,
// allowing partial updates of log fields.
//
// Authentication: Required (user_id from context)
//
// URL Parameters:
//   - id: UUID of the log (required)
//
// Request Body (JSON):
//   - content: New log content (optional, 2-10000 chars)
//   - tags: New tags array (optional, max 20 tags, each max 50 chars)
//
// PATCH Semantics:
//   - Omitted fields: Not updated (keep existing value)
//   - Provided fields: Updated to new value
//   - tags = []: Clear all tags
//   - tags = null or omitted: Keep existing tags
//
// Tags Handling:
//   - Service layer processes tags: trim, lowercase, deduplicate
//   - To clear tags, send empty array: {"tags": []}
//   - To keep tags, omit field or send null: {} or {"tags": null}
//
// Response:
//   - 200 OK: Returns LogResponse with updated log
//   - 400 Bad Request: Invalid log ID, JSON format, or validation error
//   - 401 Unauthorized: Missing or invalid authentication
//   - 404 Not Found: Log doesn't exist or belongs to another user
//
// Example Requests:
//
// Update content only:
//
//	{"content": "Updated description with more details"}
//
// Update tags only:
//
//	{"tags": ["nginx", "fix", "config"]}
//
// Clear tags:
//
//	{"tags": []}
//
// Update both:
//
//	{"content": "New content", "tags": ["new-tag"]}
func (h *LogHandler) Update(c echo.Context) error {
	// Extract user_id from context
	userID, err := middleware.GetUserIDOrError(c)
	if err != nil {
		return err
	}

	// Parse and validate log ID from URL parameter
	idParam := c.Param("id")
	logID, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid log id")
	}

	// Parse request body
	var req model.UpdateLogRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Call service (service processes tags if present)
	response, err := h.service.Update(c.Request().Context(), userID, logID, &req)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, response)
}

// Delete handles DELETE /api/v1/logs/:id
//
// Deletes a log entry. This operation is idempotent - deleting an already-deleted
// log returns 404.
//
// Authentication: Required (user_id from context)
//
// URL Parameters:
//   - id: UUID of the log (required)
//
// Response:
//   - 204 No Content: Log successfully deleted (empty response body)
//   - 400 Bad Request: Invalid log ID format
//   - 401 Unauthorized: Missing or invalid authentication
//   - 404 Not Found: Log doesn't exist or belongs to another user
//
// Security:
//   - Service layer filters by user_id to ensure users can only delete their own logs
//   - No cascade effects (logs have no dependent resources)
//
// Idempotency:
//   - Deleting the same log twice returns 404 on the second attempt
//   - This is standard RESTful behavior for DELETE operations
func (h *LogHandler) Delete(c echo.Context) error {
	// Extract user_id from context
	userID, err := middleware.GetUserIDOrError(c)
	if err != nil {
		return err
	}

	// Parse and validate log ID from URL parameter
	idParam := c.Param("id")
	logID, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid log id")
	}

	// Call service
	err = h.service.Delete(c.Request().Context(), userID, logID)
	if err != nil {
		return err
	}

	// Return 204 No Content
	return c.NoContent(http.StatusNoContent)
}
