package config

import (
	"base-common/consul"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Consul consul.Conf
	RPC    struct {
		DemoRPC string `json:"DemoRpc"`
	}
}
