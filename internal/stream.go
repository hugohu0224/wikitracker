package internal

import (
	"bufio"
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"net/http"
	"strings"
	"time"
)

func Streaming(ctx context.Context, producer sarama.SyncProducer, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading: %w", err)
			}

			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				event, err := JsonToEvent(data)
				if err != nil {
					fmt.Printf("Error parsing JSON: %v\n", err)
					continue
				}
				fmt.Printf("Received event: %+v\n", event)

				msg := &sarama.ProducerMessage{
					Topic: "wiki_events",
					Value: sarama.StringEncoder(data),
				}
				partition, offset, err := producer.SendMessage(msg)
				if err != nil {
					fmt.Printf("Error sending message to Kafka: %v\n", err)
				} else {
					fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
				}
			}
		}
	}
}
