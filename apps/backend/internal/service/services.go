package service

import (
	"ark/internal/lib/job"
	"ark/internal/repository"
	"ark/internal/server"
)

// Services holds all service layer instances
type Services struct {
	Auth  *AuthService
	Job   *job.JobService
	Asset *AssetService
	Log   *LogService
}

// NewServices creates and initializes all services with their dependencies
func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	// Initialize core services
	authService := NewAuthService(s)
	assetService := NewAssetService(repos.Asset)
	logService := NewLogService(repos.Log, repos.Asset)

	return &Services{
		Job:   s.Job,
		Auth:  authService,
		Asset: assetService,
		Log:   logService,
	}, nil
}
