package app

import (
	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.coolgame.world/go-template/base-common/app/dbM"
	"gitlab.coolgame.world/go-template/base-common/pkg/logs"
	"gitlab.coolgame.world/go-template/base-common/xxlJob"
)

func init() {
}

func logInit() {
	// 设置日志
	writer := logs.NewLogrusWriter(func(logger *logrus.Logger) {
		logger.SetFormatter(&logrus.JSONFormatter{})
	})
	logx.SetWriter(writer)
}

func InitAppServer() {
	logInit()

	logx.DisableStat()
}

type ConfigAppServer struct {
	IsDebug  bool                 `json:",optional,default=false"`
	SelectDb []dbM.SelectDbConfig `json:",optional"`
	XxlJob   xxlJob.Config        `json:",optional"`
}
