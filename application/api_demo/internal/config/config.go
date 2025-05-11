package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"useDemo/base-common/consul"
)

type Config struct {
	rest.RestConf
	Consul  consul.Conf
	DemoRPC zrpc.RpcClientConf
}
