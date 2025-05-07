package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"useDemo/application/api_demo/internal/config"
	"useDemo/application/rpc_demo/client/rpc_demo"
	"useDemo/base-common/rpc"
	"useDemo/base-common/rpc/interceptors"
)

type ServiceContext struct {
	Config  config.Config
	DemoRPC rpc_demo.RpcDemo
}

func NewServiceContext(c config.Config) *ServiceContext {
	demoRpc := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: rpc.GenRpcTarget(c.RPC.DemoRPC),
	}, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor(c.Name, c.Name)))
	return &ServiceContext{
		Config:  c,
		DemoRPC: rpc_demo.NewRpcDemo(demoRpc),
	}
}
