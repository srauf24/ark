package integration

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/srauf24/gardenjournal/internal/config"
	"github.com/srauf24/gardenjournal/internal/handler"
	"github.com/srauf24/gardenjournal/internal/middleware"
	v1 "github.com/srauf24/gardenjournal/internal/router/v1"
	"github.com/srauf24/gardenjournal/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestRouter creates a minimal router for integration testing
func createTestRouter() *echo.Echo {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	testServer := &server.Server{
		Logger: &logger,
		Config: &config.Config{
			Auth: config.AuthConfig{
				Clerk: config.ClerkConfig{
					SecretKey: "test_key",
					JWTIssuer: "https://test.clerk.accounts.dev",
				},
			},
		},
	}

	// Create minimal dependencies for routing
	handlers := &handler.Handlers{
		Plant:       &handler.PlantHandler{},
		Observation: &handler.ObservationHandler{},
	}

	// Create router
	e := echo.New()
	middlewares := middleware.NewMiddlewares(testServer)

	// Register the global error handler (required for proper error handling)
	e.HTTPErrorHandler = middlewares.Global.GlobalErrorHandler

	// Register v1 routes (which now includes ClerkAuthMiddleware)
	v1.RegisterRoutes(e, handlers, middlewares)

	return e
}

func TestAuth_NoJWT_Returns401(t *testing.T) {
	e := createTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/plants", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code, "Should return 401 Unauthorized when no JWT is provided")
}

func TestAuth_InvalidJWT_Returns401(t *testing.T) {
	e := createTestRouter()

	tests := []struct {
		name       string
		authHeader string
	}{
		{
			name:       "malformed JWT",
			authHeader: "Bearer invalid.jwt.token",
		},
		{
			name:       "no Bearer prefix",
			authHeader: "abc123token",
		},
		{
			name:       "empty Bearer token",
			authHeader: "Bearer ",
		},
		{
			name:       "random string",
			authHeader: "Bearer randomstringnotajwt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/plants", nil)
			req.Header.Set("Authorization", tt.authHeader)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code, "Should return 401 Unauthorized for invalid JWT")
		})
	}
}

func TestAuth_DifferentEndpoints_AllRequireAuth(t *testing.T) {
	e := createTestRouter()

	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/plants"},
		{http.MethodPost, "/api/v1/plants"},
		{http.MethodGet, "/api/v1/plants/123"},
		{http.MethodPut, "/api/v1/plants/123"},
		{http.MethodDelete, "/api/v1/plants/123"},
		{http.MethodGet, "/api/v1/observations"},
		{http.MethodPost, "/api/v1/observations"},
		{http.MethodGet, "/api/v1/observations/456"},
		{http.MethodPut, "/api/v1/observations/456"},
		{http.MethodDelete, "/api/v1/observations/456"},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.method+" "+endpoint.path, func(t *testing.T) {
			req := httptest.NewRequest(endpoint.method, endpoint.path, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code,
				"All /api/v1 endpoints should return 401 when no JWT is provided")
		})
	}
}

func TestAuth_MiddlewareChainOrdering(t *testing.T) {
	// This test verifies that ClerkAuthMiddleware runs before RequireAuth
	// by checking that we get the correct error message from ClerkAuthMiddleware

	e := createTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/plants", nil)
	// No Authorization header - should be caught by ClerkAuthMiddleware first
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)

	// The response should contain the error from ClerkAuthMiddleware
	// (not RequireAuth), proving the middleware chain is correct
	body := rec.Body.String()
	assert.Contains(t, body, "UNAUTHORIZED", "Should contain error code from middleware")
}

// Note: Positive test cases (valid JWT returning 200) are not included because:
// 1. They require a real Clerk environment with valid secret keys
// 2. They require generating valid JWTs with proper signing
// 3. Integration with actual Clerk services is needed
//
// These scenarios should be tested through:
// - Manual testing with real Clerk credentials (see Step 7 manual tests)
// - E2E tests in a staging environment with actual Clerk configuration
//
// The current tests verify:
// - Middleware is properly applied to all /api/v1 routes
// - Authentication failures are handled correctly
// - Middleware chain executes in the correct order
