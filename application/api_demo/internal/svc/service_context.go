package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"useDemo/application/api_demo/internal/config"
	"useDemo/application/rpc_demo/client/rpc_demo"
)

type ServiceContext struct {
	Config  config.Config
	DemoRPC rpc_demo.RpcDemo
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:  c,
		DemoRPC: rpc_demo.NewRpcDemo(zrpc.MustNewClient(c.DemoRPC)),
	}
}
