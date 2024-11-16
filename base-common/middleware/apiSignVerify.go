package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"gitlab.coolgame.world/go-template/base-common/pkg/signr"
	"gitlab.coolgame.world/go-template/base-common/result"
)

const (
	TestSignKey = "1234567890abcdefg"
	ProdSignKey = "545459z923hjhsdfa"
)

var SignTimeError = errors.New("Time too late. ")
var SignError = errors.New("Sign verify failed. ")

type SignVerifyMiddleware struct {
	signKey      string
	isOpen       bool
	OutTime      time.Duration
	noVerifyPath map[string]int
}

func NewSignVerifyMiddleware(signKey string, isOpen bool, noVerifyPath map[string]int) *SignVerifyMiddleware {
	return &SignVerifyMiddleware{
		signKey:      signKey,
		isOpen:       isOpen,
		OutTime:      1 * time.Second,
		noVerifyPath: noVerifyPath,
	}
}

func (m *SignVerifyMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if m.isOpen && m.verifyPath(r.URL.Path) {
			headSign := strings.Split(r.Header.Get("HeadSign"), ".")
			if len(headSign) != 2 {
				result.HttpErrorResult(r.Context(), w, SignError)
				return
			}
			sign, reqTimeStr := headSign[0], headSign[1]

			if err := m.checkReqTime(reqTimeStr); err != nil {
				result.HttpErrorResult(r.Context(), w, SignTimeError)
				return
			}

			if sign != signr.GenerateHeadSign(
				m.signKey,
				r.Header.Get("AuthorizationJwt"),
				r.Header.Get("Nonce"),
				r.Header.Get("Version"),
			) {
				result.HttpErrorResult(r.Context(), w, SignError)
				return
			}
		}

		next(w, r)
	}
}

func (m *SignVerifyMiddleware) checkReqTime(reqTimeStr string) error {
	//reqTime, err := strconv.ParseInt(reqTimeStr, 10, 64)
	//if err != nil {
	//	return SignTimeError
	//}
	//nowTime := time.Now()
	//if reqTime > nowTime.UnixMilli() || reqTime < nowTime.Add(-m.OutTime).UnixMilli() {
	//	return SignTimeError
	//}

	return nil
}
func (m *SignVerifyMiddleware) verifyPath(urlPath string) bool {
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
