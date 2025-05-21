package main

import (
	"context"
	"flag"
	"fmt"

	"useDemo/application/rpc_demo/internal/config"
	"useDemo/application/rpc_demo/internal/dao/schema"
	"useDemo/application/rpc_demo/internal/jobs/test"
	file_server "useDemo/application/rpc_demo/internal/server/fileservice"
	"useDemo/application/rpc_demo/internal/svc"
	"useDemo/application/rpc_demo/rpc"
	"useDemo/base-common/app"
	"useDemo/base-common/consul"
	"useDemo/base-common/middleware"
	"useDemo/base-common/snowflake"
	"useDemo/base-common/xxlJob"

	rpc_demoServer "useDemo/application/rpc_demo/internal/server/rpc_demo"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/demo_service.yaml", "the config file")

func main() {
	flag.Parse()
	app.InitAppServer()
	// migrate.Handle()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	snowflake.InitDefaultSnowflakeNode(1)
	ctx := svc.NewServiceContext(c)
	err := ctx.Dao.DB.AutoMigrate(&schema.Game{})
	if err != nil {
		panic(fmt.Sprintf("AutoMigrate failed: %v", err))
	}
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		rpc.RegisterRpcDemoServer(grpcServer, rpc_demoServer.NewRpcDemoServer(ctx))
		rpc.RegisterFileServiceServer(grpcServer, file_server.NewFileServiceServer(ctx))

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
	err = consul.Register(c.Consul, fmt.Sprintf("%s:%d", c.ServiceConf.Prometheus.Host, c.ServiceConf.Prometheus.Port))
	if err != nil {
		fmt.Printf("register consul error: %v\n", err)
	}
	fmt.Printf("Starting rpc file_server at %s...\n", c.ListenOn)
	s.Start()
}
