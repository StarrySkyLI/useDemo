package svc

import (
	"rpc_demo/internal/config"
	"rpc_demo/internal/repo"
)

type ServiceContext struct {
	Config config.Config
	repo   *repo.Repo
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		repo:   repo.NewRepo(c),
		Config: c,
	}
}
