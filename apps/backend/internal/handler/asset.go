package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"ark/internal/middleware"
	"ark/internal/model"
	"ark/internal/service"
)

// AssetHandler handles HTTP requests for asset operations
type AssetHandler struct {
	service *service.AssetService
}

// NewAssetHandler creates a new AssetHandler with the given service
func NewAssetHandler(service *service.AssetService) *AssetHandler {
	return &AssetHandler{
		service: service,
	}
}

// List handles GET /api/v1/assets
// Returns a paginated list of assets for the authenticated user
func (h *AssetHandler) List(c echo.Context) error {
	// Extract user_id from context (set by auth middleware)
	userID, err := middleware.GetUserIDOrError(c)
	if err != nil {
		return err
	}

	// Parse query parameters
	var params model.AssetQueryParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	// Set defaults for pagination
	params.SetDefaults()

	// Call service
	response, err := h.service.List(c.Request().Context(), userID, &params)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, response)
}
