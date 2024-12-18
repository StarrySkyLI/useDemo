package demo

import (
	"context"
	"rpc_demo/rpc"

	"api_demo/internal/svc"
	"api_demo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListReq) (resp *types.ListResp, err error) {
	list, err := l.svcCtx.DemoRPC.GameList(l.ctx, &rpc.GameListReq{
		Page:     req.Page,
		PageSize: req.Page_size,
	})
	if err != nil {
		return nil, err
	}
	gameList := make([]types.GameInfo, 0)
	for _, v := range list.List {
		gameList = append(gameList, types.GameInfo{
			Id:          v.Id,
			Name:        v.Name,
			Code:        v.Code,
			Data:        v.Data,
			Create_time: v.CreateTime,
			Update_time: v.UpdateTime,
		})
	}
	return &types.ListResp{
		List: gameList,
	}, nil
}
