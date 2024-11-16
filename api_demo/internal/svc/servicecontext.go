package svc

import (
	"api_demo/internal/config"
	"github.com/zeromicro/go-zero/zrpc"
	"gitlab.coolgame.world/go-template/base-common/rpc"
	"gitlab.coolgame.world/go-template/base-common/rpc/interceptors"
	"rpc_demo/rpc_demo"
	"rpc_demo/rpcdemoclient"
)

type ServiceContext struct {
	Config  config.Config
	DemoRPC rpc_demo.RpcDemoClient
}

func NewServiceContext(c config.Config) *ServiceContext {

	demoRpc := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: rpc.GenRpcTarget(c.RPC.DemoRPC),
	}, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor(c.Name, c.Name)))
	return &ServiceContext{
		Config:  c,
		DemoRPC: rpcdemoclient.NewRpcDemo(demoRpc),
	}
}
