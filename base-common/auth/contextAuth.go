package auth

import (
	"base-common/headInfo"
	"context"
)

func GetTokenUid(ctx context.Context) int64 {
	return headInfo.GetTokenUid(ctx)
}
