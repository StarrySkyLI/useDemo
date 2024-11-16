package kafkaMQ

import (
	"context"

	"github.com/segmentio/kafka-go"
	"gitlab.coolgame.world/go-template/base-common/amqp/kafka/config"
)

type ConsumerHandler func(message KafkaMessage) error

type Consumer struct {
	Reader *kafka.Reader
}

func NewConsumer(conf config.CustomerConfig) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        conf.Brokers,
		GroupID:        conf.GroupID,
		Topic:          conf.Topic,
		GroupTopics:    conf.GroupTopics,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 0,    // 禁用自动提交
		StartOffset:    kafka.FirstOffset,
	})

	return &Consumer{
		Reader: reader,
	}
}

func (c *Consumer) ConsumeMessages(handler ConsumerHandler) error {
	for {
		msg, err := c.Reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		err = handler(*NewKafkaMessage(msg))
		if err != nil {
			// 业务错误，由业务方自行处理
			continue
		}

		// 手动提交偏移量
		if err := c.Reader.CommitMessages(context.Background(), msg); err != nil {
			return err
		}
	}
}

func (c *Consumer) Close() error {
	return c.Reader.Close()
}
