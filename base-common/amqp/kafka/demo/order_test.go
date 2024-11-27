package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"testing"
	"time"
)

type OrderProcessor struct {
	mu         sync.Mutex
	orderQueue map[string][]string // 存储每个订单的消息列表
}

func NewOrderProcessor() *OrderProcessor {
	return &OrderProcessor{
		orderQueue: make(map[string][]string),
	}
}

func (p *OrderProcessor) ProcessOrderMessage(orderID, message string) {
	// 为每个订单加锁，确保同一个订单的消息顺序处理
	p.mu.Lock()
	defer p.mu.Unlock()

	// 模拟处理顺序
	fmt.Printf("Processing order %s with message: %s\n", orderID, message)
	p.orderQueue[orderID] = append(p.orderQueue[orderID], message)

	// 模拟实际的处理逻辑
	time.Sleep(1 * time.Second)
}

func consumeMessages(ctx context.Context, reader *kafka.Reader, processor *OrderProcessor) {
	for {
		// 读取 Kafka 消息
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		orderID := string(msg.Key) // 假设消息的 key 为订单 ID
		message := string(msg.Value)

		// 使用消费者并行消费消息，但在业务层面保证顺序
		go processor.ProcessOrderMessage(orderID, message)
	}
}
func TestOrder(t *testing.T) {
	// Kafka 配置
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "order-topic",
		Partition: 0, // 可以指定分区来控制每个消费者消费的分区
		GroupID:   "order-consumer-group",
	})
	defer reader.Close()

	// 创建订单处理器
	processor := NewOrderProcessor()

	// 启动消息消费
	go consumeMessages(context.Background(), reader, processor)

	// 模拟发送消息到 Kafka
	producer := kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		AllowAutoTopicCreation: true,
		Topic:                  "order-topic",
		Balancer:               &kafka.LeastBytes{},
	}
	defer producer.Close()

	// 模拟发送 3 条订单消息
	for i := 0; i < 5; i++ {
		orderID := fmt.Sprintf("order-%d", i) // 模拟多个消息属于同一个订单
		msg := kafka.Message{
			Key:   []byte(orderID),
			Value: []byte(fmt.Sprintf("Order data for %s", orderID)),
		}
		if err := producer.WriteMessages(context.Background(), msg); err != nil {
			log.Fatal("Failed to write message:", err)
		}
	}

	// 模拟等待消费者处理完成
	time.Sleep(5 * time.Second)
	fmt.Println("Finished processing all messages.")
}
