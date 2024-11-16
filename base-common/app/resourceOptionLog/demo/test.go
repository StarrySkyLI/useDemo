package main

import (
	"context"
	"fmt"

	"gitlab.coolgame.world/go-template/base-common/app/resourceOptionLog"
	"gitlab.coolgame.world/go-template/base-common/consts"
)

func main() {
	resourceOptionLog.WithConnBroker([]string{"kafka-1.leadsea.cn:9092"})
	defer resourceOptionLog.Close()

	err := resourceOptionLog.ProduceMessage(context.Background(), resourceOptionLog.LogMessage{
		UserName:   "testName14",
		UserType:   1,
		Title:      "test4",
		Status:     consts.IntBoolDbTrue,
		Describe:   "测试操作一下 用户4",
		OptionType: "财务管理4",
		LogType:    consts.SystemLoginLogs,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
