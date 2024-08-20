package main

import (
	"log"
	"wikitracker/internal"
)

func main() {
	group, err := internal.InitConsumerGroup()
	if err != nil {
		log.Fatalf("Failed to initialize consumer group: %v", err)
	}
	internal.StartConsuming(group)
}
