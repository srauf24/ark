package plant

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/srauf24/gardenjournal/internal/model"
	"github.com/srauf24/gardenjournal/internal/model/observation"
)

// Represents the row in the database and what we return as JSON to the client
type Plant struct {
	model.Base                  // expects fields like ID, CreatedAt, UpdatedAt (and/or SortOrder if your Base includes it)
	UserID      string          `json:"userId" db:"user_id"`
	Name        string          `json:"name" db:"name"`
	Species     string          `json:"species" db:"species"`
	Location    *string         `json:"location" db:"location"`
	PlantedDate *time.Time      `json:"plantedDate" db:"planted_date"`
	Notes       *string         `json:"notes" db:"notes"`
	Metadata    json.RawMessage `json:"metadata" db:"metadata"`
	SortOrder   int             `json:"sortOrder" db:"sort_order"` // keep only if not already in model.Base
}

// PopulatedPlant is the “rich” DTO for list/detail responses. Can add weather observation etc in the future
type PopulatedPlant struct {
	Plant
	Observations []observation.Observation `json:"observations" db:"observations"`
}

// Metadata holds additional, optional properties about a plant.
type Metadata struct {
	Tags []string `json:"tags"` // e.g. ["indoor", "succulent"]

	// Care and Environment
	WateringFrequency *string    `json:"wateringFrequency"` // e.g. "every 3 days"
	LastWateredAt     *time.Time `json:"lastWateredAt"`
	SunlightLevel     *string    `json:"sunlightLevel"`  // e.g. "full sun", "partial shade"
	SoilType          *string    `json:"soilType"`       // e.g. "loamy", "sandy", "clay"
	PotSizeCM         *float64   `json:"potSizeCm"`      // e.g. 15.5
	FertilizerType    *string    `json:"fertilizerType"` // e.g. "NPK 10-10-10"
	LastFertilizedAt  *time.Time `json:"lastFertilizedAt"`

	// Environmental Tracking (weather integration)
	LastWeatherSnapshotID *uuid.UUID `json:"lastWeatherSnapshotId"`
	AverageTempC          *float64   `json:"averageTempC"`
	AverageSunshineHrs    *float64   `json:"averageSunshineHrs"`

	// Growth and Health
	HealthStatus *string  `json:"healthStatus"` // e.g. "thriving", "wilting"
	GrowthStage  *string  `json:"growthStage"`  // e.g. "seedling", "blooming"
	HeightCM     *float64 `json:"heightCm"`     // most recent height reading

	// Visual / UI Enrichment
	ColorTag  *string `json:"colorTag"`  // e.g. "#4CAF50" for display in UI
	ImageURL  *string `json:"imageUrl"`  // optional link to uploaded plant image
	EmojiIcon *string `json:"emojiIcon"` // e.g. "🌵", "🌿", "🌸"

	// Future Expansion (AI Insights, Summaries)
	AIInsightSummary *string `json:"aiInsightSummary"`
}

//Future: add a populated Plant DTO when the UI starts needing richer, nested data (e.g., observations with weather snapshots).
