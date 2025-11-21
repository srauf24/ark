package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestGetUserIDOrError_Success verifies GetUserIDOrError returns user_id when present
func TestGetUserIDOrError_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedUserID := "user-123"
	c.Set(UserIDKey, expectedUserID)

	// Act
	userID, err := GetUserIDOrError(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)
}

// TestGetUserIDOrError_Missing verifies GetUserIDOrError returns error when user_id not in context
func TestGetUserIDOrError_Missing(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Don't set user_id in context

	// Act
	userID, err := GetUserIDOrError(c)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", userID)

	// Verify it's an HTTPError with 401 status
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok, "error should be *echo.HTTPError")
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

// TestGetUserIDOrError_Empty verifies GetUserIDOrError returns error when user_id is empty
func TestGetUserIDOrError_Empty(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set(UserIDKey, "") // Set empty user_id

	// Act
	userID, err := GetUserIDOrError(c)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", userID)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

// TestGetUserID_Success verifies GetUserID returns user_id when present
func TestGetUserID_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedUserID := "user-456"
	c.Set(UserIDKey, expectedUserID)

	// Act
	userID := GetUserID(c)

	// Assert
	assert.Equal(t, expectedUserID, userID)
}

// TestGetUserID_Missing verifies GetUserID returns empty string when user_id not in context
func TestGetUserID_Missing(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	userID := GetUserID(c)

	// Assert
	assert.Equal(t, "", userID)
}
