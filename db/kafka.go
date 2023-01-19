package db

import (
    "log"

    "escort-book-tracking/config"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

var singletonProducer *kafka.Producer

func NewProducer() *kafka.Producer {
    if singletonProducer != nil {
        return singletonProducer
    }

    lock.Lock()
    defer lock.Unlock()

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.InitKafkaConfig().BootstrapServers,
		"client.id":         config.InitKafkaConfig().GroupId,
	})

	if err != nil {
		log.Panic("ERROR CREATING A PRODUCER: ", err)
	}

    singletonProducer = p

	return singletonProducer
}
