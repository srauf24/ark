package v1

import (
	"ark/internal/handler"
	"ark/internal/middleware"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all v1 API routes
func RegisterRoutes(router *echo.Echo, h *handler.Handlers, m *middleware.Middlewares) {
	v1 := router.Group("/api/v1")

	// Phase 1: Apply ClerkAuthMiddleware globally to all v1 routes
	// This verifies JWT tokens and stores session claims in context
	v1.Use(m.Auth.ClerkAuthMiddleware)

	// Asset routes
	assets := v1.Group("/assets")
	assets.GET("", h.Asset.List)          // GET /api/v1/assets
	assets.POST("", h.Asset.Create)       // POST /api/v1/assets
	assets.GET("/:id", h.Asset.GetByID)   // GET /api/v1/assets/:id
	assets.PATCH("/:id", h.Asset.Update)  // PATCH /api/v1/assets/:id
	assets.DELETE("/:id", h.Asset.Delete) // DELETE /api/v1/assets/:id

	// Log routes (nested under assets for create/list)
	assets.POST("/:id/logs", h.Log.Create)     // POST /api/v1/assets/:id/logs
	assets.GET("/:id/logs", h.Log.ListByAsset) // GET /api/v1/assets/:id/logs

	// Log routes (flat for direct access)
	logs := v1.Group("/logs")
	logs.GET("/:id", h.Log.GetByID)   // GET /api/v1/logs/:id
	logs.PATCH("/:id", h.Log.Update)  // PATCH /api/v1/logs/:id
	logs.DELETE("/:id", h.Log.Delete) // DELETE /api/v1/logs/:id
}
