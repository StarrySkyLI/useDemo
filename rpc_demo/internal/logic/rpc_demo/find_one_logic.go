package rpcdemologic

import (
	"context"
	"rpc_demo/client/rpc_demo"
	"rpc_demo/internal/dao/model/game"
	"rpc_demo/rpc"
	"strconv"

	"rpc_demo/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	gameModel game.IGameModel
}

func NewFindOneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneLogic {
	return &FindOneLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		gameModel: game.NewGameModel(ctx, svcCtx.Dao, logx.WithContext(ctx)),
	}
}

func (l *FindOneLogic) FindOne(in *rpc_demo.GameInfoReq) (*rpc_demo.GameInfoRep, error) {
	info, err := l.gameModel.FindOne(in.GetId())
	if err != nil {
		return nil, err
	}

	return &rpc_demo.GameInfoRep{
		Info: &rpc.GameInfo{
			Id:         info.ID,
			Name:       info.Name,
			Code:       info.Code,
			Data:       info.Data,
			CreateTime: strconv.FormatInt(info.CreatedAt.Unix(), 10),
			UpdateTime: strconv.FormatInt(info.UpdatedAt.Unix(), 10),
		},
	}, nil
}
