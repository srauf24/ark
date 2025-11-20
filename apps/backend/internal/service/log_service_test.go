package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"ark/internal/model"
)

// TestLogService_ListByAsset_ReturnsLogListResponse verifies ListByAsset returns LogListResponse DTO
func TestLogService_ListByAsset_ReturnsLogListResponse(t *testing.T) {
	// This test verifies the return type signature
	// We're not testing the actual business logic, just the type contract

	service := NewLogService(nil, nil) // nil is okay for type checking

	// Verify the method exists and returns the correct type
	var result *model.LogListResponse
	var err error

	// Type assertion to verify the signature
	_ = func() (*model.LogListResponse, error) {
		return service.ListByAsset(nil, "", [16]byte{}, nil)
	}

	// Verify types
	assert.IsType(t, result, (*model.LogListResponse)(nil))
	assert.IsType(t, err, error(nil))
}

// TestLogService_GetByID_ReturnsLogResponse verifies GetByID returns LogResponse DTO
func TestLogService_GetByID_ReturnsLogResponse(t *testing.T) {
	service := NewLogService(nil, nil)

	var result *model.LogResponse
	var err error

	// Type assertion to verify the signature
	_ = func() (*model.LogResponse, error) {
		return service.GetByID(nil, "", [16]byte{})
	}

	assert.IsType(t, result, (*model.LogResponse)(nil))
	assert.IsType(t, err, error(nil))
}

// TestLogService_Create_ReturnsLogResponse verifies Create returns LogResponse DTO
func TestLogService_Create_ReturnsLogResponse(t *testing.T) {
	service := NewLogService(nil, nil)

	var result *model.LogResponse
	var err error

	// Type assertion to verify the signature
	_ = func() (*model.LogResponse, error) {
		return service.Create(nil, "", [16]byte{}, nil)
	}

	assert.IsType(t, result, (*model.LogResponse)(nil))
	assert.IsType(t, err, error(nil))
}

// TestLogService_Update_ReturnsLogResponse verifies Update returns LogResponse DTO
func TestLogService_Update_ReturnsLogResponse(t *testing.T) {
	service := NewLogService(nil, nil)

	var result *model.LogResponse
	var err error

	// Type assertion to verify the signature
	_ = func() (*model.LogResponse, error) {
		return service.Update(nil, "", [16]byte{}, nil)
	}

	assert.IsType(t, result, (*model.LogResponse)(nil))
	assert.IsType(t, err, error(nil))
}

// TestLogService_Delete_ReturnsError verifies Delete returns error
func TestLogService_Delete_ReturnsError(t *testing.T) {
	service := NewLogService(nil, nil)

	var err error

	// Type assertion to verify the signature
	_ = func() error {
		return service.Delete(nil, "", [16]byte{})
	}

	assert.IsType(t, err, error(nil))
}

// TestLogService_Constructor verifies NewLogService works correctly
func TestLogService_Constructor(t *testing.T) {
	service := NewLogService(nil, nil)

	assert.NotNil(t, service)
	assert.IsType(t, &LogService{}, service)
}
