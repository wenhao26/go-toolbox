package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

type ConsumerHandler struct{}

// 处理消息
func (h *ConsumerHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		fmt.Println("No Message...")
		return nil
	}

	fmt.Println("receive", m.NSQDAddress, "message:", string(m.Body))
	return nil
}

func main() {
	topicName := "topic1"
	channelName := "channel1"

	config := nsq.NewConfig()
	config.LookupdPollInterval = time.Second
	consumer, err := nsq.NewConsumer(topicName, channelName, config)
	if err != nil {
		panic(err)
	}

	consumer.SetLogger(nil, 0)
	consumer.AddHandler(&ConsumerHandler{})

	if err := consumer.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	consumer.Stop()
}
