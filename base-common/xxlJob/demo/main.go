package main

import (
	"context"
	"fmt"
	"log"

	"github.com/xxl-job/xxl-job-executor-go"
	"github.com/xxl-job/xxl-job-executor-go/example/task"
	"gitlab.coolgame.world/go-template/base-common/xxlJob"
)

func main() {
	client := xxlJob.NewClient(context.Background(), xxlJob.Config{
		Address:      "http://127.0.0.1:8080/xxl-job-admin",
		AccessToken:  "default_token",
		ExecutorPort: "9999",
		RegistryKey:  "golang-jobs",
	}).Register(
		&TestHandler{},
	)
	client.MustStart()

	select {}
}

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

func Old() {
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://127.0.0.1:8080/xxl-job-admin"),
		xxl.AccessToken("default_token"), // 请求令牌(默认为空)
		// xxl.ExecutorIp("127.0.0.1"),      // 可自动获取
		xxl.ExecutorPort("9999"),       // 默认9999（非必填）
		xxl.RegistryKey("golang-jobs"), // 执行器名称
		// xxl.SetLogger(&logger{}),       // 自定义日志
	)
	exec.Init()
	// 设置日志查看handler
	// exec.LogHandler(customLogHandle)
	// 注册任务handler
	exec.RegTask("task.test", task.Test)
	exec.RegTask("task.test2", task.Test2)
	exec.RegTask("task.panic", task.Panic)
	log.Fatal(exec.Run())
}
