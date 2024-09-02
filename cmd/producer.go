package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wikitracker/internal"
)

func main() {
	producer := internal.GetProducer()
	defer producer.Close()

	url := "https://stream.wikimedia.org/v2/stream/recentchange"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		errChan <- internal.Streaming(ctx, producer, url)
	}()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				if err == io.EOF {
					fmt.Println("connection closed. retrying in 5 seconds...")
					time.Sleep(5 * time.Second)
					go func() {
						errChan <- internal.Streaming(ctx, producer, url)
					}()
				} else {
					fmt.Println("error processing stream:", err)
					time.Sleep(5 * time.Second)
					go func() {
						errChan <- internal.Streaming(ctx, producer, url)
					}()
				}
			}
		case sig := <-signalChan:
			fmt.Printf("received signal: %v. Shutting down...\n", sig)
			cancel()
			return
		case <-ctx.Done():
			fmt.Println("context canceled, shutting down...")
			return
		}
	}
}
