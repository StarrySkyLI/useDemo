package interceptors

import (
	"context"

	"gitlab.coolgame.world/go-template/base-common/xcode"
	"google.golang.org/grpc"
)

func ServerErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		return resp, xcode.FromError(err).Err()
	}
}
