package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
	"wikitracker/pkg/models"
	"wikitracker/pkg/tools"
)

const (
	consGroup = "consumer-group-07"
	topic     = "WIKI_EDITS_COUNT_HW"
)

type Consumer struct {
	Ready chan bool
}

type countsMap struct {
	EditCounts map[string]int
	Mutex      sync.RWMutex
}

var CountsMap = &countsMap{
	EditCounts: map[string]int{},
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

func InitConsumerGroup() (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	kafkaServers := os.Getenv("KAFKA_SERVER")
	if kafkaServers == "" {
		kafkaServers = "localhost:9092"
	}

	servers := strings.Split(kafkaServers, ",")

	return sarama.NewConsumerGroup(servers, consGroup, config)
}

func StartConsuming(group sarama.ConsumerGroup) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := Consumer{
		Ready: make(chan bool),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := group.Consume(ctx, []string{topic}, &consumer); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			CountsMap.Mutex.RLock()
			fmt.Println(tools.TopK(CountsMap.EditCounts, 5))
			CountsMap.Mutex.RUnlock()
		case <-signals:
			log.Println("Interrupt is detected")
			cancel()
			wg.Wait()
			group.Close()
			return
		}
	}
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
