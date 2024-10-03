package queue

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/segmentio/kafka-go"
)

type KafkaQueue struct {
	Writer	*kafka.Writer
}

func NewKafkaQueue(ctx context.Context, addr []string, topic string) (*KafkaQueue, error) {
	writer := &kafka.Writer{
		Addr: kafka.TCP(addr...),
		Topic: topic,
		AllowAutoTopicCreation: true,
	}

	return &KafkaQueue{Writer: writer}, nil
}

func (s *KafkaQueue) WriteMessage(ctx context.Context, message []byte) error {
	err := s.Writer.WriteMessages(ctx,
		kafka.Message{Value: message},
	)
	if err != nil {
		log.Errorf("Error writing message to kafka queue: %s", err)
	}

	return err
}