package rpc

import (
	"base-common/consts"
	"base-common/headInfo"
	"context"
	"errors"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

const (
	TokenUid = consts.TokenUid
)

var NoRpcMetadata = errors.New("No rpc metadata. ")

type AppRpcMetadata struct {
	AppName  string `json:"app_name"` // 项目名称
	Business string `json:"business"` // 所属业务

	// token
	Token string `json:"token"`

	// 唯一ID
	//TranceId string `json:"trance"`

	TranceIdRole string `json:"trace_id_role"`

	/**
	具备 TokenUid 时 Token 失效
	*/
	TokenUid int64 `json:"token_uid"` // 用户ID

	// IP
	ClientIp string `json:"client_ip"`

	// 设备ID
	DeviceId string `json:"device_id"`

	// 用户浏览器
	UserAgent string `json:"user_agent"`

	ReqHost         string
	ReqPath         string
	Version         string
	Source          string
	PackageName     string
	IMEI            string
	ContentLanguage string
	BusinessCode    string `json:"business_code"` // 商户code
	OriginHost      string `json:"origin_host"`   // 来源域名
}

func newOutgoingContextRpc(ctx context.Context, data *AppRpcMetadata) context.Context {
	md := metadata.Pairs(
		consts.AppName, data.AppName,
		consts.Business, data.Business,
		consts.Token, data.Token,
		consts.TokenUid, strconv.FormatInt(data.TokenUid, 10),
		consts.ClientIp, data.ClientIp,
		consts.DeviceId, data.DeviceId,
		//consts.TranceId, data.TranceId,
		consts.UserAgent, data.UserAgent,
		consts.ReqHost, data.ReqHost,
		consts.ReqPath, data.ReqPath,
		consts.Version, data.Version,
		consts.Source, data.Source,
		consts.PackageName, data.PackageName,
		consts.IMEI, data.IMEI,
		consts.ContentLanguage, data.ContentLanguage,
		consts.BusinessCode, "",
		consts.OriginHost, data.OriginHost,
	)

	ctxNew := metadata.NewOutgoingContext(ctx, md)
	return ctxNew
}

func ContextMetadataInLog(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	var fieldList = make([]logx.LogField, 0)
	for k, v := range md {
		fieldList = append(fieldList, logx.Field(k, v))
	}
	ctxNew := logx.ContextWithFields(ctx, fieldList...)

	return ctxNew
}

func HeadInMetadata(ctx context.Context, h headInfo.Head) context.Context {
	return newOutgoingContextRpc(ctx, &AppRpcMetadata{
		AppName:         "",
		Business:        h.Business,
		Token:           h.AuthorizationJwt,
		TokenUid:        0,
		ClientIp:        h.ClientIp,
		DeviceId:        h.DeviceId,
		UserAgent:       h.UserAgent,
		ReqHost:         h.ReqHost,
		ReqPath:         h.ReqPath,
		Version:         h.Version,
		Source:          h.Source,
		PackageName:     h.PackageName,
		IMEI:            h.IMEI,
		ContentLanguage: h.ContentLanguage,
		OriginHost:      h.OriginHost,
		//TraceId:         headInfo.GetTrance(ctx),
	})
}
