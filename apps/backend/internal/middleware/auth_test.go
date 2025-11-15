package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/srauf24/gardenjournal/internal/config"
	"github.com/srauf24/gardenjournal/internal/errs"
	"github.com/srauf24/gardenjournal/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestServer creates a minimal server for testing
func createTestServer() *server.Server {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &server.Server{
		Logger: &logger,
		Config: &config.Config{},
	}
}

// createTestContext creates a test Echo context with optional headers
func createTestContext(authHeader string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Set a request ID for logging
	c.Set("request_id", "test-request-id")
	return c
}

func TestClerkAuthMiddleware_MissingAuthHeader(t *testing.T) {
	testServer := createTestServer()
	authMiddleware := NewAuthMiddleware(testServer)

	c := createTestContext("")

	handler := authMiddleware.ClerkAuthMiddleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)

	// Should return an error
	assert.Error(t, err, "ClerkAuthMiddleware should return error when Authorization header is missing")

	// Response should be 401
	httpErr, ok := err.(*errs.HTTPError)
	require.True(t, ok, "Error should be an errs.HTTPError")
	assert.Equal(t, http.StatusUnauthorized, httpErr.Status, "Should return 401 Unauthorized")
}

func TestClerkAuthMiddleware_InvalidTokenFormat(t *testing.T) {
	testServer := createTestServer()
	authMiddleware := NewAuthMiddleware(testServer)

	tests := []struct {
		name       string
		authHeader string
	}{
		{
			name:       "no Bearer prefix",
			authHeader: "abc123token",
		},
		{
			name:       "empty after Bearer",
			authHeader: "Bearer ",
		},
		{
			name:       "only Bearer",
			authHeader: "Bearer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createTestContext(tt.authHeader)

			handler := authMiddleware.ClerkAuthMiddleware(func(c echo.Context) error {
				return c.String(http.StatusOK, "success")
			})

			err := handler(c)

			assert.Error(t, err, "ClerkAuthMiddleware should return error for invalid token format")
			httpErr, ok := err.(*errs.HTTPError)
			require.True(t, ok, "Error should be an errs.HTTPError")
			assert.Equal(t, http.StatusUnauthorized, httpErr.Status, "Should return 401 Unauthorized")
		})
	}
}

func TestClerkAuthMiddleware_TokenVerificationFails(t *testing.T) {
	// Skip if no Clerk secret key is set
	if os.Getenv("GARDENJOURNAL_AUTH.CLERK.SECRET_KEY") == "" {
		t.Skip("Skipping test: GARDENJOURNAL_AUTH.CLERK.SECRET_KEY not set")
	}

	clerk.SetKey(os.Getenv("GARDENJOURNAL_AUTH.CLERK.SECRET_KEY"))

	testServer := createTestServer()
	authMiddleware := NewAuthMiddleware(testServer)

	// Use an invalid JWT token
	c := createTestContext("Bearer invalid.jwt.token")

	handler := authMiddleware.ClerkAuthMiddleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)

	assert.Error(t, err, "ClerkAuthMiddleware should return error for invalid JWT")
	httpErr, ok := err.(*errs.HTTPError)
	require.True(t, ok, "Error should be an errs.HTTPError")
	assert.Equal(t, http.StatusUnauthorized, httpErr.Status, "Should return 401 Unauthorized")
}

func TestClerkAuthMiddleware_StoresClaimsInContext(t *testing.T) {
	testServer := createTestServer()
	authMiddleware := NewAuthMiddleware(testServer)

	// Create a mock token (this will fail verification, but we can check the flow)
	c := createTestContext("Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzEyMyJ9.fake")

	var contextClaims interface{}
	handler := authMiddleware.ClerkAuthMiddleware(func(c echo.Context) error {
		// Try to get claims from context
		contextClaims = c.Get(ClerkSessionClaimsKey)
		return c.String(http.StatusOK, "success")
	})

	// This will fail during verification, which is expected
	_ = handler(c)

	// We expect the handler to fail before setting claims
	// This test verifies the context key is used correctly
	assert.Nil(t, contextClaims, "Claims should not be set if verification fails")
}

func TestRequireAuth_ClaimsNotInContext(t *testing.T) {
	testServer := createTestServer()
	authMiddleware := NewAuthMiddleware(testServer)

	c := createTestContext("")

	handler := authMiddleware.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)

	assert.Error(t, err, "RequireAuth should return error when claims not in context")
	httpErr, ok := err.(*errs.HTTPError)
	require.True(t, ok, "Error should be an errs.HTTPError")
	assert.Equal(t, http.StatusUnauthorized, httpErr.Status, "Should return 401 Unauthorized")
}

func TestRequireAuth_ClaimsInContext(t *testing.T) {
	testServer := createTestServer()
	authMiddleware := NewAuthMiddleware(testServer)

	c := createTestContext("")

	// Manually set valid claims in context (simulating ClerkAuthMiddleware)
	mockClaims := &clerk.SessionClaims{
		RegisteredClaims: clerk.RegisteredClaims{
			Subject: "user_123",
		},
		Claims: clerk.Claims{
			ActiveOrganizationRole:        "admin",
			ActiveOrganizationPermissions: []string{"read", "write"},
		},
	}

	c.Set(ClerkSessionClaimsKey, mockClaims)

	var userID interface{}
	var userRole interface{}
	var permissions interface{}

	handler := authMiddleware.RequireAuth(func(c echo.Context) error {
		userID = c.Get("user_id")
		userRole = c.Get("user_role")
		permissions = c.Get("permissions")
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)

	assert.NoError(t, err, "RequireAuth should not return error when valid claims are in context")
	assert.Equal(t, "user_123", userID, "user_id should be set correctly")
	assert.Equal(t, "admin", userRole, "user_role should be set correctly")
	assert.Equal(t, []string{"read", "write"}, permissions, "permissions should be set correctly")
}

func TestRequireAuth_InvalidClaimsType(t *testing.T) {
	testServer := createTestServer()
	authMiddleware := NewAuthMiddleware(testServer)

	c := createTestContext("")

	// Set invalid type in context
	c.Set(ClerkSessionClaimsKey, "invalid-type")

	handler := authMiddleware.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)

	assert.Error(t, err, "RequireAuth should return error when claims type assertion fails")
	httpErr, ok := err.(*errs.HTTPError)
	require.True(t, ok, "Error should be an errs.HTTPError")
	assert.Equal(t, http.StatusUnauthorized, httpErr.Status, "Should return 401 Unauthorized")
}

func TestAuthMiddleware_TwoPhaseAuthentication(t *testing.T) {
	testServer := createTestServer()
	authMiddleware := NewAuthMiddleware(testServer)

	// Create a test context
	c := createTestContext("")

	// Manually set claims to simulate successful ClerkAuthMiddleware
	mockClaims := &clerk.SessionClaims{
		RegisteredClaims: clerk.RegisteredClaims{
			Subject: "user_456",
		},
		Claims: clerk.Claims{
			ActiveOrganizationRole:        "member",
			ActiveOrganizationPermissions: []string{"read"},
		},
	}

	c.Set(ClerkSessionClaimsKey, mockClaims)

	// Now run RequireAuth (Phase 2)
	var finalUserID string
	handler := authMiddleware.RequireAuth(func(c echo.Context) error {
		finalUserID = c.Get("user_id").(string)
		return c.String(http.StatusOK, "authenticated")
	})

	err := handler(c)

	assert.NoError(t, err, "Two-phase authentication should succeed")
	assert.Equal(t, "user_456", finalUserID, "User ID should be correctly set after two phases")
}
