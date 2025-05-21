package fileservicelogic

import (
	"context"

	fileService "useDemo/application/rpc_demo/internal/service/file"
	"useDemo/application/rpc_demo/internal/svc"
	"useDemo/application/rpc_demo/rpc"
	"useDemo/application/rpc_demo/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProgressRate4SlicingPutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	fileService *fileService.FileService
}

func NewGetProgressRate4SlicingPutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProgressRate4SlicingPutLogic {
	return &GetProgressRate4SlicingPutLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		fileService: fileService.NewFileService(ctx, svcCtx),
	}
}

// get upload progress rate for slicing put
func (l *GetProgressRate4SlicingPutLogic) GetProgressRate4SlicingPut(in *rpc.GetProgressRate4SlicingPutRequest) (*rpc.GetProgressRate4SlicingPutResponse, error) {
	result, err := l.fileService.GetProgressRate4SlicingPut(in.UploadId, in.FileContext)
	if err != nil {
		return &rpc.GetProgressRate4SlicingPutResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}

	total := 0
	finished := 0
	for _, uploaded := range result {
		if uploaded {
			finished++
		}

		total++
	}

	return &rpc.GetProgressRate4SlicingPutResponse{
		Meta:         utils.GetSuccessMeta(),
		Parts:        result,
		ProgressRate: float32(finished*100) / float32(total),
	}, nil
}
