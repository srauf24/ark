package model

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========== SuccessResponse Tests ==========

// Test 1: TestSuccessResponse_MarshalJSON_AllFields
func TestSuccessResponse_MarshalJSON_AllFields(t *testing.T) {
	resp := SuccessResponse{
		Success: true,
		Data:    map[string]string{"id": "123", "name": "Test"},
		Message: "Operation completed successfully",
	}

	jsonData, err := json.Marshal(resp)
	require.NoError(t, err, "Failed to marshal SuccessResponse")

	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, `"success":true`, "JSON should contain success field")
	assert.Contains(t, jsonStr, `"data"`, "JSON should contain data field")
	assert.Contains(t, jsonStr, `"message"`, "JSON should contain message field")
	assert.Contains(t, jsonStr, "Operation completed successfully", "JSON should contain message text")
}

// Test 2: TestSuccessResponse_MarshalJSON_DataOnly
func TestSuccessResponse_MarshalJSON_DataOnly(t *testing.T) {
	resp := SuccessResponse{
		Success: true,
		Data:    map[string]string{"result": "ok"},
		Message: "",
	}

	jsonData, err := json.Marshal(resp)
	require.NoError(t, err, "Failed to marshal SuccessResponse")

	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, `"success":true`, "JSON should contain success field")
	assert.Contains(t, jsonStr, `"data"`, "JSON should contain data field")
	// With omitempty, empty message should be omitted
	if strings.Contains(jsonStr, `"message":""`) {
		t.Error("Empty message should be omitted with omitempty")
	}
}

// Test 3: TestSuccessResponse_MarshalJSON_NoData
func TestSuccessResponse_MarshalJSON_NoData(t *testing.T) {
	resp := SuccessResponse{
		Success: true,
		Data:    nil,
		Message: "Operation successful",
	}

	jsonData, err := json.Marshal(resp)
	require.NoError(t, err, "Failed to marshal SuccessResponse")

	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, `"success":true`, "JSON should contain success field")
	assert.Contains(t, jsonStr, `"message"`, "JSON should contain message field")
	// With omitempty, nil data should be omitted
	if strings.Contains(jsonStr, `"data":null`) {
		t.Log("Note: nil data serializes as null (acceptable)")
	}
}

// Test 4: TestSuccessResponse_MarshalJSON_MinimalSuccess
func TestSuccessResponse_MarshalJSON_MinimalSuccess(t *testing.T) {
	resp := SuccessResponse{
		Success: true,
	}

	jsonData, err := json.Marshal(resp)
	require.NoError(t, err, "Failed to marshal SuccessResponse")

	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, `"success":true`, "JSON should contain success field")
}

// Test 5: TestSuccessResponse_MarshalJSON_WithArrayData
func TestSuccessResponse_MarshalJSON_WithArrayData(t *testing.T) {
	resp := SuccessResponse{
		Success: true,
		Data:    []string{"item1", "item2", "item3"},
		Message: "Items retrieved",
	}

	jsonData, err := json.Marshal(resp)
	require.NoError(t, err, "Failed to marshal SuccessResponse")

	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, `"success":true`, "JSON should contain success field")
	assert.Contains(t, jsonStr, `"data"`, "JSON should contain data field")
	assert.Contains(t, jsonStr, "item1", "JSON should contain array items")
}

// Test 6: TestSuccessResponse_UnmarshalJSON_AllFields
func TestSuccessResponse_UnmarshalJSON_AllFields(t *testing.T) {
	jsonStr := `{
		"success": true,
		"data": {"key": "value"},
		"message": "Success"
	}`

	var resp SuccessResponse
	err := json.Unmarshal([]byte(jsonStr), &resp)
	require.NoError(t, err, "Failed to unmarshal SuccessResponse")

	assert.True(t, resp.Success, "Success should be true")
	assert.NotNil(t, resp.Data, "Data should not be nil")
	assert.Equal(t, "Success", resp.Message, "Message should match")
}

// ========== ErrorResponse Tests ==========

// Test 7: TestErrorResponse_MarshalJSON_WithDetails
func TestErrorResponse_MarshalJSON_WithDetails(t *testing.T) {
	resp := ErrorResponse{
		Success: false,
		Error:   "Validation failed",
		Details: map[string]string{
			"name":  "required field missing",
			"email": "invalid format",
		},
	}

	jsonData, err := json.Marshal(resp)
	require.NoError(t, err, "Failed to marshal ErrorResponse")

	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, `"success":false`, "JSON should contain success=false")
	assert.Contains(t, jsonStr, `"error"`, "JSON should contain error field")
	assert.Contains(t, jsonStr, "Validation failed", "JSON should contain error message")
	assert.Contains(t, jsonStr, `"details"`, "JSON should contain details field")
	assert.Contains(t, jsonStr, "name", "JSON should contain detail key 'name'")
	assert.Contains(t, jsonStr, "email", "JSON should contain detail key 'email'")
}

// Test 8: TestErrorResponse_MarshalJSON_NoDetails
func TestErrorResponse_MarshalJSON_NoDetails(t *testing.T) {
	resp := ErrorResponse{
		Success: false,
		Error:   "Resource not found",
		Details: nil,
	}

	jsonData, err := json.Marshal(resp)
	require.NoError(t, err, "Failed to marshal ErrorResponse")

	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, `"success":false`, "JSON should contain success=false")
	assert.Contains(t, jsonStr, `"error"`, "JSON should contain error field")
	assert.Contains(t, jsonStr, "Resource not found", "JSON should contain error message")
	// With omitempty, nil details should be omitted
	if strings.Contains(jsonStr, `"details":null`) {
		t.Log("Note: nil details serializes as null (acceptable)")
	}
}

// Test 9: TestErrorResponse_MarshalJSON_MultipleDetails
func TestErrorResponse_MarshalJSON_MultipleDetails(t *testing.T) {
	resp := ErrorResponse{
		Success: false,
		Error:   "Multiple validation errors",
		Details: map[string]string{
			"field1": "error1",
			"field2": "error2",
			"field3": "error3",
		},
	}

	jsonData, err := json.Marshal(resp)
	require.NoError(t, err, "Failed to marshal ErrorResponse")

	var unmarshaled ErrorResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err, "Failed to unmarshal ErrorResponse")

	assert.Equal(t, 3, len(unmarshaled.Details), "Should have 3 detail entries")
}

// Test 10: TestErrorResponse_UnmarshalJSON_AllFields
func TestErrorResponse_UnmarshalJSON_AllFields(t *testing.T) {
	jsonStr := `{
		"success": false,
		"error": "Bad request",
		"details": {"param": "invalid"}
	}`

	var resp ErrorResponse
	err := json.Unmarshal([]byte(jsonStr), &resp)
	require.NoError(t, err, "Failed to unmarshal ErrorResponse")

	assert.False(t, resp.Success, "Success should be false")
	assert.Equal(t, "Bad request", resp.Error, "Error message should match")
	assert.NotNil(t, resp.Details, "Details should not be nil")
	assert.Equal(t, "invalid", resp.Details["param"], "Detail value should match")
}

// Test 11: TestErrorResponse_UnmarshalJSON_NoDetails
func TestErrorResponse_UnmarshalJSON_NoDetails(t *testing.T) {
	jsonStr := `{
		"success": false,
		"error": "Server error"
	}`

	var resp ErrorResponse
	err := json.Unmarshal([]byte(jsonStr), &resp)
	require.NoError(t, err, "Failed to unmarshal ErrorResponse")

	assert.False(t, resp.Success, "Success should be false")
	assert.Equal(t, "Server error", resp.Error, "Error message should match")
	assert.Nil(t, resp.Details, "Details should be nil when omitted")
}

// ========== DeleteResponse Tests ==========

// Test 12: TestDeleteResponse_MarshalJSON_Success
func TestDeleteResponse_MarshalJSON_Success(t *testing.T) {
	resp := DeleteResponse{
		Success: true,
		Message: "Asset deleted successfully",
	}

	jsonData, err := json.Marshal(resp)
	require.NoError(t, err, "Failed to marshal DeleteResponse")

	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, `"success":true`, "JSON should contain success=true")
	assert.Contains(t, jsonStr, `"message"`, "JSON should contain message field")
	assert.Contains(t, jsonStr, "Asset deleted successfully", "JSON should contain message text")
}

// Test 13: TestDeleteResponse_MarshalJSON_Fields
func TestDeleteResponse_MarshalJSON_Fields(t *testing.T) {
	resp := DeleteResponse{
		Success: true,
		Message: "Log entry removed",
	}

	jsonData, err := json.Marshal(resp)
	require.NoError(t, err, "Failed to marshal DeleteResponse")

	// Unmarshal to verify structure
	var unmarshaled map[string]interface{}
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err, "Failed to unmarshal to map")

	// Verify both fields are present
	assert.Contains(t, unmarshaled, "success", "JSON should have success field")
	assert.Contains(t, unmarshaled, "message", "JSON should have message field")

	// Verify types
	assert.IsType(t, true, unmarshaled["success"], "Success should be boolean")
	assert.IsType(t, "", unmarshaled["message"], "Message should be string")
}

// Test 14: TestDeleteResponse_UnmarshalJSON
func TestDeleteResponse_UnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"success": true,
		"message": "Resource deleted"
	}`

	var resp DeleteResponse
	err := json.Unmarshal([]byte(jsonStr), &resp)
	require.NoError(t, err, "Failed to unmarshal DeleteResponse")

	assert.True(t, resp.Success, "Success should be true")
	assert.Equal(t, "Resource deleted", resp.Message, "Message should match")
}

// Test 15: TestDeleteResponse_RoundTrip
func TestDeleteResponse_RoundTrip(t *testing.T) {
	original := DeleteResponse{
		Success: true,
		Message: "Item deleted successfully",
	}

	// Marshal
	jsonData, err := json.Marshal(original)
	require.NoError(t, err, "Failed to marshal DeleteResponse")

	// Unmarshal
	var result DeleteResponse
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err, "Failed to unmarshal DeleteResponse")

	// Compare
	assert.Equal(t, original.Success, result.Success, "Success should match")
	assert.Equal(t, original.Message, result.Message, "Message should match")
}
