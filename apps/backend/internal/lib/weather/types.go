package weather

import "time"

// WeatherData represents current weather conditions
type WeatherData struct {
	Temperature float64   `json:"temperature"` // Temperature in Celsius
	Humidity    float64   `json:"humidity"`    // Relative humidity in %
	Condition   string    `json:"condition"`   // Weather condition description
	Timestamp   time.Time `json:"timestamp"`   // When the weather data was fetched
}

// OpenMeteoResponse represents the response from Open Meteo API
// API Documentation: https://open-meteo.com/en/docs
type OpenMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Time               string  `json:"time"`                 // ISO 8601 format
		Temperature2m      float64 `json:"temperature_2m"`       // Temperature at 2 meters in Celsius
		RelativeHumidity2m int     `json:"relative_humidity_2m"` // Relative humidity at 2 meters in %
		WeatherCode        int     `json:"weather_code"`         // WMO weather code
	} `json:"current"`
}

// WeatherCodeToCondition maps WMO weather codes to human-readable conditions
// Based on WMO code table: https://open-meteo.com/en/docs
var WeatherCodeToCondition = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Foggy",
	48: "Depositing rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	71: "Slight snow",
	73: "Moderate snow",
	75: "Heavy snow",
	77: "Snow grains",
	80: "Slight rain showers",
	81: "Moderate rain showers",
	82: "Violent rain showers",
	85: "Slight snow showers",
	86: "Heavy snow showers",
	95: "Thunderstorm",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

// GetConditionFromCode returns a human-readable weather condition from a WMO code
func GetConditionFromCode(code int) string {
	if condition, ok := WeatherCodeToCondition[code]; ok {
		return condition
	}
	return "Unknown"
}
