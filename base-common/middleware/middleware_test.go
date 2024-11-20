package middleware

import (
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"testing"
)

func TestBigKey(t *testing.T) {
	rds := redis.MustNewRedis(redis.RedisConf{
		Host: "localhost:6379",
		Type: "node",
	}, redis.WithHook(NewBigKeyHook(100)))
	rds.Set("hhh", "adasda")
	rds.Set("test", "diadasdjasidjhjkdfhadsfhasudifhuaisdhfuiashdfiuyaghsuidfhauisdhfuiashdfuiaghsdfuigasidhgfbashuidgbfhajsdgbfasbdgfhjkasd")
	// do something with rds
}
func TestRateLimitMiddleware(t *testing.T) {
	// 创建一个 rest 服务
	s := rest.MustNewServer(rest.RestConf{}) // 需要调整配置

	// 创建一个 redis
	rds := redis.MustNewRedis(redis.RedisConf{
		Host: "localhost:6379",
		Type: "node",
	}) // 需要调整配置

	// 创建一个 PeriodLimit
	var (
		seconds   = 1            // 限流周期为 1 秒
		quota     = 100          // 每个周期每个用户最多 100 次请求
		keyPrefix = "limit:uid:" // 限流 key 前缀, 会存进 redis 中，需要保证唯一性
	)
	limiter := limit.NewPeriodLimit(seconds, quota, rds, keyPrefix)

	// 添加自定义中间件
	s.Use(RateLimitMiddleware(limiter))

	s.Start()
}
