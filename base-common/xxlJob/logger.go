package xxlJob

import (
	"context"

	"github.com/zeromicro/go-zero/core/logc"
)

type Logger struct {
	Ctx context.Context
}

func (l Logger) Info(format string, a ...interface{}) {
	logc.Infof(l.Ctx, format, a...)
}
func (l Logger) Error(format string, a ...interface{}) {
	logc.Errorf(l.Ctx, format, a...)
}
