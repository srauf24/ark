package jwt

import (
	"context"
	"os"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestExtractBearerToken_ValidToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "valid bearer token with Bearer prefix",
			header:    "Bearer abc123def456",
			wantToken: "abc123def456",
			wantErr:   false,
		},
		{
			name:      "valid bearer token with lowercase bearer",
			header:    "bearer xyz789token",
			wantToken: "xyz789token",
			wantErr:   false,
		},
		{
			name:      "valid bearer token with extra spaces",
			header:    "Bearer   tokenWithSpaces  ",
			wantToken: "tokenWithSpaces",
			wantErr:   false,
		},
		{
			name:      "empty header",
			header:    "",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "missing Bearer prefix",
			header:    "abc123def456",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "only Bearer without token",
			header:    "Bearer",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "Bearer with empty token",
			header:    "Bearer   ",
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ExtractBearerToken(tt.header)

			if tt.wantErr {
				assert.Error(t, err, "ExtractBearerToken should return an error")
				assert.Equal(t, ErrInvalidToken, err, "Error should be ErrInvalidToken")
			} else {
				assert.NoError(t, err, "ExtractBearerToken should not return an error")
				assert.Equal(t, tt.wantToken, token, "Extracted token should match expected value")
			}
		})
	}
}

func TestVerifyClerkToken_InvalidTokenFormat(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()

	tests := []struct {
		name      string
		token     string
		wantErr   error
		errString string
	}{
		{
			name:      "empty token",
			token:     "",
			wantErr:   ErrInvalidToken,
			errString: "invalid token format",
		},
		{
			name:      "only spaces",
			token:     "   ",
			wantErr:   ErrInvalidToken,
			errString: "invalid token format",
		},
		{
			name:      "Bearer prefix only",
			token:     "Bearer ",
			wantErr:   ErrInvalidToken,
			errString: "invalid token format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := VerifyClerkToken(ctx, tt.token, &logger)

			assert.Error(t, err, "VerifyClerkToken should return an error for invalid token format")
			assert.Nil(t, claims, "Claims should be nil for invalid token")
			assert.ErrorIs(t, err, tt.wantErr, "Error should match expected error type")
		})
	}
}

func TestVerifyClerkToken_MalformedToken(t *testing.T) {
	// Skip if no Clerk secret key is set (this test requires Clerk SDK to be initialized)
	if os.Getenv("ARK_AUTH.CLERK.SECRET_KEY") == "" {
		t.Skip("Skipping test: ARK_AUTH.CLERK.SECRET_KEY not set")
	}

	// Set Clerk secret key for testing
	clerk.SetKey(os.Getenv("ARK_AUTH.CLERK.SECRET_KEY"))

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()

	tests := []struct {
		name      string
		token     string
		expectErr bool
	}{
		{
			name:      "completely malformed token",
			token:     "this-is-not-a-jwt",
			expectErr: true,
		},
		{
			name:      "invalid JWT structure",
			token:     "header.payload",
			expectErr: true,
		},
		{
			name:      "base64 garbage",
			token:     "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo=",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := VerifyClerkToken(ctx, tt.token, &logger)

			if tt.expectErr {
				assert.Error(t, err, "VerifyClerkToken should return an error for malformed token")
				assert.Nil(t, claims, "Claims should be nil for malformed token")
			} else {
				assert.NoError(t, err, "VerifyClerkToken should not return an error")
				assert.NotNil(t, claims, "Claims should not be nil")
			}
		})
	}
}

func TestVerifyClerkToken_WithBearerPrefix(t *testing.T) {
	// This test verifies that the function properly handles tokens with "Bearer " prefix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()

	// Using an intentionally invalid token to test prefix stripping
	tokenWithoutBearer := "invalid.jwt.token"
	tokenWithBearer := "Bearer invalid.jwt.token"

	claims1, err1 := VerifyClerkToken(ctx, tokenWithoutBearer, &logger)
	claims2, err2 := VerifyClerkToken(ctx, tokenWithBearer, &logger)

	// Both should fail (invalid token), but the error should be the same
	assert.Error(t, err1, "Should return error for invalid token without Bearer")
	assert.Error(t, err2, "Should return error for invalid token with Bearer")
	assert.Nil(t, claims1, "Claims should be nil")
	assert.Nil(t, claims2, "Claims should be nil")

	// The errors should be of the same type (both verification failures)
	assert.ErrorIs(t, err1, ErrTokenVerificationFailed)
	assert.ErrorIs(t, err2, ErrTokenVerificationFailed)
}

// Note: Tests for valid token verification and expiration are not included
// because they require:
// 1. A valid Clerk secret key and configuration
// 2. The ability to generate valid JWTs with controlled expiration
// 3. Integration with actual Clerk services
//
// These scenarios are better tested through integration tests with a real
// Clerk environment or using the actual API endpoint tests.
//
// The current unit tests focus on:
// - Input validation (empty, malformed tokens)
// - Token extraction logic (Bearer prefix handling)
// - Error handling paths
