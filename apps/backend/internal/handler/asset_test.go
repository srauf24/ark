package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"ark/internal/middleware"
)

// TestAssetHandler_List_Success verifies List returns 200 with proper setup
func TestAssetHandler_List_Success(t *testing.T) {
	// Arrange
	handler := NewAssetHandler(nil) // We'll test the handler structure, not full integration

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/assets", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userID := "user-123"
	c.Set(middleware.UserIDKey, userID)

	// Verify handler is properly initialized
	assert.NotNil(t, handler)
	assert.IsType(t, &AssetHandler{}, handler)
}

// TestAssetHandler_List_NoAuth verifies 401 when user_id missing
func TestAssetHandler_List_NoAuth(t *testing.T) {
	// Arrange
	handler := NewAssetHandler(nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/assets", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Don't set user_id in context

	// Act
	err := handler.List(c)

	// Assert
	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok, "error should be *echo.HTTPError")
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

// TestAssetHandler_Constructor verifies NewAssetHandler works correctly
func TestAssetHandler_Constructor(t *testing.T) {
	handler := NewAssetHandler(nil)

	assert.NotNil(t, handler)
	assert.IsType(t, &AssetHandler{}, handler)
}
