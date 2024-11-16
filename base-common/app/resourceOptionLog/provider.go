package resourceOptionLog

import (
	"context"
	"encoding/json"
	"errors"

	kafkaMQ "gitlab.coolgame.world/go-template/base-common/amqp/kafka"
	"gitlab.coolgame.world/go-template/base-common/amqp/kafka/config"
)

const (
	ResourceLogTopic = "resource_log_topic"
)

func GetProviderConfig(brokers []string) config.ProviderConfig {
	return config.ProviderConfig{
		Brokers: brokers,
		Topic:   ResourceLogTopic,
	}
}

var resourceLogProducer *kafkaMQ.Producer

func WithConnBroker(brokers []string) {
	WithConn(GetProviderConfig(brokers))
}

func WithConn(conf config.ProviderConfig) {
	resourceLogProducer = kafkaMQ.NewProducer(conf)
}

func Close() error {
	return resourceLogProducer.Close()
}

func ProduceMessage(ctx context.Context, msg LogMessage) error {
	if resourceLogProducer == nil {
		return errors.New("Resource provider not connect. ")
	}

	opLogMsg := GetResourceLogMessage(ctx, msg)

	mlp, err := json.Marshal(opLogMsg)
	if err != nil {
		return err
	}

	err = resourceLogProducer.ProduceMessage(&kafkaMQ.KafkaMessage{
		Trance: opLogMsg.Trace,
		Key:    []byte(opLogMsg.OptionType),
		Message: kafkaMQ.Message{
			Data: mlp,
		},
	})
	if err != nil {
		return err

	}

	return nil
}
