package kafka

import (
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var mutex sync.Mutex

func AddClient(client *websocket.Conn) {
	mutex.Lock()
	clients[client] = true
	mutex.Unlock()
}

func RemoveClient(client *websocket.Conn) {
	mutex.Lock()
	delete(clients, client)
	mutex.Unlock()
}

func BroadcastMessage(message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("WebSocket error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func StartConsumer(brokers []string, topic string) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume Kafka partition: %v", err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		log.Printf("Kafka message received: %s", string(msg.Value))

		BroadcastMessage(string(msg.Value))
	}
}
