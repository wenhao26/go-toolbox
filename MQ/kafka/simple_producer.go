package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

type Message struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	addr := []string{
		"192.168.1.216:9092",
		"192.168.1.217:9092",
		"192.168.1.218:9092",
	}
	producer, err := sarama.NewSyncProducer(addr, config)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	msg := Message{
		Id:   strconv.FormatInt(time.Now().Unix(), 10),
		Text: time.Now().String(),
	}
	data, _ := json.Marshal(msg)

	producerMsg := &sarama.ProducerMessage{
		Topic: "dev_topic_001",
		//Key:   sarama.StringEncoder("dev"),
		Value: sarama.ByteEncoder(data),
	}
	message, offset, err := producer.SendMessage(producerMsg)
	if err != nil {
		panic(err)
	}

	fmt.Println("-Partition:", message, "  -Offset:", offset)
}
