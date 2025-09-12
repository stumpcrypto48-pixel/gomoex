package main

import (
	"context"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

func main() {

	var wg sync.WaitGroup
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "GoLangTopic1",
	})
	defer w.Close()
	log.Println("Try to send message into kafka")

	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			w.WriteMessages(context.Background(), kafka.Message{
				Key:   []byte("key"),
				Value: []byte(string(i)),
			})
			wg.Done()
		}(i)
	}

	wg.Wait()
}
