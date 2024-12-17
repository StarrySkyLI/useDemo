package demo

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
