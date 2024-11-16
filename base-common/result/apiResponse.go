package result

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.coolgame.world/go-template/base-common/arLanguage"
	"net/http"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/rest/httpx"
	"gitlab.coolgame.world/go-template/base-common/headInfo"
	"gitlab.coolgame.world/go-template/base-common/pkg/aesGCM"
	"gitlab.coolgame.world/go-template/base-common/pkg/xcode"
)

func HttpSuccessResult(ctx context.Context, w http.ResponseWriter, resp interface{}) {
	resp = arLanguage.SwitchLanguage(resp, headInfo.GetContentLanguage(ctx))

	if aesGCM.IsOpenAesGcm {
		respByte, err := interfaceToBytes(resp)
		if err != nil {
			logc.Errorf(ctx, "Response interfaceToBytes fail: %s", err.Error())
			resp = nil
		} else {
			var encryptByte string
			encryptByte, err = aesGCM.Encrypt(aesGCM.EncryptKey, respByte)
			if err != nil {
				logc.Errorf(ctx, "Response encrypt fail: %s", err.Error())
			}
			resp = map[string]string{"ar_data": encryptByte}
		}
	}

	success := Success(resp, trace.TraceIDFromContext(ctx))

	go func() {
		logSucc, _ := json.Marshal(success)
		logc.Info(ctx, "ApiResponse:", fmt.Sprintf("%s", string(logSucc)))
	}()

	httpx.WriteJsonCtx(ctx, w, http.StatusOK, success)
}

func HttpSuccessResultNoAes(ctx context.Context, w http.ResponseWriter, resp interface{}) {
	resp = arLanguage.SwitchLanguage(resp, headInfo.GetContentLanguage(ctx))

	success := Success(resp, trace.TraceIDFromContext(ctx))

	go func() {
		logSucc, _ := json.Marshal(success)
		logc.Info(ctx, "ApiResponse:", fmt.Sprintf("%s", string(logSucc)))
	}()

	httpx.WriteJsonCtx(ctx, w, http.StatusOK, success)
}

// http param error
func HttpErrorResult(ctx context.Context, w http.ResponseWriter, err error) {
	var (
		xerr xcode.XCode
		code int
		msg  string
	)
	if errors.As(err, &xerr) {
		code = xerr.Code()
		msg = xerr.Error()
	} else {
		code = http.StatusBadRequest
		msg = err.Error()
	}

	msg = arLanguage.SwitchLanguageDic(msg, headInfo.GetContentLanguage(ctx)).(string)

	resp := Error(code, msg, trace.TraceIDFromContext(ctx))

	go func() {
		logSuc, _ := json.Marshal(resp)
		logc.Info(ctx, "ApiResponse:", string(logSuc))
	}()

	httpx.WriteJsonCtx(ctx, w, http.StatusOK, resp)
}

func interfaceToBytes(data interface{}) ([]byte, error) {
	switch v := data.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data: %w", err)
		}
		return jsonBytes, nil
	}
}
