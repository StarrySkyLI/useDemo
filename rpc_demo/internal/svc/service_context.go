package svc

import (
	"rpc_demo/internal/config"
	"rpc_demo/internal/dao"
)

type ServiceContext struct {
	Config config.Config
	Dao    *dao.Dao
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Dao:    dao.NewDao(c),
	}
}
