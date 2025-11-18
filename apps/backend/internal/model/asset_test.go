package model

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

// ========== Asset JSON Marshaling Tests ==========

// Test 1: TestAsset_JSONMarshaling_AllFields
func TestAsset_JSONMarshaling_AllFields(t *testing.T) {
	assetType := "server"
	hostname := "homelab-01"
	metadata := json.RawMessage(`{"cpu":"4 cores","ram":"16GB"}`)

	asset := Asset{
		ID:        uuid.New(),
		UserID:    "user_123",
		Name:      "My Server",
		Type:      &assetType,
		Hostname:  &hostname,
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(asset)
	if err != nil {
		t.Fatalf("Failed to marshal asset: %v", err)
	}

	// Verify all fields are present in JSON
	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, `"id"`) {
		t.Error("JSON should contain 'id' field")
	}
	if !strings.Contains(jsonStr, `"user_id"`) {
		t.Error("JSON should contain 'user_id' field")
	}
	if !strings.Contains(jsonStr, `"name"`) {
		t.Error("JSON should contain 'name' field")
	}
	if !strings.Contains(jsonStr, `"type"`) {
		t.Error("JSON should contain 'type' field")
	}
	if !strings.Contains(jsonStr, `"hostname"`) {
		t.Error("JSON should contain 'hostname' field")
	}
	if !strings.Contains(jsonStr, `"metadata"`) {
		t.Error("JSON should contain 'metadata' field")
	}
	if !strings.Contains(jsonStr, `"created_at"`) {
		t.Error("JSON should contain 'created_at' field")
	}
	if !strings.Contains(jsonStr, `"updated_at"`) {
		t.Error("JSON should contain 'updated_at' field")
	}
}

// Test 2: TestAsset_JSONMarshaling_NilOptionalFields
func TestAsset_JSONMarshaling_NilOptionalFields(t *testing.T) {
	asset := Asset{
		ID:        uuid.New(),
		UserID:    "user_123",
		Name:      "My Server",
		Type:      nil,
		Hostname:  nil,
		Metadata:  nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(asset)
	if err != nil {
		t.Fatalf("Failed to marshal asset: %v", err)
	}

	jsonStr := string(jsonData)
	// With omitempty, nil pointer fields should be omitted
	if strings.Contains(jsonStr, `"type"`) {
		t.Error("JSON should not contain 'type' field when nil")
	}
	if strings.Contains(jsonStr, `"hostname"`) {
		t.Error("JSON should not contain 'hostname' field when nil")
	}
	if strings.Contains(jsonStr, `"metadata"`) {
		t.Error("JSON should not contain 'metadata' field when nil")
	}
}

// Test 3: TestAsset_JSONMarshaling_EmptyMetadata
func TestAsset_JSONMarshaling_EmptyMetadata(t *testing.T) {
	metadata := json.RawMessage(`{}`)
	asset := Asset{
		ID:        uuid.New(),
		UserID:    "user_123",
		Name:      "My Server",
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(asset)
	if err != nil {
		t.Fatalf("Failed to marshal asset: %v", err)
	}

	jsonStr := string(jsonData)
	// Empty JSON object should NOT be omitted
	if !strings.Contains(jsonStr, `"metadata"`) {
		t.Error("JSON should contain 'metadata' field even when empty object")
	}
}

// Test 4: TestAsset_JSONUnmarshaling_CompleteJSON
func TestAsset_JSONUnmarshaling_CompleteJSON(t *testing.T) {
	jsonStr := `{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"user_id": "user_123",
		"name": "My Server",
		"type": "server",
		"hostname": "homelab-01",
		"metadata": {"cpu": "4 cores"},
		"created_at": "2024-03-15T10:00:00Z",
		"updated_at": "2024-03-15T11:00:00Z"
	}`

	var asset Asset
	err := json.Unmarshal([]byte(jsonStr), &asset)
	if err != nil {
		t.Fatalf("Failed to unmarshal asset: %v", err)
	}

	if asset.UserID != "user_123" {
		t.Errorf("Expected UserID 'user_123', got '%s'", asset.UserID)
	}
	if asset.Name != "My Server" {
		t.Errorf("Expected Name 'My Server', got '%s'", asset.Name)
	}
	if asset.Type == nil || *asset.Type != "server" {
		t.Error("Expected Type to be non-nil pointer to 'server'")
	}
	if asset.Hostname == nil || *asset.Hostname != "homelab-01" {
		t.Error("Expected Hostname to be non-nil pointer to 'homelab-01'")
	}
	if asset.Metadata == nil {
		t.Error("Expected Metadata to be non-nil")
	}
}

// Test 5: TestAsset_JSONUnmarshaling_MinimalJSON
func TestAsset_JSONUnmarshaling_MinimalJSON(t *testing.T) {
	jsonStr := `{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"user_id": "user_123",
		"name": "My Server",
		"created_at": "2024-03-15T10:00:00Z",
		"updated_at": "2024-03-15T11:00:00Z"
	}`

	var asset Asset
	err := json.Unmarshal([]byte(jsonStr), &asset)
	if err != nil {
		t.Fatalf("Failed to unmarshal asset: %v", err)
	}

	if asset.Type != nil {
		t.Error("Expected Type to be nil when not in JSON")
	}
	if asset.Hostname != nil {
		t.Error("Expected Hostname to be nil when not in JSON")
	}
	if asset.Metadata != nil {
		t.Error("Expected Metadata to be nil when not in JSON")
	}
}

// Test 6: TestAsset_JSONUnmarshaling_InvalidJSON
func TestAsset_JSONUnmarshaling_InvalidJSON(t *testing.T) {
	jsonStr := `{invalid json`

	var asset Asset
	err := json.Unmarshal([]byte(jsonStr), &asset)
	if err == nil {
		t.Error("Expected error when unmarshaling invalid JSON")
	}
}

// Test 7: TestAsset_JSONRoundTrip
func TestAsset_JSONRoundTrip(t *testing.T) {
	assetType := "vm"
	hostname := "test-vm"
	metadata := json.RawMessage(`{"os":"ubuntu"}`)

	original := Asset{
		ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		UserID:    "user_456",
		Name:      "Test VM",
		Type:      &assetType,
		Hostname:  &hostname,
		Metadata:  metadata,
		CreatedAt: time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 3, 15, 11, 0, 0, 0, time.UTC),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal back
	var result Asset
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Compare
	if result.ID != original.ID {
		t.Errorf("ID mismatch: %v != %v", result.ID, original.ID)
	}
	if result.UserID != original.UserID {
		t.Errorf("UserID mismatch: %v != %v", result.UserID, original.UserID)
	}
	if result.Name != original.Name {
		t.Errorf("Name mismatch: %v != %v", result.Name, original.Name)
	}
	if *result.Type != *original.Type {
		t.Errorf("Type mismatch: %v != %v", *result.Type, *original.Type)
	}
}

// ========== CreateAssetRequest Validation Tests ==========

// Test 8: TestCreateAssetRequest_Validation_Valid
func TestCreateAssetRequest_Validation_Valid(t *testing.T) {
	validate := validator.New()
	req := CreateAssetRequest{
		Name: "My Server",
		Type: stringPtr("server"),
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}
}

// Test 9: TestCreateAssetRequest_Validation_EmptyName
func TestCreateAssetRequest_Validation_EmptyName(t *testing.T) {
	validate := validator.New()
	req := CreateAssetRequest{
		Name: "",
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("Expected validation error for empty name")
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Name" && fieldError.Tag() == "required" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Name field with 'required' tag")
		}
	}
}

// Test 10: TestCreateAssetRequest_Validation_NameTooLong
func TestCreateAssetRequest_Validation_NameTooLong(t *testing.T) {
	validate := validator.New()
	// Create a 101-character string
	longName := strings.Repeat("a", 101)
	req := CreateAssetRequest{
		Name: longName,
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("Expected validation error for name > 100 chars")
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Name" && fieldError.Tag() == "max" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Name field with 'max' tag")
		}
	}
}

// Test 11: TestCreateAssetRequest_Validation_HostnameTooLong
func TestCreateAssetRequest_Validation_HostnameTooLong(t *testing.T) {
	validate := validator.New()
	// Create a 256-character string
	longHostname := strings.Repeat("a", 256)
	req := CreateAssetRequest{
		Name:     "Valid Name",
		Hostname: &longHostname,
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("Expected validation error for hostname > 255 chars")
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Hostname" && fieldError.Tag() == "max" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Hostname field with 'max' tag")
		}
	}
}

// Test 12: TestCreateAssetRequest_Validation_TypeTooLong
func TestCreateAssetRequest_Validation_TypeTooLong(t *testing.T) {
	validate := validator.New()
	// Create a 51-character string
	longType := strings.Repeat("a", 51)
	req := CreateAssetRequest{
		Name: "Valid Name",
		Type: &longType,
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("Expected validation error for type > 50 chars")
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Type" && fieldError.Tag() == "max" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Type field with 'max' tag")
		}
	}
}

// Test 13: TestCreateAssetRequest_Validation_NilOptionalFields
func TestCreateAssetRequest_Validation_NilOptionalFields(t *testing.T) {
	validate := validator.New()
	req := CreateAssetRequest{
		Name:     "Valid Name",
		Type:     nil,
		Hostname: nil,
		Metadata: nil,
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass with nil optional fields, got error: %v", err)
	}
}

// Test 14: TestCreateAssetRequest_JSONUnmarshaling_Complete
func TestCreateAssetRequest_JSONUnmarshaling_Complete(t *testing.T) {
	jsonStr := `{
		"name": "My Server",
		"type": "server",
		"hostname": "homelab-01",
		"metadata": {"cpu": "4 cores"}
	}`

	var req CreateAssetRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if req.Name != "My Server" {
		t.Errorf("Expected Name 'My Server', got '%s'", req.Name)
	}
	if req.Type == nil || *req.Type != "server" {
		t.Error("Expected Type to be non-nil pointer to 'server'")
	}
	if req.Hostname == nil || *req.Hostname != "homelab-01" {
		t.Error("Expected Hostname to be non-nil pointer to 'homelab-01'")
	}
	if req.Metadata == nil {
		t.Error("Expected Metadata to be non-nil")
	}
}

// Test 15: TestCreateAssetRequest_JSONUnmarshaling_MinimalRequired
func TestCreateAssetRequest_JSONUnmarshaling_MinimalRequired(t *testing.T) {
	jsonStr := `{"name": "Server"}`

	var req CreateAssetRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if req.Name != "Server" {
		t.Errorf("Expected Name 'Server', got '%s'", req.Name)
	}
	if req.Type != nil {
		t.Error("Expected Type to be nil")
	}
	if req.Hostname != nil {
		t.Error("Expected Hostname to be nil")
	}
	if req.Metadata != nil {
		t.Error("Expected Metadata to be nil")
	}
}

// ========== UpdateAssetRequest Validation Tests ==========

// Test 16: TestUpdateAssetRequest_Validation_AllNil
func TestUpdateAssetRequest_Validation_AllNil(t *testing.T) {
	validate := validator.New()
	req := UpdateAssetRequest{
		Name:     nil,
		Type:     nil,
		Hostname: nil,
		Metadata: nil,
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass with all nil fields (partial update), got error: %v", err)
	}
}

// Test 17: TestUpdateAssetRequest_Validation_SingleFieldUpdate
func TestUpdateAssetRequest_Validation_SingleFieldUpdate(t *testing.T) {
	validate := validator.New()
	req := UpdateAssetRequest{
		Name:     stringPtr("New Name"),
		Type:     nil,
		Hostname: nil,
		Metadata: nil,
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass for single field update, got error: %v", err)
	}
}

// Test 18: TestUpdateAssetRequest_Validation_NameTooLong
func TestUpdateAssetRequest_Validation_NameTooLong(t *testing.T) {
	validate := validator.New()
	longName := strings.Repeat("a", 101)
	req := UpdateAssetRequest{
		Name: &longName,
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("Expected validation error for name > 100 chars")
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Name" && fieldError.Tag() == "max" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Name field with 'max' tag")
		}
	}
}

// Test 19: TestUpdateAssetRequest_Validation_EmptyStringPointer
func TestUpdateAssetRequest_Validation_EmptyStringPointer(t *testing.T) {
	validate := validator.New()
	req := UpdateAssetRequest{
		Name: stringPtr(""),
	}

	// Empty string should be valid for updates (clearing field)
	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected empty string to be valid for update, got error: %v", err)
	}
}

// Test 20: TestUpdateAssetRequest_JSONUnmarshaling_PartialUpdate
func TestUpdateAssetRequest_JSONUnmarshaling_PartialUpdate(t *testing.T) {
	jsonStr := `{"name": "Updated"}`

	var req UpdateAssetRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if req.Name == nil || *req.Name != "Updated" {
		t.Error("Expected Name to be non-nil pointer to 'Updated'")
	}
	if req.Type != nil {
		t.Error("Expected Type to be nil")
	}
	if req.Hostname != nil {
		t.Error("Expected Hostname to be nil")
	}
	if req.Metadata != nil {
		t.Error("Expected Metadata to be nil")
	}
}

// Test 21: TestUpdateAssetRequest_JSONUnmarshaling_AllFields
func TestUpdateAssetRequest_JSONUnmarshaling_AllFields(t *testing.T) {
	jsonStr := `{
		"name": "Updated",
		"type": "container",
		"hostname": "new-host",
		"metadata": {"updated": true}
	}`

	var req UpdateAssetRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if req.Name == nil {
		t.Error("Expected Name to be non-nil")
	}
	if req.Type == nil {
		t.Error("Expected Type to be non-nil")
	}
	if req.Hostname == nil {
		t.Error("Expected Hostname to be non-nil")
	}
	if req.Metadata == nil {
		t.Error("Expected Metadata to be non-nil")
	}
}

// Test 22: TestUpdateAssetRequest_DistinguishNilVsEmpty
func TestUpdateAssetRequest_DistinguishNilVsEmpty(t *testing.T) {
	// Case 1: Nil means "don't update"
	req1 := UpdateAssetRequest{Name: nil}
	if req1.Name != nil {
		t.Error("Expected Name to be nil (don't update)")
	}

	// Case 2: Pointer to empty string means "set to empty"
	req2 := UpdateAssetRequest{Name: stringPtr("")}
	if req2.Name == nil {
		t.Error("Expected Name to be non-nil pointer")
	}
	if *req2.Name != "" {
		t.Errorf("Expected Name to point to empty string, got '%s'", *req2.Name)
	}

	// This test documents the semantic difference for service layer
	t.Log("req1.Name == nil means: don't update this field")
	t.Log("req2.Name == &\"\" means: update field to empty string")
}
