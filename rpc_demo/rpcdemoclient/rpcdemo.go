// Code generated by goctl. DO NOT EDIT.
// Source: rpc_demo.proto

package rpcdemoclient

import (
	"context"

	"rpc_demo/rpc_demo"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	Request  = rpc_demo.Request
	Response = rpc_demo.Response

	RpcDemo interface {
		Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
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
	client := rpc_demo.NewRpcDemoClient(m.cli.Conn())
	return client.Ping(ctx, in, opts...)
}