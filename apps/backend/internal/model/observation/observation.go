package observation

import (
	"time"

	"github.com/google/uuid"
	"ark/internal/model"
)

// Observation represents a growth record for a plant.
type Observation struct {
	model.Base
	UserID    string    `json:"userId" db:"user_id"`
	PlantID   uuid.UUID `json:"plantId" db:"plant_id"`
	Date      time.Time `json:"date" db:"date"`
	HeightCM  *float64  `json:"heightCm" db:"height_cm"`
	Notes     *string   `json:"notes" db:"notes"`
	SortOrder int       `json:"sortOrder" db:"sort_order"`
}
