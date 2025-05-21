package fileservicelogic

import (
	"context"

	fileService "useDemo/application/rpc_demo/internal/service/file"
	"useDemo/application/rpc_demo/internal/svc"
	"useDemo/application/rpc_demo/rpc"
	"useDemo/application/rpc_demo/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type MergeFilePartsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	fileService *fileService.FileService
}

func NewMergeFilePartsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MergeFilePartsLogic {
	return &MergeFilePartsLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		fileService: fileService.NewFileService(ctx, svcCtx),
	}
}

// merge a slicing uploading file
func (l *MergeFilePartsLogic) MergeFileParts(in *rpc.MergeFilePartsRequest) (*rpc.MergeFilePartsResponse, error) {
	if err := l.fileService.MergeFileParts(in); err != nil {
		return &rpc.MergeFilePartsResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}

	if err := l.fileService.ReportUploaded(in.FileContext); err != nil {
		return &rpc.MergeFilePartsResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}

	return &rpc.MergeFilePartsResponse{
		Meta: utils.GetSuccessMeta(),
	}, nil
}
