package weather

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: These are unit tests for the weather client.
// For MVP, we test the client logic with mocked HTTP responses.

func TestWeatherClient_NewClient(t *testing.T) {
	// Test that NewClient creates a client with default settings
	client := NewClient()

	require.NotNil(t, client)
	assert.Equal(t, "https://api.open-meteo.com/v1/forecast", client.baseURL)
	assert.Equal(t, 5*time.Second, client.timeout)
	assert.NotNil(t, client.httpClient)

	t.Log("NewClient creates client with correct defaults")
}

func TestWeatherClient_NewClientWithTimeout(t *testing.T) {
	// Test that NewClientWithTimeout creates a client with custom timeout
	customTimeout := 10 * time.Second
	client := NewClientWithTimeout(customTimeout)

	require.NotNil(t, client)
	assert.Equal(t, customTimeout, client.timeout)

	t.Log("NewClientWithTimeout creates client with custom timeout")
}

func TestWeatherClient_FetchWeather_Success(t *testing.T) {
	// Test successful weather fetch
	// Create a mock server
	mockResponse := OpenMeteoResponse{
		Latitude:  37.7749,
		Longitude: -122.4194,
	}
	mockResponse.Current.Time = "2024-01-15T10:00:00Z"
	mockResponse.Current.Temperature2m = 18.5
	mockResponse.Current.RelativeHumidity2m = 65
	mockResponse.Current.WeatherCode = 2 // Partly cloudy

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request parameters
		assert.Contains(t, r.URL.String(), "latitude=37.774900")
		assert.Contains(t, r.URL.String(), "longitude=-122.419400")
		assert.Contains(t, r.URL.String(), "current=temperature_2m,relative_humidity_2m,weather_code")

		// Return mock response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClient()
	client.baseURL = server.URL

	// Fetch weather
	ctx := context.Background()
	weather, err := client.FetchWeather(ctx, 37.7749, -122.4194)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, weather)
	assert.Equal(t, 18.5, weather.Temperature)
	assert.Equal(t, 65.0, weather.Humidity)
	assert.Equal(t, "Partly cloudy", weather.Condition)
	assert.False(t, weather.Timestamp.IsZero())

	t.Log("FetchWeather successfully retrieves and parses weather data")
}

func TestWeatherClient_FetchWeather_InvalidCoordinates(t *testing.T) {
	// Test that invalid coordinates return nil without error
	client := NewClient()
	ctx := context.Background()

	testCases := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{
			name: "Latitude too high",
			lat:  91.0,
			lon:  0.0,
		},
		{
			name: "Latitude too low",
			lat:  -91.0,
			lon:  0.0,
		},
		{
			name: "Longitude too high",
			lat:  0.0,
			lon:  181.0,
		},
		{
			name: "Longitude too low",
			lat:  0.0,
			lon:  -181.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			weather, err := client.FetchWeather(ctx, tc.lat, tc.lon)

			// Should return nil without error (graceful failure)
			assert.NoError(t, err)
			assert.Nil(t, weather)
		})
	}

	t.Log("FetchWeather handles invalid coordinates gracefully")
}

func TestWeatherClient_FetchWeather_NetworkError(t *testing.T) {
	// Test that network errors return nil without error
	// Use an invalid URL to trigger network error
	client := NewClient()
	client.baseURL = "http://invalid-url-that-does-not-exist.local"

	ctx := context.Background()
	weather, err := client.FetchWeather(ctx, 37.7749, -122.4194)

	// Should return nil without error (graceful failure)
	assert.NoError(t, err)
	assert.Nil(t, weather)

	t.Log("FetchWeather handles network errors gracefully")
}

func TestWeatherClient_FetchWeather_Timeout(t *testing.T) {
	// Test that timeout returns nil without error
	// Create a mock server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay longer than client timeout
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client with very short timeout
	client := NewClientWithTimeout(10 * time.Millisecond)
	client.baseURL = server.URL

	ctx := context.Background()
	weather, err := client.FetchWeather(ctx, 37.7749, -122.4194)

	// Should return nil without error (graceful failure)
	assert.NoError(t, err)
	assert.Nil(t, weather)

	t.Log("FetchWeather handles timeout gracefully")
}

func TestWeatherClient_FetchWeather_APIError(t *testing.T) {
	// Test that API errors return nil without error
	testCases := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "Bad Request",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Internal Server Error",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Service Unavailable",
			statusCode: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			}))
			defer server.Close()

			client := NewClient()
			client.baseURL = server.URL

			ctx := context.Background()
			weather, err := client.FetchWeather(ctx, 37.7749, -122.4194)

			// Should return nil without error (graceful failure)
			assert.NoError(t, err)
			assert.Nil(t, weather)
		})
	}

	t.Log("FetchWeather handles API errors gracefully")
}

func TestWeatherClient_FetchWeather_InvalidJSON(t *testing.T) {
	// Test that invalid JSON returns nil without error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	ctx := context.Background()
	weather, err := client.FetchWeather(ctx, 37.7749, -122.4194)

	// Should return nil without error (graceful failure)
	assert.NoError(t, err)
	assert.Nil(t, weather)

	t.Log("FetchWeather handles invalid JSON gracefully")
}

func TestWeatherClient_FetchWeather_ContextCancellation(t *testing.T) {
	// Test that context cancellation returns nil without error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay to ensure context is cancelled first
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	// Create a context that's immediately cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	weather, err := client.FetchWeather(ctx, 37.7749, -122.4194)

	// Should return nil without error (graceful failure)
	assert.NoError(t, err)
	assert.Nil(t, weather)

	t.Log("FetchWeather handles context cancellation gracefully")
}

func TestWeatherClient_FetchWeatherSafe(t *testing.T) {
	// Test that FetchWeatherSafe never panics
	client := NewClient()
	client.baseURL = "http://invalid-url.local"

	ctx := context.Background()

	// This should not panic even with invalid URL
	weather := client.FetchWeatherSafe(ctx, 37.7749, -122.4194)

	// Should return nil (graceful failure)
	assert.Nil(t, weather)

	t.Log("FetchWeatherSafe handles errors without panicking")
}

func TestWeatherClient_GetConditionFromCode(t *testing.T) {
	// Test weather code to condition mapping
	testCases := []struct {
		code      int
		condition string
	}{
		{0, "Clear sky"},
		{2, "Partly cloudy"},
		{61, "Slight rain"},
		{95, "Thunderstorm"},
		{999, "Unknown"}, // Invalid code
	}

	for _, tc := range testCases {
		t.Run(tc.condition, func(t *testing.T) {
			result := GetConditionFromCode(tc.code)
			assert.Equal(t, tc.condition, result)
		})
	}

	t.Log("GetConditionFromCode correctly maps weather codes")
}

func TestWeatherClient_RealCoordinates(t *testing.T) {
	// Test with various real-world coordinate examples
	testCases := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{
			name: "San Francisco, USA",
			lat:  37.7749,
			lon:  -122.4194,
		},
		{
			name: "London, UK",
			lat:  51.5074,
			lon:  -0.1278,
		},
		{
			name: "Tokyo, Japan",
			lat:  35.6762,
			lon:  139.6503,
		},
		{
			name: "Sydney, Australia",
			lat:  -33.8688,
			lon:  151.2093,
		},
		{
			name: "Equator",
			lat:  0.0,
			lon:  0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Validate coordinates are in valid range
			assert.GreaterOrEqual(t, tc.lat, -90.0)
			assert.LessOrEqual(t, tc.lat, 90.0)
			assert.GreaterOrEqual(t, tc.lon, -180.0)
			assert.LessOrEqual(t, tc.lon, 180.0)
		})
	}

	t.Log("Real-world coordinates are within valid ranges")
}

func TestWeatherClient_NonBlockingBehavior(t *testing.T) {
	// Test that weather client failures don't block the calling code
	client := NewClient()
	client.baseURL = "http://invalid-url.local"

	ctx := context.Background()

	// Start timer
	start := time.Now()

	// This should fail quickly (within timeout) and not hang
	weather, err := client.FetchWeather(ctx, 37.7749, -122.4194)

	duration := time.Since(start)

	// Should complete within reasonable time (timeout + small buffer)
	assert.Less(t, duration, 10*time.Second)

	// Should return gracefully
	assert.NoError(t, err)
	assert.Nil(t, weather)

	t.Log("FetchWeather completes quickly even on failure (non-blocking)")
}

// Integration test example (to be run manually against real API)
func TestWeatherClient_Integration_RealAPI(t *testing.T) {
	t.Skip("Integration test: requires real API call")

	// This test would make a real API call to Open Meteo
	// Only run manually when testing API integration
	client := NewClient()
	ctx := context.Background()

	weather, err := client.FetchWeather(ctx, 37.7749, -122.4194)

	require.NoError(t, err)
	if weather != nil {
		assert.NotZero(t, weather.Temperature)
		assert.NotZero(t, weather.Humidity)
		assert.NotEmpty(t, weather.Condition)
		assert.False(t, weather.Timestamp.IsZero())

		t.Logf("Weather data: %.1f°C, %.0f%% humidity, %s",
			weather.Temperature, weather.Humidity, weather.Condition)
	}
}
