package handler

import (
	"ark/internal/server"
	"ark/internal/service"
)

type Handlers struct {
	Health      *HealthHandler
	OpenAPI     *OpenAPIHandler
	Plant       *PlantHandler
	Observation *ObservationHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:      NewHealthHandler(s),
		OpenAPI:     NewOpenAPIHandler(s),
		Plant:       NewPlantHandler(s, services.Plant),
		Observation: NewObservationHandler(s, services.Observation),
	}
}
