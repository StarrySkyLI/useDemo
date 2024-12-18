package rpcdemologic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"rpc_demo/internal/dao/dto"
	"rpc_demo/internal/dao/model/game"
	"strconv"

	"rpc_demo/internal/svc"
	"rpc_demo/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GameListExportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	gameModel game.IGameModel
}

func NewGameListExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameListExportLogic {
	return &GameListExportLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		gameModel: game.NewGameModel(ctx, svcCtx.Dao, logx.WithContext(ctx)),
	}
}

func (l *GameListExportLogic) GameListExport(in *rpc.GameListReq, stream rpc.RpcDemo_GameListExportServer) error {
	results, _, err := l.gameModel.GameList(dto.GameList{
		Page:     in.Page,
		PageSize: in.PageSize,
	})
	if err != nil {
		return err
	}

	for _, v := range results {
		info := &rpc.GameInfo{
			Id:         v.ID,
			Name:       v.Name,
			Code:       v.Code,
			Data:       v.Data,
			CreateTime: strconv.FormatInt(v.CreatedAt.Unix(), 10),
			UpdateTime: strconv.FormatInt(v.UpdatedAt.Unix(), 10),
		}
		if err := stream.Send(info); err != nil {
			logc.Error(l.ctx, err)
			return errors.New("stream GameListExport")
		}
	}
	return nil
}
