package service

import (
	"ark/internal/server"

	"github.com/clerk/clerk-sdk-go/v2"
)

type AuthService struct {
	server *server.Server
}

func NewAuthService(s *server.Server) *AuthService {
	// Initialize Clerk SDK with the Clerk-specific secret key
	clerk.SetKey(s.Config.Auth.Clerk.SecretKey)
	return &AuthService{
		server: s,
	}
}
