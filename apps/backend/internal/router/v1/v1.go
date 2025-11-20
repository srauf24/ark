package v1

import (
	"ark/internal/handler"
	"ark/internal/middleware"

	"github.com/labstack/echo/v4"
)

// Package v1 implements API version 1 route registration.
//
// Authentication Strategy:
// This package uses a two-phase authentication pattern:
//
// Phase 1 (ClerkAuthMiddleware): Applied globally to all /api/v1 routes
//   - Extracts JWT from Authorization header
//   - Verifies token with Clerk SDK
//   - Stores session claims in Echo context
//   - Returns 401 if verification fails
//
// Phase 2 (RequireAuth): Applied per route or route group as needed
//   - Retrieves verified claims from context
//   - Extracts user_id and sets in context
//   - All database queries scoped to user_id (multi-tenancy)
//
// Route Structure:
//   - Asset routes: /api/v1/assets (collection and individual operations)
//   - Log routes: /api/v1/assets/:id/logs (nested for create/list)
//                 /api/v1/logs/:id (flat for individual operations)
//
// All routes require authentication via ClerkAuthMiddleware.

// RegisterRoutes registers all v1 API routes
func RegisterRoutes(router *echo.Echo, h *handler.Handlers, m *middleware.Middlewares) {
	// Create API v1 group - all routes will be prefixed with /api/v1
	v1 := router.Group("/api/v1")

	// Phase 1: Apply ClerkAuthMiddleware globally to all v1 routes
	// This verifies JWT tokens and stores session claims in context
	v1.Use(m.Auth.ClerkAuthMiddleware)

	// Asset routes - RESTful CRUD operations
	// All operations scoped to authenticated user via middleware
	assets := v1.Group("/assets")
	assets.GET("", h.Asset.List)          // GET /api/v1/assets - List user's assets
	assets.POST("", h.Asset.Create)       // POST /api/v1/assets - Create new asset
	assets.GET("/:id", h.Asset.GetByID)   // GET /api/v1/assets/:id - Get single asset
	assets.PATCH("/:id", h.Asset.Update)  // PATCH /api/v1/assets/:id - Update asset
	assets.DELETE("/:id", h.Asset.Delete) // DELETE /api/v1/assets/:id - Delete asset

	// Log routes (nested under assets for create/list)
	// These routes require asset_id in URL path
	assets.POST("/:id/logs", h.Log.Create)     // POST /api/v1/assets/:id/logs - Create log for asset
	assets.GET("/:id/logs", h.Log.ListByAsset) // GET /api/v1/assets/:id/logs - List logs for asset

	// Log routes (flat for direct access)
	// These routes operate on logs by log_id
	logs := v1.Group("/logs")
	logs.GET("/:id", h.Log.GetByID)   // GET /api/v1/logs/:id - Get single log
	logs.PATCH("/:id", h.Log.Update)  // PATCH /api/v1/logs/:id - Update log
	logs.DELETE("/:id", h.Log.Delete) // DELETE /api/v1/logs/:id - Delete log
}
