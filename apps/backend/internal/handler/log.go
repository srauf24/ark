package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"ark/internal/middleware"
	"ark/internal/model"
	"ark/internal/service"
)

// LogHandler handles HTTP requests for asset log operations
type LogHandler struct {
	service *service.LogService
}

// NewLogHandler creates a new LogHandler with the given service
func NewLogHandler(service *service.LogService) *LogHandler {
	return &LogHandler{
		service: service,
	}
}

// ListByAsset handles GET /api/v1/assets/:id/logs
// Returns a paginated list of logs for the specified asset
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
// Creates a new log for the specified asset
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
// Returns a single log by ID for the authenticated user
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
// Updates an existing log for the authenticated user
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
// Deletes a log for the authenticated user
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
