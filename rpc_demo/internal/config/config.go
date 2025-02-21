package config

import (
	"base-common/consul"
	"base-common/xxlJob"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		DataSource   string
		MaxOpenConns int `json:",default=10"`
		MaxIdleConns int `json:",default=100"`
		MaxLifetime  int `json:",default=3600"`
	}
	XxlJob   xxlJob.Config
	Consul   consul.Conf
	BizRedis redis.RedisConf
}
