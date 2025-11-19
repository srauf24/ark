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
