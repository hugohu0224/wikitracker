package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
	"wikitracker/internal"
	"wikitracker/pkg/tools"

	"github.com/IBM/sarama"
)

const (
	consGroup = "consumer-group-07"
	topic     = "WIKI_EDITS_COUNT_HW"
)

func main() {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	kafkaServers := os.Getenv("KAFKA_SERVER")
	if kafkaServers == "" {
		kafkaServers = "localhost:9092"
	}

	servers := strings.Split(kafkaServers, ",")

	group, err := sarama.NewConsumerGroup(servers, consGroup, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer group.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	ticker := time.NewTicker(2 * time.Second)

	consumer := internal.Consumer{
		Ready: make(chan bool),
	}

	go func() {
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

	for {
		select {
		case <-ticker.C:
			internal.CountsMap.Mutex.RLock()
			fmt.Println(tools.TopK(internal.CountsMap.EditCounts, 5))
			internal.CountsMap.Mutex.RUnlock()
		case <-signals:
			log.Println("Interrupt is detected")
			cancel()
			return
		}
	}
}
