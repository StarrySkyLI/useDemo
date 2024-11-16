package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gitlab.coolgame.world/go-template/base-common/pkg/aesGCM"
	"gitlab.coolgame.world/go-template/base-common/result"
)

const (
	TestDecryptKey = "key111kjjfkdsflahdf9829uihfu32hu"
	ProdDecryptKey = "pro221kjjhdf982ahfu9uifkdsfl32hu"
)

var RequestDecryptError = errors.New("Request decryption failed. ")

type RequestDecryptData struct {
	ARData string `json:"ar_data"`
}

type ApiRequestDecryptMiddleware struct {
}

func NewApiRequestDecryptMiddleware(key []byte, isOpen bool) *ApiRequestDecryptMiddleware {
	aesGCM.IsOpenAesGcm = isOpen
	aesGCM.EncryptKey = key

	return &ApiRequestDecryptMiddleware{}
}

func (m *ApiRequestDecryptMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if aesGCM.IsOpenAesGcm {
			if err := m.RequestDecrypt(r); err != nil {
				result.HttpErrorResult(r.Context(), w, err)
				return
			}
		}

		next(w, r)
	}
}

func (m *ApiRequestDecryptMiddleware) RequestDecrypt(r *http.Request) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return RequestDecryptError
	}

	if len(data) == 0 || !aesGCM.IsOpenAesGcm {
		return nil
	}

	// Decrypt the data here
	var decryptData RequestDecryptData
	if err = json.Unmarshal(data, &decryptData); err != nil {
		return RequestDecryptError
	}
	if decryptData.ARData == "" {
		r.Body = io.NopCloser(bytes.NewBuffer(data))
	} else {
		deData, err := aesGCM.Decrypt(aesGCM.EncryptKey, decryptData.ARData)
		if err != nil {
			return RequestDecryptError
		}
		r.Body = io.NopCloser(bytes.NewBuffer(deData))
	}

	return nil
}
