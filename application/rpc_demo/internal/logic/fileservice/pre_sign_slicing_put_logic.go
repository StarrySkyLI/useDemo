package fileservicelogic

import (
	"context"

	fileService "useDemo/application/rpc_demo/internal/service/file"
	"useDemo/application/rpc_demo/internal/svc"
	"useDemo/application/rpc_demo/rpc"
	"useDemo/application/rpc_demo/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSignSlicingPutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	fileService *fileService.FileService
}

func NewPreSignSlicingPutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSignSlicingPutLogic {
	return &PreSignSlicingPutLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		fileService: fileService.NewFileService(ctx, svcCtx),
	}
}

// pre sign a file url for user put it with slicing
func (l *PreSignSlicingPutLogic) PreSignSlicingPut(request *rpc.PreSignSlicingPutRequest) (*rpc.PreSignSlicingPutResponse, error) {
	fileId, existed, err := l.fileService.CheckFileExistedAndGetFile(request.FileContext)
	if err != nil {
		return &rpc.PreSignSlicingPutResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}

	if existed {
		return &rpc.PreSignSlicingPutResponse{
			Meta:     utils.GetSuccessMeta(),
			FileId:   fileId,
			Uploaded: true,
		}, nil
	}

	sf, err := l.fileService.PreSignSlicingPut(request.FileContext)
	if err != nil {
		return &rpc.PreSignSlicingPutResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}

	return &rpc.PreSignSlicingPutResponse{
		Meta:     utils.GetSuccessMeta(),
		Urls:     sf.UploadUrl,
		UploadId: sf.UploadId,
		Parts:    sf.TotalParts,
		FileId:   sf.File.ID,
	}, nil
}
