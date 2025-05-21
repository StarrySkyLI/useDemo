package fileservicelogic

import (
	"context"

	fileService "useDemo/application/rpc_demo/internal/service/file"
	"useDemo/application/rpc_demo/internal/svc"
	"useDemo/application/rpc_demo/rpc"
	"useDemo/application/rpc_demo/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileInfoByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	fileService *fileService.FileService
}

func NewGetFileInfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileInfoByIdLogic {
	return &GetFileInfoByIdLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		fileService: fileService.NewFileService(ctx, svcCtx),
	}
}

func (l *GetFileInfoByIdLogic) GetFileInfoById(in *rpc.GetFileInfoByIdRequest) (*rpc.GetFileInfoByIdResponse, error) {
	f, err := l.fileService.GetInfoById(in.DomainName, in.BizName, in.FileId)
	if err != nil {

		return &rpc.GetFileInfoByIdResponse{
			Meta: utils.GetMetaWithError(err),
		}, nil
	}

	return &rpc.GetFileInfoByIdResponse{
		Meta:       utils.GetSuccessMeta(),
		ObjectName: utils.GetObjectName(in.BizName, in.FileId),
		Hash:       f.Hash,
	}, nil
}
