package model

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

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
