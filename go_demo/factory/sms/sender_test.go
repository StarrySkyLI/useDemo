package sms

import (
	"go_demo/factory/sms/factories"
	"testing"
)

func TestSender(t *testing.T) {
	// 选择阿里云工厂

	config := factories.Config{
		MappingCode: "aliyun",
	}
	aliyunSender := factories.NewFactory(config)
	aliyunSender.SendMessage("1234567890", "这是来自阿里云的短信")
	config = factories.Config{
		MappingCode: "aws",
	}
	// 选择AWS工厂

	awsSender := factories.NewFactory(config)
	awsSender.SendMessage("0987654321", "这是来自AWS的短信")
}
