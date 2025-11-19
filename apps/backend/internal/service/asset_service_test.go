package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"ark/internal/model"
	"ark/internal/repository"
)

// TestAssetService_List_ReturnsAssetListResponse verifies List returns AssetListResponse DTO
func TestAssetService_List_ReturnsAssetListResponse(t *testing.T) {
	// This test verifies the return type signature
	// We're not testing the actual business logic, just the type contract

	service := NewAssetService(nil) // nil is okay for type checking

	// Verify the method exists and returns the correct type
	var result *model.AssetListResponse
	var err error

	// Type assertion to verify the signature
	_ = func() (*model.AssetListResponse, error) {
		return service.List(nil, "", nil)
	}

	// Verify types
	assert.IsType(t, result, (*model.AssetListResponse)(nil))
	assert.IsType(t, err, error(nil))
}

// TestAssetService_GetByID_ReturnsAssetResponse verifies GetByID returns AssetResponse DTO
func TestAssetService_GetByID_ReturnsAssetResponse(t *testing.T) {
	service := NewAssetService(nil)

	var result *model.AssetResponse
	var err error

	// Type assertion to verify the signature
	_ = func() (*model.AssetResponse, error) {
		return service.GetByID(nil, "", [16]byte{})
	}

	assert.IsType(t, result, (*model.AssetResponse)(nil))
	assert.IsType(t, err, error(nil))
}

// TestAssetService_Create_ReturnsAssetResponse verifies Create returns AssetResponse DTO
func TestAssetService_Create_ReturnsAssetResponse(t *testing.T) {
	service := NewAssetService(nil)

	var result *model.AssetResponse
	var err error

	// Type assertion to verify the signature
	_ = func() (*model.AssetResponse, error) {
		return service.Create(nil, "", nil)
	}

	assert.IsType(t, result, (*model.AssetResponse)(nil))
	assert.IsType(t, err, error(nil))
}

// TestAssetService_Update_ReturnsAssetResponse verifies Update returns AssetResponse DTO
func TestAssetService_Update_ReturnsAssetResponse(t *testing.T) {
	service := NewAssetService(nil)

	var result *model.AssetResponse
	var err error

	// Type assertion to verify the signature
	_ = func() (*model.AssetResponse, error) {
		return service.Update(nil, "", [16]byte{}, nil)
	}

	assert.IsType(t, result, (*model.AssetResponse)(nil))
	assert.IsType(t, err, error(nil))
}

// TestAssetService_Delete_ReturnsError verifies Delete returns error
func TestAssetService_Delete_ReturnsError(t *testing.T) {
	service := NewAssetService(nil)

	var err error

	// Type assertion to verify the signature
	_ = func() error {
		return service.Delete(nil, "", [16]byte{})
	}

	assert.IsType(t, err, error(nil))
}

// TestAssetService_Constructor verifies NewAssetService works correctly
func TestAssetService_Constructor(t *testing.T) {
	repo := &repository.AssetRepository{}
	service := NewAssetService(repo)

	assert.NotNil(t, service)
	assert.IsType(t, &AssetService{}, service)
}
