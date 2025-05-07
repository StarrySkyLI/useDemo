package auth

import (
	"context"
	"useDemo/base-common/headInfo"
)

func GetTokenUid(ctx context.Context) int64 {
	return headInfo.GetTokenUid(ctx)
}
