package middleware

import (
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gitlab.coolgame.world/go-template/base-common/headInfo"
	"gitlab.coolgame.world/go-template/base-common/result"
	"gitlab.coolgame.world/go-template/base-common/rpc"
	"net/http"
	"runtime/debug"
	"strings"
)

type ApiHeadOption func(m *ApiHeaderMiddleware)

func CloseVerifyOption(path map[string]int) ApiHeadOption {
	return func(m *ApiHeaderMiddleware) {
		m.noVerifyPath = path
	}
}

func WithDebugOption() ApiHeadOption {
	return func(m *ApiHeaderMiddleware) {
		m.debug = true
	}
}

type ApiHeaderMiddleware struct {
	NotVerify    bool
	debug        bool
	noVerifyPath map[string]int
}

func NewApiHeaderMiddleware(arg ...ApiHeadOption) *ApiHeaderMiddleware {
	res := &ApiHeaderMiddleware{}
	for _, o := range arg {
		o(res)
	}

	return res
}

// Deprecated use CloseVerifyOption replace.
func (m *ApiHeaderMiddleware) SetNoVerify(b bool) *ApiHeaderMiddleware {
	m.NotVerify = b
	return m
}

func (m *ApiHeaderMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// recover from panic if one occurred. Set err to nil otherwise.
			if err := recover(); err != nil {
				if m.debug {
					debug.PrintStack()
				}
				logc.Error(r.Context(), err, string(debug.Stack()))
				result.HttpErrorResult(r.Context(), w, errors.New("Server error. "))
				return
			}
		}()
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "DNT,web-token,app-token,Authorization,Accept,Origin,Keep-Alive,User-Agent,X-Mx-ReqToken,X-Data-Type,X-Auth-Token,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,token,Cookie,Content-Type,AuthorizationJwt,Version,IMEI,DeviceId,Source,Business,PackageName,Nonce,HeadSign")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h := headInfo.GetHead(r)
		if r.Method != http.MethodGet && !m.NotVerify && m.verifyPath(r.URL.Path) {
			if err := h.Verify(); err != nil {
				result.HttpErrorResult(r.Context(), w, err)
				return
			}
		}

		newCtx := headInfo.ContextHeadInLog(r.Context(), h)
		// 设置 metadata
		newCtx = rpc.HeadInMetadata(newCtx, *h)

		newReq := r.WithContext(newCtx)

		next(w, newReq)
	}
}

func (m *ApiHeaderMiddleware) verifyPath(urlPath string) bool {
	if _, ok := m.noVerifyPath[urlPath]; ok {
		return false
	}
	//todo 优化一下
	for path, _ := range m.noVerifyPath {
		if strings.HasPrefix(path, "/") && strings.HasSuffix(path, "*") {
			prefix := strings.TrimSuffix(path, "*")
			if strings.HasPrefix(urlPath, prefix) {
				return false
			}
		}
	}
	return true
}

// Handler 跨域请求处理器
func (m *ApiHeaderMiddleware) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "DNT,web-token,app-token,Authorization,Accept,Origin,Keep-Alive,User-Agent,X-Mx-ReqToken,X-Data-Type,X-Auth-Token,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,token,Cookie,Content-Type,AuthorizationJwt,Version,IMEI,DeviceId,Source,Business,PackageName,Nonce,HeadSign")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
