package svc

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/core/collection"
	"useDemo/application/rpc_demo/internal/config"
	"useDemo/application/rpc_demo/internal/dao"
	"useDemo/base-common/minio_service"

	"time"
)

const localCacheExpire = time.Second * 60

type ServiceContext struct {
	Config     config.Config
	Dao        *dao.Dao
	LocalCache *collection.Cache
	Minio      *minio_service.MinioService
}

func NewServiceContext(c config.Config) *ServiceContext {
	localCache, err := collection.NewCache(localCacheExpire)
	if err != nil {
		panic(err)
	}
	core, err := minio.NewCore(c.MinioConf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.MinioConf.AccessKeyID, c.MinioConf.SecretAccessKey, ""),
		Region: "us-east-1",
	})
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:     c,
		Dao:        dao.NewDao(c),
		LocalCache: localCache,
		Minio:      minio_service.NewMinio(core),
	}
}
