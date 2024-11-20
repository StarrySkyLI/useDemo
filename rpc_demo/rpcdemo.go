package main

import (
	"flag"
	"fmt"
	"gitlab.coolgame.world/go-template/base-common/app"
	"gitlab.coolgame.world/go-template/base-common/middleware"

	"rpc_demo/internal/config"
	"rpc_demo/internal/server"
	"rpc_demo/internal/svc"
	"rpc_demo/rpc_demo"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/rpcdemo.yaml", "the config file")

func main() {
	flag.Parse()
	// 初始化后内置调整
	app.InitAppServer()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		rpc_demo.RegisterRpcDemoServer(grpcServer, server.NewRpcDemoServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)

		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(middleware.NewRpcAuthMiddleware().Handle())
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
