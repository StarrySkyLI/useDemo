package middleware

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/mapping"
	"net"
	"strings"
)

type BigKeyHook struct {
	limitSize int
}

func NewBigKeyHook(limitSize int) *BigKeyHook {
	return &BigKeyHook{
		limitSize: limitSize,
	}
}

// DialHook 可以获取到 redis node 的地址信息，可以用于打印日志、监控等
func (h BigKeyHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		logc.Infow(ctx, "dialing", logc.Field("network", network), logc.Field("addr", addr))
		return next(ctx, network, addr)
	}
}

// ProcessHook 可以获取到 redis 命令信息
func (h BigKeyHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		// 执行前添加处理逻辑
		logc.Infow(ctx, "processing", logc.Field("cmd", cmd.String()))

		// 执行命令
		err := next(ctx, cmd)

		// 执行后添加处理逻辑
		h.cmdCheck(ctx, cmd)
		return err
	}
}

// ProcessPipelineHook 可以获取到 redis pipeline 命令信息
func (h BigKeyHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		// 执行前添加处理逻辑
		logc.Infow(ctx, "processing pipeline", logc.Field("cmds", cmds))

		// 执行命令
		err := next(ctx, cmds)

		// 执行后添加处理逻辑
		for _, cmd := range cmds {
			h.cmdCheck(ctx, cmd)
		}
		return err
	}
}

// cmdCheck 检查命令是否超过限制, 使用白名单的方式，只检查部分命令
func (h *BigKeyHook) cmdCheck(ctx context.Context, cmd redis.Cmder) {
	if h.limitSize <= 0 || len(cmd.Args()) < 2 ||
		(cmd.Err() != nil && !errors.Is(cmd.Err(), redis.Nil)) {
		return
	}
	var (
		size int
		key  = mapping.Repr(cmd.Args()[1])
	)
	switch strings.ToLower(cmd.Name()) {
	case "get":
		c, ok := cmd.(*redis.StringCmd)
		if !ok {
			return
		}
		size = len(c.Val())
	case "set", "setnx":
		if len(cmd.Args()) >= 3 {
			size = len(mapping.Repr(cmd.Args()[2]))
		}
	case "getset":
		c, ok := cmd.(*redis.StringCmd)
		if !ok {
			return
		}
		if c.Err() == nil {
			size = len(c.Val())
		} else if len(c.Args()) >= 3 {
			size = len(mapping.Repr(c.Args()[2]))
		}
	case "hgetall":
		c, ok := cmd.(*redis.MapStringStringCmd)
		if !ok {
			return
		}
		for _, v := range c.Val() {
			size += len(v)
		}
	case "hget":
		if cmd.Err() != nil {
			return
		}
		c, ok := cmd.(*redis.StringCmd)
		if !ok {
			return
		}
		if len(cmd.Args()) >= 3 {
			key += ":" + mapping.Repr(cmd.Args()[2])
		}
		size = len(c.Val())
	case "hmget":
		c, ok := cmd.(*redis.SliceCmd)
		if !ok {
			return
		}
		for _, v := range c.Val() {
			size += len(mapping.Repr(v))
		}
	case "hset", "hsetnx":
		if len(cmd.Args()) >= 4 {
			key += ":" + mapping.Repr(cmd.Args()[2])
			size = len(mapping.Repr(cmd.Args()[3]))
		}
	case "hmset":
		for i := 3; i < len(cmd.Args()); i += 2 {
			size += len(mapping.Repr(cmd.Args()[i]))
		}
	case "sadd":
		for i := 2; i < len(cmd.Args()); i++ {
			size += len(mapping.Repr(cmd.Args()[i]))
		}
	case "smembers":
		c, ok := cmd.(*redis.StringSliceCmd)
		if !ok {
			return
		}
		for _, v := range c.Val() {
			size += len(v)
		}
	case "zrange":
		switch cmd.(type) {
		case *redis.StringSliceCmd:
			for _, v := range cmd.(*redis.StringSliceCmd).Val() {
				size += len(v)
			}
		case *redis.ZSliceCmd:
			for _, v := range cmd.(*redis.ZSliceCmd).Val() {
				size += len(mapping.Repr(v.Member))
			}
		}
	case "zadd":
		for i := 3; i < len(cmd.Args()); i += 2 {
			size += len(mapping.Repr(cmd.Args()[i]))
		}
	case "zrangebyscore":
		c, ok := cmd.(*redis.ZSliceCmd)
		if !ok {
			return
		}
		for _, v := range c.Val() {
			size += len(mapping.Repr(v.Member))
		}
	default:
		return
	}
	if size > h.limitSize {
		logc.Infow(ctx, "[REDIS] BigKey limit",
			logc.Field("key", key), logc.Field("size", size))
	}
	return
}
