package main

import (
	"context"
	"fmt"
	"time"

	"gitlab.coolgame.world/go-template/base-common/configServer"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := configServer.NewClient(ctx, configServer.Config{
		Env:         "test",
		AppName:     "base-common",
		Endpoints:   []string{"localhost:2379"},
		Username:    "",
		Password:    "",
		DialTimeout: 0,
	})
	client.MustStart()

	fmt.Println(client.Get("test"))

	select {}
}
