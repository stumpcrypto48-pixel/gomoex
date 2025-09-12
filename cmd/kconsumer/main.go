package main

import (
	"context"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

func main() {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "GoLangTopic1",
		GroupID: "test_consumer_group6",
	})
	ctx := context.Background()
	var wg sync.WaitGroup

	jobs := make(chan kafka.Message, 100)
	done := make(chan kafka.Message, 100)

	// Reader goroutine
	go func() {
		for {
			m, err := consumer.ReadMessage(ctx)
			if err != nil {
				log.Print(err)
				close(jobs)
				return
			}
			jobs <- m
		}
	}()

	// Worker pool
	workerCount := 4
	for i := range workerCount {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for m := range jobs {
				log.Printf("Worker %d received message: %s | offset: %d \n", id, string(m.Value), m.Offset)
				done <- m
			}
		}(i)
	}

	// Commiter goroutine
	var commitWG sync.WaitGroup
	commitWG.Add(1)
	go func() {
		defer commitWG.Done()
		for m := range done {
			if err := consumer.CommitMessages(ctx, m); err != nil {
				log.Printf("Failed to commit message: %v \n", err)
			} else {
				log.Printf("Message commited successfully :: %v \n", m)
			}
		}
		log.Println("Committer exiting")
	}()

	wg.Wait()
	commitWG.Wait()
	close(jobs)
	close(done)

}

func errorDecorator(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		} else {
			log.Print("Gracefull ending")
		}
	}()
	fn()
}
