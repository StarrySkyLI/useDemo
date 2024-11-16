package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	RPC struct {
		DemoRPC string `json:"DemoRpc"`
	}
}
