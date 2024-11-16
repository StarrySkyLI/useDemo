package result

import (
	"net/http"
)

type ResponseBean struct {
	Code   int         `json:"code"`
	Msg    string      `json:"message"`
	Data   interface{} `json:"data"`
	Trance string      `json:"trance"`
}

func Success(data interface{}, trance string) *ResponseBean {
	if data == nil {
		data = struct{}{}
	}
	return &ResponseBean{http.StatusOK, "Success ", data, trance}
}

func Error(errCode int, errMsg string, trance string) *ResponseBean {
	return &ResponseBean{errCode, errMsg, struct{}{}, trance}
}
