package weathersnapshot

import (
	"github.com/google/uuid"
	"time"
)

func NewWeatherSnapshot(date *time.Time, city string, latitude float64, longitude float64, tempMax, precipMM, sunshineHrs *float64) *WeatherSnapshot {
	return &WeatherSnapshot{
		ID:          uuid.New(),
		Date:        date,
		City:        city,
		Latitude:    latitude,
		Longitude:   longitude,
		TempMax:     tempMax,
		PrecipMM:    precipMM,
		SunshineHrs: sunshineHrs,
	}
}

func (w *WeatherSnapshot) GetID() uuid.UUID {
	return w.ID
}

func (w *WeatherSnapshot) GetDate() *time.Time {
	return w.Date
}

func (w *WeatherSnapshot) GetCity() string {
	return w.City
}

func (w *WeatherSnapshot) GetLatitude() float64 {
	return w.Latitude
}

func (w *WeatherSnapshot) GetLongitude() float64 {
	return w.Longitude
}

func (w *WeatherSnapshot) GetTempMax() *float64 {
	return w.TempMax
}

func (w *WeatherSnapshot) GetPrecipMM() *float64 {
	return w.PrecipMM
}

func (w *WeatherSnapshot) GetSunshineHrs() *float64 {
	return w.SunshineHrs
}

func (w *WeatherSnapshot) SetID(id uuid.UUID) {
	w.ID = id
}

func (w *WeatherSnapshot) SetDate(date *time.Time) {
	w.Date = date
}

func (w *WeatherSnapshot) SetCity(city string) {
	w.City = city
}

func (w *WeatherSnapshot) SetLatitude(latitude float64) {
	w.Latitude = latitude
}

func (w *WeatherSnapshot) SetLongitude(longitude float64) {
	w.Longitude = longitude
}

func (w *WeatherSnapshot) SetTempMax(tempMax *float64) {
	w.TempMax = tempMax
}

func (w *WeatherSnapshot) SetPrecipMM(precipMM *float64) {
	w.PrecipMM = precipMM
}

func (w *WeatherSnapshot) SetSunshineHrs(sunshineHrs *float64) {
	w.SunshineHrs = sunshineHrs
}
