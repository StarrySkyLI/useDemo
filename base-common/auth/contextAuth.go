package auth

import (
	"context"
	"gitlab.coolgame.world/go-template/base-common/headInfo"
)

func GetTokenUid(ctx context.Context) int64 {
	return headInfo.GetTokenUid(ctx)
}
