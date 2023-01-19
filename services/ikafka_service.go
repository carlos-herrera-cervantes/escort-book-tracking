package services

import (
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

//go:generate mockgen -destination=./mocks/ikafka_service.go -package=mocks --build_flags=--mod=mod . IKafkaService
type IKafkaService interface {
	SendMessage(topic string, message []byte) error
}

//go:generate mockgen -destination=./mocks/kafka.go -package=mocks --build_flags=--mod=mod . IKafka
type IKafka interface {
    Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
}
