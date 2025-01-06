package main

import (
	"base-common/app"
	"base-common/consul"
	"base-common/middleware"
	"base-common/xxlJob"
	"context"
	"flag"
	"fmt"
	"rpc_demo/internal/jobs/test"
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
	// cron job
	jobCli := xxlJob.NewClient(context.Background(), c.XxlJob).Register(
		&test.TestHandler{},
	)
	jobCli.MustStart()
	// 服务注册
	err := consul.Register(c.Consul, fmt.Sprintf("%s:%d", c.ServiceConf.Prometheus.Host, c.ServiceConf.Prometheus.Port))
	if err != nil {
		fmt.Printf("register consul error: %v\n", err)
	}
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
