package weathersnapshot

import (
	"github.com/google/uuid"
	"time"
)

// WeatherSnapshot represents a weather record.
type WeatherSnapshot struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Date        *time.Time `json:"date" db:"date"`
	City        string     `json:"city" db:"city"`
	Latitude    float64    `json:"latitude" db:"latitude"`
	Longitude   float64    `json:"longitude" db:"longitude"`
	TempMax     *float64   `json:"tempMax" db:"temp_max"`
	PrecipMM    *float64   `json:"precipMm" db:"precip_mm"`
	SunshineHrs *float64   `json:"sunshineHrs" db:"sunshine_hrs"`
}
