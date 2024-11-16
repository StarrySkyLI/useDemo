package rpc

import (
	"context"
	"gitlab.coolgame.world/go-template/base-common/consts"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

func GetTokenUid(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}

	list := md.Get(consts.TokenUid)
	parseInt, err := strconv.ParseInt(strings.Join(list, ""), 10, 64)
	if err != nil {
		return 0
	}

	return parseInt
}

func GetTokenUidRole(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.TokenUidRole), "")
	return res
}

func GetJwtToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.Token), "")
	return res
}

func GetClientIp(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.ClientIp), "")
	return res
}
func GetUserAgent(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.UserAgent), "")
	return res
}

func GetBusiness(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.Business), "")
	return res
}

func GetVersion(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.Version), "")
	return res
}

func GetSource(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.Source), "")
	return res
}

func GetPackageName(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.PackageName), "")
	return res
}

func GetDeviceId(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.DeviceId), "")
	return res
}

func GetIMEI(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.IMEI), "")
	return res
}

func GetReqHost(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.ReqHost), "")
	return res
}

func GetReqPath(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.ReqPath), "")
	return res
}

func GetTrance(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).TraceID().String()
}

func GetContentLanguage(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.ContentLanguage), "")
	return res
}

func GetBusinessCode(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.BusinessCode), "")
	return res
}

func GetOriginHost(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	res := strings.Join(md.Get(consts.OriginHost), "")
	return res
}
