package configServer

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logc"
)

var ConfigServer *Client
var ConfigServerErr = errors.New("ConfigServer not init")

const SystemConfig = "system_server_config"

func MustLoad(ctx context.Context, defConf Config, v any) {
	ConfigServer = NewClient(ctx, defConf)
	ConfigServer.MustStart()

	data, ok := ConfigServer.Get(SystemConfig)
	if !ok {
		logc.Must(ConfigServerErr)
	}

	strData, ok := data.(string)
	if !ok {
		logc.Must(ConfigServerErr)
	}
	err := conf.LoadFromJsonBytes([]byte(strData), v)
	if err != nil {
		logc.Must(err)
	}
}
