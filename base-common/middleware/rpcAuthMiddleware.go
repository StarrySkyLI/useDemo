package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"gitlab.coolgame.world/go-template/base-common/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logc"
	"gitlab.coolgame.world/go-template/base-common/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type RpcAuthMiddleware struct {
	Debug bool
}

func NewRpcAuthMiddleware() *RpcAuthMiddleware {
	return &RpcAuthMiddleware{}
}

func (m *RpcAuthMiddleware) Handle() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			// recover from panic if one occurred. Set err to nil otherwise.
			if p := recover(); p != nil {
				if m.Debug {
					debug.PrintStack()
				}
				logc.Error(ctx, err, string(debug.Stack()))
				err = xcode.NewRpc(http.StatusInternalServerError, "Rpc Server error.")
				return
			}
		}()

		logc.Info(ctx, "RpcRequest:", req)

		_, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("no rpc auth metadata. ")
		}
		ctx = rpc.ContextMetadataInLog(ctx)

		resp, err = handler(ctx, req)
		if err != nil {
			logc.Error(ctx, fmt.Sprintf(info.FullMethod+",rpc错误信息：%v", err))
		}

		logc.Info(ctx, "RpcResponse:", resp)

		return resp, err
	}
}
