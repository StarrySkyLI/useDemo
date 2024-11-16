package kafkaMQ

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"

	"gitlab.coolgame.world/go-template/base-common/amqp/kafka/config"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(config config.ProviderConfig) *Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Brokers...),
		Topic:                  config.Topic,
		Balancer:               &kafka.LeastBytes{},
		MaxAttempts:            5,               // 最大重试次数
		WriteTimeout:           3 * time.Second, // 写入超时时间
		AllowAutoTopicCreation: true,
		RequiredAcks:           config.RequiredAcks,
	}

	return &Producer{
		Writer: writer,
	}
}

func (p *Producer) ProduceMessage(message *KafkaMessage) error {
	msg := kafka.Message{
		Key:     message.Key,
		Value:   message.Data,
		Headers: message.GetHeader(),
	}

	err := p.Writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.Writer.Close()
}
func (p *Producer) BatchProduceMessage(messages []*KafkaMessage) error {
	//for _, message := range messages {
	//	msg := kafka.Message{
	//		Key:     message.Key,
	//		Value:   message.Data,
	//		Headers: message.GetHeader(),
	//	}
	//	err := p.Writer.WriteMessages(context.Background(), msg)
	//	if err != nil {
	//		return err
	//	}
	//}
	//mr.MapReduceVoid[]()

	return nil
}
