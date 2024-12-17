package main

import (
	"flag"
	"fmt"
	"gitlab.coolgame.world/go-template/base-common/app"
	"gitlab.coolgame.world/go-template/base-common/middleware"

	"api_demo/internal/config"
	"api_demo/internal/handler"
	"api_demo/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/api_demo-api.yaml", "the config file")

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
	server.Use(middleware.NewApiHeaderMiddleware().Handle)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
