package repository

import "github.com/srauf24/gardenjournal/internal/server"

type Repositories struct {
	Plant       *PlantRepository
	Observation *ObservationRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Plant:       NewPlantRepository(s),
		Observation: NewObservationRepository(s),
	}
}
