package xxlJob

import (
	"github.com/xxl-job/xxl-job-executor-go"
)

type HandlerFun xxl.TaskFunc

type Handler interface {
	Handler() HandlerFun
	Pattern() string
}
