package demo

import (
	"net/http"

	"api_demo/internal/logic/demo"
	"api_demo/internal/svc"
	"api_demo/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindByIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := demo.NewFindByIdLogic(r.Context(), svcCtx)
		resp, err := l.FindById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
