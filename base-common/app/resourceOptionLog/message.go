package resourceOptionLog

import (
	"context"
	"time"

	"gitlab.coolgame.world/go-template/base-common/consts"
	"gitlab.coolgame.world/go-template/base-common/headInfo"
)

type OptionType string
type UserType int

type LogMessage struct {
	UserName   string         // 操作用户名
	UserType   UserType       // 操作用户类型
	Title      string         // 标题
	Status     consts.IntBool // 状态
	Describe   string         // 内容描述
	OptionType OptionType     // 操作类型
	LogType    consts.LogType // 日志类型
}

// 资源操作日志消息
type ResourceLogMessage struct {
	LogMessage

	Trace   string
	UserIP  string
	Browser string
	Time    time.Time
}

func GetResourceLogMessage(ctx context.Context, log LogMessage) ResourceLogMessage {
	return ResourceLogMessage{
		LogMessage: log,
		Trace:      headInfo.GetTrance(ctx),
		UserIP:     headInfo.GetClientIp(ctx),
		Browser:    headInfo.GetUserAgent(ctx),
		Time:       time.Now(),
	}
}
