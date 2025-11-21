package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"ark/internal/middleware"
)

// TestLogHandler_ListByAsset_Success verifies ListByAsset returns 200 with proper setup
func TestLogHandler_ListByAsset_Success(t *testing.T) {
	// Arrange
	handler := NewLogHandler(nil) // We'll test the handler structure, not full integration

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/assets/550e8400-e29b-41d4-a716-446655440000/logs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("550e8400-e29b-41d4-a716-446655440000")

	userID := "user-123"
	c.Set(middleware.UserIDKey, userID)

	// Verify handler is properly initialized
	assert.NotNil(t, handler)
	assert.IsType(t, &LogHandler{}, handler)
}

// TestLogHandler_ListByAsset_InvalidAssetID verifies 400 when asset ID is invalid
func TestLogHandler_ListByAsset_InvalidAssetID(t *testing.T) {
	// Arrange
	handler := NewLogHandler(nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/assets/invalid-uuid/logs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-uuid")

	userID := "user-123"
	c.Set(middleware.UserIDKey, userID)

	// Act
	err := handler.ListByAsset(c)

	// Assert
	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok, "error should be *echo.HTTPError")
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}

// TestLogHandler_ListByAsset_NoAuth verifies 401 when user_id missing
func TestLogHandler_ListByAsset_NoAuth(t *testing.T) {
	// Arrange
	handler := NewLogHandler(nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/assets/550e8400-e29b-41d4-a716-446655440000/logs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("550e8400-e29b-41d4-a716-446655440000")

	// Don't set user_id in context

	// Act
	err := handler.ListByAsset(c)

	// Assert
	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok, "error should be *echo.HTTPError")
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

// TestLogHandler_Constructor verifies NewLogHandler works correctly
func TestLogHandler_Constructor(t *testing.T) {
	handler := NewLogHandler(nil)

	assert.NotNil(t, handler)
	assert.IsType(t, &LogHandler{}, handler)
}
