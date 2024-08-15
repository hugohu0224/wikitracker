package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

type WikiEdit struct {
	Title       string `json:"TITLE"`
	EditsCount  int    `json:"EDITS_COUNT"`
	WindowStart int    `json:"WINDOW_START"`
	WindowEnd   int    `json:"WINDOW_END"`
}

func parseKey(key []byte) string {
	for i, b := range key {
		if b == 0 {
			return string(key[:i])
		}
	}
	return string(key)
}

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	kafkaServers := os.Getenv("KAFKA_SERVER")
	if kafkaServers == "" {
		kafkaServers = "localhost:9092"
	}

	servers := strings.Split(kafkaServers, ",")

	consumer, err := sarama.NewConsumer(servers, config)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer consumer.Close()

	topic := "WIKI_EDITS_COUNT"

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	editCounts := make(map[string]int)
	var mutex sync.Mutex

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var wikiEdit WikiEdit
			err := json.Unmarshal(msg.Value, &wikiEdit)
			if err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}
			key := parseKey(msg.Key)
			fmt.Println(key)
			mutex.Lock()
			editCounts[key] = wikiEdit.EditsCount
			mutex.Unlock()

		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v", err)

		case <-ticker.C:
			fmt.Println(editCounts["Q129155800"])

		case <-signals:
			log.Println("Interrupt is detected")
			return
		}
	}
}
