package factories

import (
	"fmt"
)

// AliyunSmsFactory 是一个具体的工厂类，负责创建阿里云的短信发送实现
type AliyunSmsFactory struct{}

func (f *AliyunSmsFactory) CreateSmsSender() Sender {
	return &AliyunSmsSender{}
}

// AliyunSmsSender 是阿里云的短信发送实现
type AliyunSmsSender struct{}

func (a *AliyunSmsSender) SendMessage(to, message string) error {
	// 模拟阿里云短信发送的逻辑
	fmt.Printf("阿里云短信发送给 %s: %s\n", to, message)
	return nil
}
