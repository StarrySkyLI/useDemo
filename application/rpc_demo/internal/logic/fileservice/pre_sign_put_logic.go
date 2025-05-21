package fileservicelogic

import (
	"context"

	fileService "useDemo/application/rpc_demo/internal/service/file"
	"useDemo/application/rpc_demo/internal/svc"
	"useDemo/application/rpc_demo/rpc"
	"useDemo/application/rpc_demo/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSignPutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	fileService *fileService.FileService
}

func NewPreSignPutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSignPutLogic {
	return &PreSignPutLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		fileService: fileService.NewFileService(ctx, svcCtx),
	}
}

// pre sign a file url for user put it
func (l *PreSignPutLogic) PreSignPut(in *rpc.PreSignPutRequest) (*rpc.PreSignPutResponse, error) {
	fileId, existed, err := l.fileService.CheckFileExistedAndGetFile(in.FileContext)
	if err != nil {
		return &rpc.PreSignPutResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}
	if existed {
		return &rpc.PreSignPutResponse{
			FileId: fileId,
			Meta:   utils.GetSuccessMeta(),
		}, nil
	}
	url, id, err := l.fileService.PreSignPut(in.FileContext)
	if err != nil {
		return &rpc.PreSignPutResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}

	return &rpc.PreSignPutResponse{
		Meta:   utils.GetSuccessMeta(),
		Url:    url,
		FileId: id,
	}, nil
}
