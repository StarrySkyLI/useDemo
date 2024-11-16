package rpc

import (
	"fmt"

	"github.com/zeromicro/go-zero/zrpc"
	"gitlab.coolgame.world/go-template/base-common/rpc/interceptors"
)

func GenRpcTarget(hosts string) string {
	return fmt.Sprintf("%s", hosts)
}

type Config struct {
	Host     string
	AppName  string
	Business string
}

func MustNewClient(conf Config) zrpc.Client {
	return zrpc.MustNewClient(
		zrpc.RpcClientConf{
			Target: GenRpcTarget(conf.Host),
		},
		zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor(conf.AppName, conf.Business)),
		zrpc.WithDialOption(interceptors.RetryDialOption()),
	)
}

func NewClient(conf Config) (zrpc.Client, error) {
	return zrpc.NewClient(
		zrpc.RpcClientConf{
			Target: GenRpcTarget(conf.Host),
		},
		zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor(conf.AppName, conf.Business)),
		zrpc.WithDialOption(interceptors.RetryDialOption()),
	)
}
