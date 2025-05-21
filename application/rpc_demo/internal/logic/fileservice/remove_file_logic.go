package fileservicelogic

import (
	"context"

	fileService "useDemo/application/rpc_demo/internal/service/file"
	"useDemo/application/rpc_demo/internal/svc"
	"useDemo/application/rpc_demo/rpc"
	"useDemo/application/rpc_demo/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	fileService *fileService.FileService
}

func NewRemoveFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveFileLogic {
	return &RemoveFileLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		fileService: fileService.NewFileService(ctx, svcCtx),
	}
}

// remove a file
func (l *RemoveFileLogic) RemoveFile(in *rpc.RemoveFileRequest) (*rpc.RemoveFileResponse, error) {
	if err := l.fileService.RemoveFile(in.FileContext); err != nil {
		return &rpc.RemoveFileResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}

	return &rpc.RemoveFileResponse{
		Meta: utils.GetSuccessMeta(),
	}, nil
}
