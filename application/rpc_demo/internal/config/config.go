package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"useDemo/base-common/consul"
	"useDemo/base-common/xxlJob"
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
