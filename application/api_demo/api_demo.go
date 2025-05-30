package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"useDemo/application/api_demo/internal/config"
	"useDemo/application/api_demo/internal/handler"
	"useDemo/application/api_demo/internal/svc"
	"useDemo/base-common/app"
	"useDemo/base-common/consul"
)

var configFile = flag.String("f", "etc/api_demo.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 初始化后内置调整
	app.InitAppServer()
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	// server.Use(middleware.NewApiHeaderMiddleware().Handle)
	// 服务注册
	err := consul.Register(c.Consul, fmt.Sprintf("%s:%d", c.ServiceConf.Prometheus.Host, c.ServiceConf.Prometheus.Port))
	if err != nil {
		fmt.Printf("register consul error: %v\n", err)
	}
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
