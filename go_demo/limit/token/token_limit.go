package main

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// 控制事件在一秒钟内发生频率的限流器
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	store, _ := redis.NewRedis(redis.RedisConf{
		Host: "localhost:6379",
		Type: "node",
	})
	limiter := limit.NewTokenLimiter(5, 10, store, "token_example-key")

	if limiter.AllowCtx(ctx) {
		fmt.Println("Request allowed with context")
	} else {
		fmt.Println("Request not allowed with context")
	}
}
