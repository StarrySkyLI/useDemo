package config

import (
	"github.com/segmentio/kafka-go"
	"time"
)

type ConnectConfig struct {
	Network          string        `json:"network"`       // 网络协议
	Address          string        `json:"address"`       // 链接地址
	Topic            string        `json:"topic"`         // 主题
	Partition        int           `json:"partition"`     // 分区
	SendDeadlineTime time.Duration `json:"deadline_time"` // 发送超时时间
}

type CustomerConfig struct {
	ProviderConfig

	GroupID     string
	GroupTopics []string
}

type ProviderConfig struct {
	Brokers      []string
	Topic        string
	RequiredAcks kafka.RequiredAcks `json:"RequiredAcks,omitempty,default=1"`
}
