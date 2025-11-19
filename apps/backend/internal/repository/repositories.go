package repository

import "ark/internal/server"

type Repositories struct {
	Asset *AssetRepository
	Log   *LogRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Asset: NewAssetRepository(s.DB.Pool),
		Log:   NewLogRepository(s.DB.Pool),
	}
}
