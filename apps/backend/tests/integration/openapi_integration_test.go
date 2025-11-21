package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAPI_UI_Returns200(t *testing.T) {
	// Change to backend root directory so static files can be found
	// The test runs in apps/backend/tests/integration, so we need to go up two levels
	err := os.Chdir("../..")
	require.NoError(t, err, "Should be able to change to backend root")
	defer os.Chdir("tests/integration") // Change back after test

	e := createTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "html", "Response should contain HTML")
	assert.Contains(t, rec.Body.String(), "scalar", "Response should contain scalar reference")
}

func TestOpenAPI_Spec_Returns200AndValidJSON(t *testing.T) {
	// Change to backend root directory so static files can be found
	err := os.Chdir("../..")
	require.NoError(t, err, "Should be able to change to backend root")
	defer os.Chdir("tests/integration") // Change back after test

	e := createTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/static/openapi.json", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	// Verify it's valid JSON
	var spec map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &spec)
	assert.NoError(t, err, "Response should be valid JSON")

	// Verify key fields
	info, ok := spec["info"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "ARK Asset Management API", info["title"])
}
