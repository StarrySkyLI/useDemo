package svc

import (
	"rpc_demo/internal/config"
	"rpc_demo/internal/dao"
)

type ServiceContext struct {
	Config config.Config
	dao    *dao.Dao
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		dao:    dao.NewDao(c),
	}
}
