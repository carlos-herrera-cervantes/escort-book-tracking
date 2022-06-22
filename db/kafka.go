package db

import (
	"log"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	producer *Producer
	once     sync.Once
)

type Producer struct {
	KafkaProducer *kafka.Producer
}

func initProducer() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_SERVERS"),
		"client.id":         os.Getenv("KAFKA_CLIENT_ID"),
		"acks":              "all",
	})

	if err != nil {
		log.Panic("ERROR CREATING A PRODUCER: ", err)
	}

	producer = &Producer{
		KafkaProducer: p,
	}
}

func NewProducer() *Producer {
	once.Do(initProducer)
	return producer
}
