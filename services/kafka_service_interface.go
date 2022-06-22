package services

import "context"

type IKafkaService interface {
	SendMessage(ctx context.Context, topic string, message []byte) error
}
