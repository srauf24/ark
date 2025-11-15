package jwt

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/rs/zerolog"
)

var (
	// ErrInvalidToken is returned when the token is malformed or invalid
	ErrInvalidToken = errors.New("invalid token format")

	// ErrTokenVerificationFailed is returned when token verification fails
	ErrTokenVerificationFailed = errors.New("token verification failed")

	// ErrTokenExpired is returned when the token has expired
	ErrTokenExpired = errors.New("token has expired")

	// ErrMissingClaims is returned when required claims are missing
	ErrMissingClaims = errors.New("missing required claims")
)

// VerifyClerkToken verifies a Clerk JWT token and returns the session claims.
// It performs the following validations:
// - Token format validation
// - Signature verification using Clerk SDK
// - Expiration check
// - Required claims presence
func VerifyClerkToken(ctx context.Context, token string, logger *zerolog.Logger) (*clerk.SessionClaims, error) {
	start := time.Now()

	// Validate token format
	if token == "" {
		if logger != nil {
			logger.Debug().
				Dur("duration", time.Since(start)).
				Msg("token verification failed: empty token")
		}
		return nil, ErrInvalidToken
	}

	// Remove "Bearer " prefix if present
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		if logger != nil {
			logger.Debug().
				Dur("duration", time.Since(start)).
				Msg("token verification failed: empty token after trimming")
		}
		return nil, ErrInvalidToken
	}

	// Verify the token using Clerk SDK
	claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
		Token: token,
	})

	if err != nil {
		if logger != nil {
			logger.Error().
				Err(err).
				Dur("duration", time.Since(start)).
				Msg("token verification failed")
		}

		// Check if it's an expiration error
		if strings.Contains(err.Error(), "expired") || strings.Contains(err.Error(), "exp") {
			return nil, ErrTokenExpired
		}

		return nil, fmt.Errorf("%w: %v", ErrTokenVerificationFailed, err)
	}

	// Validate that we got valid claims
	if claims == nil {
		if logger != nil {
			logger.Error().
				Dur("duration", time.Since(start)).
				Msg("token verification succeeded but claims are nil")
		}
		return nil, ErrMissingClaims
	}

	// Validate required claims
	if claims.Subject == "" {
		if logger != nil {
			logger.Error().
				Dur("duration", time.Since(start)).
				Msg("token verification succeeded but subject claim is missing")
		}
		return nil, ErrMissingClaims
	}

	if logger != nil {
		logger.Info().
			Str("user_id", claims.Subject).
			Dur("duration", time.Since(start)).
			Msg("token verification successful")
	}

	return claims, nil
}

// ExtractBearerToken extracts the token from an Authorization header value.
// It expects the format: "Bearer <token>"
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrInvalidToken
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", ErrInvalidToken
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrInvalidToken
	}

	return token, nil
}
