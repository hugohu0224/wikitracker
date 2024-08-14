package internal

import (
	"bufio"
	"fmt"
	"github.com/IBM/sarama"
	"net/http"
	"strings"
)

func Streaming(producer sarama.SyncProducer, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading: %w", err)
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			event, err := JsonToEvent(data)
			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				continue
			}
			fmt.Printf("Received event: %+v\n", event)

			msg := &sarama.ProducerMessage{
				Topic: "wikimedia_events",
				Value: sarama.StringEncoder(data),
			}
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				fmt.Println("Error sending message to Kafka:", err)
			} else {
				fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
			}
		}
	}
}
