package service

import (
	"ark/internal/lib/job"
	"ark/internal/repository"
	"ark/internal/server"
)

// Services holds all service layer instances
type Services struct {
	Auth *AuthService
	Job  *job.JobService
	// TODO: Add Asset and Log services when implemented
}

// NewServices creates and initializes all services with their dependencies
func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	// Initialize core services
	authService := NewAuthService(s)

	return &Services{
		Job:  s.Job,
		Auth: authService,
	}, nil
}
