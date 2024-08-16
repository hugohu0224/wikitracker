package internal

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"strings"
	"sync"
	"wikitracker/pkg/models"
	"wikitracker/pkg/tools"
)

type countsMap struct {
	EditCounts map[string]int
	Mutex      sync.RWMutex
}

var CountsMap = &countsMap{
	EditCounts: map[string]int{},
}

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

type Consumer struct {
	Ready chan bool
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {

		var wikiEdit models.WikiEdit
		err := json.Unmarshal(message.Value, &wikiEdit)
		if err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		key := tools.ParseKey(message.Key)
		CountsMap.Mutex.Lock()
		CountsMap.EditCounts[key] = wikiEdit.EditsCount
		CountsMap.Mutex.Unlock()
		session.MarkMessage(message, "")
	}
	return nil
}
