package batcher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	batcherSize     = 100
	batcherBuffer   = 100
	batcherWorker   = 10
	batcherInterval = time.Second

	chanCount   = 10
	bufferCount = 1024
)

type KafkaData struct {
	Uid int64 `json:"uid"`
	Pid int64 `json:"pid"`
}

func batchStart() {
	b := New(
		WithSize(batcherSize),
		WithBuffer(batcherBuffer),
		WithWorker(batcherWorker),
		WithInterval(batcherInterval),
	)
	b.Sharding = func(key string) int {
		pid, _ := strconv.ParseInt(key, 10, 64)
		return int(pid) % batcherWorker
	}
	b.Do = func(ctx context.Context, val map[string][]interface{}) {
		var msgs []*KafkaData
		for _, vs := range val {
			for _, v := range vs {
				msgs = append(msgs, v.(*KafkaData))
			}
		}
		kd, err := json.Marshal(msgs)
		if err != nil {
			logx.Errorf("Batcher.Do json.Marshal msgs: %v error: %v", msgs, err)
		}
		//kafka push数据
		fmt.Println(kd)
	}

	b.Start()
	b.Add("12", "21")
}
func TestBatcher(t *testing.T) {
	batchStart()
}

type Service struct {
	waiter   sync.WaitGroup
	msgsChan []chan *KafkaData
}

func NewService() *Service {
	s := &Service{
		msgsChan: make([]chan *KafkaData, chanCount),
	}
	for i := 0; i < chanCount; i++ {
		ch := make(chan *KafkaData, bufferCount)
		s.msgsChan[i] = ch
		s.waiter.Add(1)
		go s.consume(ch)
	}

	return s
}
func TestConsume(t *testing.T) {

	NewService()

}
func (s *Service) consume(ch chan *KafkaData) {
	defer s.waiter.Done()
	for {
		m, ok := <-ch
		if !ok {
			log.Fatal("seckill rmq exit")
		}
		fmt.Printf("consume msg: %+v\n", m)
	}
}
