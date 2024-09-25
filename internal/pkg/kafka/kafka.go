package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
)

func NewKafkaProducer(config *config.Config) (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.Kafka.Brokers})
}

func NewKafkaConsumer(groupID string, config *config.Config) (*kafka.Consumer, error) {
	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Kafka.Brokers,
		"group.id":          config.Kafka.GroupID,
		"auto.offset.reset": "earliest",
	})
}
