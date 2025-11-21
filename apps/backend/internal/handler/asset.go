package handler

import (
	"net/http"

	"github.com/google/uuid"
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

// Create handles POST /api/v1/assets
// Creates a new asset for the authenticated user
func (h *AssetHandler) Create(c echo.Context) error {
	// Extract user_id from context
	userID, err := middleware.GetUserIDOrError(c)
	if err != nil {
		return err
	}

	// Parse request body
	var req model.CreateAssetRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Call service
	response, err := h.service.Create(c.Request().Context(), userID, &req)
	if err != nil {
		return err
	}

	// Return response with 201 Created
	return c.JSON(http.StatusCreated, response)
}

// GetByID handles GET /api/v1/assets/:id
// Returns a single asset by ID for the authenticated user
func (h *AssetHandler) GetByID(c echo.Context) error {
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

	// Call service
	response, err := h.service.GetByID(c.Request().Context(), userID, assetID)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, response)
}

// Update handles PATCH /api/v1/assets/:id
// Updates an existing asset for the authenticated user
func (h *AssetHandler) Update(c echo.Context) error {
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
	var req model.UpdateAssetRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Call service
	response, err := h.service.Update(c.Request().Context(), userID, assetID, &req)
	if err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, response)
}

// Delete handles DELETE /api/v1/assets/:id
// Deletes an asset for the authenticated user
func (h *AssetHandler) Delete(c echo.Context) error {
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

	// Call service
	err = h.service.Delete(c.Request().Context(), userID, assetID)
	if err != nil {
		return err
	}

	// Return 204 No Content
	return c.NoContent(http.StatusNoContent)
}
