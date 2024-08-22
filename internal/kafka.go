package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"wikitracker/global"
	"wikitracker/pkg/models"
	"wikitracker/pkg/tools"
)

const (
	consGroup = "consumer-group-16"
	topic     = "WIKI_EDITS_COUNT_HW_2"
)

type Consumer struct {
	Ready chan bool
}

type Message struct {
	EditCounts map[string]int
	Title      string
	Url        string
	Mutex      sync.RWMutex
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
		var editInfo models.WikiEditInfo
		var editInfoKey models.WikiEditInfoKey
		err := json.Unmarshal(message.Value, &editInfo)
		if err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}
		key := tools.ParseKeyToString(message.Key)
		err = json.Unmarshal([]byte(key), &editInfoKey)
		if err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		editInfo.Title = editInfoKey.Title
		editInfo.Url = editInfoKey.Url
		timeKey := editInfo.WindowStart

		if _, exists := global.Ct.WikiEditInfo[timeKey]; !exists {
			global.Ct.Mu.Lock()
			global.Ct.WikiEditInfo[timeKey] = make(map[string]*models.WikiEditInfo)
			global.Ct.Mu.Unlock()
		}

		// testing
		global.Ct.Mu.Lock()
		global.Ct.WikiEditInfo[timeKey][editInfo.Title] = &editInfo
		global.Ct.Mu.Unlock()

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
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					log.Println("consumer group has been closed")
					return
				}
				log.Printf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-consumer.Ready
	log.Println("consumer is running")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Println("received termination signal, initiating shutdown")
		cancel()
	case <-ctx.Done():
		log.Println("context cancelled, shutting down")
	}

	wg.Wait()
	log.Println("consumer has been shut down")
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
