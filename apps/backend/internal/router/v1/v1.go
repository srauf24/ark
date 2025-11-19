package v1

import (
	"github.com/labstack/echo/v4"
	"ark/internal/handler"
	"ark/internal/middleware"
)

// RegisterRoutes registers all v1 API routes
func RegisterRoutes(router *echo.Echo, h *handler.Handlers, m *middleware.Middlewares) {
	v1 := router.Group("/api/v1")

	// Phase 1: Apply ClerkAuthMiddleware globally to all v1 routes
	// This verifies JWT tokens and stores session claims in context
	v1.Use(m.Auth.ClerkAuthMiddleware)

	// Register resource routes
	// Each protected route will use RequireAuth (Phase 2) to extract claims
	// TODO: Register asset and log routes when handlers are implemented
}
