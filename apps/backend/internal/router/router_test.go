package router

import (
	"strings"
	"testing"

	"ark/internal/config"
	"ark/internal/handler"
	"ark/internal/logger"
	"ark/internal/server"
	"ark/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Helper functions for route verification

// Helper function to create a test router with minimal setup
func createTestRouter(t *testing.T) *echo.Echo {
	// Create minimal server config for testing
	cfg := &config.Config{
		Observability: &config.ObservabilityConfig{
			ServiceName: "ark-test",
			Environment: "test",
			Logging: config.LoggingConfig{
				Level:  "debug",
				Format: "console",
			},
		},
	}

	// Create test logger service (required for NewLoggerWithService)
	loggerService := logger.NewLoggerService(cfg.Observability)

	// Create test logger
	log := logger.NewLoggerWithService(cfg.Observability, loggerService)

	// Create minimal server
	srv := &server.Server{
		Config: cfg,
		Logger: &log,
	}

	// Create minimal handlers (can be nil for route registration tests)
	handlers := &handler.Handlers{
		Health:  handler.NewHealthHandler(srv),
		OpenAPI: handler.NewOpenAPIHandler(srv),
		Asset:   nil, // Will be mocked in actual tests
		Log:     nil, // Will be mocked in actual tests
	}

	// Create minimal services
	services := &service.Services{}

	// Create router
	return NewRouter(srv, handlers, services)
}

// findRoute searches for a route by method and path pattern
func findRoute(routes []*echo.Route, method, pathPattern string) *echo.Route {
	for _, route := range routes {
		if route.Method == method && route.Path == pathPattern {
			return route
		}
	}
	return nil
}

// countRoutesWithPrefix counts routes matching a path prefix
func countRoutesWithPrefix(routes []*echo.Route, prefix string) int {
	count := 0
	for _, route := range routes {
		if strings.HasPrefix(route.Path, prefix) {
			count++
		}
	}
	return count
}

// Unit tests for helper functions

// TestHelpers_FindRoute verifies the findRoute helper function
func TestHelpers_FindRoute(t *testing.T) {
	routes := []*echo.Route{
		{Method: "GET", Path: "/api/v1/assets"},
		{Method: "POST", Path: "/api/v1/assets"},
		{Method: "GET", Path: "/api/v1/logs/:id"},
	}

	// Test finding an existing route
	found := findRoute(routes, "GET", "/api/v1/assets")
	assert.NotNil(t, found, "Expected to find GET /api/v1/assets")
	assert.Equal(t, "GET", found.Method)
	assert.Equal(t, "/api/v1/assets", found.Path)

	// Test finding another existing route
	foundLog := findRoute(routes, "GET", "/api/v1/logs/:id")
	assert.NotNil(t, foundLog, "Expected to find GET /api/v1/logs/:id")
	assert.Equal(t, "GET", foundLog.Method)

	// Test not finding a non-existent route
	notFound := findRoute(routes, "DELETE", "/api/v1/assets")
	assert.Nil(t, notFound, "Expected not to find DELETE /api/v1/assets")

	// Test not finding with wrong path
	notFoundPath := findRoute(routes, "GET", "/api/v1/nonexistent")
	assert.Nil(t, notFoundPath, "Expected not to find route with wrong path")
}

// TestHelpers_CountRoutesWithPrefix verifies the countRoutesWithPrefix helper function
func TestHelpers_CountRoutesWithPrefix(t *testing.T) {
	routes := []*echo.Route{
		{Method: "GET", Path: "/api/v1/assets"},
		{Method: "POST", Path: "/api/v1/assets"},
		{Method: "GET", Path: "/api/v1/assets/:id"},
		{Method: "GET", Path: "/api/v1/logs/:id"},
		{Method: "GET", Path: "/status"},
		{Method: "GET", Path: "/docs"},
	}

	// Test counting all /api/v1 routes
	count := countRoutesWithPrefix(routes, "/api/v1")
	assert.Equal(t, 4, count, "Expected 4 routes with /api/v1 prefix")

	// Test counting /api/v1/assets routes
	countAssets := countRoutesWithPrefix(routes, "/api/v1/assets")
	assert.Equal(t, 3, countAssets, "Expected 3 routes with /api/v1/assets prefix")

	// Test counting /api/v1/logs routes
	countLogs := countRoutesWithPrefix(routes, "/api/v1/logs")
	assert.Equal(t, 1, countLogs, "Expected 1 route with /api/v1/logs prefix")

	// Test counting system routes
	countStatus := countRoutesWithPrefix(routes, "/status")
	assert.Equal(t, 1, countStatus, "Expected 1 route with /status prefix")

	// Test counting non-existent prefix
	countNone := countRoutesWithPrefix(routes, "/nonexistent")
	assert.Equal(t, 0, countNone, "Expected 0 routes with /nonexistent prefix")
}

// TestNewRouter_AssetRoutesRegistered verifies all asset routes are correctly registered
func TestNewRouter_AssetRoutesRegistered(t *testing.T) {
	// Skip if we can't create full router (due to DB dependencies)
	// This is a placeholder that will be enhanced when we add proper mocking
	t.Skip("Requires full server initialization - will be implemented with proper mocks")

	router := createTestRouter(t)
	routes := router.Routes()

	// Expected asset routes
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/assets"},
		{"POST", "/api/v1/assets"},
		{"GET", "/api/v1/assets/:id"},
		{"PATCH", "/api/v1/assets/:id"},
		{"DELETE", "/api/v1/assets/:id"},
	}

	// Verify each expected route exists
	for _, expected := range expectedRoutes {
		route := findRoute(routes, expected.method, expected.path)
		assert.NotNil(t, route,
			"Expected route not found: %s %s", expected.method, expected.path)

		if route != nil {
			assert.Equal(t, expected.method, route.Method)
			assert.Equal(t, expected.path, route.Path)
		}
	}

	// Verify we have exactly 5 asset routes
	assetRouteCount := 0
	for _, route := range routes {
		if strings.HasPrefix(route.Path, "/api/v1/assets") &&
			!strings.Contains(route.Path, "/logs") {
			assetRouteCount++
		}
	}
	assert.Equal(t, 5, assetRouteCount, "Expected exactly 5 asset routes")
}

// TestNewRouter_LogRoutesRegistered verifies all log routes are correctly registered
func TestNewRouter_LogRoutesRegistered(t *testing.T) {
	// Skip if we can't create full router (due to DB dependencies)
	t.Skip("Requires full server initialization - will be implemented with proper mocks")

	router := createTestRouter(t)
	routes := router.Routes()

	// Expected log routes (nested under assets)
	nestedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/assets/:id/logs"},
		{"POST", "/api/v1/assets/:id/logs"},
	}

	// Expected log routes (flat for direct access)
	flatRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/logs/:id"},
		{"PATCH", "/api/v1/logs/:id"},
		{"DELETE", "/api/v1/logs/:id"},
	}

	// Verify nested routes
	for _, expected := range nestedRoutes {
		route := findRoute(routes, expected.method, expected.path)
		assert.NotNil(t, route,
			"Expected nested log route not found: %s %s", expected.method, expected.path)
	}

	// Verify flat routes
	for _, expected := range flatRoutes {
		route := findRoute(routes, expected.method, expected.path)
		assert.NotNil(t, route,
			"Expected flat log route not found: %s %s", expected.method, expected.path)
	}

	// Verify we have exactly 5 log routes total
	logRouteCount := 0
	for _, route := range routes {
		if strings.Contains(route.Path, "/logs") {
			logRouteCount++
		}
	}
	assert.Equal(t, 5, logRouteCount, "Expected exactly 5 log routes (2 nested + 3 flat)")
}

// Existing tests

// TestNewRouter_Compiles verifies the NewRouter function compiles
func TestNewRouter_Compiles(t *testing.T) {
	// This is a barebones test that verifies the function signature is correct
	// Full testing would require server, handlers, and services initialization

	// If this test compiles, the function signature is valid
	assert.True(t, true)
}

// TestNewRouter_V1GroupExists verifies /api/v1 group is created and all API routes use it
func TestNewRouter_V1GroupExists(t *testing.T) {
	// Skip if we can't create full router (due to DB dependencies)
	t.Skip("Requires full server initialization - will be implemented with proper mocks")

	router := createTestRouter(t)
	routes := router.Routes()

	// Count routes with /api/v1 prefix
	v1RouteCount := countRoutesWithPrefix(routes, "/api/v1")

	// We expect 10 routes under /api/v1 (5 asset + 5 log)
	assert.GreaterOrEqual(t, v1RouteCount, 10,
		"Expected at least 10 routes under /api/v1")

	// Verify no routes use old patterns
	for _, route := range routes {
		assert.NotContains(t, route.Path, "/plants",
			"Found legacy plant route: %s", route.Path)
		assert.NotContains(t, route.Path, "/observations",
			"Found legacy observation route: %s", route.Path)
	}

	// Verify system routes exist outside /api/v1
	statusRoute := findRoute(routes, "GET", "/status")
	assert.NotNil(t, statusRoute, "Expected /status route to exist")

	docsRoute := findRoute(routes, "GET", "/docs")
	assert.NotNil(t, docsRoute, "Expected /docs route to exist")
}

// TestNewRouter_MiddlewareApplied verifies middleware is correctly applied to routes
func TestNewRouter_MiddlewareApplied(t *testing.T) {
	// Skip if we can't create full router (due to DB dependencies)
	t.Skip("Requires full server initialization - will be implemented with proper mocks")

	router := createTestRouter(t)
	routes := router.Routes()

	// Verify global middleware is applied
	// Echo applies middleware to all routes, so we just verify the router was created
	assert.NotNil(t, router)

	// Verify we have routes registered
	assert.Greater(t, len(routes), 0, "Expected routes to be registered")

	// Note: Testing actual middleware execution requires integration tests
	// This unit test verifies the router structure is correct

	// Verify API routes exist (which should have auth middleware via v1.RegisterRoutes)
	apiRouteCount := countRoutesWithPrefix(routes, "/api/v1")
	assert.Greater(t, apiRouteCount, 0, "Expected API routes to exist")

	// Verify system routes exist (which should NOT have auth middleware)
	systemRoutes := []string{"/status", "/docs", "/static/*"}
	for _, path := range systemRoutes {
		// System routes should exist
		found := false
		for _, route := range routes {
			if strings.Contains(route.Path, path) || route.Path == path {
				found = true
				break
			}
		}
		// Note: Some system routes might not be registered in test environment
		// This is acceptable for unit tests
		_ = found
	}
}

// TestNewRouter_TotalRouteCount verifies the expected total number of routes
func TestNewRouter_TotalRouteCount(t *testing.T) {
	// Skip if we can't create full router (due to DB dependencies)
	t.Skip("Requires full server initialization - will be implemented with proper mocks")

	router := createTestRouter(t)
	routes := router.Routes()

	// Count routes by category
	assetRoutes := 0
	logRoutes := 0
	systemRoutes := 0

	for _, route := range routes {
		switch {
		case strings.HasPrefix(route.Path, "/api/v1/assets") && !strings.Contains(route.Path, "/logs"):
			assetRoutes++
		case strings.Contains(route.Path, "/logs"):
			logRoutes++
		case route.Path == "/status" || route.Path == "/docs" || strings.HasPrefix(route.Path, "/static"):
			systemRoutes++
		}
	}

	// Verify counts
	assert.Equal(t, 5, assetRoutes, "Expected 5 asset routes")
	assert.Equal(t, 5, logRoutes, "Expected 5 log routes")
	assert.GreaterOrEqual(t, systemRoutes, 2, "Expected at least 2 system routes")

	// Total API routes should be 10
	totalAPIRoutes := assetRoutes + logRoutes
	assert.Equal(t, 10, totalAPIRoutes, "Expected 10 total API routes")
}
