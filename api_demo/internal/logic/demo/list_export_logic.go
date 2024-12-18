package demo

import (
	"context"

	"api_demo/internal/svc"
	"api_demo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type List_exportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewList_exportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *List_exportLogic {
	return &List_exportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *List_exportLogic) List_export(req *types.ListReq) error {
	// todo: add your logic here and delete this line

	return nil
}
