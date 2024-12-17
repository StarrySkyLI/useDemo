package main

import (
	"flag"
	"fmt"
	"gitlab.coolgame.world/go-template/base-common/app"
	"gitlab.coolgame.world/go-template/base-common/middleware"
	"rpc_demo/rpc"

	"rpc_demo/internal/config"
	rpc_demoServer "rpc_demo/internal/server/rpc_demo"
	"rpc_demo/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/rpc_demo.yaml", "the config file")

func main() {
	flag.Parse()
	app.InitAppServer()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		rpc.RegisterRpcDemoServer(grpcServer, rpc_demoServer.NewRpcDemoServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(middleware.NewRpcAuthMiddleware().Handle())
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
