package model

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ========== AssetLog JSON Marshaling Tests ==========

// Test 1: TestAssetLog_MarshalJSON_AllFields
func TestAssetLog_MarshalJSON_AllFields(t *testing.T) {
	log := AssetLog{
		ID:        uuid.New(),
		AssetID:   uuid.New(),
		UserID:    "user_123",
		Content:   "Fixed nginx by updating /etc/nginx/nginx.conf and restarting with systemctl restart nginx",
		Tags:      []string{"nginx", "fix", "web-server"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(log)
	if err != nil {
		t.Fatalf("Failed to marshal log: %v", err)
	}

	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, `"id"`) {
		t.Error("JSON should contain 'id' field")
	}
	if !strings.Contains(jsonStr, `"asset_id"`) {
		t.Error("JSON should contain 'asset_id' field")
	}
	if !strings.Contains(jsonStr, `"user_id"`) {
		t.Error("JSON should contain 'user_id' field")
	}
	if !strings.Contains(jsonStr, `"content"`) {
		t.Error("JSON should contain 'content' field")
	}
	if !strings.Contains(jsonStr, `"tags"`) {
		t.Error("JSON should contain 'tags' field")
	}
	if !strings.Contains(jsonStr, `"created_at"`) {
		t.Error("JSON should contain 'created_at' field")
	}
	if !strings.Contains(jsonStr, `"updated_at"`) {
		t.Error("JSON should contain 'updated_at' field")
	}

	// Verify tags is an array
	if !strings.Contains(jsonStr, `"tags":["nginx","fix","web-server"]`) &&
		!strings.Contains(jsonStr, `"tags": ["nginx","fix","web-server"]`) &&
		!strings.Contains(jsonStr, `"tags":["nginx", "fix", "web-server"]`) {
		// Different formatting might occur, just check it's an array with values
		if !strings.Contains(jsonStr, `"tags":[`) {
			t.Error("Tags should be serialized as JSON array")
		}
	}
}

// Test 2: TestAssetLog_MarshalJSON_EmptyTags
func TestAssetLog_MarshalJSON_EmptyTags(t *testing.T) {
	log := AssetLog{
		ID:        uuid.New(),
		AssetID:   uuid.New(),
		UserID:    "user_123",
		Content:   "Test log",
		Tags:      []string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(log)
	if err != nil {
		t.Fatalf("Failed to marshal log: %v", err)
	}

	jsonStr := string(jsonData)
	// Empty tags array might be omitted or shown as []
	// Both are acceptable with omitempty
	if strings.Contains(jsonStr, `"tags":null`) {
		t.Error("Empty tags should not serialize as null")
	}
}

// Test 3: TestAssetLog_MarshalJSON_NilTags
func TestAssetLog_MarshalJSON_NilTags(t *testing.T) {
	log := AssetLog{
		ID:        uuid.New(),
		AssetID:   uuid.New(),
		UserID:    "user_123",
		Content:   "Test log",
		Tags:      nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(log)
	if err != nil {
		t.Fatalf("Failed to marshal log: %v", err)
	}

	jsonStr := string(jsonData)
	// With omitempty and nil slice, tags field should be omitted
	// or shown as null - both are acceptable
	_ = jsonStr // Just verify it marshals without error
}

// Test 4: TestAssetLog_UnmarshalJSON_CompletePayload
func TestAssetLog_UnmarshalJSON_CompletePayload(t *testing.T) {
	jsonStr := `{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"asset_id": "660e8400-e29b-41d4-a716-446655440001",
		"user_id": "user_123",
		"content": "Fixed nginx by restarting service",
		"tags": ["nginx", "fix"],
		"created_at": "2024-03-15T10:00:00Z",
		"updated_at": "2024-03-15T11:00:00Z"
	}`

	var log AssetLog
	err := json.Unmarshal([]byte(jsonStr), &log)
	if err != nil {
		t.Fatalf("Failed to unmarshal log: %v", err)
	}

	if log.UserID != "user_123" {
		t.Errorf("Expected UserID 'user_123', got '%s'", log.UserID)
	}
	if log.Content != "Fixed nginx by restarting service" {
		t.Errorf("Expected Content 'Fixed nginx by restarting service', got '%s'", log.Content)
	}
	if len(log.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(log.Tags))
	}
	if len(log.Tags) >= 2 {
		if log.Tags[0] != "nginx" {
			t.Errorf("Expected first tag 'nginx', got '%s'", log.Tags[0])
		}
		if log.Tags[1] != "fix" {
			t.Errorf("Expected second tag 'fix', got '%s'", log.Tags[1])
		}
	}
}

// Test 5: TestAssetLog_UnmarshalJSON_NoTags
func TestAssetLog_UnmarshalJSON_NoTags(t *testing.T) {
	jsonStr := `{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"asset_id": "660e8400-e29b-41d4-a716-446655440001",
		"user_id": "user_123",
		"content": "Fixed nginx",
		"created_at": "2024-03-15T10:00:00Z",
		"updated_at": "2024-03-15T11:00:00Z"
	}`

	var log AssetLog
	err := json.Unmarshal([]byte(jsonStr), &log)
	if err != nil {
		t.Fatalf("Failed to unmarshal log: %v", err)
	}

	// When tags field is omitted, it should be nil or empty slice
	if log.Tags != nil && len(log.Tags) > 0 {
		t.Error("Expected Tags to be nil or empty when not in JSON")
	}
}

// Test 6: TestAssetLog_UnmarshalJSON_InvalidJSON
func TestAssetLog_UnmarshalJSON_InvalidJSON(t *testing.T) {
	jsonStr := `{invalid json`

	var log AssetLog
	err := json.Unmarshal([]byte(jsonStr), &log)
	if err == nil {
		t.Error("Expected error when unmarshaling invalid JSON")
	}
}

// Test 7: TestAssetLog_Tags_MultipleValues
func TestAssetLog_Tags_MultipleValues(t *testing.T) {
	tags := []string{"tag1", "tag2", "tag3", "tag4", "tag5"}
	log := AssetLog{
		ID:        uuid.New(),
		AssetID:   uuid.New(),
		UserID:    "user_123",
		Content:   "Test",
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(log)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal back
	var result AssetLog
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify all tags preserved in order
	if len(result.Tags) != len(tags) {
		t.Errorf("Expected %d tags, got %d", len(tags), len(result.Tags))
	}
	for i, tag := range tags {
		if i < len(result.Tags) && result.Tags[i] != tag {
			t.Errorf("Tag %d mismatch: expected '%s', got '%s'", i, tag, result.Tags[i])
		}
	}
}

// Test 8: TestAssetLog_Content_LongText
func TestAssetLog_Content_LongText(t *testing.T) {
	// Create a 9000 character content string
	longContent := strings.Repeat("This is a long configuration change log entry. ", 180)

	log := AssetLog{
		ID:        uuid.New(),
		AssetID:   uuid.New(),
		UserID:    "user_123",
		Content:   longContent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(log)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal back
	var result AssetLog
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify content preserved
	if result.Content != longContent {
		t.Error("Long content not preserved correctly")
	}
	if len(result.Content) != len(longContent) {
		t.Errorf("Content length mismatch: expected %d, got %d", len(longContent), len(result.Content))
	}
}

// Test 9: TestAssetLog_UUID_Fields
func TestAssetLog_UUID_Fields(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	assetID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440001")

	log := AssetLog{
		ID:        id,
		AssetID:   assetID,
		UserID:    "user_123",
		Content:   "Test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(log)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(jsonData)
	// Verify UUIDs in standard format
	if !strings.Contains(jsonStr, "550e8400-e29b-41d4-a716-446655440000") {
		t.Error("ID UUID not in standard format in JSON")
	}
	if !strings.Contains(jsonStr, "660e8400-e29b-41d4-a716-446655440001") {
		t.Error("AssetID UUID not in standard format in JSON")
	}
}

// ========== CreateLogRequest Validation Tests ==========

// Test 10: TestCreateLogRequest_Validation_Valid_AllFields
func TestCreateLogRequest_Validation_Valid_AllFields(t *testing.T) {
	validate := validator.New()
	req := CreateLogRequest{
		Content: "Fixed nginx by restarting service",
		Tags:    []string{"nginx", "fix"},
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}
}

// Test 11: TestCreateLogRequest_Validation_Valid_ContentOnly
func TestCreateLogRequest_Validation_Valid_ContentOnly(t *testing.T) {
	validate := validator.New()
	req := CreateLogRequest{
		Content: "Fixed nginx",
		Tags:    nil,
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass with content only, got error: %v", err)
	}
}

// Test 12: TestCreateLogRequest_Validation_Invalid_EmptyContent
func TestCreateLogRequest_Validation_Invalid_EmptyContent(t *testing.T) {
	validate := validator.New()
	req := CreateLogRequest{
		Content: "",
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("Expected validation error for empty content")
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Content" && fieldError.Tag() == "required" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Content field with 'required' tag")
		}
	}
}

// Test 13: TestCreateLogRequest_Validation_Invalid_ContentTooShort
func TestCreateLogRequest_Validation_Invalid_ContentTooShort(t *testing.T) {
	validate := validator.New()
	req := CreateLogRequest{
		Content: "a", // 1 character
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("Expected validation error for content < 2 chars")
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Content" && fieldError.Tag() == "min" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Content field with 'min' tag")
		}
	}
}

// Test 14: TestCreateLogRequest_Validation_Invalid_ContentTooLong
func TestCreateLogRequest_Validation_Invalid_ContentTooLong(t *testing.T) {
	validate := validator.New()
	// Create 10,001 character string
	longContent := strings.Repeat("a", 10001)
	req := CreateLogRequest{
		Content: longContent,
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("Expected validation error for content > 10,000 chars")
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Content" && fieldError.Tag() == "max" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Content field with 'max' tag")
		}
	}
}

// Test 15: TestCreateLogRequest_Validation_Valid_ContentAtMaxLength
func TestCreateLogRequest_Validation_Valid_ContentAtMaxLength(t *testing.T) {
	validate := validator.New()
	// Exactly 10,000 characters
	maxContent := strings.Repeat("a", 10000)
	req := CreateLogRequest{
		Content: maxContent,
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass for content at max length (10,000 chars), got error: %v", err)
	}
}

// Test 16: TestCreateLogRequest_Validation_Invalid_TagTooLong
func TestCreateLogRequest_Validation_Invalid_TagTooLong(t *testing.T) {
	validate := validator.New()
	// Create tag with 51 characters
	longTag := strings.Repeat("a", 51)
	req := CreateLogRequest{
		Content: "Valid content",
		Tags:    []string{longTag},
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("Expected validation error for tag > 50 chars")
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			// The dive validation will show up on the array index
			if fieldError.Tag() == "max" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error with 'max' tag for tag length")
		}
	}
}

// Test 17: TestCreateLogRequest_Validation_Valid_MultipleTags
func TestCreateLogRequest_Validation_Valid_MultipleTags(t *testing.T) {
	validate := validator.New()
	req := CreateLogRequest{
		Content: "Test content",
		Tags:    []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass for multiple valid tags, got error: %v", err)
	}
}

// Test 18: TestCreateLogRequest_Validation_Valid_EmptyTagsArray
func TestCreateLogRequest_Validation_Valid_EmptyTagsArray(t *testing.T) {
	validate := validator.New()
	req := CreateLogRequest{
		Content: "Test content",
		Tags:    []string{},
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass for empty tags array, got error: %v", err)
	}
}

// Test 19: TestCreateLogRequest_JSONUnmarshaling_Complete
func TestCreateLogRequest_JSONUnmarshaling_Complete(t *testing.T) {
	jsonStr := `{
		"content": "Fixed nginx by restarting service",
		"tags": ["nginx", "fix"]
	}`

	var req CreateLogRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if req.Content != "Fixed nginx by restarting service" {
		t.Errorf("Expected Content 'Fixed nginx by restarting service', got '%s'", req.Content)
	}
	if len(req.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(req.Tags))
	}
}

// Test 20: TestCreateLogRequest_JSONUnmarshaling_ContentOnly
func TestCreateLogRequest_JSONUnmarshaling_ContentOnly(t *testing.T) {
	jsonStr := `{"content": "Fixed nginx"}`

	var req CreateLogRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if req.Content != "Fixed nginx" {
		t.Errorf("Expected Content 'Fixed nginx', got '%s'", req.Content)
	}
	if req.Tags != nil && len(req.Tags) > 0 {
		t.Error("Expected Tags to be nil or empty")
	}
}

// ========== UpdateLogRequest Validation Tests ==========

// Test 21: TestUpdateLogRequest_Validation_Valid_AllNil
func TestUpdateLogRequest_Validation_Valid_AllNil(t *testing.T) {
	validate := validator.New()
	req := UpdateLogRequest{
		Content: nil,
		Tags:    nil,
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass with all nil fields (partial update), got error: %v", err)
	}
}

// Test 22: TestUpdateLogRequest_Validation_Valid_OnlyContentSet
func TestUpdateLogRequest_Validation_Valid_OnlyContentSet(t *testing.T) {
	validate := validator.New()
	content := "New content"
	req := UpdateLogRequest{
		Content: &content,
		Tags:    nil,
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass for only content set, got error: %v", err)
	}
}

// Test 23: TestUpdateLogRequest_Validation_Valid_OnlyTagsSet
func TestUpdateLogRequest_Validation_Valid_OnlyTagsSet(t *testing.T) {
	validate := validator.New()
	tags := []string{"new-tag"}
	req := UpdateLogRequest{
		Content: nil,
		Tags:    &tags,
	}

	err := validate.Struct(req)
	if err != nil {
		t.Errorf("Expected validation to pass for only tags set, got error: %v", err)
	}
}

// Test 24: TestUpdateLogRequest_Validation_Invalid_ContentTooShort
func TestUpdateLogRequest_Validation_Invalid_ContentTooShort(t *testing.T) {
	validate := validator.New()
	shortContent := "a"
	req := UpdateLogRequest{
		Content: &shortContent,
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("Expected validation error for content < 2 chars")
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Content" && fieldError.Tag() == "min" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Content field with 'min' tag")
		}
	}
}

// Test 25: TestUpdateLogRequest_Validation_Invalid_ContentTooLong
func TestUpdateLogRequest_Validation_Invalid_ContentTooLong(t *testing.T) {
	validate := validator.New()
	longContent := strings.Repeat("a", 10001)
	req := UpdateLogRequest{
		Content: &longContent,
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("Expected validation error for content > 10,000 chars")
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		found := false
		for _, fieldError := range validationErrors {
			if fieldError.Field() == "Content" && fieldError.Tag() == "max" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected validation error on Content field with 'max' tag")
		}
	}
}

// Test 26: TestUpdateLogRequest_ClearTags_EmptySlice
func TestUpdateLogRequest_ClearTags_EmptySlice(t *testing.T) {
	jsonStr := `{"tags": []}`

	var req UpdateLogRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Tags should point to empty slice (not nil)
	if req.Tags == nil {
		t.Error("Expected Tags to point to empty slice, got nil")
	}
	if req.Tags != nil && len(*req.Tags) != 0 {
		t.Error("Expected Tags to point to empty slice")
	}
}

// Test 27: TestUpdateLogRequest_KeepTags_OmitField
func TestUpdateLogRequest_KeepTags_OmitField(t *testing.T) {
	jsonStr := `{"content": "Updated content"}`

	var req UpdateLogRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Tags should be nil (keep existing in DB)
	if req.Tags != nil {
		t.Error("Expected Tags to be nil when field omitted")
	}
}

// Test 28: TestUpdateLogRequest_JSONUnmarshaling_PartialFields
func TestUpdateLogRequest_JSONUnmarshaling_PartialFields(t *testing.T) {
	jsonStr := `{"content": "New content"}`

	var req UpdateLogRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if req.Content == nil {
		t.Error("Expected Content to be non-nil")
	}
	if req.Content != nil && *req.Content != "New content" {
		t.Errorf("Expected Content 'New content', got '%s'", *req.Content)
	}
	if req.Tags != nil {
		t.Error("Expected Tags to be nil")
	}
}

// Test 29: TestUpdateLogRequest_DistinguishNilVsEmpty
func TestUpdateLogRequest_DistinguishNilVsEmpty(t *testing.T) {
	// Case 1: Nil means "don't update"
	req1 := UpdateLogRequest{Tags: nil}
	if req1.Tags != nil {
		t.Error("Expected Tags to be nil (don't update)")
	}

	// Case 2: Pointer to empty slice means "clear all tags"
	emptyTags := []string{}
	req2 := UpdateLogRequest{Tags: &emptyTags}
	if req2.Tags == nil {
		t.Error("Expected Tags to be non-nil pointer")
	}
	if req2.Tags != nil && len(*req2.Tags) != 0 {
		t.Error("Expected Tags to point to empty slice")
	}

	// This test documents the semantic difference for service layer
	t.Log("req1.Tags == nil means: don't update tags")
	t.Log("req2.Tags == &[]string{} means: clear all tags")
}

// ========== LogResponse Tests ==========

// Test 30: TestNewLogResponse_AllFields
func TestNewLogResponse_AllFields(t *testing.T) {
	log := &AssetLog{
		ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		AssetID:   uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
		UserID:    "user_123",
		Content:   "Fixed nginx by restarting service",
		Tags:      []string{"nginx", "fix"},
		CreatedAt: time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 3, 15, 11, 0, 0, 0, time.UTC),
	}

	resp := NewLogResponse(log)

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}
	if resp.ID != log.ID {
		t.Errorf("ID mismatch: %v != %v", resp.ID, log.ID)
	}
	if resp.AssetID != log.AssetID {
		t.Errorf("AssetID mismatch: %v != %v", resp.AssetID, log.AssetID)
	}
	if resp.UserID != log.UserID {
		t.Errorf("UserID mismatch: %v != %v", resp.UserID, log.UserID)
	}
	if resp.Content != log.Content {
		t.Errorf("Content mismatch: %v != %v", resp.Content, log.Content)
	}
	if len(resp.Tags) != len(log.Tags) {
		t.Errorf("Tags length mismatch: %d != %d", len(resp.Tags), len(log.Tags))
	}
	for i, tag := range log.Tags {
		if i < len(resp.Tags) && resp.Tags[i] != tag {
			t.Errorf("Tag %d mismatch: %v != %v", i, resp.Tags[i], tag)
		}
	}
}

// Test 31: TestNewLogResponse_EmptyTags
func TestNewLogResponse_EmptyTags(t *testing.T) {
	log := &AssetLog{
		ID:        uuid.New(),
		AssetID:   uuid.New(),
		UserID:    "user_123",
		Content:   "Test log",
		Tags:      []string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	resp := NewLogResponse(log)

	// Empty tags array should be copied
	if resp.Tags == nil {
		t.Error("Expected Tags to be non-nil empty slice")
	}
	if len(resp.Tags) != 0 {
		t.Error("Expected Tags to be empty")
	}

	// Verify omitempty in JSON
	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	jsonStr := string(jsonData)
	// Empty array might be omitted or shown as [] with omitempty
	_ = jsonStr // Just verify it marshals without error
}

// Test 32: TestNewLogResponse_NilTags
func TestNewLogResponse_NilTags(t *testing.T) {
	log := &AssetLog{
		ID:        uuid.New(),
		AssetID:   uuid.New(),
		UserID:    "user_123",
		Content:   "Test log",
		Tags:      nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	resp := NewLogResponse(log)

	// Nil tags should remain nil
	if resp.Tags != nil && len(resp.Tags) > 0 {
		t.Error("Expected Tags to be nil or empty")
	}
}

// Test 33: TestNewLogResponse_NilInput
func TestNewLogResponse_NilInput(t *testing.T) {
	resp := NewLogResponse(nil)

	if resp != nil {
		t.Error("Expected nil response when input is nil")
	}
}

// Test 34: TestNewLogResponse_MarshalJSON
func TestNewLogResponse_MarshalJSON(t *testing.T) {
	log := &AssetLog{
		ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		AssetID:   uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
		UserID:    "user_123",
		Content:   "Test log",
		Tags:      []string{"test"},
		CreatedAt: time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 3, 15, 11, 0, 0, 0, time.UTC),
	}

	resp := NewLogResponse(log)
	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, `"id"`) {
		t.Error("JSON should contain 'id' field")
	}
	if !strings.Contains(jsonStr, `"asset_id"`) {
		t.Error("JSON should contain 'asset_id' field")
	}
	if !strings.Contains(jsonStr, `"user_id"`) {
		t.Error("JSON should contain 'user_id' field")
	}
	if !strings.Contains(jsonStr, `"content"`) {
		t.Error("JSON should contain 'content' field")
	}
	if !strings.Contains(jsonStr, `"tags"`) {
		t.Error("JSON should contain 'tags' field")
	}
	if !strings.Contains(jsonStr, `"created_at"`) {
		t.Error("JSON should contain 'created_at' field")
	}
	if !strings.Contains(jsonStr, `"updated_at"`) {
		t.Error("JSON should contain 'updated_at' field")
	}
}

// ========== LogListResponse Tests ==========

// Test 35: TestNewLogListResponse_MultipleLogs
func TestNewLogListResponse_MultipleLogs(t *testing.T) {
	logs := []*AssetLog{
		{
			ID:        uuid.New(),
			AssetID:   uuid.New(),
			UserID:    "user_123",
			Content:   "Log 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			AssetID:   uuid.New(),
			UserID:    "user_123",
			Content:   "Log 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			AssetID:   uuid.New(),
			UserID:    "user_123",
			Content:   "Log 3",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			AssetID:   uuid.New(),
			UserID:    "user_123",
			Content:   "Log 4",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			AssetID:   uuid.New(),
			UserID:    "user_123",
			Content:   "Log 5",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	resp := NewLogListResponse(logs, 100, 50, 0)

	if len(resp.Logs) != 5 {
		t.Errorf("Expected 5 logs, got %d", len(resp.Logs))
	}
	if resp.Total != 100 {
		t.Errorf("Expected Total=100, got %d", resp.Total)
	}
	if resp.Limit != 50 {
		t.Errorf("Expected Limit=50, got %d", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("Expected Offset=0, got %d", resp.Offset)
	}

	// Verify all logs were converted
	for i, logResp := range resp.Logs {
		if logResp.Content != logs[i].Content {
			t.Errorf("Log %d: Content mismatch", i)
		}
	}
}

// Test 36: TestNewLogListResponse_EmptySlice
func TestNewLogListResponse_EmptySlice(t *testing.T) {
	resp := NewLogListResponse([]*AssetLog{}, 0, 50, 0)

	if resp.Logs == nil {
		t.Error("Expected Logs to be non-nil empty slice, not nil")
	}
	if len(resp.Logs) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(resp.Logs))
	}
	if resp.Total != 0 {
		t.Errorf("Expected Total=0, got %d", resp.Total)
	}

	// Verify JSON serialization
	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(jsonData)
	// Empty slice should serialize as [] not null
	if strings.Contains(jsonStr, `"logs":null`) {
		t.Error("Logs should serialize as empty array [], not null")
	}
	if !strings.Contains(jsonStr, `"logs":[]`) {
		t.Error("Logs should serialize as empty array []")
	}
}

// Test 37: TestNewLogListResponse_NilSlice
func TestNewLogListResponse_NilSlice(t *testing.T) {
	resp := NewLogListResponse(nil, 0, 50, 0)

	if resp.Logs == nil {
		t.Error("Expected Logs to be non-nil empty slice, not nil")
	}
	if len(resp.Logs) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(resp.Logs))
	}
}

// Test 38: TestNewLogListResponse_SingleLog
func TestNewLogListResponse_SingleLog(t *testing.T) {
	logs := []*AssetLog{
		{
			ID:        uuid.New(),
			AssetID:   uuid.New(),
			UserID:    "user_123",
			Content:   "Single log",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	resp := NewLogListResponse(logs, 1, 50, 0)

	if len(resp.Logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(resp.Logs))
	}
	if resp.Logs[0].Content != "Single log" {
		t.Error("Log not converted correctly")
	}
}

// Test 39: TestNewLogListResponse_PaginationMetadata
func TestNewLogListResponse_PaginationMetadata(t *testing.T) {
	logs := []*AssetLog{
		{ID: uuid.New(), AssetID: uuid.New(), UserID: "user_123", Content: "Log", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	resp := NewLogListResponse(logs, 200, 50, 100)

	if resp.Total != 200 {
		t.Errorf("Expected Total=200, got %d", resp.Total)
	}
	if resp.Limit != 50 {
		t.Errorf("Expected Limit=50, got %d", resp.Limit)
	}
	if resp.Offset != 100 {
		t.Errorf("Expected Offset=100, got %d", resp.Offset)
	}
}

// Test 40: TestNewLogListResponse_SkipsNilLogs
func TestNewLogListResponse_SkipsNilLogs(t *testing.T) {
	logs := []*AssetLog{
		{ID: uuid.New(), AssetID: uuid.New(), UserID: "user_123", Content: "Log 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		nil, // nil log should be skipped
		{ID: uuid.New(), AssetID: uuid.New(), UserID: "user_123", Content: "Log 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	resp := NewLogListResponse(logs, 100, 50, 0)

	// Should only have 2 logs (nil skipped)
	if len(resp.Logs) != 2 {
		t.Errorf("Expected 2 logs (nil skipped), got %d", len(resp.Logs))
	}
}

// Test 41: TestNewLogListResponse_PreservesTags
func TestNewLogListResponse_PreservesTags(t *testing.T) {
	logs := []*AssetLog{
		{
			ID:        uuid.New(),
			AssetID:   uuid.New(),
			UserID:    "user_123",
			Content:   "Log with tags",
			Tags:      []string{"tag1", "tag2"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			AssetID:   uuid.New(),
			UserID:    "user_123",
			Content:   "Log without tags",
			Tags:      nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	resp := NewLogListResponse(logs, 2, 50, 0)

	if len(resp.Logs) != 2 {
		t.Fatalf("Expected 2 logs, got %d", len(resp.Logs))
	}

	// First log should have tags
	if len(resp.Logs[0].Tags) != 2 {
		t.Errorf("Expected first log to have 2 tags, got %d", len(resp.Logs[0].Tags))
	}

	// Second log should have nil/empty tags
	if resp.Logs[1].Tags != nil && len(resp.Logs[1].Tags) > 0 {
		t.Error("Expected second log to have nil or empty tags")
	}
}
