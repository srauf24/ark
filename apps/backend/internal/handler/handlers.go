package handler

import (
	"ark/internal/server"
	"ark/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
	Asset   *AssetHandler
	Log     *LogHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
		Asset:   NewAssetHandler(services.Asset),
		Log:     NewLogHandler(services.Log),
	}
}
