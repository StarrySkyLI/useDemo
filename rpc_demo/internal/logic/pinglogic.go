package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"

	"rpc_demo/internal/svc"
	"rpc_demo/rpc_demo"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *rpc_demo.Request) (*rpc_demo.Response, error) {
	fmt.Println(in.Ping)
	logc.Info(l.ctx, in.Ping)

	return &rpc_demo.Response{
		Pong: "pong",
	}, nil
}
