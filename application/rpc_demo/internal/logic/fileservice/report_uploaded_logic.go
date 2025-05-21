package fileservicelogic

import (
	"context"

	fileService "useDemo/application/rpc_demo/internal/service/file"
	"useDemo/application/rpc_demo/internal/svc"
	"useDemo/application/rpc_demo/rpc"
	"useDemo/application/rpc_demo/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportUploadedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	fileService *fileService.FileService
}

func NewReportUploadedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportUploadedLogic {
	return &ReportUploadedLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		fileService: fileService.NewFileService(ctx, svcCtx),
	}
}

// report a file has been uploaded
func (l *ReportUploadedLogic) ReportUploaded(in *rpc.ReportUploadedRequest) (*rpc.ReportUploadedResponse, error) {
	err := l.fileService.ReportUploaded(in.FileContext)
	if err != nil {
		return &rpc.ReportUploadedResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}
	url, err := l.fileService.PreGetFile(in.FileContext)
	if err != nil {
		return &rpc.ReportUploadedResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}

	return &rpc.ReportUploadedResponse{
		Meta: utils.GetSuccessMeta(),
		Url:  url,
	}, nil
}
