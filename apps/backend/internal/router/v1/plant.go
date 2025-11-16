package v1

import (
	"github.com/labstack/echo/v4"
	"ark/internal/handler"
	"ark/internal/middleware"
)

func registerPlantRoutes(r *echo.Group, h *handler.Handlers, auth *middleware.AuthMiddleware) {
	// Plant operations
	plants := r.Group("/plants")

	// Phase 2: Apply RequireAuth to extract user claims from context
	// (Phase 1: ClerkAuthMiddleware already verified JWT at /api/v1 level)
	plants.Use(auth.RequireAuth)

	// Collection operations
	plants.POST("", h.Plant.CreatePlant)
	plants.GET("", h.Plant.GetPlants)

	// Individual plant operations
	dynamicPlant := plants.Group("/:id")
	dynamicPlant.GET("", h.Plant.GetPlantByID)
	dynamicPlant.PUT("", h.Plant.UpdatePlant)
	dynamicPlant.DELETE("", h.Plant.DeletePlant)
}
