package config

import (
	"github.com/zeromicro/go-zero/rest"
	"gitlab.coolgame.world/go-template/base-common/consul"
)

type Config struct {
	rest.RestConf
	Consul consul.Conf
	RPC    struct {
		DemoRPC string `json:"DemoRpc"`
	}
}
