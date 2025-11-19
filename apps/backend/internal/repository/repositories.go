package repository

import "ark/internal/server"

type Repositories struct {
	Asset *AssetRepository
	// TODO: Add Log repository when implemented
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Asset: NewAssetRepository(s.DB.Pool),
	}
}
