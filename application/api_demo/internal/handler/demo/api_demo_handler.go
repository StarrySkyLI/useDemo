package demo

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"useDemo/application/api_demo/internal/logic/demo"
	"useDemo/application/api_demo/internal/svc"
	"useDemo/application/api_demo/internal/types"
)

func Api_demoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := demo.NewApi_demoLogic(r.Context(), svcCtx)
		resp, err := l.Api_demo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
