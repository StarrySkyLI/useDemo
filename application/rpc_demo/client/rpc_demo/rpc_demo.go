// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.3
// Source: demo_service.proto

package rpc_demo

import (
	"context"

	"useDemo/application/rpc_demo/rpc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	FileContext                        = rpc.FileContext
	GameInfo                           = rpc.GameInfo
	GameInfoRep                        = rpc.GameInfoRep
	GameInfoReq                        = rpc.GameInfoReq
	GameListRep                        = rpc.GameListRep
	GameListReq                        = rpc.GameListReq
	GetFileInfoByIdRequest             = rpc.GetFileInfoByIdRequest
	GetFileInfoByIdResponse            = rpc.GetFileInfoByIdResponse
	GetProgressRate4SlicingPutRequest  = rpc.GetProgressRate4SlicingPutRequest
	GetProgressRate4SlicingPutResponse = rpc.GetProgressRate4SlicingPutResponse
	MergeFilePartsRequest              = rpc.MergeFilePartsRequest
	MergeFilePartsResponse             = rpc.MergeFilePartsResponse
	Metadata                           = rpc.Metadata
	PaginationRequest                  = rpc.PaginationRequest
	PaginationResponse                 = rpc.PaginationResponse
	PreSignGetRequest                  = rpc.PreSignGetRequest
	PreSignGetResponse                 = rpc.PreSignGetResponse
	PreSignPutRequest                  = rpc.PreSignPutRequest
	PreSignPutResponse                 = rpc.PreSignPutResponse
	PreSignSlicingPutRequest           = rpc.PreSignSlicingPutRequest
	PreSignSlicingPutResponse          = rpc.PreSignSlicingPutResponse
	RemoveFileRequest                  = rpc.RemoveFileRequest
	RemoveFileResponse                 = rpc.RemoveFileResponse
	ReportUploadedFilePartsRequest     = rpc.ReportUploadedFilePartsRequest
	ReportUploadedFilePartsResponse    = rpc.ReportUploadedFilePartsResponse
	ReportUploadedRequest              = rpc.ReportUploadedRequest
	ReportUploadedResponse             = rpc.ReportUploadedResponse
	Request                            = rpc.Request
	Response                           = rpc.Response
	SearchField                        = rpc.SearchField
	SearchRequest                      = rpc.SearchRequest
	SortField                          = rpc.SortField

	RpcDemo interface {
		Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
		FindOne(ctx context.Context, in *GameInfoReq, opts ...grpc.CallOption) (*GameInfoRep, error)
		GameList(ctx context.Context, in *GameListReq, opts ...grpc.CallOption) (*GameListRep, error)
		GameListExport(ctx context.Context, in *GameListReq, opts ...grpc.CallOption) (rpc.RpcDemo_GameListExportClient, error)
	}

	defaultRpcDemo struct {
		cli zrpc.Client
	}
)

func NewRpcDemo(cli zrpc.Client) RpcDemo {
	return &defaultRpcDemo{
		cli: cli,
	}
}

func (m *defaultRpcDemo) Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	client := rpc.NewRpcDemoClient(m.cli.Conn())
	return client.Ping(ctx, in, opts...)
}

func (m *defaultRpcDemo) FindOne(ctx context.Context, in *GameInfoReq, opts ...grpc.CallOption) (*GameInfoRep, error) {
	client := rpc.NewRpcDemoClient(m.cli.Conn())
	return client.FindOne(ctx, in, opts...)
}

func (m *defaultRpcDemo) GameList(ctx context.Context, in *GameListReq, opts ...grpc.CallOption) (*GameListRep, error) {
	client := rpc.NewRpcDemoClient(m.cli.Conn())
	return client.GameList(ctx, in, opts...)
}

func (m *defaultRpcDemo) GameListExport(ctx context.Context, in *GameListReq, opts ...grpc.CallOption) (rpc.RpcDemo_GameListExportClient, error) {
	client := rpc.NewRpcDemoClient(m.cli.Conn())
	return client.GameListExport(ctx, in, opts...)
}
