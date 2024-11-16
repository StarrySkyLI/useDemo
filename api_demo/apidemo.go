package main

import (
	"api_demo/internal/config"
	"api_demo/internal/handler"
	"api_demo/internal/svc"
	"flag"
	"fmt"
	"gitlab.coolgame.world/go-template/base-common/app"
	"gitlab.coolgame.world/go-template/base-common/middleware"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/apidemo-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	// 初始化后内置调整
	app.InitAppServer()
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	server.Use(middleware.NewApiHeaderMiddleware().Handle)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
