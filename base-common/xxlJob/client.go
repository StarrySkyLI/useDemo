package xxlJob

import (
	"context"

	"github.com/xxl-job/xxl-job-executor-go"
	"github.com/zeromicro/go-zero/core/logc"
)

type Config struct {
	Address      string `json:",optional"`
	AccessToken  string `json:",optional"` // 请求令牌(默认为空)
	ExecutorPort string `json:",optional"` // 默认9999（非必填）
	RegistryKey  string `json:",optional"` // 执行器名称
}

type Client struct {
	conf    Config
	xxlExe  xxl.Executor
	handler []Handler
	ctx     context.Context
}

func NewClient(ctx context.Context, c Config) *Client {
	return &Client{
		conf: c,
		ctx:  ctx,
	}
}

func (c *Client) MustStart() {
	c.xxlExe = xxl.NewExecutor(
		xxl.ServerAddr(c.conf.Address),
		xxl.AccessToken(c.conf.AccessToken),   // 请求令牌(默认为空)
		xxl.ExecutorPort(c.conf.ExecutorPort), // 默认9999（非必填）
		xxl.RegistryKey(c.conf.RegistryKey),   // 执行器名称
		xxl.SetLogger(&Logger{Ctx: c.ctx}),    // 自定义日志
	)
	c.xxlExe.Init()

	for _, h := range c.handler {
		c.xxlExe.RegTask(h.Pattern(), xxl.TaskFunc(h.Handler()))
	}

	go func() {
		if err := c.xxlExe.Run(); err != nil {
			logc.Must(err)
		}
	}()
}

func (c *Client) Register(handler ...Handler) *Client {
	for _, h := range handler {
		c.handler = append(c.handler, h)
	}

	return c
}

func (c *Client) Stop() {
	c.xxlExe.Stop()
}
