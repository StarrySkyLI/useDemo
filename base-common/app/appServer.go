package app

import (
	"base-common/app/dbM"
	"base-common/pkg/logs"
	"base-common/xxlJob"
	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/logx"
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
