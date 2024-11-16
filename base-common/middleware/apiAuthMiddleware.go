package middleware

import (
	"gitlab.coolgame.world/go-template/base-common/result"
	"net/http"
)

type MiddlewareOption func(m *UserAgentMiddleware)

func WithCheckOption(check CheckRequestTokenFunc) MiddlewareOption {
	return func(m *UserAgentMiddleware) {
		m.check = check
	}
}

type UserAgentMiddleware struct {
	check CheckRequestTokenFunc
}

func NewUserAgentMiddleware(ops ...MiddlewareOption) *UserAgentMiddleware {
	res := &UserAgentMiddleware{}
	for _, op := range ops {
		op(res)
	}
	return res
}

func (m *UserAgentMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newReq, err := MustAuthTokenRequest(r, m.check)
		if err != nil {
			result.HttpErrorResult(r.Context(), w, err)
			return
		}
		next(w, newReq)
	}
}
