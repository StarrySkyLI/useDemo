package demo

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"useDemo/application/api_demo/internal/svc"
	"useDemo/application/api_demo/internal/types"
)

type Api_demoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApi_demoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Api_demoLogic {
	return &Api_demoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Api_demoLogic) Api_demo(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
