package repository

import "ark/internal/server"

type Repositories struct {
	// TODO: Add Asset and Log repositories when implemented
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}
