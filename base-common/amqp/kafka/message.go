package kafkaMQ

import (
	"strings"

	"github.com/segmentio/kafka-go"
	"gitlab.coolgame.world/go-template/base-common/consts"
)

type Handler string

type Message struct {
	Data []byte `json:"data"`
}

type KafkaMessage struct {
	Trance string `json:"trance"`
	Key    []byte `json:"key"` //  = Kafka key
	Message
}

func NewKafkaMessage(msg kafka.Message) *KafkaMessage {
	var res = &KafkaMessage{
		Trance: "",
		Key:    msg.Key,
		Message: Message{
			Data: msg.Value,
		},
	}

	for i := 0; i < len(msg.Headers); i++ {
		if strings.Compare(msg.Headers[i].Key, consts.Trance) == 0 {
			res.Trance = string(msg.Headers[i].Value)
		}
	}

	return res
}

func (msg *KafkaMessage) GetHeader() []kafka.Header {
	var res []kafka.Header

	res = append(res, kafka.Header{
		Key:   consts.Trance,
		Value: []byte(msg.Trance),
	})

	return res
}
