package dao

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"useDemo/application/rpc_demo/internal/config"
	"useDemo/base-common/middleware"
	"useDemo/base-common/orm"
)

type Dao struct {
	DB       *orm.DB
	BizRedis *redis.Redis
	config   config.Config
	ctx      context.Context
}

func NewDao(c config.Config) *Dao {
	db := orm.MustNewMysql(&orm.Config{
		DSN:          c.DB.DataSource,
		MaxOpenConns: c.DB.MaxOpenConns,
		MaxIdleConns: c.DB.MaxIdleConns,
		MaxLifetime:  c.DB.MaxLifetime,
	})

	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	}, redis.WithHook(middleware.NewBigKeyHook(100)))
	return &Dao{
		DB:       db,
		BizRedis: rds,
		config:   c,
		ctx:      context.Background(),
	}
}
