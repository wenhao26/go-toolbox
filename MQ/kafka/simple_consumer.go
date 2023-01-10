package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	// https://github.com/bsm/sarama-cluster

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	//config.Version = sarama.V2_1_1_0
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second * 5
	// 这个参数是控制消费组初始消费的位置，默认是OffsetNewest，即从最新的位置开始消费
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := sarama.NewClient([]string{
		"192.168.1.216:9092",
		"192.168.1.217:9092",
		"192.168.1.218:9092",
	}, config)
	if err != nil {
		log.Printf("消费者错误1: %s \n", err.Error())
		return
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Printf("消费者错误2: %s \n", err.Error())
		return
	}
	defer consumer.Close()

	topic := "dev_topic_001"
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("无法获取分区列表：", err)
		return
	}
	fmt.Println("分区列表：", partitionList)

	// 循环读取分区
	var wg sync.WaitGroup
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("无法启动分区[%d]的使用者：%s\n", partition, err)
			return
		}
		defer pc.Close()

		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
	}
	wg.Wait()
}
