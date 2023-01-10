package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

type Message2 struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

var addr = []string{
	"192.168.1.216:9092",
	"192.168.1.217:9092",
	"192.168.1.218:9092",
}

func main() {
	// 配置
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 是否等待成功和失败后的响应，只有RequiredAcks设置不为NoReponse这里才有用
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	// 随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	// 使用配置，创建异步生产者
	producer, err := sarama.NewAsyncProducer(addr, config)
	if err != nil {
		log.Printf("异步生产者错误: %s \n", err.Error())
		return
	}
	defer producer.AsyncClose()

	// 循环判断哪个通道发送过来的数据
	go func(p sarama.AsyncProducer) {
		for {
			select {
			case ret := <-p.Successes():
				fmt.Println("offset: ", ret.Offset, "partitions: ", ret.Partition, "timestamp: ", ret.Timestamp.String())
			case ret := <-p.Errors():
				fmt.Println("error: ", ret.Error())
			}
		}
	}(producer)

	for {
		time.Sleep(time.Second * 2)

		// 创建消息
		value, _ := json.Marshal(Message2{
			Id:   strconv.FormatInt(time.Now().Unix(), 10),
			Text: time.Now().String(),
		})
		producerMsg := &sarama.ProducerMessage{
			Topic: "dev_topic_001",
			//Key:   sarama.StringEncoder("dev"),
			Value: sarama.ByteEncoder(value),
		}

		// 使用通道发送
		producer.Input() <- producerMsg
	}
}
