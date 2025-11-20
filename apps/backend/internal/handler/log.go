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
