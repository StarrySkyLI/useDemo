package main

import (
	"fmt"
	"log"
	"time"

	kafkaMQ "gitlab.coolgame.world/go-template/base-common/amqp/kafka"
	"gitlab.coolgame.world/go-template/base-common/amqp/kafka/config"
)

func main() {
	producer := kafkaMQ.NewProducer(config.ProviderConfig{
		Brokers: []string{"kafka-1.leadsea.cn:9092"},
		Topic:   "testDemo",
	})
	defer producer.Close()

	consumer := kafkaMQ.NewConsumer(config.CustomerConfig{
		ProviderConfig: config.ProviderConfig{
			Brokers: []string{"kafka-1.leadsea.cn:9092"},
			Topic:   "testDemo",
		},
		GroupID: "t2", // 相同 GroupID的消费者只会消费一次
	})
	defer consumer.Close()

	// Produce a message
	go func() {
		for i := 0; ; i++ {
			itemData := fmt.Sprintf("e %d", i)
			err := producer.ProduceMessage(&kafkaMQ.KafkaMessage{
				Trance: fmt.Sprintf("tranceId-%d", i),
				Key:    []byte("samePage"),
				// Key:      []byte(itemData), // 顺序不同
				Message: kafkaMQ.Message{
					Data: []byte(itemData),
				},
			})
			if err != nil {
				log.Fatalf("could not produce message: %v", err)
			}
		}
	}()

	// Consume messages
	go func() {
		consumer.ConsumeMessages(func(message kafkaMQ.KafkaMessage) error {
			log.Printf("received message: key=%s, value=%s, tranceId=%s", message.Key, string(message.Data), message.Trance)
			time.Sleep(100 * time.Millisecond)
			return nil
		})
	}()
	// Consume messages
	go func() {
		consumer.ConsumeMessages(func(message kafkaMQ.KafkaMessage) error {
			log.Printf("received message: key=%s, value=%s, tranceId=%s", message.Key, string(message.Data), message.Trance)
			time.Sleep(100 * time.Millisecond)
			return nil
		})
	}()

	// Wait to consume messages
	select {}
}
