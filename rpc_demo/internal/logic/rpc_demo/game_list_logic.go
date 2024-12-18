package rpcdemologic

import (
	"context"
	"rpc_demo/client/rpc_demo"
	"rpc_demo/internal/dao/dto"
	"rpc_demo/internal/dao/model/game"
	"rpc_demo/rpc"
	"strconv"

	"rpc_demo/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GameListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	gameModel game.IGameModel
}

func NewGameListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameListLogic {
	return &GameListLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		gameModel: game.NewGameModel(ctx, svcCtx.Dao, logx.WithContext(ctx)),
	}
}

func (l *GameListLogic) GameList(in *rpc_demo.GameListReq) (*rpc_demo.GameListRep, error) {
	results, _, err := l.gameModel.GameList(dto.GameList{
		Page:     in.Page,
		PageSize: in.PageSize,
	})
	if err != nil {
		return nil, err
	}
	var list []*rpc.GameInfo
	for _, v := range results {
		list = append(list, &rpc.GameInfo{
			Id:         v.ID,
			Name:       v.Name,
			Code:       v.Code,
			Data:       v.Data,
			CreateTime: strconv.FormatInt(v.CreatedAt.Unix(), 10),
			UpdateTime: strconv.FormatInt(v.UpdatedAt.Unix(), 10),
		})
	}

	return &rpc_demo.GameListRep{
		List: list,
	}, nil
}
