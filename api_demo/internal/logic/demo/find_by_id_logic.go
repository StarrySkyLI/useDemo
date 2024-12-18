package demo

import (
	"context"
	"rpc_demo/rpc"

	"api_demo/internal/svc"
	"api_demo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindByIdLogic) FindById(req *types.FindByIdReq) (resp *types.FindByIdResp, err error) {
	one, err := l.svcCtx.DemoRPC.FindOne(l.ctx, &rpc.GameInfoReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err

	}

	return &types.FindByIdResp{
		Id:          one.Info.Id,
		Name:        one.Info.Name,
		Code:        one.Info.Code,
		Data:        one.Info.Data,
		Create_time: one.Info.CreateTime,
		Update_time: one.Info.UpdateTime,
	}, nil
}
