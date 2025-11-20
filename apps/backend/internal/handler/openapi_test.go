package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 1: Verify OpenAPI file is valid JSON
func TestOpenAPISpecIsValidJSON(t *testing.T) {
	// Read the OpenAPI spec file
	data, err := os.ReadFile("../../static/openapi.json")
	require.NoError(t, err, "Should be able to read openapi.json file")

	// Verify it's valid JSON
	var spec map[string]interface{}
	err = json.Unmarshal(data, &spec)
	require.NoError(t, err, "OpenAPI spec should be valid JSON")

	// Verify it has required OpenAPI fields
	assert.Contains(t, spec, "openapi", "Should have openapi version field")
	assert.Contains(t, spec, "info", "Should have info field")
	assert.Contains(t, spec, "paths", "Should have paths field")
}

// Test 2: Verify OpenAPI spec has correct title
func TestOpenAPISpecTitle(t *testing.T) {
	// Read the OpenAPI spec file
	data, err := os.ReadFile("../../static/openapi.json")
	require.NoError(t, err, "Should be able to read openapi.json file")

	var spec map[string]interface{}
	err = json.Unmarshal(data, &spec)
	require.NoError(t, err, "OpenAPI spec should be valid JSON")

	// Check info section
	info, ok := spec["info"].(map[string]interface{})
	require.True(t, ok, "info should be an object")

	// Verify title
	title, ok := info["title"].(string)
	require.True(t, ok, "title should be a string")
	assert.Equal(t, "ARK Asset Management API", title, "Title should be 'ARK Asset Management API'")

	// Verify description mentions ARK/homelab
	description, ok := info["description"].(string)
	require.True(t, ok, "description should be a string")
	assert.Contains(t, strings.ToLower(description), "ark", "Description should mention ARK")
	assert.Contains(t, strings.ToLower(description), "homelab", "Description should mention homelab")
}

// Test 3: Verify no legacy references in spec
func TestOpenAPISpecNoLegacyReferences(t *testing.T) {
	// Read the OpenAPI spec file
	data, err := os.ReadFile("../../static/openapi.json")
	require.NoError(t, err, "Should be able to read openapi.json file")

	specContent := strings.ToLower(string(data))

	// Check for legacy terms (case-insensitive)
	assert.NotContains(t, specContent, "garden", "Should not contain 'garden'")
	assert.NotContains(t, specContent, "plant", "Should not contain 'plant'")
	assert.NotContains(t, specContent, "observation", "Should not contain 'observation'")
	assert.NotContains(t, specContent, "/api/v1/plants", "Should not contain /api/v1/plants endpoint")
	assert.NotContains(t, specContent, "/api/v1/observations", "Should not contain /api/v1/observations endpoint")
}

// Test 4: Verify all expected endpoints are documented
func TestOpenAPISpecHasRequiredEndpoints(t *testing.T) {
	// Read the OpenAPI spec file
	data, err := os.ReadFile("../../static/openapi.json")
	require.NoError(t, err, "Should be able to read openapi.json file")

	var spec map[string]interface{}
	err = json.Unmarshal(data, &spec)
	require.NoError(t, err, "OpenAPI spec should be valid JSON")

	paths, ok := spec["paths"].(map[string]interface{})
	require.True(t, ok, "paths should be an object")

	// Required endpoints
	requiredEndpoints := []string{
		"/api/status",              // Health check
		"/api/v1/assets",           // Asset list/create
		"/api/v1/assets/{id}",      // Asset get/update/delete
		"/api/v1/assets/{id}/logs", // Log list/create for asset
		"/api/v1/logs/{id}",        // Log get/update/delete
	}

	for _, endpoint := range requiredEndpoints {
		assert.Contains(t, paths, endpoint, "Should document endpoint: %s", endpoint)
	}
}

// Test 5: Verify handler returns correct content type
func TestOpenAPIHandler_ContentType(t *testing.T) {
	// Create handler
	handler := NewOpenAPIHandler()

	// Create test request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call handler
	err := handler.ServeOpenAPISpec(c)
	require.NoError(t, err, "Handler should not return error")

	// Verify content type
	assert.Equal(t, "application/json; charset=UTF-8", rec.Header().Get("Content-Type"))
}

// Test 6: Verify handler returns 200 status
func TestOpenAPIHandler_StatusCode(t *testing.T) {
	// Create handler
	handler := NewOpenAPIHandler()

	// Create test request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call handler
	err := handler.ServeOpenAPISpec(c)
	require.NoError(t, err, "Handler should not return error")

	// Verify status code
	assert.Equal(t, http.StatusOK, rec.Code)
}

// Test 7: Verify handler returns valid JSON
func TestOpenAPIHandler_ValidJSON(t *testing.T) {
	// Create handler
	handler := NewOpenAPIHandler()

	// Create test request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call handler
	err := handler.ServeOpenAPISpec(c)
	require.NoError(t, err, "Handler should not return error")

	// Verify response is valid JSON
	var spec map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &spec)
	require.NoError(t, err, "Response should be valid JSON")
}

// Test 8: Verify handler response matches static file
func TestOpenAPIHandler_MatchesStaticFile(t *testing.T) {
	// Read static file
	staticData, err := os.ReadFile("../../static/openapi.json")
	require.NoError(t, err, "Should be able to read static file")

	var staticSpec map[string]interface{}
	err = json.Unmarshal(staticData, &staticSpec)
	require.NoError(t, err, "Static file should be valid JSON")

	// Create handler
	handler := NewOpenAPIHandler()

	// Create test request
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call handler
	err = handler.ServeOpenAPISpec(c)
	require.NoError(t, err, "Handler should not return error")

	// Parse handler response
	var handlerSpec map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &handlerSpec)
	require.NoError(t, err, "Handler response should be valid JSON")

	// Verify key fields match
	staticInfo := staticSpec["info"].(map[string]interface{})
	handlerInfo := handlerSpec["info"].(map[string]interface{})
	assert.Equal(t, staticInfo["title"], handlerInfo["title"], "Titles should match")
	assert.Equal(t, staticInfo["version"], handlerInfo["version"], "Versions should match")
}
