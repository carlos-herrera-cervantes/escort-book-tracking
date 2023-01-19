package config

import (
    "os"
    "sync"
)

type kafkaConfig struct {
	BootstrapServers string
	GroupId string
	Topics topic
}

type topic struct {
	OperationTopic string
}

var singletonKafkaConfig *kafkaConfig
var lock = &sync.Mutex{}

func InitKafkaConfig() *kafkaConfig {
	if singletonKafkaConfig != nil {
		return singletonKafkaConfig
	}

	lock.Lock()
	defer lock.Unlock()

	singletonKafkaConfig = &kafkaConfig{
	    BootstrapServers: os.Getenv("KAFKA_SERVERS"),
	    GroupId: os.Getenv("KAFKA_CLIENT_ID"),
	    Topics: topic{
	        OperationTopic: "operations-statistics",
        },
    }

	return singletonKafkaConfig
}
