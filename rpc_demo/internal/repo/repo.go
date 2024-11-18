package repo

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gitlab.coolgame.world/go-template/base-common/orm"
	"rpc_demo/internal/config"
)

type Repo struct {
	DB       *orm.DB
	BizRedis *redis.Redis
	config   config.Config
}

func NewRepo(c config.Config) *Repo {
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
	})
	return &Repo{
		DB:       db,
		BizRedis: rds,
		config:   c,
	}
}
