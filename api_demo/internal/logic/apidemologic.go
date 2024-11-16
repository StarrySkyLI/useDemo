package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"rpc_demo/rpc_demo"

	"api_demo/internal/svc"
	"api_demo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Api_demoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApi_demoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Api_demoLogic {
	return &Api_demoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Api_demoLogic) Api_demo(req *types.Request) (resp *types.Response, err error) {

	fmt.Println(l.ctx)
	ping, err := l.svcCtx.DemoRPC.Ping(l.ctx, &rpc_demo.Request{
		Ping: "ping",
	})
	if err != nil {
		return nil, err
	}
	logc.Info(l.ctx, "logc")

	fmt.Println("-----------------", ping)

	return &types.Response{
		Message: ping.Pong,
	}, nil
}
