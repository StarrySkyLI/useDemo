package rpcdemologic

import (
	"context"
	"rpc_demo/client/rpc_demo"

	"rpc_demo/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneLogic {
	return &FindOneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindOneLogic) FindOne(in *rpc_demo.GameInfoReq) (*rpc_demo.GameInfoRep, error) {
	// todo: add your logic here and delete this line

	return &rpc_demo.GameInfoRep{}, nil
}
