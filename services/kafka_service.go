package services

import (
	"context"
	"escort-book-tracking/db"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaService struct {
	Producer *db.Producer
}

func (k *KafkaService) SendMessage(ctx context.Context, topic string, message []byte) error {
	if err := k.Producer.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}, nil); err != nil {
		return err
	}

	return nil
}
