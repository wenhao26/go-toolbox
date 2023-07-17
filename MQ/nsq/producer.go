package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

func main() {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}

	topicName := "topic1"

	go func() {
		for {
			t := time.Now().String()
			messageBody := []byte("message:" + t)
			err = producer.Publish(topicName, messageBody)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Sent successfully")
			time.Sleep(2e9)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	producer.Stop()
}
