package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_WithValidClerkConfig(t *testing.T) {
	// Set up required environment variables
	envVars := map[string]string{
		"ARK_PRIMARY.ENV":                 "test",
		"ARK_SERVER.PORT":                 "8080",
		"ARK_SERVER.READ_TIMEOUT":         "30",
		"ARK_SERVER.WRITE_TIMEOUT":        "30",
		"ARK_SERVER.IDLE_TIMEOUT":         "60",
		"ARK_SERVER.CORS_ALLOWED_ORIGINS": "http://localhost:3000",
		"ARK_DATABASE.HOST":               "localhost",
		"ARK_DATABASE.PORT":               "5432",
		"ARK_DATABASE.USER":               "postgres",
		"ARK_DATABASE.PASSWORD":           "password",
		"ARK_DATABASE.NAME":               "testdb",
		"ARK_DATABASE.SSL_MODE":           "disable",
		"ARK_DATABASE.MAX_OPEN_CONNS":     "25",
		"ARK_DATABASE.MAX_IDLE_CONNS":     "25",
		"ARK_DATABASE.CONN_MAX_LIFETIME":  "300",
		"ARK_DATABASE.CONN_MAX_IDLE_TIME": "300",
		"ARK_AUTH.SECRET_KEY":             "secret",
		"ARK_AUTH.CLERK.SECRET_KEY":       "sk_test_1234567890",
		"ARK_AUTH.CLERK.JWT_ISSUER":       "https://test-app.clerk.accounts.dev",
		"ARK_INTEGRATION.RESEND_API_KEY":  "re_test_key",
		"ARK_REDIS.ADDRESS":               "localhost:6379",
	}

	// Set environment variables
	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		// Clean up environment variables
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Load config
	cfg, err := LoadConfig()

	// Assertions
	require.NoError(t, err, "LoadConfig should not return an error with valid configuration")
	require.NotNil(t, cfg, "Config should not be nil")

	// Verify Clerk configuration is loaded correctly
	assert.Equal(t, "sk_test_1234567890", cfg.Auth.Clerk.SecretKey, "Clerk secret key should be loaded correctly")
	assert.Equal(t, "https://test-app.clerk.accounts.dev", cfg.Auth.Clerk.JWTIssuer, "Clerk JWT issuer should be loaded correctly")
	assert.Empty(t, cfg.Auth.Clerk.PEMPublicKey, "PEM public key should be empty when not provided")
}

func TestLoadConfig_WithClerkPEMPublicKey(t *testing.T) {
	// Set up required environment variables including PEM key
	envVars := map[string]string{
		"ARK_PRIMARY.ENV":                 "test",
		"ARK_SERVER.PORT":                 "8080",
		"ARK_SERVER.READ_TIMEOUT":         "30",
		"ARK_SERVER.WRITE_TIMEOUT":        "30",
		"ARK_SERVER.IDLE_TIMEOUT":         "60",
		"ARK_SERVER.CORS_ALLOWED_ORIGINS": "http://localhost:3000",
		"ARK_DATABASE.HOST":               "localhost",
		"ARK_DATABASE.PORT":               "5432",
		"ARK_DATABASE.USER":               "postgres",
		"ARK_DATABASE.PASSWORD":           "password",
		"ARK_DATABASE.NAME":               "testdb",
		"ARK_DATABASE.SSL_MODE":           "disable",
		"ARK_DATABASE.MAX_OPEN_CONNS":     "25",
		"ARK_DATABASE.MAX_IDLE_CONNS":     "25",
		"ARK_DATABASE.CONN_MAX_LIFETIME":  "300",
		"ARK_DATABASE.CONN_MAX_IDLE_TIME": "300",
		"ARK_AUTH.SECRET_KEY":             "secret",
		"ARK_AUTH.CLERK.SECRET_KEY":       "sk_test_1234567890",
		"ARK_AUTH.CLERK.JWT_ISSUER":       "https://test-app.clerk.accounts.dev",
		"ARK_AUTH.CLERK.PEM_PUBLIC_KEY":   "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA\n-----END PUBLIC KEY-----",
		"ARK_INTEGRATION.RESEND_API_KEY":  "re_test_key",
		"ARK_REDIS.ADDRESS":               "localhost:6379",
	}

	// Set environment variables
	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		// Clean up environment variables
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Load config
	cfg, err := LoadConfig()

	// Assertions
	require.NoError(t, err, "LoadConfig should not return an error")
	require.NotNil(t, cfg, "Config should not be nil")

	// Verify PEM public key is loaded
	assert.NotEmpty(t, cfg.Auth.Clerk.PEMPublicKey, "PEM public key should be loaded when provided")
	assert.Contains(t, cfg.Auth.Clerk.PEMPublicKey, "BEGIN PUBLIC KEY", "PEM public key should contain valid PEM format")
}

// Note: Negative test cases (missing/invalid config) are not included here
// because LoadConfig() calls logger.Fatal() which executes os.Exit(1)
// instead of panicking. This cannot be reliably tested with standard Go testing.
// The validation is still in place and will fail at runtime if configuration is invalid.
