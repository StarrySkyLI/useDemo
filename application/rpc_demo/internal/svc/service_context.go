package svc

import (
	"github.com/zeromicro/go-zero/core/collection"
	"useDemo/application/rpc_demo/internal/config"
	"useDemo/application/rpc_demo/internal/dao"

	"time"
)

const localCacheExpire = time.Duration(time.Second * 60)

type ServiceContext struct {
	Config     config.Config
	Dao        *dao.Dao
	LocalCache *collection.Cache
}

func NewServiceContext(c config.Config) *ServiceContext {
	localCache, err := collection.NewCache(localCacheExpire)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:     c,
		Dao:        dao.NewDao(c),
		LocalCache: localCache,
	}
}
