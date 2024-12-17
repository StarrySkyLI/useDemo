package rpcdemologic

import (
	"context"
	"rpc_demo/client/rpc_demo"

	"rpc_demo/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GameListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGameListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameListLogic {
	return &GameListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GameListLogic) GameList(in *rpc_demo.GameListReq) (*rpc_demo.GameListRep, error) {
	// todo: add your logic here and delete this line

	return &rpc_demo.GameListRep{}, nil
}
