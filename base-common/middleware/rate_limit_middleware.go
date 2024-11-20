package middleware

import (
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

//local limit = tonumber(ARGV[1])   -- 限流的最大次数
//local window = tonumber(ARGV[2])  -- 限流时间窗口的秒数
//local current = redis.call("INCRBY", KEYS[1], 1)  -- 计数+1
//if current == 1 then
//redis.call("expire", KEYS[1], window)  -- 如果计数值是 1，这表示初始创建
//-- 设置过期时间为窗口时长
//end
//if current < limit then
//return 1  -- 如果当前计数值小于限流值，则允许请求，返回 1
//elseif current == limit then
//return 2  -- 如果当前计数值等于限流值，返回 2，表示达到了限流阈值
//else
//return 0  -- 如果当前计数值大于限流值，返回 0，表示超过了限流阈值
//end

// RateLimitMiddleware 用户限流中间件
func RateLimitMiddleware(limiter *limit.PeriodLimit) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			uid := getUserIDFromRequest(request)
			// 如果没有用户 ID，不进行限流
			if uid == "" {
				next(writer, request)
			}

			// 限流
			v, err := limiter.Take(uid)
			if err != nil {
				logx.Errorf("take limit failed: %v", err)
				next(writer, request)
			}

			// 如果超过限流次数，返回 429
			if v != limit.Allowed {
				http.Error(writer, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next(writer, request)
		}
	}
}

// getUserIDFromRequest extracts user ID from the request
func getUserIDFromRequest(r *http.Request) string {
	// 实现你的用户ID提取逻辑
	return r.Header.Get("X-User-ID")
}
