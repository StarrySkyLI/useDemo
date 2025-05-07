package rpcdemologic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"useDemo/application/rpc_demo/client/rpc_demo"
	"useDemo/application/rpc_demo/internal/svc"
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
	// todo: add your logic here and delete this line

	return &rpc_demo.Response{}, nil
}
