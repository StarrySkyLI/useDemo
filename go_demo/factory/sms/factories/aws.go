package factories

import (
	"fmt"
)

// AwsSmsFactory 是一个具体的工厂类，负责创建AWS的短信发送实现
type AwsSmsFactory struct{}

func (f *AwsSmsFactory) CreateSmsSender() Sender {
	return &AwsSmsSender{}
}

// AwsSmsSender 是AWS的短信发送实现
type AwsSmsSender struct{}

func (a *AwsSmsSender) SendMessage(to, message string) error {
	// 模拟AWS短信发送的逻辑
	fmt.Printf("AWS短信发送给 %s: %s\n", to, message)
	return nil
}
