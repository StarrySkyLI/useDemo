package demo

import (
	"errors"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"useDemo/application/api_demo/internal/logic/demo"
	"useDemo/application/api_demo/internal/svc"
	"useDemo/application/api_demo/internal/types"
	"useDemo/application/rpc_demo/rpc"
	"useDemo/base-common/excel"
	"useDemo/base-common/result"
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
