package kafka

import (
	"github.com/IBM/sarama"
	"log"
)

var producer sarama.SyncProducer

func InitKafkaProducer(brokerList []string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	producer, err = sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to start IBM producer: %v", err)
		return err
	}
	return nil
}

func SendMessage(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}
