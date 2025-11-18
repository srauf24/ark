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
