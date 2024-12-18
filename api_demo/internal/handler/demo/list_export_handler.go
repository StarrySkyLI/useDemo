package demo

import (
	"errors"
	"gitlab.coolgame.world/go-template/base-common/excel"
	"gitlab.coolgame.world/go-template/base-common/result"
	"io"
	"net/http"
	"rpc_demo/rpc"

	"api_demo/internal/logic/demo"
	"api_demo/internal/svc"
	"api_demo/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func List_exportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := demo.NewList_exportLogic(r.Context(), svcCtx)
		err := l.List_export(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		client := svcCtx.DemoRPC
		stream, err := client.GameListExport(r.Context(), &rpc.GameListReq{
			Page:     req.Page,
			PageSize: req.Page_size,
		})
		if err != nil {
			result.HttpErrorResult(r.Context(), w, err)

		}
		excelWriter, err := excel.NewExcelWriter()
		if err != nil {
			result.HttpErrorResult(r.Context(), w, err)
			return
		}
		defer excelWriter.Close()
		excelWriter.SetHeaders([]string{
			"会员ID",
			"name",
			"code",
			"data",
			"create",
			"update",
		})
		index := 2
		for {
			data, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				result.HttpErrorResult(r.Context(), w, err)
				return
			}

			excelWriter.AddRow(index, []interface{}{
				data.Id,
				data.Name,
				data.Code,
				data.Data,
				data.CreateTime,
				data.UpdateTime,
			})
			index++
		}

		excelWriter.Flush()

		// download
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=users.xlsx")
		excelWriter.Write(w)
	}
}
