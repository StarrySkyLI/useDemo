package factories

type Config struct {
	MappingCode string
}

// Sender 是一个统一的短信发送接口
type Sender interface {
	SendMessage(to, message string) error
}

// SmsFactory 是一个抽象工厂接口
type SmsFactory interface {
	CreateSmsSender() Sender
}

func NewFactory(c Config) Sender {
	if c.MappingCode == "" {
		c.MappingCode = "aws"
	}
	var factory SmsFactory
	switch c.MappingCode {
	case "aliyun":
		factory = &AliyunSmsFactory{}
	case "aws":
		factory = &AwsSmsFactory{}
	default:
		return nil
	}
	return factory.CreateSmsSender()
}
