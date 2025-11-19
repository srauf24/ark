package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Asset represents a homelab asset (server, VM, container, etc.)
type Asset struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	UserID    string          `json:"user_id" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	Type      *string         `json:"type,omitempty" db:"type"`
	Hostname  *string         `json:"hostname,omitempty" db:"hostname"`
	Metadata  json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateAssetRequest is the DTO for creating a new asset
type CreateAssetRequest struct {
	Name     string           `json:"name" validate:"required,max=100"`
	Type     *string          `json:"type,omitempty" validate:"omitempty,max=50"`
	Hostname *string          `json:"hostname,omitempty" validate:"omitempty,max=255"`
	Metadata *json.RawMessage `json:"metadata,omitempty"`
}

// UpdateAssetRequest is the DTO for updating an existing asset
type UpdateAssetRequest struct {
	Name     *string          `json:"name,omitempty" validate:"omitempty,max=100"`
	Type     *string          `json:"type,omitempty" validate:"omitempty,max=50"`
	Hostname *string          `json:"hostname,omitempty" validate:"omitempty,max=255"`
	Metadata *json.RawMessage `json:"metadata,omitempty"`
}

// AssetResponse is the DTO for single asset responses
type AssetResponse struct {
	ID        uuid.UUID       `json:"id"`
	UserID    string          `json:"user_id"`
	Name      string          `json:"name"`
	Type      *string         `json:"type,omitempty"`
	Hostname  *string         `json:"hostname,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// NewAssetResponse converts an Asset domain model to AssetResponse DTO
func NewAssetResponse(asset *Asset) *AssetResponse {
	if asset == nil {
		return nil
	}

	return &AssetResponse{
		ID:        asset.ID,
		UserID:    asset.UserID,
		Name:      asset.Name,
		Type:      asset.Type,
		Hostname:  asset.Hostname,
		Metadata:  asset.Metadata,
		CreatedAt: asset.CreatedAt,
		UpdatedAt: asset.UpdatedAt,
	}
}

// AssetListResponse is the DTO for paginated asset list responses
type AssetListResponse struct {
	Assets []AssetResponse `json:"assets"`
	Total  int64           `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// NewAssetListResponse converts a slice of Assets to AssetListResponse with pagination metadata
func NewAssetListResponse(assets []*Asset, total int64, limit, offset int) *AssetListResponse {
	responses := make([]AssetResponse, 0, len(assets))
	for _, asset := range assets {
		if resp := NewAssetResponse(asset); resp != nil {
			responses = append(responses, *resp)
		}
	}

	return &AssetListResponse{
		Assets: responses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
}
