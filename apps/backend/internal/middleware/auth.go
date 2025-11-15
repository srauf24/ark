package middleware

import (
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/labstack/echo/v4"
	"github.com/srauf24/gardenjournal/internal/errs"
	"github.com/srauf24/gardenjournal/internal/lib/jwt"
	"github.com/srauf24/gardenjournal/internal/server"
)

const (
	// ClerkSessionClaimsKey is the context key for storing Clerk session claims
	ClerkSessionClaimsKey = "clerk_session_claims"
)

type AuthMiddleware struct {
	server *server.Server
}

func NewAuthMiddleware(s *server.Server) *AuthMiddleware {
	return &AuthMiddleware{
		server: s,
	}
}

// ClerkAuthMiddleware is Phase 1 of authentication.
// It extracts the JWT from the Authorization header, verifies it using Clerk SDK,
// and stores the session claims in the Echo context for later use.
func (auth *AuthMiddleware) ClerkAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		requestID := GetRequestID(c)

		// Extract Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			auth.server.Logger.Debug().
				Str("function", "ClerkAuthMiddleware").
				Str("request_id", requestID).
				Dur("duration", time.Since(start)).
				Msg("missing Authorization header")
			return errs.NewUnauthorizedError("Missing authorization header", false)
		}

		// Extract Bearer token from header
		token, err := jwt.ExtractBearerToken(authHeader)
		if err != nil {
			auth.server.Logger.Debug().
				Err(err).
				Str("function", "ClerkAuthMiddleware").
				Str("request_id", requestID).
				Dur("duration", time.Since(start)).
				Msg("invalid Authorization header format")
			return errs.NewUnauthorizedError("Invalid authorization header format", false)
		}

		// Verify the JWT token using Clerk SDK
		claims, err := jwt.VerifyClerkToken(c.Request().Context(), token, auth.server.Logger)
		if err != nil {
			auth.server.Logger.Error().
				Err(err).
				Str("function", "ClerkAuthMiddleware").
				Str("request_id", requestID).
				Dur("duration", time.Since(start)).
				Msg("JWT verification failed")
			return errs.NewUnauthorizedError("Invalid or expired token", false)
		}

		// Store claims in Echo context for RequireAuth to use
		c.Set(ClerkSessionClaimsKey, claims)

		auth.server.Logger.Debug().
			Str("function", "ClerkAuthMiddleware").
			Str("user_id", claims.Subject).
			Str("request_id", requestID).
			Dur("duration", time.Since(start)).
			Msg("JWT verified and claims stored in context")

		return next(c)
	}
}

// RequireAuth is Phase 2 of authentication.
// It retrieves the verified session claims from the Echo context (set by ClerkAuthMiddleware)
// and populates user-specific data (user_id, user_role, permissions) into the context.
func (auth *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		requestID := GetRequestID(c)

		// Retrieve claims from Echo context
		claimsInterface := c.Get(ClerkSessionClaimsKey)
		if claimsInterface == nil {
			auth.server.Logger.Error().
				Str("function", "RequireAuth").
				Str("request_id", requestID).
				Dur("duration", time.Since(start)).
				Msg("could not get session claims from context")
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		// Type assert to *clerk.SessionClaims
		claims, ok := claimsInterface.(*clerk.SessionClaims)
		if !ok {
			auth.server.Logger.Error().
				Str("function", "RequireAuth").
				Str("request_id", requestID).
				Dur("duration", time.Since(start)).
				Msg("session claims type assertion failed")
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		// Set user context data for downstream handlers
		c.Set("user_id", claims.Subject)
		c.Set("user_role", claims.ActiveOrganizationRole)
		c.Set("permissions", claims.Claims.ActiveOrganizationPermissions)

		auth.server.Logger.Info().
			Str("function", "RequireAuth").
			Str("user_id", claims.Subject).
			Str("request_id", requestID).
			Dur("duration", time.Since(start)).
			Msg("user authenticated successfully")

		return next(c)
	}
}
