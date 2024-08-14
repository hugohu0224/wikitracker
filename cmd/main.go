package main

import (
	"fmt"
	"io"
	"time"
	"wikitracker/internal"
)

func main() {
	producer := internal.GetProducer()
	defer producer.Close()

	url := "https://stream.wikimedia.org/v2/stream/recentchange"

	for {
		err := internal.Streaming(producer, url)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed. Retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
				continue
			}
			fmt.Println("Error processing stream:", err)
			time.Sleep(5 * time.Second)
		}
	}
}
