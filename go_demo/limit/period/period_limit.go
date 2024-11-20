package main

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/limit"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// 一段时间内限制请求
func main() {
	store := redis.New("localhost:6379")
	limiter := limit.NewPeriodLimit(60, 10, store, "period_example-Key")

	result, err := limiter.Take("user1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch result {
	case limit.Allowed:
		fmt.Println("Request allowed")
	case limit.HitQuota:
		fmt.Println("Hit the quota")
	case limit.OverQuota:
		fmt.Println("Over the quota")
	default:
		fmt.Println("Unknown status")
	}
}
