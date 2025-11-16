package v1

import (
	"github.com/labstack/echo/v4"
	"ark/internal/handler"
	"ark/internal/middleware"
)

func registerObservationRoutes(r *echo.Group, h *handler.Handlers, auth *middleware.AuthMiddleware) {
	// Observation operations
	observations := r.Group("/observations")

	// Phase 2: Apply RequireAuth to extract user claims from context
	// (Phase 1: ClerkAuthMiddleware already verified JWT at /api/v1 level)
	observations.Use(auth.RequireAuth)

	// Collection operations
	observations.POST("", h.Observation.CreateObservation)
	observations.GET("", h.Observation.GetObservations)

	// Individual observation operations
	dynamicObservation := observations.Group("/:id")
	dynamicObservation.GET("", h.Observation.GetObservationByID)
	dynamicObservation.PUT("", h.Observation.UpdateObservation)
	dynamicObservation.DELETE("", h.Observation.DeleteObservation)
}
