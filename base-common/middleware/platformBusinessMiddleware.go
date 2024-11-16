package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"gitlab.coolgame.world/go-template/base-common/consts"
	"gitlab.coolgame.world/go-template/base-common/headInfo"
	"gitlab.coolgame.world/go-template/base-common/pkg/xcode"
	"gitlab.coolgame.world/go-template/base-common/result"
	"google.golang.org/grpc/metadata"
)

type GetPlatformBusinessFunc func(ctx context.Context, dns string) string

type PlatformBusinessMiddlewareOption func(m *PlatformBusinessMiddleware)

func WithGetPlatformBusinessFuncOption(getF GetPlatformBusinessFunc) PlatformBusinessMiddlewareOption {
	return func(m *PlatformBusinessMiddleware) {
		m.getCodeF = getF
	}
}
func WithPlatformBusinessDebugOption() PlatformBusinessMiddlewareOption {
	return func(m *PlatformBusinessMiddleware) {
		m.isDebug = true
	}
}

type PlatformBusinessMiddleware struct {
	getCodeF GetPlatformBusinessFunc
	isDebug  bool
}

func NewPlatformBusinessMiddleware(opt ...PlatformBusinessMiddlewareOption) *PlatformBusinessMiddleware {
	res := &PlatformBusinessMiddleware{}
	for _, o := range opt {
		o(res)
	}

	return res
}

func (m *PlatformBusinessMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newReq, err := MustGetPlatformBusinessCode(r, m.getCodeF, m.isDebug)
		if err != nil {
			result.HttpErrorResult(r.Context(), w, err)
			return
		}
		next(w, newReq)
	}
}

func MustGetPlatformBusinessCode(r *http.Request, getCodeF GetPlatformBusinessFunc, isDebug bool) (*http.Request, error) {
	ctx := r.Context()
	dnsStr := headInfo.GetOriginHost(ctx)
	if dnsStr != "" {
		var bCode string

		bCode = consts.BusinessCodeDefaultValue
		//if isDebug {
		//	bCode = consts.BusinessCodeDefaultValue
		//} else {
		//	bCode = getBusinessCodeByCache(dnsStr)
		//}

		if bCode == "" && getCodeF != nil {
			bCode = getCodeF(r.Context(), dnsStr)
		}
		if bCode != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, consts.BusinessCode, bCode)
			go func() {
				setBusinessCodeByCache(dnsStr, bCode)
			}()
		} else {
			return r, xcode.PlatformCodeNotFound
		}
	}

	newReq := r.WithContext(ctx)

	return newReq, nil
}

var cacheBusinessCode = cache.New(10*time.Hour, 1*time.Minute)

func getBusinessCodeByCache(dns string) string {
	cRes, ok := cacheBusinessCode.Get(dns)
	if !ok {
		return ""
	}

	if res, ok := cRes.(string); !ok {
		return ""
	} else {
		return res
	}
}

func setBusinessCodeByCache(dns string, code string) {
	cacheBusinessCode.Set(dns, code, 1*time.Hour)
}
