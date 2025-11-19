package model

import (
	"time"

	"github.com/google/uuid"
)

// AssetLog represents a configuration change or troubleshooting log for an asset
type AssetLog struct {
	ID        uuid.UUID `json:"id" db:"id"`
	AssetID   uuid.UUID `json:"asset_id" db:"asset_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	Tags      []string  `json:"tags,omitempty" db:"tags"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateLogRequest is the DTO for creating a new log entry
type CreateLogRequest struct {
	Content string   `json:"content" validate:"required,min=2,max=10000"`
	Tags    []string `json:"tags,omitempty" validate:"omitempty,dive,max=50"`
}

// UpdateLogRequest is the DTO for updating an existing log entry
type UpdateLogRequest struct {
	Content *string   `json:"content,omitempty" validate:"omitempty,min=2,max=10000"`
	Tags    *[]string `json:"tags,omitempty" validate:"omitempty,dive,max=50"`
}
