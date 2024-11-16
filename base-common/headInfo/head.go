package headInfo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gitlab.coolgame.world/go-template/base-common/consts"
	"net"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	HeaderToken = consts.HeaderToken
)

type Head struct {
	AuthorizationJwt string // 用户token
	Version          string // APP版本
	IMEI             string // APP 手机设备IMEI
	DeviceId         string // 设备号
	Source           string // 来源渠道	* Android * Ios * Pc
	Business         string // 业务名称，例如 AR 【ar项目】
	PackageName      string // APP包名字
	Nonce            string // 随机64位数+拼接毫秒级时间戳
	HeadSign         string // 接口签名
	ClientIp         string // 客户端IP
	UserAgent        string // 用户浏览器
	ReqHost          string // 请求地址
	ReqPath          string // 请求路径
	ContentLanguage  string // 系统语言
	OriginHost       string // 来源域名
}

func GetHead(r *http.Request) *Head {
	header := r.Header

	headMarshal, _ := json.Marshal(header)
	logc.Info(context.Background(), "request head data:"+string(headMarshal))

	return &Head{
		AuthorizationJwt: strings.Trim(header.Get(consts.HeaderToken), " "),
		Version:          strings.Trim(header.Get("Version"), " "),
		IMEI:             strings.Trim(header.Get("IMEI"), " "),
		DeviceId:         strings.Trim(header.Get("Device_id"), " "),
		Source:           strings.Trim(header.Get("Source"), " "),
		Business:         strings.Trim(header.Get("Business"), " "),
		PackageName:      strings.Trim(header.Get("PackageName"), " "),
		Nonce:            strings.Trim(header.Get("Nonce"), " "),
		HeadSign:         strings.Trim(header.Get("HeadSign"), " "),
		ClientIp:         getClientIP(r),
		UserAgent:        strings.Trim(header.Get("User-Agent"), " "),
		ReqHost:          r.Host,
		ReqPath:          r.URL.Path,
		ContentLanguage:  strings.Trim(header.Get("Content-Language"), " "),
		OriginHost:       strings.Trim(header.Get("Origin"), " "),
	}
}

func (h *Head) Verify() error {
	if h.Version == "" {
		return errors.New("Head Version empty ")
	}
	if h.Source == "" {
		return errors.New("Head Source empty ")
	}
	if h.Business == "" {
		return errors.New("Head Business empty ")
	}

	return nil
}

func (h *Head) String() string {
	data, _ := json.Marshal(h)
	return string(data)
}

func ContextHeadInLog(ctx context.Context, h *Head) context.Context {
	ctxNew := logx.ContextWithFields(ctx,
		logx.Field(consts.HeaderToken, h.AuthorizationJwt),
		logx.Field("Version", h.Version),
		logx.Field("IMEI", h.IMEI),
		logx.Field("DeviceId", h.DeviceId),
		logx.Field("Source", h.Source),
		logx.Field("PackageName", h.PackageName),
		logx.Field("Nonce", h.Nonce),
		logx.Field("HeadSign", h.HeadSign),
		logx.Field("ClientIp", h.ClientIp),
		logx.Field("UserAgent", h.UserAgent),
		logx.Field("UserAgent", h.UserAgent),
		logx.Field("ContentLanguage", h.ContentLanguage),
		logx.Field("Origin", h.OriginHost),
	)

	return ctxNew
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("x_forwarded_realip")
	if ip == "" {
		ip = r.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func GetFullHead(r *http.Request) map[string][]string {
	headers := make(map[string][]string)

	for k, v := range r.Header {
		headers[k] = v
	}

	return headers
}
