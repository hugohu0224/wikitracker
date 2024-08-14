package internal

import (
	"fmt"
	"github.com/IBM/sarama"
	"os"
	"strings"
)

func GetProducer() sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	kafkaServers := os.Getenv("KAFKA_SERVER")
	if kafkaServers == "" {
		kafkaServers = "localhost:9092"
	}

	servers := strings.Split(kafkaServers, ",")

	producer, err := sarama.NewSyncProducer(servers, config)
	if err != nil {
		fmt.Println("Error creating Kafka producer:", err)
		return nil
	}
	return producer
}
