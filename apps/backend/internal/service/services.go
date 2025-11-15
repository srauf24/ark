package service

import (
	"github.com/srauf24/gardenjournal/internal/lib/job"
	"github.com/srauf24/gardenjournal/internal/lib/weather"
	"github.com/srauf24/gardenjournal/internal/repository"
	"github.com/srauf24/gardenjournal/internal/server"
)

// Services holds all service layer instances
type Services struct {
	Auth        *AuthService
	Job         *job.JobService
	Plant       *PlantService
	Observation *ObservationService
}

// NewServices creates and initializes all services with their dependencies
func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	// Initialize core services
	authService := NewAuthService(s)

	// Create weather client for observation enrichment (future use)
	weatherClient := weather.NewClient()

	// Create domain services with repository dependencies
	plantService := NewPlantService(repos.Plant)
	observationService := NewObservationService(
		repos.Observation,
		repos.Plant,
		weatherClient,
	)

	return &Services{
		Job:         s.Job,
		Auth:        authService,
		Plant:       plantService,
		Observation: observationService,
	}, nil
}
