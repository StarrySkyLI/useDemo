package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer client.Close()
	ctx := context.Background()

	// 示例 1: 基础管道
	pipelineDemo(client, ctx)

	// 示例 2: 带结果返回的管道
	pipelineWithResultsDemo(client, ctx)

	// 示例 3: 事务管道
	transactionPipelineDemo(client, ctx)
}

func pipelineDemo(client *redis.Client, ctx context.Context) {
	start := time.Now()
	pipe := client.Pipeline()
	for i := 0; i < 1000; i++ {
		pipe.Incr(ctx, "pipeline_counter")
	}
	if _, err := pipe.Exec(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Pipeline 耗时:", time.Since(start))
}

func pipelineWithResultsDemo(client *redis.Client, ctx context.Context) {
	pipe := client.Pipeline()
	incr := pipe.Incr(ctx, "result_counter")
	get := pipe.Get(ctx, "result_counter")

	if _, err := pipe.Exec(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Incr 结果:", incr.Val())
	fmt.Println("Get 结果:", get.Val())
}

func transactionPipelineDemo(client *redis.Client, ctx context.Context) {
	pipe := client.TxPipeline()
	pipe.Set(ctx, "tx_key", "tx_value", 0)
	get := pipe.Get(ctx, "tx_key")

	if _, err := pipe.Exec(ctx); err != nil {
		panic(err)
	}

	fmt.Println("事务 Get 结果:", get.Val())
}
