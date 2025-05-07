package middleware

import (
	"net/http"
	"strconv"

	"useDemo/base-common/consts"
	"useDemo/base-common/pkg/xcode"

	"google.golang.org/grpc/metadata"
	"useDemo/base-common/headInfo"
)

type CheckRequestTokenFunc func(r *http.Request, token string) int64

func AuthTokenRequest(r *http.Request, checkToken CheckRequestTokenFunc) *http.Request {
	ctx := r.Context()
	token := headInfo.GetJwtToken(ctx)
	if token != "" {
		var tokenUid int64 = 0
		if checkToken != nil {
			tokenUid = checkToken(r, token)
		}
		if tokenUid > 0 {
			ctx = metadata.AppendToOutgoingContext(ctx, consts.TokenUid, strconv.FormatInt(tokenUid, 10))
		}
	}

	newReq := r.WithContext(ctx)

	return newReq
}

func MustAuthTokenRequest(r *http.Request, checkToken CheckRequestTokenFunc) (*http.Request, error) {
	ctx := r.Context()
	token := headInfo.GetJwtToken(ctx)
	if token != "" {
		var tokenUid int64 = 0
		if checkToken != nil {
			tokenUid = checkToken(r, token)
		}
		if tokenUid > 0 {
			ctx = metadata.AppendToOutgoingContext(ctx, consts.TokenUid, strconv.FormatInt(tokenUid, 10))
		} else {
			return r, xcode.UserNotFound
		}
	}

	newReq := r.WithContext(ctx)

	return newReq, nil
}
