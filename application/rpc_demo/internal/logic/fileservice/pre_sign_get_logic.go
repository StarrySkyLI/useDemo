package fileservicelogic

import (
	"context"

	fileService "useDemo/application/rpc_demo/internal/service/file"
	"useDemo/application/rpc_demo/internal/svc"
	"useDemo/application/rpc_demo/rpc"
	"useDemo/application/rpc_demo/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSignGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	fileService *fileService.FileService
}

func NewPreSignGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSignGetLogic {
	return &PreSignGetLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		fileService: fileService.NewFileService(ctx, svcCtx),
	}
}

// pre sign a file url for user get it
func (l *PreSignGetLogic) PreSignGet(in *rpc.PreSignGetRequest) (*rpc.PreSignGetResponse, error) {
	url, err := l.fileService.PreSignGet(in.FileContext)
	if err != nil {
		return &rpc.PreSignGetResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}

	return &rpc.PreSignGetResponse{
		Meta: utils.GetSuccessMeta(),
		Url:  url,
	}, nil
}
