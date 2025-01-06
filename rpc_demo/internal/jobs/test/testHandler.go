package test

import (
	"base-common/xxlJob"
	"context"
	"fmt"

	"github.com/xxl-job/xxl-job-executor-go"
)

type TestHandler struct {
}

func (h TestHandler) Handler() xxlJob.HandlerFun {
	return func(cxt context.Context, param *xxl.RunReq) string {
		fmt.Println("test one task" + param.ExecutorHandler + " param：" + param.ExecutorParams + " log_id:" + xxl.Int64ToStr(param.LogID))
		return "test finish..."
	}
}

func (h TestHandler) Pattern() string {
	return "task.test"
}
