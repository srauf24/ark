package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client provides weather data from Open Meteo API
type Client struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// NewClient creates a new weather client with default settings
func NewClient() *Client {
	return &Client{
		baseURL: "https://api.open-meteo.com/v1/forecast",
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		timeout: 5 * time.Second,
	}
}

// NewClientWithTimeout creates a new weather client with a custom timeout
func NewClientWithTimeout(timeout time.Duration) *Client {
	return &Client{
		baseURL: "https://api.open-meteo.com/v1/forecast",
		httpClient: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// FetchWeather fetches current weather data for given coordinates
// Returns nil (not error) if the fetch fails - this is intentional to not block log entry creation
func (c *Client) FetchWeather(ctx context.Context, lat, lon float64) (*WeatherData, error) {
	// Validate coordinates
	if lat < -90 || lat > 90 {
		// Return nil without error - don't block log entry creation
		return nil, nil
	}
	if lon < -180 || lon > 180 {
		// Return nil without error - don't block log entry creation
		return nil, nil
	}

	// Build request URL with required parameters
	// Open Meteo API docs: https://open-meteo.com/en/docs
	url := fmt.Sprintf("%s?latitude=%.6f&longitude=%.6f&current=temperature_2m,relative_humidity_2m,weather_code",
		c.baseURL, lat, lon)

	// Create request with context for cancellation
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		// Return nil without error - don't block log entry creation
		return nil, nil
	}

	// Set user agent
	req.Header.Set("User-Agent", "ARK/1.0")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Network error, timeout, etc. - return nil without error
		return nil, nil
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		// API error - return nil without error
		return nil, nil
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Failed to read response - return nil without error
		return nil, nil
	}

	// Parse response
	var apiResponse OpenMeteoResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		// Failed to parse JSON - return nil without error
		return nil, nil
	}

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, apiResponse.Current.Time)
	if err != nil {
		// Use current time if parsing fails
		timestamp = time.Now()
	}

	// Build weather data
	weatherData := &WeatherData{
		Temperature: apiResponse.Current.Temperature2m,
		Humidity:    float64(apiResponse.Current.RelativeHumidity2m),
		Condition:   GetConditionFromCode(apiResponse.Current.WeatherCode),
		Timestamp:   timestamp,
	}

	return weatherData, nil
}

// FetchWeatherSafe is a wrapper that ensures FetchWeather never panics
// and always returns gracefully, even in case of unexpected errors
func (c *Client) FetchWeatherSafe(ctx context.Context, lat, lon float64) *WeatherData {
	defer func() {
		if r := recover(); r != nil {
			// Recovered from panic - log it if needed but don't crash
			// For MVP, we just return nil
		}
	}()

	weather, _ := c.FetchWeather(ctx, lat, lon)
	return weather
}
